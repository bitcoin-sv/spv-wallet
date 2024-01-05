package config

import (
	"errors"
	"slices"
)

// Validate checks the configuration for specific rules
func (n *NodesConfig) Validate() error {
	if n == nil {
		return errors.New("nodes are not configured")
	}

	err := n.Protocol.Validate()
	if err != nil {
		return err
	}

	if n.Apis == nil || len(n.Apis) == 0 {
		return errors.New("no miner apis configured")
	}

	// check if at least one mapi url is configured
	if n.Protocol == NodesProtocolMapi {
		found := slices.IndexFunc(n.Apis, func(el *MinerAPI) bool {
			return el.MapiURL != ""
		})
		if found == -1 {
			return errors.New("no mapi urls configured")
		}
	}

	// check if at least one arc url is configured
	if n.Protocol == NodesProtocolArc {
		found := slices.IndexFunc(n.Apis, func(el *MinerAPI) bool {
			return el.ArcURL != ""
		})
		if found == -1 {
			return errors.New("no arc urls configured")
		}
	}

	return nil
}
