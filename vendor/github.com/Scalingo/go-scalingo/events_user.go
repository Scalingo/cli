package scalingo

import "fmt"

type EventAuthorizeGithubType struct {
	Event
	TypeData EventAuthorizeGithubTypeData `json:"type_data"`
}

func (ev *EventAuthorizeGithubType) String() string {
	return fmt.Sprintf("Github account '%s' has been authorized", ev.TypeData.GithubUser.Login)
}

type EventAuthorizeGithubTypeData struct {
	GithubUser struct {
		ID    int    `json:"id"`
		Login string `json:"login"`
		Email string `json:"email"`
		Name  string `json:"name"`
	} `json:"github_user"`
}

type EventRevokeGithubType struct {
	Event
}

func (ev *EventRevokeGithubType) String() string {
	return fmt.Sprintf("Github authrization has been revoked")
}

type EventNewKeyType struct {
	Event
	TypeData EventNewKeyTypeData `json:"type_data"`
}

func (ev *EventNewKeyType) String() string {
	return fmt.Sprintf("name '%s' with fingerprint %s", ev.TypeData.Name, ev.TypeData.Fingerprint)
}

type EventNewKeyTypeData struct {
	Name        string `json:"name"`
	Fingerprint string `json:"fingerprint"`
}

type EventDeleteKeyType struct {
	Event
	TypeData EventDeleteKeyTypeData `json:"type_data"`
}

func (ev *EventDeleteKeyType) String() string {
	return fmt.Sprintf("name '%s'", ev.TypeData.Name)
}

type EventDeleteKeyTypeData struct {
	Name string `json:"name"`
}
