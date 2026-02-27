package scalingo

import (
	"context"

	"github.com/Scalingo/go-scalingo/v9/http"
	"github.com/Scalingo/go-utils/errors/v3"
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
		return errors.Wrap(ctx, err, "sign up user")
	}
	defer res.Body.Close()

	return nil
}
