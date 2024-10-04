package engine

import (
	"context"
	"testing"

	"github.com/bitcoin-sv/go-paymail"
	"github.com/bitcoin-sv/spv-wallet/engine/datastore"
	"github.com/mrz1836/go-cachestore"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// todo: finish unit tests!

// TestClient_Debug will test the method Debug()
func TestClient_Debug(t *testing.T) {
	t.Parallel()

	t.Run("load basic with debug", func(t *testing.T) {
		tc, err := NewClient(
			context.Background(),
			DefaultClientOpts(false, true)...,
		)
		require.NoError(t, err)
		require.NotNil(t, tc)
		defer CloseClient(context.Background(), t, tc)

		assert.Equal(t, false, tc.IsDebug())

		tc.Debug(true)

		assert.Equal(t, true, tc.IsDebug())
		assert.Equal(t, true, tc.Cachestore().IsDebug())
		assert.Equal(t, true, tc.Datastore().IsDebug())
	})
}

// TestClient_IsDebug will test the method IsDebug()
func TestClient_IsDebug(t *testing.T) {
	t.Parallel()

	t.Run("basic debug checks", func(t *testing.T) {
		tc, err := NewClient(
			context.Background(),
			DefaultClientOpts(false, true)...,
		)
		require.NoError(t, err)
		require.NotNil(t, tc)
		defer CloseClient(context.Background(), t, tc)

		assert.Equal(t, false, tc.IsDebug())

		tc.Debug(true)

		assert.Equal(t, true, tc.IsDebug())
	})
}

// TestClient_Version will test the method Version()
func TestClient_Version(t *testing.T) {
	t.Parallel()

	t.Run("check version", func(t *testing.T) {
		tc, err := NewClient(
			context.Background(),
			DefaultClientOpts(false, true)...,
		)
		require.NoError(t, err)
		require.NotNil(t, tc)
		defer CloseClient(context.Background(), t, tc)

		assert.Equal(t, version, tc.Version())
	})
}

// TestClient_Cachestore will test the method Cachestore()
func TestClient_Cachestore(t *testing.T) {
	t.Parallel()

	t.Run("no options, panic", func(t *testing.T) {
		assert.Panics(t, func() {
			c := new(Client)
			assert.Nil(t, c.Cachestore())
		})
	})

	t.Run("valid cachestore", func(t *testing.T) {
		tc, err := NewClient(
			context.Background(),
			DefaultClientOpts(false, true)...,
		)
		require.NoError(t, err)
		defer CloseClient(context.Background(), t, tc)

		assert.NotNil(t, tc.Cachestore())
		assert.IsType(t, &cachestore.Client{}, tc.Cachestore())
	})
}

// TestClient_Datastore will test the method Datastore()
func TestClient_Datastore(t *testing.T) {
	t.Parallel()

	t.Run("no options, panic", func(t *testing.T) {
		assert.Panics(t, func() {
			c := new(Client)
			assert.Nil(t, c.Datastore())
		})
	})

	t.Run("valid datastore", func(t *testing.T) {
		tc, err := NewClient(
			context.Background(),
			DefaultClientOpts(false, true)...,
		)
		require.NoError(t, err)
		defer CloseClient(context.Background(), t, tc)

		assert.NotNil(t, tc.Datastore())
		assert.IsType(t, &datastore.Client{}, tc.Datastore())
	})
}

// TestClient_PaymailClient will test the method PaymailClient()
func TestClient_PaymailClient(t *testing.T) {
	t.Parallel()

	t.Run("no options, panic", func(t *testing.T) {
		assert.Panics(t, func() {
			c := new(Client)
			assert.Nil(t, c.PaymailClient())
		})
	})

	t.Run("valid paymail client", func(t *testing.T) {
		tc, err := NewClient(
			context.Background(),
			DefaultClientOpts(false, true)...,
		)
		require.NoError(t, err)
		defer CloseClient(context.Background(), t, tc)

		assert.NotNil(t, tc.PaymailClient())
		assert.IsType(t, &paymail.Client{}, tc.PaymailClient())
	})
}

// TestClient_GetPaymailConfig will test the method GetPaymailConfig()
func TestClient_GetPaymailConfig(t *testing.T) {
	t.Parallel()

	t.Run("no options, panic", func(t *testing.T) {
		assert.Panics(t, func() {
			c := new(Client)
			assert.Nil(t, c.GetPaymailConfig())
		})
	})

	t.Run("valid paymail server config", func(t *testing.T) {
		opts := DefaultClientOpts(false, true)
		opts = append(opts, WithPaymailSupport(
			[]string{"example.com"},
			defaultSenderPaymail,
			false, false,
		))

		tc, err := NewClient(context.Background(), opts...)
		require.NoError(t, err)
		defer CloseClient(context.Background(), t, tc)

		assert.NotNil(t, tc.GetPaymailConfig())
		assert.IsType(t, &PaymailServerOptions{}, tc.GetPaymailConfig())
	})
}

// TestPaymailOptions_Client will test the method Client()
func TestPaymailOptions_Client(t *testing.T) {
	t.Parallel()

	t.Run("no client", func(t *testing.T) {
		p := new(paymailOptions)
		assert.Nil(t, p.Client())
	})

	t.Run("valid paymail client", func(t *testing.T) {
		tc, err := NewClient(
			context.Background(),
			DefaultClientOpts(false, true)...,
		)
		require.NoError(t, err)
		assert.NotNil(t, tc.PaymailClient())
		defer CloseClient(context.Background(), t, tc)

		assert.IsType(t, &paymail.Client{}, tc.PaymailClient())
		assert.NotNil(t, tc.PaymailClient())
		assert.IsType(t, &paymail.Client{}, tc.PaymailClient())
	})
}

// TestPaymailOptions_FromSender will test the method FromSender()
func TestPaymailOptions_FromSender(t *testing.T) {
	t.Parallel()

	t.Run("no sender, use default", func(t *testing.T) {
		p := &paymailOptions{
			serverConfig: &PaymailServerOptions{},
		}
		assert.Equal(t, defaultSenderPaymail, p.FromSender())
	})

	t.Run("custom sender set", func(t *testing.T) {
		p := &paymailOptions{
			serverConfig: &PaymailServerOptions{
				DefaultFromPaymail: "from@domain.com",
			},
		}
		assert.Equal(t, "from@domain.com", p.FromSender())
	})
}

// TestPaymailOptions_ServerConfig will test the method ServerConfig()
func TestPaymailOptions_ServerConfig(t *testing.T) {
	// t.Parallel()

	t.Run("no server config", func(t *testing.T) {
		p := new(paymailOptions)
		assert.Nil(t, p.ServerConfig())
	})

	t.Run("valid server config", func(t *testing.T) {
		logger := zerolog.Nop()
		opts := DefaultClientOpts(false, true)
		opts = append(opts, WithPaymailSupport(
			[]string{"example.com"},
			defaultSenderPaymail,
			false, false,
		),
			WithLogger(&logger))

		tc, err := NewClient(context.Background(), opts...)
		require.NoError(t, err)
		defer CloseClient(context.Background(), t, tc)

		assert.NotNil(t, tc.GetPaymailConfig())
		assert.IsType(t, &PaymailServerOptions{}, tc.GetPaymailConfig())
	})
}
