package apps

import (
	"encoding/json"
	"fmt"

	"strconv"
	"strings"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/cli/utils"
	"github.com/Scalingo/go-scalingo"
	"github.com/Scalingo/go-scalingo/debug"
	"github.com/Scalingo/go-scalingo/http"
	"github.com/Scalingo/go-utils/errors"
	"gopkg.in/errgo.v1"
)

type ScaleRes struct {
	Containers []scalingo.ContainerType `json:"containers"`
}

func Scale(app string, sync bool, types []string) error {
	var (
		size        string
		containers  []scalingo.ContainerType
		modificator byte
		err         error
	)

	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}
	scaleParams := &scalingo.AppsScaleParams{}
	typesWithAutoscaler := []string{}
	autoscalers, err := c.AutoscalersList(app)
	if err != nil {
		return errgo.NoteMask(err, "fail to list the autoscalers")
	}

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

		for _, a := range autoscalers {
			if a.ContainerType == typeName {
				typesWithAutoscaler = append(typesWithAutoscaler, typeName)
				break
			}
		}

		newContainerConfig := scalingo.ContainerType{Name: typeName, Size: size}
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

	if len(typesWithAutoscaler) > 0 {
		io.Warning(autoscaleDisableMessage(typesWithAutoscaler))

		var confirm string
		fmt.Scanln(&confirm)
		if confirm != "y" && confirm != "Y" {
			return errgo.New("You didn't confirm, aborting…")
		}
	}

	res, err := c.AppsScale(app, scaleParams)
	if err != nil {
		if !utils.IsPaymentRequiredAndFreeTrialExceededError(err) {
			reqestFailedError, ok := errors.ErrgoRoot(err).(*http.RequestFailedError)
			if !ok {
				return errgo.Notef(err, "request to scale the application failed")
			}

			// In case of unprocessable, format and return a clear error
			if reqestFailedError.Code == 422 {
				return formatContainerTypesError(c, app, reqestFailedError)
			}

			return errgo.Mask(err)
		}
		// If error is Payment Required and user tries to exceed its free trial
		return utils.AskAndStopFreeTrial(c, func() error {
			return Scale(app, sync, types)
		})
	}
	defer res.Body.Close()

	var scaleRes ScaleRes
	err = json.NewDecoder(res.Body).Decode(&scaleRes)
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

func formatContainerTypesError(c *scalingo.Client, app string, requestFailedError *http.RequestFailedError) error {
	containerTypes, err := c.AppsPs(app)
	if err != nil {
		debug.Println(fmt.Sprintf("Failed to get container types: %s", err.Error()))
		return requestFailedError
	}

	if len(containerTypes) == 0 {
		return errgo.New("You have no container type yet.\nPlease refer to the documentation to deploy your application.")
	}

	var containerTypesName string
	for _, containerType := range containerTypes {
		if containerTypesName == "" {
			containerTypesName = "'" + containerType.Name + "'"
			continue
		}
		containerTypesName = containerTypesName + ", '" + containerType.Name + "'"
	}

	errMessage := requestFailedError.APIError.Error() + `

Your available container `
	if len(containerTypes) > 1 {
		errMessage += "types are"
	} else {
		errMessage += "type is"
	}
	errMessage += ": " + containerTypesName + `.

Example of usage:
scalingo --app my-app scale web:2 worker:1
scalingo --app my-app scale web:1:XL
scalingo --app my-app scale web:+1 worker:-1

Use 'scalingo scale --help' for more information`
	return errgo.New(errMessage)
}

func autoscaleDisableMessage(typesWithAutoscaler []string) string {
	if len(typesWithAutoscaler) <= 0 {
		return ""
	}
	msg := "An autoscaler is configured for " + strings.Join(typesWithAutoscaler, ", ") + " container"
	if len(typesWithAutoscaler) > 1 {
		msg += "s"
	}
	msg += ". Manually scaling "
	if len(typesWithAutoscaler) > 1 {
		msg += "them"
	} else {
		msg += "it"
	}
	msg += " will disable the autoscaler. Do you confirm? (y/N)"
	return msg
}
