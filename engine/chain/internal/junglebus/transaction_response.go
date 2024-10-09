package junglebus

// TransactionResponse represents a transaction received from the junglebus.gorillapool external service.
type TransactionResponse struct {
	Transaction []byte `json:"transaction"`
	// There are more fields in the struct, but we only need (support) the Transaction field.
}
