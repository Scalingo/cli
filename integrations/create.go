package integrations

import (
	"fmt"

	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
)

func Create(scmType string, link string, token string) error {
	if scmType == "github" || scmType == "gitlab" {
		fmt.Printf("Please follow this url for create the %s integration :\n", scmType)
		fmt.Printf("===> %s/users/%s/link\n", config.C.ScalingoAuthUrl, scmType)
		return nil
	}

	if link == "" && token == "" {
		return errgo.New("URL and Token is empty")
	}
	if link == "" {
		return errgo.New("URL is empty")
	}
	if token == "" {
		return errgo.New("Token is empty")
	}

	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	_, err = c.IntegrationsCreate(scmType, link, token)
	if err != nil {
		return errgo.Mask(err)
	}

	fmt.Printf("Integration '%s' has been added.\n", scmType)
	return nil
}
