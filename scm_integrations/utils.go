package scm_integrations

import (
	"github.com/Scalingo/go-scalingo"
)

func checkIfIntegrationAlreadyExist(c *scalingo.Client, id string) bool {
	integrations, _ := c.SCMIntegrationsShow(id)
	if integrations != nil {
		return true
	}
	return false
}
