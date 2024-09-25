package config

import (
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
)

// Validate checks the configuration for specific rules
func (n *ARCConfig) Validate() error {
	if n == nil {
		return spverrors.Newf("nodes are not configured")
	}

	if n.URL == "" {
		return spverrors.Newf("node url is not configured")
	}

	if !n.UseFeeQuotes && n.FeeUnit == nil {
		return spverrors.Newf("fee unit is not configured, define nodes.fee_unit or set nodes.use_fee_quotes")
	}

	return nil
}
