package auth

import (
	"appsdeck/api"
	"code.google.com/p/gopass"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func Authenticate() (*AuthConfig, error) {
	fmt.Print("Email : ")
	var email string
	_, err := fmt.Scanln(&email)
	if err != nil {
		return nil, err
	}

	password, err := gopass.GetPass("Password : ")
	if err != nil {
		return nil, err
	}

	res, err := api.Login(email, password)
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
		return nil, fmt.Errorf("Wrong email or password.")
	}

	return nil, fmt.Errorf("Wrong answer from server : %s", res.Status)
}
