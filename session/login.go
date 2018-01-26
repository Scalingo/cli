package session

import (
	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/debug"
	"github.com/Scalingo/cli/io"
	netssh "github.com/Scalingo/cli/net/ssh"
	"github.com/Scalingo/go-scalingo"
	"github.com/pkg/errors"
	"gopkg.in/errgo.v1"
)

type LoginOpts struct {
	ApiKey      string
	Ssh         bool
	SshIdentity string
}

func Login(opts LoginOpts) error {
	if opts.SshIdentity == "" {
		opts.SshIdentity = "ssh-agent"
	}

	if opts.ApiKey != "" {
		return loginWithApiKey(opts.ApiKey)
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
	_, err := config.Auth()
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}
	return nil
}

func loginWithApiKey(apiKey string) error {
	c := config.ScalingoUnauthenticatedClient()
	c.TokenGenerator = scalingo.NewStaticTokenGenerator(apiKey)
	user, err := c.Self()
	if err != nil {
		return errgo.Mask(err)
	}

	io.Statusf("Hello %s, nice to see you!\n", user.Username)

	user.AuthenticationToken = apiKey
	err = config.SetCurrentUser(user)
	if err != nil {
		return errgo.Mask(err)
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

	_, err = channel.SendRequest("auth@scalingo.com", false, []byte{})
	if err != nil {
		return errgo.Notef(err, "SSH authentication request fails")
	}
	req := <-reqs
	if req == nil {
		return errgo.Newf("invalid response from auth request")
	}
	if req.Type != "auth@scalingo.com" {
		return errgo.Newf("invalid response from SSH server, type is %v", req.Type)
	}
	payload := req.Payload

	if len(payload) == 0 {
		return errgo.Newf("invalid response from SSH server")
	}
	return loginWithApiKey(string(payload))
}
