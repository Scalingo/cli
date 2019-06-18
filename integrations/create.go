package integrations

import (
	"fmt"
	"net/http"
	"net/url"

	"gopkg.in/errgo.v1"
)


func Create(scm_type string, link string, token string) error {
	switch scm_type {
	case "github", "gitlab":

		fmt.Printf("Integration '%s' has been added.\n", scm_type)
		return nil
	case "github-enterprise", "gitlab-self-hosted":
		if link == "" {
			return errgo.New("URL is empty")
		}

		u, err := url.Parse(link)
		if err != nil || u.Scheme == "" || u.Host == "" {
			return errgo.New("URL format is invalid, valid format is : https://gitlab.domain.com")
		}

		resp, err := http.Get(link)
		if err != nil {
			return errgo.New("Failed to access to the integration instance")
		}

		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return errgo.New("Failed to access to the integration instance")
		}

		if token == "" {
			return errgo.New("Token is empty")
		}

		fmt.Printf("Integration '%s' has been added.\n", scm_type)
		return nil
	default:
		return errgo.New("Type don't exist or is empty")
	}
}
