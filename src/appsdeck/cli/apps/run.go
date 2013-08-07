package apps

import (
	"appsdeck/cli/api"
	"fmt"
	"io/ioutil"
)

func Run(app string, command []string) error {
	res, err := api.Run(app, command)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	buffer, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(buffer))
	return nil
}
