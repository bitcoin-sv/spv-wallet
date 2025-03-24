package datastore

import (
	"database/sql"
	"testing"
	"time"

	customtypes "github.com/bitcoin-sv/spv-wallet/engine/datastore/customtypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testTableName = "test_model"
)

type testModel struct {
	Field string `json:"field"`
}

// GetModelName will return a model name
func (t *testModel) GetModelName() string {
	return testModelName
}

// GetModelTableName will return a table name
func (t *testModel) GetModelTableName() string {
	return testTableName
}

type badModel struct {
	Field string `json:"field"`
}

func TestGetModelStringAttribute(t *testing.T) {
	t.Parallel()

	type TestModel struct {
		StringField string `json:"string_field"`
		ID          string `json:"id"`
	}

	t.Run("valid string attribute", func(t *testing.T) {
		m := &TestModel{
			StringField: "test",
			ID:          "12345678",
		}
		field1 := GetModelStringAttribute(m, "StringField")
		id := GetModelStringAttribute(m, sqlIDFieldProper)
		assert.Equal(t, "test", *field1)
		assert.Equal(t, "12345678", *id)
	})

	t.Run("nil input", func(t *testing.T) {
		id := GetModelStringAttribute(nil, sqlIDFieldProper)
		assert.Nil(t, id)
	})

	t.Run("invalid type", func(t *testing.T) {
		id := GetModelStringAttribute("invalid-type", sqlIDFieldProper)
		assert.Nil(t, id)
	})
}

func TestGetModelBoolAttribute(t *testing.T) {
	t.Parallel()

	type TestModel struct {
		BoolField bool   `json:"bool_field"`
		ID        string `json:"id"`
	}

	t.Run("valid bool attribute", func(t *testing.T) {
		m := &TestModel{
			BoolField: true,
		}
		field1 := GetModelBoolAttribute(m, "BoolField")
		assert.True(t, *field1)
	})

	t.Run("nil input", func(t *testing.T) {
		val := GetModelBoolAttribute(nil, "BoolField")
		assert.Nil(t, val)
	})

	t.Run("invalid type", func(t *testing.T) {
		assert.Panics(t, func() {
			val := GetModelBoolAttribute("invalid-type", "BoolField")
			assert.Nil(t, val)
		})
	})
}

func TestGetModelUnset(t *testing.T) {
	t.Parallel()

	type TestModel struct {
		NullableTime   customtypes.NullTime   `json:"nullable_time"`
		NullableString customtypes.NullString `json:"nullable_string"`
		Internal       string                 `json:"-"`
	}

	t.Run("basic test", func(t *testing.T) {
		ty := make(map[string]bool)
		m := &TestModel{
			NullableTime: customtypes.NullTime{NullTime: sql.NullTime{
				Time:  time.Time{},
				Valid: false,
			}},
			NullableString: customtypes.NullString{NullString: sql.NullString{
				String: "",
				Valid:  false,
			}},
			Internal: "test",
		}
		un := GetModelUnset(m)
		assert.IsType(t, ty, un)
		assert.True(t, un["nullable_time"])
		assert.True(t, un["nullable_string"])
		assert.False(t, un["internal"])
	})
}

func TestIsModelSlice(t *testing.T) {
	t.Parallel()

	t.Run("valid slices", func(t *testing.T) {
		s := []string{"test"}
		assert.True(t, IsModelSlice(s))

		i := []int{1}
		assert.True(t, IsModelSlice(i))

		in := []interface{}{"test"}
		assert.True(t, IsModelSlice(in))

		ptr := []string{"test"}
		assert.True(t, IsModelSlice(&ptr))
	})

	t.Run("not a slice", func(t *testing.T) {
		s := "string"
		assert.False(t, IsModelSlice(s))

		i := 1
		assert.False(t, IsModelSlice(i))
	})
}

func TestGetModelName(t *testing.T) {
	t.Parallel()

	t.Run("model is nil", func(t *testing.T) {
		name := GetModelName(nil)
		require.Nil(t, name)
	})

	t.Run("model is set - pointer", func(t *testing.T) {
		tm := &testModel{Field: testModelName}
		name := GetModelName(tm)
		assert.Equal(t, testModelName, *name)
	})

	t.Run("model is set - value", func(t *testing.T) {
		tm := testModel{Field: testModelName}
		name := GetModelName(tm)
		assert.Equal(t, testModelName, *name)
	})

	t.Run("models are set - value", func(t *testing.T) {
		tm := &[]testModel{{Field: testModelName}}
		name := GetModelName(tm)
		assert.Equal(t, testModelName, *name)
	})

	t.Run("model does not have method - pointer", func(t *testing.T) {
		tm := &badModel{}
		name := GetModelName(tm)
		assert.Nil(t, name)
	})

	t.Run("model does not have method - value", func(t *testing.T) {
		tm := badModel{}
		name := GetModelName(tm)
		assert.Nil(t, name)
	})
}

func TestGetModelTableName(t *testing.T) {
	t.Parallel()

	t.Run("model is nil", func(t *testing.T) {
		name := GetModelTableName(nil)
		require.Nil(t, name)
	})

	t.Run("model is set - pointer", func(t *testing.T) {
		tm := &testModel{Field: testTableName}
		name := GetModelTableName(tm)
		assert.Equal(t, testTableName, *name)
	})

	t.Run("model is set - value", func(t *testing.T) {
		tm := testModel{Field: testTableName}
		name := GetModelTableName(tm)
		assert.Equal(t, testTableName, *name)
	})

	t.Run("models are set - value", func(t *testing.T) {
		tm := &[]testModel{{Field: testModelName}}
		name := GetModelTableName(tm)
		assert.Equal(t, testModelName, *name)
	})

	t.Run("model does not have method - pointer", func(t *testing.T) {
		tm := &badModel{}
		name := GetModelTableName(tm)
		assert.Nil(t, name)
	})

	t.Run("model does not have method - value", func(t *testing.T) {
		tm := badModel{}
		name := GetModelTableName(tm)
		assert.Nil(t, name)
	})
}

func TestGetModelType(t *testing.T) {
	t.Parallel()

	type modelExample struct {
		Field string `json:"field"`
	}

	t.Run("panic - nil model", func(t *testing.T) {
		assert.Panics(t, func() {
			modelType := GetModelType(nil)
			require.NotNil(t, modelType)
		})
	})

	t.Run("default type", func(t *testing.T) {
		m := new(modelExample)
		modelType := GetModelType(m)
		assert.NotNil(t, modelType)
	})
}

func TestStringInSlice(t *testing.T) {
	t.Parallel()

	t.Run("nil / empty", func(t *testing.T) {
		assert.False(t, StringInSlice("test", []string{}))
		assert.False(t, StringInSlice("test", nil))
	})

	t.Run("slices", func(t *testing.T) {
		slice := []string{"test", "test1", "test2"}
		assert.True(t, StringInSlice("test", slice))
		assert.True(t, StringInSlice("test1", slice))
		assert.True(t, StringInSlice("test2", slice))
		assert.False(t, StringInSlice("test3", slice))
	})
}
