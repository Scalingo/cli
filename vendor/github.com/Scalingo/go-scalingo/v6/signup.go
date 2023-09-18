package scalingo

import (
	"context"

	"gopkg.in/errgo.v1"

	"github.com/Scalingo/go-scalingo/v6/http"
)

type SignUpService interface {
	SignUp(ctx context.Context, email, password string) error
}

var _ SignUpService = (*Client)(nil)

func (c *Client) SignUp(ctx context.Context, email, password string) error {
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
	res, err := c.ScalingoAPI().Do(ctx, req)
	if err != nil {
		return errgo.Mask(err)
	}
	defer res.Body.Close()

	return nil
}
