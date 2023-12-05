package config

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/mrz1836/go-cachestore"
	"github.com/mrz1836/go-datastore"
	"github.com/newrelic/go-agent/v3/newrelic"
)

// Validate checks the configuration for specific rules
func (a *AppConfig) Validate(txn *newrelic.Transaction) error {
	var err error
	defer txn.StartSegment("config_validation").End()

	if err = a.Authentication.Validate(); err != nil {
		return err
	}

	if err = a.Cachestore.Validate(); err != nil {
		return err
	}

	if err = a.Datastore.Validate(); err != nil {
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

	if err = a.validateCachestore(); err != nil {
		return err
	}

	if err = a.validateDatastore(); err != nil {
		return err
	}

	return validation.ValidateStruct(a,
		validation.Field(&a.Environment, validation.Required, validation.In(environments...)),
		validation.Field(&a.WorkingDirectory, validation.Required),
	)
}

// validateDatastore will check the datastore and validate basic requirements
func (a *AppConfig) validateDatastore() error {
	if a.Datastore.Engine == datastore.SQLite {
		if a.SQLite == nil {
			return errors.New("missing sqlite config")
		}
	} else if a.Datastore.Engine == datastore.MySQL || a.Datastore.Engine == datastore.PostgreSQL {
		if a.SQL == nil {
			return errors.New("missing sql config")
		} else if len(a.SQL.Host) == 0 {
			return errors.New("missing sql host")
		} else if len(a.SQL.User) == 0 {
			return errors.New("missing sql username")
		} else if len(a.SQL.Name) == 0 {
			return errors.New("missing sql db name")
		}
	} else if a.Datastore.Engine == datastore.MongoDB {
		if a.Mongo == nil {
			return errors.New("missing mongo config")
		} else if len(a.Mongo.URI) == 0 {
			return errors.New("missing mongo uri")
		} else if len(a.Mongo.DatabaseName) == 0 {
			return errors.New("missing mongo database name")
		}
	}
	return nil
}

// validateCachestore will check the cachestore and validate basic requirements
func (a *AppConfig) validateCachestore() error {
	if a.Cachestore.Engine == cachestore.Redis {
		if a.Redis == nil {
			return errors.New("missing redis config")
		} else if len(a.Redis.URL) == 0 {
			return errors.New("missing redis url")
		}
	}
	return nil
}
