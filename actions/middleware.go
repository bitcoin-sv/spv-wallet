package actions

import "github.com/bitcoin-sv/spv-wallet/config"

// Action is the configuration for the actions and related services
type Action struct {
	AppConfig *config.AppConfig
	Services  *config.AppServices
}
