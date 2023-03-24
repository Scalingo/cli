package scmintegrations

import (
	"context"

	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-scalingo/v6"
)

type CreateArgs struct {
	SCMType scalingo.SCMType
	URL     string
	Token   string
}

func Create(ctx context.Context, args CreateArgs) error {
	c, err := config.ScalingoAuthClient(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	integrationExist := checkIfIntegrationAlreadyExist(ctx, c, args.SCMType.Str())
	if integrationExist {
		io.Statusf("SCM Integration '%s' is already linked with your Scalingo account.\n", scalingo.SCMTypeDisplay[args.SCMType])
		return nil
	}

	if args.SCMType == scalingo.SCMGithubType || args.SCMType == scalingo.SCMGitlabType {
		io.Statusf("Please follow this URL to create the %s SCM integration:\n", scalingo.SCMTypeDisplay[args.SCMType])
		io.Statusf("%s/users/%s/link\n", config.C.ScalingoAuthURL, args.SCMType)
		return nil
	}

	_, err = c.SCMIntegrationsCreate(ctx, args.SCMType, args.URL, args.Token)
	if err != nil {
		return errgo.Notef(err, "fail to create the SCM integration")
	}

	io.Statusf("Your Scalingo account has been linked to your %s account.\n", scalingo.SCMTypeDisplay[args.SCMType])
	return nil
}
