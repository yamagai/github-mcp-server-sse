package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/yamagai/github-mcp-server-sse/common"
	"github.com/yamagai/github-mcp-server-sse/operations"
)

// GitHubMCPServer はGitHub MCP Serverのラッパー構造体
type GitHubMCPServer struct {
	server *server.MCPServer
}

// NewGitHubMCPServer は新しいGitHub MCP Serverを作成します
func NewGitHubMCPServer() *GitHubMCPServer {
	// MCPサーバーの作成
	s := server.NewMCPServer(
		"github-mcp-server",
		common.VERSION,
	)

	// リポジトリ検索ツール
	searchReposTool := mcp.NewTool("search_repositories",
		mcp.WithDescription("GitHub リポジトリを検索します"),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Description("検索クエリ"),
		),
		mcp.WithNumber("page",
			mcp.Description("ページ番号"),
		),
		mcp.WithNumber("per_page",
			mcp.Description("1ページあたりの結果数"),
		),
	)

	// リポジトリ作成ツール
	createRepoTool := mcp.NewTool("create_repository",
		mcp.WithDescription("新しいGitHubリポジトリを作成します"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("リポジトリ名"),
		),
		mcp.WithString("description",
			mcp.Description("リポジトリの説明"),
		),
		mcp.WithBoolean("private",
			mcp.Description("プライベートリポジトリかどうか"),
		),
		mcp.WithBoolean("auto_init",
			mcp.Description("READMEファイルを自動生成するかどうか"),
		),
	)

	// ファイル取得ツール
	getFileTool := mcp.NewTool("get_file_contents",
		mcp.WithDescription("GitHubリポジトリからファイルの内容を取得します"),
		mcp.WithString("owner",
			mcp.Required(),
			mcp.Description("リポジトリオーナー"),
		),
		mcp.WithString("repo",
			mcp.Required(),
			mcp.Description("リポジトリ名"),
		),
		mcp.WithString("path",
			mcp.Required(),
			mcp.Description("ファイルパス"),
		),
		mcp.WithString("branch",
			mcp.Description("ブランチ名 (省略時はデフォルトブランチ)"),
		),
	)

	// ファイル作成・更新ツール
	createOrUpdateFileTool := mcp.NewTool("create_or_update_file",
		mcp.WithDescription("GitHubリポジトリにファイルを作成または更新します"),
		mcp.WithString("owner",
			mcp.Required(),
			mcp.Description("リポジトリオーナー"),
		),
		mcp.WithString("repo",
			mcp.Required(),
			mcp.Description("リポジトリ名"),
		),
		mcp.WithString("path",
			mcp.Required(),
			mcp.Description("ファイルパス"),
		),
		mcp.WithString("content",
			mcp.Required(),
			mcp.Description("ファイルの内容"),
		),
		mcp.WithString("message",
			mcp.Required(),
			mcp.Description("コミットメッセージ"),
		),
		mcp.WithString("branch",
			mcp.Description("ブランチ名 (省略時はデフォルトブランチ)"),
		),
		mcp.WithString("sha",
			mcp.Description("更新する場合のファイルのSHA"),
		),
	)

	// 複数ファイルプッシュツール
	pushFilesTool := mcp.NewTool("push_files",
		mcp.WithDescription("複数のファイルを一度にGitHubリポジトリにプッシュします"),
		mcp.WithString("owner",
			mcp.Required(),
			mcp.Description("リポジトリオーナー"),
		),
		mcp.WithString("repo",
			mcp.Required(),
			mcp.Description("リポジトリ名"),
		),
		mcp.WithString("branch",
			mcp.Required(),
			mcp.Description("ブランチ名"),
		),
		mcp.WithArray("files",
			mcp.Required(),
			mcp.Description("ファイル操作の配列"),
		),
		mcp.WithString("message",
			mcp.Required(),
			mcp.Description("コミットメッセージ"),
		),
	)

	// リポジトリフォークツール
	forkRepoTool := mcp.NewTool("fork_repository",
		mcp.WithDescription("GitHubリポジトリをフォークします"),
		mcp.WithString("owner",
			mcp.Required(),
			mcp.Description("元のリポジトリオーナー"),
		),
		mcp.WithString("repo",
			mcp.Required(),
			mcp.Description("元のリポジトリ名"),
		),
		mcp.WithString("organization",
			mcp.Description("フォーク先の組織名 (省略時は個人アカウント)"),
		),
	)

	// Pull Request取得ツール
	getPRTool := mcp.NewTool("get_pull_request",
		mcp.WithDescription("GitHubリポジトリからPull Requestの詳細を取得します"),
		mcp.WithString("owner",
			mcp.Required(),
			mcp.Description("リポジトリオーナー"),
		),
		mcp.WithString("repo",
			mcp.Required(),
			mcp.Description("リポジトリ名"),
		),
		mcp.WithNumber("pull_number",
			mcp.Required(),
			mcp.Description("取得するPull Requestの番号"),
		),
	)

	// Pull Request作成ツール
	createPRTool := mcp.NewTool("create_pull_request",
		mcp.WithDescription("GitHubリポジトリに新しいPull Requestを作成します"),
		mcp.WithString("owner",
			mcp.Required(),
			mcp.Description("リポジトリオーナー"),
		),
		mcp.WithString("repo",
			mcp.Required(),
			mcp.Description("リポジトリ名"),
		),
		mcp.WithString("title",
			mcp.Required(),
			mcp.Description("Pull Requestのタイトル"),
		),
		mcp.WithString("body",
			mcp.Description("Pull Requestの説明"),
		),
		mcp.WithString("head",
			mcp.Required(),
			mcp.Description("変更を含むブランチ（例：'feature'）"),
		),
		mcp.WithString("base",
			mcp.Required(),
			mcp.Description("変更をマージするブランチ（例：'main'）"),
		),
		mcp.WithBoolean("draft",
			mcp.Description("ドラフトPull Requestとして作成するかどうか"),
		),
		mcp.WithBoolean("maintainer_can_modify",
			mcp.Description("メンテナーが変更を加えられるようにするかどうか"),
		),
	)

	// Pull Requestレビュー作成ツール
	createPRReviewTool := mcp.NewTool("create_pull_request_review",
		mcp.WithDescription("Pull Requestにレビューを作成します"),
		mcp.WithString("owner",
			mcp.Required(),
			mcp.Description("リポジトリオーナー"),
		),
		mcp.WithString("repo",
			mcp.Required(),
			mcp.Description("リポジトリ名"),
		),
		mcp.WithNumber("pull_number",
			mcp.Required(),
			mcp.Description("レビューするPull Requestの番号"),
		),
		mcp.WithString("body",
			mcp.Description("レビューのコメント本文"),
		),
		mcp.WithString("event",
			mcp.Required(),
			mcp.Description("レビューイベント (APPROVE, REQUEST_CHANGES, COMMENT)"),
		),
		mcp.WithString("commit_id",
			mcp.Description("レビューする特定のコミットID（省略時は最新コミット）"),
		),
	)

	// ツールハンドラーの登録
	s.AddTool(searchReposTool, handleSearchRepositories)
	s.AddTool(createRepoTool, handleCreateRepository)
	s.AddTool(getFileTool, handleGetFileContents)
	s.AddTool(createOrUpdateFileTool, handleCreateOrUpdateFile)
	s.AddTool(pushFilesTool, handlePushFiles)
	s.AddTool(forkRepoTool, handleForkRepository)
	s.AddTool(getPRTool, handleGetPullRequest)
	s.AddTool(createPRTool, handleCreatePullRequest)
	s.AddTool(createPRReviewTool, handleCreatePullRequestReview)

	return &GitHubMCPServer{
		server: s,
	}
}

