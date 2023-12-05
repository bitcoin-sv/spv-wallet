package config

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

// Validate checks the configuration for specific rules
func (s *ServerConfig) Validate() error {

	// Set defaults
	if s.IdleTimeout.String() == "0s" {
		s.IdleTimeout = DefaultHTTPRequestIdleTimeout
	}
	if s.ReadTimeout.String() == "0s" {
		s.ReadTimeout = DefaultHTTPRequestReadTimeout
	}
	if s.WriteTimeout.String() == "0s" {
		s.WriteTimeout = DefaultHTTPRequestWriteTimeout
	}

	return validation.ValidateStruct(s,
		validation.Field(&s.IdleTimeout, validation.Required),
		validation.Field(&s.ReadTimeout, validation.Required),
		validation.Field(&s.WriteTimeout, validation.Required),
		validation.Field(&s.Port, validation.Required, is.Digit, validation.Length(2, 6)),
	)
}
