package outlines

import (
	"context"

	primitives "github.com/bitcoin-sv/go-sdk/primitives/ec"
	sdk "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/bsv"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction"
	bsvmodel "github.com/bitcoin-sv/spv-wallet/models/bsv"
)

// PaymailAddressService is a component that provides methods for working with paymail address.
type PaymailAddressService interface {
	HasPaymailAddress(ctx context.Context, userID string, address string) (bool, error)
	GetDefaultPaymailAddress(ctx context.Context, userID string) (string, error)
}

// UTXOSelector is a component that provides methods for selecting UTXOs of given user to fund a transaction.
type UTXOSelector interface {
	Select(ctx context.Context, tx *sdk.Transaction, userID string) ([]*UTXO, bsvmodel.Satoshis, error)
}

// Service is a service for creating transaction outlines.
type Service interface {
	CreateBEEF(ctx context.Context, spec *TransactionSpec) (*Transaction, error)
	CreateRawTx(ctx context.Context, spec *TransactionSpec) (*Transaction, error)
}

// UsersService is a service for working with users.
type UsersService interface {
	GetPubKey(ctx context.Context, userID string) (*primitives.PublicKey, error)
}

// TransactionBEEFService provides functionality to generate a BEEF-encoded
// hex string from a given Bitcoin transaction.
type TransactionBEEFService interface {
	PrepareBEEF(ctx context.Context, tx *sdk.Transaction) (string, error)
}

// UTXO represents an unspent transaction output.
type UTXO struct {
	TxID string
	Vout uint32
	bsvmodel.CustomInstructions
}

// Transaction represents a transaction outline.
type Transaction struct {
	Hex         bsv.TxHex
	Annotations transaction.Annotations
}
