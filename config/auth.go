package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/Scalingo/cli/Godeps/_workspace/src/github.com/Scalingo/go-scalingo"
	"github.com/Scalingo/cli/Godeps/_workspace/src/github.com/Scalingo/go-scalingo/users"
	"github.com/Scalingo/cli/Godeps/_workspace/src/gopkg.in/errgo.v1"
	"github.com/Scalingo/cli/term"
)

type CliAuthenticator struct{}

type AuthConfigData struct {
	LastUpdate        time.Time              `json:"last_update"`
	AuthConfigPerHost map[string]*users.User `json:"auth_config_data"`
}

var Authenticator = &CliAuthenticator{}

func Auth() (*users.User, error) {
	var user *users.User
	var err error

	fmt.Fprintln(os.Stderr, "You need to be authenticated to use Scalingo client.\nNo account ? → https://scalingo.com")
	if C.DisableInteractive {
		err = errors.New("Fail to login (interactive mode disabled)")
	} else {
		for i := 0; i < 3; i++ {
			user, err = tryAuth()
			if err == nil {
				break
			} else if errgo.Cause(err) == io.EOF {
				return nil, errors.New("canceled by user")
			} else {
				fmt.Printf("Fail to login (%d/3): %v\n", i+1, err)
			}
		}
	}
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}

	fmt.Printf("Hello %s, nice to see you !\n\n", user.Username)
	err = Authenticator.StoreAuth(user)
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}

	return user, nil
}

func (a *CliAuthenticator) StoreAuth(user *users.User) error {
	authConfig, err := existingAuth()
	if err != nil {
		return errgo.Mask(err)
	}

	authConfig.AuthConfigPerHost[C.apiHost] = user
	authConfig.LastUpdate = time.Now()

	return writeAuthFile(authConfig)
}

func (a *CliAuthenticator) LoadAuth() (*users.User, error) {
	file, err := os.OpenFile(C.AuthFile, os.O_RDONLY, 0644)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}
	defer file.Close()

	var authConfig AuthConfigData
	dec := json.NewDecoder(file)
	if err := dec.Decode(&authConfig); err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}

	if user, ok := authConfig.AuthConfigPerHost[C.apiHost]; !ok {
		return Auth()
	} else {
		if user == nil {
			return Auth()
		}
		return user, nil
	}
}

func (a *CliAuthenticator) RemoveAuth() error {
	authConfig, err := existingAuth()
	if err != nil {
		return errgo.Mask(err)
	}
	if _, ok := authConfig.AuthConfigPerHost[C.apiHost]; ok {
		delete(authConfig.AuthConfigPerHost, C.apiHost)
	}

	return writeAuthFile(authConfig)
}

func tryAuth() (*users.User, error) {
	var login string
	var err error

	for login == "" {
		fmt.Print("Username or email: ")
		_, err := fmt.Scanln(&login)
		if err != nil {
			if strings.Contains(err.Error(), "unexpected newline") {
				continue
			}
			return nil, errgo.Mask(err, errgo.Any)
		}
		login = strings.TrimRight(login, "\n")
	}

	password, err := term.Password("Password: ")
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}

	user, err := scalingo.AuthUser(login, password)
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}

	return user, nil
}

func writeAuthFile(auth *AuthConfigData) error {
	file, err := os.OpenFile(C.AuthFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0700)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	if err := enc.Encode(auth); err != nil {
		return errgo.Mask(err, errgo.Any)
	}
	return nil

}

func existingAuth() (*AuthConfigData, error) {
	authConfig := &AuthConfigData{}
	content, err := ioutil.ReadFile(C.AuthFile)
	if err == nil {
		// We don't care of the error
		json.Unmarshal(content, &authConfig)
	} else if os.IsNotExist(err) {
		authConfig.AuthConfigPerHost = make(map[string]*users.User)
	} else {
		return nil, errgo.Mask(err)
	}
	return authConfig, nil
}
