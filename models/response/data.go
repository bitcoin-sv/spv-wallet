package response

// Data is a response model for data stored in outputs (e.g. OP_RETURN).
type Data struct {
	ID string `json:"id"`

	Blob string `json:"blob"`
}
