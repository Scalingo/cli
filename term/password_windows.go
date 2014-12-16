package term

import (
	"fmt"
)

func Password(prompt string) (string, error) {
	fmt.Print(prompt)
	var pass string
	_, err := fmt.Scanln(&pass)
	return pass, err
}
