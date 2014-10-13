package apps

import "github.com/Scalingo/cli/api"

func Scale() error {
	res, _ := api.AppsList()
	defer res.Body.Close()

	return nil
}
