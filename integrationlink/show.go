package integrationlink

import (
	"context"
	"os"
	"text/template"

	"github.com/fatih/color"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/cli/utils"
	scalingo "github.com/Scalingo/go-scalingo/v6"
	"github.com/Scalingo/go-scalingo/v6/http"
	"github.com/Scalingo/go-utils/errors/v2"
)

const integrationSettingsTemplateText = `
{{- title "Application:"}} {{ .AppName }} ({{ .AppID }})
{{title "Integration:"}} {{ .SCMType }} ({{ .IntegrationID }})
{{title "Linker:"}} {{.LinkerUsername}}

{{title "Repository:"}} {{ .RepositoryOwner }}/{{ .RepositoryName }}
{{title "Automatic deployment:"}} {{ enabledIcon .AutoDeploy }}{{ if .AutoDeploy }}, {{ .Branch }}{{ end }}
{{title "Automatic deployment of review apps:"}} {{ enabledIcon .AutoDeployReviewApps }}
{{- if .AutoDeployReviewApps }}
	{{title "Allow creation from forks:"}} {{ enabledIcon .ReviewAppsFromForks }}
	{{title "Automatic destroy on close:"}} {{ enabledIcon .DestroyOnClose }}{{ if .DestroyOnClose }}, {{ if eq .HoursBeforeDestroyOnClose 0 -}} instantly {{- else -}} {{ .HoursBeforeDestroyOnClose }}h{{ end }}{{ end }}
	{{title "Automatic destroy on stale:"}} {{ enabledIcon .DestroyOnStale }}{{ if .DestroyOnStale }}, {{ if eq .HoursBeforeDestroyOnStale 0 -}} instantly {{- else -}} {{ .HoursBeforeDestroyOnStale }}h{{ end }}{{ end }}
{{- end}}
`

type IntegrationSettings struct {
	AppName, AppID, SCMType, IntegrationID, LinkerUsername, RepositoryOwner, RepositoryName, Branch string
	AutoDeploy, AutoDeployReviewApps, ReviewAppsFromForks, DestroyOnClose, DestroyOnStale           bool
	HoursBeforeDestroyOnClose, HoursBeforeDestroyOnStale                                            uint
}

func title(textToFormat string) string {
	return color.New(color.FgYellow).Sprint(textToFormat)
}

func enabledIcon(enabled bool) string {
	if enabled {
		return color.GreenString(utils.Success)
	}

	return color.RedString(utils.Error)
}

func Show(ctx context.Context, app string) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Notef(ctx, err, "fail to get Scalingo client")
	}

	repoLink, err := c.SCMRepoLinkShow(ctx, app)
	if scerr, ok := errors.RootCause(err).(*http.RequestFailedError); ok && scerr.Code == 404 {
		io.Statusf("Your app '%s' has no integration link.\n", app)
		return nil
	}
	if err != nil {
		return errors.Notef(ctx, err, "fail to get integration link for this app")
	}

	templater, err := template.New("integrationSettings").
		Funcs(template.FuncMap{
			"enabledIcon": enabledIcon,
			"title":       title,
		}).
		Parse(integrationSettingsTemplateText)

	if err != nil {
		return errors.Notef(ctx, err, "invalid integration settings template")
	}

	settings := IntegrationSettings{
		AppName:                   app,
		AppID:                     repoLink.AppID,
		SCMType:                   scalingo.SCMTypeDisplay[repoLink.SCMType],
		IntegrationID:             repoLink.AuthIntegrationUUID,
		LinkerUsername:            repoLink.Linker.Username,
		RepositoryOwner:           repoLink.Owner,
		RepositoryName:            repoLink.Repo,
		Branch:                    repoLink.Branch,
		AutoDeploy:                repoLink.AutoDeployEnabled,
		AutoDeployReviewApps:      repoLink.DeployReviewAppsEnabled,
		ReviewAppsFromForks:       repoLink.AutomaticCreationFromForksAllowed,
		DestroyOnClose:            repoLink.DeleteOnCloseEnabled,
		HoursBeforeDestroyOnClose: repoLink.HoursBeforeDeleteOnClose,
		DestroyOnStale:            repoLink.DeleteStaleEnabled,
		HoursBeforeDestroyOnStale: repoLink.HoursBeforeDeleteStale,
	}

	err = templater.Execute(os.Stdout, settings)

	if err != nil {
		return errors.Notef(ctx, err, "failed rendering integration settings")
	}

	return nil
}
