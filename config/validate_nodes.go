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

	if n.Protocol == NodesProtocolMapi {
		// check if at least one mapi url is configured
		anyMapiNode := slices.IndexFunc(n.Apis, func(el *MinerAPI) bool {
			return isMapiNode(el)
		})
		if anyMapiNode == -1 {
			return errors.New("no mapi urls configured")
		}

		wrongMapiNode := slices.IndexFunc(n.Apis, func(el *MinerAPI) bool {
			return isMapiNode(el) && el.MinerID == ""
		})
		if wrongMapiNode != -1 {
			return errors.New("mapi url configured without miner id")
		}

		// check if MinerIDs for mAPI nodes are unique
		ids := make(map[string]bool)
		for _, el := range n.Apis {
			if isMapiNode(el) {
				if _, ok := ids[el.MinerID]; ok {
					return errors.New("miner ids are not unique")
				}
				ids[el.MinerID] = true
			}
		}
	}

	// check if at least one arc url is configured
	if n.Protocol == NodesProtocolArc {
		found := slices.IndexFunc(n.Apis, func(el *MinerAPI) bool {
			return isArcNode(el)
		})
		if found == -1 {
			return errors.New("no arc urls configured")
		}
	}

	return nil
}

func isMapiNode(node *MinerAPI) bool {
	return node.MapiURL != ""
}

func isArcNode(node *MinerAPI) bool {
	return node.ArcURL != ""
}
