package scalingo

import (
	"context"
	"time"

	"gopkg.in/errgo.v1"

	httpclient "github.com/Scalingo/go-scalingo/v8/http"
)

type OperationsService interface {
	OperationsShow(ctx context.Context, app, opID string) (*Operation, error)
}

var _ OperationsService = (*Client)(nil)

type OperationStatus string

const (
	OperationStatusPending OperationStatus = "pending"
	OperationStatusDone    OperationStatus = "done"
	OperationStatusRunning OperationStatus = "running"
	OperationStatusError   OperationStatus = "error"
)

type OperationType string

const (
	OperationTypeScale       OperationType = "scale"
	OperationTypeStart       OperationType = "start"
	OperationTypeStartOneOff OperationType = "start-one-off"
)

type OperationResponse struct {
	Op Operation `json:"operation"`
}

type OperationStartOneOffData struct {
	AttachURL   string `json:"attach_url"`
	ContainerID string `json:"container_id"`
}

type Operation struct {
	ID              string                   `json:"id"`
	AppID           string                   `json:"app_id"`
	CreatedAt       time.Time                `json:"created_at"`
	FinishedAt      time.Time                `json:"finished_at"`
	Status          OperationStatus          `json:"status"`
	Type            OperationType            `json:"type"`
	Error           string                   `json:"error"`
	StartOneOffData OperationStartOneOffData `json:"start_one_off_data"`
}

func (op *Operation) ElapsedDuration() float64 {
	return op.FinishedAt.Sub(op.CreatedAt).Seconds()
}

func (c *Client) OperationsShowFromURL(ctx context.Context, url string) (*Operation, error) {
	var opRes OperationResponse
	err := c.ScalingoAPI().DoRequest(ctx, &httpclient.APIRequest{
		Method: "GET", URL: url,
	}, &opRes)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	return &opRes.Op, nil
}

func (c *Client) OperationsShow(ctx context.Context, app, opID string) (*Operation, error) {
	var opRes OperationResponse
	err := c.ScalingoAPI().SubresourceGet(ctx, "apps", app, "operations", opID, nil, &opRes)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	return &opRes.Op, nil
}
