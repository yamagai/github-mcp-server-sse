package common

import (
	"fmt"
	"time"
)

// GitHubError は GitHub API からのエラーを表します
type GitHubError struct {
	Message string
	Status  int
}

func (e *GitHubError) Error() string {
	return fmt.Sprintf("GitHub API Error: %s (Status: %d)", e.Message, e.Status)
}

// GitHubValidationError はバリデーションエラーを表します
type GitHubValidationError struct {
	GitHubError
	Response interface{}
}

// GitHubResourceNotFoundError はリソースが見つからない場合のエラーを表します
type GitHubResourceNotFoundError struct {
	GitHubError
}

// GitHubAuthenticationError は認証エラーを表します
type GitHubAuthenticationError struct {
	GitHubError
}

// GitHubPermissionError は権限エラーを表します
type GitHubPermissionError struct {
	GitHubError
}

// GitHubRateLimitError はレート制限エラーを表します
type GitHubRateLimitError struct {
	GitHubError
	ResetAt time.Time
}

// GitHubConflictError は競合エラーを表します
type GitHubConflictError struct {
	GitHubError
}

// IsGitHubError は指定されたエラーがGitHubエラーかどうかを判断します
func IsGitHubError(err error) bool {
	switch err.(type) {
	case *GitHubError, *GitHubValidationError, *GitHubResourceNotFoundError,
		*GitHubAuthenticationError, *GitHubPermissionError, *GitHubRateLimitError,
		*GitHubConflictError:
		return true
	default:
		return false
	}
}

// FormatGitHubError はGitHubエラーを人間が読みやすい形式にフォーマットします
func FormatGitHubError(err error) string {
	if !IsGitHubError(err) {
		return err.Error()
	}

	message := "GitHub API Error"

	switch e := err.(type) {
	case *GitHubValidationError:
		message = fmt.Sprintf("Validation Error: %s", e.Message)
		if e.Response != nil {
			message += fmt.Sprintf("\nDetails: %v", e.Response)
		}
	case *GitHubResourceNotFoundError:
		message = fmt.Sprintf("Not Found: %s", e.Message)
	case *GitHubAuthenticationError:
		message = fmt.Sprintf("Authentication Failed: %s", e.Message)
	case *GitHubPermissionError:
		message = fmt.Sprintf("Permission Denied: %s", e.Message)
	case *GitHubRateLimitError:
		message = fmt.Sprintf("Rate Limit Exceeded: %s\nResets at: %s", e.Message, e.ResetAt.Format(time.RFC3339))
	case *GitHubConflictError:
		message = fmt.Sprintf("Conflict: %s", e.Message)
	case *GitHubError:
		message = fmt.Sprintf("GitHub API Error: %s", e.Message)
	}

	return message
}
