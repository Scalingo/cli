package apps

import (
	"fmt"
)

func List() error {
	apps, err := All()
	if err != nil {
		return err
	}

	fmt.Println("List of your apps :\n")
	for _, app := range apps {
		fmt.Println(app)
	}

	return nil
}
