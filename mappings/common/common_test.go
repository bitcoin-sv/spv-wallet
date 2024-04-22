package common

import (
	"database/sql"
	"fmt"
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/engine/datastore/customtypes"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestMapToContract_NewlyCreatedRecord(t *testing.T) {
	currentTimestamp := time.Now().UTC()
	engineModel := engine.Model{
		CreatedAt: currentTimestamp,
		UpdatedAt: currentTimestamp,
	}

	commonModel := MapToContract(&engineModel)
	fmt.Printf("Struct: %+v", commonModel)
	assert.Equal(t, engineModel.CreatedAt, commonModel.CreatedAt)
	assert.Equal(t, engineModel.UpdatedAt, commonModel.UpdatedAt)
	assert.Nil(t, commonModel.DeletedAt)
}

func TestMapToContract_DeletedAtFieldSet(t *testing.T) {
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

	commonModel := MapToContract(&engineModel)
	fmt.Printf("Struct: %+v", commonModel)
	assert.Equal(t, engineModel.CreatedAt, commonModel.CreatedAt)
	assert.Equal(t, engineModel.UpdatedAt, commonModel.UpdatedAt)
	assert.Equal(t, engineModel.DeletedAt.Time, *commonModel.DeletedAt)
}
