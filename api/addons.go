package api

import "net/http"

func AddonsList() (*http.Response, error) {
	req := map[string]interface{}{
		"auth":     false,
		"method":   "GET",
		"endpoint": "/addons",
	}
	return Do(req)
}

func AddonPlansList(addon string) (*http.Response, error) {
	req := map[string]interface{}{
		"auth":     false,
		"method":   "GET",
		"endpoint": "/addons/" + addon + "/plans",
	}
	return Do(req)
}
