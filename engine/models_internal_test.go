package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestModelSetRecordTime(t *testing.T) {
	t.Parallel()

	t.Run("empty model", func(t *testing.T) {
		m := new(Model)
		assert.Equal(t, true, m.CreatedAt.IsZero())
		assert.Equal(t, true, m.UpdatedAt.IsZero())
	})

	t.Run("set created at time", func(t *testing.T) {
		m := new(Model)
		m.SetRecordTime(true)
		assert.Equal(t, false, m.CreatedAt.IsZero())
		assert.Equal(t, true, m.UpdatedAt.IsZero())
	})

	t.Run("set updated at time", func(t *testing.T) {
		m := new(Model)
		m.SetRecordTime(false)
		assert.Equal(t, true, m.CreatedAt.IsZero())
		assert.Equal(t, false, m.UpdatedAt.IsZero())
	})

	t.Run("set both times", func(t *testing.T) {
		m := new(Model)
		m.SetRecordTime(false)
		m.SetRecordTime(true)
		assert.Equal(t, false, m.CreatedAt.IsZero())
		assert.Equal(t, false, m.UpdatedAt.IsZero())
	})
}

func TestModelNew(t *testing.T) {
	t.Parallel()

	t.Run("New model", func(t *testing.T) {
		m := new(Model)
		assert.Equal(t, false, m.IsNew())
	})

	t.Run("set New flag", func(t *testing.T) {
		m := new(Model)
		m.New()
		assert.Equal(t, true, m.IsNew())
	})
}

func TestModelGetOptions(t *testing.T) {
	// t.Parallel()

	t.Run("base model", func(t *testing.T) {
		m := new(Model)
		opts := m.GetOptions(false)
		assert.Equal(t, 0, len(opts))
	})

	t.Run("new record model", func(t *testing.T) {
		m := new(Model)
		opts := m.GetOptions(true)
		assert.Equal(t, 1, len(opts))
	})
}

func TestModel_IsNew(t *testing.T) {
	t.Parallel()

	t.Run("base model", func(t *testing.T) {
		m := new(Model)
		assert.Equal(t, false, m.IsNew())
	})

	t.Run("New model", func(t *testing.T) {
		m := new(Model)
		m.New()
		assert.Equal(t, true, m.IsNew())
	})
}

func TestModel_RawXpub(t *testing.T) {
	m := new(Model)
	m.rawXpubKey = "xpub661MyMwAqRbcFqp1qzrF2AryEo4X8W1CNSAiT7t2wgXxkbt8nSrdZFYQeD19aTeiPtpAHDGtNUBxgFAg5d2GMzbAiVEsP8DJPLgTQ2LvZTz"
	assert.Equal(t, "xpub661MyMwAqRbcFqp1qzrF2AryEo4X8W1CNSAiT7t2wgXxkbt8nSrdZFYQeD19aTeiPtpAHDGtNUBxgFAg5d2GMzbAiVEsP8DJPLgTQ2LvZTz", m.RawXpub())
}

func TestModel_Name(t *testing.T) {
	t.Parallel()

	t.Run("base model", func(t *testing.T) {
		m := new(Model)
		assert.Equal(t, "", m.Name())
	})

	t.Run("set model name", func(t *testing.T) {
		m := new(Model)
		m.name = ModelXPub
		assert.Equal(t, "xpub", m.Name())
	})
}
