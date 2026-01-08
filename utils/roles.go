package utils

import (
	"github.com/Scalingo/go-scalingo/v9"
)

const (
	roleOwner        Role = "owner"
	roleCollaborator Role = "collaborator"
)

type Role string

func AppRole(user *scalingo.User, app *scalingo.App) Role {
	if user.Email == app.Owner.Email {
		return roleOwner
	}
	return roleCollaborator
}
