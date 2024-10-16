package config

import "github.com/bitcoin-sv/spv-wallet/engine/spverrors"

// Validate checks the configuration for specific rules
func (c *AppConfig) Validate() error {
	var err error

	if err = c.Authentication.Validate(); err != nil {
		return err
	}

	if err = c.Cache.Validate(); err != nil {
		return err
	}

	if err = c.Db.Validate(); err != nil {
		return err
	}

	if err = c.Paymail.Validate(); err != nil {
		return err
	}

	if err = c.BHS.Validate(); err != nil {
		return err
	}

	if err = c.Server.Validate(); err != nil {
		return err
	}

	if err = c.ARC.Validate(); err != nil {
		return err
	}

	if c.CustomFeeUnit != nil && (c.CustomFeeUnit.Bytes <= 0 || c.CustomFeeUnit.Satoshis < 0) {
		return spverrors.Newf("invalid custom fee unit - satoshis: %d; bytes: %d", c.CustomFeeUnit.Satoshis, c.CustomFeeUnit.Bytes)
	}

	return nil
}
