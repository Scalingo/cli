package scalingo

import "gopkg.in/errgo.v1"

func (c *Client) SignUp(email, password string) error {
	req := &APIRequest{
		Client:   c,
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
