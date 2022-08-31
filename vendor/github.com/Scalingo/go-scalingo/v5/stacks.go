package scalingo

import (
	"context"
	"time"

	"gopkg.in/errgo.v1"

	httpclient "github.com/Scalingo/go-scalingo/v5/http"
)

type Stack struct {
	ID           string
	CreatedAt    time.Time
	Name         string
	Description  string
	BaseImage    string
	Default      bool
	DeprecatedAt time.Time
}

// This is to properly manage the deprecation date. It is retrieved from the API
// as only the date part (YYYY-MM-DD). Go-lang cannot unmarshal it directly into
// a time.Time, so it is considered a string and converted into a time.Time later
// (cf. func jsonStackToStack)
type jsonStack struct {
	ID           string    `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	BaseImage    string    `json:"base_image"`
	Default      bool      `json:"default"`
	DeprecatedAt string    `json:"deprecated_at"`
}

type StacksService interface {
	StacksList(ctx context.Context) ([]Stack, error)
}

var _ StacksService = (*Client)(nil)

func (c *Client) StacksList(ctx context.Context) ([]Stack, error) {
	req := &httpclient.APIRequest{
		Endpoint: "/features/stacks",
	}

	resmap := map[string][]jsonStack{}
	err := c.ScalingoAPI().DoRequest(ctx, req, &resmap)
	if err != nil {
		return nil, errgo.Notef(err, "fail to request Scalingo API")
	}

	return jsonStackToStack(resmap["stacks"]), nil
}

func jsonStackToStack(s []jsonStack) []Stack {
	var stacks []Stack

	for _, stack := range s {
		deprecationDate, _ := time.Parse("2006-01-02", stack.DeprecatedAt)

		stacks = append(stacks, Stack{
			stack.ID,
			stack.CreatedAt,
			stack.Name,
			stack.Description,
			stack.BaseImage,
			stack.Default,
			deprecationDate,
		})
	}

	return stacks
}

func (s *Stack) IsDeprecated() bool {
	if s.DeprecatedAt.IsZero() {
		return false
	}

	return time.Now().After(s.DeprecatedAt)
}
