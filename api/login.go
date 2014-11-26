package api

import (
	"net/http"
)

func Login(email, password string) (*http.Response, error) {
	req := map[string]interface{}{
		"auth":     false,
		"method":   "POST",
		"endpoint": "/users/sign_in",
		"expected": Statuses{201, 401},
		"params": map[string]interface{}{
			"user": map[string]string{
				"login":    email,
				"password": password,
			},
		},
	}
	return Do(req)
}
