FROM golang:1.24-bullseye AS builder

WORKDIR /app

# 依存関係ファイルをコピー
COPY go.mod go.sum ./

# 依存関係のダウンロード
RUN go mod download

# ソースコードをコピー
COPY . .

# アプリケーションのビルド
RUN CGO_ENABLED=0 GOOS=linux go build -o github-mcp-server-sse .

# 本番用の小さなイメージを使用
FROM debian:bullseye-slim

WORKDIR /app

# ビルドしたバイナリをコピー
COPY --from=builder /app/github-mcp-server-sse .

# SSE用のポートを公開
EXPOSE 8080

# SSEモードでサーバーを起動（環境変数PORTを使用）
ENTRYPOINT ["sh", "-c", "./github-mcp-server-sse -t sse -p 8080"] 
