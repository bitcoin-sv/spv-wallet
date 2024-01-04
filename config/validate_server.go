package config

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation"
)

// Validate checks the configuration for specific rules
func (s *ServerConfig) Validate() error {
	// Set defaults
	if s.IdleTimeout.String() == "0s" {
		return errors.New("Idle timeout needs to be set for server configuration")
	}
	if s.ReadTimeout.String() == "0s" {
		return errors.New("Read timeout needs to be set for server configuration")
	}
	if s.WriteTimeout.String() == "0s" {
		return errors.New("Write timeout needs to be set for server configuration")
	}

	if s.Port < 10 || s.Port > 999999 {
		return errors.New("Server port outside of bounds")
	}

	return validation.ValidateStruct(s,
		validation.Field(&s.IdleTimeout, validation.Required),
		validation.Field(&s.ReadTimeout, validation.Required),
		validation.Field(&s.WriteTimeout, validation.Required),
	)
}
