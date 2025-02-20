package testabilities

import (
	"encoding/json"
	trx "github.com/bitcoin-sv/go-sdk/transaction"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"github.com/go-resty/resty/v2"
)

const transactionsOutlinesRecordURL = "/api/v2/transactions/outlines/record"

// TransactionScenario provides higher-level methods for transaction testing
type TransactionScenario interface {
	// SendsToInternal simulates sending to an internal recipient
	SendsToInternal(recipient fixtures.User, amount bsv.Satoshis) *TransactionResult

	// ReceivesFromExternal simulates receiving from an external user
	ReceivesFromExternal(amount bsv.Satoshis, ref string) *TransactionResult
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

func (ts *transactionScenario) ReceivesFromExternal(amount bsv.Satoshis, ref string) *TransactionResult {
	externalSender := fixtures.ExternalFaucet

	txSpec := fixtures.GivenTX(ts.t).
		WithSender(externalSender).
		WithRecipient(ts.user).
		WithInput(uint64(amount)).
		WithP2PKHOutput(uint64(amount))

	ts.app.ARC().WillRespondForBroadcastWithSeenOnNetwork(txSpec.ID())

	client := ts.app.HttpClient().ForAnonymous()
	res, _ := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]any{
			"hex":    txSpec.BEEF(),
			"format": "BEEF",
			"annotations": map[string]any{
				"outputs": map[string]any{
					"0": map[string]any{
						"bucket": "bsv",
						"paymail": map[string]any{
							"receiver":  ts.user.DefaultPaymail(),
							"reference": ref,
							"sender":    externalSender.DefaultPaymail(),
						},
					},
				},
			},
		}).
		Post(transactionsOutlinesRecordURL)

	return &TransactionResult{
		TxSpec:    txSpec,
		Response:  res,
		Reference: ref,
		TxID:      txSpec.ID(),
	}
}

func (ts *transactionScenario) SendsToInternal(recipient fixtures.User, amount bsv.Satoshis) *TransactionResult {
	reference := "z2def6eb-8c3a-414f-9b27-e03f8415c9d3"

	outlineURL := "/api/v2/transactions/outlines"
	outlineClient := ts.app.HttpClient().ForGivenUser(ts.user)
	outlineRes, _ := outlineClient.R().
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]any{
			"outputs": []map[string]any{
				{
					"paymail": map[string]any{
						"to":       recipient.DefaultPaymail().Address(),
						"satoshis": uint64(amount),
					},
				},
			},
		}).
		Post(outlineURL)

	var outlineResponse struct {
		Hex         string                 `json:"hex"`
		Format      string                 `json:"format"`
		Annotations map[string]interface{} `json:"annotations"`
	}
	_ = json.Unmarshal(outlineRes.Body(), &outlineResponse)

	tx, err := trx.NewTransactionFromBEEFHex(outlineResponse.Hex)
	if err != nil {
		ts.t.Fatalf("Failed to parse transaction: %v", err)
	}

	for i := range tx.Inputs {
		tx.Inputs[i].UnlockingScriptTemplate = ts.user.P2PKHUnlockingScriptTemplate()
	}

	err = tx.Sign()
	if err != nil {
		ts.t.Fatalf("Failed to sign transaction: %v", err)
	}

	signedBeefHex, err := tx.BEEFHex()
	if err != nil {
		ts.t.Fatalf("Failed to get BEEF hex: %v", err)
	}

	ts.app.ARC().WillRespondForBroadcastWithSeenOnNetwork(tx.TxID().String())

	recordClient := ts.app.HttpClient().ForGivenUser(ts.user)
	recordRes, _ := recordClient.R().
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]any{
			"hex":    signedBeefHex,
			"format": "BEEF",
			"annotations": map[string]any{
				"outputs": map[string]any{
					"0": map[string]any{
						"bucket": "bsv",
						"paymail": map[string]any{
							"receiver":  recipient.DefaultPaymail().Address(),
							"reference": reference,
							"sender":    ts.user.DefaultPaymail().Address(),
						},
					},
				},
			},
		}).
		Post(transactionsOutlinesRecordURL)

	txSpec := fixtures.GivenTX(ts.t).
		WithSender(ts.user).
		WithRecipient(recipient).
		WithInput(uint64(amount)).
		WithP2PKHOutput(uint64(amount))

	return &TransactionResult{
		TxSpec:    txSpec,
		Response:  recordRes,
		Reference: reference,
		TxID:      tx.TxID().String(),
	}
}
