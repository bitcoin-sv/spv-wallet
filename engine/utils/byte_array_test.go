package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToByteArray(t *testing.T) {
	t.Run("should convert string to byte array", func(t *testing.T) {
		expected := []byte("test")
		actual, err := ToByteArray("test")
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, expected, actual)
	})

	t.Run("should leave byte array", func(t *testing.T) {
		expected := []byte("test")
		actual, err := ToByteArray([]byte("test"))
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, expected, actual)
	})

	t.Run("should return error for unsupported type", func(t *testing.T) {
		_, err := ToByteArray(1)
		assert.EqualError(t, err, "unsupported type: int")
	})
}

func TestStrOrBytesToString(t *testing.T) {
	t.Run("should convert byte array to string", func(t *testing.T) {
		expected := "test"
		actual, err := StrOrBytesToString([]byte("test"))
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, expected, actual)
	})

	t.Run("should leave string", func(t *testing.T) {
		expected := "test"
		actual, err := StrOrBytesToString("test")
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, expected, actual)
	})

	t.Run("should return error for unsupported type", func(t *testing.T) {
		_, err := StrOrBytesToString(1)
		assert.EqualError(t, err, "unsupported type: int")
	})
}
