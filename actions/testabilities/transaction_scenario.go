package testabilities

import (
	"encoding/json"
	"fmt"
	"github.com/bitcoin-sv/go-sdk/script"
	trx "github.com/bitcoin-sv/go-sdk/transaction"
	chainmodels "github.com/bitcoin-sv/spv-wallet/engine/chain/models"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"github.com/go-resty/resty/v2"
)

const transactionsOutlinesRecordURL = "/api/v2/transactions/outlines/record"

// Rename this -> just given.User.ReceivesFromExternal

// TransactionScenario provides higher-level methods for transaction testing
type TransactionScenario interface {
	// SendsToInternal simulates sending to an internal recipient
	SendsToInternal(recipient fixtures.User, amount bsv.Satoshis) *TransactionResult

	// ReceivesFromExternal simulates receiving from an external user
	ReceivesFromExternal(amount bsv.Satoshis) *TransactionResult
}

// TransactionResult contains all information about a completed transaction
type TransactionResult struct {
	TxSpec    fixtures.GivenTXSpec
	Response  *resty.Response
	Reference string
	TxID      string
}

type transactionScenario struct {
	app  SPVWalletApplicationFixture
	user fixtures.User
	t    testing.TB
}

// NewTransactionScenario creates a new transaction scenario helper for the given user
func NewTransactionScenario(app SPVWalletApplicationFixture, user fixtures.User, t testing.TB) TransactionScenario {
	return &transactionScenario{
		app:  app,
		user: user,
		t:    t,
	}
}

func (ts *transactionScenario) ReceivesFromExternal(amount bsv.Satoshis) *TransactionResult {
	client := ts.app.HttpClient().ForAnonymous()
	_, then := NewOf(ts.app, ts.t)

	requestBody := map[string]any{
		"satoshis": uint64(amount),
	}

	destRes, _ := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(requestBody).
		Post(fmt.Sprintf(
			"https://example.com/v1/bsvalias/p2p-payment-destination/%s",
			ts.user.DefaultPaymail(),
		))

	then.Response(destRes).IsOK().WithJSONMatching(`{
		"outputs": [
		  {
			"address": "{{ matchAddress }}",
			"satoshis": {{ .satoshis }},
			"script": "{{ matchHex }}"
		  }
		],
		"reference": "{{ matchHexWithLength 32 }}"
	}`, map[string]any{
		"satoshis": uint64(amount),
	})

	getter := then.Response(destRes).JSONValue()
	reference := getter.GetString("reference")
	lockingScript, err := script.NewFromHex(getter.GetString("outputs[0]/script"))
	require.NoError(ts.t, err)

	// Step 2: Call beef capability
	txSpec := fixtures.GivenTX(ts.t).
		WithInput(uint64(amount+1)).
		WithOutputScript(uint64(amount), lockingScript)

	ts.app.ARC().WillRespondForBroadcast(200, &chainmodels.TXInfo{
		TxID:     txSpec.ID(),
		TXStatus: chainmodels.SeenOnNetwork,
	})

	ts.app.BHS().WillRespondForMerkleRootsVerify(200, &chainmodels.MerkleRootsConfirmations{
		ConfirmationState: chainmodels.MRConfirmed,
	})

	beefRes, _ := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]any{
			"beef":      txSpec.BEEF(),
			"reference": reference,
			"metadata": map[string]any{
				"sender": fixtures.ExternalFaucet.DefaultPaymail(),
			},
		}).
		Post(fmt.Sprintf(
			"https://example.com/v1/bsvalias/beef/%s",
			ts.user.DefaultPaymail(),
		))

	then.Response(beefRes).IsOK()

	return &TransactionResult{
		TxSpec:    txSpec,
		Response:  beefRes,
		Reference: reference,
		TxID:      txSpec.ID(),
	}
}

func (ts *transactionScenario) SendsToInternal(recipient fixtures.User, amount bsv.Satoshis) *TransactionResult {
	_, then := NewOf(ts.app, ts.t)
	// replace with proper http client + check for unauthorized
	client := ts.app.HttpClient().ForAnonymous()

	requestBody := map[string]any{
		"satoshis": uint64(amount),
	}

	destRes, _ := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(requestBody).
		Post(fmt.Sprintf(
			"https://example.com/v1/bsvalias/p2p-payment-destination/%s",
			recipient.DefaultPaymail(),
		))

	then.Response(destRes).IsOK()

	getter := then.Response(destRes).JSONValue()
	reference := getter.GetString("reference")
	lockingScript, err := script.NewFromHex(getter.GetString("outputs[0]/script"))
	require.NoError(ts.t, err)

	outlineClient := ts.app.HttpClient().ForGivenUser(ts.user)
	outlineBody := map[string]any{
		"outputs": []map[string]any{
			{
				"type":     "paymail",
				"to":       recipient.DefaultPaymail(),
				"satoshis": uint64(amount),
			},
		},
	}

	outlineRes, _ := outlineClient.R().
		SetHeader("Content-Type", "application/json").
		SetBody(outlineBody).
		Post("/api/v2/transactions/outlines")

	then.Response(outlineRes).IsOK()

	var outline struct {
		Hex         string                 `json:"hex"`
		Format      string                 `json:"format"`
		Annotations map[string]interface{} `json:"annotations"`
	}
	err = json.Unmarshal(outlineRes.Body(), &outline)
	require.NoError(ts.t, err)

	tx, err := trx.NewTransactionFromBEEFHex(outline.Hex)
	require.NoError(ts.t, err)

	for i := range tx.Inputs {
		// pass custom instructions from outline annotations
		tx.Inputs[i].UnlockingScriptTemplate = ts.user.P2PKHUnlockingScriptTemplate()
	}

	err = tx.Sign()
	require.NoError(ts.t, err)

	signedBeefHex, err := tx.BEEFHex()
	require.NoError(ts.t, err)

	ts.app.ARC().WillRespondForBroadcastWithSeenOnNetwork(tx.TxID().String())

	recordRes, _ := outlineClient.R().
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]any{
			"hex":    signedBeefHex,
			"format": "BEEF",
			"annotations": map[string]any{
				"outputs": map[string]any{
					"0": map[string]any{
						"bucket": "bsv",
						"paymail": map[string]any{
							"receiver":  recipient.DefaultPaymail(),
							"reference": reference,
							"sender":    ts.user.DefaultPaymail(),
						},
					},
				},
			},
		}).
		Post("/api/v2/transactions/outlines/record")

	then.Response(recordRes).IsCreated()

	txSpec := fixtures.GivenTX(ts.t).
		WithSender(ts.user).
		WithRecipient(recipient).
		WithInput(uint64(amount+1)).
		WithOutputScript(uint64(amount), lockingScript)

	return &TransactionResult{
		TxSpec:    txSpec,
		Response:  recordRes,
		Reference: reference,
		TxID:      tx.TxID().String(),
	}
}
