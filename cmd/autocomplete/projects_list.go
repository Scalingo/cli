package autocomplete

import (
	"context"
	"fmt"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/go-utils/errors/v2"
)

func ProjectsGenericListAutoComplete(ctx context.Context) error {
	client, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "get Scalingo client")
	}
	projects, err := client.ProjectsList(ctx)
	if err == nil {
		for _, proj := range projects {
			fmt.Println(proj.ID)
		}
	}

	return nil
}
