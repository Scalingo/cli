package scalingo

import (
	"github.com/Scalingo/go-scalingo/v4/http"
	"gopkg.in/errgo.v1"
)

type SignUpService interface {
	SignUp(email, password string) error
}

var _ SignUpService = (*Client)(nil)

func (c *Client) SignUp(email, password string) error {
	req := &http.APIRequest{
		NoAuth:   true,
		Method:   "POST",
		Endpoint: "/users",
		Expected: http.Statuses{201},
		Params: map[string]interface{}{
			"user": map[string]string{
				"email":    email,
				"password": password,
			},
		},
	}
	_, err := c.ScalingoAPI().Do(req)
	if err != nil {
		return errgo.Mask(err)
	}

	return nil
}
