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

type ScaleUnprocessableEntity struct {
	Errors map[string]map[string][]string `json:"errors"`
}

func (err ScaleUnprocessableEntity) Error() string {
	var errMsg string
	for typ, errors := range err.Errors {
		errArray := make([]string, 0, len(err.Errors))
		errType := fmt.Sprintf("Container type '%v' is invalid:\n", typ)
		for attr, attrErrs := range errors {
			errArray = append(errArray, fmt.Sprintf("  %s → %s", attr, strings.Join(attrErrs, ", ")))
		}
		errMsg += errType + strings.Join(errArray, "\n")
	}
	return errMsg
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

	if res.StatusCode == 422 {
		var scaleUnprocessableEntity ScaleUnprocessableEntity
		err = api.ParseJSON(res, &scaleUnprocessableEntity)
		if err != nil {
			return errgo.Mask(err)
		}
		return scaleUnprocessableEntity
	}

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
