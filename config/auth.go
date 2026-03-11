package config

import (
	"context"
	"encoding/json"
	stderrors "errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/Scalingo/cli/config/auth"
	appio "github.com/Scalingo/cli/io"
	"github.com/Scalingo/cli/term"
	"github.com/Scalingo/go-scalingo/v11"
	scalingohttp "github.com/Scalingo/go-scalingo/v11/http"
	"github.com/Scalingo/go-utils/errors/v3"
)

var (
	ErrAuthenticationFailed = stderrors.New("authentication failed")
)

type CliAuthenticator struct{}

var (
	ErrUnauthenticated = stderrors.New("user unauthenticated")
)

func Auth(ctx context.Context) (*scalingo.User, string, error) {
	var user *scalingo.User
	var tokens string
	var err error

	if C.DisableInteractive {
		err = errors.New(ctx, "Fail to login (interactive mode disabled)")
	} else {
		for i := range 3 {
			user, tokens, err = tryAuth(ctx)
			if err == nil {
				break
			} else if errors.Is(err, io.EOF) {
				return nil, "", errors.New(ctx, "canceled by user")
			} else {
				appio.Errorf("Fail to login (%d/3): %v\n\n", i+1, err)
			}
		}
	}
	if err == ErrAuthenticationFailed {
		return nil, "", errors.New(ctx, "Forgot your password? https://auth.scalingo.com")
	}
	if err != nil {
		return nil, "", errors.Wrapf(ctx, err, "authentication failed")
	}

	fmt.Print("\n")
	appio.Statusf("Hello %s, nice to see you!\n\n", user.Username)
	err = SetCurrentUser(ctx, user, tokens)
	if err != nil {
		return nil, "", errors.Wrapf(ctx, err, "set current user")
	}

	return user, tokens, nil
}

func SetCurrentUser(ctx context.Context, user *scalingo.User, token string) error {
	authenticator := &CliAuthenticator{}
	err := authenticator.StoreAuth(ctx, user, token)
	if err != nil {
		return errors.Wrapf(ctx, err, "store credentials")
	}
	return nil
}

func (a *CliAuthenticator) StoreAuth(ctx context.Context, user *scalingo.User, token string) error {
	authConfig, err := existingAuth(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "read existing authentication config")
	}

	var c auth.ConfigPerHostV2
	err = json.Unmarshal(authConfig.AuthConfigPerHost, &c)
	if err != nil {
		fmt.Println("Auth: error while reading auth file. Recreating a new one.")
		c = make(auth.ConfigPerHostV2)
	}

	authHost, err := a.authHost(ctx)
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to get authentication service host")
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
		return errors.Wrapf(ctx, err, "fail to marshal the configuration to JSON")
	}

	authConfig.AuthConfigPerHost = json.RawMessage(buffer)
	return writeAuthFile(ctx, authConfig)
}

func (a *CliAuthenticator) LoadAuth(ctx context.Context) (*scalingo.User, *auth.UserToken, error) {
	file, err := os.OpenFile(C.AuthFile, os.O_RDONLY, 0600)
	if os.IsNotExist(err) {
		return nil, nil, ErrUnauthenticated
	}
	if err != nil {
		return nil, nil, errors.Wrapf(ctx, err, "open authentication file for read")
	}

	var authConfig auth.ConfigData
	if err := json.NewDecoder(file).Decode(&authConfig); err != nil {
		file.Close()
		return nil, nil, errors.Wrapf(ctx, err, "decode authentication file")
	}
	file.Close()

	if authConfig.AuthDataVersion != auth.ConfigVersionV2 && authConfig.AuthDataVersion != auth.ConfigVersionV21 {
		err = writeAuthFile(ctx, &authConfig)
		if err != nil {
			return nil, nil, errors.Wrapf(ctx, err, "fail to update to authv2")
		}
		return nil, nil, ErrUnauthenticated
	}

	var configPerHost auth.ConfigPerHostV2
	err = json.Unmarshal(authConfig.AuthConfigPerHost, &configPerHost)
	if err != nil {
		return nil, nil, errors.Wrapf(ctx, err, "unmarshal authentication config per host")
	}

	if authConfig.AuthDataVersion == auth.ConfigVersionV2 {
		authConfig.AuthDataVersion = auth.ConfigVersionV21
		configPerHost["auth.scalingo.com"] = configPerHost["api.scalingo.com"]
		delete(configPerHost, "api.scalingo.com")
		buffer, err := json.Marshal(&configPerHost)
		if err != nil {
			return nil, nil, errors.Wrapf(ctx, err, "migrate auth config v2.0 to v2.1")
		}
		authConfig.AuthConfigPerHost = json.RawMessage(buffer)
		err = writeAuthFile(ctx, &authConfig)
		if err != nil {
			return nil, nil, errors.Wrapf(ctx, err, "migrate auth config v2.0 to v2.1")
		}
	}

	authHost, err := a.authHost(ctx)
	if err != nil {
		return nil, nil, errors.Wrapf(ctx, err, "fail to get authentication service host")
	}

	creds, ok := configPerHost[authHost]
	if !ok || creds.User == nil {
		return nil, nil, ErrUnauthenticated
	}
	return creds.User, creds.Tokens, nil
}

