package scalingo

import (
	"context"
	"time"

	"gopkg.in/errgo.v1"

	"github.com/Scalingo/go-scalingo/v4/http"
)

type SCMRepoLinkService interface {
	SCMRepoLinkList(ctx context.Context, opts PaginationOpts) ([]*SCMRepoLink, PaginationMeta, error)
	SCMRepoLinkShow(ctx context.Context, app string) (*SCMRepoLink, error)
	SCMRepoLinkCreate(ctx context.Context, app string, params SCMRepoLinkCreateParams) (*SCMRepoLink, error)
	SCMRepoLinkUpdate(ctx context.Context, app string, params SCMRepoLinkUpdateParams) (*SCMRepoLink, error)
	SCMRepoLinkDelete(ctx context.Context, app string) error

	SCMRepoLinkManualDeploy(ctx context.Context, app, branch string) error
	SCMRepoLinkManualReviewApp(ctx context.Context, app, pullRequestID string) error
	SCMRepoLinkDeployments(ctx context.Context, app string) ([]*Deployment, error)
	SCMRepoLinkReviewApps(ctx context.Context, app string) ([]*ReviewApp, error)
}

type SCMRepoLinkCreateParams struct {
	Source                   *string `json:"source,omitempty"`
	Branch                   *string `json:"branch,omitempty"`
	AuthIntegrationUUID      *string `json:"auth_integration_uuid,omitempty"`
	AutoDeployEnabled        *bool   `json:"auto_deploy_enabled,omitempty"`
	DeployReviewAppsEnabled  *bool   `json:"deploy_review_apps_enabled,omitempty"`
	DestroyOnCloseEnabled    *bool   `json:"delete_on_close_enabled,omitempty"`
	HoursBeforeDeleteOnClose *uint   `json:"hours_before_delete_on_close,omitempty"`
	DestroyStaleEnabled      *bool   `json:"delete_stale_enabled,omitempty"`
	HoursBeforeDeleteStale   *uint   `json:"hours_before_delete_stale,omitempty"`
}

type SCMRepoLinkUpdateParams struct {
	Branch                   *string `json:"branch,omitempty"`
	AutoDeployEnabled        *bool   `json:"auto_deploy_enabled,omitempty"`
	DeployReviewAppsEnabled  *bool   `json:"deploy_review_apps_enabled,omitempty"`
	DestroyOnCloseEnabled    *bool   `json:"delete_on_close_enabled,omitempty"`
	HoursBeforeDeleteOnClose *uint   `json:"hours_before_delete_on_close,omitempty"`
	DestroyStaleEnabled      *bool   `json:"delete_stale_enabled,omitempty"`
	HoursBeforeDeleteStale   *uint   `json:"hours_before_delete_stale,omitempty"`
}

type SCMRepoLink struct {
	ID                       string            `json:"id"`
	AppID                    string            `json:"app_id"`
	Linker                   SCMRepoLinkLinker `json:"linker"`
	Owner                    string            `json:"owner"`
	Repo                     string            `json:"repo"`
	Branch                   string            `json:"branch"`
	SCMType                  SCMType           `json:"scm_type"`
	CreatedAt                time.Time         `json:"created_at"`
	UpdatedAt                time.Time         `json:"updated_at"`
	AutoDeployEnabled        bool              `json:"auto_deploy_enabled"`
	AuthIntegrationUUID      string            `json:"auth_integration_uuid"`
	DeployReviewAppsEnabled  bool              `json:"deploy_review_apps_enabled"`
	DeleteOnCloseEnabled     bool              `json:"delete_on_close_enabled"`
	DeleteStaleEnabled       bool              `json:"delete_stale_enabled"`
	HoursBeforeDeleteOnClose uint              `json:"hours_before_delete_on_close"`
	HoursBeforeDeleteStale   uint              `json:"hours_before_delete_stale"`
	LastAutoDeployAt         time.Time         `json:"last_auto_deploy_at"`
}

type SCMRepoLinkLinker struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	ID       string `json:"id"`
}

type SCMRepoLinksResponse struct {
	SCMRepoLinks []*SCMRepoLink `json:"scm_repo_links"`
	Meta         struct {
		PaginationMeta PaginationMeta `json:"pagination"`
	}
}

type ScmRepoLinkResponse struct {
	SCMRepoLink *SCMRepoLink `json:"scm_repo_link"`
}

type SCMRepoLinkDeploymentsResponse struct {
	Deployments []*Deployment `json:"deployments"`
}

type SCMRepoLinkReviewAppsResponse struct {
	ReviewApps []*ReviewApp `json:"review_apps"`
}

var _ SCMRepoLinkService = (*Client)(nil)

