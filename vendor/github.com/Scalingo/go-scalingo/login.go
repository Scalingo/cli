package scalingo

import (
	"fmt"

	"gopkg.in/errgo.v1"
)

type LoginError struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type LoginResponse struct {
	AuthenticationToken string `json:"authentication_token"`
	User                *User  `json:"user"`
}

func (err *LoginError) Error() string {
	return err.Message
}

func (c *Client) Login(email, password string) (*LoginResponse, error) {
	fmt.Println("[GO-SCALINGO] You are using the Login method. This method is deprecated, please use the OAuth flow")
	req := &APIRequest{
		Client:   c,
		NoAuth:   true,
		Method:   "POST",
		Endpoint: "/users/sign_in",
		Expected: Statuses{201, 401},
		Params: map[string]interface{}{
			"user": map[string]string{
				"login":    email,
				"password": password,
			},
		},
	}
	res, err := req.Do()
	if err != nil {
		return nil, errgo.NoteMask(err, "fail to login", errgo.Any)
	}
	defer res.Body.Close()

	var loginRes LoginResponse
	err = ParseJSON(res, &loginRes)
	if err != nil {
		return nil, errgo.NoteMask(err, "invalid response from server", errgo.Any)
	}
	return &loginRes, nil
}
