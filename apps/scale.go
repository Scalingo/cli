package apps

import (
	"fmt"

	"github.com/Scalingo/cli/Godeps/_workspace/src/gopkg.in/errgo.v1"
	"github.com/Scalingo/cli/api"
	"github.com/Scalingo/cli/io"
)
import (
	"strconv"
	"strings"
)

type ScaleRes struct {
	Containers []api.Container `json:"containers"`
}

func Scale(app string, sync bool, types []string) error {
	scaleParams := &api.AppsScaleParams{}

	for _, t := range types {
		splitT := strings.Split(t, ":")
		if len(splitT) != 2 {
			return errgo.Newf("%s is invalid, format is <type>:<amount>", t)
		}
		typeName, typeAmount := splitT[0], splitT[1]
		amount, err := strconv.ParseInt(typeAmount, 10, 32)
		if err != nil {
			return errgo.Newf("%s in %s should be an integer", typeAmount, t)
		}
		scaleParams.Containers = append(scaleParams.Containers, api.Container{Name: typeName, Amount: int(amount)})
	}

	res, err := api.AppsScale(app, scaleParams)
	if err != nil {
		return errgo.Mask(err)
	}
	defer res.Body.Close()

	var scaleRes ScaleRes
	err = api.ParseJSON(res, &scaleRes)
	if err != nil {
		return errgo.Mask(err)
	}

	fmt.Printf("You application is being scaled to:\n")
	for _, ct := range scaleRes.Containers {
		fmt.Println(io.Indent(fmt.Sprintf("%s: %d", ct.Name, ct.Amount), 2))
	}

	if !sync {
		return nil
	}

	err = handleOperation(app, res)
	if err != nil {
		return errgo.Mask(err)
	}

	fmt.Println("Your application has been scaled.")
	return nil
}
