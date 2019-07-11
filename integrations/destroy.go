package integrations

import (
	"fmt"

	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
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
		i, err := integrationByName(c, integration)
		if err != nil {
			return errgo.Mask(err)
		}

		id = i.ID
		name = i.ScmType
	} else {
		i, err := integrationByUUID(c, integration)
		if err != nil {
			return errgo.Mask(err)
		}

		id = integration
		name = i.ScmType
	}

	err = c.IntegrationsDestroy(id)
	if err != nil {
		return errgo.Mask(err)
	}

	fmt.Printf("Integration '%s' has been deleted.\n", name)
	return nil
}
