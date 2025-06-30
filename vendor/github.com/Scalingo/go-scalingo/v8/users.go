package scalingo

import (
	"context"
	"encoding/json"

	"gopkg.in/errgo.v1"

	"github.com/Scalingo/go-scalingo/v8/http"
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
	req := &http.APIRequest{
		Endpoint: "/users/self",
	}
	res, err := c.AuthAPI().Do(ctx, req)
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}
	defer res.Body.Close()

	var u SelfResponse
	err = json.NewDecoder(res.Body).Decode(&u)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	return u.User, nil
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

	req := &http.APIRequest{
		Method:   "PATCH",
		Endpoint: "/account/profile",
		Params: map[string]interface{}{
			"user": params,
		},
		Expected: http.Statuses{200},
	}
	res, err := c.AuthAPI().Do(ctx, req)
	if err != nil {
		return nil, errgo.Notef(err, "fail to execute the query to update the user")
	}
	defer res.Body.Close()

	var u UpdateUserResponse
	err = json.NewDecoder(res.Body).Decode(&u)
	if err != nil {
		return nil, errgo.Notef(err, "fail to decode response of the query to update the user")
	}

	return u.User, nil
}

func (c *Client) UserStopFreeTrial(ctx context.Context) error {
	req := &http.APIRequest{
		Method:   "POST",
		Endpoint: "/users/stop_free_trial",
		Params:   map[string]interface{}{},
		Expected: http.Statuses{200},
	}

	res, err := c.AuthAPI().Do(ctx, req)
	if err != nil {
		return errgo.Notef(err, "fail to execute the query to stop user free trial")
	}
	defer res.Body.Close()

	return nil
}
