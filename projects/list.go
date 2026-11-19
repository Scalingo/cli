package projects

import (
	"context"
	stdio "io"
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-scalingo/v11"
	"github.com/Scalingo/go-utils/errors/v3"
)

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
	listedOrganizations := make(map[string]bool, len(organizations))
	for _, organization := range organizations {
		if listedOrganizations[organization.ID] {
			continue
		}
		listedOrganizations[organization.ID] = true

		organizationProjects, err := projectsService.ProjectsListByOrganization(ctx, organization.ID)
		if err != nil {
			return nil, errors.Wrapf(ctx, err, "list projects of organization %s", organization.Slug)
		}
		for _, project := range organizationProjects {
			projects = append(projects, listedProject{project: project, organizationSlug: organization.Slug})
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
