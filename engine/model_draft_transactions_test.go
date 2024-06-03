package engine

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/bitcoin-sv/spv-wallet/engine/chainstate"
	"github.com/bitcoin-sv/spv-wallet/engine/utils"
	"github.com/bitcoinschema/go-bitcoin/v2"
	"github.com/jarcoal/httpmock"
	"github.com/libsv/go-bk/bec"
	"github.com/libsv/go-bk/bip32"
	"github.com/libsv/go-bt/v2"
	"github.com/libsv/go-bt/v2/bscript"
	"github.com/libsv/go-bt/v2/sighash"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testDraftLockingScript = "76a9140692ed78f6988968ce612f620894997cc7edf1ad88ac"
)

// TestDraftTransaction_newDraftTransaction will test the method newDraftTransaction()
func TestDraftTransaction_newDraftTransaction(t *testing.T) {
	t.Parallel()

	t.Run("nil config, panic", func(t *testing.T) {
		assert.Panics(t, func() {
			draftTx, err := newDraftTransaction(
				testXPub, nil, New(),
			)
			require.NotNil(t, draftTx)
			require.NoError(t, err)
		})
	})

	t.Run("valid config", func(t *testing.T) {
		expires := time.Now().UTC().Add(defaultDraftTxExpiresIn)
		draftTx, err := newDraftTransaction(
			testXPub, &TransactionConfig{
				FeeUnit: chainstate.MockDefaultFee,
			}, New(),
		)
		require.NoError(t, err)
		require.NotNil(t, draftTx)
		assert.NotEqual(t, "", draftTx.ID)
		assert.Equal(t, 64, len(draftTx.ID))
		assert.WithinDurationf(t, expires, draftTx.ExpiresAt, 1*time.Second, "within 1 second")
		assert.Equal(t, DraftStatusDraft, draftTx.Status)
		assert.Equal(t, testXPubID, draftTx.XpubID)
		assert.Equal(t, *chainstate.MockDefaultFee, *draftTx.Configuration.FeeUnit)
	})
}

// TestDraftTransaction_GetModelName will test the method GetModelName()
func TestDraftTransaction_GetModelName(t *testing.T) {
	t.Parallel()

	t.Run("model name", func(t *testing.T) {
		draftTx, err := newDraftTransaction(testXPub, &TransactionConfig{
			FeeUnit: chainstate.MockDefaultFee,
		}, New())
		require.NoError(t, err)
		require.NotNil(t, draftTx)
		assert.Equal(t, ModelDraftTransaction.String(), draftTx.GetModelName())
	})
}

// TestDraftTransaction_getOutputSatoshis tests getting the output satoshis for the destinations
func TestDraftTransaction_getOutputSatoshis(t *testing.T) {
	t.Run("1 change destination", func(t *testing.T) {
		draftTx, err := newDraftTransaction(
			testXPub, &TransactionConfig{
				ChangeDestinations: []*Destination{{
					LockingScript: testLockingScript,
				}},
				FeeUnit: chainstate.MockDefaultFee,
			},
		)
		require.NoError(t, err)
		changSatoshis, err := draftTx.getChangeSatoshis(1000000)
		require.NoError(t, err)
		assert.Len(t, changSatoshis, 1)
		assert.Equal(t, uint64(1000000), changSatoshis[testLockingScript])
	})

	t.Run("1 change destination using default", func(t *testing.T) {
		draftTx, err := newDraftTransaction(
			testXPub, &TransactionConfig{
				ChangeDestinationsStrategy: ChangeStrategyDefault,
				ChangeDestinations: []*Destination{{
					LockingScript: testLockingScript,
				}},
				FeeUnit: chainstate.MockDefaultFee,
			},
		)
		require.NoError(t, err)
		changSatoshis, err := draftTx.getChangeSatoshis(1000000)
		require.NoError(t, err)
		assert.Len(t, changSatoshis, 1)
		assert.Equal(t, uint64(1000000), changSatoshis[testLockingScript])
	})

	t.Run("2 change destinations", func(t *testing.T) {
		draftTx, err := newDraftTransaction(
			testXPub, &TransactionConfig{
				ChangeDestinations: []*Destination{{
					LockingScript: testLockingScript,
				}, {
					LockingScript: testTxInScriptPubKey,
				}},
				FeeUnit: chainstate.MockDefaultFee,
			},
		)
		require.NoError(t, err)
		changSatoshis, err := draftTx.getChangeSatoshis(1000001)
		require.NoError(t, err)
		assert.Len(t, changSatoshis, 2)
		assert.Equal(t, uint64(500000), changSatoshis[testLockingScript])
		assert.Equal(t, uint64(500001), changSatoshis[testTxInScriptPubKey])
	})

	t.Run("3 change destinations - random", func(t *testing.T) {
		draftTx, err := newDraftTransaction(
			testXPub, &TransactionConfig{
				ChangeDestinationsStrategy: ChangeStrategyRandom,
				ChangeDestinations: []*Destination{{
					LockingScript: testLockingScript,
				}, {
					LockingScript: testTxInScriptPubKey,
				}, {
					LockingScript: testTxScriptPubKey1,
				}},
				FeeUnit: chainstate.MockDefaultFee,
			},
		)
		require.NoError(t, err)
		satoshis := uint64(1000001)
		changSatoshis, err := draftTx.getChangeSatoshis(satoshis)
		require.NoError(t, err)
		assert.Len(t, changSatoshis, 3)
		totalSatoshis := uint64(0)
		for _, s := range changSatoshis {
			totalSatoshis += s
		}
		assert.Equal(t, totalSatoshis, satoshis)
	})
}

// TestDraftTransaction_setChangeDestinations sets the given of change destinations on the draft transaction
func TestDraftTransaction_setChangeDestinations(t *testing.T) {
	t.Run("1 change destination", func(t *testing.T) {
		ctx, client, deferMe := CreateTestSQLiteClient(t, false, false, withTaskManagerMockup())
		defer deferMe()
		xPub := newXpub(testXPub, append(client.DefaultModelOptions(), New())...)
		err := xPub.Save(ctx)
		require.NoError(t, err)

		draftTx, err := newDraftTransaction(testXPub, &TransactionConfig{
			Outputs: []*TransactionOutput{{
				To:       testExternalAddress,
				Satoshis: 1000,
			}},
		}, append(client.DefaultModelOptions(), New())...)
		require.NoError(t, err)

		err = draftTx.setChangeDestinations(ctx, 1)
		require.NoError(t, err)
		assert.Len(t, draftTx.Configuration.ChangeDestinations, 1)
	})

	t.Run("5 change destinations", func(t *testing.T) {
		ctx, client, deferMe := CreateTestSQLiteClient(t, false, false, withTaskManagerMockup())
		defer deferMe()
		xPub := newXpub(testXPub, append(client.DefaultModelOptions(), New())...)
		err := xPub.Save(ctx)
		require.NoError(t, err)

		draftTx, err := newDraftTransaction(testXPub, &TransactionConfig{
			Outputs: []*TransactionOutput{{
				To:       testExternalAddress,
				Satoshis: 1000,
			}},
		}, append(client.DefaultModelOptions(), New())...)
		require.NoError(t, err)

		err = draftTx.setChangeDestinations(ctx, 5)
		require.NoError(t, err)
		assert.Len(t, draftTx.Configuration.ChangeDestinations, 5)
	})
}

// TestDraftTransaction_getDraftTransactionID tests getting the draft transaction by draft id
func TestDraftTransaction_getDraftTransactionID(t *testing.T) {
	t.Run("not found", func(t *testing.T) {
		ctx, client, deferMe := CreateTestSQLiteClient(t, false, false, withTaskManagerMockup())
		defer deferMe()
		draftTx, err := getDraftTransactionID(ctx, testXPubID, testDraftID, client.DefaultModelOptions()...)
		require.NoError(t, err)
		assert.Nil(t, draftTx)
	})

	t.Run("found by draft id", func(t *testing.T) {
		ctx, client, deferMe := CreateTestSQLiteClient(t, false, false, withTaskManagerMockup())
		defer deferMe()
		draftTransaction, err := newDraftTransaction(testXPub, &TransactionConfig{}, client.DefaultModelOptions()...)
		require.NoError(t, err)
		err = draftTransaction.Save(ctx)
		require.NoError(t, err)

		var draftTx *DraftTransaction
		draftTx, err = getDraftTransactionID(ctx, testXPubID, draftTransaction.ID, client.DefaultModelOptions()...)
		require.NoError(t, err)
		assert.Equal(t, 64, len(draftTx.GetID()))
		assert.Equal(t, testXPubID, draftTx.XpubID)
	})
}

