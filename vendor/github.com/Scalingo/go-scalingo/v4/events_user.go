package scalingo

import "fmt"

type EventNewIntegrationType struct {
	Event
	TypeData EventNewIntegrationTypeData `json:"type_data"`
}

func (ev *EventNewIntegrationType) String() string {
	integrationType, ok := SCMTypeDisplay[ev.TypeData.IntegrationType]
	if !ok {
		integrationType = string(ev.TypeData.IntegrationType)
	}
	msg := fmt.Sprintf("%s", integrationType)

	if ev.TypeData.IntegrationType == SCMGithubEnterpriseType ||
		ev.TypeData.IntegrationType == SCMGitlabSelfHostedType {
		msg = fmt.Sprintf("%s (%s)", msg, ev.TypeData.Data.URL)
	}

	return fmt.Sprintf("%s account '%s' has been authorized", msg, ev.TypeData.Data.Login)
}

type EventNewIntegrationTypeData struct {
	IntegrationID   string  `json:"integration_id"`
	IntegrationType SCMType `json:"integration_type"`
	Data            struct {
		Login     string `json:"login"`
		AvatarURL string `json:"avatar_url"`
		URL       string `json:"url"`
	} `json:"data"`
}

type EventDeleteIntegrationType struct {
	Event
	TypeData EventDeleteIntegrationTypeData `json:"type_data"`
}

func (ev *EventDeleteIntegrationType) String() string {
	integrationType, ok := SCMTypeDisplay[ev.TypeData.IntegrationType]
	if !ok {
		integrationType = string(ev.TypeData.IntegrationType)
	}
	msg := fmt.Sprintf("%s", integrationType)

	if ev.TypeData.IntegrationType == SCMGithubEnterpriseType ||
		ev.TypeData.IntegrationType == SCMGitlabSelfHostedType {
		msg = fmt.Sprintf("%s (%s)", msg, ev.TypeData.Data.URL)
	}

	return fmt.Sprintf("%s account '%s' has been revoked", msg, ev.TypeData.Data.Login)
}

type EventDeleteIntegrationTypeData struct {
	IntegrationID   string  `json:"integration_id"`
	IntegrationType SCMType `json:"integration_type"`
	Data            struct {
		Login     string `json:"login"`
		AvatarURL string `json:"avatar_url"`
		URL       string `json:"url"`
	} `json:"data"`
}

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
