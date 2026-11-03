package table

import (
	"context"
	"os"

	"github.com/olekukonko/tablewriter"

	"github.com/Scalingo/cli/internal/boundaries/out/renderer"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-scalingo/v11"
)

type organizationsListRenderer struct {
	organizations []scalingo.Organization
}

func NewOrganizationsList() renderer.Renderer[[]scalingo.Organization] {
	return &organizationsListRenderer{}
}

func (r *organizationsListRenderer) Render(_ context.Context) error {
	if len(r.organizations) == 0 {
		io.Status("You are not a member of any organization.")
		return nil
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.Header([]string{"ID", "SLUG", "ROLE"})

	for _, organization := range r.organizations {
		_ = t.Append([]string{organization.ID, organization.Slug, organization.Role})
	}

	return t.Render()
}

func (r *organizationsListRenderer) SetData(_ context.Context, organizations []scalingo.Organization) {
	r.organizations = organizations
}
