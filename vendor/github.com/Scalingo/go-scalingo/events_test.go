package scalingo

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/Scalingo/go-scalingo/test"
)

var eventsListCases = map[string]struct {
	PaginationOpts PaginationOpts
	App            string
	Body           string
	Code           int
	EventsCount    int
	ReqURL         string
	ReqMethod      string
	Error          error
}{
	"test get list of events with no event": {
		App:         "app-1",
		Body:        `{"events": [], "meta": {"pagination": {"prev_page": 1, "current_page": 1, "next_page": 1, "total_pages": 1, "total_count": 0}}}`,
		EventsCount: 0,
		ReqURL:      "https://api.scalingo.com/v1/apps/app-1/events?page=0&per_page=0",
		ReqMethod:   "GET",
	},
	"test get list of events with 1 event": {
		App:         "app-1",
		EventsCount: 1,
		Body: `{"events": [{
		"type": "run",
		"type_data": {
			"command": "bundle exec rails console"
		}
	}], "meta": {"pagination": {"prev_page": 1, "current_page": 1, "next_page": 1, "total_pages": 1, "total_count": 0}}}`,
	},
}

func TestEventsList(t *testing.T) {
	for msg, c := range eventsListCases {
		t.Run(msg, func(t *testing.T) {
			hc := test.NewHTTPClient()
			tg := newtokenGeneratorMock()
			tg.setAccessToken("token")
			client := NewClient(ClientConfig{
				TokenGenerator: tg,
			})
			client.backendConfiguration.httpClient = hc

			res := new(http.Response)
			res.Body = ioutil.NopCloser(strings.NewReader(c.Body))
			res.StatusCode = 200
			hc.SetResponseData(res)

			events, _, err := client.EventsList(c.App, c.PaginationOpts)
			if len(hc.Calls) != 1 {
				t.Errorf("expected 1 http request, got %v", len(hc.Calls))
			}
			if len(events) != c.EventsCount {
				t.Errorf("expected %d event, got %v", c.EventsCount, len(events))
			}
			if err != c.Error {
				t.Errorf("expected '%v' error got %v", c.Error, err)
			}

			req := hc.Calls[0]
			if c.ReqURL != "" && req.URL.String() != c.ReqURL {
				t.Errorf("expected request to URL %v, got %v", c.ReqURL, req.URL.String())
			}
			if c.ReqMethod != "" && req.Method != c.ReqMethod {
				t.Errorf("expected request with verb %v, got %v", c.ReqMethod, req.Method)
			}
		})
	}
}
