package scalingo

import (
	"encoding/json"

	"github.com/Scalingo/cli/Godeps/_workspace/src/github.com/Scalingo/go-scalingo/users"
	"github.com/Scalingo/cli/Godeps/_workspace/src/gopkg.in/errgo.v1"
)

type SelfResults struct {
	User *users.User `json:"user"`
}

func Self() (*users.User, error) {
	req := &APIRequest{
		Endpoint: "/users/self",
	}
	res, err := req.Do()
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}
	defer res.Body.Close()
	var u *users.User
	err = json.NewDecoder(res.Body).Decode(&u)
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}
	return u, nil
}
