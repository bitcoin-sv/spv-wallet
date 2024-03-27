package datastore

import (
	"context"
	"database/sql"
	"testing"
	"time"

	customtypes "github.com/bitcoin-sv/spv-wallet/engine/datastore/customtypes"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// Test_whereObject test the SQL where selector
func Test_whereSlice(t *testing.T) {
	t.Parallel()

	t.Run("MySQL", func(t *testing.T) {
		query := whereSlice(MySQL, fieldInIDs, "id_1")
		expected := `JSON_CONTAINS(` + fieldInIDs + `, CAST('["id_1"]' AS JSON))`
		assert.Equal(t, expected, query)
	})

	t.Run("Postgres", func(t *testing.T) {
		query := whereSlice(PostgreSQL, fieldInIDs, "id_1")
		expected := fieldInIDs + `::jsonb @> '["id_1"]'`
		assert.Equal(t, expected, query)
	})

	t.Run("SQLite", func(t *testing.T) {
		query := whereSlice(SQLite, fieldInIDs, "id_1")
		expected := `EXISTS (SELECT 1 FROM json_each(` + fieldInIDs + `) WHERE value = "id_1")`
		assert.Equal(t, expected, query)
	})
}

// Test_processConditions test the SQL where selectors
func Test_processConditions(t *testing.T) {
	t.Parallel()

	dateField := dateCreatedAt
	uniqueField := "unique_field_name"

	conditions := map[string]interface{}{
		dateField: map[string]interface{}{
			conditionGreaterThan: customtypes.NullTime{NullTime: sql.NullTime{
				Valid: true,
				Time:  time.Date(2022, 4, 4, 15, 12, 37, 651387237, time.UTC),
			}},
		},
		uniqueField: map[string]interface{}{
			conditionExists: true,
		},
	}

	t.Run("MySQL", func(t *testing.T) {
		client, deferFunc := testClient(context.Background(), t)
		defer deferFunc()
		tx := &mockSQLCtx{
			WhereClauses: make([]interface{}, 0),
			Vars:         make(map[string]interface{}),
		}
		varNum := 0
		_ = processConditions(client, tx, conditions, MySQL, &varNum, nil)
		// assert.Equal(t, "created_at > @var0", tx.WhereClauses[0])
		assert.Contains(t, tx.WhereClauses, dateField+" > @var0")
		// assert.Equal(t, "unique_field_name IS NOT NULL", tx.WhereClauses[1])
		assert.Contains(t, tx.WhereClauses, uniqueField+" IS NOT NULL")
		assert.Equal(t, "2022-04-04 15:12:37", tx.Vars["var0"])
	})

	t.Run("Postgres", func(t *testing.T) {
		client, deferFunc := testClient(context.Background(), t)
		defer deferFunc()
		tx := &mockSQLCtx{
			WhereClauses: make([]interface{}, 0),
			Vars:         make(map[string]interface{}),
		}
		varNum := 0
		_ = processConditions(client, tx, conditions, PostgreSQL, &varNum, nil)
		// assert.Equal(t, "created_at > @var0", tx.WhereClauses[0])
		assert.Contains(t, tx.WhereClauses, dateField+" > @var0")
		// assert.Equal(t, "unique_field_name IS NOT NULL", tx.WhereClauses[1])
		assert.Contains(t, tx.WhereClauses, uniqueField+" IS NOT NULL")
		assert.Equal(t, "2022-04-04T15:12:37Z", tx.Vars["var0"])
	})

	t.Run("SQLite", func(t *testing.T) {
		client, deferFunc := testClient(context.Background(), t)
		defer deferFunc()
		tx := &mockSQLCtx{
			WhereClauses: make([]interface{}, 0),
			Vars:         make(map[string]interface{}),
		}
		varNum := 0
		_ = processConditions(client, tx, conditions, SQLite, &varNum, nil)
		// assert.Equal(t, "created_at > @var0", tx.WhereClauses[0])
		assert.Contains(t, tx.WhereClauses, dateField+" > @var0")
		// assert.Equal(t, "unique_field_name IS NOT NULL", tx.WhereClauses[1])
		assert.Contains(t, tx.WhereClauses, uniqueField+" IS NOT NULL")
		assert.Equal(t, "2022-04-04T15:12:37.651Z", tx.Vars["var0"])
	})
}

// Test_whereObject test the SQL where selector
func Test_whereObject(t *testing.T) {
	t.Parallel()

	t.Run("MySQL", func(t *testing.T) {
		metadata := map[string]interface{}{
			"test_key": "test-value",
		}
		query := whereObject(MySQL, metadataField, metadata)
		expected := "JSON_EXTRACT(" + metadataField + ", '$.test_key') = \"test-value\""
		assert.Equal(t, expected, query)

		metadata = map[string]interface{}{
			"test_key": "test-'value'",
		}
		query = whereObject(MySQL, metadataField, metadata)
		expected = "JSON_EXTRACT(" + metadataField + ", '$.test_key') = \"test-\\'value\\'\""
		assert.Equal(t, expected, query)

		metadata = map[string]interface{}{
			"test_key1": "test-value",
			"test_key2": "test-value2",
		}
		query = whereObject(MySQL, metadataField, metadata)

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
		query = whereObject(MySQL, "object_metadata", objectMetadata)

		assert.Contains(t, []string{
			"(JSON_EXTRACT(object_metadata, '$.testId.test_key1') = \"test-value\" AND JSON_EXTRACT(object_metadata, '$.testId.test_key2') = \"test-value2\")",
			"(JSON_EXTRACT(object_metadata, '$.testId.test_key2') = \"test-value2\" AND JSON_EXTRACT(object_metadata, '$.testId.test_key1') = \"test-value\")",
		}, query)

		// NOTE: the order of the items can change, hence the query order can change
		// assert.Equal(t, expected, query)
	})

	t.Run("Postgres", func(t *testing.T) {
		metadata := map[string]interface{}{
			"test_key": "test-value",
		}
		query := whereObject(PostgreSQL, metadataField, metadata)
		expected := metadataField + "::jsonb @> '{\"test_key\":\"test-value\"}'::jsonb"
		assert.Equal(t, expected, query)

		metadata = map[string]interface{}{
			"test_key": "test-'value'",
		}
		query = whereObject(PostgreSQL, metadataField, metadata)
		expected = metadataField + "::jsonb @> '{\"test_key\":\"test-\\'value\\'\"}'::jsonb"
		assert.Equal(t, expected, query)

		metadata = map[string]interface{}{
			"test_key1": "test-value",
			"test_key2": "test-value2",
		}
		query = whereObject(PostgreSQL, metadataField, metadata)

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
		query = whereObject(PostgreSQL, "object_metadata", objectMetadata)
		assert.Contains(t, []string{
			"object_metadata::jsonb @> '{\"testId\":{\"test_key1\":\"test-value\",\"test_key2\":\"test-value2\"}}'::jsonb",
			"object_metadata::jsonb @> '{\"testId\":{\"test_key2\":\"test-value2\",\"test_key1\":\"test-value\"}}'::jsonb",
		}, query)

		// NOTE: the order of the items can change, hence the query order can change
		// assert.Equal(t, expected, query)
	})

	t.Run("SQLite", func(t *testing.T) {
		metadata := map[string]interface{}{
			"test_key": "test-value",
		}
		query := whereObject(SQLite, metadataField, metadata)
		expected := "JSON_EXTRACT(" + metadataField + ", '$.test_key') = \"test-value\""
		assert.Equal(t, expected, query)

		metadata = map[string]interface{}{
			"test_key": "test-'value'",
		}
		query = whereObject(SQLite, metadataField, metadata)
		expected = "JSON_EXTRACT(" + metadataField + ", '$.test_key') = \"test-\\'value\\'\""
		assert.Equal(t, expected, query)

		metadata = map[string]interface{}{
			"test_key1": "test-value",
			"test_key2": "test-value2",
		}
		query = whereObject(SQLite, metadataField, metadata)
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
		query = whereObject(SQLite, "object_metadata", objectMetadata)
		assert.Contains(t, []string{
			"(JSON_EXTRACT(object_metadata, '$.testId.test_key1') = \"test-value\" AND JSON_EXTRACT(object_metadata, '$.testId.test_key2') = \"test-value2\")",
			"(JSON_EXTRACT(object_metadata, '$.testId.test_key2') = \"test-value2\" AND JSON_EXTRACT(object_metadata, '$.testId.test_key1') = \"test-value\")",
		}, query)
		// NOTE: the order of the items can change, hence the query order can change
		// assert.Equal(t, expected, query)
	})
}

// mockSQLCtx is used to mock the SQL
type mockSQLCtx struct {
	WhereClauses []interface{}
	Vars         map[string]interface{}
}

func (f *mockSQLCtx) Where(query interface{}, args ...interface{}) {
	f.WhereClauses = append(f.WhereClauses, query)
	if len(args) > 0 {
		for _, variables := range args {
			for key, value := range variables.(map[string]interface{}) {
				f.Vars[key] = value
			}
		}
	}
}

func (f *mockSQLCtx) getGormTx() *gorm.DB {
	return nil
}

// TestCustomWhere will test the method CustomWhere()
func TestCustomWhere(t *testing.T) {
	t.Parallel()

	t.Run("SQLite empty select", func(t *testing.T) {
		client, deferFunc := testClient(context.Background(), t)
		defer deferFunc()
		tx := mockSQLCtx{
			WhereClauses: make([]interface{}, 0),
			Vars:         make(map[string]interface{}),
		}
		conditions := map[string]interface{}{}
		_ = client.CustomWhere(&tx, conditions, SQLite)
		assert.Equal(t, []interface{}{}, tx.WhereClauses)
	})

	t.Run("SQLite simple select", func(t *testing.T) {
		client, deferFunc := testClient(context.Background(), t)
		defer deferFunc()
		tx := mockSQLCtx{
			WhereClauses: make([]interface{}, 0),
			Vars:         make(map[string]interface{}),
		}
		conditions := map[string]interface{}{
			sqlIDFieldProper: "testID",
		}
		_ = client.CustomWhere(&tx, conditions, SQLite)
		assert.Len(t, tx.WhereClauses, 1)
		assert.Equal(t, sqlIDFieldProper+" = @var0", tx.WhereClauses[0])
		assert.Equal(t, "testID", tx.Vars["var0"])
	})

	t.Run("SQLite "+conditionOr, func(t *testing.T) {
		arrayField1 := fieldInIDs
		arrayField2 := fieldOutIDs

		client, deferFunc := testClient(context.Background(), t, WithCustomFields([]string{arrayField1, arrayField2}, nil))
		defer deferFunc()
		tx := mockSQLCtx{
			WhereClauses: make([]interface{}, 0),
			Vars:         make(map[string]interface{}),
		}
		conditions := map[string]interface{}{
			conditionOr: []map[string]interface{}{{
				arrayField1: "value_id",
			}, {
				arrayField2: "value_id",
			}},
		}
		_ = client.CustomWhere(&tx, conditions, SQLite)
		assert.Len(t, tx.WhereClauses, 1)
		assert.Equal(t, " ( (EXISTS (SELECT 1 FROM json_each("+arrayField1+") WHERE value = \"value_id\")) OR (EXISTS (SELECT 1 FROM json_each("+arrayField2+") WHERE value = \"value_id\")) ) ", tx.WhereClauses[0])
	})

	t.Run("MySQL "+conditionOr, func(t *testing.T) {
		arrayField1 := fieldInIDs
		arrayField2 := fieldOutIDs

		client, deferFunc := testClient(context.Background(), t, WithCustomFields([]string{arrayField1, arrayField2}, nil))
		defer deferFunc()
		tx := mockSQLCtx{
			WhereClauses: make([]interface{}, 0),
			Vars:         make(map[string]interface{}),
		}
		conditions := map[string]interface{}{
			conditionOr: []map[string]interface{}{{
				arrayField1: "value_id",
			}, {
				arrayField2: "value_id",
			}},
		}
		_ = client.CustomWhere(&tx, conditions, MySQL)
		assert.Len(t, tx.WhereClauses, 1)
		assert.Equal(t, " ( (JSON_CONTAINS("+arrayField1+", CAST('[\"value_id\"]' AS JSON))) OR (JSON_CONTAINS("+arrayField2+", CAST('[\"value_id\"]' AS JSON))) ) ", tx.WhereClauses[0])
	})

	t.Run("PostgreSQL "+conditionOr, func(t *testing.T) {
		arrayField1 := fieldInIDs
		arrayField2 := fieldOutIDs

		client, deferFunc := testClient(context.Background(), t, WithCustomFields([]string{arrayField1, arrayField2}, nil))
		defer deferFunc()
		tx := mockSQLCtx{
			WhereClauses: make([]interface{}, 0),
			Vars:         make(map[string]interface{}),
		}
		conditions := map[string]interface{}{
			conditionOr: []map[string]interface{}{{
				arrayField1: "value_id",
			}, {
				arrayField2: "value_id",
			}},
		}
		_ = client.CustomWhere(&tx, conditions, PostgreSQL)
		assert.Len(t, tx.WhereClauses, 1)
		assert.Equal(t, " ( ("+arrayField1+"::jsonb @> '[\"value_id\"]') OR ("+arrayField2+"::jsonb @> '[\"value_id\"]') ) ", tx.WhereClauses[0])
	})

	t.Run("SQLite "+metadataField, func(t *testing.T) {
		client, deferFunc := testClient(context.Background(), t)
		defer deferFunc()
		tx := mockSQLCtx{
			WhereClauses: make([]interface{}, 0),
			Vars:         make(map[string]interface{}),
		}
		conditions := map[string]interface{}{
			metadataField: map[string]interface{}{
				"field_name": "field_value",
			},
		}
		_ = client.CustomWhere(&tx, conditions, SQLite)
		assert.Len(t, tx.WhereClauses, 1)
		assert.Equal(t, "JSON_EXTRACT("+metadataField+", '$.field_name') = \"field_value\"", tx.WhereClauses[0])
	})

	t.Run("MySQL "+metadataField, func(t *testing.T) {
		client, deferFunc := testClient(context.Background(), t)
		defer deferFunc()
		tx := mockSQLCtx{
			WhereClauses: make([]interface{}, 0),
			Vars:         make(map[string]interface{}),
		}
		conditions := map[string]interface{}{
			metadataField: map[string]interface{}{
				"field_name": "field_value",
			},
		}
		_ = client.CustomWhere(&tx, conditions, MySQL)
		assert.Len(t, tx.WhereClauses, 1)
		assert.Equal(t, "JSON_EXTRACT("+metadataField+", '$.field_name') = \"field_value\"", tx.WhereClauses[0])
	})

	t.Run("PostgreSQL "+metadataField, func(t *testing.T) {
		client, deferFunc := testClient(context.Background(), t)
		defer deferFunc()
		tx := mockSQLCtx{
			WhereClauses: make([]interface{}, 0),
			Vars:         make(map[string]interface{}),
		}
		conditions := map[string]interface{}{
			metadataField: map[string]interface{}{
				"field_name": "field_value",
			},
		}
		_ = client.CustomWhere(&tx, conditions, PostgreSQL)
		assert.Len(t, tx.WhereClauses, 1)
		assert.Equal(t, metadataField+"::jsonb @> '{\"field_name\":\"field_value\"}'::jsonb", tx.WhereClauses[0])
	})

	t.Run("SQLite "+conditionAnd, func(t *testing.T) {
		arrayField1 := fieldInIDs
		arrayField2 := fieldOutIDs

		client, deferFunc := testClient(context.Background(), t, WithCustomFields([]string{arrayField1, arrayField2}, nil))
		defer deferFunc()
		tx := mockSQLCtx{
			WhereClauses: make([]interface{}, 0),
			Vars:         make(map[string]interface{}),
		}
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
		_ = client.CustomWhere(&tx, conditions, SQLite)
		assert.Len(t, tx.WhereClauses, 1)
		assert.Equal(t, " ( reference_id = @var0 AND number = @var1 AND  ( (EXISTS (SELECT 1 FROM json_each("+arrayField1+") WHERE value = \"value_id\")) OR (EXISTS (SELECT 1 FROM json_each("+arrayField2+") WHERE value = \"value_id\")) )  ) ", tx.WhereClauses[0])
		assert.Equal(t, "reference", tx.Vars["var0"])
		assert.Equal(t, 12, tx.Vars["var1"])
	})

	t.Run("MySQL "+conditionAnd, func(t *testing.T) {
		arrayField1 := fieldInIDs
		arrayField2 := fieldOutIDs

		client, deferFunc := testClient(context.Background(), t, WithCustomFields([]string{arrayField1, arrayField2}, nil))
		defer deferFunc()
		tx := mockSQLCtx{
			WhereClauses: make([]interface{}, 0),
			Vars:         make(map[string]interface{}),
		}
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
		_ = client.CustomWhere(&tx, conditions, MySQL)
		assert.Len(t, tx.WhereClauses, 1)
		assert.Equal(t, " ( reference_id = @var0 AND number = @var1 AND  ( (JSON_CONTAINS("+arrayField1+", CAST('[\"value_id\"]' AS JSON))) OR (JSON_CONTAINS("+arrayField2+", CAST('[\"value_id\"]' AS JSON))) )  ) ", tx.WhereClauses[0])
		assert.Equal(t, "reference", tx.Vars["var0"])
		assert.Equal(t, 12, tx.Vars["var1"])
	})

	t.Run("PostgreSQL "+conditionAnd, func(t *testing.T) {
		arrayField1 := fieldInIDs
		arrayField2 := fieldOutIDs

		client, deferFunc := testClient(context.Background(), t, WithCustomFields([]string{arrayField1, arrayField2}, nil))
		defer deferFunc()
		tx := mockSQLCtx{
			WhereClauses: make([]interface{}, 0),
			Vars:         make(map[string]interface{}),
		}
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
		_ = client.CustomWhere(&tx, conditions, PostgreSQL)
		assert.Len(t, tx.WhereClauses, 1)
		assert.Equal(t, " ( reference_id = @var0 AND number = @var1 AND  ( ("+arrayField1+"::jsonb @> '[\"value_id\"]') OR ("+arrayField2+"::jsonb @> '[\"value_id\"]') )  ) ", tx.WhereClauses[0])
		assert.Equal(t, "reference", tx.Vars["var0"])
		assert.Equal(t, 12, tx.Vars["var1"])
	})

	t.Run("Where "+conditionGreaterThan, func(t *testing.T) {
		client, deferFunc := testClient(context.Background(), t)
		defer deferFunc()
		tx := mockSQLCtx{
			WhereClauses: make([]interface{}, 0),
			Vars:         make(map[string]interface{}),
		}
		conditions := map[string]interface{}{
			"amount": map[string]interface{}{
				conditionGreaterThan: 502,
			},
		}
		_ = client.CustomWhere(&tx, conditions, PostgreSQL) // all the same
		assert.Len(t, tx.WhereClauses, 1)
		assert.Equal(t, "amount > @var0", tx.WhereClauses[0])
		assert.Equal(t, 502, tx.Vars["var0"])
	})

	t.Run("Where "+conditionGreaterThan+" "+conditionLessThan, func(t *testing.T) {
		client, deferFunc := testClient(context.Background(), t)
		defer deferFunc()
		tx := mockSQLCtx{
			WhereClauses: make([]interface{}, 0),
			Vars:         make(map[string]interface{}),
		}
		conditions := map[string]interface{}{
			conditionAnd: []map[string]interface{}{{
				"amount": map[string]interface{}{
					conditionLessThan: 503,
				},
			}, {
				"amount": map[string]interface{}{
					conditionGreaterThan: 203,
				},
			}},
		}
		_ = client.CustomWhere(&tx, conditions, PostgreSQL) // all the same
		assert.Len(t, tx.WhereClauses, 1)
		assert.Equal(t, " ( amount < @var0 AND amount > @var1 ) ", tx.WhereClauses[0])
		assert.Equal(t, 503, tx.Vars["var0"])
		assert.Equal(t, 203, tx.Vars["var1"])
	})

	t.Run("Where "+conditionGreaterThanOrEqual+" "+conditionLessThanOrEqual, func(t *testing.T) {
		client, deferFunc := testClient(context.Background(), t)
		defer deferFunc()
		tx := mockSQLCtx{
			WhereClauses: make([]interface{}, 0),
			Vars:         make(map[string]interface{}),
		}
		conditions := map[string]interface{}{
			conditionOr: []map[string]interface{}{{
				"amount": map[string]interface{}{
					conditionLessThanOrEqual: 203,
				},
			}, {
				"amount": map[string]interface{}{
					conditionGreaterThanOrEqual: 1203,
				},
			}},
		}
		_ = client.CustomWhere(&tx, conditions, PostgreSQL) // all the same
		assert.Len(t, tx.WhereClauses, 1)
		assert.Equal(t, " ( (amount <= @var0) OR (amount >= @var1) ) ", tx.WhereClauses[0])
		assert.Equal(t, 203, tx.Vars["var0"])
		assert.Equal(t, 1203, tx.Vars["var1"])
	})

	t.Run("Where "+conditionOr+" "+conditionAnd+" "+conditionOr+" "+conditionGreaterThanOrEqual+" "+conditionLessThanOrEqual, func(t *testing.T) {
		client, deferFunc := testClient(context.Background(), t)
		defer deferFunc()
		tx := mockSQLCtx{
			WhereClauses: make([]interface{}, 0),
			Vars:         make(map[string]interface{}),
		}
		conditions := map[string]interface{}{
			conditionOr: []map[string]interface{}{{
				conditionAnd: []map[string]interface{}{{
					"amount": map[string]interface{}{
						conditionLessThanOrEqual: 203,
					},
				}, {
					conditionOr: []map[string]interface{}{{
						"amount": map[string]interface{}{
							conditionGreaterThanOrEqual: 1203,
						},
					}, {
						"value": map[string]interface{}{
							conditionGreaterThanOrEqual: 2203,
						},
					}},
				}},
			}, {
				conditionAnd: []map[string]interface{}{{
					"amount": map[string]interface{}{
						conditionGreaterThanOrEqual: 3203,
					},
				}, {
					"value": map[string]interface{}{
						conditionGreaterThanOrEqual: 4203,
					},
				}},
			}},
		}
		_ = client.CustomWhere(&tx, conditions, PostgreSQL) // all the same
		assert.Len(t, tx.WhereClauses, 1)
		assert.Equal(t, " ( ( ( amount <= @var0 AND  ( (amount >= @var1) OR (value >= @var2) )  ) ) OR ( ( amount >= @var3 AND value >= @var4 ) ) ) ", tx.WhereClauses[0])
		assert.Equal(t, 203, tx.Vars["var0"])
		assert.Equal(t, 1203, tx.Vars["var1"])
		assert.Equal(t, 2203, tx.Vars["var2"])
		assert.Equal(t, 3203, tx.Vars["var3"])
		assert.Equal(t, 4203, tx.Vars["var4"])
	})

	t.Run("Where "+conditionAnd+" "+conditionOr+" "+conditionOr+" "+conditionGreaterThanOrEqual+" "+conditionLessThanOrEqual, func(t *testing.T) {
		client, deferFunc := testClient(context.Background(), t)
		defer deferFunc()
		tx := mockSQLCtx{
			WhereClauses: make([]interface{}, 0),
			Vars:         make(map[string]interface{}),
		}
		conditions := map[string]interface{}{
			conditionAnd: []map[string]interface{}{{
				conditionAnd: []map[string]interface{}{{
					"amount": map[string]interface{}{
						conditionLessThanOrEqual:    203,
						conditionGreaterThanOrEqual: 103,
					},
				}, {
					conditionOr: []map[string]interface{}{{
						"amount": map[string]interface{}{
							conditionGreaterThanOrEqual: 1203,
						},
					}, {
						"value": map[string]interface{}{
							conditionGreaterThanOrEqual: 2203,
						},
					}},
				}},
			}, {
				conditionOr: []map[string]interface{}{{
					"amount": map[string]interface{}{
						conditionGreaterThanOrEqual: 3203,
					},
				}, {
					"value": map[string]interface{}{
						conditionGreaterThanOrEqual: 4203,
					},
				}},
			}},
		}
		_ = client.CustomWhere(&tx, conditions, PostgreSQL) // all the same
		assert.Len(t, tx.WhereClauses, 1)
		assert.Contains(t, []string{
			" (  ( amount <= @var0 AND amount >= @var1 AND  ( (amount >= @var2) OR (value >= @var3) )  )  AND  ( (amount >= @var4) OR (value >= @var5) )  ) ",
			" (  ( amount >= @var0 AND amount <= @var1 AND  ( (amount >= @var2) OR (value >= @var3) )  )  AND  ( (amount >= @var4) OR (value >= @var5) )  ) ",
		}, tx.WhereClauses[0])

		// assert.Equal(t, " (  ( amount <= @var0 AND amount >= @var1 AND  ( (amount >= @var2) OR (value >= @var3) )  )  AND  ( (amount >= @var4) OR (value >= @var5) )  ) ", tx.WhereClauses[0])

		assert.Contains(t, []int{203, 103}, tx.Vars["var0"])
		assert.Contains(t, []int{203, 103}, tx.Vars["var1"])
		// assert.Equal(t, 203, tx.Vars["var0"])
		// assert.Equal(t, 103, tx.Vars["var1"])
		assert.Equal(t, 1203, tx.Vars["var2"])
		assert.Equal(t, 2203, tx.Vars["var3"])
		assert.Equal(t, 3203, tx.Vars["var4"])
		assert.Equal(t, 4203, tx.Vars["var5"])
	})
}

// Test_escapeDBString will test the method escapeDBString()
func Test_escapeDBString(t *testing.T) {
	t.Parallel()

	str := escapeDBString(`SELECT * FROM 'table' WHERE 'field'=1;`)
	assert.Equal(t, `SELECT * FROM \'table\' WHERE \'field\'=1;`, str)
}
