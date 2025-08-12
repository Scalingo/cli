package autocomplete

import (
	"fmt"

	"github.com/urfave/cli/v3"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/go-utils/errors/v2"
)

func ProjectsGenericListAutoComplete(c *cli.Context) error {
	client, err := config.ScalingoClient(c.Context)
	if err != nil {
		return errors.Wrap(c.Context, err, "get Scalingo client")
	}
	projects, err := client.ProjectsList(c.Context)
	if err == nil {
		for _, proj := range projects {
			fmt.Println(proj.ID)
		}
	}

	return nil
}
