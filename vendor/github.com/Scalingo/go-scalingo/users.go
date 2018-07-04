package scalingo

import "gopkg.in/errgo.v1"

type UsersService interface {
	Self() (*User, error)
	UpdateUser(params UpdateUserParams) (*User, error)
}

var _ UsersService = (*Client)(nil)

type User struct {
	ID                  string          `json:"id"`
	Username            string          `json:"username"`
	Fullname            string          `json:"fullname"`
	Email               string          `json:"email"`
	Flags               map[string]bool `json:"flags"`
	AuthenticationToken string          `json:"authentication_token"`
}

type SelfResponse struct {
	User *User `json:"user"`
}

func (c *Client) Self() (*User, error) {
	req := &APIRequest{
		Client:   c,
		Endpoint: "/users/self",
	}
	res, err := req.Do()
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}
	defer res.Body.Close()

	var u SelfResponse
	err = ParseJSON(res, &u)
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
	req := &APIRequest{
		Client:   c,
		Method:   "PATCH",
		Endpoint: "/account/profile",
		Params: map[string]interface{}{
			"user": params,
		},
		Expected: Statuses{200},
	}
	res, err := req.Do()
	if err != nil {
		return nil, errgo.Mask(err)
	}
	defer res.Body.Close()

	var u UpdateUserResponse
	err = ParseJSON(res, &u)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	return u.User, nil
}
