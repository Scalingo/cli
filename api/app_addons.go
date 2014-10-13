package api

import "net/http"

func AddonResourcesList(app string) (*http.Response, error) {
	req := map[string]interface{}{
		"method":   "GET",
		"endpoint": "/apps/" + app + "/addons",
	}
	return Do(req)
}
