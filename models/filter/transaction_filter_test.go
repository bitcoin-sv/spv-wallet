package filter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransactionFilter(t *testing.T) {
	t.Parallel()

	t.Run("default filter", func(t *testing.T) {
		filter := TransactionFilter{}
		dbConditions := filter.ToDbConditions()

		assert.Equal(t, 1, len(dbConditions))
		assert.Nil(t, dbConditions["deleted_at"])
	})

	t.Run("empty filter with include deleted", func(t *testing.T) {
		filter := fromJSON[TransactionFilter](`{
			"include_deleted": true
		}`)
		dbConditions := filter.ToDbConditions()

		assert.Equal(t, 0, len(dbConditions))
	})

	t.Run("with hex", func(t *testing.T) {
		filter := fromJSON[TransactionFilter](`{
			"hex": "test",
			"include_deleted": true
		}`)
		dbConditions := filter.ToDbConditions()

		assert.Equal(t, 1, len(dbConditions))
		assert.Equal(t, "test", dbConditions["hex"])
	})

	t.Run("with block_height", func(t *testing.T) {
		filter := fromJSON[TransactionFilter](`{
			"block_height": 100,
			"include_deleted": true
		}`)
		dbConditions := filter.ToDbConditions()

		assert.Equal(t, 1, len(dbConditions))
		assert.Equal(t, uint64(100), dbConditions["block_height"])
	})
}
