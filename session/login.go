package session

import (
	"context"
	"fmt"
	"os"

	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	netssh "github.com/Scalingo/cli/net/ssh"
	"github.com/Scalingo/go-scalingo/v7"
	"github.com/Scalingo/go-scalingo/v7/debug"
	"github.com/Scalingo/go-utils/errors/v2"
)

type LoginOpts struct {
	APIToken     string
	PasswordOnly bool
	SSH          bool
	SSHIdentity  string
}

func Login(ctx context.Context, opts LoginOpts) error {
	if opts.SSHIdentity == "" {
		opts.SSHIdentity = "ssh-agent"
	}

	if opts.APIToken != "" {
		return loginWithToken(ctx, opts.APIToken)
	}

	if !opts.PasswordOnly {
		io.Info("Trying login with SSH…")
		err := loginWithSSH(ctx, opts.SSHIdentity)
		if err != nil {
			config.C.Logger.Printf("SSH connection failed: %+v\n", err)
			io.Error("SSH connection failed.")
			if opts.SSH {
				if errors.RootCause(err) == netssh.ErrNoAuthSucceed {
					return errors.Wrapf(ctx, err, "please use the flag '--ssh-identity /path/to/private/key' to specify your private key")
				}
				return errors.Wrapf(ctx, err, "fail to login with SSH")
			}
		} else {
			return nil
		}
	}

	io.Info("Trying login with user/password:\n")
	return loginWithUserAndPassword(ctx)
}

func loginWithUserAndPassword(ctx context.Context) error {
	_, _, err := config.Auth(ctx)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}
	return nil
}

func loginWithToken(ctx context.Context, token string) error {
	err := finalizeLogin(ctx, token)
	if err != nil {
		return errors.Wrapf(ctx, err, "token invalid")
	}
	return nil
}

func loginWithSSH(ctx context.Context, identity string) error {
	host := config.C.ScalingoSSHHost
	if host == "" {
		regions, err := config.EnsureRegionsCache(ctx, config.C, config.GetRegionOpts{
			SkipAuth: true,
		})
		if err != nil {
			return errors.Wrapf(ctx, err, "fail to ensure region cache")
		}

		defaultRegion, err := regions.Default()
		if err != nil {
			return errors.Wrapf(ctx, err, "fail to find default region")
		}

		host = defaultRegion.SSH
	}

	debug.Printf("Login through SSH, Host: %s Identity:%s\n", host, identity)
	client, _, err := netssh.Connect(ctx, netssh.ConnectOpts{
		Host:     host,
		Identity: identity,
	})
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to connect to SSH server")
	}
	channel, reqs, err := client.OpenChannel("session", []byte{})
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to open SSH channel")
	}

	defer client.Close()

	_, err = channel.SendRequest("auth.v2@scalingo.com", false, []byte{})
	if err != nil {
		return errors.Wrapf(ctx, err, "SSH authentication request fails")
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
		return errors.Wrapf(ctx, err, "fail to get current hostname")
	}

	c, err := config.ScalingoUnauthenticatedAuthClient(ctx)
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to create an unauthenticated Scalingo client")
	}
	token, err := c.TokenCreateWithLogin(ctx, scalingo.TokenCreateParams{
		Name: fmt.Sprintf("Scalingo CLI - %s", hostname),
	}, scalingo.LoginParams{
		JWT: string(payload),
	})
	if err != nil {
		return errgo.NoteMask(err, "fail to create API token", errgo.Any)
	}

	err = finalizeLogin(ctx, token.Token)
	if err != nil {
		return errgo.NoteMask(err, "fail to finalize login", errgo.Any)
	}
	return nil
}

func finalizeLogin(ctx context.Context, token string) error {
	c, err := config.ScalingoAuthClientFromToken(ctx, token)
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to create an authenticated Scalingo client using the API token")
	}
	user, err := c.Self(ctx)
	if err != nil {
		return errgo.Mask(err)
	}

	io.Statusf("Hello %s, nice to see you!\n", user.Username)

	err = config.SetCurrentUser(ctx, user, token)
	if err != nil {
		return errgo.Mask(err)
	}
	return nil
}
