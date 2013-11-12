package api

import (
	"net/http"
)

func Login(email, password string) (*http.Response, error) {
	req := map[string]interface{}{
		"method":   "POST",
		"endpoint": "/users/sign_in",
		"params": map[string]interface{}{
			"user": map[string]string{
				"email":    email,
				"password": password,
			},
		},
	}
	return Do(req)
}
