package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestServerConfig_Validate will test the method Validate()
func TestServerConfig_Validate(t *testing.T) {
	t.Parallel()

	defaultAppConfig := getDefaultAppConfig()
	idleTimeout := defaultAppConfig.Server.IdleTimeout
	readTimeout := defaultAppConfig.Server.ReadTimeout
	writeTimeout := defaultAppConfig.Server.WriteTimeout

	t.Run("port is required", func(t *testing.T) {
		s := ServerConfig{
			IdleTimeout:  idleTimeout,
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
			Port:         "",
		}
		err := s.Validate()
		assert.Error(t, err)
	})

	t.Run("port is too big", func(t *testing.T) {
		s := ServerConfig{
			IdleTimeout:  idleTimeout,
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
			Port:         "1234567",
		}
		err := s.Validate()
		assert.Error(t, err)
	})

	t.Run("valid server config", func(t *testing.T) {
		s := ServerConfig{
			IdleTimeout:  idleTimeout,
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
			Port:         "3000",
		}
		err := s.Validate()
		assert.NoError(t, err)
	})

	t.Run("default timeouts", func(t *testing.T) {
		s := ServerConfig{
			IdleTimeout:  0,
			ReadTimeout:  0,
			WriteTimeout: 0,
			Port:         "3000",
		}
		err := s.Validate()
		assert.Error(t, err)
	})
}
