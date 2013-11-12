package apps

import (
	"appsdeck/api"
	"appsdeck/debug"
	"bytes"
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

	if res.StatusCode == 500 {
		res.Body.Close()
		return nil, fmt.Errorf("Server Internal Error")
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

	if debug.Enable {
		beautifulJson := bytes.NewBuffer(make([]byte, len(buffer)))
		json.Indent(beautifulJson, buffer, "", "  ")
		debug.Println("[API Response]", beautifulJson.String())
	}

	if err := json.Unmarshal(buffer, out); err != nil {
		return err
	}
	return nil
}
