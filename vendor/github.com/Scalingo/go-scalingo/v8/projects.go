package scalingo

import (
	"context"
	"time"

	"github.com/Scalingo/go-utils/errors/v2"
)

const projectResource = "projects"

type ProjectsService interface {
	ProjectsList(ctx context.Context) ([]Project, error)
	ProjectAdd(ctx context.Context, params ProjectAddParams) (Project, error)
	ProjectUpdate(ctx context.Context, projectID string, params ProjectUpdateParams) (Project, error)
	ProjectGet(ctx context.Context, projectID string) (Project, error)
	ProjectDelete(ctx context.Context, projectID string) error
}

var _ ProjectsService = (*Client)(nil)

type Project struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Default   bool      `json:"default"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Owner     Owner     `json:"owner"`
}

type ProjectsRes struct {
	Projects []Project `json:"projects"`
}

type ProjectRes struct {
	Project Project `json:"project"`
}

type ProjectAddParams struct {
	Name    string `json:"name"`
	Default bool   `json:"default"`
}

type projectAddParamsPayload struct {
	Project ProjectAddParams `json:"project"`
}

type ProjectUpdateParams struct {
	Name    *string `json:"name"`
	Default *bool   `json:"default"`
}

type projectUpdateParamsPayload struct {
	Project ProjectUpdateParams `json:"project"`
}

func (c *Client) ProjectsList(ctx context.Context) ([]Project, error) {
	var projectsRes ProjectsRes
	err := c.ScalingoAPI().ResourceList(ctx, projectResource, nil, &projectsRes)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "list projects")
	}

	return projectsRes.Projects, nil
}

func (c *Client) ProjectAdd(ctx context.Context, params ProjectAddParams) (Project, error) {
	var projectRes ProjectRes
	err := c.ScalingoAPI().ResourceAdd(ctx, projectResource, projectAddParamsPayload{Project: params}, &projectRes)
	if err != nil {
		return Project{}, errors.Wrap(ctx, err, "add project")
	}

	return projectRes.Project, nil
}

func (c *Client) ProjectUpdate(ctx context.Context, projectID string, params ProjectUpdateParams) (Project, error) {
	var projectRes ProjectRes
	err := c.ScalingoAPI().ResourceUpdate(ctx, projectResource, projectID, projectUpdateParamsPayload{Project: params}, &projectRes)
	if err != nil {
		return Project{}, errors.Wrap(ctx, err, "update project")
	}

	return projectRes.Project, nil
}

func (c *Client) ProjectGet(ctx context.Context, projectID string) (Project, error) {
	var projectRes ProjectRes
	err := c.ScalingoAPI().ResourceGet(ctx, projectResource, projectID, nil, &projectRes)
	if err != nil {
		return Project{}, errors.Wrap(ctx, err, "get project")
	}

	return projectRes.Project, nil
}

func (c *Client) ProjectDelete(ctx context.Context, projectID string) error {
	err := c.ScalingoAPI().ResourceDelete(ctx, projectResource, projectID)
	if err != nil {
		return errors.Wrap(ctx, err, "delete project")
	}

	return nil
}
