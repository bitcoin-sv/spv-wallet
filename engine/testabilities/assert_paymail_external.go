package testabilities

import (
	"github.com/bitcoin-sv/spv-wallet/engine/tester/jsonrequire"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/paymailmock"
	"github.com/stretchr/testify/require"
	"testing"
)

type externalClientAssertions struct {
	t           testing.TB
	require     *require.Assertions
	mockPaymail *paymailmock.PaymailClientMock
}

func (e *externalClientAssertions) ReceivedBeefTransaction(sender, beef, reference string) {
	urlRegex := "beef"
	details := e.mockPaymail.GetCallByRegex(urlRegex)
	e.require.NotNil(details, "Expected call to %s", urlRegex)

	jsonrequire.Match(e.t, `
		{
			"beef": "{{ .beef }}",
			"decodedBeef": null,
			"hex": "",
			"metadata": {
				"sender": "{{ .sender }}"
			},
			"reference": "{{ .reference }}"
		}
	`, map[string]any{
		"beef":      beef,
		"sender":    sender,
		"reference": reference,
	}, string(details.RequestBody))
}