// ServeSSE はSSEモードでサーバーを起動します
func (s *GitHubMCPServer) ServeSSE(addr string) *server.SSEServer {
	return server.NewSSEServer(s.server,
		server.WithBaseURL(fmt.Sprintf("http://%s", addr)),
		server.WithSSEContextFunc(common.AuthTokenFromRequest),
	)
}

// ServeStdio はStdioモードでサーバーを起動します
func (s *GitHubMCPServer) ServeStdio() error {
	return server.ServeStdio(s.server, server.WithStdioContextFunc(common.AuthTokenFromEnv))
}

func main() {
	// コマンドライン引数の処理
	var transport string
	flag.StringVar(&transport, "t", "stdio", "トランスポートタイプ (stdio または sse)")
	flag.StringVar(&transport, "transport", "stdio", "トランスポートタイプ (stdio または sse)")

	var port string
	flag.StringVar(&port, "p", "8080", "SSEサーバーのポート番号")
	flag.StringVar(&port, "port", "8080", "SSEサーバーのポート番号")

	flag.Parse()

	// GitHubMCPServerの作成
	s := NewGitHubMCPServer()

	// 指定されたトランスポートタイプでサーバーを起動
	switch transport {
	case "stdio":
		log.Printf("GitHub MCP Server を標準入出力モードで起動します")
		if err := s.ServeStdio(); err != nil {
			log.Fatalf("サーバーエラー: %v", err)
		}
	case "sse":
		addr := fmt.Sprintf("localhost:%s", port)
		log.Printf("GitHub MCP Server をSSEモードで起動します (アドレス: %s)", addr)
		sseServer := s.ServeSSE(addr)
		if err := sseServer.Start(addr); err != nil {
			log.Fatalf("サーバーエラー: %v", err)
		}
	default:
		log.Fatalf("無効なトランスポートタイプ: %s (stdio または sse を指定してください)", transport)
	}
}

