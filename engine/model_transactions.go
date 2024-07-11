package engine

import (
	"context"

	"github.com/bitcoin-sv/spv-wallet/engine/chainstate"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/utils"
	"github.com/libsv/go-bt/v2"
)

// TransactionBase is the same fields share between multiple transaction models
type TransactionBase struct {
	ID  string `json:"id" toml:"id" yaml:"id" gorm:"<-:create;type:char(64);primaryKey;comment:This is the unique id (hash of the transaction hex)" bson:"_id"`
	Hex string `json:"hex" toml:"hex" yaml:"hex" gorm:"<-:create;type:text;comment:This is the raw transaction hex" bson:"hex"`

	// Private for internal use
	parsedTx *bt.Tx `gorm:"-" bson:"-"` // The go-bt version of the transaction
}

// TransactionDirection String describing the direction of the transaction (in / out)
type TransactionDirection string

const (
	// TransactionDirectionIn The transaction is coming in to the wallet of the xpub
	TransactionDirectionIn TransactionDirection = "incoming"

	// TransactionDirectionOut The transaction is going out of to the wallet of the xpub
	TransactionDirectionOut TransactionDirection = "outgoing"

	// TransactionDirectionReconcile The transaction is an internal reconciliation transaction
	TransactionDirectionReconcile TransactionDirection = "reconcile"
)

// Transaction is an object representing the BitCoin transaction
//
// Gorm related models & indexes: https://gorm.io/docs/models.html - https://gorm.io/docs/indexes.html
type Transaction struct {
	// Base model
	Model `bson:",inline"`

	// Standard transaction model base fields
	TransactionBase `bson:",inline"`

	// Model specific fields
	XpubInIDs       IDs             `json:"xpub_in_ids,omitempty" toml:"xpub_in_ids" yaml:"xpub_in_ids" gorm:"<-;type:json" bson:"xpub_in_ids,omitempty"`
	XpubOutIDs      IDs             `json:"xpub_out_ids,omitempty" toml:"xpub_out_ids" yaml:"xpub_out_ids" gorm:"<-;type:json" bson:"xpub_out_ids,omitempty"`
	BlockHash       string          `json:"block_hash" toml:"block_hash" yaml:"block_hash" gorm:"<-;type:char(64);comment:This is the related block when the transaction was mined" bson:"block_hash,omitempty"`
	BlockHeight     uint64          `json:"block_height" toml:"block_height" yaml:"block_height" gorm:"<-;type:bigint;comment:This is the related block when the transaction was mined" bson:"block_height,omitempty"`
	Fee             uint64          `json:"fee" toml:"fee" yaml:"fee" gorm:"<-create;type:bigint" bson:"fee,omitempty"`
	NumberOfInputs  uint32          `json:"number_of_inputs" toml:"number_of_inputs" yaml:"number_of_inputs" gorm:"<-;type:int" bson:"number_of_inputs,omitempty"`
	NumberOfOutputs uint32          `json:"number_of_outputs" toml:"number_of_outputs" yaml:"number_of_outputs" gorm:"<-;type:int" bson:"number_of_outputs,omitempty"`
	DraftID         string          `json:"draft_id" toml:"draft_id" yaml:"draft_id" gorm:"<-create;type:varchar(64);index;comment:This is the related draft id" bson:"draft_id,omitempty"`
	TotalValue      uint64          `json:"total_value" toml:"total_value" yaml:"total_value" gorm:"<-create;type:bigint" bson:"total_value,omitempty"`
	XpubMetadata    XpubMetadata    `json:"-" toml:"xpub_metadata" gorm:"<-;type:json;xpub_id specific metadata" bson:"xpub_metadata,omitempty"`
	XpubOutputValue XpubOutputValue `json:"-" toml:"xpub_output_value" gorm:"<-;type:json;xpub_id specific value" bson:"xpub_output_value,omitempty"`
	BUMP            BUMP            `json:"bump" toml:"bump" yaml:"bump" gorm:"<-;type:text;comment:BSV Unified Merkle Path (BUMP) Format" bson:"bump,omitempty"`
	TxStatus        string          `json:"txStatus" toml:"txStatus" yaml:"txStatus" gorm:"<-;type:varchar(64);comment:TxStatus retrieved from Arc API." bson:"txStatus,omitempty"`

	// Virtual Fields
	OutputValue int64                `json:"output_value" toml:"-" yaml:"-" gorm:"-" bson:"-,omitempty"`
	Status      SyncStatus           `json:"status" toml:"-" yaml:"-" gorm:"-" bson:"-"`
	Direction   TransactionDirection `json:"direction" toml:"-" yaml:"-" gorm:"-" bson:"-"`
	// Confirmations  uint64       `json:"-" toml:"-" yaml:"-" gorm:"-" bson:"-"`

	// Private for internal use
	draftTransaction   *DraftTransaction    `gorm:"-" bson:"-"` // Related draft transaction for processing and recording
	syncTransaction    *SyncTransaction     `gorm:"-" bson:"-"` // Related record if broadcast config is detected (create new recordNew)
	transactionService transactionInterface `gorm:"-" bson:"-"` // Used for interfacing methods
	utxos              []Utxo               `gorm:"-" bson:"-"` // json:"destinations,omitempty"
	XPubID             string               `gorm:"-" bson:"-"` // XPub of the user registering this transaction
	beforeCreateCalled bool                 `gorm:"-" bson:"-"` // Private information that the transaction lifecycle method BeforeCreate was already called
}

