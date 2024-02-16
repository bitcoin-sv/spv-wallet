package mappings

import (
	"github.com/bitcoin-sv/spv-wallet/engine"
	spvwalletmodels "github.com/bitcoin-sv/spv-wallet/models"
)

// MapToAdminStatsContract will map the model from spv-wallet to the spv-wallet-models contract
func MapToAdminStatsContract(s *engine.AdminStats) *spvwalletmodels.AdminStats {
	if s == nil {
		return nil
	}

	return &spvwalletmodels.AdminStats{
		Balance:            s.Balance,
		Destinations:       s.Destinations,
		PaymailAddresses:   s.PaymailAddresses,
		Transactions:       s.Transactions,
		TransactionsPerDay: s.TransactionsPerDay,
		Utxos:              s.Utxos,
		UtxosPerType:       s.UtxosPerType,
		XPubs:              s.XPubs,
	}
}