// TestDraftTransaction_createTransaction create a transaction hex
func TestDraftTransaction_createTransaction(t *testing.T) {
	t.Run("empty transaction", func(t *testing.T) {
		ctx, client, deferMe := CreateTestSQLiteClient(t, false, false, withTaskManagerMockup())
		defer deferMe()
		draftTransaction, err := newDraftTransaction(testXPub, &TransactionConfig{}, append(client.DefaultModelOptions(), New())...)
		require.NoError(t, err)

		err = draftTransaction.createTransactionHex(ctx)
		require.ErrorIs(t, err, ErrMissingTransactionOutputs)
	})

	t.Run("transaction with no utxos", func(t *testing.T) {
		ctx, client, deferMe := CreateTestSQLiteClient(t, false, false, withTaskManagerMockup())
		defer deferMe()
		draftTransaction, err := newDraftTransaction(testXPub, &TransactionConfig{
			Outputs: []*TransactionOutput{{
				To:       testExternalAddress,
				Satoshis: 1000,
			}},
		}, append(client.DefaultModelOptions(), New())...)
		require.NoError(t, err)

		err = draftTransaction.createTransactionHex(ctx)
		require.ErrorIs(t, err, ErrNotEnoughUtxos)
	})

	t.Run("transaction with utxos", func(t *testing.T) {
		ctx, client, deferMe := CreateTestSQLiteClient(t, false, true, withTaskManagerMockup())
		defer deferMe()
		xPub := newXpub(testXPub, append(client.DefaultModelOptions(), New())...)
		err := xPub.Save(ctx)
		require.NoError(t, err)

		destination := newDestination(testXPubID, testLockingScript,
			append(client.DefaultModelOptions(), New())...)
		err = destination.Save(ctx)
		require.NoError(t, err)

		utxo := newUtxo(testXPubID, testTxID, testLockingScript, 0, 100000,
			append(client.DefaultModelOptions(), New())...)
		err = utxo.Save(ctx)
		require.NoError(t, err)

		transaction, err := txFromHex(testTxHex, append(client.DefaultModelOptions(), New())...)
		require.NoError(t, err)

		err = transaction.Save(ctx)
		require.NoError(t, err)

		draftTransaction, err := newDraftTransaction(testXPub, &TransactionConfig{
			Outputs: []*TransactionOutput{{
				To:       testExternalAddress,
				Satoshis: 1000,
			}},
		}, append(client.DefaultModelOptions(), New())...)
		require.NoError(t, err)

		err = draftTransaction.createTransactionHex(ctx)
		require.NoError(t, err)
		assert.Equal(t, testXPubID, draftTransaction.XpubID)
		assert.Equal(t, DraftStatusDraft, draftTransaction.Status)

		assert.Equal(t, testXPubID, draftTransaction.Configuration.ChangeDestinations[0].XpubID)
		assert.Equal(t, draftTransaction.ID, draftTransaction.Configuration.ChangeDestinations[0].DraftID)
		assert.Equal(t, uint64(98988), draftTransaction.Configuration.ChangeSatoshis)

		assert.Equal(t, uint64(12), draftTransaction.Configuration.Fee)
		assert.Equal(t, *chainstate.MockDefaultFee, *draftTransaction.Configuration.FeeUnit)

		assert.Equal(t, 1, len(draftTransaction.Configuration.Inputs))
		assert.Equal(t, testLockingScript, draftTransaction.Configuration.Inputs[0].ScriptPubKey)
		assert.Equal(t, uint64(100000), draftTransaction.Configuration.Inputs[0].Satoshis)

		assert.Equal(t, 2, len(draftTransaction.Configuration.Outputs))
		assert.Equal(t, uint64(1000), draftTransaction.Configuration.Outputs[0].Satoshis)
		assert.Equal(t, uint64(98988), draftTransaction.Configuration.Outputs[1].Satoshis)
		assert.Equal(t, draftTransaction.Configuration.ChangeDestinations[0].LockingScript, draftTransaction.Configuration.Outputs[1].Scripts[0].Script)

		var btTx *bt.Tx
		btTx, err = bt.NewTxFromString(draftTransaction.Hex)
		require.NoError(t, err)

		assert.Equal(t, 1, len(btTx.Inputs))
		assert.Equal(t, testTxID, hex.EncodeToString(btTx.Inputs[0].PreviousTxID()))
		assert.Equal(t, uint32(0), btTx.Inputs[0].PreviousTxOutIndex)

		assert.Equal(t, 2, len(btTx.Outputs))
		assert.Equal(t, uint64(1000), btTx.Outputs[0].Satoshis)
		assert.Equal(t, draftTransaction.Configuration.Outputs[0].Scripts[0].Script, btTx.Outputs[0].LockingScript.String())

		assert.Equal(t, uint64(98988), btTx.Outputs[1].Satoshis)
		assert.Equal(t, draftTransaction.Configuration.Outputs[1].Scripts[0].Script, btTx.Outputs[1].LockingScript.String())

		var gUtxo *Utxo
		gUtxo, err = getUtxo(ctx, testTxID, 0, client.DefaultModelOptions()...)
		require.NoError(t, err)
		assert.True(t, gUtxo.DraftID.Valid)
		assert.Equal(t, draftTransaction.ID, gUtxo.DraftID.String)
		assert.True(t, gUtxo.ReservedAt.Valid)
	})

	t.Run("send to all", func(t *testing.T) {
		ctx, client, deferMe := CreateTestSQLiteClient(t, false, true, withTaskManagerMockup())
		defer deferMe()
		xPub := newXpub(testXPub, append(client.DefaultModelOptions(), New())...)
		err := xPub.Save(ctx)
		require.NoError(t, err)

		destination := newDestination(testXPubID, testLockingScript,
			append(client.DefaultModelOptions(), New())...)
		err = destination.Save(ctx)
		require.NoError(t, err)

		utxo := newUtxo(testXPubID, testTxID, testLockingScript, 0, 100000,
			append(client.DefaultModelOptions(), New())...)
		err = utxo.Save(ctx)
		require.NoError(t, err)

		transaction, err := txFromHex(testTxHex, append(client.DefaultModelOptions(), New())...)
		require.NoError(t, err)

		err = transaction.Save(ctx)
		require.NoError(t, err)

		draftTransaction, err := newDraftTransaction(testXPub, &TransactionConfig{
			SendAllTo: &TransactionOutput{To: testExternalAddress},
		}, append(client.DefaultModelOptions(), New())...)
		require.NoError(t, err)

		err = draftTransaction.createTransactionHex(ctx)
		require.NoError(t, err)
		assert.Equal(t, testXPubID, draftTransaction.XpubID)
		assert.Equal(t, DraftStatusDraft, draftTransaction.Status)
		assert.Equal(t, testExternalAddress, draftTransaction.Configuration.SendAllTo.To)
		assert.Equal(t, uint64(10), draftTransaction.Configuration.Fee)
		assert.Len(t, draftTransaction.Configuration.Inputs, 1)
		assert.Len(t, draftTransaction.Configuration.Outputs, 1)
		assert.Equal(t, testExternalAddress, draftTransaction.Configuration.Outputs[0].To)
		assert.Equal(t, uint64(99990), draftTransaction.Configuration.Outputs[0].Satoshis)
		assert.Len(t, draftTransaction.Configuration.Outputs[0].Scripts, 1)
		assert.Equal(t, testExternalAddress, draftTransaction.Configuration.Outputs[0].Scripts[0].Address)
		assert.Equal(t, uint64(99990), draftTransaction.Configuration.Outputs[0].Scripts[0].Satoshis)
	})

	t.Run("fee calculation - MAP", func(t *testing.T) {
		ctx, client, deferMe := CreateTestSQLiteClient(t, false, true, withTaskManagerMockup())
		defer deferMe()
		xPub := newXpub(testXPub, append(client.DefaultModelOptions(), New())...)
		err := xPub.Save(ctx)
		require.NoError(t, err)

		destination := newDestination(testXPubID, testLockingScript,
			append(client.DefaultModelOptions(), New())...)
		err = destination.Save(ctx)
		require.NoError(t, err)

		utxo := newUtxo(testXPubID, testTxID, testLockingScript, 0, 100000,
			append(client.DefaultModelOptions(), New())...)
		err = utxo.Save(ctx)
		require.NoError(t, err)

		transaction, err := txFromHex(testTxHex, append(client.DefaultModelOptions(), New())...)
		require.NoError(t, err)

		err = transaction.Save(ctx)
		require.NoError(t, err)

		draftTransaction, err := newDraftTransaction(testXPub, &TransactionConfig{
			Outputs: []*TransactionOutput{{
				To:       testExternalAddress,
				Satoshis: 1000,
			}, {
				OpReturn: &OpReturn{
					Map: &MapProtocol{
						App: "tonicpow_staging",
						Keys: map[string]interface{}{
							"offer_config_id":  "336",
							"offer_session_id": "4f06c11358e6586e67c77467c252a8be9187211f704de2627e4824945f31f07e",
						},
						Type: "offer_clicks",
					},
				},
			}},
		}, append(client.DefaultModelOptions(), New())...)
		require.NoError(t, err)

		err = draftTransaction.createTransactionHex(ctx)
		require.NoError(t, err)
		fee := draftTransaction.Configuration.Fee
		calculateFee := draftTransaction.estimateFee(draftTransaction.Configuration.FeeUnit, 0)
		assert.Equal(t, fee, calculateFee)
		txFee := uint64(0)
		for _, input := range draftTransaction.Configuration.Inputs {
			txFee += input.Satoshis
		}
		for _, output := range draftTransaction.Configuration.Outputs {
			txFee -= output.Satoshis
		}
		assert.Equal(t, fee, txFee)
	})

	t.Run("fee calculation - MAP 2", func(t *testing.T) {
		ctx, client, deferMe := CreateTestSQLiteClient(t, false, true, withTaskManagerMockup())
		defer deferMe()
		xPub := newXpub(testXPub, append(client.DefaultModelOptions(), New())...)
		err := xPub.Save(ctx)
		require.NoError(t, err)

		destination := newDestination(testXPubID, testLockingScript,
			append(client.DefaultModelOptions(), New())...)
		err = destination.Save(ctx)
		require.NoError(t, err)

		utxo := newUtxo(testXPubID, testTxID, testLockingScript, 0, 11239,
			append(client.DefaultModelOptions(), New())...)
		err = utxo.Save(ctx)
		require.NoError(t, err)

		utxo = newUtxo(testXPubID, testTxID, testLockingScript, 1, 251725,
			append(client.DefaultModelOptions(), New())...)
		err = utxo.Save(ctx)
		require.NoError(t, err)

		transaction, err := txFromHex(testTxHex, append(client.DefaultModelOptions(), New())...)
		require.NoError(t, err)

		err = transaction.Save(ctx)
		require.NoError(t, err)

		draftTransaction, err := newDraftTransaction(testXPub, &TransactionConfig{
			Outputs: []*TransactionOutput{{
				To:       testExternalAddress,
				Satoshis: 12166,
			}, {
				To:       testExternalAddress,
				Satoshis: 1217,
			}, {
				OpReturn: &OpReturn{
					Map: &MapProtocol{
						App: "tonicpow_staging",
						Keys: map[string]interface{}{
							"offer_conversion_config_id": "79e650cf5e76938f1e96ea0086f1707fe2f28da39f2460bb2626b738d0958fe4",
							"offer_session_id":           "4f06c11358e6586e67c77467c252a8be9187211f704de2627e4824945f31f07e",
						},
						Type: "offer_conversion",
					},
				},
			}},
		}, append(client.DefaultModelOptions(), New())...)
		require.NoError(t, err)

		err = draftTransaction.createTransactionHex(ctx)
		require.NoError(t, err)
		fee := draftTransaction.Configuration.Fee
		calculateFee := draftTransaction.estimateFee(draftTransaction.Configuration.FeeUnit, 0)
		assert.Equal(t, fee, calculateFee)
		txFee := uint64(0)
		for _, input := range draftTransaction.Configuration.Inputs {
			txFee += input.Satoshis
		}
		for _, output := range draftTransaction.Configuration.Outputs {
			txFee -= output.Satoshis
		}
		assert.Equal(t, fee, txFee)
	})

	t.Run("fee calculation - tonicpow", func(t *testing.T) {
		ctx, client, deferMe := CreateTestSQLiteClient(t, false, true, withTaskManagerMockup())
		defer deferMe()
		xPub := newXpub(testXPub, append(client.DefaultModelOptions(), New())...)
		err := xPub.Save(ctx)
		require.NoError(t, err)

		destination := newDestination(testXPubID, testLockingScript,
			append(client.DefaultModelOptions(), New())...)
		err = destination.Save(ctx)
		require.NoError(t, err)

		utxo := newUtxo(testXPubID, testTxID, testLockingScript, 1, 114216,
			append(client.DefaultModelOptions(), New())...)
		err = utxo.Save(ctx)
		require.NoError(t, err)

		transaction, err := txFromHex(testTxHex, append(client.DefaultModelOptions(), New())...)
		require.NoError(t, err)

		err = transaction.Save(ctx)
		require.NoError(t, err)

		draftTransaction, err := newDraftTransaction(testXPub, &TransactionConfig{
			FeeUnit: &utils.FeeUnit{
				Satoshis: 5,
				Bytes:    100,
			},
			Outputs: []*TransactionOutput{{
				OpReturn: &OpReturn{
					Map: &MapProtocol{
						App:  "tonicpow_staging",
						Type: "offer_conversion",
						Keys: map[string]interface{}{
							"offer_conversion_config_id": "05384d8d8678e560426e1fb81e6723b6",
							"offer_session_id":           "74a66f09450bd0bcccd5a26714cbebdb20d6ba7860d668562182f4c2512460a4",
						},
					},
				},
			}},
		}, append(client.DefaultModelOptions(), New())...)
		require.NoError(t, err)

		err = draftTransaction.createTransactionHex(ctx)
		require.NoError(t, err)
		fee := draftTransaction.Configuration.Fee
		calculateFee := draftTransaction.estimateFee(draftTransaction.Configuration.FeeUnit, 0)
		assert.Equal(t, fee, calculateFee)
		txFee := uint64(0)
		for _, input := range draftTransaction.Configuration.Inputs {
			txFee += input.Satoshis
		}
		for _, output := range draftTransaction.Configuration.Outputs {
			txFee -= output.Satoshis
		}
		assert.Equal(t, fee, txFee)
	})

	t.Run("send to all - multiple utxos", func(t *testing.T) {
		ctx, client, deferMe := CreateTestSQLiteClient(t, false, true, withTaskManagerMockup())
		defer deferMe()
		xPub := newXpub(testXPub, append(client.DefaultModelOptions(), New())...)
		err := xPub.Save(ctx)
		require.NoError(t, err)

		destination := newDestination(testXPubID, testLockingScript,
			append(client.DefaultModelOptions(), New())...)
		err = destination.Save(ctx)
		require.NoError(t, err)

		utxo := newUtxo(testXPubID, testTxID, testLockingScript, 0, 100000,
			append(client.DefaultModelOptions(), New())...)
		err = utxo.Save(ctx)
		require.NoError(t, err)
		utxo = newUtxo(testXPubID, testTxID, testLockingScript, 1, 110000,
			append(client.DefaultModelOptions(), New())...)
		err = utxo.Save(ctx)
		require.NoError(t, err)
		utxo = newUtxo(testXPubID, testTxID, testLockingScript, 2, 130000,
			append(client.DefaultModelOptions(), New())...)
		err = utxo.Save(ctx)
		require.NoError(t, err)

		transaction, err := txFromHex(testTxHex, append(client.DefaultModelOptions(), New())...)
		require.NoError(t, err)

		err = transaction.Save(ctx)
		require.NoError(t, err)

		draftTransaction, err := newDraftTransaction(testXPub, &TransactionConfig{
			SendAllTo: &TransactionOutput{To: testExternalAddress},
		}, append(client.DefaultModelOptions(), New())...)
		require.NoError(t, err)

		err = draftTransaction.createTransactionHex(ctx)
		require.NoError(t, err)
		assert.Equal(t, testXPubID, draftTransaction.XpubID)
		assert.Equal(t, DraftStatusDraft, draftTransaction.Status)
		assert.Equal(t, testExternalAddress, draftTransaction.Configuration.SendAllTo.To)
		assert.Equal(t, uint64(25), draftTransaction.Configuration.Fee)
		assert.Len(t, draftTransaction.Configuration.Inputs, 3)
		assert.Len(t, draftTransaction.Configuration.Outputs, 1)
		assert.Equal(t, testExternalAddress, draftTransaction.Configuration.Outputs[0].To)
		assert.Equal(t, uint64(339975), draftTransaction.Configuration.Outputs[0].Satoshis)
		assert.Len(t, draftTransaction.Configuration.Outputs[0].Scripts, 1)
		assert.Equal(t, testExternalAddress, draftTransaction.Configuration.Outputs[0].Scripts[0].Address)
		assert.Equal(t, uint64(339975), draftTransaction.Configuration.Outputs[0].Scripts[0].Satoshis)
	})

	t.Run("send to all - selected utxos", func(t *testing.T) {
		ctx, client, deferMe := CreateTestSQLiteClient(t, false, true, withTaskManagerMockup())
		defer deferMe()
		xPub := newXpub(testXPub, append(client.DefaultModelOptions(), New())...)
		err := xPub.Save(ctx)
		require.NoError(t, err)

		destination := newDestination(testXPubID, testLockingScript,
			append(client.DefaultModelOptions(), New())...)
		err = destination.Save(ctx)
		require.NoError(t, err)

		utxo := newUtxo(testXPubID, testTxID, testLockingScript, 0, 100000,
			append(client.DefaultModelOptions(), New())...)
		err = utxo.Save(ctx)
		require.NoError(t, err)
		utxo = newUtxo(testXPubID, testTxID, testLockingScript, 1, 110000,
			append(client.DefaultModelOptions(), New())...)
		err = utxo.Save(ctx)
		require.NoError(t, err)
		utxo = newUtxo(testXPubID, testTxID, testLockingScript, 2, 130000,
			append(client.DefaultModelOptions(), New())...)
		err = utxo.Save(ctx)
		require.NoError(t, err)

		transaction, err := txFromHex(testTxHex, append(client.DefaultModelOptions(), New())...)
		require.NoError(t, err)

		err = transaction.Save(ctx)
		require.NoError(t, err)

		draftTransaction, err := newDraftTransaction(testXPub, &TransactionConfig{
			SendAllTo: &TransactionOutput{To: testExternalAddress},
			FromUtxos: []*UtxoPointer{{
				TransactionID: testTxID,
				OutputIndex:   1,
			}, {
				TransactionID: testTxID,
				OutputIndex:   2,
			}},
		}, append(client.DefaultModelOptions(), New())...)
		require.NoError(t, err)

		err = draftTransaction.createTransactionHex(ctx)
		require.NoError(t, err)
		assert.Equal(t, testXPubID, draftTransaction.XpubID)
		assert.Equal(t, DraftStatusDraft, draftTransaction.Status)
		assert.Equal(t, testExternalAddress, draftTransaction.Configuration.SendAllTo.To)
		assert.Equal(t, uint64(17), draftTransaction.Configuration.Fee)
		assert.Len(t, draftTransaction.Configuration.Inputs, 2)
		assert.Len(t, draftTransaction.Configuration.Outputs, 1)
		assert.Equal(t, testExternalAddress, draftTransaction.Configuration.Outputs[0].To)
		assert.Equal(t, uint64(239983), draftTransaction.Configuration.Outputs[0].Satoshis)
		assert.Len(t, draftTransaction.Configuration.Outputs[0].Scripts, 1)
		assert.Equal(t, testExternalAddress, draftTransaction.Configuration.Outputs[0].Scripts[0].Address)
		assert.Equal(t, uint64(239983), draftTransaction.Configuration.Outputs[0].Scripts[0].Satoshis)
	})

	t.Run("include utxos - tokens", func(t *testing.T) {
		ctx, client, deferMe := CreateTestSQLiteClient(t, false, true, withTaskManagerMockup())
		defer deferMe()
		xPub := newXpub(testXPub, append(client.DefaultModelOptions(), New())...)
		err := xPub.Save(ctx)
		require.NoError(t, err)

		destination := newDestination(testXPubID, testLockingScript,
			append(client.DefaultModelOptions(), New())...)
		err = destination.Save(ctx)
		require.NoError(t, err)

		destination = newDestination(testXPubID, testSTASScriptPubKey,
			append(client.DefaultModelOptions(), New())...)
		err = destination.Save(ctx)
		require.NoError(t, err)

		utxo := newUtxo(testXPubID, testTxID, testLockingScript, 0, 100000,
			append(client.DefaultModelOptions(), New())...)
		err = utxo.Save(ctx)
		require.NoError(t, err)
		utxo = newUtxo(testXPubID, testTxID, testLockingScript, 1, 110000,
			append(client.DefaultModelOptions(), New())...)
		err = utxo.Save(ctx)
		require.NoError(t, err)
		utxo = newUtxo(testXPubID, testTxID, testSTASLockingScript, 2, 564,
			append(client.DefaultModelOptions(), New())...)
		err = utxo.Save(ctx)
		require.NoError(t, err)

		transaction, err := txFromHex(testTxHex, append(client.DefaultModelOptions(), New())...)
		require.NoError(t, err)

		err = transaction.Save(ctx)
		require.NoError(t, err)

		// todo how to make sure we do not unwittingly destroy tokens ?
		draftTransaction, err := newDraftTransaction(testXPub, &TransactionConfig{
			IncludeUtxos: []*UtxoPointer{{
				TransactionID: testTxID,
				OutputIndex:   2,
			}},
			Outputs: []*TransactionOutput{{
				To:       testExternalAddress,
				Satoshis: 1000,
			}, {
				Script:   testSTASLockingScript, // send token to the same destination
				Satoshis: 564,
			}},
		}, append(client.DefaultModelOptions(), New())...)
		require.NoError(t, err)

		err = draftTransaction.createTransactionHex(ctx)
		require.NoError(t, err)
		assert.Equal(t, testXPubID, draftTransaction.XpubID)
		assert.Equal(t, DraftStatusDraft, draftTransaction.Status)
		assert.Equal(t, uint64(120), draftTransaction.Configuration.Fee)
		assert.Len(t, draftTransaction.Configuration.Inputs, 2)
		assert.Len(t, draftTransaction.Configuration.Outputs, 3)

		assert.Equal(t, testSTASLockingScript, draftTransaction.Configuration.Inputs[0].ScriptPubKey)
		assert.Equal(t, uint64(564), draftTransaction.Configuration.Inputs[0].Satoshis)

		assert.Equal(t, testExternalAddress, draftTransaction.Configuration.Outputs[0].To)
		assert.Equal(t, uint64(1000), draftTransaction.Configuration.Outputs[0].Satoshis)
		assert.Len(t, draftTransaction.Configuration.Outputs[0].Scripts, 1)
		assert.Equal(t, testExternalAddress, draftTransaction.Configuration.Outputs[0].Scripts[0].Address)
		assert.Equal(t, uint64(1000), draftTransaction.Configuration.Outputs[0].Scripts[0].Satoshis)

		assert.Equal(t, "", draftTransaction.Configuration.Outputs[1].To)
		assert.Equal(t, uint64(564), draftTransaction.Configuration.Outputs[1].Satoshis)
		assert.Equal(t, testSTASLockingScript, draftTransaction.Configuration.Outputs[1].Script)
		assert.Len(t, draftTransaction.Configuration.Outputs[1].Scripts, 1)
		assert.Equal(t, "", draftTransaction.Configuration.Outputs[1].Scripts[0].Address)
		assert.Equal(t, uint64(564), draftTransaction.Configuration.Outputs[1].Scripts[0].Satoshis)
		assert.Equal(t, testSTASLockingScript, draftTransaction.Configuration.Outputs[1].Scripts[0].Script)

		assert.Equal(t, uint64(98880), draftTransaction.Configuration.Outputs[2].Satoshis)
	})

	t.Run("SendAllTo", func(t *testing.T) {
		ctx, client, deferMe := initSimpleTestCase(t)
		defer deferMe()

		draftTransaction, err := newDraftTransaction(testXPub, &TransactionConfig{
			FromUtxos: []*UtxoPointer{{
				TransactionID: testTxID,
				OutputIndex:   0,
			}},
			SendAllTo: &TransactionOutput{
				To: testExternalAddress,
			},
		}, append(client.DefaultModelOptions(), New())...)
		require.NoError(t, err)

		err = draftTransaction.createTransactionHex(ctx)
		require.NoError(t, err)
		assert.Len(t, draftTransaction.Configuration.Outputs, 1)
		assert.Equal(t, testExternalAddress, draftTransaction.Configuration.Outputs[0].To)
		assert.Equal(t, uint64(99990), draftTransaction.Configuration.Outputs[0].Satoshis)
		assert.Equal(t, uint64(10), draftTransaction.Configuration.Fee)
	})

	t.Run("SendAllTo + output", func(t *testing.T) {
		ctx, client, deferMe := initSimpleTestCase(t)
		defer deferMe()

		draftTransaction, err := newDraftTransaction(testXPub, &TransactionConfig{
			FromUtxos: []*UtxoPointer{{
				TransactionID: testTxID,
				OutputIndex:   0,
			}},
			SendAllTo: &TransactionOutput{
				To: testExternalAddress,
			},
			Outputs: []*TransactionOutput{{
				To:       testExternalAddress,
				Satoshis: 1000,
			}},
		}, append(client.DefaultModelOptions(), New())...)
		require.NoError(t, err)

		err = draftTransaction.createTransactionHex(ctx)
		require.NoError(t, err)
		assert.Len(t, draftTransaction.Configuration.Outputs, 2)
		assert.Equal(t, testExternalAddress, draftTransaction.Configuration.Outputs[0].To)
		assert.Equal(t, uint64(98988), draftTransaction.Configuration.Outputs[0].Satoshis)
		assert.Equal(t, uint64(12), draftTransaction.Configuration.Fee)
		assert.Equal(t, testExternalAddress, draftTransaction.Configuration.Outputs[1].To)
		assert.Equal(t, uint64(1000), draftTransaction.Configuration.Outputs[1].Satoshis)
	})

	t.Run("SendAllTo + output + op_return", func(t *testing.T) {
		ctx, client, deferMe := initSimpleTestCase(t)
		defer deferMe()

		draftTransaction, err := newDraftTransaction(testXPub, &TransactionConfig{
			FromUtxos: []*UtxoPointer{{
				TransactionID: testTxID,
				OutputIndex:   0,
			}},
			SendAllTo: &TransactionOutput{
				To: testExternalAddress,
			},
			Outputs: []*TransactionOutput{{
				To:       testExternalAddress,
				Satoshis: 1000,
			}, {
				OpReturn: &OpReturn{
					Map: &MapProtocol{
						App:  "social",
						Type: "post",
						Keys: map[string]interface{}{
							"title": "Hello World!",
						},
					},
				},
			}},
		}, append(client.DefaultModelOptions(), New())...)
		require.NoError(t, err)

		err = draftTransaction.createTransactionHex(ctx)
		require.NoError(t, err)
		assert.Len(t, draftTransaction.Configuration.Outputs, 3)
		assert.Equal(t, testExternalAddress, draftTransaction.Configuration.Outputs[0].To)
		assert.Equal(t, uint64(98984), draftTransaction.Configuration.Outputs[0].Satoshis)
		assert.Equal(t, uint64(16), draftTransaction.Configuration.Fee)
		assert.Equal(t, testExternalAddress, draftTransaction.Configuration.Outputs[1].To)
		assert.Equal(t, uint64(1000), draftTransaction.Configuration.Outputs[1].Satoshis)
		assert.Equal(t, "", draftTransaction.Configuration.Outputs[2].To)
		assert.Equal(t, uint64(0), draftTransaction.Configuration.Outputs[2].Satoshis)
	})

	t.Run("SendAllTo + 2 utxos", func(t *testing.T) {
		p := newTestPaymailClient(t, []string{"handcash.io"})
		ctx, client, deferMe := CreateTestSQLiteClient(t, false, true,
			withTaskManagerMockup(),
			WithPaymailClient(p),
		)
		defer deferMe()

		httpmock.Reset()
		mockValidResponse(http.StatusOK, true, "handcash.io")
		httpmock.RegisterResponder(http.MethodPost, "https://handcash.io/api/v1/bsvalias/p2p-payment-destination/mrzz@handcash.io",
			httpmock.NewStringResponder(
				200,
				`{"outputs": [{"script": "76a9143e2d1d795f8acaa7957045cc59376177eb04a3c588ac","satoshis": 100}],"reference": "z0bac4ec-6f15-42de-9ef4-e60bfdabf4f7"}`,
			),
		)

		xPub := newXpub(testXPub, append(client.DefaultModelOptions(), New())...)
		err := xPub.Save(ctx)
		require.NoError(t, err)

		destination := newDestination(testXPubID, testLockingScript,
			append(client.DefaultModelOptions(), New())...)
		err = destination.Save(ctx)
		require.NoError(t, err)

		utxo := newUtxo(testXPubID, testTxID, testLockingScript, 1, 13483,
			append(client.DefaultModelOptions(), New())...)
		err = utxo.Save(ctx)
		require.NoError(t, err)

		utxo = newUtxo(testXPubID, testTxID, testLockingScript, 0, 13483,
			append(client.DefaultModelOptions(), New())...)
		err = utxo.Save(ctx)
		require.NoError(t, err)

		transaction, err := txFromHex(testTxHex, append(client.DefaultModelOptions(), New())...)
		require.NoError(t, err)

		err = transaction.Save(ctx)
		require.NoError(t, err)

		draftTransaction, err := newDraftTransaction(testXPub, &TransactionConfig{
			FromUtxos: []*UtxoPointer{{
				TransactionID: testTxID,
				OutputIndex:   1,
			}, {
				TransactionID: testTxID,
				OutputIndex:   0,
			}},
			SendAllTo: &TransactionOutput{
				To: "mrzz@handcash.io",
			},
			Outputs: []*TransactionOutput{{
				OpReturn: &OpReturn{
					Map: &MapProtocol{
						App:  "social",
						Type: "post",
						Keys: map[string]interface{}{
							"title": "Hello World!",
						},
					},
				},
			}},
		}, append(client.DefaultModelOptions(), New())...)
		require.NoError(t, err)

		err = draftTransaction.createTransactionHex(ctx)
		require.NoError(t, err)
		assert.Len(t, draftTransaction.Configuration.Outputs, 2)
		assert.Equal(t, "mrzz@handcash.io", draftTransaction.Configuration.Outputs[0].To)
		assert.Equal(t, uint64(26944), draftTransaction.Configuration.Outputs[0].Satoshis)
		assert.Equal(t, "16fkwYn8feXEbK7iCTg5KMx9Rx9GzZ9HuE", draftTransaction.Configuration.Outputs[0].Scripts[0].Address)
		assert.Equal(t, "76a9143e2d1d795f8acaa7957045cc59376177eb04a3c588ac", draftTransaction.Configuration.Outputs[0].Scripts[0].Script)
		assert.Equal(t, uint64(26944), draftTransaction.Configuration.Outputs[0].Scripts[0].Satoshis)
		assert.Equal(t, uint64(22), draftTransaction.Configuration.Fee)
		assert.Equal(t, "", draftTransaction.Configuration.Outputs[1].To)
		assert.Equal(t, uint64(0), draftTransaction.Configuration.Outputs[1].Satoshis)
	})

	t.Run("duplicate inputs", func(t *testing.T) {
		ctx, client, deferMe := initSimpleTestCase(t)
		defer deferMe()

		opts := append(client.DefaultModelOptions(), New())
		utxo := newUtxo(testXPubID, testTxID, testLockingScript, 12, 1225, opts...)
		err := utxo.Save(ctx)
		require.NoError(t, err)

		draftTransaction, err := newDraftTransaction(testXPub, &TransactionConfig{
			FromUtxos: []*UtxoPointer{{
				TransactionID: utxo.TransactionID,
				OutputIndex:   utxo.OutputIndex,
			}, {
				TransactionID: utxo.TransactionID,
				OutputIndex:   utxo.OutputIndex,
			}},
			Outputs: []*TransactionOutput{{
				To:       testExternalAddress,
				Satoshis: 1500,
			}},
		}, append(client.DefaultModelOptions(), New())...)
		require.NoError(t, err)

		err = draftTransaction.createTransactionHex(ctx)
		require.ErrorIs(t, err, ErrDuplicateUTXOs)
	})
}