func (c *Client) SCMRepoLinkList(ctx context.Context, opts PaginationOpts) ([]*SCMRepoLink, PaginationMeta, error) {
	var res SCMRepoLinksResponse
	err := c.ScalingoAPI().ResourceList(ctx, "scm_repo_links", opts.ToMap(), &res)
	if err != nil {
		return nil, PaginationMeta{}, errgo.Notef(err, "fail to list SCM repo links")
	}
	return res.SCMRepoLinks, res.Meta.PaginationMeta, nil
}

func (c *Client) SCMRepoLinkShow(ctx context.Context, app string) (*SCMRepoLink, error) {
	var res ScmRepoLinkResponse
	err := c.ScalingoAPI().DoRequest(ctx, &http.APIRequest{
		Method:   "GET",
		Endpoint: "/apps/" + app + "/scm_repo_link",
		Expected: http.Statuses{200},
	}, &res)
	if err != nil {
		return nil, errgo.Notef(err, "fail to get this SCM repo link")
	}
	return res.SCMRepoLink, nil
}

func (c *Client) SCMRepoLinkCreate(ctx context.Context, app string, params SCMRepoLinkCreateParams) (*SCMRepoLink, error) {
	var res ScmRepoLinkResponse
	err := c.ScalingoAPI().DoRequest(ctx, &http.APIRequest{
		Method:   "POST",
		Endpoint: "/apps/" + app + "/scm_repo_link",
		Expected: http.Statuses{201},
		Params:   map[string]SCMRepoLinkCreateParams{"scm_repo_link": params},
	}, &res)
	if err != nil {
		return nil, errgo.Notef(err, "fail to create the SCM repo link")
	}

	return res.SCMRepoLink, nil
}

func (c *Client) SCMRepoLinkUpdate(ctx context.Context, app string, params SCMRepoLinkUpdateParams) (*SCMRepoLink, error) {
	var res ScmRepoLinkResponse
	err := c.ScalingoAPI().DoRequest(ctx, &http.APIRequest{
		Method:   "PATCH",
		Endpoint: "/apps/" + app + "/scm_repo_link",
		Expected: http.Statuses{200},
		Params:   map[string]SCMRepoLinkUpdateParams{"scm_repo_link": params},
	}, &res)
	if err != nil {
		return nil, errgo.Notef(err, "fail to update this SCM repo link")
	}

	return res.SCMRepoLink, nil
}

func (c *Client) SCMRepoLinkDelete(ctx context.Context, app string) error {
	res, err := c.ScalingoAPI().Do(ctx, &http.APIRequest{
		Method:   "DELETE",
		Endpoint: "/apps/" + app + "/scm_repo_link",
		Expected: http.Statuses{204},
	})
	if err != nil {
		return errgo.Notef(err, "fail to delete this SCM repo link")
	}
	defer res.Body.Close()

	return nil
}

func (c *Client) SCMRepoLinkManualDeploy(ctx context.Context, app, branch string) error {
	res, err := c.ScalingoAPI().Do(ctx, &http.APIRequest{
		Method:   "POST",
		Endpoint: "/apps/" + app + "/scm_repo_link/manual_deploy",
		Expected: http.Statuses{200},
		Params:   map[string]string{"branch": branch},
	})
	if err != nil {
		return errgo.Notef(err, "fail to trigger manual app deployment")
	}
	defer res.Body.Close()

	return nil
}

func (c *Client) SCMRepoLinkManualReviewApp(ctx context.Context, app, pullRequestID string) error {
	_, err := c.ScalingoAPI().Do(ctx, &http.APIRequest{
		Method:   "POST",
		Endpoint: "/apps/" + app + "/scm_repo_link/manual_review_app",
		Expected: http.Statuses{200},
		Params:   map[string]string{"pull_request_id": pullRequestID},
	})
	if err != nil {
		return errgo.Notef(err, "fail to trigger manual review app deployment")
	}
	return nil
}

func (c *Client) SCMRepoLinkDeployments(ctx context.Context, app string) ([]*Deployment, error) {
	var res SCMRepoLinkDeploymentsResponse

	err := c.ScalingoAPI().DoRequest(ctx, &http.APIRequest{
		Method:   "GET",
		Endpoint: "/apps/" + app + "/scm_repo_link/deployments",
		Expected: http.Statuses{200},
	}, &res)
	if err != nil {
		return nil, errgo.Notef(err, "fail to list deployments of this SCM repo link")
	}
	return res.Deployments, nil
}

func (c *Client) SCMRepoLinkReviewApps(ctx context.Context, app string) ([]*ReviewApp, error) {
	var res SCMRepoLinkReviewAppsResponse

	err := c.ScalingoAPI().DoRequest(ctx, &http.APIRequest{
		Method:   "GET",
		Endpoint: "/apps/" + app + "/scm_repo_link/review_apps",
		Expected: http.Statuses{200},
	}, &res)
	if err != nil {
		return nil, errgo.Notef(err, "fail to list review apps of this SCM repo link")
	}
	return res.ReviewApps, nil
}
