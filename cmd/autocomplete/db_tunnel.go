package autocomplete

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/urfave/cli/v3"
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
)

func DbTunnelAutoComplete(ctx context.Context, c *cli.Command) error {
	appName := CurrentAppCompletion(c)
	if appName == "" {
		return nil
	}

	lastArg := ""
	if len(os.Args) > 2 {
		lastArg = os.Args[len(os.Args)-2]
	}

	if !strings.HasPrefix(lastArg, "-") {
		client, err := config.ScalingoClient(ctx)
		if err != nil {
			return errgo.Notef(err, "fail to get Scalingo client")
		}
		variables, err := client.VariablesList(ctx, appName)
		if err == nil {
			for _, v := range variables {
				if matched, err := regexp.Match("SCALINGO_.*_URL", []byte(v.Name)); matched && err == nil {
					fmt.Println(v.Name)
				}
			}
		}
	}

	return nil
}