// handleSearchRepositories はリポジトリ検索リクエストを処理します
func handleSearchRepositories(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// GitHubトークンの取得
	token, err := common.GetAuthTokenFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// パラメータの解析
	query, ok := request.Params.Arguments["query"].(string)
	if !ok {
		return nil, fmt.Errorf("query must be a string")
	}

	page := 1
	if p, ok := request.Params.Arguments["page"].(float64); ok {
		page = int(p)
	}

	perPage := 30
	if pp, ok := request.Params.Arguments["per_page"].(float64); ok {
		perPage = int(pp)
	}

	// リポジトリ検索の実行
	result, err := operations.SearchRepositories(operations.SearchRepositoriesOptions{
		Query:   query,
		Page:    page,
		PerPage: perPage,
	}, token)
	if err != nil {
		return nil, err
	}

	// JSON形式で結果を返す
	jsonResult, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(string(jsonResult)), nil
}

// handleCreateRepository はリポジトリ作成リクエストを処理します
func handleCreateRepository(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// GitHubトークンの取得
	token, err := common.GetAuthTokenFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// パラメータの解析
	name, ok := request.Params.Arguments["name"].(string)
	if !ok {
		return nil, fmt.Errorf("name must be a string")
	}

	description := ""
	if desc, ok := request.Params.Arguments["description"].(string); ok {
		description = desc
	}

	private := false
	if priv, ok := request.Params.Arguments["private"].(bool); ok {
		private = priv
	}

	autoInit := false
	if ai, ok := request.Params.Arguments["auto_init"].(bool); ok {
		autoInit = ai
	}

	// リポジトリ作成の実行
	result, err := operations.CreateRepository(operations.CreateRepositoryOptions{
		Name:        name,
		Description: description,
		Private:     private,
		AutoInit:    autoInit,
	}, token)
	if err != nil {
		return nil, err
	}

	// JSON形式で結果を返す
	jsonResult, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(string(jsonResult)), nil
}

// handleGetFileContents はファイル取得リクエストを処理します
func handleGetFileContents(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// GitHubトークンの取得
	token, err := common.GetAuthTokenFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// パラメータの解析
	owner, ok := request.Params.Arguments["owner"].(string)
	if !ok {
		return nil, fmt.Errorf("owner must be a string")
	}

	repo, ok := request.Params.Arguments["repo"].(string)
	if !ok {
		return nil, fmt.Errorf("repo must be a string")
	}

	path, ok := request.Params.Arguments["path"].(string)
	if !ok {
		return nil, fmt.Errorf("path must be a string")
	}

	branch := ""
	if b, ok := request.Params.Arguments["branch"].(string); ok {
		branch = b
	}

	// ファイル取得の実行
	result, err := operations.GetFileContents(operations.GetFileContentOptions{
		Owner:  owner,
		Repo:   repo,
		Path:   path,
		Branch: branch,
	}, token)
	if err != nil {
		return nil, err
	}

	// JSON形式で結果を返す
	jsonResult, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(string(jsonResult)), nil
}

