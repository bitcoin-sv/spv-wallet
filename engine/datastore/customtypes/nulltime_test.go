package customtypes

import (
	"bytes"
	"database/sql/driver"
	"strconv"
	"testing"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNullTime(t *testing.T) {
	t.Parallel()

	t.Run("empty time", func(t *testing.T) {
		nt := new(NullTime)
		assert.False(t, nt.Valid)
		assert.True(t, driver.IsValue(nt.Time))
	})

	t.Run("time", func(t *testing.T) {
		nt := new(NullTime)
		testTime := time.Now().UTC()
		nt.Time = testTime
		nt.Valid = true
		assert.True(t, nt.Valid)
		assert.True(t, driver.IsValue(nt.Time))
		assert.Equal(t, testTime, nt.Time)
	})
}

func TestMarshalNullTime(t *testing.T) {
	t.Parallel()

	t.Run("empty time", func(t *testing.T) {
		nt := new(NullTime)
		marshaller := MarshalNullTime(*nt)
		assert.Equal(t, graphql.Null, marshaller)
	})

	t.Run("time", func(t *testing.T) {
		nt := new(NullTime)
		testTime := time.Now().UTC()
		nt.Time = testTime
		nt.Valid = true
		marshaller := MarshalNullTime(*nt)
		var b bytes.Buffer
		marshaller.MarshalGQL(&b)
		assert.Equal(t, strconv.Quote(testTime.Format(time.RFC3339Nano)), b.String())
	})
}

func TestUnmarshalNullTime(t *testing.T) {
	t.Parallel()

	t.Run("nil time", func(t *testing.T) {
		nt, err := UnmarshalNullTime(nil)
		require.NoError(t, err)
		assert.IsType(t, NullTime{}, nt)
		assert.False(t, nt.Valid)
	})

	t.Run("time", func(t *testing.T) {
		val := time.Time{}
		nt, err := UnmarshalNullTime(val)
		require.Error(t, err)
		assert.IsType(t, NullTime{}, nt)
		assert.False(t, nt.Valid)
	})

	t.Run("time string", func(t *testing.T) {
		testTime := time.Now().UTC()
		testTime = testTime.Add(-1 * time.Duration(testTime.Nanosecond()))
		str := testTime.Format(time.RFC3339)
		nt, err := UnmarshalNullTime(str)
		require.NoError(t, err)
		assert.IsType(t, NullTime{}, nt)
		assert.True(t, nt.Valid)
		assert.True(t, testTime.Equal(nt.Time))
	})
}

func TestNullTime_IsZero(t *testing.T) {
	t.Parallel()

	t.Run("nil time", func(t *testing.T) {
		nt := new(NullTime)
		assert.True(t, nt.IsZero())
	})

	t.Run("time", func(t *testing.T) {
		nt := time.Now().UTC()
		assert.False(t, nt.IsZero())
	})
}

func TestNullTime_MarshalJSON(t *testing.T) {
	t.Parallel()

	t.Run("empty time", func(t *testing.T) {
		nt := new(NullTime)
		marshaled, err := nt.MarshalJSON()
		require.NoError(t, err)
		assert.Equal(t, "null", string(marshaled))
	})

	t.Run("time", func(t *testing.T) {
		nt := new(NullTime)
		testTime := time.Now().UTC()
		testTime = testTime.Add(-1 * time.Duration(testTime.Nanosecond()))
		nt.Time = testTime
		nt.Valid = true
		marshaled, err := nt.MarshalJSON()
		require.NoError(t, err)
		assert.Equal(t, strconv.Quote(testTime.Format(time.RFC3339)), string(marshaled))
	})
}

func TestNullTime_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	t.Run("nil time", func(t *testing.T) {
		nt := new(NullTime)
		err := nt.UnmarshalJSON(nil)
		require.NoError(t, err)
		assert.IsType(t, &NullTime{}, nt)
		assert.False(t, nt.Valid)
	})

	t.Run("empty time", func(t *testing.T) {
		nt := new(NullTime)
		err := nt.UnmarshalJSON([]byte(""))
		require.Error(t, err)
		assert.IsType(t, &NullTime{}, nt)
		assert.False(t, nt.Valid)
	})

	t.Run("empty string time", func(t *testing.T) {
		nt := new(NullTime)
		err := nt.UnmarshalJSON([]byte("\"\""))
		require.NoError(t, err)
		assert.IsType(t, &NullTime{}, nt)
		assert.False(t, nt.Valid)
	})

	t.Run("time string", func(t *testing.T) {
		testTime := time.Now().UTC()
		testTime = testTime.Add(-1 * time.Duration(testTime.Nanosecond()))
		str := testTime.Format(time.RFC3339)
		nt := new(NullTime)
		err := nt.UnmarshalJSON([]byte("\"" + str + "\""))
		require.NoError(t, err)
		assert.IsType(t, &NullTime{}, nt)
		assert.True(t, nt.Valid)
		assert.True(t, testTime.Equal(nt.Time))
	})
}
