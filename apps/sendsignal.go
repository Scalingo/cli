package apps

import (
	"context"
	"fmt"
	"strings"

	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/go-scalingo/v6"
)

func isContainerInList(containers []scalingo.Container, containerToCheck scalingo.Container) bool {
	for _, container := range containers {
		if container.ID == containerToCheck.ID {
			return true
		}
	}
	return false
}

func keepUniqueContainersWithPrefixes(containers []scalingo.Container, prefixes []string) []scalingo.Container {
	containersToStop := make([]scalingo.Container, 0)

	for _, prefix := range prefixes {
		hasMatched := false
		for _, container := range containers {
			if strings.HasPrefix(container.Label, prefix) && !isContainerInList(containersToStop, container) {
				containersToStop = append(containersToStop, container)
				hasMatched = true
			}
		}
		if !hasMatched {
			fmt.Printf("-----X The prefix '%v' did not match any container\n", prefix)
		}
	}

	return containersToStop
}

func keepUniqueContainersWithNames(containers []scalingo.Container, names []string) []scalingo.Container {
	containersToStop := make([]scalingo.Container, 0)

	for _, name := range names {
		hasMatched := false
		for _, container := range containers {
			if container.Label == name && !isContainerInList(containersToStop, container) {
				containersToStop = append(containersToStop, container)
				hasMatched = true
			}
		}
		if !hasMatched {
			fmt.Printf("-----X The name '%v' did not match any container\n", name)
		}
	}

	return containersToStop
}

func SendSignal(ctx context.Context, appName string, signal string, isPrefixList bool, containerNames []string) error {
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
		return errgo.Notef(err, "fail to list the application containers to get the ID of the container to stop")
	}

	var containersToKill []scalingo.Container
	if isPrefixList {
		containersToKill = keepUniqueContainersWithPrefixes(containers, containerNames)
	} else {
		containersToKill = keepUniqueContainersWithNames(containers, containerNames)
	}

	for _, container := range containersToKill {
		err := c.ContainersKill(ctx, appName, signal, container.ID)
		if err != nil {
			return errgo.Notef(err, "fail to send signal to container")
		}
		fmt.Printf("-----> Sent signal '%v' to '%v' container.\n", signal, container.Label)
	}
	return nil
}
