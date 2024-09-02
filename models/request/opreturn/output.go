package opreturn

type Output struct {
	DataType DataType `json:"dataType,omitempty"`
	Data     []string `json:"data"`
}

func (o Output) GetType() string {
	return "op_return"
}
