package apps

import (
	"context"
	"maps"
	"strings"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-scalingo/v11"
	"github.com/Scalingo/go-utils/errors/v3"
)

func keepUniqueContainersWithType(ctx context.Context, containers []scalingo.Container, typeName string) (map[string]scalingo.Container, error) {
	containersToKill := map[string]scalingo.Container{}

	hasMatched := false
	prefix := typeName + "-"
	for _, container := range containers {
		if strings.HasPrefix(container.Label, prefix) {
			containersToKill[container.Label] = container
			hasMatched = true
		}
	}
	if !hasMatched {
		return containersToKill, errors.Newf(ctx, "'%v' did not match any container", typeName)
	}

	return containersToKill, nil
}

func keepUniqueContainersWithNames(ctx context.Context, containers []scalingo.Container, names []string) map[string]scalingo.Container {
	containersToKill := map[string]scalingo.Container{}

	for _, name := range names {
		for _, container := range containers {
			if container.Label == name {
				containersToKill[name] = container
			}
		}
		if _, ok := containersToKill[name]; !ok {
			containersToKillWithType, err := keepUniqueContainersWithType(ctx, containers, name)
			if err != nil {
				io.Error(err.Error())
				continue
			}

			maps.Copy(containersToKill, containersToKillWithType)
		}
	}

	return containersToKill
}

func SendSignal(ctx context.Context, appName string, signal string, containerNames []string) error {
	if len(containerNames) == 0 {
		return errors.New(ctx, "at least one container name should be given")
	}
	if signal == "" {
		return errors.New(ctx, "signal must not be empty")
	}

	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to get Scalingo client to send signal to application containers")
	}

	containers, err := c.AppsContainersPs(ctx, appName)
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to list the application containers to get the ID of the container to send the signal")
	}

	containersToKill := keepUniqueContainersWithNames(ctx, containers, containerNames)

	for _, container := range containersToKill {
		err := c.ContainersKill(ctx, appName, signal, container.ID)
		if err != nil {
			io.Errorf("Fail to send signal to container '%v' because of: %v\n", container.Label, err)
			continue
		}
		io.Statusf("Sent signal '%v' to '%v' container.\n", signal, container.Label)
	}
	return nil
}