// TestDraftTransaction_setChangeDestination setting the change destination
func TestDraftTransaction_setChangeDestination(t *testing.T) {
	t.Run("missing xpub", func(t *testing.T) {
		ctx, client, deferMe := CreateTestSQLiteClient(t, false, false, withTaskManagerMockup())
		defer deferMe()
		draftTransaction := &DraftTransaction{
			Model: *NewBaseModel(
				ModelDraftTransaction,
				append(client.DefaultModelOptions(), WithXPub(testXPub))...,
			),
			Configuration: TransactionConfig{
				ChangeDestinations: nil,
				FeeUnit: &utils.FeeUnit{
					Satoshis: 5,
					Bytes:    10,
				},
			},
		}

		_, err := draftTransaction.setChangeDestination(ctx, 100, 200)
		require.ErrorIs(t, err, ErrMissingXpub)
	})

	t.Run("set valid destination", func(t *testing.T) {
		ctx, client, deferMe := CreateTestSQLiteClient(t, false, false, withTaskManagerMockup())
		defer deferMe()
		xPub := newXpub(testXPub, append(client.DefaultModelOptions(), New())...)
		xPub.NextExternalNum = 121
		xPub.NextInternalNum = 12
		err := xPub.Save(ctx)
		require.NoError(t, err)

		draftTransaction := &DraftTransaction{
			Model: *NewBaseModel(
				ModelDraftTransaction,
				append(client.DefaultModelOptions(), WithXPub(testXPub))...,
			),
			Configuration: TransactionConfig{
				ChangeDestinations: nil,
				FeeUnit: &utils.FeeUnit{
					Satoshis: 5,
					Bytes:    10,
				},
			},
		}

		var newFee uint64
		newFee, err = draftTransaction.setChangeDestination(ctx, 100, 0)
		require.NoError(t, err)
		assert.Equal(t, uint64(23), newFee)
		assert.Equal(t, uint64(77), draftTransaction.Configuration.ChangeSatoshis)
		assert.Equal(t, testXPubID, draftTransaction.Configuration.ChangeDestinations[0].XpubID)
		assert.Equal(t, uint32(1), draftTransaction.Configuration.ChangeDestinations[0].Chain)
		assert.Equal(t, uint32(12), draftTransaction.Configuration.ChangeDestinations[0].Num)
		assert.Equal(t, utils.ScriptTypePubKeyHash, draftTransaction.Configuration.ChangeDestinations[0].Type)
		assert.Equal(t, uint64(77), draftTransaction.Configuration.Outputs[0].Satoshis)
	})

	t.Run("use existing output", func(t *testing.T) {
		ctx, client, deferMe := CreateTestSQLiteClient(t, false, false, withTaskManagerMockup())
		defer deferMe()
		xPub := newXpub(testXPub, append(client.DefaultModelOptions(), New())...)
		err := xPub.Save(ctx)
		require.NoError(t, err)

		draftTransaction := &DraftTransaction{
			Model: *NewBaseModel(
				ModelDraftTransaction,
				append(client.DefaultModelOptions(), WithXPub(testXPub))...,
			),
			Configuration: TransactionConfig{
				Outputs: []*TransactionOutput{{
					To:           testExternalAddress,
					Satoshis:     1000,
					UseForChange: true,
				}},
				ChangeDestinations: []*Destination{{
					ID: testDestinationID,
				}},
			},
		}

		var newFee uint64
		newFee, err = draftTransaction.setChangeDestination(ctx, 100, 0)
		require.NoError(t, err)
		assert.Equal(t, uint64(0), newFee)
		assert.Equal(t, uint64(100), draftTransaction.Configuration.ChangeSatoshis)
		assert.Nil(t, draftTransaction.Configuration.ChangeDestinations)
		// 100 sats added to the output
		assert.Equal(t, uint64(1100), draftTransaction.Configuration.Outputs[0].Satoshis)
	})

	t.Run("use existing outputs", func(t *testing.T) {
		ctx, client, deferMe := CreateTestSQLiteClient(t, false, false, withTaskManagerMockup())
		defer deferMe()
		xPub := newXpub(testXPub, append(client.DefaultModelOptions(), New())...)
		err := xPub.Save(ctx)
		require.NoError(t, err)

		draftTransaction := &DraftTransaction{
			Model: *NewBaseModel(
				ModelDraftTransaction,
				append(client.DefaultModelOptions(), WithXPub(testXPub))...,
			),
			Configuration: TransactionConfig{
				Outputs: []*TransactionOutput{{
					To:           testExternalAddress,
					Satoshis:     1000,
					UseForChange: true,
				}, {
					To:       testExternalAddress,
					Satoshis: 1000,
				}, {
					To:           testExternalAddress,
					Satoshis:     1000,
					UseForChange: true,
				}},
				ChangeDestinations: []*Destination{{
					ID: testDestinationID,
				}},
			},
		}

		var newFee uint64
		newFee, err = draftTransaction.setChangeDestination(ctx, 251, 100)
		require.NoError(t, err)
		assert.Equal(t, uint64(100), newFee) // fee should not change
		assert.Equal(t, uint64(251), draftTransaction.Configuration.ChangeSatoshis)
		assert.Nil(t, draftTransaction.Configuration.ChangeDestinations)
		assert.Equal(t, uint64(1126), draftTransaction.Configuration.Outputs[0].Satoshis)
		assert.Equal(t, uint64(1000), draftTransaction.Configuration.Outputs[1].Satoshis)
		assert.Equal(t, uint64(1125), draftTransaction.Configuration.Outputs[2].Satoshis)
	})
}

