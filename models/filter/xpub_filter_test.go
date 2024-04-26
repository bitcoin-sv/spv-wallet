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
			"include_deleted": true
		}`)
		dbConditions := filter.ToDbConditions()

		assert.Equal(t, 0, len(dbConditions))
	})

	t.Run("with id", func(t *testing.T) {
		filter := fromJSON[XpubFilter](`{
			"id": "test",
			"include_deleted": true
		}`)
		dbConditions := filter.ToDbConditions()

		assert.Equal(t, 1, len(dbConditions))
		assert.Equal(t, "test", dbConditions["id"])
	})

	t.Run("with nextInternalNum", func(t *testing.T) {
		filter := fromJSON[XpubFilter](`{
			"nextInternalNum": 100,
			"include_deleted": true
		}`)
		dbConditions := filter.ToDbConditions()

		assert.Equal(t, 1, len(dbConditions))
		assert.Equal(t, uint32(100), dbConditions["next_internal_num"])
	})
}
