package filter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDestinationFilter(t *testing.T) {
	t.Parallel()

	t.Run("default filter", func(t *testing.T) {
		filter := DestinationFilter{}
		dbConditions := filter.ToDbConditions()

		assert.Equal(t, 1, len(dbConditions))
		assert.Nil(t, dbConditions["deleted_at"])
	})

	t.Run("empty filter with include deleted", func(t *testing.T) {
		filter := fromJSON[DestinationFilter](`{
			"includeDeleted": true
		}`)
		dbConditions := filter.ToDbConditions()

		assert.Equal(t, 0, len(dbConditions))
	})

	t.Run("with full CreatedRange", func(t *testing.T) {
		filter := fromJSON[DestinationFilter](`{
			"createdRange": {
				"from": "2024-02-26T11:01:28Z",
				"to": "2024-02-25T11:01:28Z"
			},
			"includeDeleted": true
		}`)

		dbConditions := filter.ToDbConditions()

		assert.Equal(t, 1, len(dbConditions))
		assert.NotNil(t, dbConditions["created_at"])
	})

	t.Run("with empty CreatedRange", func(t *testing.T) {
		filter := fromJSON[DestinationFilter](`{
			"lockingScript": "test",
			"address": "test",
			"draftId": "test",
			"includeDeleted": true
		}`)
		dbConditions := filter.ToDbConditions()

		assert.Equal(t, 3, len(dbConditions))
		assert.NotNil(t, dbConditions["locking_script"])
		assert.NotNil(t, dbConditions["address"])
		assert.NotNil(t, dbConditions["draft_id"])
	})
}
