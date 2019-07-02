package appdetect

import (
	"fmt"
	"os"

	"github.com/Scalingo/go-scalingo/debug"
	"github.com/urfave/cli"
)

func CurrentApp(c *cli.Context) string {
	var appName string
	if c.GlobalString("app") != "<name>" {
		appName = c.GlobalString("app")
	} else if c.String("app") != "<name>" {
		appName = c.String("app")
	} else if os.Getenv("SCALINGO_APP") != "" {
		appName = os.Getenv("SCALINGO_APP")
	} else if dir, ok := DetectGit(); ok {
		appName, _ = ScalingoRepo(dir, c.GlobalString("remote"))
	}
	if appName == "" {
		fmt.Println("Unable to find the application name, please use --app flag.")
		os.Exit(1)
	}

	debug.Println("[AppDetect] App name is", appName)
	return appName
}
