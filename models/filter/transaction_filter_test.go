package filter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransactionFilter(t *testing.T) {
	t.Parallel()

	t.Run("default filter", func(t *testing.T) {
		filter := TransactionFilter{}
		dbConditions, _ := filter.ToDbConditions()

		assert.Equal(t, 1, len(dbConditions))
		assert.Nil(t, dbConditions["deleted_at"])
	})

	t.Run("empty filter with include deleted", func(t *testing.T) {
		filter := fromJSON[TransactionFilter](`{
			"include_deleted": true
		}`)
		dbConditions, _ := filter.ToDbConditions()

		assert.Equal(t, 0, len(dbConditions))
	})

	t.Run("with hex", func(t *testing.T) {
		filter := fromJSON[TransactionFilter](`{
			"hex": "test",
			"include_deleted": true
		}`)
		dbConditions, _ := filter.ToDbConditions()

		assert.Equal(t, 1, len(dbConditions))
		assert.Equal(t, "test", dbConditions["hex"])
	})

	t.Run("with block_height", func(t *testing.T) {
		filter := fromJSON[TransactionFilter](`{
			"block_height": 100,
			"include_deleted": true
		}`)
		dbConditions, _ := filter.ToDbConditions()

		assert.Equal(t, 1, len(dbConditions))
		assert.Equal(t, uint64(100), dbConditions["block_height"])
	})

	t.Run("with correct direction", func(t *testing.T) {
		filter := fromJSON[TransactionFilter](`{
			"direction": "incoming",
			"include_deleted": true
		}`)
		dbConditions, _ := filter.ToDbConditions()

		assert.Equal(t, 1, len(dbConditions))
		assert.Equal(t, "incoming", dbConditions["direction"])
	})

	t.Run("with wrong direction", func(t *testing.T) {
		filter := fromJSON[TransactionFilter](`{
			"direction": "wrong_direction",
			"include_deleted": true
		}`)
		_, err := filter.ToDbConditions()

		assert.Error(t, err)
	})
}
