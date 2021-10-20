package scalingo

import (
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/go-scalingo/v4/http"
)

type EventCategory struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Position int    `json:"position"`
}

type EventType struct {
	ID          string `json:"id"`
	CategoryID  string `json:"category_id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
	Template    string `json:"template"`
}

func (c *Client) EventTypesList() ([]EventType, error) {
	req := &http.APIRequest{
		Endpoint: "/event_types",
	}

	var res map[string][]EventType
	err := c.ScalingoAPI().DoRequest(req, &res)
	if err != nil {
		return nil, errgo.Notef(err, "fail to make request to Scalingo API")
	}

	return res["event_types"], nil
}

func (c *Client) EventCategoriesList() ([]EventCategory, error) {
	req := &http.APIRequest{
		Endpoint: "/event_categories",
	}

	var res map[string][]EventCategory
	err := c.ScalingoAPI().DoRequest(req, &res)
	if err != nil {
		return nil, errgo.Notef(err, "fail to make request to Scalingo API")
	}

	return res["event_categories"], nil
}
