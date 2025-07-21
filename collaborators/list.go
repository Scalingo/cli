package collaborators

import (
	"context"
	"os"

	"github.com/olekukonko/tablewriter"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/go-utils/errors/v2"
)

const (
	owner        = "owner"
	collaborator = "collaborator"
	limited      = "limited collaborator"
)

func List(ctx context.Context, app string) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "get Scalingo client")
	}

	collaborators, err := c.CollaboratorsList(ctx, app)
	if err != nil {
		return errors.Wrap(ctx, err, "list collaborators")
	}

	scapp, err := c.AppsShow(ctx, app)
	if err != nil {
		return errors.Wrap(ctx, err, "get application information")
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.Header([]string{"Email", "Username", "Status", "Role"})

	_ = t.Append([]string{scapp.Owner.Email, scapp.Owner.Username, owner, owner})
	for _, collaborator := range collaborators {
		_ = t.Append([]string{collaborator.Email, collaborator.Username,
			string(collaborator.Status), collaboratorToTole(collaborator.IsLimited)})
	}
	_ = t.Render()
	return nil
}

// collaboratorToTole converts a collaborator into a human-readable role.
func collaboratorToTole(isLimited bool) string {
	if isLimited {
		return limited
	}

	return collaborator
}
