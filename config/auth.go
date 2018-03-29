package config

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/Scalingo/cli/config/auth"
	appio "github.com/Scalingo/cli/io"
	"github.com/Scalingo/cli/term"
	"github.com/Scalingo/go-scalingo"
	"github.com/pkg/errors"
	"gopkg.in/errgo.v1"
)

var (
	ErrAuthenticationFailed = errors.New("authentication failed")
)

type CliAuthenticator struct{}

var (
	Authenticator      = &CliAuthenticator{}
	ErrUnauthenticated = errgo.New("user unauthenticated")
)

func Auth() (*scalingo.User, *scalingo.APITokenGenerator, error) {
	var user *scalingo.User
	var tokens *scalingo.APITokenGenerator
	var err error

	if C.DisableInteractive {
		err = errors.New("Fail to login (interactive mode disabled)")
	} else {
		for i := 0; i < 3; i++ {
			user, tokens, err = tryAuth()
			if err == nil {
				break
			} else if errgo.Cause(err) == io.EOF {
				return nil, nil, errors.New("canceled by user")
			} else {
				appio.Errorf("Fail to login (%d/3): %v\n\n", i+1, err)
			}
		}
	}
	if err == ErrAuthenticationFailed {
		return nil, nil, errors.New("Forgot your password? https://auth.scalingo.com")
	}
	if err != nil {
		return nil, nil, errgo.Mask(err, errgo.Any)
	}

	fmt.Print("\n")
	appio.Statusf("Hello %s, nice to see you!\n\n", user.Username)
	err = SetCurrentUser(user, tokens)
	if err != nil {
		return nil, nil, errgo.Mask(err, errgo.Any)
	}

	return user, tokens, nil
}

func SetCurrentUser(user *scalingo.User, generator *scalingo.APITokenGenerator) error {
	C.TokenGenerator = generator
	AuthenticatedUser = user
	err := Authenticator.StoreAuth(user, generator)
	if err != nil {
		return errgo.Notef(err, "fail to store user credentials")
	}
	return nil
}

func (a *CliAuthenticator) StoreAuth(user *scalingo.User, generator *scalingo.APITokenGenerator) error {
	authConfig, err := existingAuth()
	if err != nil {
		return errgo.Mask(err)
	}

	var c auth.ConfigPerHostV2
	err = json.Unmarshal(authConfig.AuthConfigPerHost, &c)
	if err != nil {
		fmt.Println("Auth: error while reading auth file. Recreating a new one.")
		c = make(auth.ConfigPerHostV2)
	}

	c[C.apiHost] = &auth.CredentialsData{
		TokenGenerator: generator,
		User:           user,
	}

	authConfig.LastUpdate = time.Now()
	authConfig.AuthDataVersion = auth.ConfigVersionV2

	buffer, err := json.Marshal(&c)
	if err != nil {
		return errgo.Mask(err)
	}

	authConfig.AuthConfigPerHost = json.RawMessage(buffer)
	return writeAuthFile(authConfig)
}

func (a *CliAuthenticator) LoadAuth() (*scalingo.User, *scalingo.APITokenGenerator, error) {
	file, err := os.OpenFile(C.AuthFile, os.O_RDONLY, 0600)
	if os.IsNotExist(err) {
		return nil, nil, ErrUnauthenticated
	}
	if err != nil {
		return nil, nil, errgo.Mask(err, errgo.Any)
	}

	var authConfig auth.ConfigData
	if err := json.NewDecoder(file).Decode(&authConfig); err != nil {
		file.Close()
		return nil, nil, errgo.Mask(err, errgo.Any)
	}
	file.Close()

	if authConfig.AuthDataVersion != auth.ConfigVersionV2 {
		err = writeAuthFile(&authConfig)
		if err != nil {
			return nil, nil, errgo.NoteMask(err, "fail to update to authv2", errgo.Any)
		}
		return nil, nil, ErrUnauthenticated
	}

	var configPerHost auth.ConfigPerHostV2
	err = json.Unmarshal(authConfig.AuthConfigPerHost, &configPerHost)
	if err != nil {
		return nil, nil, errgo.Mask(err)
	}

	if creds, ok := configPerHost[C.apiHost]; !ok {
		return nil, nil, ErrUnauthenticated
	} else {
		if creds == nil {
			return nil, nil, ErrUnauthenticated
		}
		return creds.User, creds.TokenGenerator, nil
	}
}

func (a *CliAuthenticator) RemoveAuth() error {
	authConfig, err := existingAuth()
	if err != nil {
		return errgo.Mask(err)
	}

	var c auth.ConfigPerHostV2
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

func tryAuth() (*scalingo.User, *scalingo.APITokenGenerator, error) {
	var login string
	var err error

	for login == "" {
		appio.Infof("Username or email: ")
		_, err := fmt.Scanln(&login)
		if err != nil {
			if strings.Contains(err.Error(), "unexpected newline") {
				continue
			}
			return nil, nil, errgo.Mask(err, errgo.Any)
		}
		login = strings.TrimRight(login, "\n")
	}

	password, err := term.Password("       Password: ")
	if err != nil {
		return nil, nil, errgo.Mask(err, errgo.Any)
	}

	otpRequired := false
	retryAuth := true

	c := ScalingoUnauthenticatedClient()

	loginParams := scalingo.LoginParams{}
	var apiToken scalingo.Token
	for retryAuth {
		var otp string
		if otpRequired {
			otp, err = term.Password("       OTP: ")
			if err != nil {
				return nil, nil, errgo.NoteMask(err, "fail to get otp", errgo.Any)
			}
		}

		loginParams.Identifier = login
		loginParams.Password = password

		if otpRequired {
			loginParams.OTP = otp
		}

		hostname, err := os.Hostname()
		if err != nil {
			return nil, nil, errgo.Notef(err, "fail to get current hostname")
		}

		apiToken, err = c.TokenCreateWithLogin(scalingo.TokenCreateParams{
			Name: fmt.Sprintf("Scalingo CLI - %s", hostname),
		}, loginParams)
		if err != nil {
			if !otpRequired && errgo.Cause(err) == scalingo.ErrOTPRequired {
				otpRequired = true
			} else {
				return nil, nil, errgo.NoteMask(err, "fail to create API token", errgo.Any)
			}
		} else {
			retryAuth = false
		}
	}

	generator := c.GetAPITokenGenerator(apiToken.Token)
	C.TokenGenerator = generator

	client := ScalingoClient()
	userInformation, err := client.Self()
	if err != nil {
		return nil, nil, errgo.NoteMask(err, "fail to get account data", errgo.Any)
	}

	return userInformation, generator, nil
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
		config := make(auth.ConfigPerHostV2)
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
