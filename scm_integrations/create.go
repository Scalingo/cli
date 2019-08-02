package scm_integrations

import (
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-scalingo"
)

type CreateArgs struct {
	SCMType scalingo.SCMType
	URL     string
	Token   string
}

func Create(args CreateArgs) error {
	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	integrationExist, err := checkIfIntegrationAlreadyExist(c, args.SCMType.Str())
	if err != nil {
		return errgo.Notef(err, "fail to check if SCM integration exist")
	}
	if integrationExist {
		io.Statusf("SCM Integration '%s' is already linked with your Scalingo account.\n", args.SCMType)
		return nil
	}

	if args.SCMType == scalingo.SCMGithubType || args.SCMType == scalingo.SCMGitlabType {
		io.Statusf("Please follow this URL to create the %s SCM integration :\n", args.SCMType)
		io.Statusf("%s/users/%s/link\n", config.C.ScalingoAuthUrl, args.SCMType)
		return nil
	}

	_, err = c.SCMIntegrationsCreate(args.SCMType, args.URL, args.Token)
	if err != nil {
		return errgo.Notef(err, "fail to create the SCM integration")
	}

	io.Statusf("Your Scalingo account has been linked to your '%s' account.\n", args.SCMType)
	return nil
}
