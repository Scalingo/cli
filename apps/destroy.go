package apps

import (
	"github.com/Appsdeck/appsdeck/api"
	"fmt"
)

func Destroy(id string) {
	res, err := api.AppsDestroy(id)
	if err != nil {
		fmt.Println("Fail to create app:", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		fmt.Printf("App identified by %v has not been found\n", id)
	} else if res.StatusCode == 204 {
		fmt.Printf("App %s has been deleted\n", id)
	}
}
