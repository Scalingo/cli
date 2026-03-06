package json

import (
	"context"
	"encoding/json"
	"os"

	"github.com/Scalingo/cli/apps"
	"github.com/Scalingo/go-scalingo/v10"
	"github.com/Scalingo/go-utils/errors/v3"
)

type appsListRenderer struct {
}

type appsListResponse struct {
	Apps []*scalingo.App `json:"apps"`
}

func NewAppsList() apps.ListRenderer {
	return appsListRenderer{}
}

func (r appsListRenderer) Render(ctx context.Context, apps []*scalingo.App) error {
	err := json.NewEncoder(os.Stdout).Encode(appsListResponse{Apps: apps})
	if err != nil {
		return errors.Wrap(ctx, err, "encode apps list to JSON")
	}
	return nil
}
