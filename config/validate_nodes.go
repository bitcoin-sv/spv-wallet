package config

import (
	"slices"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
)

// Validate checks the configuration for specific rules
func (n *NodesConfig) Validate() error {
	if n == nil {
		return spverrors.Newf("nodes are not configured")
	}

	if len(n.Apis) == 0 {
		return spverrors.Newf("no miner apis configured")
	}

	// check if at least one arc url is configured
	found := slices.IndexFunc(n.Apis, func(el *ArcAPI) bool {
		return el.ArcURL != ""
	})
	if found == -1 {
		return spverrors.Newf("no arc urls configured")
	}

	if !n.UseFeeQuotes && n.FeeUnit == nil {
		return spverrors.Newf("fee unit is not configured, define nodes.fee_unit or set nodes.use_fee_quotes")
	}

	return nil
}
