package api

import (
	"net/http"
)

func Login(email, password string) (*http.Response, error) {
	req := map[string]interface{}{
		"auth":     false,
		"method":   "POST",
		"endpoint": "/users/sign_in",
		"params": map[string]interface{}{
			"user": map[string]string{
				"login":    email,
				"password": password,
			},
		},
	}
	return Do(req)
}
