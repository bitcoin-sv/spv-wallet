package engine

import "context"

// transactionInterface is used for extending or mocking transaction methods
type transactionInterface interface {
	getDestinationByLockingScript(ctx context.Context, lockingScript string, opts ...ModelOps) (*Destination, error)
	getDestinationByAddress(ctx context.Context, address string, opts ...ModelOps) (*Destination, error)
	getUtxo(ctx context.Context, txID string, index uint32, opts ...ModelOps) (*Utxo, error)
}

// transactionService is an obj using transactionInterface
type transactionService struct{}

// getDestinationByLockingScript will get a destination by locking script
func (x transactionService) getDestinationByLockingScript(ctx context.Context,
	lockingScript string, opts ...ModelOps,
) (*Destination, error) {
	return getDestinationByLockingScript(ctx, lockingScript, opts...)
}

// getDestinationByAddress will get a destination by address
func (x transactionService) getDestinationByAddress(ctx context.Context,
	address string, opts ...ModelOps,
) (*Destination, error) {
	return getDestinationByAddress(ctx, address, opts...)
}

// getUtxo will get an utxo given the conditions
func (x transactionService) getUtxo(ctx context.Context, txID string, index uint32,
	opts ...ModelOps,
) (*Utxo, error) {
	return getUtxo(ctx, txID, index, opts...)
}
