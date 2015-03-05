package apps

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/Scalingo/cli/Godeps/_workspace/src/gopkg.in/errgo.v1"
	"github.com/Scalingo/cli/debug"
)

func ReadJson(body io.ReadCloser, out interface{}) error {
	defer body.Close()
	buffer, err := ioutil.ReadAll(body)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	if debug.Enable {
		beautifulJson := bytes.NewBuffer(make([]byte, len(buffer)))
		json.Indent(beautifulJson, buffer, "", "  ")
		debug.Println("[API Response]", beautifulJson.String())
	}

	if err := json.Unmarshal(buffer, out); err != nil {
		return errgo.Mask(err, errgo.Any)
	}
	return nil
}
