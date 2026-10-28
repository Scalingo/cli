package json

import (
	"context"
	"encoding/json"
	"os"

	"github.com/Scalingo/cli/internal/boundaries/out/renderer"
	"github.com/Scalingo/go-scalingo/v11"
	"github.com/Scalingo/go-utils/errors/v3"
)

type orgsListRenderer struct {
	organizations []scalingo.Organization
}

type orgsListResponse struct {
	Organizations []scalingo.Organization `json:"organizations"`
}

func NewOrgsList() renderer.Renderer[[]scalingo.Organization] {
	return &orgsListRenderer{}
}

func (r *orgsListRenderer) Render(ctx context.Context) error {
	err := json.NewEncoder(os.Stdout).Encode(orgsListResponse{Organizations: r.organizations})
	if err != nil {
		return errors.Wrap(ctx, err, "encode organizations list to JSON")
	}
	return nil
}

func (r *orgsListRenderer) SetData(_ context.Context, organizations []scalingo.Organization) {
	r.organizations = organizations
}
