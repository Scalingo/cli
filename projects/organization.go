package projects

import (
	"context"

	"github.com/Scalingo/go-scalingo/v11"
	"github.com/Scalingo/go-scalingo/v11/debug"
)

type projectGetter interface {
	ProjectGet(ctx context.Context, projectID string) (scalingo.Project, error)
}

// OrganizationProjectSlug returns "<organization slug>/<project name>" for a project of an organization.
// It returns an empty string when the project has no ID, belongs to the personal scope, or when its
// organization cannot be resolved. A failure must not fail the caller, which keeps its own output.
// It performs at most one project read and one organizations read.
func OrganizationProjectSlug(ctx context.Context, projectsService projectGetter, organizationsService scalingo.OrganizationsService, projectID string) string {
	if projectID == "" {
		return ""
	}

	project, err := projectsService.ProjectGet(ctx, projectID)
	if err != nil {
		debug.Printf("[Projects] Fail to get the project %s: %v\n", projectID, err)
		return ""
	}
	if project.OrganizationID == "" {
		return ""
	}

	organizations, err := organizationsService.OrganizationsList(ctx)
	if err != nil {
		debug.Printf("[Projects] Fail to list the organizations: %v\n", err)
		return ""
	}

	for _, organization := range organizations {
		if organization.ID == project.OrganizationID && organization.Slug != "" {
			return organization.Slug + "/" + project.Name
		}
	}

	return ""
}
