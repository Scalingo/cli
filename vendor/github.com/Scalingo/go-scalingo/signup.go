package scalingo

import "gopkg.in/errgo.v1"

type SignUpService interface {
	SignUp(email, password string) error
}

type SignUpClient struct {
	*backendConfiguration
}

func (c *SignUpClient) SignUp(email, password string) error {
	req := &APIRequest{
		Client:   c.backendConfiguration,
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
