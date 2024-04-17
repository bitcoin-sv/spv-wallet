package engine

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"

	"github.com/bitcoin-sv/spv-wallet/engine/datastore"
	"github.com/bitcoin-sv/spv-wallet/engine/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
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

	return json.Unmarshal(byteValue, &m)
}

// Value return json value, implement driver.Valuer interface
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

	return json.Unmarshal(byteValue, &x)
}

// Value return json value, implement driver.Valuer interface
func (x XpubMetadata) Value() (driver.Value, error) {
	if x == nil {
		return nil, nil
	}
	marshal, err := json.Marshal(x)
	if err != nil {
		return nil, err
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

// MarshalBSONValue method is called by bson.Marshal in Mongo for type = Metadata
func (m *Metadata) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if m == nil || len(*m) == 0 {
		return bson.TypeNull, nil, nil
	}

	metadata := make([]map[string]interface{}, 0)
	for key, value := range *m {
		metadata = append(metadata, map[string]interface{}{
			"k": key,
			"v": value,
		})
	}

	return bson.MarshalValue(metadata)
}

// UnmarshalBSONValue method is called by bson.Unmarshal in Mongo for type = Metadata
func (m *Metadata) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	raw := bson.RawValue{Type: t, Value: data}

	if raw.Value == nil {
		return nil
	}

	var uMap []map[string]interface{}
	if err := raw.Unmarshal(&uMap); err != nil {
		return err
	}

	*m = make(Metadata)
	for _, meta := range uMap {
		key := meta["k"].(string)
		(*m)[key] = meta["v"]
	}

	return nil
}

// MarshalBSONValue method is called by bson.Marshal in Mongo for type = XpubMetadata
func (x *XpubMetadata) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if x == nil || len(*x) == 0 {
		return bson.TypeNull, nil, nil
	}

	metadata := make([]map[string]interface{}, 0)
	for xPubKey, meta := range *x {
		for key, value := range meta {
			metadata = append(metadata, map[string]interface{}{
				"x": xPubKey,
				"k": key,
				"v": value,
			})
		}
	}

	return bson.MarshalValue(metadata)
}

// UnmarshalBSONValue method is called by bson.Unmarshal in Mongo for type = XpubMetadata
func (x *XpubMetadata) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	raw := bson.RawValue{Type: t, Value: data}

	if raw.Value == nil {
		return nil
	}

	var uMap []map[string]interface{}
	if err := raw.Unmarshal(&uMap); err != nil {
		return err
	}

	*x = make(XpubMetadata)
	for _, meta := range uMap {
		xPub := meta["x"].(string)
		key := meta["k"].(string)
		if (*x)[xPub] == nil {
			(*x)[xPub] = make(Metadata)
		}
		(*x)[xPub][key] = meta["v"]
	}

	return nil
}
