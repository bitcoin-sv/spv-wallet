package config

import (
	"errors"

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

	if err = a.Db.Datastore.Validate(); err != nil {
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

	return nil
}

// validateDatastore will check the datastore and validate basic requirements
func (a *AppConfig) validateDatastore() error {
	if a.Db.Datastore.Engine == datastore.SQLite {
		if a.Db.SQLite == nil {
			return errors.New("missing sqlite config")
		}
	} else if a.Db.Datastore.Engine == datastore.MySQL || a.Db.Datastore.Engine == datastore.PostgreSQL {
		if a.Db.SQL == nil {
			return errors.New("missing sql config")
		} else if len(a.Db.SQL.Host) == 0 {
			return errors.New("missing sql host")
		} else if len(a.Db.SQL.User) == 0 {
			return errors.New("missing sql username")
		} else if len(a.Db.SQL.Name) == 0 {
			return errors.New("missing sql db name")
		}
	} else if a.Db.Datastore.Engine == datastore.MongoDB {
		if a.Db.Mongo == nil {
			return errors.New("missing mongo config")
		} else if len(a.Db.Mongo.URI) == 0 {
			return errors.New("missing mongo uri")
		} else if len(a.Db.Mongo.DatabaseName) == 0 {
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
