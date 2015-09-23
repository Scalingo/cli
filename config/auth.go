package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"

	"github.com/Scalingo/cli/Godeps/_workspace/src/gopkg.in/errgo.v1"
	"github.com/Scalingo/cli/users"
)

type CliAuthenticator struct{}

type AuthConfigData struct {
	LastUpdate        time.Time              `json:"last_update"`
	AuthConfigPerHost map[string]*users.User `json:"auth_config_data"`
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
		return nil, nil
	} else {
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
