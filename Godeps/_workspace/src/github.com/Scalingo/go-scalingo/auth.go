package scalingo

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Scalingo/cli/Godeps/_workspace/src/gopkg.in/errgo.v1"
)

var (
	ErrLoginAborted  = errors.New("canceled by user.")
	ApiAuthenticator Authenticator
	ApiUrl           string
	ApiVersion       string
)

type Authenticator interface {
	LoadAuth() (*User, error)
	StoreAuth(user *User) error
	RemoveAuth() error
}

type LoginError struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type LoginResponse struct {
	AuthenticationToken string `json:"authentication_token"`
	User                *User  `json:"user"`
}

func (err *LoginError) Error() string {
	return err.Message
}

func AuthFromConfig() (*User, error) {
	user, err := ApiAuthenticator.LoadAuth()
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}
	return user, nil
}

func AuthUser(login, passwd string) (*User, error) {
	res, err := Login(login, passwd)
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}
	defer res.Body.Close()

	if res.StatusCode == 401 {
		var lErr *LoginError
		err = json.NewDecoder(res.Body).Decode(&lErr)
		if err != nil {
			return nil, errgo.Mask(err, errgo.Any)
		}
		return nil, errgo.Mask(lErr, errgo.Any)
	}
	if res.StatusCode != 201 {
		return nil, fmt.Errorf("%s %v invalid status %s", res.Request.Method, res.Request.URL, res.Status)
	}

	loginRes := &LoginResponse{}

	err = ApiAuthenticator.StoreAuth(loginRes.User)
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}

	err = json.NewDecoder(res.Body).Decode(&loginRes)
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}

	loginRes.User.AuthenticationToken = loginRes.AuthenticationToken
	return loginRes.User, nil
}
