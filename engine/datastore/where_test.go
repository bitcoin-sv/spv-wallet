package datastore

import (
	"database/sql"
	"testing"
	"time"

	"github.com/bitcoin-sv/spv-wallet/engine/datastore/customtypes"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

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
}

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
	t.Run("SQLite nested metadata", func(t *testing.T) {
		client, gdb := mockClient(SQLite)
		conditions := Metadata{
			metadataField: Metadata{
				"p2p_tx_metadata": Metadata{
					"note": "test",
				},
			},
		}

		raw := gdb.ToSQL(func(tx *gorm.DB) *gorm.DB {
			tx, err := ApplyCustomWhere(client, tx, conditions, mockObject{})
			assert.NoError(t, err)
			return tx.First(&mockObject{})
		})

		assert.Contains(t, raw, `JSON_EXTRACT(metadata, "$.p2p_tx_metadata.note") = "test"`)
	})

	t.Run("PostgreSQL nested metadata", func(t *testing.T) {
		client, gdb := mockClient(PostgreSQL)
		conditions := Metadata{
			metadataField: Metadata{
				"p2p_tx_metadata": Metadata{
					"note": "test",
				},
			},
		}

		raw := gdb.ToSQL(func(tx *gorm.DB) *gorm.DB {
			tx, err := ApplyCustomWhere(client, tx, conditions, mockObject{})
			assert.NoError(t, err)
			return tx.First(&mockObject{})
		})

		assert.Contains(t, raw, `metadata::jsonb @> '{"p2p_tx_metadata":{"note":"test"}}'`)
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

func Test_isEmptyCondition(t *testing.T) {
	t.Parallel()

	t.Run("nil ptr to map", func(t *testing.T) {
		var condition *map[string]interface{}
		assert.True(t, isEmptyCondition(condition))
	})

	t.Run("not-nil ptr to nil-map", func(t *testing.T) {
		var theMap map[string]interface{}
		condition := &theMap
		assert.True(t, isEmptyCondition(condition))
	})

	t.Run("not-nil ptr to empty map", func(t *testing.T) {
		theMap := map[string]interface{}{}
		condition := &theMap
		assert.True(t, isEmptyCondition(condition))
	})

	t.Run("not-nil ptr to not-empty map", func(t *testing.T) {
		theMap := map[string]interface{}{
			"key": 123,
		}
		condition := &theMap
		assert.False(t, isEmptyCondition(condition))
	})

	t.Run("nil ptr to int", func(t *testing.T) {
		var condition *int
		assert.True(t, isEmptyCondition(condition))
	})

	t.Run("not nil ptr to int", func(t *testing.T) {
		theInt := 123
		condition := &theInt
		assert.False(t, isEmptyCondition(condition))
	})

	t.Run("just int", func(t *testing.T) {
		assert.False(t, isEmptyCondition(123))
	})

	t.Run("nil ptr to slice", func(t *testing.T) {
		var condition *[]interface{}
		assert.True(t, isEmptyCondition(condition))
	})

	t.Run("not-nil ptr to nil-slice", func(t *testing.T) {
		var theSlice []interface{}
		var condition = &theSlice
		assert.True(t, isEmptyCondition(condition))
	})

	t.Run("not-nil ptr to empty slice", func(t *testing.T) {
		theSlice := []interface{}{}
		condition := &theSlice
		assert.True(t, isEmptyCondition(condition))
	})
}
