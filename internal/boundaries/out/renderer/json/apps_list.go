package json

import (
	"context"
	"encoding/json"
	"os"

	"github.com/Scalingo/cli/apps"
	"github.com/Scalingo/cli/internal/boundaries/out/renderer"
	"github.com/Scalingo/go-scalingo/v11"
	"github.com/Scalingo/go-utils/errors/v3"
)

type appsListRenderer struct {
	apps []*scalingo.App
}

type appsListResponse struct {
	Apps []*scalingo.App `json:"apps"`
}

func NewAppsList() renderer.Renderer[[]apps.ListedApp] {
	return &appsListRenderer{}
}

func (r *appsListRenderer) Render(ctx context.Context) error {
	err := json.NewEncoder(os.Stdout).Encode(appsListResponse{Apps: r.apps})
	if err != nil {
		return errors.Wrap(ctx, err, "encode apps list to JSON")
	}
	return nil
}

func (r *appsListRenderer) SetData(ctx context.Context, listedApps []apps.ListedApp) {
	if listedApps == nil {
		r.apps = nil
		return
	}

	r.apps = make([]*scalingo.App, 0, len(listedApps))
	for _, listed := range listedApps {
		r.apps = append(r.apps, listed.App)
	}
}
