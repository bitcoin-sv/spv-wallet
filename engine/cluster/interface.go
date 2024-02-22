package cluster

import "github.com/rs/zerolog"

// Coordinator the coordinators supported in cluster mode
type Coordinator string

var (
	// CoordinatorRedis definition
	CoordinatorRedis Coordinator = "redis"

	// CoordinatorMemory definition - use only in single server setups of SPV Wallet Engine!
	CoordinatorMemory Coordinator = "memory"
)

// Channel all keys used in cluster coordinator
type Channel string

var (
	// DestinationNew is a message sent when a new destination is created
	DestinationNew Channel = "new-destination"
)

// ClientInterface interface for the internal pub/sub functionality for clusters
type ClientInterface interface {
	pubSubService
	IsDebug() bool
	GetClusterPrefix() string
}

type pubSubService interface {
	Logger() *zerolog.Logger
	Subscribe(channel Channel, callback func(data string)) (func() error, error)
	Publish(channel Channel, data string) error
}
