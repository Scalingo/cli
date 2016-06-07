package apps

import (
	"fmt"

	"strconv"
	"strings"

	"github.com/Scalingo/go-scalingo"
	"gopkg.in/errgo.v1"
	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/debug"
	"github.com/Scalingo/cli/io"
)

type ScaleRes struct {
	Containers []scalingo.Container `json:"containers"`
}

func Scale(app string, sync bool, types []string) error {
	var (
		size        string
		containers  []scalingo.Container
		modificator byte
		err         error
	)

	scaleParams := &scalingo.AppsScaleParams{}

	for _, t := range types {
		splitT := strings.Split(t, ":")
		if len(splitT) != 2 && len(splitT) != 3 {
			return errgo.Newf("%s is invalid, format is <type>:<amount>[:<size>]", t)
		}
		typeName, typeAmount := splitT[0], splitT[1]
		if len(splitT) == 3 {
			size = splitT[2]
		}

		if typeAmount[0] == '-' || typeAmount[0] == '+' {
			modificator = typeAmount[0]
			typeAmount = typeAmount[1:]
			if size != "" {
				return errgo.Newf("%s is invalid, can't use relative modificator with size, change the size first", t)
			}
			if containers == nil {
				c := config.ScalingoClient()
				containers, err = c.AppsPs(app)
				if err != nil {
					return errgo.Notef(err, "fail to get list of running containers")
				}
				debug.Println("get container list", containers)
			}
		}

		amount, err := strconv.ParseInt(typeAmount, 10, 32)
		if err != nil {
			return errgo.Newf("%s in %s should be an integer", typeAmount, t)
		}

		newContainerConfig := scalingo.Container{Name: typeName, Size: size}
		if modificator != 0 {
			for _, container := range containers {
				if container.Name == typeName {
					if modificator == '-' {
						newContainerConfig.Amount = container.Amount - int(amount)
					} else if modificator == '+' {
						newContainerConfig.Amount = container.Amount + int(amount)
					}
					break
				}
			}
		} else {
			newContainerConfig.Amount = int(amount)
		}

		scaleParams.Containers = append(scaleParams.Containers, newContainerConfig)
	}

	c := config.ScalingoClient()
	res, err := c.AppsScale(app, scaleParams)
	if err != nil {
		return errgo.Mask(err)
	}
	defer res.Body.Close()

	var scaleRes ScaleRes
	err = scalingo.ParseJSON(res, &scaleRes)
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
