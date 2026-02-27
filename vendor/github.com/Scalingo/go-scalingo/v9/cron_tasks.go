package scalingo

import (
	"context"
	"time"

	"github.com/Scalingo/go-utils/errors/v3"
)

type CronTasksService interface {
	CronTasksGet(ctx context.Context, app string) (CronTasks, error)
}

var _ CronTasksService = (*Client)(nil)

type Job struct {
	Command           string    `json:"command"`
	Size              string    `json:"size,omitempty"`
	LastExecutionDate time.Time `json:"last_execution_date,omitempty"`
	NextExecutionDate time.Time `json:"next_execution_date,omitempty"`
}

type CronTasks struct {
	Jobs []Job `json:"jobs"`
}

func (c *Client) CronTasksGet(ctx context.Context, app string) (CronTasks, error) {
	resp := CronTasks{}
	err := c.ScalingoAPI().SubresourceList(ctx, "apps", app, "cron_tasks", nil, &resp)
	if err != nil {
		return CronTasks{}, errors.Wrap(ctx, err, "get cron tasks")
	}
	return resp, nil
}
