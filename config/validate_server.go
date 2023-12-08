package config

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

// Validate checks the configuration for specific rules
func (s *ServerConfig) Validate() error {
	// Set defaults
	if s.IdleTimeout.String() == "0s" {
		duration, err := time.ParseDuration(ServerIdleTimeoutDefault)
		if err != nil {
			return err
		}
		s.IdleTimeout = duration
	}
	if s.ReadTimeout.String() == "0s" {
		duration, err := time.ParseDuration(ServerReadTimeoutDefault)
		if err != nil {
			return err
		}
		s.ReadTimeout = duration
	}
	if s.WriteTimeout.String() == "0s" {
		duration, err := time.ParseDuration(ServerWriteTimeoutDefault)
		if err != nil {
			return err
		}
		s.WriteTimeout = duration
	}

	return validation.ValidateStruct(s,
		validation.Field(&s.IdleTimeout, validation.Required),
		validation.Field(&s.ReadTimeout, validation.Required),
		validation.Field(&s.WriteTimeout, validation.Required),
		validation.Field(&s.Port, validation.Required, is.Digit, validation.Length(2, 6)),
	)
}