// TransactionGetter interface for getting transactions by their IDs
type TransactionGetter interface {
	GetTransactionsByIDs(ctx context.Context, txIDs []string) ([]*Transaction, error)
}

func emptyTx(opts ...ModelOps) *Transaction {
	return &Transaction{
		TransactionBase:    TransactionBase{},
		Model:              *NewBaseModel(ModelTransaction, opts...),
		Status:             statusComplete,
		transactionService: transactionService{},
		XpubOutputValue:    map[string]int64{},
	}
}

// baseTxFromHex creates the standard transaction model base
func baseTxFromHex(hex string, opts ...ModelOps) (*Transaction, error) {
	var btTx *bt.Tx
	var err error

	if btTx, err = bt.NewTxFromString(hex); err != nil {
		return nil, err
	}

	tx := emptyTx(opts...)
	tx.ID = btTx.TxID()
	tx.Hex = hex
	tx.parsedTx = btTx

	return tx, nil
}

// txFromHex will start a new transaction model
func txFromHex(txHex string, opts ...ModelOps) (*Transaction, error) {
	tx, err := baseTxFromHex(txHex, opts...)
	if err != nil {
		return nil, err
	}

	// Set xPub ID
	tx.setXPubID()

	return tx, nil
}

// newTransactionWithDraftID will start a new transaction model and set the draft ID
func newTransactionWithDraftID(txHex, draftID string, opts ...ModelOps) (*Transaction, error) {
	tx, err := txFromHex(txHex, opts...)
	if err != nil {
		return nil, err
	}

	tx.DraftID = draftID

	return tx, nil
}

func txFromBtTx(btTx *bt.Tx, opts ...ModelOps) *Transaction {
	tx := emptyTx(opts...)
	tx.ID = btTx.TxID()
	tx.Hex = btTx.String()
	tx.parsedTx = btTx

	return tx
}

// setXPubID will set the xPub ID on the model
func (m *Transaction) setXPubID() {
	if len(m.rawXpubKey) > 0 && len(m.XPubID) == 0 {
		m.XPubID = utils.Hash(m.rawXpubKey)
	}
}

// UpdateTransactionMetadata will update the transaction metadata by xPubID
func (m *Transaction) UpdateTransactionMetadata(xPubID string, metadata Metadata) error {
	if xPubID == "" {
		return spverrors.ErrXpubIDMisMatch
	}

	// transaction metadata is saved per xPubID
	if m.XpubMetadata == nil {
		m.XpubMetadata = make(XpubMetadata)
	}
	if m.XpubMetadata[xPubID] == nil {
		m.XpubMetadata[xPubID] = make(Metadata)
	}

	for key, value := range metadata {
		if value == nil {
			delete(m.XpubMetadata[xPubID], key)
		} else {
			m.XpubMetadata[xPubID][key] = value
		}
	}

	return nil
}

// GetModelName will get the name of the current model
func (m *Transaction) GetModelName() string {
	return ModelTransaction.String()
}

// GetID will get the ID
func (m *Transaction) GetID() string {
	return m.ID
}

