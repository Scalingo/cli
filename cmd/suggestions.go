package cmd

import (
	"fmt"
	"strings"

	"github.com/Scalingo/cli/Godeps/_workspace/src/github.com/Scalingo/codegangsta-cli"
)

func ShowSuggestions(c *cli.Context) {
	var suggestions []string
	cmdName := c.Args().First()

	startRange := 3
	if 3 >= len(cmdName) {
		startRange = len(cmdName) % 4
	}
	endRange := len(cmdName) - startRange

	if startRange >= 0 {
		for _, cmd := range c.App.Commands {
			if strings.HasPrefix(cmd.Name, cmdName[:startRange]) {
				suggestions = append(suggestions, cmd.Name)
			} else if strings.HasSuffix(cmd.Name, cmdName[endRange:]) {
				suggestions = append(suggestions, cmd.Name)
			}
		}
	}

	if len(c.Args()) > 1 {
		alikeCmd := ""
		matches := 0
		tmpMatches := 0
		for _, cmd := range c.App.Commands {
			for _, a := range c.Args() {
				for _, f := range cmd.Flags {
					for _, s := range strings.Split(f.String(), ", ") {
						if s == a {
							tmpMatches += 1
						}
					}
				}
			}
			if tmpMatches > matches {
				alikeCmd = cmd.Name
				matches = tmpMatches
			}
			tmpMatches = 0
		}

		if alikeCmd != "" {
			fullCmd := "scalingo " + alikeCmd + " " + strings.Join(c.Args()[1:], " ")
			suggestions = append(suggestions, fullCmd)
		}
	}

	if len(suggestions) > 0 {
		fmt.Println("You might be looking for:")
		for _, s := range suggestions {
			fmt.Printf("  - '%s'\n", s)
		}
	}
}
