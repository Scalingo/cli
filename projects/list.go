package projects

import (
	"context"
	stdio "io"
	"os"
	"strconv"
	"sync"

	"github.com/olekukonko/tablewriter"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-scalingo/v11"
	"github.com/Scalingo/go-utils/errors/v3"
)

const maxConcurrentRequests = 5

type projectsLister interface {
	ProjectsList(ctx context.Context) ([]scalingo.Project, error)
	ProjectsListByOrganization(ctx context.Context, organizationID string) ([]scalingo.Project, error)
}

// listedProject is a project with the slug of its organization, empty for a personal project.
type listedProject struct {
	project          scalingo.Project
	organizationSlug string
}

func List(ctx context.Context) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "get Scalingo client")
	}

	projects, err := listWithOrganizations(ctx, c, c)
	if err != nil {
		return err
	}

	io.Warning("This command only displays projects where you are the owner")

	return render(os.Stdout, projects)
}

// listWithOrganizations returns the personal projects followed by the projects of each organization of the user.
func listWithOrganizations(ctx context.Context, projectsService projectsLister, organizationsService scalingo.OrganizationsService) ([]listedProject, error) {
	personalProjects, err := projectsService.ProjectsList(ctx)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "list projects")
	}

	organizations, err := organizationsService.OrganizationsList(ctx)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "list organizations")
	}

	projects := make([]listedProject, 0, len(personalProjects))
	for _, project := range personalProjects {
		projects = append(projects, listedProject{project: project})
	}

	// The authentication service can return the same organization several times.
	uniqueOrganizations := make([]scalingo.Organization, 0, len(organizations))
	listedOrganizations := make(map[string]bool, len(organizations))
	for _, organization := range organizations {
		if !listedOrganizations[organization.ID] {
			listedOrganizations[organization.ID] = true
			uniqueOrganizations = append(uniqueOrganizations, organization)
		}
	}

	organizationsProjects, err := listProjectsOfOrganizations(ctx, projectsService, uniqueOrganizations)
	if err != nil {
		return nil, err
	}
	for i, organizationProjects := range organizationsProjects {
		for _, project := range organizationProjects {
			projects = append(projects, listedProject{project: project, organizationSlug: uniqueOrganizations[i].Slug})
		}
	}

	return projects, nil
}

// listProjectsOfOrganizations requests the projects of the organizations concurrently and returns them in the order of the organizations.
// The SDK client is only safe for concurrent use once it has completed a first request (its token and HTTP client are cached then).
func listProjectsOfOrganizations(ctx context.Context, service projectsLister, organizations []scalingo.Organization) ([][]scalingo.Project, error) {
	projects := make([][]scalingo.Project, len(organizations))
	errs := make([]error, len(organizations))
	slots := make(chan struct{}, maxConcurrentRequests)

	var wg sync.WaitGroup
	for i, organization := range organizations {
		wg.Go(func() {
			slots <- struct{}{}
			defer func() { <-slots }()
			projects[i], errs[i] = service.ProjectsListByOrganization(ctx, organization.ID)
		})
	}
	wg.Wait()

	for i, err := range errs {
		if err != nil {
			return nil, errors.Wrapf(ctx, err, "list projects of organization %s", organizations[i].Slug)
		}
	}

	return projects, nil
}

func render(w stdio.Writer, projects []listedProject) error {
	t := tablewriter.NewWriter(w)
	t.Header([]string{"Name", "Organization", "Default", "ID"})

	for _, listed := range projects {
		project := listed.project
		_ = t.Append([]string{project.Name, listed.organizationSlug, strconv.FormatBool(project.Default), project.ID})
	}

	return t.Render()
}