// TestDraftTransaction_getInputsFromUtxos getting bt.UTXOs from SPV Wallet Engine Utxos
func TestDraftTransaction_getInputsFromUtxos(t *testing.T) {
	t.Run("invalid lockingScript", func(t *testing.T) {
		draftTransaction := &DraftTransaction{}

		reservedUtxos := []*Utxo{{
			UtxoPointer: UtxoPointer{
				OutputIndex:   123,
				TransactionID: testTxID,
			},
			Satoshis:     124235,
			ScriptPubKey: "testLockingScript",
		}}
		inputUtxos, satoshisReserved, err := draftTransaction.getInputsFromUtxos(reservedUtxos)
		require.ErrorIs(t, err, ErrInvalidLockingScript)
		assert.Nil(t, inputUtxos)
		assert.Equal(t, uint64(0), satoshisReserved)
	})

	t.Run("invalid transactionId", func(t *testing.T) {
		draftTransaction := &DraftTransaction{}

		reservedUtxos := []*Utxo{{
			UtxoPointer: UtxoPointer{
				OutputIndex:   123,
				TransactionID: "testTxID",
			},
			Satoshis:     124235,
			ScriptPubKey: testLockingScript,
		}}
		inputUtxos, satoshisReserved, err := draftTransaction.getInputsFromUtxos(reservedUtxos)
		require.ErrorIs(t, err, ErrInvalidTransactionID)
		assert.Nil(t, inputUtxos)
		assert.Equal(t, uint64(0), satoshisReserved)
	})

	t.Run("get valid", func(t *testing.T) {
		draftTransaction := &DraftTransaction{}

		reservedUtxos := []*Utxo{{
			UtxoPointer: UtxoPointer{
				OutputIndex:   123,
				TransactionID: testTxID,
			},
			Satoshis:     124235,
			ScriptPubKey: testLockingScript,
		}}
		inputUtxos, satoshisReserved, err := draftTransaction.getInputsFromUtxos(reservedUtxos)
		require.NoError(t, err)
		assert.Equal(t, uint64(124235), satoshisReserved)
		assert.Equal(t, 1, len(*inputUtxos))
		assert.Equal(t, testTxID, hex.EncodeToString((*inputUtxos)[0].TxID))
		assert.Equal(t, uint32(123), (*inputUtxos)[0].Vout)
		assert.Equal(t, testLockingScript, (*inputUtxos)[0].LockingScript.String())
		assert.Equal(t, uint64(124235), (*inputUtxos)[0].Satoshis)
	})

	t.Run("get multi", func(t *testing.T) {
		draftTransaction := &DraftTransaction{}

		reservedUtxos := []*Utxo{{
			UtxoPointer: UtxoPointer{
				OutputIndex:   124,
				TransactionID: testTxID,
			},
			Satoshis:     52313,
			ScriptPubKey: testLockingScript,
		}, {
			UtxoPointer: UtxoPointer{
				OutputIndex:   123,
				TransactionID: testTxID,
			},
			Satoshis:     124235,
			ScriptPubKey: testLockingScript,
		}}
		inputUtxos, satoshisReserved, err := draftTransaction.getInputsFromUtxos(reservedUtxos)
		require.NoError(t, err)
		assert.Equal(t, uint64(124235+52313), satoshisReserved)
		assert.Equal(t, 2, len(*inputUtxos))

		assert.Equal(t, testTxID, hex.EncodeToString((*inputUtxos)[0].TxID))
		assert.Equal(t, uint32(124), (*inputUtxos)[0].Vout)
		assert.Equal(t, testLockingScript, (*inputUtxos)[0].LockingScript.String())
		assert.Equal(t, uint64(52313), (*inputUtxos)[0].Satoshis)

		assert.Equal(t, testTxID, hex.EncodeToString((*inputUtxos)[1].TxID))
		assert.Equal(t, uint32(123), (*inputUtxos)[1].Vout)
		assert.Equal(t, testLockingScript, (*inputUtxos)[1].LockingScript.String())
		assert.Equal(t, uint64(124235), (*inputUtxos)[1].Satoshis)
	})
}

