package testabilities

import (
	"fmt"
	"testing"

	"github.com/bitcoin-sv/go-sdk/script"
	trx "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/actions/testabilities"
	chainmodels "github.com/bitcoin-sv/spv-wallet/engine/chain/models"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"github.com/stretchr/testify/require"
)

const (
	transactionRecordURL = "/api/v2/transactions"
	newOutlineURL        = "/api/v2/transactions/outlines"
)

type IntegrationTestAction interface {
	Alice() ActorsActions
	ARC() ARCActions
}

type ActorsActions interface {
	ReceivesFromExternal(amount bsv.Satoshis) (txID string)
	SendsTo(recipient *fixtures.User, amount bsv.Satoshis) (txID string)
}

type ARCActions interface {
	Callbacks(txInfo chainmodels.TXInfo)
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

func (a *actions) ARC() ARCActions {
	return &arcActions{
		t:       a.t,
		fixture: a.fixture,
	}
}

type user struct {
	fixtures.User
	app testabilities.SPVWalletApplicationFixture
	t   testing.TB
}

// ReceivesFromExternal simulates receiving funds from an external source
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

	txSpec := fixtures.GivenTX(u.t).
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

// SendsTo simulates sending funds to another user
func (u *user) SendsTo(recipient *fixtures.User, amount bsv.Satoshis) string {
	_, then := testabilities.NewOf(u.app, u.t)

	outlineClient := u.app.HttpClient().ForGivenUser(u.User)
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
		Post(newOutlineURL)

	then.Response(outlineRes).
		IsOK().
		WithJSONMatching(`{
          "format": "BEEF",
          "hex": "{{ matchBEEF }}",
          "annotations": {{ anything }}
       }`, nil)

	getter := then.Response(outlineRes).JSONValue()

	tx, err := trx.NewTransactionFromBEEFHex(getter.GetString("hex"))
	require.NoError(u.t, err)

	inputAnnotations := map[string]struct {
		CustomInstructions bsv.CustomInstructions `json:"customInstructions"`
	}{}
	getter.GetAsType("annotations/inputs", &inputAnnotations)

	for i, input := range tx.Inputs {
		var customInstr bsv.CustomInstructions
		if annotation, ok := inputAnnotations[fmt.Sprintf("%d", i)]; ok {
			customInstr = annotation.CustomInstructions
		}
		input.UnlockingScriptTemplate = u.P2PKHUnlockingScriptTemplate(customInstr...)
	}

	err = tx.Sign()
	require.NoError(u.t, err)

	signedBeefHex, err := tx.BEEFHex()
	require.NoError(u.t, err)

	u.app.ARC().WillRespondForBroadcastWithSeenOnNetwork(tx.TxID().String())

	recordRes, _ := outlineClient.R().
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]any{
			"hex":         signedBeefHex,
			"format":      "BEEF",
			"annotations": getter.GetField("annotations"),
		}).
		Post(transactionRecordURL)

	then.Response(recordRes).IsCreated()

	return tx.TxID().String()
}
