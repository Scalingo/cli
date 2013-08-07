package apps

import (
	"appsdeck/cli/api"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func All() ([]App, error) {
	res, err := api.AppsList()
	if err != nil {
		return nil, err
	}
	buffer, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 401 {
		return nil, fmt.Errorf("Unauthorized")
	}

	apps := []App{}
	if err := json.Unmarshal(buffer, &apps); err != nil {
		return nil, err
	}
	return apps, nil
}

