package api

import "gopkg.in/errgo.v1"

func SignUp(email, password string) error {
	req := map[string]interface{}{
		"auth":     false,
		"method":   "POST",
		"endpoint": "/users",
		"expected": Statuses{201},
		"params": map[string]interface{}{
			"user": map[string]string{
				"email":    email,
				"password": password,
			},
		},
	}
	_, err := Do(req)
	if err != nil {
		return errgo.Mask(err)
	}

	return nil
}
