package operations

import (
	"context"

	"github.com/google/go-github/v70/github"
	"golang.org/x/oauth2"
)

// Repository はGitHubリポジトリを表します
type Repository struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	FullName    string `json:"full_name"`
	Description string `json:"description"`
	Private     bool   `json:"private"`
	HTMLURL     string `json:"html_url"`
	CloneURL    string `json:"clone_url"`
	SSHURL      string `json:"ssh_url"`
	Fork        bool   `json:"fork"`
}

// SearchRepositoriesOptions は検索オプションを表します
type SearchRepositoriesOptions struct {
	Query   string `json:"query"`
	Page    int    `json:"page,omitempty"`
	PerPage int    `json:"per_page,omitempty"`
}

// SearchRepositoriesResult は検索結果を表します
type SearchRepositoriesResult struct {
	TotalCount int          `json:"total_count"`
	Items      []Repository `json:"items"`
}

// CreateRepositoryOptions はリポジトリ作成オプションを表します
type CreateRepositoryOptions struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Private     bool   `json:"private,omitempty"`
	AutoInit    bool   `json:"auto_init,omitempty"`
}

// ForkRepositoryOptions はフォークオプションを表します
type ForkRepositoryOptions struct {
	Owner        string `json:"owner"`
	Repo         string `json:"repo"`
	Organization string `json:"organization,omitempty"`
}

// getGitHubClient は認証済みのGitHubクライアントを作成します
func getGitHubClient(ctx context.Context, token string) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}

// SearchRepositories はGitHubリポジトリを検索します
func SearchRepositories(options SearchRepositoriesOptions, token string) (*SearchRepositoriesResult, error) {
	ctx := context.Background()
	client := getGitHubClient(ctx, token)

	// 検索オプションの設定
	opts := &github.SearchOptions{
		ListOptions: github.ListOptions{
			Page:    options.Page,
			PerPage: options.PerPage,
		},
	}

	// GitHub APIを呼び出してリポジトリを検索
	repos, _, err := client.Search.Repositories(ctx, options.Query, opts)
	if err != nil {
		return nil, err
	}

	// 結果をマッピング
	var items []Repository
	for _, repo := range repos.Repositories {
		items = append(items, Repository{
			ID:          int(repo.GetID()),
			Name:        repo.GetName(),
			FullName:    repo.GetFullName(),
			Description: repo.GetDescription(),
			Private:     repo.GetPrivate(),
			HTMLURL:     repo.GetHTMLURL(),
			CloneURL:    repo.GetCloneURL(),
			SSHURL:      repo.GetSSHURL(),
			Fork:        repo.GetFork(),
		})
	}

	return &SearchRepositoriesResult{
		TotalCount: int(repos.GetTotal()),
		Items:      items,
	}, nil
}

// CreateRepository は新しいリポジトリを作成します
func CreateRepository(options CreateRepositoryOptions, token string) (*Repository, error) {
	ctx := context.Background()
	client := getGitHubClient(ctx, token)

	// リポジトリ作成リクエストの設定
	repo := &github.Repository{
		Name:        github.String(options.Name),
		Description: github.String(options.Description),
		Private:     github.Bool(options.Private),
		AutoInit:    github.Bool(options.AutoInit),
	}

	// GitHub APIを呼び出してリポジトリを作成
	newRepo, _, err := client.Repositories.Create(ctx, "", repo)
	if err != nil {
		return nil, err
	}

	// 結果をマッピング
	return &Repository{
		ID:          int(newRepo.GetID()),
		Name:        newRepo.GetName(),
		FullName:    newRepo.GetFullName(),
		Description: newRepo.GetDescription(),
		Private:     newRepo.GetPrivate(),
		HTMLURL:     newRepo.GetHTMLURL(),
		CloneURL:    newRepo.GetCloneURL(),
		SSHURL:      newRepo.GetSSHURL(),
		Fork:        newRepo.GetFork(),
	}, nil
}

// ForkRepository はリポジトリをフォークします
func ForkRepository(options ForkRepositoryOptions, token string) (*Repository, error) {
	ctx := context.Background()
	client := getGitHubClient(ctx, token)

	// フォークオプションの設定
	forkOpts := &github.RepositoryCreateForkOptions{}
	if options.Organization != "" {
		forkOpts.Organization = options.Organization
	}

	// GitHub APIを呼び出してリポジトリをフォーク
	newRepo, _, err := client.Repositories.CreateFork(ctx, options.Owner, options.Repo, forkOpts)
	if err != nil {
		return nil, err
	}

	// 結果をマッピング
	return &Repository{
		ID:          int(newRepo.GetID()),
		Name:        newRepo.GetName(),
		FullName:    newRepo.GetFullName(),
		Description: newRepo.GetDescription(),
		Private:     newRepo.GetPrivate(),
		HTMLURL:     newRepo.GetHTMLURL(),
		CloneURL:    newRepo.GetCloneURL(),
		SSHURL:      newRepo.GetSSHURL(),
		Fork:        newRepo.GetFork(),
	}, nil
}
