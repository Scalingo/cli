package projects

import (
	"context"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-scalingo/v11"
	"github.com/Scalingo/go-utils/errors/v3"
)

// Add creates a personal project, or a project in the organization identified by organizationSlug when it is not empty.
func Add(ctx context.Context, params scalingo.ProjectAddParams, organizationSlug string) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "get Scalingo client")
	}

	project, err := add(ctx, c, params, organizationSlug)
	if err != nil {
		return err
	}

	if organizationSlug == "" {
		io.Status(project.Name, "has been created\n")
		return nil
	}

	io.Statusf("%s has been created in organization %s\n", project.Name, organizationSlug)

	return nil
}

// add resolves the organization slug to its ID with a single organizations read, then creates the project.
func add(ctx context.Context, client scalingo.API, params scalingo.ProjectAddParams, organizationSlug string) (scalingo.Project, error) {
	if organizationSlug != "" {
		organizations, err := client.OrganizationsList(ctx)
		if err != nil {
			return scalingo.Project{}, errors.Wrap(ctx, err, "list organizations")
		}

		for _, organization := range organizations {
			if organization.Slug == organizationSlug {
				params.OrganizationID = organization.ID
				break
			}
		}
		if params.OrganizationID == "" {
			return scalingo.Project{}, errors.Newf(ctx, "organization %q not found", organizationSlug)
		}
	}

	project, err := client.ProjectAdd(ctx, params)
	if err != nil {
		return scalingo.Project{}, errors.Wrap(ctx, err, "add project")
	}

	return project, nil
}
