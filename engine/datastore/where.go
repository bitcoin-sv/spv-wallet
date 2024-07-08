package datastore

import (
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"gorm.io/gorm"
)

// ApplyCustomWhere adds conditions to the gorm db instance
// it returns a tx of type *gorm.DB with a model and conditions applied
func ApplyCustomWhere(client ClientInterface, gdb *gorm.DB, conditions map[string]interface{}, model interface{}) (tx *gorm.DB, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = spverrors.Newf("error processing conditions, %v", r)
		}
	}()

	tx = gdb.Model(model)

	builder := &whereBuilder{
		client: client,
		tx:     tx,
		varNum: 0,
	}

	builder.processConditions(tx, conditions, nil)
	return
}
