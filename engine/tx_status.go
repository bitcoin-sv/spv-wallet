package engine

type TxStatus string

const (
	TxStatusCreated     TxStatus = "CREATED"
	TxStatusSent        TxStatus = "SENT"
	TxStatusBroadcasted TxStatus = "BROADCASTED"
	TxStatusMined       TxStatus = "MINED"
	TxStatusReverted    TxStatus = "REVERTED"
	TxStatusProblematic TxStatus = "PROBLEMATIC"
)
