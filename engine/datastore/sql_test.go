package datastore

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestClient_getSourceDatabase will test the method getSourceDatabase()
func TestClient_getSourceDatabase(t *testing.T) {
	t.Run("single write db", func(t *testing.T) {
		source, configs := getSourceDatabase(
			[]*SQLConfig{
				{
					CommonConfig: CommonConfig{
						Debug:                 true,
						MaxConnectionIdleTime: 10 * time.Second,
						MaxConnectionTime:     10 * time.Second,
						MaxIdleConnections:    1,
						MaxOpenConnections:    1,
						TablePrefix:           "test",
					},
					Driver:    MySQL.String(),
					Host:      "host-write.domain.com",
					Name:      "db_name",
					Password:  "test",
					Port:      defaultMySQLPort,
					Replica:   false,
					TimeZone:  defaultTimeZone,
					TxTimeout: defaultDatabaseTxTimeout,
					User:      "test",
				},
			},
		)
		require.NotNil(t, source)
		require.Empty(t, len(configs))
		assert.False(t, source.Replica)
		assert.Equal(t, "host-write.domain.com", source.Host)
	})

	t.Run("read vs write", func(t *testing.T) {
		source, configs := getSourceDatabase(
			[]*SQLConfig{
				{
					CommonConfig: CommonConfig{
						Debug:                 true,
						MaxConnectionIdleTime: 10 * time.Second,
						MaxConnectionTime:     10 * time.Second,
						MaxIdleConnections:    1,
						MaxOpenConnections:    1,
						TablePrefix:           "test",
					},
					Driver:    MySQL.String(),
					Host:      "host-write.domain.com",
					Name:      "db_name",
					Password:  "test",
					Port:      defaultMySQLPort,
					Replica:   false,
					TimeZone:  defaultTimeZone,
					TxTimeout: defaultDatabaseTxTimeout,
					User:      "test",
				},
				{
					CommonConfig: CommonConfig{
						Debug:                 true,
						MaxConnectionIdleTime: 10 * time.Second,
						MaxConnectionTime:     10 * time.Second,
						MaxIdleConnections:    1,
						MaxOpenConnections:    1,
						TablePrefix:           "test",
					},
					Driver:    MySQL.String(),
					Host:      "host-read.domain.com",
					Name:      "db_name",
					Password:  "test",
					Port:      defaultMySQLPort,
					Replica:   true,
					TimeZone:  defaultTimeZone,
					TxTimeout: defaultDatabaseTxTimeout,
					User:      "test",
				},
			},
		)
		require.NotNil(t, source)

		assert.False(t, source.Replica)
		assert.Equal(t, "host-write.domain.com", source.Host)

		assert.Len(t, configs, 1)
		assert.True(t, configs[0].Replica)
		assert.Equal(t, "host-read.domain.com", configs[0].Host)
	})

	t.Run("only replica, no source found", func(t *testing.T) {
		source, configs := getSourceDatabase(
			[]*SQLConfig{
				{
					CommonConfig: CommonConfig{
						Debug:                 true,
						MaxConnectionIdleTime: 10 * time.Second,
						MaxConnectionTime:     10 * time.Second,
						MaxIdleConnections:    1,
						MaxOpenConnections:    1,
						TablePrefix:           "test",
					},
					Driver:    MySQL.String(),
					Host:      "host-read.domain.com",
					Name:      "db_name",
					Password:  "test",
					Port:      defaultMySQLPort,
					Replica:   true,
					TimeZone:  defaultTimeZone,
					TxTimeout: defaultDatabaseTxTimeout,
					User:      "test",
				},
			},
		)
		require.Nil(t, source)
		assert.Len(t, configs, 1)
		assert.True(t, configs[0].Replica)
		assert.Equal(t, "host-read.domain.com", configs[0].Host)
	})
}
