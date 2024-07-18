package config

import (
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	validation "github.com/go-ozzo/ozzo-validation"
)

// Validate checks the configuration for specific rules
func (s *ServerConfig) Validate() error {
	// Set defaults
	if s.IdleTimeout.String() == "0s" {
		return spverrors.Newf("idle timeout needs to be set for server configuration")
	}
	if s.ReadTimeout.String() == "0s" {
		return spverrors.Newf("read timeout needs to be set for server configuration")
	}
	if s.WriteTimeout.String() == "0s" {
		return spverrors.Newf("write timeout needs to be set for server configuration")
	}

	if s.Port < 10 || s.Port > 65535 {
		return spverrors.Newf("server port outside of bounds")
	}

	return validation.ValidateStruct(s,
		validation.Field(&s.IdleTimeout, validation.Required),
		validation.Field(&s.ReadTimeout, validation.Required),
		validation.Field(&s.WriteTimeout, validation.Required),
		validation.Field(&s.Port, validation.Required),
	)
}
