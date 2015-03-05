package api

import "github.com/Scalingo/cli/Godeps/_workspace/src/gopkg.in/errgo.v1"

func SignUp(email, password string) error {
	req := &APIRequest{
		NoAuth:   true,
		Method:   "POST",
		Endpoint: "/users",
		Expected: Statuses{201},
		Params: map[string]interface{}{
			"user": map[string]string{
				"email":    email,
				"password": password,
			},
		},
	}
	_, err := req.Do()
	if err != nil {
		return errgo.Mask(err)
	}

	return nil
}
