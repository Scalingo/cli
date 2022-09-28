package scalingo

import (
	"context"
	"time"

	"gopkg.in/errgo.v1"

	httpclient "github.com/Scalingo/go-scalingo/v6/http"
)

type DeprecationDate struct {
	time.Time
}

type Stack struct {
	ID           string          `json:"id"`
	CreatedAt    time.Time       `json:"created_at"`
	Name         string          `json:"name"`
	Description  string          `json:"description"`
	BaseImage    string          `json:"base_image"`
	Default      bool            `json:"default"`
	DeprecatedAt DeprecationDate `json:"deprecated_at,omitempty"`
}

type StacksService interface {
	StacksList(ctx context.Context) ([]Stack, error)
}

var _ StacksService = (*Client)(nil)

func (c *Client) StacksList(ctx context.Context) ([]Stack, error) {
	req := &httpclient.APIRequest{
		Endpoint: "/features/stacks",
	}

	resmap := map[string][]Stack{}
	err := c.ScalingoAPI().DoRequest(ctx, req, &resmap)
	if err != nil {
		return nil, errgo.Notef(err, "fail to request Scalingo API")
	}

	return resmap["stacks"], nil
}

// The regional API returns a date formatted as "2006-01-02"
// Go standard library does not unmarshal that format
func (deprecationDate *DeprecationDate) UnmarshalJSON(b []byte) error {
	s := string(b)

	if s == "null" {
		// When there is no deprecation date for a stack, the json will look like: {"deprecated_at": null}
		return nil
	}

	t, err := time.Parse(`"2006-01-02"`, s)
	if err != nil {
		return err
	}

	deprecationDate.Time = t
	return nil
}

func (s *Stack) IsDeprecated() bool {
	if s.DeprecatedAt.IsZero() {
		return false
	}

	return time.Now().After(s.DeprecatedAt.Time)
}
