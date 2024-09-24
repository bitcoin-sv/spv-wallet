package config

import "github.com/bitcoin-sv/spv-wallet/engine/spverrors"

// Validate checks the configuration for specific rules
func (b *BHSConfig) Validate() error {
	if b == nil {
		return spverrors.Newf("bhs config is required")
	}

	if b.URL == "" {
		return spverrors.Newf("bhs url is required")
	}

	return nil
}
