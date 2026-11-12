package dbdr

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	scalingo "github.com/Scalingo/go-scalingo/v11"
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

func TestFindDatabaseByNameOrIDPrioritizesID(t *testing.T) {
	nameMatch := scalingo.DatabaseNG{ID: "other-id", Name: "shared-identifier"}
	idMatch := scalingo.DatabaseNG{ID: "shared-identifier", Name: "id-match"}

	result, found := findDatabaseByNameOrID(
		[]scalingo.DatabaseNG{nameMatch, idMatch},
		"shared-identifier",
	)

	require.True(t, found)
	assert.Equal(t, idMatch, result)
}

func TestDatabaseAddon(t *testing.T) {
	database := scalingo.DatabaseNG{Database: scalingo.Database{ResourceID: "database-resource-id"}}
	matchingAddon := &scalingo.Addon{ID: "database-addon-id", ResourceID: "database-resource-id"}

	t.Run("uses the only linked addon", func(t *testing.T) {
		result, found := databaseAddon(database, []*scalingo.Addon{matchingAddon})

		require.True(t, found)
		assert.Equal(t, matchingAddon, result)
	})

	t.Run("matches the database resource among multiple addons", func(t *testing.T) {
		result, found := databaseAddon(database, []*scalingo.Addon{
			{ID: "other-addon-id", ResourceID: "other-resource-id"},
			matchingAddon,
		})

		require.True(t, found)
		assert.Equal(t, matchingAddon, result)
	})

	t.Run("reports no matching addon", func(t *testing.T) {
		result, found := databaseAddon(database, []*scalingo.Addon{
			{ID: "other-addon-id", ResourceID: "other-resource-id"},
			{ID: "another-addon-id", ResourceID: "another-resource-id"},
		})

		assert.False(t, found)
		assert.Nil(t, result)
	})
}