// TestDraftTransaction_AfterUpdated after updated tests
func TestDraftTransaction_AfterUpdated(t *testing.T) {
	t.Run("cancel draft - update utxo reservation", func(t *testing.T) {
		ctx, client, deferMe := CreateTestSQLiteClient(t, false, false)
		defer deferMe()
		reservationDraftID, _ := utils.RandomHex(32)

		utxo := newUtxo(testXPubID, testTxID, testLockingScript, 0, 100000,
			append(client.DefaultModelOptions(), New())...)
		utxo.DraftID.Valid = true
		utxo.DraftID.String = reservationDraftID
		utxo.ReservedAt.Valid = true
		utxo.ReservedAt.Time = time.Now().UTC()
		err := utxo.Save(ctx)
		require.NoError(t, err)

		var gUtxo *Utxo
		gUtxo, err = getUtxo(ctx, testTxID, 0, client.DefaultModelOptions()...)
		require.NoError(t, err)
		assert.True(t, gUtxo.DraftID.Valid)
		assert.Equal(t, reservationDraftID, gUtxo.DraftID.String)
		assert.True(t, gUtxo.ReservedAt.Valid)

		draftTransaction := &DraftTransaction{
			Model: *NewBaseModel(
				ModelDraftTransaction,
				client.DefaultModelOptions()...,
			),
			TransactionBase: TransactionBase{ID: reservationDraftID},
			Configuration:   TransactionConfig{},
			Status:          DraftStatusCanceled,
		}

		err = draftTransaction.AfterUpdated(ctx)
		require.NoError(t, err)

		var gUtxo2 *Utxo
		gUtxo2, err = getUtxo(ctx, testTxID, 0, client.DefaultModelOptions()...)
		require.NoError(t, err)
		assert.False(t, gUtxo2.DraftID.Valid)
		assert.False(t, gUtxo2.ReservedAt.Valid)
	})
}

