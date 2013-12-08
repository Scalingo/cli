package auth

import (
	"appsdeck/config"
	"appsdeck/httpclient"
	"bytes"
	"encoding/json"
	"net/http"
)

func login(email, passwd string) (*http.Response, error) {
	paramsMap := map[string]interface{}{
		"user": map[string]string{
			"email":    email,
			"password": passwd,
		},
	}

	params, _ := json.Marshal(&paramsMap)
	paramsReader := bytes.NewReader(params)
	req, _ := http.NewRequest("POST", config.C["APPSDECK_API"]+"/users/sign_in", paramsReader)

	return httpclient.Do(req)
}
