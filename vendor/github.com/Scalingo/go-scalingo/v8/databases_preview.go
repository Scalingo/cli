package scalingo

import (
	"context"
	stderrors "errors"
	"fmt"

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
	DatabaseInfo
	Database *Database `json:"database,omitempty"`
	App      App       `json:"app"`
}

type databaseNgResponse struct {
	// App      App          `json:"app"` 	// Thoses fields will be removed
	// Addon    Addon        `json:"addon"` // Thoses fields will be removed
	Database DatabaseInfo `json:"database"`
}

type DatabaseInfo struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	ProjectID  string `json:"project_id"`
	Technology string `json:"technology"`
	Plan       string `json:"plan"`
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
	var listResp []databaseNgResponse

	err := c.parent.ScalingoAPI().ResourceList(ctx, databasesResource, nil, &listResp)
	if err != nil {
		return res, errors.Wrap(ctx, err, "list databases")
	}

	for _, response := range listResp {
		databaseNG, err := c.formatDatabaseNG(ctx, response)
		if err != nil {
			return res, errors.Wrap(ctx, err, "populate databaseNG")
		}

		res = append(res, databaseNG)
	}

	return res, nil
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

// formatDatabaseNG populates a DatabaseNG without using the App and Addon from the databases endpoints.
func (c *PreviewClient) formatDatabaseNG(ctx context.Context, response databaseNgResponse) (DatabaseNG, error) {
	var res DatabaseNG

	res.DatabaseInfo = response.Database

	addons, err := c.parent.AddonsList(ctx, response.Database.ID)
	if err != nil {
		return res, errors.Wrap(ctx, err, "list addons")
	}

	appPtr, err := c.parent.AppsShow(ctx, response.Database.Name)
	if err != nil {
		return res, errors.Wrap(ctx, err, "show app")
	}
	res.App = *appPtr

	database, err := c.parent.DatabaseShow(ctx, response.Database.ID, addons[0].ID)
	if err != nil {
		fmt.Printf("Addons probably deleted for app: %+v\n", res.DatabaseInfo.Name)
	}
	res.Database = &database

	return res, nil
}
