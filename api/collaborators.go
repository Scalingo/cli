package api

import "github.com/Scalingo/cli/Godeps/_workspace/src/gopkg.in/errgo.v1"

type Collaborator struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Status   string `json:"status"`
}

type CollaboratorsRes struct {
	Collaborators []Collaborator `json:"collaborators"`
}

type CollaboratorRes struct {
	Collaborator Collaborator `json:"collaborator"`
}

func CollaboratorsList(app string) ([]Collaborator, error) {
	var collaboratorsRes CollaboratorsRes
	err := subresourceList(app, "collaborators", nil, &collaboratorsRes)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	return collaboratorsRes.Collaborators, nil
}

func CollaboratorAdd(app string, email string) (Collaborator, error) {
	var collaboratorRes CollaboratorRes
	err := subresourceAdd(app, "collaborators", CollaboratorRes{
		Collaborator: Collaborator{Email: email},
	}, &collaboratorRes)
	if err != nil {
		return Collaborator{}, errgo.Mask(err)
	}
	return collaboratorRes.Collaborator, nil
}

func CollaboratorRemove(app string, id string) error {
	return subresourceDelete(app, "collaborators", id)
}
