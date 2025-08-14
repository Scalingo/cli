package autocomplete

import (
	"fmt"

	"github.com/urfave/cli/v3"
)

func HelpAutoComplete(c *cli.Command) error {
	for cmd := range c.Commands {
		fmt.Println(c.Commands[cmd].Name)
	}

	return nil
}
