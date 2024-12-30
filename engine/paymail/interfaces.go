package paymail

import "github.com/bitcoin-sv/spv-wallet/engine/database"

// Repository is an interface for the paymail repository
type Repository interface {
	GetPaymailByAlias(alias, domain string) (*database.Paymail, error)
}
