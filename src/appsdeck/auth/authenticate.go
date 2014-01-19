package auth

import (
	"code.google.com/p/gopass"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func Authenticate() (*AuthConfig, error) {
	fmt.Print("Username or email: ")
	var login string
	_, err := fmt.Scanln(&login)
	if err != nil {
		return nil, err
	}

	password, err := gopass.GetPass("Password: ")
	if err != nil {
		return nil, err
	}

	res, err := loginUser(login, password)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusOK:
		var authConfig AuthConfig
		if err := json.Unmarshal(body, &authConfig); err != nil {
			return nil, err
		}
		return &authConfig, StoreAuth(&authConfig)
	case http.StatusUnauthorized:
		return nil, fmt.Errorf("Wrong identifier or password.")
	}

	return nil, fmt.Errorf("Wrong answer from server : %s", res.Status)
}
