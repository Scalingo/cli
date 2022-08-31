package scalingo

import (
	"context"
	"errors"
	"time"

	"github.com/Scalingo/go-scalingo/v4/http"
)

var (
	GithubLinkNotFoundErr = errors.New("No github SCM Repo Link")
)

type GithubLinkService interface {
	GithubLinkShow(ctx context.Context, app string) (*GithubLink, error)
	GithubLinkAdd(ctx context.Context, app string, params GithubLinkParams) (*GithubLink, error)
	GithubLinkUpdate(ctx context.Context, app, id string, params GithubLinkParams) (*GithubLink, error)
	GithubLinkDelete(ctx context.Context, app string, id string) error
	GithubLinkManualDeploy(ctx context.Context, app, id, branch string) error
}

type GithubLinkParams struct {
	GithubSource             *string `json:"github_source,omitempty"`
	GithubBranch             *string `json:"github_branch,omitempty"`
	AutoDeployEnabled        *bool   `json:"auto_deploy_enabled,omitempty"`
	DeployReviewAppsEnabled  *bool   `json:"deploy_review_apps_enabled,omitempty"`
	DestroyOnCloseEnabled    *bool   `json:"delete_on_close_enabled,omitempty"`
	HoursBeforeDeleteOnClose *uint   `json:"hours_before_delete_on_close,omitempty"`
	DestroyStaleEnabled      *bool   `json:"delete_stale_enabled,omitempty"`
	HoursBeforeDeleteStale   *uint   `json:"hours_before_delete_stale,omitempty"`
}

type GithubLink struct {
	ID                       string           `json:"id"`
	Linker                   GithubLinkLinker `json:"linker"`
	CreatedAt                time.Time        `json:"created_at"`
	UpdatedAt                time.Time        `json:"updated_at"`
	GithubSource             string           `json:"github_source"`
	GithubBranch             string           `json:"github_branch"`
	AutoDeployEnabled        bool             `json:"auto_deploy_enabled"`
	GithubIntegrationUUID    string           `json:"github_integration_uuid"`
	DeployReviewAppsEnabled  bool             `json:"deploy_review_apps_enabled"`
	DestroyOnCloseEnabled    bool             `json:"delete_on_close_enabled"`
	DestroyOnStaleEnabled    bool             `json:"delete_stale_enabled"`
	HoursBeforeDeleteOnClose uint             `json:"hours_before_delete_on_close"`
	HoursBeforeDeleteStale   uint             `json:"hours_before_delete_stale"`
	LastAutoDeployAt         time.Time        `json:"last_auto_deploy_at"`
}

type GithubLinkLinker struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	ID       string `json:"id"`
}

type GithubLinkResponse struct {
	GithubLink *GithubLink `json:"github_repo_link"`
}

type GithubLinksResponse struct {
	GithubLinks []*GithubLink `json:"github_repo_links"`
}

var _ GithubLinkService = (*Client)(nil)

func (c *Client) GithubLinkShow(ctx context.Context, app string) (*GithubLink, error) {
	var link GithubLinksResponse
	err := c.ScalingoAPI().SubresourceList(ctx, "apps", app, "github_repo_links", nil, &link)
	if err != nil {
		return nil, err
	}

	if len(link.GithubLinks) == 0 {
		return nil, GithubLinkNotFoundErr
	}

	return link.GithubLinks[0], nil
}

func (c *Client) GithubLinkAdd(ctx context.Context, app string, params GithubLinkParams) (*GithubLink, error) {
	linkParams := map[string]GithubLinkParams{
		"github_repo_link": params,
	}

	var link GithubLinkResponse
	err := c.ScalingoAPI().SubresourceAdd(ctx, "apps", app, "github_repo_links", linkParams, &link)
	if err != nil {
		return nil, err
	}
	return link.GithubLink, nil
}

func (c *Client) GithubLinkUpdate(ctx context.Context, app, id string, params GithubLinkParams) (*GithubLink, error) {
	linkParams := map[string]GithubLinkParams{
		"github_repo_link": params,
	}

	var link GithubLinkResponse
	err := c.ScalingoAPI().SubresourceUpdate(ctx, "apps", app, "github_repo_links", id, linkParams, &link)
	if err != nil {
		return nil, err
	}
	return link.GithubLink, nil
}

func (c *Client) GithubLinkDelete(ctx context.Context, app, id string) error {
	return c.ScalingoAPI().SubresourceDelete(ctx, "apps", app, "github_repo_links", id)
}

func (c *Client) GithubLinkManualDeploy(ctx context.Context, app, id, branch string) error {
	req := &http.APIRequest{
		Method:   "POST",
		Endpoint: "/apps/" + app + "/github_repo_links/" + id + "/manual_deploy",
		Expected: http.Statuses{200},
		Params: map[string]string{
			"branch": branch,
		},
	}
	res, err := c.ScalingoAPI().Do(ctx, req)
	defer res.Body.Close()

	return err
}
