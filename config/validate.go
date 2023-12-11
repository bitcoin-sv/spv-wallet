package config

import (
	"github.com/newrelic/go-agent/v3/newrelic"
)

// Validate checks the configuration for specific rules
func (a *AppConfig) Validate(txn *newrelic.Transaction) error {
	var err error
	defer txn.StartSegment("config_validation").End()

	if err = a.Authentication.Validate(); err != nil {
		return err
	}

	if err = a.Cache.Validate(); err != nil {
		return err
	}

	if err = a.Db.Validate(); err != nil {
		return err
	}

	if err = a.NewRelic.Validate(); err != nil {
		return err
	}

	if err = a.Paymail.Validate(); err != nil {
		return err
	}

	if err = a.Server.Validate(); err != nil {
		return err
	}

	return nil
}
