package tester

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestPaymailMockClient will test the method PaymailMockClient()
func TestPaymailMockClient(t *testing.T) {
	t.Run("valid client", func(t *testing.T) {
		client, err := PaymailMockClient([]string{testDomain})
		require.NoError(t, err)
		require.NotNil(t, client)
	})

	t.Run("no domain", func(t *testing.T) {
		client, err := PaymailMockClient(nil)
		require.NoError(t, err)
		require.NotNil(t, client)
	})
}
