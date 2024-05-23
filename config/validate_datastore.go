package config

import (
	"errors"

	"github.com/bitcoin-sv/spv-wallet/engine/datastore"
)

// Validate checks the configuration for specific rules
func (d *DbConfig) Validate() error {
	if d.Datastore.Engine == datastore.Empty || d.Datastore.Engine == "" {
		return errors.New("missing a valid datastore engine")
	}

	if d.Datastore.Engine == datastore.SQLite {
		if d.SQLite == nil {
			return errors.New("missing sqlite config")
		}
	} else if d.Datastore.Engine == datastore.PostgreSQL {
		if d.SQL == nil {
			return errors.New("missing sql config")
		} else if len(d.SQL.Host) == 0 {
			return errors.New("missing sql host")
		} else if len(d.SQL.User) == 0 {
			return errors.New("missing sql username")
		} else if len(d.SQL.Name) == 0 {
			return errors.New("missing sql db name")
		}
	} else if d.Datastore.Engine == datastore.MongoDB {
		if d.Mongo == nil {
			return errors.New("missing mongo config")
		} else if len(d.Mongo.URI) == 0 {
			return errors.New("missing mongo uri")
		} else if len(d.Mongo.DatabaseName) == 0 {
			return errors.New("missing mongo database name")
		}
	}
	return nil
}
