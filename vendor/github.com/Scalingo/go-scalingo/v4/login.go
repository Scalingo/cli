package scalingo

import (
	"encoding/json"
	"fmt"

	"github.com/Scalingo/go-scalingo/v4/http"

	"gopkg.in/errgo.v1"
)

type LoginService interface {
	Login(email, password string) (*LoginResponse, error)
}

var _ LoginService = (*Client)(nil)

type LoginError struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type LoginResponse struct {
	User *User `json:"user"`
}

func (err *LoginError) Error() string {
	return err.Message
}

func (c *Client) Login(email, password string) (*LoginResponse, error) {
	fmt.Println("[GO-SCALINGO] You are using the Login method. This method is deprecated, please use the OAuth flow")
	req := &http.APIRequest{
		NoAuth:   true,
		Method:   "POST",
		Endpoint: "/users/sign_in",
		Expected: http.Statuses{201, 401},
		Params: map[string]interface{}{
			"user": map[string]string{
				"login":    email,
				"password": password,
			},
		},
	}
	res, err := c.ScalingoAPI().Do(req)
	if err != nil {
		return nil, errgo.NoteMask(err, "fail to login", errgo.Any)
	}
	defer res.Body.Close()

	var loginRes LoginResponse
	err = json.NewDecoder(res.Body).Decode(&loginRes)
	if err != nil {
		return nil, errgo.NoteMask(err, "invalid response from server", errgo.Any)
	}
	return &loginRes, nil
}
