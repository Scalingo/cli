package scalingo

import (
	"gopkg.in/errgo.v1"
)

type CronTasksService interface {
	CronTasksGet(app string) (CronTasks, error)
}

var _ CronTasksService = (*Client)(nil)

type Job struct {
	Command string `json:"command"`
	Size    string `json:"size,omitempty"`
}

type CronTasks struct {
	Jobs []Job `json:"jobs"`
}

func (c *Client) CronTasksGet(app string) (CronTasks, error) {
	resp := CronTasks{}
	err := c.ScalingoAPI().SubresourceList("apps", app, "cron_tasks", nil, &resp)
	if err != nil {
		return CronTasks{}, errgo.Notef(err, "fail to get cron tasks")
	}
	return resp, nil
}
