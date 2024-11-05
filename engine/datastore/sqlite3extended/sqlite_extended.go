// Package sqlite3extended is a workaround for disabled by default math functions in sqlite3
// Unfortunately those math functions can be only enabled with build tag,
// Which means that we would need to force everyone, who is running any command like go build|run|test,
// to include also -tag "sqlite_math_functions"
// which looks like pretty big overhead
// and potential source of many registered issues just because someone overlooked that he needs to set this tag.
package sqlite3extended

import (
	"database/sql"
	"math"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	sqlite "github.com/mattn/go-sqlite3"
)

// NAME is the name of the driver registered by this package.
const NAME = "sqlite3_extended"

func init() {
	sql.Register(NAME, &sqlite.SQLiteDriver{
		ConnectHook: func(conn *sqlite.SQLiteConn) error {
			if err := conn.RegisterFunc("ceil", math.Ceil, true); err != nil {
				return spverrors.Wrapf(err, "error when registering ceil function in sqlite")
			}
			return nil
		},
	})
}
