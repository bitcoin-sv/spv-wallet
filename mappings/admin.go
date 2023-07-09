package mappings

import (
	"github.com/BuxOrg/bux"
	buxmodels "github.com/BuxOrg/bux-models"
)

func MapToAdminStatsContract(s *bux.AdminStats) *buxmodels.AdminStats {
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
