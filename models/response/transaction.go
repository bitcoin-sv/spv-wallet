package response

import "github.com/bitcoin-sv/spv-wallet/models/common"

// Transaction is a model that represents a transaction.
type Transaction struct {
	// Model is a common model that contains common fields for all models.
	common.Model
	// ID is a transaction id.
	ID string `json:"id" example:"01d0d0067652f684c6acb3683763f353fce55f6496521c7d99e71e1d27e53f5c"`
	// Hex is a transaction hex.
	Hex string `json:"hex" example:"0100000002..."`
	// XpubInIDs is a slice of xpub input ids.
	XpubInIDs []string `json:"xpubInIds" example:"bb8593f85ef8056a77026ad415f02128f3768906de53e9e8bf8749fe2d66cf50"`
	// XpubOutIDs is a slice of xpub output ids.
	XpubOutIDs []string `json:"xpubOutIds" example:"2075eca10bf2688b38cd7fdad6c24562463a9a26ae505d66c480fd53165dbaa2"`
	// BlockHash is a block hash that transaction is in.
	BlockHash string `json:"blockHash" example:"0000000000000000046e81025ca6cfbd2f45c7331f650c77edc99a14d5a1f0d0"`
	// BlockHeight is a block height that transaction is in.
	BlockHeight uint64 `json:"blockHeight" example:"833505"`
	// Fee is a transaction fee.
	Fee uint64 `json:"fee" example:"1"`
	// NumberOfInputs is a number of transaction inputs.
	NumberOfInputs uint32 `json:"numberOfInputs" example:"3"`
	// NumberOfOutputs is a number of transaction outputs.
	NumberOfOutputs uint32 `json:"numberOfOutputs" example:"2"`
	// DraftID is a transaction related draft id.
	DraftID string `json:"draftId" example:"b356f7fa00cd3f20cce6c21d704cd13e871d28d714a5ebd0532f5a0e0cde63f7"`
	// TotalValue is a total input value.
	TotalValue uint64 `json:"totalValue" example:"51"`
	// OutputValue is a total output value.
	OutputValue int64 `json:"outputValue,omitempty" example:"50"`
	// Outputs represents all spv-wallet-transaction outputs. Will be shown only for admin.
	Outputs map[string]int64 `json:"outputs,omitempty" example:"92640954841510a9d95f7737a43075f22ebf7255976549de4c52e8f3faf57470:-51,9d07977d2fc14402426288a6010b4cdf7d91b61461acfb75af050b209d2d07ba:50"`
	// Status is a transaction status.
	Status string `json:"status" example:"MINED"`
	// TransactionDirection is a transaction direction (incoming/outgoing).
	TransactionDirection string `json:"direction" example:"outgoing"`
}
