package session

import (
	"fmt"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/debug"
	netssh "github.com/Scalingo/cli/net/ssh"
	"gopkg.in/errgo.v1"
)

type LoginOpts struct {
	ApiKey      string
	Ssh         bool
	SshIdentity string
}

func Login(opts LoginOpts) error {
	if opts.ApiKey != "" && opts.Ssh {
		return errgo.New("only use --api-key or --ssh")
	}

	if opts.SshIdentity != "ssh-agent" && !opts.Ssh {
		return errgo.New("you can't use --ssh-identify without having --ssh")
	}

	if opts.ApiKey != "" {
		return loginWithApiKey(opts.ApiKey)
	}

	if opts.Ssh {
		return loginWithSsh(opts.SshIdentity)
	}

	return loginWithUserAndPassword()
}

func loginWithUserAndPassword() error {
	if config.AuthenticatedUser != nil {
		fmt.Printf("You are already logged as %s (%s)!\n", config.AuthenticatedUser.Email, config.AuthenticatedUser.Username)
		return nil
	}

	_, err := config.Authenticator.LoadAuth()
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}
	return nil
}

func loginWithApiKey(apiKey string) error {
	c := config.ScalingoUnauthenticatedClient()
	c.APIToken = apiKey
	user, err := c.Self()
	if err != nil {
		return errgo.Mask(err)
	}

	fmt.Printf("Hello %s, nice to see you!\n", user.Username)
	user.AuthenticationToken = apiKey

	err = config.Authenticator.StoreAuth(user)
	if err != nil {
		return errgo.Mask(err)
	}
	return nil
}

func loginWithSsh(identity string) error {
	debug.Println("Login through SSH, identity:", identity)
	client, _, err := netssh.Connect(identity)
	if err != nil {
		return errgo.Notef(err, "fail to connect to SSH server")
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
