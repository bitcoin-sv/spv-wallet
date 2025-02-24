package testabilities

import (
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

const transactionsOutlinesRecordURL = "/api/v2/transactions"

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

	then.Response(outlineRes).
		IsOK().
		WithJSONMatching(`{
			"format": "BEEF",
			"hex": "{{ matchBEEF }}",
			"annotations": {{ anything }}
		}`, nil)

	/*
		example response:
		{
		  "annotations" : {
		    "inputs" : {
		      "0" : {
		        "customInstructions" : [ {
		          "instruction" : "1-paymail_pki-sender@example.com_0",
		          "type" : "type42"
		        }, {
		          "instruction" : "1-destination-c046de7781f5c3160ddf96c4e9f620fd",
		          "type" : "type42"
		        } ]
		      }
		    },
		    "outputs" : {
		      "0" : {
		        "bucket" : "bsv",
		        "paymail" : {
		          "receiver" : "recipient@example.com",
		          "reference" : "6910562d06b0bd4cd9757cb328013683",
		          "sender" : "sender@example.com"
		        }
		      },
		      "1" : {
		        "bucket" : "bsv",
		        "customInstructions" : [ {
		          "instruction" : "1-destination-86baf3192970615c6f29a692095cfe91",
		          "type" : "type42"
		        } ]
		      }
		    }
		  },
		  "format" : "BEEF",
		  "hex" : "0100beef01fde8030101000094b15d18662b860d88982425f7e9f3d636417c00adb9c2134e0f3d2a2ce680e80301000000012e3f4683e173b40a20527fe5719633ba070df649983614886e90e45aecf2ac56000000006b483045022100e546b3d86802d85832a45bee2d12ae48836ddea113ec7eb5768178305f2ef3eb02201dd563486507ea9d0e306fc7588a1e7bcbbd201dc2e4b77001f0594aa5ce031e4121034d2d6d23fbcb6eefe3e80c47044e36797dcb80d0ac5e96e732ef03c3c550a116ffffffff010b000000000000001976a91494677c56fa2968644c90a517214338b4139899ce88ac000000000100010000000194b15d18662b860d88982425f7e9f3d636417c00adb9c2134e0f3d2a2ce680e8000000006a473044022006d3784daf1d0fb5144a6d40e6d83ea3c68dc89965ef993a444b45e0985fc4090220556c89008fd03a95f0d190e9bc0d449624f7a8f5de0cbe417ef32bb87f6093444121034d2d6d23fbcb6eefe3e80c47044e36797dcb80d0ac5e96e732ef03c3c550a116ffffffff010a000000000000001976a9148c1aa410ea92abf5e7da6aee0e2d91352bc75d7d88ac0000000000010000000110c22c211df903188b53d3db038fedf3b1006ff950e7ba22f6e1d553c4ba61380000000000000000000205000000000000001976a914d3043a219efdf10eda7cefe08306fd8933895e0488ac04000000000000001976a914249e02525df53a9778fe80f5243d59135a3a9fcb88ac0000000000"
		}
	*/

	getter := then.Response(outlineRes).JSONValue()

	tx, err := trx.NewTransactionFromBEEFHex(getter.GetString("hex"))
	require.NoError(ts.t, err)

	var customInstr bsv.CustomInstructions
	getter.GetAsType("annotations/inputs/0/customInstructions", &customInstr)
	tx.Inputs[0].UnlockingScriptTemplate = ts.user.P2PKHUnlockingScriptTemplate(customInstr...)

	err = tx.Sign()
	require.NoError(ts.t, err)

	signedBeefHex, err := tx.BEEFHex()
	require.NoError(ts.t, err)

	ts.app.ARC().WillRespondForBroadcastWithSeenOnNetwork(tx.TxID().String())

	recordRes, _ := outlineClient.R().
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]any{
			"hex":         signedBeefHex,
			"format":      "BEEF",
			"annotations": getter.GetField("annotations"),
		}).
		Post(transactionsOutlinesRecordURL)

	then.Response(recordRes).IsCreated()

	return &TransactionResult{
		//TxSpec:   txSpec,
		Response: recordRes,
		//Reference: reference,
		TxID: tx.TxID().String(),
	}
}
