package apps

import (
	"appsdeck/cli/api"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
)

func All() ([]App, error) {
	res, err := api.AppsList()
	if err != nil {
		return nil, err
	}

	if res.StatusCode == 401 {
		res.Body.Close()
		return nil, fmt.Errorf("Unauthorized")
	}

	apps := []App{}
	err = ReadJson(res.Body, &apps)
	return apps, err
}

func ReadJson(body io.ReadCloser, out interface{}) error {
	defer body.Close()
	buffer, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(buffer, out); err != nil {
		return err
	}
	return nil
}
