package scalingo

import (
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/go-scalingo/v4/http"
)

type SCMType string

// Type of SCM integrations
const (
	SCMGithubType           SCMType = "github"             // GitHub
	SCMGithubEnterpriseType SCMType = "github-enterprise"  // GitHub Enterprise (private instance)
	SCMGitlabType           SCMType = "gitlab"             // GitLab.com
	SCMGitlabSelfHostedType SCMType = "gitlab-self-hosted" // GitLab self-hosted (private instance)
)

var SCMTypeDisplay = map[SCMType]string{
	SCMGithubType:           "GitHub",
	SCMGitlabType:           "GitLab",
	SCMGithubEnterpriseType: "GitHub Enterprise",
	SCMGitlabSelfHostedType: "GitLab self-hosted",
}

func (t SCMType) Str() string {
	return string(t)
}

type SCMIntegrationsService interface {
	SCMIntegrationsList() ([]SCMIntegration, error)
	SCMIntegrationsShow(id string) (*SCMIntegration, error)
	SCMIntegrationsCreate(scmType SCMType, url string, accessToken string) (*SCMIntegration, error)
	SCMIntegrationsDelete(id string) error
	SCMIntegrationsImportKeys(id string) ([]Key, error)
}

var _ SCMIntegrationsService = (*Client)(nil)

type SCMIntegration struct {
	ID          string  `json:"id"`
	SCMType     SCMType `json:"scm_type"`
	URL         string  `json:"url,omitempty"`
	AccessToken string  `json:"access_token"`
	Uid         string  `json:"uid"`
	Username    string  `json:"username"`
	Email       string  `json:"email"`
	AvatarURL   string  `json:"avatar_url"`
	ProfileURL  string  `json:"profile_url"`
}

type SCMIntegrationParams struct {
	SCMType     SCMType `json:"scm_type"`
	URL         string  `json:"url,omitempty"`
	AccessToken string  `json:"access_token"`
}

type SCMIntegrationRes struct {
	SCMIntegration SCMIntegration `json:"scm_integration"`
}

type SCMIntegrationsRes struct {
	SCMIntegrations []SCMIntegration `json:"scm_integrations"`
}

type SCMIntegrationParamsReq struct {
	SCMIntegrationParams SCMIntegrationParams `json:"scm_integration"`
}

func (c *Client) SCMIntegrationsList() ([]SCMIntegration, error) {
	var res SCMIntegrationsRes

	err := c.AuthAPI().ResourceList("scm_integrations", nil, &res)
	if err != nil {
		return nil, errgo.Notef(err, "fail to list SCM integration")
	}
	return res.SCMIntegrations, nil
}

func (c *Client) SCMIntegrationsShow(id string) (*SCMIntegration, error) {
	var res SCMIntegrationRes

	err := c.AuthAPI().ResourceGet("scm_integrations", id, nil, &res)
	if err != nil {
		return nil, errgo.Notef(err, "fail to get this SCM integration")
	}
	return &res.SCMIntegration, nil
}

func (c *Client) SCMIntegrationsCreate(scmType SCMType, url string, accessToken string) (*SCMIntegration, error) {
	payload := SCMIntegrationParamsReq{SCMIntegrationParams{
		SCMType:     scmType,
		URL:         url,
		AccessToken: accessToken,
	}}
	var res SCMIntegrationRes

	err := c.AuthAPI().ResourceAdd("scm_integrations", payload, &res)
	if err != nil {
		return nil, errgo.Notef(err, "fail to create the SCM integration")
	}

	return &res.SCMIntegration, nil
}

func (c *Client) SCMIntegrationsDelete(id string) error {
	err := c.AuthAPI().ResourceDelete("scm_integrations", id)
	if err != nil {
		return errgo.Notef(err, "fail to delete this SCM integration")
	}
	return nil
}

func (c *Client) SCMIntegrationsImportKeys(id string) ([]Key, error) {
	var res KeysRes

	var err = c.AuthAPI().DoRequest(&http.APIRequest{
		Method:   "POST",
		Endpoint: "/scm_integrations/" + id + "/import_keys",
		Params:   nil,
		Expected: http.Statuses{201},
	}, &res)
	if err != nil {
		return nil, errgo.Notef(err, "fail to import ssh keys from this SCM integration")
	}
	return res.Keys, nil
}
