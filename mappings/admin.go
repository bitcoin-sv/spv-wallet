package mappings

import (
	"github.com/BuxOrg/bux"
	buxmodels "github.com/BuxOrg/bux-models"
)

// MapToAdminStatsContract will map the model from bux to the bux-models contract
func MapToAdminStatsContract(s *bux.AdminStats) *buxmodels.AdminStats {
	if s == nil {
		return nil
	}

	return &buxmodels.AdminStats{
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
