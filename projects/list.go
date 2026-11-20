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

	projects, err := listWithOrganizations(ctx, c)
	if err != nil {
		return errors.Wrap(ctx, err, "list projects with organizations")
	}

	io.Warning("This command only displays projects where you are the owner")

	err = render(os.Stdout, projects)
	if err != nil {
		return errors.Wrap(ctx, err, "render listed projects")
	}

	return nil
}

// listWithOrganizations returns the personal projects followed by the projects of each organization of the user.
func listWithOrganizations(ctx context.Context, client scalingo.API) ([]listedProject, error) {
	personalProjects, err := client.ProjectsList(ctx)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "list projects")
	}

	organizations, err := client.OrganizationsList(ctx)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "list organizations")
	}

	projects := make([]listedProject, 0, len(personalProjects))
	for _, project := range personalProjects {
		projects = append(projects, listedProject{project: project})
	}

	for _, organization := range organizations {
		organizationProjects, err := client.ProjectsListByOrganization(ctx, organization.ID)
		if err != nil {
			return nil, errors.Wrap(ctx, err, "list organization projects")
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
