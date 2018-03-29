package session

import (
	"fmt"
	"os"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/debug"
	"github.com/Scalingo/cli/io"
	netssh "github.com/Scalingo/cli/net/ssh"
	"github.com/Scalingo/go-scalingo"
	"github.com/pkg/errors"
	"gopkg.in/errgo.v1"
)

type LoginOpts struct {
	APIToken    string
	Ssh         bool
	SshIdentity string
}

func Login(opts LoginOpts) error {
	if opts.SshIdentity == "" {
		opts.SshIdentity = "ssh-agent"
	}

	if opts.APIToken != "" {
		return loginWithToken(opts.APIToken)
	}

	if config.AuthenticatedUser != nil {
		io.Statusf("You are already logged as %s (%s)!\n", config.AuthenticatedUser.Email, config.AuthenticatedUser.Username)
		return nil
	}
	io.Status("Currently unauthenticated")
	io.Info("Trying login with SSH…")
	err := loginWithSsh(opts.SshIdentity)
	if err != nil {
		config.C.Logger.Printf("SSH connection failed: %+v\n", err)
		io.Error("SSH connection failed.")
		if opts.Ssh {
			if errors.Cause(err) == netssh.ErrNoAuthSucceed {
				return errors.Wrap(err, "please use the flag '--ssh-identity /path/to/private/key' to specify your private key")
			}
			if err != nil {
				return errors.Wrap(err, "fail to login wish SSH")
			}
		} else {
			io.Info("Trying login with user/password:\n")
			return loginWithUserAndPassword()
		}
	}

	return nil
}

func loginWithUserAndPassword() error {
	_, _, err := config.Auth()
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}
	return nil
}

func loginWithToken(token string) error {
	err := finalizeLogin(token)
	if err != nil {
		return errgo.Notef(err, "token invalid")
	}
	return nil
}

func loginWithSsh(identity string) error {
	debug.Println("Login through SSH, identity:", identity)
	client, _, err := netssh.Connect(identity)
	if err != nil {
		return errors.Wrap(err, "fail to connect to SSH server")
	}
	channel, reqs, err := client.OpenChannel("session", []byte{})
	if err != nil {
		return errgo.Notef(err, "fail to open SSH channel")
	}

	defer client.Close()

	_, err = channel.SendRequest("auth.v2@scalingo.com", false, []byte{})
	if err != nil {
		return errgo.Notef(err, "SSH authentication request fails")
	}
	req := <-reqs
	if req == nil {
		return errgo.Newf("invalid response from auth request")
	}
	if req.Type != "auth.v2@scalingo.com" {
		return errgo.Newf("invalid response from SSH server, type is %v", req.Type)
	}
	payload := req.Payload

	if len(payload) == 0 {
		return errgo.Newf("invalid response from SSH server")
	}

	hostname, err := os.Hostname()
	if err != nil {
		return errgo.Notef(err, "fail to get current hostname")
	}

	c := config.ScalingoUnauthenticatedClient()
	token, err := c.TokenCreateWithLogin(scalingo.TokenCreateParams{
		Name: fmt.Sprintf("Scalingo CLI - %s", hostname),
	}, scalingo.LoginParams{
		JWT: string(payload),
	})
	if err != nil {
		return errgo.NoteMask(err, "fail to create API token", errgo.Any)
	}

	err = finalizeLogin(token.Token)
	if err != nil {
		return errgo.NoteMask(err, "fail to finalize login", errgo.Any)
	}
	return nil
}

func finalizeLogin(token string) error {
	c := config.ScalingoUnauthenticatedClient()
	generator := c.GetAPITokenGenerator(token)
	config.C.TokenGenerator = generator

	c = config.ScalingoClient()
	user, err := c.Self()
	if err != nil {
		return errgo.Mask(err)
	}

	io.Statusf("Hello %s, nice to see you!\n", user.Username)

	err = config.SetCurrentUser(user, generator)
	if err != nil {
		return errgo.Mask(err)
	}
	return nil

}
