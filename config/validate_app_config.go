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

	if c.CustomFeeUnit != nil {
		if c.CustomFeeUnit.Bytes <= 0 {
			return spverrors.Newf("invalid custom fee unit - bytes value is equal or less than zero: %d", c.CustomFeeUnit.Bytes)
		}
		if c.CustomFeeUnit.Satoshis < 0 {
			return spverrors.Newf("invalid custom fee unit - satoshis value is less than zero: %d", c.CustomFeeUnit.Satoshis)
		}
	}

	return nil
}