// handleCreateOrUpdateFile はファイル作成・更新リクエストを処理します
func handleCreateOrUpdateFile(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// GitHubトークンの取得
	token, err := common.GetAuthTokenFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// パラメータの解析
	owner, ok := request.Params.Arguments["owner"].(string)
	if !ok {
		return nil, fmt.Errorf("owner must be a string")
	}

	repo, ok := request.Params.Arguments["repo"].(string)
	if !ok {
		return nil, fmt.Errorf("repo must be a string")
	}

	path, ok := request.Params.Arguments["path"].(string)
	if !ok {
		return nil, fmt.Errorf("path must be a string")
	}

	content, ok := request.Params.Arguments["content"].(string)
	if !ok {
		return nil, fmt.Errorf("content must be a string")
	}

	message, ok := request.Params.Arguments["message"].(string)
	if !ok {
		return nil, fmt.Errorf("message must be a string")
	}

	branch := ""
	if b, ok := request.Params.Arguments["branch"].(string); ok {
		branch = b
	}

	sha := ""
	if s, ok := request.Params.Arguments["sha"].(string); ok {
		sha = s
	}

	// ファイル作成・更新の実行
	result, err := operations.CreateOrUpdateFile(operations.CreateOrUpdateFileOptions{
		Owner:   owner,
		Repo:    repo,
		Path:    path,
		Content: content,
		Message: message,
		Branch:  branch,
		SHA:     sha,
	}, token)
	if err != nil {
		return nil, err
	}

	// JSON形式で結果を返す
	jsonResult, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(string(jsonResult)), nil
}

// handlePushFiles は複数ファイルプッシュリクエストを処理します
func handlePushFiles(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// GitHubトークンの取得
	token, err := common.GetAuthTokenFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// パラメータの解析
	owner, ok := request.Params.Arguments["owner"].(string)
	if !ok {
		return nil, fmt.Errorf("owner must be a string")
	}

	repo, ok := request.Params.Arguments["repo"].(string)
	if !ok {
		return nil, fmt.Errorf("repo must be a string")
	}

	branch, ok := request.Params.Arguments["branch"].(string)
	if !ok {
		return nil, fmt.Errorf("branch must be a string")
	}

	filesRaw, ok := request.Params.Arguments["files"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("files must be an array")
	}

	message, ok := request.Params.Arguments["message"].(string)
	if !ok {
		return nil, fmt.Errorf("message must be a string")
	}

	// ファイル操作の変換
	files := make([]operations.FileOperation, 0, len(filesRaw))
	for _, f := range filesRaw {
		fileMap, ok := f.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("each file must be an object")
		}

		path, ok := fileMap["path"].(string)
		if !ok {
			return nil, fmt.Errorf("file path must be a string")
		}

		content, ok := fileMap["content"].(string)
		if !ok {
			return nil, fmt.Errorf("file content must be a string")
		}

		sha := ""
		if s, ok := fileMap["sha"].(string); ok {
			sha = s
		}

		files = append(files, operations.FileOperation{
			Path:    path,
			Content: content,
			SHA:     sha,
		})
	}

	// ファイル更新の実行
	result, err := operations.PushFiles(operations.PushFilesOptions{
		Owner:   owner,
		Repo:    repo,
		Branch:  branch,
		Files:   files,
		Message: message,
	}, token)
	if err != nil {
		return nil, err
	}

	// JSON形式で結果を返す
	jsonResult, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(string(jsonResult)), nil
}

// handleForkRepository はリポジトリフォークリクエストを処理します
func handleForkRepository(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// GitHubトークンの取得
	token, err := common.GetAuthTokenFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// パラメータの解析
	owner, ok := request.Params.Arguments["owner"].(string)
	if !ok {
		return nil, fmt.Errorf("owner must be a string")
	}

	repo, ok := request.Params.Arguments["repo"].(string)
	if !ok {
		return nil, fmt.Errorf("repo must be a string")
	}

	organization := ""
	if org, ok := request.Params.Arguments["organization"].(string); ok {
		organization = org
	}

	// リポジトリフォークの実行
	result, err := operations.ForkRepository(operations.ForkRepositoryOptions{
		Owner:        owner,
		Repo:         repo,
		Organization: organization,
	}, token)
	if err != nil {
		return nil, err
	}

	// JSON形式で結果を返す
	jsonResult, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(string(jsonResult)), nil
}

