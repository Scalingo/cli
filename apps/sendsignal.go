package apps

import (
	"context"
	"fmt"
	"strings"

	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-scalingo/v8"
	"github.com/Scalingo/go-utils/errors/v2"
)

func keepUniqueContainersWithType(containers []scalingo.Container, typeName string) (map[string]scalingo.Container, error) {
	containersToKill := map[string]scalingo.Container{}

	hasMatched := false
	prefix := fmt.Sprintf("%s-", typeName)
	for _, container := range containers {
		if strings.HasPrefix(container.Label, prefix) {
			containersToKill[container.Label] = container
			hasMatched = true
		}
	}
	if !hasMatched {
		return containersToKill, errgo.Newf("'%v' did not match any container", typeName)
	}

	return containersToKill, nil
}

func keepUniqueContainersWithNames(containers []scalingo.Container, names []string) map[string]scalingo.Container {
	containersToKill := map[string]scalingo.Container{}

	for _, name := range names {
		for _, container := range containers {
			if container.Label == name {
				containersToKill[name] = container
			}
		}
		if _, ok := containersToKill[name]; !ok {
			containersToKillWithType, err := keepUniqueContainersWithType(containers, name)
			if err != nil {
				io.Error(err.Error())
				continue
			}

			for k, v := range containersToKillWithType {
				containersToKill[k] = v
			}
		}
	}

	return containersToKill
}

func SendSignal(ctx context.Context, appName string, signal string, containerNames []string) error {
	if len(containerNames) == 0 {
		return errgo.New("at least one container name should be given")
	}
	if signal == "" {
		return errgo.New("signal must not be empty")
	}

	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client to send signal to application containers")
	}

	containers, err := c.AppsContainersPs(ctx, appName)
	if err != nil {
		return errgo.Notef(err, "fail to list the application containers to get the ID of the container to send the signal")
	}

	containersToKill := keepUniqueContainersWithNames(containers, containerNames)

	for _, container := range containersToKill {
		err := c.ContainersKill(ctx, appName, signal, container.ID)
		if err != nil {
			rootError := errors.RootCause(err)
			io.Errorf("Fail to send signal to container '%v' because of: %v\n", container.Label, rootError)
			continue
		}
		io.Statusf("Sent signal '%v' to '%v' container.\n", signal, container.Label)
	}
	return nil
}