// TestDraftTransaction_addIncludeUtxos will test the method addIncludeUtxos()
func TestDraftTransaction_addIncludeUtxos(t *testing.T) {
	t.Run("no includeUtxos", func(t *testing.T) {
		ctx := context.Background()
		draft := &DraftTransaction{
			Configuration: TransactionConfig{},
		}
		includeUtxoSatoshis, err := draft.addIncludeUtxos(ctx)
		require.NoError(t, err)
		assert.Len(t, draft.Configuration.Inputs, 0)
		assert.Equal(t, uint64(0), includeUtxoSatoshis)
	})

	t.Run("with includeUtxos", func(t *testing.T) {
		ctx, client, deferMe := CreateTestSQLiteClient(t, false, false)
		defer deferMe()

		destination := newDestination(testXPubID, testSTASScriptPubKey, client.DefaultModelOptions(New())...)
		err := destination.Save(ctx)
		require.NoError(t, err)

		utxo := newUtxo(testXPubID, testSTAStxID, testSTASLockingScript, 0, 550, client.DefaultModelOptions(New())...)
		err = utxo.Save(ctx)
		require.NoError(t, err)

		draft, err := newDraftTransaction("", &TransactionConfig{
			IncludeUtxos: []*UtxoPointer{{
				testSTAStxID,
				0,
			}},
		}, client.DefaultModelOptions(New())...)
		require.NoError(t, err)

		var includeUtxoSatoshis uint64
		includeUtxoSatoshis, err = draft.addIncludeUtxos(ctx)
		require.NoError(t, err)
		assert.Len(t, draft.Configuration.Inputs, 1)
		assert.Equal(t, uint64(550), includeUtxoSatoshis)
		assert.Equal(t, testSTASLockingScript, draft.Configuration.Inputs[0].ScriptPubKey)
		assert.Equal(t, testSTASScriptPubKey, draft.Configuration.Inputs[0].Destination.LockingScript)
	})
}

