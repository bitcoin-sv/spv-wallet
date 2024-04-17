package cluster

import (
	"context"

	"github.com/rs/zerolog"
)

// MemoryPubSub struct
type MemoryPubSub struct {
	ctx       context.Context
	callbacks map[string]func(data string)
	debug     bool
	logger    *zerolog.Logger
	prefix    string
}

// NewMemoryPubSub create a new memory pub/sub client
// this is the default (mock) implementation for the internal pub/sub communications on standalone servers
// for clusters, use another solution, like RedisPubSub
func NewMemoryPubSub(ctx context.Context) (*MemoryPubSub, error) {

	return &MemoryPubSub{
		ctx:       ctx,
		callbacks: make(map[string]func(data string)),
	}, nil
}

// Logger returns the logger to use
func (m *MemoryPubSub) Logger() *zerolog.Logger {
	return m.logger
}

// Subscribe to a channel
func (m *MemoryPubSub) Subscribe(channel Channel, callback func(data string)) (func() error, error) {

	channelName := m.prefix + string(channel)
	m.callbacks[channelName] = callback

	return func() error {
		delete(m.callbacks, channelName)
		return nil
	}, nil
}

// Publish to a channel
func (m *MemoryPubSub) Publish(channel Channel, data string) error {

	channelName := m.prefix + string(channel)
	callback, ok := m.callbacks[channelName]
	if ok {
		callback(data)
	}

	return nil
}
