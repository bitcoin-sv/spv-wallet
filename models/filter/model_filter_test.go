package filter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestModelFilter(t *testing.T) {
	t.Parallel()

	t.Run("default filter", func(t *testing.T) {
		filter := ModelFilter{}
		dbConditions := filter.ToDbConditions()

		assert.Equal(t, 1, len(dbConditions))
		assert.Nil(t, dbConditions["deleted_at"])
	})

	t.Run("empty filter with include deleted", func(t *testing.T) {
		filter := ModelFilter{
			IncludeDeleted: ptr(true),
		}
		dbConditions := filter.ToDbConditions()

		assert.Equal(t, 0, len(dbConditions))
	})

	t.Run("with full CreatedRange", func(t *testing.T) {
		filter := ModelFilter{
			CreatedRange: &TimeRange{
				From: ptr(time.Now()),
				To:   ptr(time.Now()),
			},
		}
		dbConditions := filter.ToDbConditions()

		assert.Equal(t, 2, len(dbConditions))
	})

	t.Run("with empty CreatedRange", func(t *testing.T) {
		filter := ModelFilter{
			CreatedRange:   &TimeRange{},
			IncludeDeleted: ptr(true),
		}
		dbConditions := filter.ToDbConditions()

		assert.Equal(t, 0, len(dbConditions))
	})
}
