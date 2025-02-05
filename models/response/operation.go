package response

import "time"

// Operation represents a user's operation on a transaction.
type Operation struct {
	CreatedAt    time.Time `json:"createdAt" example:"2024-02-26T11:00:28.069911Z"`
	Value        int64     `json:"value" example:"1234"`
	TxID         string    `json:"txID" example:"bb8593f85ef8056a77026ad415f02128f3768906de53e9e8bf8749fe2d66cf50"`
	Type         string    `json:"type" example:"incoming" enums:"incoming,outgoing"`
	Counterparty string    `json:"counterparty" example:"alice@example.com"`
	TxStatus     string    `json:"txStatus" example:"MINED" enums:"CREATED,BROADCASTED,MINED,REVERTED,PROBLEMATIC"`
}
