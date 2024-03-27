package engine

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/bitcoin-sv/spv-wallet/engine/datastore"
	"github.com/bitcoin-sv/spv-wallet/engine/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// IDs are string ids saved as an array
type IDs []string

// GormDataType type in gorm
func (i IDs) GormDataType() string {
	return gormTypeText
}

// Scan scan value into JSON, implements sql.Scanner interface
func (i *IDs) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	byteValue, err := utils.ToByteArray(value)
	if err != nil {
		return nil
	}

	return json.Unmarshal(byteValue, &i)
}

// Value return json value, implement driver.Valuer interface
func (i IDs) Value() (driver.Value, error) {
	if i == nil {
		return nil, nil
	}
	marshal, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}

	return string(marshal), nil
}

// GormDBDataType the gorm data type for metadata
func (IDs) GormDBDataType(db *gorm.DB, _ *schema.Field) string {
	if db.Dialector.Name() == datastore.Postgres {
		return datastore.JSONB
	}
	return datastore.JSON
}
