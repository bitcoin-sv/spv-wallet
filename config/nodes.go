package config

import (
	"errors"

	broadcastclient "github.com/bitcoin-sv/go-broadcast-client/broadcast/broadcast-client"
)

// NodesProtocol is the protocol/api_type used to communicate with the miners
type NodesProtocol string

const (
	// NodesProtocolArc represents the arc protocol provided by go-broadcast-client
	NodesProtocolArc NodesProtocol = "arc"
)

// Validate whether the protocol is known
func (n NodesProtocol) Validate() error {
	switch n {
	case NodesProtocolArc:
		return nil
	default:
		return errors.New("invalid nodes protocol")
	}
}

func (nodes *NodesConfig) toBroadcastClientArc() []*broadcastclient.ArcClientConfig {
	minerApis := []*broadcastclient.ArcClientConfig{}
	if nodes.Apis != nil {
		for _, cfg := range nodes.Apis {
			if cfg.ArcURL == "" {
				continue
			}

			minerApis = append(minerApis, &broadcastclient.ArcClientConfig{
				Token:        cfg.Token,
				APIUrl:       cfg.ArcURL,
				DeploymentID: nodes.DeploymentID,
			})
		}
	}
	return minerApis
}
