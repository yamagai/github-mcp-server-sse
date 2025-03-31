# GitHub MCP Server

Model Context Protocol (MCP) を使用したGitHub統合サーバーです。このサーバーはLLM（大規模言語モデル）アプリケーションがGitHub APIと通信するための標準化されたインターフェースを提供します。

## 機能

このサーバーは以下の機能を提供します：

- リポジトリの検索
- リポジトリの作成
- ファイルの内容の取得
- ファイルの作成・更新
- 複数ファイルの一括プッシュ
- リポジトリのフォーク
- Pull Requestの作成
- Pull Requestの詳細取得
- Pull Requestへのレビュー追加

## インストール

```
go get github.com/yamagai/github-mcp-server-sse
```

## 使用方法

### 標準入出力モード (Stdio)

```bash
# 環境変数でGitHubトークンを設定
export GITHUB_TOKEN=your_github_token

# サーバーを起動（デフォルトは標準入出力モード）
github-mcp

# または明示的に標準入出力モードを指定
github-mcp -t stdio
```

### SSEモード

```bash
# 環境変数でGitHubトークンを設定（または各リクエストのAuthorizationヘッダーで指定）
export GITHUB_TOKEN=your_github_token

# デフォルトポート(8080)でSSEサーバーを起動
github-mcp -t sse

# カスタムポートでSSEサーバーを起動
github-mcp -t sse -p 3000
```

SSEモードでは、HTTPリクエストのAuthorizationヘッダーにGitHubトークンを含めることもできます：

```bash
curl -H "Authorization: Bearer YOUR_GITHUB_TOKEN" http://localhost:8080/events
```

## ツール一覧

| ツール名 | 説明 |
|---------|------|
| search_repositories | GitHubリポジトリを検索します |
| create_repository | 新しいGitHubリポジトリを作成します |
| get_file_contents | GitHubリポジトリからファイルの内容を取得します |
| create_or_update_file | GitHubリポジトリにファイルを作成または更新します |
| push_files | 複数のファイルを一度にGitHubリポジトリにプッシュします |
| fork_repository | GitHubリポジトリをフォークします |
| create_pull_request | GitHubリポジトリに新しいPull Requestを作成します |
| get_pull_request | GitHubリポジトリからPull Requestの詳細を取得します |
| create_pull_request_review | Pull Requestにレビューを作成します |

## 開発

```bash
# 依存関係の解決
go mod tidy

# ビルド
go build

# テスト
go test ./...
```

## コマンドラインオプション

| オプション | 短縮形 | 説明 | デフォルト値 |
|------------|--------|------|------------|
| --transport | -t | 使用するトランスポートタイプ (stdio または sse) | stdio |
| --port | -p | SSEサーバーのポート番号 | 8080 |

## 参考

このプロジェクトは[MCP-Go](https://github.com/mark3labs/mcp-go)ライブラリを使用しています。詳細なドキュメントについては、そちらを参照してください。

## ライセンス

MIT 
