package apps

import (
	"os"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/utils"
	"github.com/Scalingo/go-scalingo/v4"
	"github.com/olekukonko/tablewriter"
	"gopkg.in/errgo.v1"
)

func Ps(app string) error {
	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client to list the application containers")
	}

	containers, err := c.AppsContainersPs(app)
	if err != nil {
		return errgo.Notef(err, "fail to list the application containers")
	}

	containerTypesSlice, err := c.AppsContainerTypes(app)
	if err != nil {
		return errgo.Notef(err, "fail to list the application container types")
	}
	containerTypes := containerTypesSliceToMap(containerTypesSlice)

	t := tablewriter.NewWriter(os.Stdout)
	t.SetHeader([]string{"Name", "Status", "Command", "Size", "Created At"})

	for _, container := range containers {
		command := container.Command
		if command == "" {
			command = containerTypes[container.Type].Command
		}
		t.Append([]string{container.Label, container.State, command, container.ContainerSize.HumanName, container.CreatedAt.Format(utils.TimeFormat)})
	}
	t.Render()
	return nil
}

// containerTypesSliceToMap takes a slice of container types, and returns a map with the index being the type name
func containerTypesSliceToMap(containerTypesSlice []scalingo.ContainerType) map[string]scalingo.ContainerType {
	containerTypes := map[string]scalingo.ContainerType{}
	for _, containerType := range containerTypesSlice {
		containerTypes[containerType.Name] = containerType
	}
	return containerTypes
}