// TestDraftTransaction_addOutputsToTx will test the method addOutputsToTx()
func TestDraftTransaction_addOutputsToTx(t *testing.T) {
	t.Run("no output", func(t *testing.T) {
		draft := &DraftTransaction{
			Configuration: TransactionConfig{
				Outputs: []*TransactionOutput{{
					Satoshis: 0,
				}},
			},
		}
		tx := bt.NewTx()
		err := draft.addOutputsToTx(tx)
		require.NoError(t, err)
	})

	t.Run("no output", func(t *testing.T) {
		draft := &DraftTransaction{
			Configuration: TransactionConfig{
				Outputs: []*TransactionOutput{{
					Scripts: []*ScriptOutput{{
						Satoshis: 0,
						Script:   testDraftLockingScript,
					}},
				}},
			},
		}
		tx := bt.NewTx()
		err := draft.addOutputsToTx(tx)
		require.ErrorIs(t, err, ErrOutputValueTooLow)
		assert.Len(t, tx.Outputs, 0)
	})

	t.Run("normal address", func(t *testing.T) {
		draft := &DraftTransaction{
			Configuration: TransactionConfig{
				Outputs: []*TransactionOutput{{
					Scripts: []*ScriptOutput{{
						Satoshis: 1000,
						Script:   testDraftLockingScript,
					}},
				}},
			},
		}
		tx := bt.NewTx()
		err := draft.addOutputsToTx(tx)
		require.NoError(t, err)
		assert.Len(t, tx.Outputs, 1)
		assert.Equal(t, uint64(1000), tx.Outputs[0].Satoshis)
		assert.Equal(t, testDraftLockingScript, tx.Outputs[0].LockingScript.String())
	})

	t.Run("op return", func(t *testing.T) {
		draft := &DraftTransaction{
			Configuration: TransactionConfig{
				Outputs: []*TransactionOutput{{
					Scripts: []*ScriptOutput{{
						Satoshis:   0,
						Script:     testDraftLockingScript,
						ScriptType: utils.ScriptTypeNullData,
					}},
				}},
			},
		}
		tx := bt.NewTx()
		err := draft.addOutputsToTx(tx)
		require.NoError(t, err)
		assert.Len(t, tx.Outputs, 1)
		assert.Equal(t, uint64(0), tx.Outputs[0].Satoshis)
		assert.Equal(t, testDraftLockingScript, tx.Outputs[0].LockingScript.String())
	})

	t.Run("op return", func(t *testing.T) {
		draft := &DraftTransaction{
			Configuration: TransactionConfig{
				Outputs: []*TransactionOutput{{
					Scripts: []*ScriptOutput{{
						Satoshis:   1000,
						Script:     testDraftLockingScript,
						ScriptType: utils.ScriptTypeNullData,
					}},
				}},
			},
		}
		tx := bt.NewTx()
		err := draft.addOutputsToTx(tx)
		require.ErrorIs(t, err, ErrInvalidOpReturnOutput)
	})
}

