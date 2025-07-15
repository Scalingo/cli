package scalingo

import (
	"context"

	"gopkg.in/errgo.v1"
)

type CollaboratorStatus string

const (
	CollaboratorStatusPending  CollaboratorStatus = "pending"
	CollaboratorStatusAccepted CollaboratorStatus = "accepted"
	CollaboratorStatusDeleted  CollaboratorStatus = "user account deleted"
)

type CollaboratorsService interface {
	CollaboratorsList(ctx context.Context, app string) ([]Collaborator, error)
	CollaboratorAdd(ctx context.Context, app string, params CollaboratorAddParams) (Collaborator, error)
	CollaboratorRemove(ctx context.Context, app, collaboratorID string) error
	CollaboratorUpdate(ctx context.Context, app, collaboratorID string, params CollaboratorUpdateParams) (Collaborator, error)
}

var _ CollaboratorsService = (*Client)(nil)

type Collaborator struct {
	ID        string             `json:"id"`
	AppID     string             `json:"app_id"`
	Username  string             `json:"username"`
	Email     string             `json:"email"`
	Status    CollaboratorStatus `json:"status"`
	UserID    string             `json:"user_id"`
	IsLimited bool               `json:"is_limited"`
}

type CollaboratorsRes struct {
	Collaborators []Collaborator `json:"collaborators"`
}

type CollaboratorRes struct {
	Collaborator Collaborator `json:"collaborator"`
}

type CollaboratorAddParams struct {
	Email     string `json:"email"`
	IsLimited bool   `json:"is_limited"`
}

type CollaboratorAddParamsPayload struct {
	Collaborator CollaboratorAddParams `json:"collaborator"`
}

type CollaboratorUpdateParams struct {
	IsLimited bool `json:"is_limited"`
}

type CollaboratorUpdateParamsPayload struct {
	Collaborator CollaboratorUpdateParams `json:"collaborator"`
}

func (c *Client) CollaboratorsList(ctx context.Context, app string) ([]Collaborator, error) {
	var collaboratorsRes CollaboratorsRes
	err := c.ScalingoAPI().SubresourceList(ctx, "apps", app, "collaborators", nil, &collaboratorsRes)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	return collaboratorsRes.Collaborators, nil
}

func (c *Client) CollaboratorAdd(ctx context.Context, app string, params CollaboratorAddParams) (Collaborator, error) {
	var collaboratorRes CollaboratorRes
	err := c.ScalingoAPI().SubresourceAdd(ctx, "apps", app, "collaborators", CollaboratorAddParamsPayload{params}, &collaboratorRes)
	if err != nil {
		return Collaborator{}, errgo.Mask(err)
	}
	return collaboratorRes.Collaborator, nil
}

func (c *Client) CollaboratorRemove(ctx context.Context, app, collaboratorID string) error {
	return c.ScalingoAPI().SubresourceDelete(ctx, "apps", app, "collaborators", collaboratorID)
}

func (c *Client) CollaboratorUpdate(ctx context.Context, app, collaboratorID string, params CollaboratorUpdateParams) (Collaborator, error) {
	var collaboratorRes CollaboratorRes
	err := c.ScalingoAPI().SubresourceUpdate(ctx, "apps", app, "collaborators", collaboratorID, CollaboratorUpdateParamsPayload{params}, &collaboratorRes)
	if err != nil {
		return Collaborator{}, errgo.Mask(err)
	}
	return collaboratorRes.Collaborator, nil
}
