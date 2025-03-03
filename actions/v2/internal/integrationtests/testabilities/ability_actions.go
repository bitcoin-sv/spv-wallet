package testabilities

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/bitcoin-sv/go-sdk/script"
	trx "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/actions/testabilities"
	chainmodels "github.com/bitcoin-sv/spv-wallet/engine/chain/models"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures/txtestability"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"github.com/stretchr/testify/require"
)

const (
	transactionRecordURL = "/api/v2/transactions"
	newOutlineURL        = "/api/v2/transactions/outlines"
)

// outline represents a transaction outline that can be modified
type outline struct {
	outputs []map[string]any
}

// OutlineBuilder interface defines methods for building a transaction outline
type OutlineBuilder interface {
	WithOpReturnOutput(data []string) OutlineBuilder
	WithPaymailOutput(recipient *fixtures.User, amount bsv.Satoshis) OutlineBuilder
	SignsAndRecord() string
}

type IntegrationTestAction interface {
	Alice() ActorsActions
	Bob() ActorsActions
	Charlie() ActorsActions
}

type ActorsActions interface {
	ReceivesFromExternal(amount bsv.Satoshis) (txID string)
	SendsFundsTo(recipient *fixtures.User, amount bsv.Satoshis) string
	SendsData(data []string) string

	CreatesOutline() OutlineBuilder
}

type actions struct {
	t       testing.TB
	fixture *fixture
}

func newActions(t testing.TB, given *fixture) IntegrationTestAction {
	return &actions{
		t:       t,
		fixture: given,
	}
}

func (a *actions) Alice() ActorsActions {
	return a.fixture.alice
}

func (a *actions) Bob() ActorsActions {
	return a.fixture.bob
}

func (a *actions) Charlie() ActorsActions {
	return a.fixture.charlie
}

type user struct {
	fixtures.User
	txFixture      txtestability.TransactionsFixtures
	app            testabilities.SPVWalletApplicationFixture
	t              testing.TB
	currentOutline *outline
}

// resetOutline clears the current outline
func (u *user) resetOutline() {
	u.currentOutline = nil
}

// CreatesOutline initializes a new outline
func (u *user) CreatesOutline() OutlineBuilder {
	u.currentOutline = &outline{
		outputs: make([]map[string]any, 0),
	}
	return u
}

// WithPaymailOutput adds a paymail output to the outline
func (u *user) WithPaymailOutput(recipient *fixtures.User, amount bsv.Satoshis) OutlineBuilder {
	if u.currentOutline == nil {
		u.CreatesOutline()
	}

	output := map[string]any{
		"type":     "paymail",
		"to":       recipient.DefaultPaymail(),
		"satoshis": uint64(amount),
	}

	u.currentOutline.outputs = append(u.currentOutline.outputs, output)
	return u
}

// WithOpReturnOutput adds an OP_RETURN output to the outline
func (u *user) WithOpReturnOutput(data []string) OutlineBuilder {
	if u.currentOutline == nil {
		u.CreatesOutline()
	}

	output := map[string]any{
		"type": "op_return",
		"data": data,
	}

	u.currentOutline.outputs = append(u.currentOutline.outputs, output)
	return u
}

// SignsAndRecord creates, signs and records the transaction
func (u *user) SignsAndRecord() string {
	if u.currentOutline == nil || len(u.currentOutline.outputs) == 0 {
		u.t.Fatal("No outputs added to outline")
	}

	_, then := testabilities.NewOf(u.app, u.t)

	outlineClient := u.app.HttpClient().ForGivenUser(u.User)
	outlineBody := map[string]any{
		"outputs": u.currentOutline.outputs,
	}

	outlineRes, _ := outlineClient.R().
		SetHeader("Content-Type", "application/json").
		SetBody(outlineBody).
		Post(newOutlineURL)

	then.Response(outlineRes).
		IsOK().
		WithJSONMatching(`{
          "format": "BEEF",
          "hex": "{{ matchBEEF }}",
          "annotations": {{ anything }}
       }`, nil)

	getter := then.Response(outlineRes).JSONValue()

	hex := getter.GetString("hex")
	annotations := make(map[string]any)
	getter.GetAsType("annotations", &annotations)

	tx, err := trx.NewTransactionFromBEEFHex(hex)
	require.NoError(u.t, err)

	inputAnnotations := map[string]struct {
		CustomInstructions bsv.CustomInstructions `json:"customInstructions"`
	}{}

	inputs, ok := annotations["inputs"]
	if ok {
		inputsJSON, err := json.Marshal(inputs)
		require.NoError(u.t, err)

		err = json.Unmarshal(inputsJSON, &inputAnnotations)
		require.NoError(u.t, err)
	}

	for i, input := range tx.Inputs {
		var customInstr bsv.CustomInstructions
		if annotation, ok := inputAnnotations[fmt.Sprintf("%d", i)]; ok {
			customInstr = annotation.CustomInstructions
		}
		input.UnlockingScriptTemplate = u.P2PKHUnlockingScriptTemplate(customInstr...)
	}

	err = tx.Sign()
	require.NoError(u.t, err)

	signedHex, err := tx.BEEFHex()
	require.NoError(u.t, err)

	u.app.ARC().WillRespondForBroadcastWithSeenOnNetwork(tx.TxID().String())

	recordClient := u.app.HttpClient().ForGivenUser(u.User)
	recordRes, _ := recordClient.R().
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]any{
			"hex":         signedHex,
			"format":      "BEEF",
			"annotations": annotations,
		}).
		Post(transactionRecordURL)

	then.Response(recordRes).IsCreated()

	txID := tx.TxID().String()

	u.resetOutline()

	return txID
}

// ReceivesFromExternal receives funds from external source
func (u *user) ReceivesFromExternal(amount bsv.Satoshis) string {
	client := u.app.HttpClient().ForAnonymous()
	_, then := testabilities.NewOf(u.app, u.t)

	requestBody := map[string]any{
		"satoshis": uint64(amount),
	}

	destRes, _ := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(requestBody).
		Post(fmt.Sprintf(
			"https://example.com/v1/bsvalias/p2p-payment-destination/%s",
			u.DefaultPaymail(),
		))

	then.Response(destRes).IsOK()

	getter := then.Response(destRes).JSONValue()
	reference := getter.GetString("reference")
	lockingScript, err := script.NewFromHex(getter.GetString("outputs[0]/script"))
	require.NoError(u.t, err)

	txSpec := u.txFixture.Tx().
		WithInput(uint64(amount+1)).
		WithOutputScript(uint64(amount), lockingScript)

	u.app.ARC().WillRespondForBroadcast(200, &chainmodels.TXInfo{
		TxID:     txSpec.ID(),
		TXStatus: chainmodels.SeenOnNetwork,
	})

	u.app.BHS().WillRespondForMerkleRootsVerify(200, &chainmodels.MerkleRootsConfirmations{
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
			u.DefaultPaymail(),
		))

	then.Response(beefRes).IsOK()

	return txSpec.ID()
}

// SendsFundsTo sends funds to the recipient
func (u *user) SendsFundsTo(recipient *fixtures.User, amount bsv.Satoshis) string {
	return u.CreatesOutline().
		WithPaymailOutput(recipient, amount).
		SignsAndRecord()
}

// SendsData sends data output
func (u *user) SendsData(data []string) string {
	return u.CreatesOutline().
		WithOpReturnOutput(data).
		SignsAndRecord()
}
