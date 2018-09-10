package scalingo

import (
	"time"

	"gopkg.in/errgo.v1"
)

type OperationsService interface {
	OperationsShow(app string, opID string) (*Operation, error)
}

var _ OperationsService = (*Client)(nil)

type OperationResponse struct {
	Op Operation `json:"operation"`
}

type Operation struct {
	ID         string    `json:"id"`
	AppID      string    `json:"app_id"`
	CreatedAt  time.Time `json:"created_at"`
	FinishedAt time.Time `json:"finished_at"`
	Status     string    `json:"status"`
	Type       string    `json:"type"`
	Error      string    `json:"error"`
}

func (op *Operation) ElapsedDuration() float64 {
	return op.FinishedAt.Sub(op.CreatedAt).Seconds()
}

func (c *Client) OperationsShow(app string, opID string) (*Operation, error) {
	var opRes OperationResponse
	err := c.subresourceGet(app, "operations", opID, nil, &opRes)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	return &opRes.Op, nil
}
