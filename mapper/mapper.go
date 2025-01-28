package mapper

type Mapper struct {
	errors error
}

func New() *Mapper {
	return &Mapper{}
}

func (m *Mapper) Errors() error {
	return m.errors
}

func MapSlice[Source any, Output any](source []Source, itemParser func(source Source) Output) []Output {
	result := make([]Output, 0, len(source))
	for _, item := range source {
		result = append(result, itemParser(item))
	}
	return result
}
