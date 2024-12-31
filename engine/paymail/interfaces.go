package paymail

import (
	"context"
	"github.com/bitcoin-sv/spv-wallet/engine/database"
)

// Repository is an interface for the paymail repository
type Repository interface {
	GetPaymailByAlias(alias, domain string) (*database.Paymail, error)
	SaveAddress(ctx context.Context, userRow *database.User, addressRow *database.Address) error
}
