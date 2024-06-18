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

	if n.Apis == nil || len(n.Apis) == 0 {
		return errors.New("no miner apis configured")
	}

	// check if at least one arc url is configured
	found := slices.IndexFunc(n.Apis, func(el *ArcAPI) bool {
		return isArcNode(el)
	})
	if found == -1 {
		return errors.New("no arc urls configured")
	}

	if !n.UseFeeQuotes && n.FeeUnit == nil {
		return errors.New("fee unit is not configured, define nodes.fee_unit or set nodes.use_fee_quotes")
	}

	return nil
}

func isArcNode(node *ArcAPI) bool {
	return node.ArcURL != ""
}
