package scalingo

import (
	"encoding/json"

	"github.com/Scalingo/go-scalingo/http"
	"gopkg.in/errgo.v1"
)

type UsersService interface {
	Self() (*User, error)
	UpdateUser(params UpdateUserParams) (*User, error)
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

func (c *Client) Self() (*User, error) {
	req := &http.APIRequest{
		Endpoint: "/users/self",
	}
	res, err := c.AuthAPI().Do(req)
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
	StopFreeTrial bool   `json:"stop_free_trial,omitempty"`
	Password      string `json:"password,omitempty"`
	Email         string `json:"email,omitempty"`
}

type UpdateUserResponse struct {
	User *User `json:"user"`
}

func (c *Client) UpdateUser(params UpdateUserParams) (*User, error) {
	req := &http.APIRequest{
		Method:   "PATCH",
		Endpoint: "/account/profile",
		Params: map[string]interface{}{
			"user": params,
		},
		Expected: http.Statuses{200},
	}
	res, err := c.AuthAPI().Do(req)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	defer res.Body.Close()

	var u UpdateUserResponse
	err = json.NewDecoder(res.Body).Decode(&u)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	return u.User, nil
}
