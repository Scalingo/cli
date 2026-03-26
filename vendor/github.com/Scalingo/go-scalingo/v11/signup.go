package scalingo

import (
	"context"
	"net/http"

	httpclient "github.com/Scalingo/go-scalingo/v11/http"
	"github.com/Scalingo/go-utils/errors/v3"
)

type SignUpService interface {
	SignUp(ctx context.Context, email, password string) error
}

var _ SignUpService = (*Client)(nil)

func (c *Client) SignUp(ctx context.Context, email, password string) error {
	req := &httpclient.APIRequest{
		NoAuth:   true,
		Method:   http.MethodPost,
		Endpoint: "/users",
		Expected: httpclient.Statuses{http.StatusCreated},
		Params: map[string]any{
			"user": map[string]string{
				"email":    email,
				"password": password,
			},
		},
	}
	err := c.ScalingoAPI().DoRequest(ctx, req, nil)
	if err != nil {
		return errors.Wrap(ctx, err, "sign up user")
	}

	return nil
}
