package api

import (
	"encoding/json"
	"fmt"
	"strings"

	"code.google.com/p/gopass"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/users"
)

type LoginError struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func (err *LoginError) Error() string {
	return err.Message
}

func AuthFromConfig() (*users.User, error) {
	user, err := config.LoadAuth()
	if err != nil {
		return nil, err
	}
	return user, nil
}

func Auth() (*users.User, error) {
	var user *users.User
	var err error
	for i := 0; i < 3; i++ {
		user, err = tryAuth()
		if err == nil {
			break
		} else {
			fmt.Printf("Fail to login (%d/3): %v\n", i+1, err)
		}
	}
	if err != nil {
		return nil, err
	}

	fmt.Printf("Hello %s, nice to see you !\n\n", user.Username)
	err = config.StoreAuth(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func tryAuth() (*users.User, error) {
	var login string
	for login == "" {
		fmt.Print("Username or email: ")
		_, err := fmt.Scanln(&login)
		if err != nil {
			if strings.Contains(err.Error(), "unexpected newline") {
				continue
			}
			return nil, err
		}
	}

	password, err := gopass.GetPass("Password: ")
	if err != nil {
		return nil, err
	}

	user, err := AuthUser(login, password)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func AuthUser(login, passwd string) (*users.User, error) {
	res, err := Login(login, passwd)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 401 {
		var lErr *LoginError
		err = json.NewDecoder(res.Body).Decode(&lErr)
		if err != nil {
			return nil, err
		}
		return nil, lErr
	}
	if res.StatusCode != 201 {
		return nil, fmt.Errorf("%s %v invalid status %s", res.Request.Method, res.Request.URL, res.Status)
	}

	var tokenMap map[string]string

	err = json.NewDecoder(res.Body).Decode(&tokenMap)
	if err != nil {
		return nil, err
	}

	token := tokenMap["authentication_token"]

	params := map[string]interface{}{
		"auth":     false,
		"token":    token,
		"method":   "GET",
		"endpoint": "/users/self",
	}

	res, err = Do(params)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("%s %v invalid status %s", res.Request.Method, res.Request.URL, res.Status)
	}

	var selfRes *SelfResults
	err = json.NewDecoder(res.Body).Decode(&selfRes)
	if err != nil {
		return nil, err
	}
	selfRes.User.AuthToken = token

	return selfRes.User, nil
}
