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
	var size string
	scaleParams := &api.AppsScaleParams{}

	for _, t := range types {
		splitT := strings.Split(t, ":")
		if len(splitT) != 2 && len(splitT) != 3 {
			return errgo.Newf("%s is invalid, format is <type>:<amount>[:<size>]", t)
		}
		typeName, typeAmount := splitT[0], splitT[1]
		if len(splitT) == 3 {
			size = splitT[2]
		}

		amount, err := strconv.ParseInt(typeAmount, 10, 32)
		if err != nil {
			return errgo.Newf("%s in %s should be an integer", typeAmount, t)
		}
		scaleParams.Containers = append(scaleParams.Containers, api.Container{Name: typeName, Amount: int(amount), Size: size})
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

	fmt.Printf("Your application is being scaled to:\n")
	for _, ct := range scaleRes.Containers {
		fmt.Println(io.Indent(fmt.Sprintf("%s: %d - %s", ct.Name, ct.Amount, ct.Size), 2))
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
