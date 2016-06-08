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

	"github.com/Scalingo/go-scalingo"
	"gopkg.in/errgo.v1"
	"github.com/Scalingo/cli/config/auth"
	"github.com/Scalingo/cli/debug"
	"github.com/Scalingo/cli/term"
)

type CliAuthenticator struct{}

var Authenticator = &CliAuthenticator{}

func Auth() (*scalingo.User, error) {
	var user *scalingo.User
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

	fmt.Printf("Hello %s, nice to see you!\n\n", user.Username)
	err = Authenticator.StoreAuth(user)
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}

	return user, nil
}

func (a *CliAuthenticator) StoreAuth(user *scalingo.User) error {
	authConfig, err := existingAuth()
	if err != nil {
		return errgo.Mask(err)
	}

	var c auth.ConfigPerHostV1
	err = json.Unmarshal(authConfig.AuthConfigPerHost, &c)
	if err != nil {
		return errgo.Mask(err)
	}

	c[C.apiHost] = user
	authConfig.LastUpdate = time.Now()

	buffer, err := json.Marshal(&c)
	if err != nil {
		return errgo.Mask(err)
	}

	authConfig.AuthConfigPerHost = json.RawMessage(buffer)
	return writeAuthFile(authConfig)
}

func (a *CliAuthenticator) LoadAuth() (*scalingo.User, error) {
	file, err := os.OpenFile(C.AuthFile, os.O_RDONLY, 0600)
	if os.IsNotExist(err) {
		return Auth()
	}
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}

	var authConfig auth.ConfigData
	if err := json.NewDecoder(file).Decode(&authConfig); err != nil {
		file.Close()
		return nil, errgo.Mask(err, errgo.Any)
	}
	file.Close()

	if authConfig.AuthDataVersion == "" {
		debug.Println("auth config should be updated")
		err = authConfig.MigrateToV1()
		if err != nil {
			return nil, errgo.Mask(err)
		}
		err = writeAuthFile(&authConfig)
		if err != nil {
			return nil, errgo.Mask(err)
		}
		debug.Println("auth config has been updated to V1")
	}

	var configPerHost auth.ConfigPerHostV1
	err = json.Unmarshal(authConfig.AuthConfigPerHost, &configPerHost)
	if err != nil {
		return nil, errgo.Mask(err)
	}

	if user, ok := configPerHost[C.apiHost]; !ok {
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

	var c auth.ConfigPerHostV1
	err = json.Unmarshal(authConfig.AuthConfigPerHost, &c)
	if err != nil {
		return errgo.Mask(err)
	}

	if _, ok := c[C.apiHost]; ok {
		delete(c, C.apiHost)
	}

	buffer, err := json.Marshal(&c)
	if err != nil {
		return errgo.Mask(err)
	}

	authConfig.AuthConfigPerHost = json.RawMessage(buffer)
	return writeAuthFile(authConfig)
}

func tryAuth() (*scalingo.User, error) {
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

	c := ScalingoUnauthenticatedClient()
	res, err := c.Login(login, password)
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}
	return res.User, nil
}

func writeAuthFile(authConfig *auth.ConfigData) error {
	file, err := os.OpenFile(C.AuthFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0700)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	if err := enc.Encode(authConfig); err != nil {
		return errgo.Mask(err, errgo.Any)
	}
	return nil

}

func existingAuth() (*auth.ConfigData, error) {
	authConfig := auth.NewConfigData()
	content, err := ioutil.ReadFile(C.AuthFile)
	if err == nil {
		// We don't care of the error
		json.Unmarshal(content, &authConfig)
	} else if os.IsNotExist(err) {
		config := make(auth.ConfigPerHostV1)
		buffer, err := json.Marshal(&config)
		if err != nil {
			return nil, errgo.Mask(err)
		}
		authConfig.AuthConfigPerHost = json.RawMessage(buffer)
	} else {
		return nil, errgo.Mask(err)
	}
	return authConfig, nil
}
