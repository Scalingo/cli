package scm_integrations

import (
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
)

type CreateArgs struct {
	ScmType string
	Url     string
	Token   string
}

func Create(args CreateArgs) error {
	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	integrationExist, err := checkIfIntegrationAlreadyExist(c, args.ScmType)
	if err != nil {
		return errgo.Notef(err, "fail to check if integration exist")
	}
	if integrationExist {
		io.Statusf("Integration '%s' is already linked with your Scalingo account.\n", args.ScmType)
		return nil
	}

	if args.ScmType == "github" || args.ScmType == "gitlab" {
		io.Statusf("Please follow this url to create the %s integration :\n", args.ScmType)
		io.Statusf("%s/users/%s/link\n", config.C.ScalingoAuthUrl, args.ScmType)
		return nil
	}

	_, err = c.IntegrationsCreate(args.ScmType, args.Url, args.Token)
	if err != nil {
		return errgo.Notef(err, "fail to create integration")
	}

	io.Statusf("Integration '%s' has been added.\n", args.ScmType)
	return nil
}
