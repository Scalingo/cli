package scalingo

import (
	"time"

	"github.com/Scalingo/go-scalingo/http"
)

type ScmRepoLinkService interface {
	ScmRepoLinkShow(app string) (*ScmRepoLink, error)
	ScmRepoLinkCreate(app string, params ScmRepoLinkParams) (*ScmRepoLink, error)
	ScmRepoLinkUpdate(app string, params ScmRepoLinkParams) (*ScmRepoLink, error)
	ScmRepoLinkDelete(app string) error

	ScmRepoLinkManualDeploy(app, branch string) error
	ScmRepoLinkManualReviewApp(app, pullRequestId string) error
	ScmRepoLinkDeployments(app string) ([]*Deployment, error)
	ScmRepoLinkReviewApps(app string) ([]*ReviewApp, error)
}

type ScmRepoLinkParams struct {
	Source                   *string `json:"source,omitempty"`
	Branch                   *string `json:"branch,omitempty"`
	AuthIntegrationID        *string `json:"auth_integration_id,omitempty"`
	ScmIntegrationUUID       *string `json:"scm_integration_uuid,omitempty"`
	AutoDeployEnabled        *bool   `json:"auto_deploy_enabled,omitempty"`
	DeployReviewAppsEnabled  *bool   `json:"deploy_review_apps_enabled,omitempty"`
	DestroyOnCloseEnabled    *bool   `json:"delete_on_close_enabled,omitempty"`
	HoursBeforeDeleteOnClose *uint   `json:"hours_before_delete_on_close,omitempty"`
	DestroyStaleEnabled      *bool   `json:"delete_stale_enabled,omitempty"`
	HoursBeforeDeleteStale   *uint   `json:"hours_before_delete_stale,omitempty"`
}

type ScmRepoLink struct {
	ID                       string            `json:"id"`
	AppID                    string            `json:"app_id"`
	Linker                   ScmRepoLinkLinker `json:"linker"`
	Owner                    string            `json:"owner"`
	Repo                     string            `json:"repo"`
	Branch                   string            `json:"branch"`
	CreatedAt                time.Time         `json:"created_at"`
	UpdatedAt                time.Time         `json:"updated_at"`
	AutoDeployEnabled        bool              `json:"auto_deploy_enabled"`
	ScmIntegrationUUID       string            `json:"scm_integration_uuid"`
	AuthIntegrationID        string            `json:"auth_integration_id"`
	DeployReviewAppsEnabled  bool              `json:"deploy_review_apps_enabled"`
	DeleteOnCloseEnabled     bool              `json:"delete_on_close_enabled"`
	DeleteStaleEnabled       bool              `json:"delete_stale_enabled"`
	HoursBeforeDeleteOnClose uint              `json:"hours_before_delete_on_close"`
	HoursBeforeDeleteStale   uint              `json:"hours_before_delete_stale"`
	LastAutoDeployAt         time.Time         `json:"last_auto_deploy_at"`
}

type ScmRepoLinkLinker struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	ID       string `json:"id"`
}

type ScmRepoLinkResponse struct {
	ScmRepoLink *ScmRepoLink `json:"scm_repo_link"`
}

type ScmRepoLinkDeploymentsResponse struct {
	Deployments []*Deployment `json:"deployments"`
}

type ScmRepoLinkReviewAppsResponse struct {
	ReviewApps []*ReviewApp `json:"review_apps"`
}

var _ ScmRepoLinkService = (*Client)(nil)

func (c *Client) ScmRepoLinkShow(app string) (*ScmRepoLink, error) {
	var res ScmRepoLinkResponse
	err := c.ScalingoAPI().DoRequest(&http.APIRequest{
		Method:   "GET",
		Endpoint: "/apps/" + app + "/scm_repo_link",
		Expected: http.Statuses{200},
	}, &res)
	if err != nil {
		return nil, err
	}
	return res.ScmRepoLink, nil
}

func (c *Client) ScmRepoLinkCreate(app string, params ScmRepoLinkParams) (*ScmRepoLink, error) {
	var res ScmRepoLinkResponse
	err := c.ScalingoAPI().DoRequest(&http.APIRequest{
		Method:   "POST",
		Endpoint: "/apps/" + app + "/scm_repo_link",
		Expected: http.Statuses{201},
		Params:   map[string]ScmRepoLinkParams{"scm_repo_link": params},
	}, &res)
	if err != nil {
		return nil, err
	}

	return res.ScmRepoLink, nil
}

func (c *Client) ScmRepoLinkUpdate(app string, params ScmRepoLinkParams) (*ScmRepoLink, error) {
	var res ScmRepoLinkResponse
	err := c.ScalingoAPI().DoRequest(&http.APIRequest{
		Method:   "UPDATE",
		Endpoint: "/apps/" + app + "/scm_repo_link",
		Expected: http.Statuses{200},
		Params:   map[string]ScmRepoLinkParams{"scm_repo_link": params},
	}, &res)
	if err != nil {
		return nil, err
	}

	return res.ScmRepoLink, nil
}

func (c *Client) ScmRepoLinkDelete(app string) error {
	_, err := c.ScalingoAPI().Do(&http.APIRequest{
		Method:   "DELETE",
		Endpoint: "/apps/" + app + "/scm_repo_link",
		Expected: http.Statuses{204},
	})
	return err
}

func (c *Client) ScmRepoLinkManualDeploy(app, branch string) error {
	_, err := c.ScalingoAPI().Do(&http.APIRequest{
		Method:   "POST",
		Endpoint: "/apps/" + app + "/scm_repo_link/manual_deploy",
		Expected: http.Statuses{200},
		Params:   map[string]string{"branch": branch},
	})
	return err
}

func (c *Client) ScmRepoLinkManualReviewApp(app, pullRequestId string) error {
	_, err := c.ScalingoAPI().Do(&http.APIRequest{
		Method:   "POST",
		Endpoint: "/apps/" + app + "/scm_repo_link/manual_review_app",
		Expected: http.Statuses{200},
		Params:   map[string]string{"pull_request_id": pullRequestId},
	})
	return err
}

func (c *Client) ScmRepoLinkDeployments(app string) ([]*Deployment, error) {
	var res ScmRepoLinkDeploymentsResponse

	err := c.ScalingoAPI().DoRequest(&http.APIRequest{
		Method:   "GET",
		Endpoint: "/apps/" + app + "/scm_repo_link/deployments",
		Expected: http.Statuses{200},
	}, &res)
	if err != nil {
		return nil, err
	}
	return res.Deployments, nil
}

func (c *Client) ScmRepoLinkReviewApps(app string) ([]*ReviewApp, error) {
	var res ScmRepoLinkReviewAppsResponse

	err := c.ScalingoAPI().DoRequest(&http.APIRequest{
		Method:   "GET",
		Endpoint: "/apps/" + app + "/scm_repo_link/review_apps",
		Expected: http.Statuses{200},
	}, &res)
	if err != nil {
		return nil, err
	}
	return res.ReviewApps, nil
}
