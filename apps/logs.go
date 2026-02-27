package apps

import (
	"context"
	"strings"

	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/logs"
	"github.com/Scalingo/go-scalingo/v9"
)

func Logs(ctx context.Context, appName string, stream bool, n int, filter string) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	err = checkFilter(ctx, c, appName, filter)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	logsRes, err := c.LogsURL(ctx, appName)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	if err = logs.Dump(ctx, logsRes.LogsURL, n, filter); err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	if stream {
		if err = logs.Stream(ctx, logsRes.LogsURL, filter); err != nil {
			return errgo.Mask(err, errgo.Any)
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
		return errgo.Mask(err)
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
			return errgo.Newf(
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
