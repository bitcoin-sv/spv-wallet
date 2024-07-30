package common

import (
	"database/sql"
	"testing"
	"time"

	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/engine/datastore/customtypes"
	"github.com/stretchr/testify/assert"
)

func TestMapToOldContract_NewlyCreatedRecord(t *testing.T) {
	currentTimestamp := time.Now().UTC()
	engineModel := engine.Model{
		CreatedAt: currentTimestamp,
		UpdatedAt: currentTimestamp,
	}

	commonModel := MapToOldContract(&engineModel)
	assert.Equal(t, engineModel.CreatedAt, commonModel.CreatedAt)
	assert.Equal(t, engineModel.UpdatedAt, commonModel.UpdatedAt)
	assert.Nil(t, commonModel.DeletedAt)
}

func TestMapToOldContract_DeletedAtFieldSet(t *testing.T) {
	currentTimestamp := time.Now().UTC()

	engineModel := engine.Model{
		CreatedAt: currentTimestamp,
		UpdatedAt: currentTimestamp,
		DeletedAt: customtypes.NullTime{
			NullTime: sql.NullTime{
				Time:  currentTimestamp,
				Valid: true,
			},
		},
	}

	commonModel := MapToOldContract(&engineModel)
	assert.Equal(t, engineModel.CreatedAt, commonModel.CreatedAt)
	assert.Equal(t, engineModel.UpdatedAt, commonModel.UpdatedAt)
	assert.Equal(t, engineModel.DeletedAt.Time, *commonModel.DeletedAt)
}
