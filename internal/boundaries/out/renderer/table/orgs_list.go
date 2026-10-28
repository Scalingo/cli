package table

import (
	"context"
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"

	"github.com/Scalingo/cli/internal/boundaries/out/renderer"
	"github.com/Scalingo/go-scalingo/v11"
	"github.com/Scalingo/go-utils/errors/v3"
)

type orgsListRenderer struct {
	organizations []scalingo.Organization
}

func NewOrgsList() renderer.Renderer[[]scalingo.Organization] {
	return &orgsListRenderer{}
}

func (r *orgsListRenderer) Render(ctx context.Context) error {
	if len(r.organizations) == 0 {
		fmt.Println("You are not a member of any organization.")
		return nil
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.Header([]string{"ID", "SLUG", "ROLE"})

	for _, organization := range r.organizations {
		err := t.Append([]string{organization.ID, organization.Slug, organization.Role})
		if err != nil {
			return errors.Wrap(ctx, err, "append organization to table")
		}
	}

	return t.Render()
}

func (r *orgsListRenderer) SetData(_ context.Context, organizations []scalingo.Organization) {
	r.organizations = organizations
}
