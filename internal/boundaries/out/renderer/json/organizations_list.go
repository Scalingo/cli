package json

import (
	"context"
	"encoding/json"
	"os"

	"github.com/Scalingo/cli/internal/boundaries/out/renderer"
	"github.com/Scalingo/go-scalingo/v11"
	"github.com/Scalingo/go-utils/errors/v3"
)

type organizationsListRenderer struct {
	organizations []scalingo.Organization
}

type organizationOutput struct {
	ID   string `json:"id"`
	Slug string `json:"slug"`
	Role string `json:"role"`
}

type organizationsListResponse struct {
	Organizations []organizationOutput `json:"organizations"`
}

func NewOrganizationsList() renderer.Renderer[[]scalingo.Organization] {
	return &organizationsListRenderer{}
}

func (r *organizationsListRenderer) Render(ctx context.Context) error {
	organizations := make([]organizationOutput, 0, len(r.organizations))
	for _, organization := range r.organizations {
		organizations = append(organizations, organizationOutput{ID: organization.ID, Slug: organization.Slug, Role: organization.Role})
	}

	err := json.NewEncoder(os.Stdout).Encode(organizationsListResponse{Organizations: organizations})
	if err != nil {
		return errors.Wrap(ctx, err, "encode organizations list to JSON")
	}
	return nil
}

func (r *organizationsListRenderer) SetData(_ context.Context, organizations []scalingo.Organization) {
	r.organizations = organizations
}
