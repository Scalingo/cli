package apps

import (
	"fmt"
)

func List() error {
	apps, err := All()
	if err != nil {
		return err
	}

	fmt.Printf("List of your apps :\n\n")
	for _, app := range apps {
		fmt.Println(app)
	}

	return nil
}
