package taskmanager

// Factory is the different types of task factories that are supported
type Factory string

// Supported factories
const (
	FactoryEmpty  Factory = "empty"
	FactoryMemory Factory = "memory"
	FactoryRedis  Factory = "redis"
)

// String is the string version of factory
func (f Factory) String() string {
	return string(f)
}

// IsEmpty will return true if the factory is not set
func (f Factory) IsEmpty() bool {
	return f == FactoryEmpty
}
