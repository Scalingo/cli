package config

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config/auth"
	appio "github.com/Scalingo/cli/io"
	"github.com/Scalingo/cli/term"
	"github.com/Scalingo/go-scalingo/v4"
	scalingoerrors "github.com/Scalingo/go-utils/errors"
)

var (
	ErrAuthenticationFailed = errors.New("authentication failed")
)

type CliAuthenticator struct{}

var (
	ErrUnauthenticated = errgo.New("user unauthenticated")
)

func Auth() (*scalingo.User, string, error) {
	var user *scalingo.User
	var tokens string
	var err error

	if C.DisableInteractive {
		err = errors.New("Fail to login (interactive mode disabled)")
	} else {
		for i := 0; i < 3; i++ {
			user, tokens, err = tryAuth()
			if err == nil {
				break
			} else if scalingoerrors.ErrgoRoot(err) == io.EOF {
				return nil, "", errors.New("canceled by user")
			} else {
				appio.Errorf("Fail to login (%d/3): %v\n\n", i+1, err)
			}
		}
	}
	if err == ErrAuthenticationFailed {
		return nil, "", errors.New("Forgot your password? https://auth.scalingo.com")
	}
	if err != nil {
		return nil, "", errgo.Mask(err, errgo.Any)
	}

	fmt.Print("\n")
	appio.Statusf("Hello %s, nice to see you!\n\n", user.Username)
	err = SetCurrentUser(user, tokens)
	if err != nil {
		return nil, "", errgo.Mask(err, errgo.Any)
	}

	return user, tokens, nil
}

func SetCurrentUser(user *scalingo.User, token string) error {
	authenticator := &CliAuthenticator{}
	err := authenticator.StoreAuth(user, token)
	if err != nil {
		return errgo.Notef(err, "fail to store user credentials")
	}
	return nil
}

func (a *CliAuthenticator) StoreAuth(user *scalingo.User, token string) error {
	authConfig, err := existingAuth()
	if err != nil {
		return err
	}

	var c auth.ConfigPerHostV2
	err = json.Unmarshal(authConfig.AuthConfigPerHost, &c)
	if err != nil {
		fmt.Println("Auth: error while reading auth file. Recreating a new one.")
		c = make(auth.ConfigPerHostV2)
	}

	authHost, err := a.authHost()
	if err != nil {
		return errgo.Notef(err, "fail to get authentication service host")
	}

	c[authHost] = auth.CredentialsData{
		Tokens: &auth.UserToken{
			Token: token,
		},
		User: user,
	}

	authConfig.LastUpdate = time.Now()
	authConfig.AuthDataVersion = auth.ConfigVersionV21

	buffer, err := json.Marshal(&c)
	if err != nil {
		return errgo.Notef(err, "fail to marshal the configuration to JSON")
	}

	authConfig.AuthConfigPerHost = json.RawMessage(buffer)
	return writeAuthFile(authConfig)
}

func (a *CliAuthenticator) LoadAuth() (*scalingo.User, *auth.UserToken, error) {
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

	if authConfig.AuthDataVersion != auth.ConfigVersionV2 && authConfig.AuthDataVersion != auth.ConfigVersionV21 {
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

	if authConfig.AuthDataVersion == auth.ConfigVersionV2 {
		authConfig.AuthDataVersion = auth.ConfigVersionV21
		configPerHost["auth.scalingo.com"] = configPerHost["api.scalingo.com"]
		delete(configPerHost, "api.scalingo.com")
		buffer, err := json.Marshal(&configPerHost)
		if err != nil {
			return nil, nil, errgo.Notef(err, "Fail to migrate auth config v2.0 to v2.1")
		}
		authConfig.AuthConfigPerHost = json.RawMessage(buffer)
		err = writeAuthFile(&authConfig)
		if err != nil {
			return nil, nil, errgo.Notef(err, "Fail to migrate auth config v2.0 to v2.1")
		}
	}

	authHost, err := a.authHost()
	if err != nil {
		return nil, nil, errgo.Notef(err, "fail to get authentication service host")
	}

	creds, ok := configPerHost[authHost]
	if !ok || creds.User == nil {
		return nil, nil, ErrUnauthenticated
	}
	return creds.User, creds.Tokens, nil
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

	authHost, err := a.authHost()
	if err != nil {
		return errgo.Notef(err, "fail to get authentication service host")
	}
	if _, ok := c[authHost]; ok {
		delete(c, authHost)
	}

	buffer, err := json.Marshal(&c)
	if err != nil {
		return errgo.Mask(err)
	}

	authConfig.AuthConfigPerHost = json.RawMessage(buffer)
	return writeAuthFile(authConfig)
}

func (a *CliAuthenticator) authHost() (string, error) {
	u, err := url.Parse(C.ScalingoAuthUrl)
	if err != nil {
		return "", errgo.Notef(err, "fail to parse auth URL: %v", C.ScalingoAuthUrl)
	}
	return strings.Split(u.Host, ":")[0], nil
}

func tryAuth() (*scalingo.User, string, error) {
	var login string
	var err error

	for login == "" {
		appio.Infof("Username or email: ")
		_, err := fmt.Scanln(&login)
		if err != nil {
			if strings.Contains(err.Error(), "unexpected newline") {
				continue
			}
			return nil, "", errgo.Mask(err, errgo.Any)
		}
		login = strings.TrimRight(login, "\n")
	}

	password, err := term.Password("       Password: ")
	if err != nil {
		return nil, "", errgo.Mask(err, errgo.Any)
	}

	otpRequired := false
	retryAuth := true

	c, err := ScalingoUnauthenticatedAuthClient()
	if err != nil {
		return nil, "", errgo.Notef(err, "fail to create an unauthenticated Scalingo client")
	}

	loginParams := scalingo.LoginParams{}
	var apiToken scalingo.Token
	for retryAuth {
		var otp string
		if otpRequired {
			otp, err = term.Password("       OTP: ")
			if err != nil {
				return nil, "", errgo.NoteMask(err, "fail to get otp", errgo.Any)
			}
		}

		loginParams.Identifier = login
		loginParams.Password = password

		if otpRequired {
			loginParams.OTP = otp
		}

		hostname, err := os.Hostname()
		if err != nil {
			return nil, "", errgo.Notef(err, "fail to get current hostname")
		}

		apiToken, err = c.TokenCreateWithLogin(scalingo.TokenCreateParams{
			Name: fmt.Sprintf("Scalingo CLI - %s", hostname),
		}, loginParams)
		if err != nil {
			if !otpRequired && scalingoerrors.ErrgoRoot(err) == scalingo.ErrOTPRequired {
				otpRequired = true
			} else {
				return nil, "", errgo.NoteMask(err, "fail to create API token", errgo.Any)
			}
		} else {
			retryAuth = false
		}
	}

	client, err := ScalingoAuthClientFromToken(apiToken.Token)
	if err != nil {
		return nil, "", errgo.Notef(err, "fail to create an authenticated Scalingo client using the API token")
	}
	userInformation, err := client.Self()
	if err != nil {
		return nil, "", errgo.Notef(err, "fail to get account data")
	}

	return userInformation, apiToken.Token, nil
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
	content, err := os.ReadFile(C.AuthFile)
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
