package scalingo

import "gopkg.in/errgo.v1"

type CollaboratorStatus string

const (
	CollaboratorStatusPending  CollaboratorStatus = "pending"
	CollaboratorStatusAccepted CollaboratorStatus = "accepted"
	CollaboratorStatusDeleted  CollaboratorStatus = "user account deleted"
)

type CollaboratorsService interface {
	CollaboratorsList(app string) ([]Collaborator, error)
	CollaboratorAdd(app string, email string) (Collaborator, error)
	CollaboratorRemove(app string, id string) error
}

var _ CollaboratorsService = (*Client)(nil)

type Collaborator struct {
	ID       string             `json:"id"`
	AppID    string             `json:"app_id"`
	Username string             `json:"username"`
	Email    string             `json:"email"`
	Status   CollaboratorStatus `json:"status"`
	UserID   string             `json:"user_id"`
}

type CollaboratorsRes struct {
	Collaborators []Collaborator `json:"collaborators"`
}

type CollaboratorRes struct {
	Collaborator Collaborator `json:"collaborator"`
}

func (c *Client) CollaboratorsList(app string) ([]Collaborator, error) {
	var collaboratorsRes CollaboratorsRes
	err := c.ScalingoAPI().SubresourceList("apps", app, "collaborators", nil, &collaboratorsRes)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	return collaboratorsRes.Collaborators, nil
}

func (c *Client) CollaboratorAdd(app string, email string) (Collaborator, error) {
	var collaboratorRes CollaboratorRes
	err := c.ScalingoAPI().SubresourceAdd("apps", app, "collaborators", CollaboratorRes{
		Collaborator: Collaborator{Email: email},
	}, &collaboratorRes)
	if err != nil {
		return Collaborator{}, errgo.Mask(err)
	}
	return collaboratorRes.Collaborator, nil
}

func (c *Client) CollaboratorRemove(app string, id string) error {
	return c.ScalingoAPI().SubresourceDelete("apps", app, "collaborators", id)
}
