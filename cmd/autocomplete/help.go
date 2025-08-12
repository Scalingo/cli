package autocomplete

import (
	"fmt"

	"github.com/urfave/cli/v3"
)

func HelpAutoComplete(c *cli.Context) error {
	for cmd := range c.App.Commands {
		fmt.Println(c.App.Commands[cmd].Name)
	}

	return nil
}