func (a *CliAuthenticator) RemoveAuth(ctx context.Context) error {
	authConfig, err := existingAuth(ctx)
	if err != nil {
		return errors.Wrapf(ctx, err, "get authentication config for removal")
	}

	var c auth.ConfigPerHostV2
	err = json.Unmarshal(authConfig.AuthConfigPerHost, &c)
	if err != nil {
		return errors.Wrapf(ctx, err, "unmarshal authentication config per host")
	}

	authHost, err := a.authHost(ctx)
	if err != nil {
		return errors.Wrapf(ctx, err, "get authentication service host")
	}

	delete(c, authHost)

	buffer, err := json.Marshal(&c)
	if err != nil {
		return errors.Wrapf(ctx, err, "marshal cleaned authentication config")
	}

	authConfig.AuthConfigPerHost = json.RawMessage(buffer)
	return writeAuthFile(ctx, authConfig)
}

func (a *CliAuthenticator) authHost(ctx context.Context) (string, error) {
	u, err := url.Parse(C.ScalingoAuthURL)
	if err != nil {
		return "", errors.Wrapf(ctx, err, "parse auth URL: %v", C.ScalingoAuthURL)
	}
	return strings.Split(u.Host, ":")[0], nil
}

func tryAuth(ctx context.Context) (*scalingo.User, string, error) {
	var login string
	var err error

	for login == "" {
		appio.Infof("Username or email: ")
		_, err := fmt.Scanln(&login)
		if err != nil {
			if strings.Contains(err.Error(), "unexpected newline") {
				continue
			}
			return nil, "", errors.Wrapf(ctx, err, "read username")
		}
		login = strings.TrimRight(login, "\n")
	}

	password, err := term.Password(ctx, "       Password: ")
	if err != nil {
		return nil, "", errors.Wrapf(ctx, err, "read password")
	}
	fmt.Printf("\n")

	otpRequired := false
	retryAuth := true

	c, err := ScalingoUnauthenticatedAuthClient(ctx)
	if err != nil {
		return nil, "", errors.Wrapf(ctx, err, "fail to create an unauthenticated Scalingo client")
	}

	loginParams := scalingo.LoginParams{}
	var apiToken scalingo.Token
	for retryAuth {
		loginParams.Identifier = login
		loginParams.Password = password

		var otp string
		if otpRequired {
			appio.Infof("OTP: ")
			fmt.Scan(&otp)
			loginParams.OTP = otp
		}

		hostname, err := os.Hostname()
		if err != nil {
			return nil, "", errors.Wrapf(ctx, err, "fail to get current hostname")
		}

		apiToken, err = c.TokenCreateWithLogin(ctx, scalingo.TokenCreateParams{
			Name: "Scalingo CLI - " + hostname,
		}, loginParams)
		if err != nil {
			if !otpRequired && scalingohttp.IsOTPRequired(err) {
				otpRequired = true
			} else {
				return nil, "", errors.Wrapf(ctx, err, "fail to create API token")
			}
		} else {
			retryAuth = false
		}
	}

	client, err := ScalingoAuthClientFromToken(ctx, apiToken.Token)
	if err != nil {
		return nil, "", errors.Wrapf(ctx, err, "fail to create an authenticated Scalingo client using the API token")
	}
	userInformation, err := client.Self(ctx)
	if err != nil {
		return nil, "", errors.Wrapf(ctx, err, "fail to get account data")
	}

	return userInformation, apiToken.Token, nil
}

func writeAuthFile(ctx context.Context, authConfig *auth.ConfigData) error {
	file, err := os.OpenFile(C.AuthFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0700)
	if err != nil {
		return errors.Wrapf(ctx, err, "open authentication file for writing")
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	if err := enc.Encode(authConfig); err != nil {
		return errors.Wrapf(ctx, err, "encode authentication file")
	}
	return nil
}

func existingAuth(ctx context.Context) (*auth.ConfigData, error) {
	authConfig := auth.NewConfigData()
	content, err := os.ReadFile(C.AuthFile)
	if err == nil {
		// We don't care of the error
		json.Unmarshal(content, &authConfig)
	} else if os.IsNotExist(err) {
		config := make(auth.ConfigPerHostV2)
		buffer, err := json.Marshal(&config)
		if err != nil {
			return nil, errors.Wrapf(ctx, err, "encode non-existing authentication file")
		}
		authConfig.AuthConfigPerHost = json.RawMessage(buffer)
	} else {
		return nil, errors.Wrapf(ctx, err, "read existing authentication file")
	}
	return authConfig, nil
}
