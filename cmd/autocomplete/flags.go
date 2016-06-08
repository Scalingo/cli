package autocomplete

import (
	"fmt"
	"os"
	"strings"

	"github.com/Scalingo/codegangsta-cli"
)

func getFlagByName(lastArg string, flags []cli.Flag) (bool, cli.Flag) {
	found := false
	var lastFlag cli.Flag

	for _, lastFlag = range flags {
		names := GetFlagNames(lastFlag)
		i := 0
		for i = range names {
			if names[i] == lastArg {
				found = true
				break
			}
		}
		if names[i] == lastArg {
			break
		}
	}

	return found, lastFlag
}

func CountFlags(flags []string) int {
	count := 0

	for i := range os.Args {
		for _, f := range flags {
			if os.Args[i] == f {
				count = count + 1
			}
		}
	}
	return count
}

func GetFlagNames(flag cli.Flag) []string {
	names := strings.Split(cli.GetFlagName(flag), ",")

	for i := range names {
		if i == 0 {
			names[i] = "--" + strings.TrimSpace(names[i])
		} else {
			names[i] = "-" + strings.TrimSpace(names[i])
		}
	}
	return names
}

func DisplayFlags(flags []cli.Flag) {
	lastArg := os.Args[len(os.Args)-2]

	found, lastFlag := getFlagByName(lastArg, flags)

	isBoolFlag := false
	switch lastFlag.(type) {
	case cli.BoolFlag, cli.BoolTFlag:
		isBoolFlag = true && found
	}

	if !strings.HasPrefix(lastArg, "-") || isBoolFlag {
		for _, flag := range flags {
			names := GetFlagNames(flag)

			isSliceFlag := false
			switch flag.(type) {
			case cli.IntSliceFlag, cli.StringSliceFlag:
				isSliceFlag = true
			}
			if CountFlags(names) == 0 || isSliceFlag {
				for i := range names {
					fmt.Println(names[i])
				}
			}
		}
	}
}

func FlagsAutoComplete(c *cli.Context, flag string) bool {
	switch flag {
	case "-r", "--remote":
		return CountFlags([]string{"-r", "--remote"}) == 1 && FlagRemoteAutoComplete(c)
	case "-a", "--app":
		return CountFlags([]string{"-a", "--app"}) == 1 && FlagAppAutoComplete(c)
	}

	return false
}
