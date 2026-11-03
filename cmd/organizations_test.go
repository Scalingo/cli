package cmd

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrganizationsCommand(t *testing.T) {
	assert.True(t, slices.Contains(globalCommands, &organizationsCommand))
	assert.False(t, slices.Contains(regionalCommands, &organizationsCommand))

	assert.True(t, organizationsCommand.HasName("orgs"))
	assert.True(t, organizationsCommand.HasName("organizations"))
}
