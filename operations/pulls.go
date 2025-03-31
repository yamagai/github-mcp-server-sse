package operations

import (
	"context"
	"time"

	"github.com/google/go-github/v70/github"
)

// PullRequest はGitHub Pull Requestを表します
type PullRequest struct {
	ID                  int       `json:"id"`
	Number              int       `json:"number"`
	State               string    `json:"state"`
	Title               string    `json:"title"`
	Body                string    `json:"body"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
	ClosedAt            time.Time `json:"closed_at,omitempty"`
	MergedAt            time.Time `json:"merged_at,omitempty"`
	MergeCommitSHA      string    `json:"merge_commit_sha,omitempty"`
	User                User      `json:"user"`
	HTMLURL             string    `json:"html_url"`
	DiffURL             string    `json:"diff_url"`
	PatchURL            string    `json:"patch_url"`
	Base                Ref       `json:"base"`
	Head                Ref       `json:"head"`
	Merged              bool      `json:"merged"`
	Mergeable           bool      `json:"mergeable"`
	MergeableState      string    `json:"mergeable_state"`
	Comments            int       `json:"comments"`
	ReviewComments      int       `json:"review_comments"`
	Commits             int       `json:"commits"`
	Additions           int       `json:"additions"`
	Deletions           int       `json:"deletions"`
	ChangedFiles        int       `json:"changed_files"`
	Draft               bool      `json:"draft"`
	RequestedReviewers  []User    `json:"requested_reviewers"`
	MaintainerCanModify bool      `json:"maintainer_can_modify"`
}

// User はGitHubユーザーを表します
type User struct {
	Login     string `json:"login"`
	ID        int    `json:"id"`
	AvatarURL string `json:"avatar_url"`
	HTMLURL   string `json:"html_url"`
	Type      string `json:"type"`
}

// Ref はブランチの参照を表します
type Ref struct {
	Label string     `json:"label"`
	Ref   string     `json:"ref"`
	SHA   string     `json:"sha"`
	User  User       `json:"user"`
	Repo  Repository `json:"repo"`
}

// CreatePullRequestOptions はPull Request作成オプションを表します
type CreatePullRequestOptions struct {
	Owner               string `json:"owner"`
	Repo                string `json:"repo"`
	Title               string `json:"title"`
	Body                string `json:"body,omitempty"`
	Head                string `json:"head"`
	Base                string `json:"base"`
	Draft               bool   `json:"draft,omitempty"`
	MaintainerCanModify bool   `json:"maintainer_can_modify,omitempty"`
}

// PullRequestReviewOptions はPull Requestレビューオプションを表します
type PullRequestReviewOptions struct {
	Owner      string `json:"owner"`
	Repo       string `json:"repo"`
	PullNumber int    `json:"pull_number"`
	Body       string `json:"body,omitempty"`
	Event      string `json:"event"` // APPROVE, REQUEST_CHANGES, COMMENT
	CommitID   string `json:"commit_id,omitempty"`
}

// PullRequestReview はレビュー結果を表します
type PullRequestReview struct {
	ID          int       `json:"id"`
	User        User      `json:"user"`
	Body        string    `json:"body"`
	State       string    `json:"state"`
	HTMLURL     string    `json:"html_url"`
	CommitID    string    `json:"commit_id"`
	SubmittedAt time.Time `json:"submitted_at"`
}

// GetPullRequestOptions はPull Request取得オプションを表します
type GetPullRequestOptions struct {
	Owner      string `json:"owner"`
	Repo       string `json:"repo"`
	PullNumber int    `json:"pull_number"`
}

// mapGitHubUserToUser はGitHubユーザーをUserモデルに変換します
func mapGitHubUserToUser(ghUser *github.User) User {
	if ghUser == nil {
		return User{}
	}
	return User{
		Login:     ghUser.GetLogin(),
		ID:        int(ghUser.GetID()),
		AvatarURL: ghUser.GetAvatarURL(),
		HTMLURL:   ghUser.GetHTMLURL(),
		Type:      ghUser.GetType(),
	}
}

// mapGitHubRepositoryToRepo はGitHubリポジトリをRepositoryモデルに変換します
func mapGitHubRepositoryToRepo(ghRepo *github.Repository) Repository {
	if ghRepo == nil {
		return Repository{}
	}
	return Repository{
		ID:          int(ghRepo.GetID()),
		Name:        ghRepo.GetName(),
		FullName:    ghRepo.GetFullName(),
		Description: ghRepo.GetDescription(),
		Private:     ghRepo.GetPrivate(),
		HTMLURL:     ghRepo.GetHTMLURL(),
		CloneURL:    ghRepo.GetCloneURL(),
		SSHURL:      ghRepo.GetSSHURL(),
		Fork:        ghRepo.GetFork(),
	}
}

// mapTimestamp はgithub.Timestampをtime.Timeに変換します
func mapTimestamp(timestamp *github.Timestamp) time.Time {
	if timestamp == nil {
		return time.Time{}
	}
	return timestamp.Time
}

// CreatePullRequest は新しいPull Requestを作成します
func CreatePullRequest(options CreatePullRequestOptions, token string) (*PullRequest, error) {
	ctx := context.Background()
	client := getGitHubClient(ctx, token)

	// Pull Request作成リクエストの設定
	newPR := &github.NewPullRequest{
		Title:               github.String(options.Title),
		Head:                github.String(options.Head),
		Base:                github.String(options.Base),
		Body:                github.String(options.Body),
		MaintainerCanModify: github.Bool(options.MaintainerCanModify),
		Draft:               github.Bool(options.Draft),
	}

	// GitHub APIを呼び出してPull Requestを作成
	pr, _, err := client.PullRequests.Create(ctx, options.Owner, options.Repo, newPR)
	if err != nil {
		return nil, err
	}

	// 結果をマッピング
	result := &PullRequest{
		ID:                  int(pr.GetID()),
		Number:              pr.GetNumber(),
		State:               pr.GetState(),
		Title:               pr.GetTitle(),
		Body:                pr.GetBody(),
		CreatedAt:           mapTimestamp(pr.CreatedAt),
		UpdatedAt:           mapTimestamp(pr.UpdatedAt),
		ClosedAt:            mapTimestamp(pr.ClosedAt),
		MergedAt:            mapTimestamp(pr.MergedAt),
		MergeCommitSHA:      pr.GetMergeCommitSHA(),
		User:                mapGitHubUserToUser(pr.User),
		HTMLURL:             pr.GetHTMLURL(),
		DiffURL:             pr.GetDiffURL(),
		PatchURL:            pr.GetPatchURL(),
		Merged:              pr.GetMerged(),
		Mergeable:           pr.GetMergeable(),
		MergeableState:      pr.GetMergeableState(),
		Comments:            pr.GetComments(),
		ReviewComments:      pr.GetReviewComments(),
		Commits:             pr.GetCommits(),
		Additions:           pr.GetAdditions(),
		Deletions:           pr.GetDeletions(),
		ChangedFiles:        pr.GetChangedFiles(),
		Draft:               pr.GetDraft(),
		MaintainerCanModify: pr.GetMaintainerCanModify(),
	}

	// ベースとヘッドの設定
	if pr.Base != nil {
		result.Base = Ref{
			Label: pr.Base.GetLabel(),
			Ref:   pr.Base.GetRef(),
			SHA:   pr.Base.GetSHA(),
			User:  mapGitHubUserToUser(pr.Base.User),
			Repo:  mapGitHubRepositoryToRepo(pr.Base.Repo),
		}
	}

	if pr.Head != nil {
		result.Head = Ref{
			Label: pr.Head.GetLabel(),
			Ref:   pr.Head.GetRef(),
			SHA:   pr.Head.GetSHA(),
			User:  mapGitHubUserToUser(pr.Head.User),
			Repo:  mapGitHubRepositoryToRepo(pr.Head.Repo),
		}
	}

	// レビューア情報があれば設定
	if pr.RequestedReviewers != nil {
		result.RequestedReviewers = make([]User, 0, len(pr.RequestedReviewers))
		for _, reviewer := range pr.RequestedReviewers {
			result.RequestedReviewers = append(result.RequestedReviewers, mapGitHubUserToUser(reviewer))
		}
	}

	return result, nil
}

// GetPullRequest はPull Requestの詳細を取得します
func GetPullRequest(options GetPullRequestOptions, token string) (*PullRequest, error) {
	ctx := context.Background()
	client := getGitHubClient(ctx, token)

	// GitHub APIを呼び出してPull Requestを取得
	pr, _, err := client.PullRequests.Get(ctx, options.Owner, options.Repo, options.PullNumber)
	if err != nil {
		return nil, err
	}

	// 結果をマッピング
	result := &PullRequest{
		ID:                  int(pr.GetID()),
		Number:              pr.GetNumber(),
		State:               pr.GetState(),
		Title:               pr.GetTitle(),
		Body:                pr.GetBody(),
		CreatedAt:           mapTimestamp(pr.CreatedAt),
		UpdatedAt:           mapTimestamp(pr.UpdatedAt),
		ClosedAt:            mapTimestamp(pr.ClosedAt),
		MergedAt:            mapTimestamp(pr.MergedAt),
		MergeCommitSHA:      pr.GetMergeCommitSHA(),
		User:                mapGitHubUserToUser(pr.User),
		HTMLURL:             pr.GetHTMLURL(),
		DiffURL:             pr.GetDiffURL(),
		PatchURL:            pr.GetPatchURL(),
		Merged:              pr.GetMerged(),
		Mergeable:           pr.GetMergeable(),
		MergeableState:      pr.GetMergeableState(),
		Comments:            pr.GetComments(),
		ReviewComments:      pr.GetReviewComments(),
		Commits:             pr.GetCommits(),
		Additions:           pr.GetAdditions(),
		Deletions:           pr.GetDeletions(),
		ChangedFiles:        pr.GetChangedFiles(),
		Draft:               pr.GetDraft(),
		MaintainerCanModify: pr.GetMaintainerCanModify(),
	}

	// ベースとヘッドの設定
	if pr.Base != nil {
		result.Base = Ref{
			Label: pr.Base.GetLabel(),
			Ref:   pr.Base.GetRef(),
			SHA:   pr.Base.GetSHA(),
			User:  mapGitHubUserToUser(pr.Base.User),
			Repo:  mapGitHubRepositoryToRepo(pr.Base.Repo),
		}
	}

	if pr.Head != nil {
		result.Head = Ref{
			Label: pr.Head.GetLabel(),
			Ref:   pr.Head.GetRef(),
			SHA:   pr.Head.GetSHA(),
			User:  mapGitHubUserToUser(pr.Head.User),
			Repo:  mapGitHubRepositoryToRepo(pr.Head.Repo),
		}
	}

	// レビューア情報があれば設定
	if pr.RequestedReviewers != nil {
		result.RequestedReviewers = make([]User, 0, len(pr.RequestedReviewers))
		for _, reviewer := range pr.RequestedReviewers {
			result.RequestedReviewers = append(result.RequestedReviewers, mapGitHubUserToUser(reviewer))
		}
	}

	return result, nil
}

// CreatePullRequestReview はPull Requestにレビューを作成します
func CreatePullRequestReview(options PullRequestReviewOptions, token string) (*PullRequestReview, error) {
	ctx := context.Background()
	client := getGitHubClient(ctx, token)

	// レビュー作成リクエストの設定
	review := &github.PullRequestReviewRequest{
		Body:     github.String(options.Body),
		Event:    github.String(options.Event),
		CommitID: github.String(options.CommitID),
	}

	// GitHub APIを呼び出してレビューを作成
	prReview, _, err := client.PullRequests.CreateReview(
		ctx,
		options.Owner,
		options.Repo,
		options.PullNumber,
		review,
	)
	if err != nil {
		return nil, err
	}

	// 結果をマッピング
	return &PullRequestReview{
		ID:          int(prReview.GetID()),
		User:        mapGitHubUserToUser(prReview.User),
		Body:        prReview.GetBody(),
		State:       prReview.GetState(),
		HTMLURL:     prReview.GetHTMLURL(),
		CommitID:    prReview.GetCommitID(),
		SubmittedAt: mapTimestamp(prReview.SubmittedAt),
	}, nil
}
