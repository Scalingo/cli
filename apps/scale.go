package apps

import (
	"encoding/json"
	"fmt"
	"github.com/Scalingo/cli/api"
	"github.com/Scalingo/cli/io"
	"gopkg.in/errgo.v1"
)
import (
	"strconv"
	"strings"
)

type ScaleRes struct {
	ContainerTypes []api.ContainerType `json:"container_types"`
}

func Scale(app string, types []string) error {
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
		scaleParams.Scale = append(scaleParams.Scale, api.ContainerType{Name: typeName, Amount: int(amount)})
	}

	res, err := api.AppsScale(app, scaleParams)
	if err != nil {
		return errgo.Mask(err)
	}
	defer res.Body.Close()

	var scaleRes ScaleRes
	err = json.NewDecoder(res.Body).Decode(&scaleRes)
	if err != nil {
		return errgo.Mask(err)
	}

	fmt.Printf("You application is being scaled to:\n")
	for _, ct := range scaleRes.ContainerTypes {
		fmt.Println(io.Indent(fmt.Sprintf("%s: %d", ct.Name, ct.Amount), 2))
	}
	return nil
}
