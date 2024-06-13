package config

import (
	broadcastclient "github.com/bitcoin-sv/go-broadcast-client/broadcast/broadcast-client"
)

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
