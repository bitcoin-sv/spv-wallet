package filter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPaymailFilter(t *testing.T) {
	t.Parallel()

	t.Run("default filter", func(t *testing.T) {
		filter := PaymailFilter{}
		dbConditions := filter.ToDbConditions()

		assert.Equal(t, 1, len(dbConditions))
		assert.Nil(t, dbConditions["deleted_at"])
	})

	t.Run("empty filter with include deleted", func(t *testing.T) {
		filter := fromJSON[PaymailFilter](`{
			"include_deleted": true
		}`)
		dbConditions := filter.ToDbConditions()

		assert.Equal(t, 0, len(dbConditions))
	})

	t.Run("with alias", func(t *testing.T) {
		filter := fromJSON[PaymailFilter](`{
			"alias": "example",
			"include_deleted": true
		}`)
		dbConditions := filter.ToDbConditions()

		assert.Equal(t, 1, len(dbConditions))
		assert.Equal(t, "example", dbConditions["alias"])
	})

	t.Run("with publicName", func(t *testing.T) {
		filter := fromJSON[PaymailFilter](`{
			"publicName": "pubName",
			"include_deleted": true
		}`)
		dbConditions := filter.ToDbConditions()

		assert.Equal(t, 1, len(dbConditions))
		assert.Equal(t, "pubName", dbConditions["public_name"])
	})
}
