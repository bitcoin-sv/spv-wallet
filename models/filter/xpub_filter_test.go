package filter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestXpubFilter(t *testing.T) {
	t.Parallel()

	t.Run("default filter", func(t *testing.T) {
		filter := XpubFilter{}
		dbConditions := filter.ToDbConditions()

		assert.Equal(t, 1, len(dbConditions))
		assert.Nil(t, dbConditions["deleted_at"])
	})

	t.Run("empty filter with include deleted", func(t *testing.T) {
		filter := fromJSON[XpubFilter](`{
			"includeDeleted": true
		}`)
		dbConditions := filter.ToDbConditions()

		assert.Equal(t, 0, len(dbConditions))
	})

	t.Run("with id", func(t *testing.T) {
		filter := fromJSON[XpubFilter](`{
			"id": "test",
			"includeDeleted": true
		}`)
		dbConditions := filter.ToDbConditions()

		assert.Equal(t, 1, len(dbConditions))
		assert.Equal(t, "test", dbConditions["id"])
	})

	t.Run("with currentBalance", func(t *testing.T) {
		filter := fromJSON[XpubFilter](`{
			"currentBalance": 100,
			"includeDeleted": true
		}`)
		dbConditions := filter.ToDbConditions()

		assert.Equal(t, 1, len(dbConditions))
		assert.Equal(t, uint64(100), dbConditions["current_balance"])
	})
}
