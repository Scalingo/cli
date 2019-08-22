package scalingo

import "fmt"

type EventAuthorizeGithubType struct {
	Event
	TypeData EventAuthorizeGithubTypeData `json:"type_data"`
}

func (ev *EventAuthorizeGithubType) String() string {
	return fmt.Sprintf("GitHub account '%s' has been authorized", ev.TypeData.GithubUser.Login)
}

type EventAuthorizeGithubTypeData struct {
	GithubUser struct {
		Login     string `json:"login"`
		AvatarURL string `json:"avatar_url"`
	} `json:"github_user"`
}

type EventAuthorizeGitLabType struct {
	Event
	TypeData EventAuthorizeGitLabTypeData `json:"type_data"`
}

func (ev *EventAuthorizeGitLabType) String() string {
	return fmt.Sprintf("GitLab account '%s' has been authorized", ev.TypeData.GitLabUser.Login)
}

type EventAuthorizeGitLabTypeData struct {
	GitLabUser struct {
		Login     string `json:"login"`
		AvatarURL string `json:"avatar_url"`
	} `json:"gitlab_user"`
}

type EventRevokeGitLabType struct {
	Event
}

func (ev *EventRevokeGitLabType) String() string {
	return fmt.Sprintf("GitLab authorization has been revoked")
}

type EventRevokeGithubType struct {
	Event
}

func (ev *EventRevokeGithubType) String() string {
	return fmt.Sprintf("GitHub authorization has been revoked")
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
