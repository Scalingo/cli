package scalingo

import (
	"context"
	stderrors "errors"

	"github.com/Scalingo/go-scalingo/v8/debug"
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
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	ProjectID  string   `json:"project_id"`
	Technology string   `json:"technology"`
	Plan       string   `json:"plan"`
	Database   Database `json:"database"`
	App        App      `json:"app"`
}

type databaseListItem struct {
	Database DatabaseNG `json:"database"`
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
	var listRes []databaseListItem

	err := c.parent.ScalingoAPI().ResourceList(ctx, databasesResource, nil, &listRes)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "list databases")
	}

	databasesNG := make([]DatabaseNG, len(listRes))

	for i, apiResponse := range listRes {
		dbNG, err := c.populateAPIResponse(ctx, apiResponse)
		if err != nil {
			return nil, errors.Wrap(ctx, err, "populate database NG")
		}

		databasesNG[i] = dbNG
	}

	return databasesNG, nil
}

// DatabaseShow currently uses appID as the database identifier.
func (c *PreviewClient) DatabaseShow(ctx context.Context, appID string) (DatabaseNG, error) {
	var res DatabaseNG

	res, err := c.searchDatabase(ctx, appID)
	if err != nil {
		return res, errors.Wrap(ctx, err, "search database")
	}

	return res, nil
}

// DatabaseDestroy currently uses appID as the database identifier.
func (c *PreviewClient) DatabaseDestroy(ctx context.Context, appID string) error {
	database, err := c.searchDatabase(ctx, appID)
	if err != nil {
		return errors.Wrap(ctx, err, "search database")
	}

	err = c.parent.AppsDestroy(ctx, database.Name, database.Name)
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
		if databaseNG.ID == appID {
			return databaseNG, nil
		}
	}
	return res, ErrDatabaseNotFound
}

// populateAPIResponse populates a DatabaseNG without using the App and Addon from the databases endpoints.
func (c *PreviewClient) populateAPIResponse(ctx context.Context, apiResponse databaseListItem) (DatabaseNG, error) {
	databaseNG := DatabaseNG{
		ID:         apiResponse.Database.ID,
		Name:       apiResponse.Database.Name,
		ProjectID:  apiResponse.Database.ProjectID,
		Technology: apiResponse.Database.Technology,
		Plan:       apiResponse.Database.Plan,
	}

	addons, err := c.parent.AddonsList(ctx, apiResponse.Database.ID)
	if err != nil {
		return databaseNG, errors.Wrap(ctx, err, "list addons")
	}

	if len(addons) == 0 {
		return databaseNG, errors.New(ctx, "no addons found for database")
	}

	app, err := c.parent.AppsShow(ctx, apiResponse.Database.Name)
	if err != nil {
		return databaseNG, errors.Wrap(ctx, err, "show app")
	}
	databaseNG.App = *app

	databaseNG.Database, err = c.parent.DatabaseShow(ctx, apiResponse.Database.ID, addons[0].ID)
	if err != nil {
		debug.Printf("Addon has been removed from app: %+v\n", databaseNG.Name)
	}

	return databaseNG, nil
}
