package projects

import (
	"context"
	"time"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/go-scalingo/v11"
	"github.com/Scalingo/go-scalingo/v11/debug"
)

// enrichmentTimeout bounds each HTTP request made to enrich the confirmation of a resource creation. The
// resource is already created when these requests are sent, so a slow API must not delay the rest of the command.
// The SDK builds its requests without the context, so this timeout is set on the HTTP client and the same
// duration is also set on the context to skip the organizations read once it is expired. At most 4 requests
// are sent: the token exchange and the read of the project, then the same for the organizations.
const enrichmentTimeout = 2 * time.Second

type projectGetter interface {
	ProjectGet(ctx context.Context, projectID string) (scalingo.Project, error)
}

type clientFactory func(ctx context.Context, timeout time.Duration) (*scalingo.Client, error)

// CreatedProjectSlug returns "<organization slug>/<project name>" for the project of a resource that has just
// been created, with a short timeout. It returns an empty string when the project is not in an organization or
// cannot be resolved, so that the caller keeps its historical confirmation.
func CreatedProjectSlug(ctx context.Context, projectID string) string {
	return createdProjectSlug(ctx, config.ScalingoClientWithTimeout, enrichmentTimeout, projectID)
}

func createdProjectSlug(ctx context.Context, newClient clientFactory, timeout time.Duration, projectID string) string {
	if projectID == "" {
		return ""
	}

	// A dedicated client is required because the timeout of a shared client cannot be changed.
	c, err := newClient(ctx, timeout)
	if err != nil {
		debug.Printf("[Projects] Fail to get a Scalingo client: %v\n", err)
		return ""
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return OrganizationProjectSlug(ctx, c, c, projectID)
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
	if ctx.Err() != nil {
		debug.Printf("[Projects] Skip the organizations read: %v\n", ctx.Err())
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
