package operations

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-github/v70/github"
)

// FileContent はGitHubファイルの内容を表します
type FileContent struct {
	Type        string `json:"type"`
	Encoding    string `json:"encoding,omitempty"`
	Size        int    `json:"size"`
	Name        string `json:"name"`
	Path        string `json:"path"`
	Content     string `json:"content,omitempty"`
	SHA         string `json:"sha"`
	URL         string `json:"url"`
	HTMLURL     string `json:"html_url"`
	DownloadURL string `json:"download_url,omitempty"`
}

// FileOperation はファイル操作のタイプを表します
type FileOperation struct {
	Path    string `json:"path"`
	Content string `json:"content"`
	SHA     string `json:"sha,omitempty"`
}

// GetFileContentOptions はファイル取得オプションを表します
type GetFileContentOptions struct {
	Owner  string `json:"owner"`
	Repo   string `json:"repo"`
	Path   string `json:"path"`
	Branch string `json:"branch,omitempty"`
}

// CreateOrUpdateFileOptions はファイル作成・更新オプションを表します
type CreateOrUpdateFileOptions struct {
	Owner   string `json:"owner"`
	Repo    string `json:"repo"`
	Path    string `json:"path"`
	Content string `json:"content"`
	Message string `json:"message"`
	Branch  string `json:"branch,omitempty"`
	SHA     string `json:"sha,omitempty"`
}

// PushFilesOptions は複数ファイル更新オプションを表します
type PushFilesOptions struct {
	Owner   string          `json:"owner"`
	Repo    string          `json:"repo"`
	Branch  string          `json:"branch"`
	Files   []FileOperation `json:"files"`
	Message string          `json:"message"`
}

// CommitResult はコミット結果を表します
type CommitResult struct {
	SHA    string `json:"sha"`
	URL    string `json:"url"`
	Author struct {
		Name  string    `json:"name"`
		Email string    `json:"email"`
		Date  time.Time `json:"date"`
	} `json:"author"`
	Message string `json:"message"`
}

// GetFileContents はファイルの内容を取得します
func GetFileContents(options GetFileContentOptions, token string) (*FileContent, error) {
	ctx := context.Background()
	client := getGitHubClient(ctx, token)

	// ファイル取得オプションの設定
	opts := &github.RepositoryContentGetOptions{}
	if options.Branch != "" {
		opts.Ref = options.Branch
	}

	// GitHub APIを呼び出してファイル内容を取得
	fileContent, _, _, err := client.Repositories.GetContents(
		ctx,
		options.Owner,
		options.Repo,
		options.Path,
		opts,
	)
	if err != nil {
		return nil, err
	}

	// content, err := fileContent.GetContent()
	content, decodeErr := fileContent.GetContent()
	if decodeErr != nil {
		return nil, fmt.Errorf("ファイル内容のデコードに失敗: %v", decodeErr)
	}

	// 結果をマッピング
	return &FileContent{
		Type:        fileContent.GetType(),
		Encoding:    fileContent.GetEncoding(),
		Size:        fileContent.GetSize(),
		Name:        fileContent.GetName(),
		Path:        fileContent.GetPath(),
		Content:     content,
		SHA:         fileContent.GetSHA(),
		URL:         fileContent.GetURL(),
		HTMLURL:     fileContent.GetHTMLURL(),
		DownloadURL: fileContent.GetDownloadURL(),
	}, nil
}

