package engine

import (
	"context"
)

// AdminStats are statistics about the SPV Wallet server
type AdminStats struct {
	Balance            int64                  `json:"balance"`
	Destinations       int64                  `json:"destinations"`
	PaymailAddresses   int64                  `json:"paymail_addresses"`
	Transactions       int64                  `json:"transactions"`
	TransactionsPerDay map[string]interface{} `json:"transactions_per_day"`
	Utxos              int64                  `json:"utxos"`
	UtxosPerType       map[string]interface{} `json:"utxos_per_type"`
	XPubs              int64                  `json:"xpubs"`
}

// GetStats will get stats for the SPV Wallet Console (admin)
func (c *Client) GetStats(ctx context.Context, opts ...ModelOps) (*AdminStats, error) {

	// Set the default model options
	defaultOpts := c.DefaultModelOptions(opts...)

	var (
		destinationsCount   int64
		err                 error
		paymailAddressCount int64
		transactionsCount   int64
		transactionsPerDay  map[string]interface{}
		utxosCount          int64
		utxosPerType        map[string]interface{}
		xpubsCount          int64
	)

	// Get the destination count
	if destinationsCount, err = getDestinationsCount(
		ctx, nil, nil, defaultOpts...,
	); err != nil {
		return nil, err
	}

	// Get the transaction count
	if transactionsCount, err = getTransactionsCount(
		ctx, nil, nil, defaultOpts...,
	); err != nil {
		return nil, err
	}

	// Get the paymail address count
	conditions := map[string]interface{}{
		"deleted_at": nil,
	}
	if paymailAddressCount, err = getPaymailAddressesCount(
		ctx, nil, conditions, defaultOpts...,
	); err != nil {
		return nil, err
	}

	// Get the utxo count
	if utxosCount, err = getUtxosCount(
		ctx, nil, nil, defaultOpts...,
	); err != nil {
		return nil, err
	}

	// Get the xpub count
	if xpubsCount, err = getXPubsCount(
		ctx, nil, nil, defaultOpts...,
	); err != nil {
		return nil, err
	}

	// Get the transactions per day count
	if transactionsPerDay, err = getTransactionsAggregate(
		ctx, nil, nil, "created_at", defaultOpts...,
	); err != nil {
		return nil, err
	}

	// Get the utxos per day count
	if utxosPerType, err = getUtxosAggregate(
		ctx, nil, nil, "type", defaultOpts...,
	); err != nil {
		return nil, err
	}

	// Return the statistics
	return &AdminStats{
		Balance:            0,
		Destinations:       destinationsCount,
		PaymailAddresses:   paymailAddressCount,
		Transactions:       transactionsCount,
		TransactionsPerDay: transactionsPerDay,
		Utxos:              utxosCount,
		UtxosPerType:       utxosPerType,
		XPubs:              xpubsCount,
	}, nil
}
