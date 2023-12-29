package config

import (
	"errors"
	"strconv"

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

	port := strconv.Itoa(s.Port)

	return validation.ValidateStruct(s,
		validation.Field(&s.IdleTimeout, validation.Required),
		validation.Field(&s.ReadTimeout, validation.Required),
		validation.Field(&s.WriteTimeout, validation.Required),
		validation.Field(&port, validation.Required, validation.Length(2, 6)),
	)
}