// handleGetPullRequest はPull Request取得リクエストを処理します
func handleGetPullRequest(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// GitHubトークンの取得
	token, err := common.GetAuthTokenFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// パラメータの解析
	owner, ok := request.Params.Arguments["owner"].(string)
	if !ok {
		return nil, fmt.Errorf("owner must be a string")
	}

	repo, ok := request.Params.Arguments["repo"].(string)
	if !ok {
		return nil, fmt.Errorf("repo must be a string")
	}

	pullNumberFloat, ok := request.Params.Arguments["pull_number"].(float64)
	if !ok {
		return nil, fmt.Errorf("pull_number must be a number")
	}
	pullNumber := int(pullNumberFloat)

	// Pull Request取得の実行
	result, err := operations.GetPullRequest(operations.GetPullRequestOptions{
		Owner:      owner,
		Repo:       repo,
		PullNumber: pullNumber,
	}, token)
	if err != nil {
		return nil, err
	}

	// JSON形式で結果を返す
	jsonResult, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(string(jsonResult)), nil
}

// handleCreatePullRequest はPull Request作成リクエストを処理します
func handleCreatePullRequest(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// GitHubトークンの取得
	token, err := common.GetAuthTokenFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// パラメータの解析
	owner, ok := request.Params.Arguments["owner"].(string)
	if !ok {
		return nil, fmt.Errorf("owner must be a string")
	}

	repo, ok := request.Params.Arguments["repo"].(string)
	if !ok {
		return nil, fmt.Errorf("repo must be a string")
	}

	title, ok := request.Params.Arguments["title"].(string)
	if !ok {
		return nil, fmt.Errorf("title must be a string")
	}

	head, ok := request.Params.Arguments["head"].(string)
	if !ok {
		return nil, fmt.Errorf("head must be a string")
	}

	base, ok := request.Params.Arguments["base"].(string)
	if !ok {
		return nil, fmt.Errorf("base must be a string")
	}

	body := ""
	if b, ok := request.Params.Arguments["body"].(string); ok {
		body = b
	}

	draft := false
	if d, ok := request.Params.Arguments["draft"].(bool); ok {
		draft = d
	}

	maintainerCanModify := false
	if m, ok := request.Params.Arguments["maintainer_can_modify"].(bool); ok {
		maintainerCanModify = m
	}

	// Pull Request作成の実行
	result, err := operations.CreatePullRequest(operations.CreatePullRequestOptions{
		Owner:               owner,
		Repo:                repo,
		Title:               title,
		Body:                body,
		Head:                head,
		Base:                base,
		Draft:               draft,
		MaintainerCanModify: maintainerCanModify,
	}, token)
	if err != nil {
		return nil, err
	}

	// JSON形式で結果を返す
	jsonResult, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(string(jsonResult)), nil
}

// handleCreatePullRequestReview はPull Requestレビュー作成リクエストを処理します
func handleCreatePullRequestReview(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// GitHubトークンの取得
	token, err := common.GetAuthTokenFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// パラメータの解析
	owner, ok := request.Params.Arguments["owner"].(string)
	if !ok {
		return nil, fmt.Errorf("owner must be a string")
	}

	repo, ok := request.Params.Arguments["repo"].(string)
	if !ok {
		return nil, fmt.Errorf("repo must be a string")
	}

	pullNumberFloat, ok := request.Params.Arguments["pull_number"].(float64)
	if !ok {
		return nil, fmt.Errorf("pull_number must be a number")
	}
	pullNumber := int(pullNumberFloat)

	event, ok := request.Params.Arguments["event"].(string)
	if !ok {
		return nil, fmt.Errorf("event must be a string")
	}

	body := ""
	if b, ok := request.Params.Arguments["body"].(string); ok {
		body = b
	}

	commitID := ""
	if c, ok := request.Params.Arguments["commit_id"].(string); ok {
		commitID = c
	}

	// Pull Requestレビュー作成の実行
	result, err := operations.CreatePullRequestReview(operations.PullRequestReviewOptions{
		Owner:      owner,
		Repo:       repo,
		PullNumber: pullNumber,
		Body:       body,
		Event:      event,
		CommitID:   commitID,
	}, token)
	if err != nil {
		return nil, err
	}

	// JSON形式で結果を返す
	jsonResult, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(string(jsonResult)), nil
}