// CreateOrUpdateFile はファイルを作成または更新します
func CreateOrUpdateFile(options CreateOrUpdateFileOptions, token string) (*CommitResult, error) {
	ctx := context.Background()
	client := getGitHubClient(ctx, token)

	// ファイル作成・更新リクエストの設定
	opts := &github.RepositoryContentFileOptions{
		Message: github.Ptr(options.Message),
		Content: []byte(options.Content),
	}

	if options.Branch != "" {
		opts.Branch = github.Ptr(options.Branch)
	}

	if options.SHA != "" {
		opts.SHA = github.Ptr(options.SHA)
	}

	// GitHub APIを呼び出してファイルを作成または更新
	commit, _, err := client.Repositories.CreateFile(
		ctx,
		options.Owner,
		options.Repo,
		options.Path,
		opts,
	)
	if err != nil {
		return nil, err
	}

	// 結果をマッピング
	result := &CommitResult{
		SHA:     commit.GetSHA(),
		URL:     commit.GetURL(),
		Message: options.Message,
	}

	// 著者情報があれば設定
	if commit.GetAuthor() != nil {
		result.Author.Name = commit.Author.GetName()
		result.Author.Email = commit.Author.GetEmail()
		if commit.Author.Date != nil {
			result.Author.Date = commit.Author.Date.Time
		}
	}

	return result, nil
}

// PushFiles は複数のファイルを一度にプッシュします
func PushFiles(options PushFilesOptions, token string) (*CommitResult, error) {
	// 注意: GitHub APIは一度に複数ファイルを更新する直接的なエンドポイントを提供していません
	// そのため、このような実装が必要になります
	ctx := context.Background()
	client := getGitHubClient(ctx, token)

	// 現在のブランチの最新コミットSHAを取得
	ref, _, err := client.Git.GetRef(ctx, options.Owner, options.Repo, "refs/heads/"+options.Branch)
	if err != nil {
		return nil, fmt.Errorf("ブランチの取得に失敗: %v", err)
	}
	baseTreeSHA := ref.Object.GetSHA()

	// ベースとなるツリーを取得
	baseCommit, _, err := client.Git.GetCommit(ctx, options.Owner, options.Repo, baseTreeSHA)
	if err != nil {
		return nil, fmt.Errorf("コミットの取得に失敗: %v", err)
	}
	baseTreeSHA = baseCommit.Tree.GetSHA()

	// 新しいツリーのエントリを作成
	entries := make([]*github.TreeEntry, 0, len(options.Files))
	for _, file := range options.Files {
		// ファイル内容をBase64からデコード (必要な場合)
		content := file.Content
		mode := "100644" // 通常ファイルのモード
		entryType := "blob"

		entry := &github.TreeEntry{
			Path:    github.String(file.Path),
			Mode:    github.String(mode),
			Type:    github.String(entryType),
			Content: github.String(content),
		}
		entries = append(entries, entry)
	}

	// 新しいツリーを作成
	newTree, _, err := client.Git.CreateTree(ctx, options.Owner, options.Repo, baseTreeSHA, entries)
	if err != nil {
		return nil, fmt.Errorf("ツリーの作成に失敗: %v", err)
	}

	// 新しいコミットを作成
	newCommit, _, err := client.Git.CreateCommit(ctx, options.Owner, options.Repo, &github.Commit{
		Message: github.String(options.Message),
		Tree:    newTree,
		Parents: []*github.Commit{{SHA: github.String(baseTreeSHA)}},
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("コミットの作成に失敗: %v", err)
	}

	// リファレンスを更新
	_, _, err = client.Git.UpdateRef(ctx, options.Owner, options.Repo, &github.Reference{
		Ref: github.String("refs/heads/" + options.Branch),
		Object: &github.GitObject{
			SHA: newCommit.SHA,
		},
	}, false)
	if err != nil {
		return nil, fmt.Errorf("リファレンスの更新に失敗: %v", err)
	}

	// 結果をマッピング
	result := &CommitResult{
		SHA:     newCommit.GetSHA(),
		URL:     newCommit.GetURL(),
		Message: options.Message,
	}

	// 著者情報があれば設定
	if newCommit.GetAuthor() != nil {
		result.Author.Name = newCommit.Author.GetName()
		result.Author.Email = newCommit.Author.GetEmail()
		if newCommit.Author.Date != nil {
			result.Author.Date = newCommit.Author.Date.Time
		}
	}

	return result, nil
}
