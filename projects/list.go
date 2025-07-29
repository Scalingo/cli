package projects

import (
	"context"
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-utils/errors/v2"
)

func List(ctx context.Context) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "get Scalingo client")
	}

	projects, err := c.ProjectsList(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "list projects")
	}

	io.Warning("This command only displays projects where you are the owner")

	t := tablewriter.NewWriter(os.Stdout)
	t.Header([]string{"Name", "Default", "ID"})

	for _, project := range projects {
		_ = t.Append([]string{project.Name, strconv.FormatBool(project.Default), project.ID})
	}
	_ = t.Render()

	return nil
}
