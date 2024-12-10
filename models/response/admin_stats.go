package response

// AdminStats is a model that represents admin stats.
type AdminStats struct {
	// Balance is a total balance of all xpubs.
	Balance int64 `json:"balance"`
	// Destinations is a total number of destinations.
	Destinations int64 `json:"destinations"`
	// PaymailAddresses is a total number of paymail addresses.
	PaymailAddresses int64 `json:"paymailAddresses"`
	// Transactions is a total number of committed transactions.
	Transactions int64 `json:"transactions"`
	// TransactionsPerDay is a total number of committed transactions per day.
	TransactionsPerDay map[string]interface{} `json:"transactionsPerDay"`
	// Utxos is a total number of utxos.
	Utxos int64 `json:"utxos"`
	// UtxosPerType are utxos grouped by type.
	UtxosPerType map[string]interface{} `json:"utxosPerType"`
	// Xpubs is a total number of xpubs.
	XPubs int64 `json:"xpubs"`
}
