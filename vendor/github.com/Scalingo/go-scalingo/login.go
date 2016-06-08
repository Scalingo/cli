package scalingo

import "gopkg.in/errgo.v1"

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
		return nil, errgo.Mask(err)
	}
	defer res.Body.Close()

	var loginRes LoginResponse
	err = ParseJSON(res, &loginRes)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	return &loginRes, nil
}
