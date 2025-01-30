package opreturn

// Output is a struct that represents the output containing OP_RETURN script
type Output struct {
	DataType DataType `json:"dataType,omitempty"`
	Data     []string `json:"data"`
}

// GetType returns the string typename of the output
func (o Output) GetType() string {
	return "op_return"
}

// Output is a struct that represents the output containing OP_RETURN script
type OutputT struct {
	DataType DataType `json:"dataType,omitempty"`
	Data     []string `json:"data"`
}

// GetType returns the string typename of the output
func (o OutputT) GetType() string {
	return "opT_return"
}
