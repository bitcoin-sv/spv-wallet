package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewRelicConfig_Validate will test the method Validate()
func TestNodesConfig_Validate(t *testing.T) {
	t.Parallel()

	t.Run("valid default nodes config", func(t *testing.T) {
		n := getNodesDefaults()
		assert.NoError(t, n.Validate())
	})

	t.Run("wrong protocol", func(t *testing.T) {
		n := getNodesDefaults()
		n.Protocol = "wrong"
		assert.Error(t, n.Validate())
	})

	t.Run("empty list of apis", func(t *testing.T) {
		n := getNodesDefaults()

		n.Apis = nil
		assert.Error(t, n.Validate())

		n.Apis = []*MinerAPI{}
		assert.Error(t, n.Validate())
	})

	t.Run("no arc url", func(t *testing.T) {
		n := getNodesDefaults()

		n.Protocol = NodesProtocolArc
		n.Apis[0].ArcURL = ""
		assert.Error(t, n.Validate())
	})
}
