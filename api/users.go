package api

import (
	"encoding/json"

	"github.com/Scalingo/cli/users"
)

type SelfResults struct {
	User *users.User `json:"user"`
}

func Self() (*users.User, error) {
	req := map[string]interface{}{
		"method":   "GET",
		"endpoint": "/users/self",
	}
	res, err := Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var u *users.User
	err = json.NewDecoder(res.Body).Decode(&u)
	if err != nil {
		return nil, err
	}
	return u, nil
}
