package scm_integrations

import (
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/go-scalingo"
)

func checkIfIntegrationAlreadyExist(c *scalingo.Client, id string) (bool, error) {
	integrations, err := c.SCMIntegrationsShow(id)
	if err != nil {
		return false, errgo.Notef(err, "fail to show SCM integrations")
	}
	if integrations != nil {
		return true, nil
	}
	return false, nil
}
