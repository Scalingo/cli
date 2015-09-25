package scalingo

import (
	"net/http"
)

func Login(email, password string) (*http.Response, error) {
	req := &APIRequest{
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
	return req.Do()
}