func createDraftTransactionFromHex(hex string, inInfo []interface{}, feeUnit *utils.FeeUnit) (*DraftTransaction, *bt.Tx, error) {
	tx, err := bt.NewTxFromString(hex)
	if err != nil {
		return nil, nil, err
	}

	feePaid := uint64(0)

	inputs := make([]*TransactionInput, 0)
	for txIndex := range tx.Inputs {
		in := inInfo[txIndex].(map[string]interface{})
		satoshis := uint64(in["satoshis"].(float64))
		lockingScript := in["locking_script"].(string)
		input := TransactionInput{
			Utxo: Utxo{
				UtxoPointer: UtxoPointer{
					OutputIndex:   uint32(txIndex),
					TransactionID: tx.TxID(),
				},
				XpubID:       testXPubID,
				Satoshis:     satoshis,
				ScriptPubKey: lockingScript,
				Type:         utils.ScriptTypePubKeyHash,
			},
			Destination: Destination{
				XpubID:        testXPubID,
				LockingScript: lockingScript,
				Type:          utils.ScriptTypePubKeyHash,
				Chain:         0,
				Num:           uint32(txIndex),
				Address:       testExternalAddress,
				DraftID:       testDraftID,
			},
		}
		feePaid += input.Satoshis

		inputs = append(inputs, &input)
	}

	outputs := make([]*TransactionOutput, 0)
	for _, txOutput := range tx.Outputs {
		output := TransactionOutput{
			Satoshis: txOutput.Satoshis,
			Scripts: []*ScriptOutput{{
				Address: testExternalAddress,
				Script:  txOutput.LockingScript.String(),
			}},
			To: testExternalAddress,
		}
		feePaid -= output.Satoshis

		outputs = append(outputs, &output)
	}

	draftConfig := &TransactionConfig{
		ChangeDestinations: []*Destination{{}}, // set to not nil, otherwise will be overwritten when processing
		Fee:                0,
		FeeUnit:            feeUnit,
		Inputs:             inputs,
		Outputs:            outputs,
	}

	draftTransaction, err := newDraftTransaction(testXPub, draftConfig)
	if err != nil {
		return nil, nil, err
	}

	return draftTransaction, tx, nil
}

func TestDraftTransaction_estimateFees(t *testing.T) {
	t.Run("json data", func(t *testing.T) {
		jsonFile, err := os.Open("./tests/model_draft_transactions_test.json")
		require.NoError(t, err)
		defer func() {
			_ = jsonFile.Close()
		}()

		byteValue, bErr := ioutil.ReadAll(jsonFile)
		require.NoError(t, bErr)

		var testData map[string]interface{}
		err = json.Unmarshal(byteValue, &testData)
		require.NoError(t, err)

		for _, inTx := range testData["rawTransactions"].([]interface{}) {
			in := inTx.(map[string]interface{})
			txID := in["txId"].(string)
			var feeUnit *utils.FeeUnit
			if _, ok := in["feeUnit"]; ok {
				b, _ := json.Marshal(in["feeUnit"])
				_ = json.Unmarshal(b, &feeUnit)
			} else {
				feeUnit = chainstate.MockDefaultFee
			}
			draftTransaction, tx, err2 := createDraftTransactionFromHex(in["hex"].(string), in["inputs"].([]interface{}), feeUnit)
			require.NoError(t, err2)
			assert.Equal(t, txID, tx.TxID())
			assert.IsType(t, DraftTransaction{}, *draftTransaction)
			assert.IsType(t, bt.Tx{}, *tx)

			realFee := uint64(0)
			for _, input := range in["inputs"].([]interface{}) {
				i := input.(map[string]interface{})
				realFee += uint64(i["satoshis"].(float64))
			}
			for _, output := range tx.Outputs {
				realFee -= output.Satoshis
			}

			realSize := uint64(float64(len(in["hex"].(string))) / 2)
			sizeEstimate := draftTransaction.estimateSize()
			feeEstimate := draftTransaction.estimateFee(feeUnit, 0)
			assert.GreaterOrEqual(t, sizeEstimate, realSize)
			assert.GreaterOrEqual(t, feeEstimate, realFee)
		}
	})
}

func TestDraftTransaction_SignInputs(t *testing.T) {
	ctx, client, deferMe := CreateTestSQLiteClient(t, true, true)
	defer deferMe()

	xPrivString := "xprv9s21ZrQH143K31pvNoYNcRZjtdJXnNVEc5NmBbgJmEg27YWbZVL7jTLQhPELqAR7tcJTnF9AJLwVN5w3ABZvrfeDLm4vnBDw76bkx8a2NxK"
	xPrivHD, err := bitcoin.GenerateHDKeyFromString(xPrivString)
	require.NoError(t, err)
	xPubHD, _ := xPrivHD.Neuter()
	xPubID := utils.Hash(xPubHD.String())

	xPub := newXpub(xPubHD.String(), client.DefaultModelOptions(New())...)
	err = xPub.Save(ctx)
	require.NoError(t, err)

	// Derive the child key (chain)
	var chainKey *bip32.ExtendedKey
	if chainKey, err = xPrivHD.Child(
		0,
	); err != nil {
		return
	}

	// Derive the child key (num)
	var numKey *bip32.ExtendedKey
	if numKey, err = chainKey.Child(
		0,
	); err != nil {
		return
	}

	// Get the private key
	var privateKey *bec.PrivateKey
	if privateKey, err = bitcoin.GetPrivateKeyFromHDKey(
		numKey,
	); err != nil {
		return
	}

	// create a destination for the utxo
	lockingScript := "76a91447868e6b13de36e2739d8f2a9e0e0a323ad9b8ff88ac"
	destination := newDestination(xPubID, lockingScript, append(client.DefaultModelOptions(), New())...)
	err = destination.Save(ctx)
	require.NoError(t, err)

	// create a utxo with enough output for all our tests
	utxo := newUtxo(xPubID, testTxID, lockingScript, 0, 12229, client.DefaultModelOptions(New())...)
	err = utxo.Save(ctx)
	require.NoError(t, err)

	transaction, err := txFromHex(testTxHex, append(client.DefaultModelOptions(), New())...)
	require.NoError(t, err)

	err = transaction.Save(ctx)
	require.NoError(t, err)

	tests := []struct {
		name    string
		config  *TransactionConfig
		xPriv   *bip32.ExtendedKey
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "sign 1",
			config: &TransactionConfig{
				SendAllTo: &TransactionOutput{To: "1AqYEDUf16CHaD2guBLHHhosfV2AyYJLz"},
			},
			xPriv:   xPrivHD,
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := newDraftTransaction(xPub.rawXpubKey, tt.config, client.DefaultModelOptions(New())...)
			require.NoError(t, err)
			err = m.createTransactionHex(ctx)
			require.NoError(t, err)

			var gotSignedHex string
			gotSignedHex, err = m.SignInputs(tt.xPriv)
			if !tt.wantErr(t, err, fmt.Sprintf("SignInputs(%v)", tt.xPriv)) {
				return
			}

			var tx *bt.Tx
			tx, err = bt.NewTxFromString(gotSignedHex)

			var ls *bscript.Script
			if ls, err = bscript.NewFromHexString(
				lockingScript,
			); err != nil {
				return
			}
			tx.Inputs[0].PreviousTxScript = ls
			tx.Inputs[0].PreviousTxSatoshis = 12229

			require.NoError(t, err)
			assert.True(t, tx.Version > 0)
			for _, input := range tx.Inputs {
				var unlocker string
				unlocker, err = input.UnlockingScript.ToASM()
				require.NoError(t, err)
				scriptParts := strings.Split(unlocker, " ")
				pubKey := hex.EncodeToString(privateKey.PubKey().SerialiseCompressed())

				var hash []byte
				hash, err = tx.CalcInputSignatureHash(0, sighash.AllForkID)
				require.NoError(t, err)

				var hash32 [32]byte
				copy(hash32[:], hash)
				var verified bool
				verified, err = bitcoin.VerifyMessageDER(hash32, pubKey, scriptParts[0])
				require.NoError(t, err)
				assert.True(t, verified)
			}
		})
	}
}

func initSimpleTestCase(t *testing.T) (context.Context, ClientInterface, func()) {
	ctx, client, deferMe := CreateTestSQLiteClient(t, false, true, withTaskManagerMockup())

	xPub := newXpub(testXPub, append(client.DefaultModelOptions(), New())...)
	xPub.CurrentBalance = 100000
	err := xPub.Save(ctx)
	require.NoError(t, err)

	destination := newDestination(testXPubID, testLockingScript,
		append(client.DefaultModelOptions(), New())...)
	err = destination.Save(ctx)
	require.NoError(t, err)

	utxo := newUtxo(testXPubID, testTxID, testLockingScript, 0, 100000,
		append(client.DefaultModelOptions(), New())...)
	err = utxo.Save(ctx)
	require.NoError(t, err)

	transaction, err := txFromHex(testTxHex, append(client.DefaultModelOptions(), New())...)
	require.NoError(t, err)

	err = transaction.processUtxos(ctx)
	require.NoError(t, err)

	err = transaction.Save(ctx)
	require.NoError(t, err)

	return ctx, client, deferMe
}
