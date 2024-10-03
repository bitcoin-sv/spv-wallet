package engine

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"

	"github.com/bitcoin-sv/spv-wallet/engine/datastore"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

const (
	// MetadataField is the field name used for metadata (params)
	MetadataField = "metadata"
)

// Metadata is an object representing the metadata about the related record (standard across all tables)
//
// Gorm related models & indexes: https://gorm.io/docs/models.html - https://gorm.io/docs/indexes.html
type Metadata map[string]interface{}

// XpubMetadata XpubId specific metadata
type XpubMetadata map[string]Metadata

// GormDataType type in gorm
func (m Metadata) GormDataType() string {
	return gormTypeText
}

// Scan scan value into Json, implements sql.Scanner interface
func (m *Metadata) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	byteValue, err := utils.ToByteArray(value)
	if err != nil || bytes.Equal(byteValue, []byte("")) || bytes.Equal(byteValue, []byte("\"\"")) {
		return nil
	}

	err = json.Unmarshal(byteValue, &m)
	return spverrors.Wrapf(err, "failed to parse Metadata from JSON, data: %v", value)
}

// Value return json value, implement driver.Valuer interface
func (m Metadata) Value() (driver.Value, error) {
	if m == nil {
		return nil, nil
	}
	marshal, err := json.Marshal(m)
	if err != nil {
		return nil, spverrors.Wrapf(err, "failed to convert Metadata to JSON, data: %v", m)
	}

	return string(marshal), nil
}

// GormDBDataType the gorm data type for metadata
func (Metadata) GormDBDataType(db *gorm.DB, _ *schema.Field) string {
	if db.Dialector.Name() == datastore.Postgres {
		return datastore.JSONB
	}
	return datastore.JSON
}

// Scan scan value into Json, implements sql.Scanner interface
func (x *XpubMetadata) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	byteValue, err := utils.ToByteArray(value)
	if err != nil || bytes.Equal(byteValue, []byte("")) || bytes.Equal(byteValue, []byte("\"\"")) {
		return nil
	}

	err = json.Unmarshal(byteValue, &x)
	return spverrors.Wrapf(err, "failed to parse XpubMetadata from JSON, data: %v", value)
}

// Value return json value, implement driver.Valuer interface
func (x XpubMetadata) Value() (driver.Value, error) {
	if x == nil {
		return nil, nil
	}
	marshal, err := json.Marshal(x)
	if err != nil {
		return nil, spverrors.Wrapf(err, "failed to convert XpubMetadata to JSON, data: %v", x)
	}

	return string(marshal), nil
}

// GormDBDataType the gorm data type for metadata
func (XpubMetadata) GormDBDataType(db *gorm.DB, _ *schema.Field) string {
	if db.Dialector.Name() == datastore.Postgres {
		return datastore.JSONB
	}
	return datastore.JSON
}
