package apps

import (
	"context"
	"strings"

	"github.com/Scalingo/go-utils/errors/v3"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/logs"
	"github.com/Scalingo/go-scalingo/v10"
)

type LogsRes struct {
	LogsURL string        `json:"logs_url"`
	App     *scalingo.App `json:"app"`
}

func Logs(ctx context.Context, appName string, stream bool, n int, filter string) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to get Scalingo client")
	}

	err = checkFilter(ctx, c, appName, filter)
	if err != nil {
		return errors.Wrap(ctx, err, "check logs filter")
	}

	logsURLRes, err := c.LogsURL(ctx, appName)
	if err != nil {
		return errors.Wrapf(ctx, err, "fetch logs URL for app %s", appName)
	}

	err = logs.Dump(ctx, logsURLRes.LogsURL, n, filter)
	if err != nil {
		return errors.Wrap(ctx, err, "dump application logs")
	}

	if stream {
		err := logs.Stream(ctx, logsURLRes.LogsURL, filter)
		if err != nil {
			return errors.Wrap(ctx, err, "stream application logs")
		}
	}

	return nil
}

func checkFilter(ctx context.Context, c *scalingo.Client, appName string, filter string) error {
	if filter == "" {
		return nil
	}

	if strings.HasPrefix(filter, "one-off-") || strings.HasPrefix(filter, "postdeploy-") {
		return nil
	}

	if filter == "router" {
		return nil
	}

	processes, err := c.AppsContainerTypes(ctx, appName)
	if err != nil {
		return errors.Wrapf(ctx, err, "list container types for app %s", appName)
	}

	filters := strings.Split(filter, "|")
	for _, f := range filters {
		ctName := ""

		for _, ct := range processes {
			ctName = ct.Name
			if strings.HasPrefix(f, ctName+"-") || f == ctName {
				break
			}
		}

		if !strings.HasPrefix(f, ctName+"-") && f != ctName {
			return errors.Newf(ctx,
				"%s is not a valid container filter\n\nEXAMPLES:\n"+
					"\"scalingo logs -F web\": logs of every web containers\n"+
					"\"scalingo logs -F web-1\": logs of web container 1\n"+
					"\"scalingo logs -F router\": only router logs\n"+
					"\"scalingo logs -F web|worker\": logs of every web and worker containers\n",
				f)
		}
	}

	return nil
}
