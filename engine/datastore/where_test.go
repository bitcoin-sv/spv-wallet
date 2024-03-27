package datastore

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	customtypes "github.com/bitcoin-sv/spv-wallet/engine/datastore/customtypes"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
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

type mockObject struct {
	ID              string
	CreatedAt       time.Time
	UniqueFieldName string
	Number          int
	ReferenceID     string
}

// Test_whereObject test the SQL where selector
func Test_whereSlice(t *testing.T) {
	t.Parallel()

	t.Run("MySQL", func(t *testing.T) {
		client, gdb := mockClient(MySQL)
		builder := makeWhereBuilder(client, gdb, mockObject{})
		query := builder.whereSlice(fieldInIDs, "id_1")
		expected := `JSON_CONTAINS(` + fieldInIDs + `, CAST('["id_1"]' AS JSON))`
		assert.Equal(t, expected, query)
	})

	t.Run("Postgres", func(t *testing.T) {
		client, gdb := mockClient(PostgreSQL)
		builder := makeWhereBuilder(client, gdb, mockObject{})
		query := builder.whereSlice(fieldInIDs, "id_1")
		expected := fieldInIDs + `::jsonb @> '["id_1"]'`
		assert.Equal(t, expected, query)
	})

	t.Run("SQLite", func(t *testing.T) {
		client, gdb := mockClient(SQLite)
		builder := makeWhereBuilder(client, gdb, mockObject{})
		query := builder.whereSlice(fieldInIDs, "id_1")
		expected := `EXISTS (SELECT 1 FROM json_each(` + fieldInIDs + `) WHERE value = "id_1")`
		assert.Equal(t, expected, query)
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

// Test_whereObject test the SQL where selector
func Test_whereObject(t *testing.T) {
	t.Parallel()

	t.Run("MySQL", func(t *testing.T) {
		client, gdb := mockClient(MySQL)
		builder := makeWhereBuilder(client, gdb, mockObject{})

		metadata := map[string]interface{}{
			"test_key": "test-value",
		}
		query := builder.whereObject(metadataField, metadata)
		expected := "JSON_EXTRACT(" + metadataField + ", '$.test_key') = \"test-value\""
		assert.Equal(t, expected, query)

		metadata = map[string]interface{}{
			"test_key": "test-'value'",
		}
		query = builder.whereObject(metadataField, metadata)
		expected = "JSON_EXTRACT(" + metadataField + ", '$.test_key') = \"test-\\'value\\'\""
		assert.Equal(t, expected, query)

		metadata = map[string]interface{}{
			"test_key1": "test-value",
			"test_key2": "test-value2",
		}
		query = builder.whereObject(metadataField, metadata)

		assert.Contains(t, []string{
			"(JSON_EXTRACT(" + metadataField + ", '$.test_key1') = \"test-value\" AND JSON_EXTRACT(" + metadataField + ", '$.test_key2') = \"test-value2\")",
			"(JSON_EXTRACT(" + metadataField + ", '$.test_key2') = \"test-value2\" AND JSON_EXTRACT(" + metadataField + ", '$.test_key1') = \"test-value\")",
		}, query)

		// NOTE: the order of the items can change, hence the query order can change
		// assert.Equal(t, expected, query)

		objectMetadata := map[string]interface{}{
			"testId": map[string]interface{}{
				"test_key1": "test-value",
				"test_key2": "test-value2",
			},
		}
		query = builder.whereObject("object_metadata", objectMetadata)

		assert.Contains(t, []string{
			"(JSON_EXTRACT(object_metadata, '$.testId.test_key1') = \"test-value\" AND JSON_EXTRACT(object_metadata, '$.testId.test_key2') = \"test-value2\")",
			"(JSON_EXTRACT(object_metadata, '$.testId.test_key2') = \"test-value2\" AND JSON_EXTRACT(object_metadata, '$.testId.test_key1') = \"test-value\")",
		}, query)

		// NOTE: the order of the items can change, hence the query order can change
		// assert.Equal(t, expected, query)
	})

	t.Run("Postgres", func(t *testing.T) {
		client, gdb := mockClient(PostgreSQL)
		builder := makeWhereBuilder(client, gdb, mockObject{})

		metadata := map[string]interface{}{
			"test_key": "test-value",
		}
		query := builder.whereObject(metadataField, metadata)
		expected := metadataField + "::jsonb @> '{\"test_key\":\"test-value\"}'::jsonb"
		assert.Equal(t, expected, query)

		metadata = map[string]interface{}{
			"test_key": "test-'value'",
		}
		query = builder.whereObject(metadataField, metadata)
		expected = metadataField + "::jsonb @> '{\"test_key\":\"test-\\'value\\'\"}'::jsonb"
		assert.Equal(t, expected, query)

		metadata = map[string]interface{}{
			"test_key1": "test-value",
			"test_key2": "test-value2",
		}
		query = builder.whereObject(metadataField, metadata)

		assert.Contains(t, []string{
			"(" + metadataField + "::jsonb @> '{\"test_key1\":\"test-value\"}'::jsonb AND " + metadataField + "::jsonb @> '{\"test_key2\":\"test-value2\"}'::jsonb)",
			"(" + metadataField + "::jsonb @> '{\"test_key2\":\"test-value2\"}'::jsonb AND " + metadataField + "::jsonb @> '{\"test_key1\":\"test-value\"}'::jsonb)",
		}, query)

		// NOTE: the order of the items can change, hence the query order can change
		// assert.Equal(t, expected, query)

		objectMetadata := map[string]interface{}{
			"testId": map[string]interface{}{
				"test_key1": "test-value",
				"test_key2": "test-value2",
			},
		}
		query = builder.whereObject("object_metadata", objectMetadata)
		assert.Contains(t, []string{
			"object_metadata::jsonb @> '{\"testId\":{\"test_key1\":\"test-value\",\"test_key2\":\"test-value2\"}}'::jsonb",
			"object_metadata::jsonb @> '{\"testId\":{\"test_key2\":\"test-value2\",\"test_key1\":\"test-value\"}}'::jsonb",
		}, query)

		// NOTE: the order of the items can change, hence the query order can change
		// assert.Equal(t, expected, query)
	})

	t.Run("SQLite", func(t *testing.T) {
		client, gdb := mockClient(SQLite)
		builder := makeWhereBuilder(client, gdb, mockObject{})

		metadata := map[string]interface{}{
			"test_key": "test-value",
		}
		query := builder.whereObject(metadataField, metadata)
		expected := "JSON_EXTRACT(" + metadataField + ", '$.test_key') = \"test-value\""
		assert.Equal(t, expected, query)

		metadata = map[string]interface{}{
			"test_key": "test-'value'",
		}
		query = builder.whereObject(metadataField, metadata)
		expected = "JSON_EXTRACT(" + metadataField + ", '$.test_key') = \"test-\\'value\\'\""
		assert.Equal(t, expected, query)

		metadata = map[string]interface{}{
			"test_key1": "test-value",
			"test_key2": "test-value2",
		}
		query = builder.whereObject(metadataField, metadata)
		assert.Contains(t, []string{
			"(JSON_EXTRACT(" + metadataField + ", '$.test_key1') = \"test-value\" AND JSON_EXTRACT(" + metadataField + ", '$.test_key2') = \"test-value2\")",
			"(JSON_EXTRACT(" + metadataField + ", '$.test_key2') = \"test-value2\" AND JSON_EXTRACT(" + metadataField + ", '$.test_key1') = \"test-value\")",
		}, query)

		// NOTE: the order of the items can change, hence the query order can change
		// assert.Equal(t, expected, query)

		objectMetadata := map[string]interface{}{
			"testId": map[string]interface{}{
				"test_key1": "test-value",
				"test_key2": "test-value2",
			},
		}
		query = builder.whereObject("object_metadata", objectMetadata)
		assert.Contains(t, []string{
			"(JSON_EXTRACT(object_metadata, '$.testId.test_key1') = \"test-value\" AND JSON_EXTRACT(object_metadata, '$.testId.test_key2') = \"test-value2\")",
			"(JSON_EXTRACT(object_metadata, '$.testId.test_key2') = \"test-value2\" AND JSON_EXTRACT(object_metadata, '$.testId.test_key1') = \"test-value\")",
		}, query)
		// NOTE: the order of the items can change, hence the query order can change
		// assert.Equal(t, expected, query)
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
		arrayField1 := fieldInIDs
		arrayField2 := fieldOutIDs

		client, gdb := mockClient(SQLite)
		WithCustomFields([]string{arrayField1, arrayField2}, nil)(client.options)

		conditions := map[string]interface{}{
			conditionOr: []map[string]interface{}{{
				arrayField1: "value_id",
			}, {
				arrayField2: "value_id",
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
		arrayField1 := fieldInIDs
		arrayField2 := fieldOutIDs

		client, gdb := mockClient(PostgreSQL)
		WithCustomFields([]string{arrayField1, arrayField2}, nil)(client.options)

		conditions := map[string]interface{}{
			conditionOr: []map[string]interface{}{{
				arrayField1: "value_id",
			}, {
				arrayField2: "value_id",
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

		assert.Contains(t, raw, "JSON_EXTRACT(metadata, '$.field_name') = \"field_value\"")
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

		assert.Contains(t, raw, "JSON_EXTRACT(metadata, '$.field_name') = \"field_value\"")
	})

	t.Run("PostgreSQL metadata", func(t *testing.T) {
		client, gdb := mockClient(PostgreSQL)
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

		assert.Contains(t, raw, "metadata::jsonb @> '{\"field_name\":\"field_value\"}'::jsonb")
	})

	t.Run("SQLite $and", func(t *testing.T) {
		arrayField1 := fieldInIDs
		arrayField2 := fieldOutIDs

		client, gdb := mockClient(SQLite)
		WithCustomFields([]string{arrayField1, arrayField2}, nil)(client.options)

		conditions := map[string]interface{}{
			conditionAnd: []map[string]interface{}{{
				"reference_id": "reference",
			}, {
				"number": 12,
			}, {
				conditionOr: []map[string]interface{}{{
					arrayField1: "value_id",
				}, {
					arrayField2: "value_id",
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
		arrayField1 := fieldInIDs
		arrayField2 := fieldOutIDs

		client, gdb := mockClient(MySQL)
		WithCustomFields([]string{arrayField1, arrayField2}, nil)(client.options)

		conditions := map[string]interface{}{
			conditionAnd: []map[string]interface{}{{
				"reference_id": "reference",
			}, {
				"number": 12,
			}, {
				conditionOr: []map[string]interface{}{{
					arrayField1: "value_id",
				}, {
					arrayField2: "value_id",
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
		arrayField1 := fieldInIDs
		arrayField2 := fieldOutIDs

		client, gdb := mockClient(PostgreSQL)
		WithCustomFields([]string{arrayField1, arrayField2}, nil)(client.options)

		conditions := map[string]interface{}{
			conditionAnd: []map[string]interface{}{{
				"reference_id": "reference",
			}, {
				"number": 12,
			}, {
				conditionOr: []map[string]interface{}{{
					arrayField1: "value_id",
				}, {
					arrayField2: "value_id",
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
			metadataField: map[string]interface{}{
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
