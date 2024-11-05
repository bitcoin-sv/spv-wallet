package record

import "github.com/rs/zerolog"

// Service for recording transactions
type Service struct {
	repo        Repository
	broadcaster Broadcaster
	logger      zerolog.Logger
}

// NewService creates a new service for transactions
func NewService(logger zerolog.Logger, repo Repository, broadcaster Broadcaster) *Service {
	return &Service{
		repo:        repo,
		broadcaster: broadcaster,
		logger:      logger,
	}
}
