package api

import (
	"encoding/json"
	"io/ioutil"
	"time"

	"gopkg.in/errgo.v1"
)

type OperationResponse struct {
	Op Operation `json:"operation"`
}

type Operation struct {
	ID         string    `json:"id"`
	CreatedAt  time.Time `json:"created_at"`
	FinishedAt time.Time `json:"finished_at"`
	Status     string    `json:"status"`
	Type       string    `json:"type"`
	Error      string    `json:"error"`
}

func (op *Operation) ElapsedDuration() float64 {
	return op.FinishedAt.Sub(op.CreatedAt).Seconds()
}

func OperationsShow(app string, opID string) (*Operation, error) {
	req := map[string]interface{}{
		"method":   "GET",
		"endpoint": "/apps/" + app + "/operations/" + opID,
		"expected": Statuses{200},
	}
	res, err := Do(req)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	opRes := OperationResponse{}
	err = json.Unmarshal(body, &opRes)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	return &opRes.Op, nil
}
