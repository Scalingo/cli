package scalingo

import "gopkg.in/errgo.v1"

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
