package datastore

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/bitcoin-sv/spv-wallet/engine/tester/tgorm"
	"github.com/bitcoin-sv/spv-wallet/engine/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type mockObject struct {
	ID              string
	CreatedAt       time.Time
	UniqueFieldName string
	Number          int
	ReferenceID     string
	Metadata        Metadata
	FieldInIDs      IDs
	FieldOutIDs     IDs
}

func mockClient(engine Engine) (*Client, *gorm.DB) {
	clientInterface, _ := NewClient()
	client, _ := clientInterface.(*Client)
	client.options.engine = engine
	gdb := tgorm.GormDBForPrintingSQL(tgorm.DBType(engine))
	return client, gdb
}

type Metadata map[string]interface{}

func (m Metadata) GormDataType() string {
	return "text"
}

func (m *Metadata) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	byteValue, err := utils.ToByteArray(value)
	if err != nil || bytes.Equal(byteValue, []byte("")) || bytes.Equal(byteValue, []byte("\"\"")) {
		return nil
	}

	return json.Unmarshal(byteValue, &m)
}

func (m Metadata) Value() (driver.Value, error) {
	if m == nil {
		return nil, nil
	}
	marshal, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	return string(marshal), nil
}

func (Metadata) GormDBDataType(db *gorm.DB, _ *schema.Field) string {
	if db.Dialector.Name() == Postgres {
		return JSONB
	}
	return JSON
}

type IDs []string

// GormDataType type in gorm
func (i IDs) GormDataType() string {
	return "text"
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
	if db.Dialector.Name() == Postgres {
		return JSONB
	}
	return JSON
}
