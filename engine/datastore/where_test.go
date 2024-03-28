package datastore

import (
	"bytes"
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	customtypes "github.com/bitcoin-sv/spv-wallet/engine/datastore/customtypes"
	"github.com/bitcoin-sv/spv-wallet/engine/utils"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func mockDialector(engine Engine) gorm.Dialector {
	mockDb, _, _ := sqlmock.New()
	switch engine {
	case MySQL:
		return mysql.New(mysql.Config{
			Conn:                      mockDb,
			SkipInitializeWithVersion: true,
			DriverName:                "mysql",
		})
	case PostgreSQL:
		return postgres.New(postgres.Config{
			Conn:       mockDb,
			DriverName: "postgres",
		})
	case SQLite:
		return sqlite.Open("file::memory:?cache=shared")
	case MongoDB, Empty:
		// the where builder is not applicable for MongoDB
		return nil
	default:
		return nil
	}
}

func mockClient(engine Engine) (*Client, *gorm.DB) {
	clientInterface, _ := NewClient(context.Background())
	client, _ := clientInterface.(*Client)
	client.options.engine = engine
	dialector := mockDialector(engine)
	gdb, _ := gorm.Open(dialector, &gorm.Config{})
	return client, gdb
}

func makeWhereBuilder(client *Client, gdb *gorm.DB, model interface{}) *whereBuilder {
	return &whereBuilder{
		client: client,
		tx:     gdb.Model(model),
		varNum: 0,
	}
}

const (
	// MetadataField is the field name used for metadata (params)
	MetadataField = "metadata"
)

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

// Test_whereObject test the SQL where selector
func Test_whereObject(t *testing.T) {
	t.Parallel()

	conditions := map[string]interface{}{
		"metadata": Metadata{
			"domain": "test-domain",
		},
	}

	t.Run("Postgres", func(t *testing.T) {
		client, gdb := mockClient(PostgreSQL)

		raw := gdb.ToSQL(func(tx *gorm.DB) *gorm.DB {
			tx, err := ApplyCustomWhere(client, tx, conditions, mockObject{})
			assert.NoError(t, err)
			return tx.First(&mockObject{})
		})

		assert.Contains(t, raw, "metadata::jsonb @>")
		assert.Contains(t, raw, `'{"domain":"test-domain"}'`)
	})

	t.Run("SQLite", func(t *testing.T) {
		client, gdb := mockClient(SQLite)

		raw := gdb.ToSQL(func(tx *gorm.DB) *gorm.DB {
			tx, err := ApplyCustomWhere(client, tx, conditions, mockObject{})
			assert.NoError(t, err)
			return tx.First(&mockObject{})
		})

		assert.Contains(t, raw, "JSON_EXTRACT(metadata")
		assert.Contains(t, raw, `"$.domain"`)
		assert.Contains(t, raw, `"test-domain"`)
	})

	t.Run("MySQL", func(t *testing.T) {
		client, gdb := mockClient(MySQL)

		raw := gdb.ToSQL(func(tx *gorm.DB) *gorm.DB {
			tx, err := ApplyCustomWhere(client, tx, conditions, mockObject{})
			assert.NoError(t, err)
			return tx.First(&mockObject{})
		})

		assert.Contains(t, raw, "JSON_EXTRACT(metadata")
		assert.Contains(t, raw, "'$.domain'")
		assert.Contains(t, raw, "'test-domain'")
	})
}

// Test_whereObject test the SQL where selector
func Test_whereSlice(t *testing.T) {
	t.Parallel()

	conditions := map[string]interface{}{
		"field_in_ids": "test",
	}

	t.Run("Postgres", func(t *testing.T) {
		client, gdb := mockClient(PostgreSQL)
		WithCustomFields([]string{"field_in_ids"}, nil)(client.options)

		raw := gdb.ToSQL(func(tx *gorm.DB) *gorm.DB {
			tx, err := ApplyCustomWhere(client, tx, conditions, mockObject{})
			assert.NoError(t, err)
			return tx.First(&mockObject{})
		})
		// produced SQL:
		// SELECT * FROM "mock_objects" WHERE field_in_ids::jsonb @> '["test"]' ORDER BY "mock_objects"."id" LIMIT 1

		assert.Contains(t, raw, "field_in_ids::jsonb @>")
		assert.Contains(t, raw, `'["test"]'`)
	})

	t.Run("MySQL", func(t *testing.T) {
		client, gdb := mockClient(MySQL)
		WithCustomFields([]string{"field_in_ids"}, nil)(client.options)

		raw := gdb.ToSQL(func(tx *gorm.DB) *gorm.DB {
			tx, err := ApplyCustomWhere(client, tx, conditions, mockObject{})
			assert.NoError(t, err)
			return tx.First(&mockObject{})
		})
		// produced SQL:
		// SELECT * FROM `mock_objects` WHERE JSON_CONTAINS(field_in_ids, CAST('["test"]' AS JSON)) ORDER BY `mock_objects`.`id` LIMIT 1

		assert.Contains(t, raw, "JSON_CONTAINS(field_in_ids")
		assert.Contains(t, raw, `'["test"]'`)
	})

	t.Run("SQLite", func(t *testing.T) {
		client, gdb := mockClient(SQLite)
		WithCustomFields([]string{"field_in_ids"}, nil)(client.options)

		raw := gdb.ToSQL(func(tx *gorm.DB) *gorm.DB {
			tx, err := ApplyCustomWhere(client, tx, conditions, mockObject{})
			assert.NoError(t, err)
			return tx.First(&mockObject{})
		})
		// produced SQL:
		// SELECT * FROM `mock_objects` WHERE EXISTS (SELECT 1 FROM json_each(field_in_ids) WHERE value = "test") ORDER BY `mock_objects`.`id` LIMIT 1

		assert.Contains(t, raw, "json_each(field_in_ids)")
		assert.Contains(t, raw, `"test"`)
	})
}

// Test_processConditions test the SQL where selectors
func Test_processConditions(t *testing.T) {
	t.Parallel()

	theTime := time.Date(2022, 4, 4, 15, 12, 37, 651387237, time.UTC)
	nullTime := sql.NullTime{
		Valid: true,
		Time:  theTime,
	}

	conditions := map[string]interface{}{
		"created_at": map[string]interface{}{
			conditionGreaterThan: customtypes.NullTime{NullTime: nullTime},
		},
		"unique_field_name": map[string]interface{}{
			conditionExists: true,
		},
	}

	t.Run("MySQL", func(t *testing.T) {
		client, gdb := mockClient(MySQL)

		raw := gdb.ToSQL(func(tx *gorm.DB) *gorm.DB {
			tx, err := ApplyCustomWhere(client, tx, conditions, mockObject{})
			assert.NoError(t, err)
			return tx.First(&mockObject{})
		})

		assert.Contains(t, raw, "2022-04-04 15:12:37")
		assert.Contains(t, raw, "AND")
		assert.Regexp(t, "(.+)unique_field_name(.+)IS NOT NULL", raw)
	})

	t.Run("Postgres", func(t *testing.T) {
		client, gdb := mockClient(PostgreSQL)

		raw := gdb.ToSQL(func(tx *gorm.DB) *gorm.DB {
			tx, err := ApplyCustomWhere(client, tx, conditions, mockObject{})
			assert.NoError(t, err)
			return tx.First(&mockObject{})
		})

		assert.Contains(t, raw, "2022-04-04T15:12:37Z")
		assert.Contains(t, raw, "AND")
		assert.Regexp(t, "(.+)unique_field_name(.+)IS NOT NULL", raw)
	})

	t.Run("SQLite", func(t *testing.T) {
		client, gdb := mockClient(SQLite)

		raw := gdb.ToSQL(func(tx *gorm.DB) *gorm.DB {
			tx, err := ApplyCustomWhere(client, tx, conditions, mockObject{})
			assert.NoError(t, err)
			return tx.First(&mockObject{})
		})

		assert.Contains(t, raw, "2022-04-04T15:12:37.651Z")
		assert.Contains(t, raw, "AND")
		assert.Regexp(t, "(.+)unique_field_name(.+)IS NOT NULL", raw)
	})
}

// TestCustomWhere will test the method CustomWhere()
func TestCustomWhere(t *testing.T) {
	t.Parallel()

	t.Run("SQLite empty select", func(t *testing.T) {
		client, gdb := mockClient(SQLite)

		conditions := map[string]interface{}{}

		raw := gdb.ToSQL(func(tx *gorm.DB) *gorm.DB {
			tx, err := ApplyCustomWhere(client, tx, conditions, mockObject{})
			assert.NoError(t, err)
			return tx.First(&mockObject{})
		})

		assert.Regexp(t, "SELECT(.+)FROM(.+)ORDER BY(.+)LIMIT 1", raw)
		assert.NotContains(t, raw, "WHERE")
	})

	t.Run("SQLite simple select", func(t *testing.T) {
		client, gdb := mockClient(SQLite)

		conditions := map[string]interface{}{
			sqlIDFieldProper: "testID",
		}

		raw := gdb.ToSQL(func(tx *gorm.DB) *gorm.DB {
			tx, err := ApplyCustomWhere(client, tx, conditions, mockObject{})
			assert.NoError(t, err)
			return tx.First(&mockObject{})
		})

		assert.Regexp(t, "SELECT(.+)FROM(.+)WHERE(.+)id(.*)\\=(.*)testID(.+)ORDER BY(.+)LIMIT 1", raw)
	})

	t.Run("SQLite $or in json", func(t *testing.T) {
		client, gdb := mockClient(SQLite)
		WithCustomFields([]string{"field_in_ids", "field_out_ids"}, nil)(client.options)

		conditions := map[string]interface{}{
			conditionOr: []map[string]interface{}{{
				"field_in_ids": "value_id",
			}, {
				"field_out_ids": "value_id",
			}},
		}

		raw := gdb.ToSQL(func(tx *gorm.DB) *gorm.DB {
			tx, err := ApplyCustomWhere(client, tx, conditions, mockObject{})
			assert.NoError(t, err)
			return tx.First(&mockObject{})
		})

		assert.Contains(t, raw, "json_each(field_in_ids) WHERE value = \"value_id\"")
		assert.Contains(t, raw, "json_each(field_out_ids) WHERE value = \"value_id\"")
		assert.Regexp(t, "SELECT(.+)FROM(.+)WHERE(.+)EXISTS(.+)SELECT 1(.+)FROM(.+)json_each(.+)WHERE(.+)OR(.+)EXISTS(.+)SELECT 1(.+)FROM(.+)json_each(.+)WHERE(.+)ORDER BY(.+)LIMIT 1", raw)
	})

	t.Run("PostgreSQL $or in json", func(t *testing.T) {
		client, gdb := mockClient(PostgreSQL)
		WithCustomFields([]string{"field_in_ids", "field_out_ids"}, nil)(client.options)

		conditions := map[string]interface{}{
			conditionOr: []map[string]interface{}{{
				"field_in_ids": "value_id",
			}, {
				"field_out_ids": "value_id",
			}},
		}

		raw := gdb.ToSQL(func(tx *gorm.DB) *gorm.DB {
			tx, err := ApplyCustomWhere(client, tx, conditions, mockObject{})
			assert.NoError(t, err)
			return tx.First(&mockObject{})
		})

		assert.Contains(t, raw, "field_in_ids::jsonb @> '[\"value_id\"]'")
		assert.Contains(t, raw, "field_out_ids::jsonb @> '[\"value_id\"]")
		assert.Regexp(t, "SELECT(.+)FROM(.+)WHERE(.+)field_(in|out)_ids(.+)OR(.+)field_(in|out)_ids(.+)ORDER BY(.+)LIMIT 1", raw)
	})

	t.Run("SQLite metadata", func(t *testing.T) {
		client, gdb := mockClient(SQLite)
		conditions := map[string]interface{}{
			metadataField: map[string]interface{}{
				"field_name": "field_value",
			},
		}

		raw := gdb.ToSQL(func(tx *gorm.DB) *gorm.DB {
			tx, err := ApplyCustomWhere(client, tx, conditions, mockObject{})
			assert.NoError(t, err)
			return tx.First(&mockObject{})
		})

		assert.Contains(t, raw, `JSON_EXTRACT(metadata, "$.field_name") = "field_value"`)
	})

	t.Run("MySQL metadata", func(t *testing.T) {
		client, gdb := mockClient(MySQL)
		conditions := map[string]interface{}{
			metadataField: map[string]interface{}{
				"field_name": "field_value",
			},
		}

		raw := gdb.ToSQL(func(tx *gorm.DB) *gorm.DB {
			tx, err := ApplyCustomWhere(client, tx, conditions, mockObject{})
			assert.NoError(t, err)
			return tx.First(&mockObject{})
		})

		assert.Contains(t, raw, "JSON_EXTRACT(metadata, '$.field_name') = 'field_value'")
	})

	t.Run("PostgreSQL metadata", func(t *testing.T) {
		client, gdb := mockClient(PostgreSQL)
		conditions := map[string]interface{}{
			metadataField: Metadata{
				"field_name": "field_value",
			},
		}

		raw := gdb.ToSQL(func(tx *gorm.DB) *gorm.DB {
			tx, err := ApplyCustomWhere(client, tx, conditions, mockObject{})
			assert.NoError(t, err)
			return tx.First(&mockObject{})
		})

		assert.Contains(t, raw, `metadata::jsonb @> '{"field_name":"field_value"}'`)
	})

	t.Run("SQLite $and", func(t *testing.T) {
		client, gdb := mockClient(SQLite)
		WithCustomFields([]string{"field_in_ids", "field_out_ids"}, nil)(client.options)

		conditions := map[string]interface{}{
			conditionAnd: []map[string]interface{}{{
				"reference_id": "reference",
			}, {
				"number": 12,
			}, {
				conditionOr: []map[string]interface{}{{
					"field_in_ids": "value_id",
				}, {
					"field_out_ids": "value_id",
				}},
			}},
		}

		raw := gdb.ToSQL(func(tx *gorm.DB) *gorm.DB {
			tx, err := ApplyCustomWhere(client, tx, conditions, mockObject{})
			assert.NoError(t, err)
			return tx.First(&mockObject{})
		})

		assert.Regexp(t, "reference_id(.*)\\=(.*)reference", raw)
		assert.Regexp(t, "number(.*)\\=(.*)12", raw)
		assert.Regexp(t, "AND(.*)AND", raw)
	})

	t.Run("MySQL $and", func(t *testing.T) {
		client, gdb := mockClient(MySQL)
		WithCustomFields([]string{"field_in_ids", "field_out_ids"}, nil)(client.options)

		conditions := map[string]interface{}{
			conditionAnd: []map[string]interface{}{{
				"reference_id": "reference",
			}, {
				"number": 12,
			}, {
				conditionOr: []map[string]interface{}{{
					"field_in_ids": "value_id",
				}, {
					"field_out_ids": "value_id",
				}},
			}},
		}

		raw := gdb.ToSQL(func(tx *gorm.DB) *gorm.DB {
			tx, err := ApplyCustomWhere(client, tx, conditions, mockObject{})
			assert.NoError(t, err)
			return tx.First(&mockObject{})
		})

		assert.Regexp(t, "reference_id(.*)\\=(.*)reference", raw)
		assert.Regexp(t, "number(.*)\\=(.*)12", raw)
		assert.Regexp(t, "AND(.*)AND", raw)
	})

	t.Run("PostgreSQL $and", func(t *testing.T) {
		client, gdb := mockClient(PostgreSQL)
		WithCustomFields([]string{"field_in_ids", "field_out_ids"}, nil)(client.options)

		conditions := map[string]interface{}{
			conditionAnd: []map[string]interface{}{{
				"reference_id": "reference",
			}, {
				"number": 12,
			}, {
				conditionOr: []map[string]interface{}{{
					"field_in_ids": "value_id",
				}, {
					"field_out_ids": "value_id",
				}},
			}},
		}

		raw := gdb.ToSQL(func(tx *gorm.DB) *gorm.DB {
			tx, err := ApplyCustomWhere(client, tx, conditions, mockObject{})
			assert.NoError(t, err)
			return tx.First(&mockObject{})
		})

		assert.Regexp(t, "reference_id(.*)\\=(.*)reference", raw)
		assert.Regexp(t, "number(.*)\\=(.*)12", raw)
		assert.Regexp(t, "AND(.*)AND", raw)
	})

	t.Run("Where $gt", func(t *testing.T) {
		client, gdb := mockClient(PostgreSQL)

		conditions := map[string]interface{}{
			"number": map[string]interface{}{
				conditionGreaterThan: 502,
			},
		}

		raw := gdb.ToSQL(func(tx *gorm.DB) *gorm.DB {
			tx, err := ApplyCustomWhere(client, tx, conditions, mockObject{})
			assert.NoError(t, err)
			return tx.First(&mockObject{})
		})

		assert.Regexp(t, "number(.*)\\>(.*)502", raw)
	})

	t.Run("Where $and $gt $lt", func(t *testing.T) {
		client, gdb := mockClient(PostgreSQL)

		conditions := map[string]interface{}{
			conditionAnd: []map[string]interface{}{{
				"number": map[string]interface{}{
					conditionLessThan: 503,
				},
			}, {
				"number": map[string]interface{}{
					conditionGreaterThan: 203,
				},
			}},
		}

		raw := gdb.ToSQL(func(tx *gorm.DB) *gorm.DB {
			tx, err := ApplyCustomWhere(client, tx, conditions, mockObject{})
			assert.NoError(t, err)
			return tx.First(&mockObject{})
		})

		// the order may vary
		assert.Regexp(t, "number(.*)\\>(.*)203", raw)
		assert.Regexp(t, "number(.*)\\<(.*)503", raw)

		assert.Regexp(t, "number(.*)[\\d](.*)AND(.*)number(.*)[\\d]", raw)
	})

	t.Run("Where $or $gte $lte", func(t *testing.T) {
		client, gdb := mockClient(PostgreSQL)

		conditions := map[string]interface{}{
			conditionOr: []map[string]interface{}{{
				"number": map[string]interface{}{
					conditionLessThanOrEqual: 503,
				},
			}, {
				"number": map[string]interface{}{
					conditionGreaterThanOrEqual: 203,
				},
			}},
		}

		raw := gdb.ToSQL(func(tx *gorm.DB) *gorm.DB {
			tx, err := ApplyCustomWhere(client, tx, conditions, mockObject{})
			assert.NoError(t, err)
			return tx.First(&mockObject{})
		})

		// the order may vary
		assert.Regexp(t, "number(.*)\\>\\=(.*)203", raw)
		assert.Regexp(t, "number(.*)\\<\\=(.*)503", raw)

		assert.Regexp(t, "number(.*)[\\d](.*)OR(.*)number(.*)[\\d]", raw)
	})

	t.Run("Where $or $and $or $gte $lte", func(t *testing.T) {
		client, gdb := mockClient(PostgreSQL)

		conditions := map[string]interface{}{
			conditionOr: []map[string]interface{}{{
				conditionAnd: []map[string]interface{}{{
					"number": map[string]interface{}{
						conditionLessThanOrEqual: 203,
					},
				}, {
					conditionOr: []map[string]interface{}{{
						"number": map[string]interface{}{
							conditionGreaterThanOrEqual: 1203,
						},
					}, {
						"number": map[string]interface{}{
							conditionGreaterThanOrEqual: 2203,
						},
					}},
				}},
			}, {
				conditionAnd: []map[string]interface{}{{
					"number": map[string]interface{}{
						conditionGreaterThanOrEqual: 3203,
					},
				}, {
					"number": map[string]interface{}{
						conditionGreaterThanOrEqual: 4203,
					},
				}},
			}},
		}

		raw := gdb.ToSQL(func(tx *gorm.DB) *gorm.DB {
			tx, err := ApplyCustomWhere(client, tx, conditions, mockObject{})
			assert.NoError(t, err)
			return tx.First(&mockObject{})
		})

		assert.Regexp(t, "number(.*)\\>\\=(.*)203", raw)
		assert.Regexp(t, "number(.*)\\>\\=(.*)1203", raw)
		assert.Regexp(t, "number(.*)\\<\\=(.*)2203", raw)
		assert.Regexp(t, "number(.*)\\>\\=(.*)3203", raw)
		assert.Regexp(t, "number(.*)\\<\\=(.*)4203", raw)
		assert.Regexp(t, "AND(.+)OR(.+)AND", raw)
	})
}

func Test_sqlInjectionSafety(t *testing.T) {
	t.Parallel()

	t.Run("injection as simple key", func(t *testing.T) {
		client, gdb := mockClient(PostgreSQL)

		conditions := map[string]interface{}{
			"1=1 --": 12,
		}

		gdb.ToSQL(func(tx *gorm.DB) *gorm.DB {
			tx, err := ApplyCustomWhere(client, tx, conditions, mockObject{})
			assert.Error(t, err)
			return tx.First(&mockObject{})
		})
	})

	t.Run("injection in key as conditionExists", func(t *testing.T) {
		client, gdb := mockClient(PostgreSQL)

		conditions := map[string]interface{}{
			"1=1 OR unique_field_name": map[string]interface{}{
				conditionExists: true,
			},
		}

		gdb.ToSQL(func(tx *gorm.DB) *gorm.DB {
			tx, err := ApplyCustomWhere(client, tx, conditions, mockObject{})
			assert.Error(t, err)
			return tx.First(&mockObject{})
		})
	})

	t.Run("injection in metadata", func(t *testing.T) {
		client, gdb := mockClient(PostgreSQL)
		conditions := map[string]interface{}{
			metadataField: Metadata{
				"1=1; DELETE FROM users": "field_value",
			},
		}

		raw := gdb.ToSQL(func(tx *gorm.DB) *gorm.DB {
			tx, err := ApplyCustomWhere(client, tx, conditions, mockObject{})
			assert.NoError(t, err)
			return tx.First(&mockObject{})
		})

		assert.Contains(t, raw, `'{"1=1; DELETE FROM users":"field_value"}'`)
	})
}

// Test_escapeDBString will test the method escapeDBString()
func Test_escapeDBString(t *testing.T) {
	t.Parallel()

	str := escapeDBString(`SELECT * FROM 'table' WHERE 'field'=1;`)
	assert.Equal(t, `SELECT * FROM \'table\' WHERE \'field\'=1;`, str)
}
