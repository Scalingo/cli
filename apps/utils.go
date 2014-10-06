package apps

import (
	"github.com/Scalingo/cli/debug"
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
)

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
