package auth

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/Appsdeck/appsdeck/config"
	"github.com/Appsdeck/appsdeck/httpclient"
)

func loginUser(login, passwd string) (*http.Response, error) {
	paramsMap := map[string]interface{}{
		"user": map[string]string{
			"login":    login,
			"password": passwd,
		},
	}

	params, _ := json.Marshal(&paramsMap)
	paramsReader := bytes.NewReader(params)
	req, _ := http.NewRequest("POST", config.C["APPSDECK_API"]+"/users/sign_in", paramsReader)

	return httpclient.Do(req)
}
