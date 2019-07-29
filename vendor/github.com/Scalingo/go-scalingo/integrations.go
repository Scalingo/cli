package scalingo

import (
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/go-scalingo/http"
)

type IntegrationsService interface {
	IntegrationsList() ([]Integration, error)
	IntegrationsCreate(scmType string, url string, accessToken string) (*Integration, error)
	IntegrationsDestroy(id string) error
}

var _ IntegrationsService = (*Client)(nil)

type Integration struct {
	ID          string `json:"id,omitempty"`
	ScmType     string `json:"scm_type"`
	Url         string `json:"url"`
	AccessToken string `json:"access_token"`
	Uid         string `json:"uid,omitempty"`
	Username    string `json:"username,omitempty"`
	Email       string `json:"email,omitempty"`
	AvatarUrl   string `json:"avatar_url,omitempty"`
	ProfileUrl  string `json:"profile_url,omitempty"`
}

type IntegrationRes struct {
	Integration Integration `json:"integration"`
}

type IntegrationsRes struct {
	Integrations []Integration `json:"integrations"`
}

func (c *Client) IntegrationsList() ([]Integration, error) {
	var res IntegrationsRes

	err := c.AuthAPI().ResourceList("integrations", nil, &res)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	return res.Integrations, nil
}

func (c *Client) IntegrationsShow(id string) (*Integration, error) {
	var res IntegrationRes

	err := c.AuthAPI().ResourceGet("integrations", id, nil, &res)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	return &res.Integration, nil
}

func (c *Client) IntegrationsCreate(scmType string, url string, accessToken string) (*Integration, error) {
	payload := IntegrationRes{Integration{
		ScmType:     scmType,
		Url:         url,
		AccessToken: accessToken,
	}}
	var res IntegrationRes

	err := c.AuthAPI().ResourceAdd("integrations", payload, &res)
	if err != nil {
		return nil, errgo.Mask(err)
	}

	return &res.Integration, nil
}

func (c *Client) IntegrationsDestroy(id string) error {
	err := c.AuthAPI().ResourceDelete("integrations", id)
	if err != nil {
		return errgo.Mask(err)
	}
	return nil
}

func (c *Client) IntegrationsImportKeys(id string) ([]Key, error) {
	var res KeysRes

	var err = c.AuthAPI().DoRequest(&http.APIRequest{
		Method:   "POST",
		Endpoint: "/integrations/" + id + "/import_keys",
		Params:   nil,
		Expected: http.Statuses{201},
	}, &res)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	return res.Keys, nil
}
