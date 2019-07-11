package integrations

import (
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/go-scalingo"
)

func IntegrationByName(c *scalingo.Client, name string) (*scalingo.Integration, error) {
	integrations, err := c.IntegrationsList()
	if err != nil {
		return nil, errgo.Mask(err)
	}

	for _, i := range integrations {
		if i.ScmType == name {
			return &i, nil
		}
	}

	return nil, errgo.New("not linked integration or unknown integration : '" + name + "'")
}

func IntegrationByUUID(c *scalingo.Client, uuid string) (*scalingo.Integration, error) {
	integration, err := c.IntegrationsShow(uuid)
	if err != nil {
		return nil, errgo.New("not linked integration or unknown integration")
	}
	return integration, nil
}

func checkIfIntegrationAlreadyExist(c *scalingo.Client, name string) (bool, error) {
	integrations, err := c.IntegrationsList()
	if err != nil {
		return false, errgo.Notef(err, "fail to get integrations list")
	}

	for _, i := range integrations {
		if i.ScmType == name {
			return true, nil
		}
	}

	return false, nil
}
