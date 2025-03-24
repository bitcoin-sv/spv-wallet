package clienterr

// ClientErrorDefinition is a definition of a client error.
type ClientErrorDefinition struct {
	title    string
	typeName string
	httpCode int
}

// Detailed is a shortcut for creating a new client error with a detailed message.
func (c ClientErrorDefinition) Detailed(errType string, detail string, args ...any) *Builder {
	return c.New().Detailed(errType, detail, args...)
}

// New creates a new client error builder.
func (c ClientErrorDefinition) New() *Builder {
	b := &Builder{from: &c}
	b.problemDetails.Status = c.httpCode
	b.problemDetails.Title = c.title

	b.problemDetails.Type = c.typeName
	return b
}
