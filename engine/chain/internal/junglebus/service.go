package junglebus

import (
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog"
)

// Service for Junglebus/Gorillapool requests.
type Service struct {
	logger     zerolog.Logger
	httpClient *resty.Client
}

// NewJunglebusService creates a new arc service.
func NewJunglebusService(logger zerolog.Logger, httpClient *resty.Client) *Service {
	return &Service{
		logger:     logger,
		httpClient: httpClient,
	}
}
