package base

import "github.com/bitcoin-sv/spv-wallet/config"

// APIBase represents server with base API endpoints
type APIBase struct {
	config *config.AppConfig
}

// NewAPIBase creates a new APIBase with config
func NewAPIBase(config *config.AppConfig) APIBase {
	return APIBase{config}
}
