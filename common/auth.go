package common

import (
	"context"
	"fmt"
	"net/http"
	"os"
)

// authKey は認証トークンを保存するためのコンテキストキー
type authKey struct{}

// WithAuthToken はコンテキストに認証トークンを追加します
func WithAuthToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, authKey{}, token)
}

// AuthTokenFromRequest はHTTPリクエストから認証トークンを抽出してコンテキストに追加します
func AuthTokenFromRequest(ctx context.Context, r *http.Request) context.Context {
	return WithAuthToken(ctx, r.Header.Get("Authorization"))
}

// AuthTokenFromEnv は環境変数から認証トークンを抽出してコンテキストに追加します
func AuthTokenFromEnv(ctx context.Context) context.Context {
	return WithAuthToken(ctx, os.Getenv("GITHUB_TOKEN"))
}

// GetAuthTokenFromContext はコンテキストから認証トークンを取得します
func GetAuthTokenFromContext(ctx context.Context) (string, error) {
	token, ok := ctx.Value(authKey{}).(string)
	if !ok || token == "" {
		return "", fmt.Errorf("認証トークンがありません。環境変数GITHUB_TOKENを設定するか、Authorizationヘッダーを指定してください")
	}
	return token, nil
}
