package dbng

import (
	"testing"

	scalingo "github.com/Scalingo/go-scalingo/v11"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFindDatabaseByNameOrID(t *testing.T) {
	database := scalingo.DatabaseNG{ID: "6acf7671db2f9ad09ad5a692", Name: "binioudbdr2"}
	databases := []scalingo.DatabaseNG{database}

	for _, identifier := range []string{database.ID, database.Name} {
		t.Run(identifier, func(t *testing.T) {
			result, found := findDatabaseByNameOrID(databases, identifier)

			require.True(t, found)
			assert.Equal(t, database, result)
		})
	}
}

func TestFindDatabaseByNameOrIDNotFound(t *testing.T) {
	result, found := findDatabaseByNameOrID(
		[]scalingo.DatabaseNG{{ID: "database-id", Name: "database-name"}},
		"missing",
	)

	assert.False(t, found)
	assert.Empty(t, result)
}
