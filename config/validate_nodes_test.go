package config

import (
	"testing"

	"github.com/BuxOrg/bux-server/logging"
	"github.com/stretchr/testify/assert"
)

// TestNewRelicConfig_Validate will test the method Validate()
func TestNodesConfig_Validate(t *testing.T) {
	t.Parallel()
	defLogger := logging.GetDefaultLogger()

	t.Run("valid default nodes config", func(t *testing.T) {
		n := getNodesDefaults(defLogger)
		assert.NoError(t, n.Validate())
	})

	t.Run("wrong protocol", func(t *testing.T) {
		n := getNodesDefaults(defLogger)
		n.Protocol = "wrong"
		assert.Error(t, n.Validate())
	})

	t.Run("empty list of apis", func(t *testing.T) {
		n := getNodesDefaults(defLogger)

		n.Apis = nil
		assert.Error(t, n.Validate())

		n.Apis = []*MinerAPI{}
		assert.Error(t, n.Validate())
	})

	t.Run("no mapi url", func(t *testing.T) {
		n := getNodesDefaults(defLogger)

		n.Apis = []*MinerAPI{
			{
				MapiURL: "",
			},
		}
		assert.Error(t, n.Validate())
	})

	t.Run("no mapi url", func(t *testing.T) {
		n := getNodesDefaults(defLogger)

		n.Protocol = NodesProtocolMapi
		n.Apis[0].MapiURL = ""
		assert.Error(t, n.Validate())
	})

	t.Run("no arc url", func(t *testing.T) {
		n := getNodesDefaults(defLogger)

		n.Protocol = NodesProtocolArc
		n.Apis[0].ArcURL = ""
		assert.Error(t, n.Validate())
	})

	t.Run("mapi url without miner id", func(t *testing.T) {
		n := getNodesDefaults(defLogger)

		n.Protocol = NodesProtocolMapi
		n.Apis[0].MapiURL = "http://localhost"
		n.Apis[0].MinerID = ""
		assert.Error(t, n.Validate())
	})

	t.Run("mapi url with the same miner id", func(t *testing.T) {
		n := getNodesDefaults(defLogger)

		n.Protocol = NodesProtocolMapi
		n.Apis[0].MapiURL = "http://localhost"
		n.Apis[0].MinerID = "miner1"
		n.Apis = append(n.Apis, &MinerAPI{
			MapiURL: "http://localhost",
			MinerID: "miner1",
		})

		assert.Error(t, n.Validate())
	})
}
