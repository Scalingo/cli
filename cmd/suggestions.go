package cmd

import (
	"fmt"
	"strings"

	"github.com/Scalingo/codegangsta-cli"
)

func ShowSuggestions(c *cli.Context) {
	var suggestions []string
	cmdName := c.Args().First()

	startRange := 2
	if 2 >= len(cmdName) {
		startRange = len(cmdName) % 3
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

	if len(c.Args()) > 0 && len(suggestions) > 0 {
		fmt.Println("You might be looking for:")
		for _, s := range suggestions {
			fmt.Printf("  - '%s'\n", s)
		}
	}
}
