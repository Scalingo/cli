package scalingo

import (
	"context"
	stderrors "errors"

	"github.com/Scalingo/go-utils/errors/v2"
)

const databasesResource = "databases"

var ErrDatabaseNotFound = stderrors.New("database not found")

type PreviewAPI interface {
	DatabasesPreviewService
}

type DatabasesPreviewService interface {
	DatabaseCreate(ctx context.Context, params DatabaseCreateParams) (DatabaseNG, error)
	DatabasesList(ctx context.Context) ([]DatabaseNG, error)
	DatabaseShow(ctx context.Context, appID string) (DatabaseNG, error)
	DatabaseDestroy(ctx context.Context, appID string) error
}

var _ DatabasesPreviewService = (*PreviewClient)(nil)

// DatabaseNG stands for Database Next Generation.
type DatabaseNG struct {
	App      App       `json:"app"`
	Addon    Addon     `json:"addon"`
	Database *Database `json:"database,omitempty"`
}

type PreviewClient struct {
	parent *Client
}

func NewPreviewClient(parent *Client) *PreviewClient {
	return &PreviewClient{
		parent: parent,
	}
}

type DatabaseCreateParams struct {
	AddonProviderID string `json:"addon_provider_id"`
	PlanID          string `json:"plan_id"`
	Name            string `json:"name"`
	ProjectID       string `json:"project_id,omitempty"`
}

func (c *PreviewClient) DatabaseCreate(ctx context.Context, params DatabaseCreateParams) (DatabaseNG, error) {
	var res DatabaseNG

	err := c.parent.ScalingoAPI().ResourceAdd(ctx, databasesResource, params, &res)
	if err != nil {
		return res, errors.Wrap(ctx, err, "create database")
	}
	return res, nil
}

func (c *PreviewClient) DatabasesList(ctx context.Context) ([]DatabaseNG, error) {
	var res []DatabaseNG

	err := c.parent.ScalingoAPI().ResourceList(ctx, databasesResource, nil, &res)
	if err != nil {
		return res, errors.Wrap(ctx, err, "list databases")
	}
	return res, nil
}

// DatabaseShow currently uses appID as the database identifier.
func (c *PreviewClient) DatabaseShow(ctx context.Context, appID string) (DatabaseNG, error) {
	var res DatabaseNG

	databaseNG, err := c.searchDatabase(ctx, appID)
	if err != nil {
		return res, errors.Wrap(ctx, err, "search database")
	}

	database, err := c.parent.DatabaseShow(ctx, databaseNG.App.ID, databaseNG.Addon.ID)
	if err != nil {
		return res, errors.Wrap(ctx, err, "show database")
	}

	res.App = databaseNG.App
	res.Addon = databaseNG.Addon
	res.Database = &database

	return res, nil
}

// DatabaseDestroy currently uses appID as the database identifier.
func (c *PreviewClient) DatabaseDestroy(ctx context.Context, appID string) error {
	database, err := c.searchDatabase(ctx, appID)
	if err != nil {
		return errors.Wrap(ctx, err, "search database")
	}

	appName := database.App.Name

	err = c.parent.AppsDestroy(ctx, appName, appName)
	if err != nil {
		return errors.Wrap(ctx, err, "destroy database app")
	}
	return nil
}

// searchDatabase performs a linear search through DatabasesList method result.
func (c *PreviewClient) searchDatabase(ctx context.Context, appID string) (DatabaseNG, error) {
	var res DatabaseNG

	databases, err := c.DatabasesList(ctx)
	if err != nil {
		return res, errors.Wrap(ctx, err, "list databases")
	}

	for _, databaseNG := range databases {
		if databaseNG.App.ID == appID {
			return databaseNG, nil
		}
	}
	return res, ErrDatabaseNotFound
}
