package api

import (
	"net/http"
	"net/url"

	"github.com/Scalingo/cli/Godeps/_workspace/src/gopkg.in/errgo.v1"
)

func LogsURL(app string) (*http.Response, error) {
	req := &APIRequest{
		Endpoint: "/apps/" + app + "/logs",
	}
	return req.Do()
}

func Logs(logsURL string, stream bool, n int) (*http.Response, error) {
	u, err := url.Parse(logsURL)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	req := &APIRequest{
		NoAuth:   true,
		Expected: Statuses{200, 404},
		URL:      u.Scheme + "://" + u.Host,
		Endpoint: u.Path,
		Params: map[string]interface{}{
			"token":  u.Query().Get("token"),
			"stream": stream,
			"n":      n,
		},
	}
	return req.Do()
}
