package scalingo

import (
	"context"
	"net/http"

	httpclient "github.com/Scalingo/go-scalingo/v11/http"
	"github.com/Scalingo/go-utils/errors/v3"
)

type UsersService interface {
	Self(context.Context) (*User, error)
	UpdateUser(context.Context, UpdateUserParams) (*User, error)
	UserStopFreeTrial(context.Context) error
}

var _ UsersService = (*Client)(nil)

type User struct {
	ID       string          `json:"id"`
	Username string          `json:"username"`
	Fullname string          `json:"fullname"`
	Email    string          `json:"email"`
	Flags    map[string]bool `json:"flags"`
}

type SelfResponse struct {
	User *User `json:"user"`
}

func (c *Client) Self(ctx context.Context) (*User, error) {
	var selfRes SelfResponse
	req := &httpclient.APIRequest{
		Endpoint: "/users/self",
	}
	err := c.AuthAPI().DoRequest(ctx, req, &selfRes)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "get current user")
	}
	return selfRes.User, nil
}

type UpdateUserParams struct {
	Password string `json:"password,omitempty"`
	Email    string `json:"email,omitempty"`
}

type UpdateUserResponse struct {
	User *User `json:"user"`
}

func (c *Client) UpdateUser(ctx context.Context, params UpdateUserParams) (*User, error) {
	if params.Password == "" && params.Email == "" {
		return nil, nil
	}

	req := &httpclient.APIRequest{
		Method:   http.MethodPatch,
		Endpoint: "/account/profile",
		Params: map[string]any{
			"user": params,
		},
		Expected: httpclient.Statuses{http.StatusOK},
	}
	var updateUserRes UpdateUserResponse
	err := c.AuthAPI().DoRequest(ctx, req, &updateUserRes)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "execute the query to update the user")
	}

	return updateUserRes.User, nil
}

func (c *Client) UserStopFreeTrial(ctx context.Context) error {
	req := &httpclient.APIRequest{
		Method:   http.MethodPost,
		Endpoint: "/users/stop_free_trial",
		Params:   map[string]any{},
		Expected: httpclient.Statuses{http.StatusOK},
	}

	err := c.AuthAPI().DoRequest(ctx, req, nil)
	if err != nil {
		return errors.Wrap(ctx, err, "execute the query to stop user free trial")
	}

	return nil
}