// setID will set the ID from the transaction hex
func (m *Transaction) setID() (err error) {
	// Parse the hex (if not already parsed)
	if m.TransactionBase.parsedTx == nil {
		if m.TransactionBase.parsedTx, err = bt.NewTxFromString(m.Hex); err != nil {
			return
		}
	}

	// Set the true transaction ID
	m.ID = m.TransactionBase.parsedTx.TxID()

	return
}

// getValue calculates the value of the transaction
func (m *Transaction) getValues() (outputValue uint64, fee uint64) {
	// Parse the outputs
	for _, output := range m.TransactionBase.parsedTx.Outputs {
		outputValue += output.Satoshis
	}

	// Remove the "change" from the transaction if found
	// todo: this will NOT work for an "external" tx that is coming into our system?
	if m.draftTransaction != nil {
		outputValue -= m.draftTransaction.Configuration.ChangeSatoshis
		fee = m.draftTransaction.Configuration.Fee
	} else { // external transaction

		// todo: missing inputs in some tests?
		var inputValue uint64
		for _, input := range m.TransactionBase.parsedTx.Inputs {
			inputValue += input.PreviousTxSatoshis
		}

		if inputValue > 0 {
			fee = inputValue - outputValue
			outputValue -= fee
		}

		// todo: outputs we know are accumulated
	}

	// remove the fee from the value
	if outputValue > fee {
		outputValue -= fee
	}

	return
}

func (m *Transaction) setChainInfo(txInfo *chainstate.TransactionInfo) {
	m.BlockHash = txInfo.BlockHash
	m.BlockHeight = uint64(txInfo.BlockHeight)
	m.TxStatus = txInfo.TxStatus.String()
	m.setBUMP(txInfo)
}

// Converts from bc.BUMP to our BUMP struct in Transaction model
func (m *Transaction) setBUMP(txInfo *chainstate.TransactionInfo) {
	if txInfo.BUMP != nil {
		m.BUMP = bcBumpToBUMP(txInfo.BUMP)
	} else {
		m.client.Logger().Error().Msg("No BUMP found")
	}
}

func (m *Transaction) isMined() bool {
	return m.BlockHash != ""
}

// IsXpubAssociated will check if this key is associated to this transaction
func (m *Transaction) IsXpubAssociated(rawXpubKey string) bool {
	// Hash the raw key
	xPubID := utils.Hash(rawXpubKey)
	return m.IsXpubIDAssociated(xPubID)
}

// IsXpubIDAssociated will check if an xPub ID is associated
func (m *Transaction) IsXpubIDAssociated(xPubID string) bool {
	if len(xPubID) == 0 {
		return false
	}

	// On the input side
	for _, id := range m.XpubInIDs {
		if id == xPubID {
			return true
		}
	}

	// On the output side
	for _, id := range m.XpubOutIDs {
		if id == xPubID {
			return true
		}
	}
	return false
}

// Display filter the model for display
func (m *Transaction) Display() interface{} {
	// In case it was not set
	m.setXPubID()

	if len(m.XpubMetadata) > 0 && len(m.XpubMetadata[m.XPubID]) > 0 {
		if m.Metadata == nil {
			m.Metadata = make(Metadata)
		}
		for key, value := range m.XpubMetadata[m.XPubID] {
			m.Metadata[key] = value
		}
	}

	m.OutputValue = int64(0)
	if len(m.XpubOutputValue) > 0 && m.XpubOutputValue[m.XPubID] != 0 {
		m.OutputValue = m.XpubOutputValue[m.XPubID]
	}

	if m.OutputValue > 0 {
		m.Direction = TransactionDirectionIn
	} else {
		m.Direction = TransactionDirectionOut
	}

	m.XpubInIDs = nil
	m.XpubOutIDs = nil
	m.XpubMetadata = nil
	m.XpubOutputValue = nil
	return m
}

// hasOneKnownDestination will check if the transaction has at least one known destination
//
// This is used to validate if an external transaction should be recorded into the engine
func (m *TransactionBase) hasOneKnownDestination(ctx context.Context, client ClientInterface) bool {
	// todo: this can be optimized searching X records at a time vs loop->query->loop->query
	for _, output := range m.parsedTx.Outputs {
		lockingScript := output.LockingScript.String()
		destination, err := getDestinationWithCache(ctx, client, "", "", lockingScript)

		if err != nil {
			client.Logger().Error().Str("txID", m.ID).Msgf("error getting destination: %s", err.Error())
			continue
		} else if destination != nil && destination.LockingScript == lockingScript {
			return true
		}
	}
	return false
}
