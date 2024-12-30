package paymail

import "github.com/bitcoin-sv/spv-wallet/engine/database"

type Repository interface {
	GetPaymailByAlias(alias, domain string) (*database.Paymail, error)
}
