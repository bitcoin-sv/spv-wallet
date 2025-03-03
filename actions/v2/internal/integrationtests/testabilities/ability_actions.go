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

type outlineBuilder struct {
	outputs []map[string]any
}

type outlineResult struct {
	hex         string
	annotations map[string]any
}

func (r *outlineResult) reset() {
	r.hex = ""
	r.annotations = nil
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

	AddPaymailOutput(recipient *fixtures.User, amount bsv.Satoshis)
	AddOpReturnOutput(data []string)
	CreateOutline()
	SignOutline()
	SendTransaction() string
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
	outlineBuilder *outlineBuilder
	currentOutline *outlineResult
}

func (u *user) initOutlineBuilder() {
	u.outlineBuilder = &outlineBuilder{
		outputs: make([]map[string]any, 0),
	}
}

// AddPaymailOutput adds a paymail output to the outline
func (u *user) AddPaymailOutput(recipient *fixtures.User, amount bsv.Satoshis) {
	if u.outlineBuilder == nil {
		u.initOutlineBuilder()
	}

	output := map[string]any{
		"type":     "paymail",
		"to":       recipient.DefaultPaymail(),
		"satoshis": uint64(amount),
	}

	u.outlineBuilder.outputs = append(u.outlineBuilder.outputs, output)
}

// AddOpReturnOutput adds an OP_RETURN output to the outline
func (u *user) AddOpReturnOutput(data []string) {
	if u.outlineBuilder == nil {
		u.initOutlineBuilder()
	}

	output := map[string]any{
		"type": "op_return",
		"data": data,
	}

	u.outlineBuilder.outputs = append(u.outlineBuilder.outputs, output)
}

// CreateOutline creates an outline from the current outputs
func (u *user) CreateOutline() {
	if u.outlineBuilder == nil || len(u.outlineBuilder.outputs) == 0 {
		u.t.Fatal("No outputs added to outline builder")
	}

	_, then := testabilities.NewOf(u.app, u.t)

	outlineClient := u.app.HttpClient().ForGivenUser(u.User)
	outlineBody := map[string]any{
		"outputs": u.outlineBuilder.outputs,
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

	u.currentOutline = &outlineResult{
		hex: getter.GetString("hex"),
	}

	getter.GetAsType("annotations", &u.currentOutline.annotations)

	u.outlineBuilder = nil
}

// SignOutline signs the current outline
func (u *user) SignOutline() {
	if u.currentOutline == nil {
		u.t.Fatal("No outline available to sign")
	}

	tx, err := trx.NewTransactionFromBEEFHex(u.currentOutline.hex)
	require.NoError(u.t, err)

	inputAnnotations := map[string]struct {
		CustomInstructions bsv.CustomInstructions `json:"customInstructions"`
	}{}

	inputs, ok := u.currentOutline.annotations["inputs"]
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

	u.currentOutline.hex = signedHex
}

// SendTransaction sends the signed transaction
func (u *user) SendTransaction() string {
	if u.currentOutline == nil {
		u.t.Fatal("No signed transaction available to send")
	}

	tx, err := trx.NewTransactionFromBEEFHex(u.currentOutline.hex)
	require.NoError(u.t, err)

	_, then := testabilities.NewOf(u.app, u.t)

	u.app.ARC().WillRespondForBroadcastWithSeenOnNetwork(tx.TxID().String())

	recordClient := u.app.HttpClient().ForGivenUser(u.User)
	recordRes, _ := recordClient.R().
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]any{
			"hex":         u.currentOutline.hex,
			"format":      "BEEF",
			"annotations": u.currentOutline.annotations,
		}).
		Post(transactionRecordURL)

	then.Response(recordRes).IsCreated()

	txID := tx.TxID().String()
	u.currentOutline.reset()

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
	u.AddPaymailOutput(recipient, amount)
	u.CreateOutline()
	u.SignOutline()

	return u.SendTransaction()
}

// SendsData sends data output
func (u *user) SendsData(data []string) string {
	u.AddOpReturnOutput(data)
	u.CreateOutline()
	u.SignOutline()

	return u.SendTransaction()
}
