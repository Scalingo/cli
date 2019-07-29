package scalingo

import (
	"time"
)

type PullRequest struct {
	Number     int       `json:"number"`
	BranchName string    `json:"branch_name"`
	Title      string    `json:"title"`
	Url        string    `json:"url"`
	HtmlUrl    string    `json:"html_url"`
	Ref        string    `json:"ref"`
	BaseRef    string    `json:"base_ref"`
	CreatedAt  time.Time `json:"created_at"`
	ClosedAt   time.Time `json:"closed_at"`
}

type ReviewApp struct {
	ID                  string       `json:"id"`
	RepoLinkID          string       `json:"repo_link_id"`
	AppID               string       `json:"app_id"`
	AppName             string       `json:"app_name"`
	ParentAppID         string       `json:"parent_app_id"`
	ParentAppName       string       `json:"parent_app_name"`
	CreatedAt           time.Time    `json:"created_at"`
	StaleDeletionDate   time.Time    `json:"stale_deletion_date"`
	OnCloseDeletionDate time.Time    `json:"on_close_deletion_date"`
	PullRequest         *PullRequest `json:"pull_request"`
	LastDeployment      *Deployment  `json:"last_deployment"`
}
