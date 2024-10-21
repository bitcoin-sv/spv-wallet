package config

import "github.com/bitcoin-sv/spv-wallet/engine/spverrors"

// Validate validates the custom fee unit configuration
func (cf *FeeUnitConfig) Validate() error {
	if cf == nil {
		return nil
	}

	if cf.Bytes <= 0 {
		return spverrors.Newf("invalid custom fee unit - bytes value is equal or less than zero: %d", cf.Bytes)
	}
	if cf.Satoshis < 0 {
		return spverrors.Newf("invalid custom fee unit - satoshis value is less than zero: %d", cf.Satoshis)
	}
	return nil
}
