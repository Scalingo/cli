package integrations

import (
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/go-scalingo"
)

func integrationByName(c *scalingo.Client, name string) (*scalingo.Integration, error) {
	integrations, err := c.IntegrationsList()
	if err != nil {
		return nil, errgo.Mask(err)
	}

	for _, k := range integrations {
		if k.ScmType == name {
			return &k, nil
		}
	}

	return nil, errgo.New("not linked integration or unknown integration : '" + name + "'")
}

func integrationByUUID(c *scalingo.Client, uuid string) (*scalingo.Integration, error) {
	integration, err := c.IntegrationsShow(uuid)
	if err != nil {
		return nil, errgo.New("not linked integration or unknown integration")
	}
	return integration, nil
}
