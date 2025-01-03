package utils

import (
	"iter"
	"maps"

	trx "github.com/bitcoin-sv/go-sdk/transaction"
)

// CollectAncestors Gets the ancestors (up to mined) of provided transaction with no duplicates and excluding itself
func CollectAncestors(tx *trx.Transaction) iter.Seq[*trx.Transaction] {
	stack := make([]*trx.Transaction, 0)

	it := func(yield func(string, *trx.Transaction) bool) {
		if tx.MerklePath != nil {
			// empty result if provided transaction is already mined
			return
		}
		stack = append(stack, tx)
		for len(stack) > 0 {
			t := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			for _, input := range t.Inputs {
				if input.SourceTransaction == nil {
					continue
				}

				// collect source transaction of the input
				yield(input.SourceTransaction.TxID().String(), input.SourceTransaction)

				if input.SourceTransaction.MerklePath != nil {
					// don't process source transaction if it's already mined
					continue
				}
				stack = append(stack, input.SourceTransaction)
			}
		}
	}

	return maps.Values(maps.Collect(it))
}
