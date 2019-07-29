package scm_integrations

import (
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/cli/utils"
)

func Destroy(integration string) error {
	var id string
	var name string

	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	if !utils.IsUUID(integration) {
		i, err := IntegrationByName(c, integration)
		if err != nil {
			return errgo.Notef(err, "fail to get integration by name")
		}

		id = i.ID
		name = i.ScmType
	} else {
		i, err := IntegrationByUUID(c, integration)
		if err != nil {
			return errgo.Notef(err, "fail to get integration by uuid")
		}

		id = integration
		name = i.ScmType
	}

	err = c.IntegrationsDestroy(id)
	if err != nil {
		return errgo.Notef(err, "fail to destroy integration")
	}

	io.Statusf("Integration '%s' has been deleted.\n", name)
	return nil
}
