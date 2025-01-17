package fixtures

import (
	"testing"

	"github.com/bitcoin-sv/go-sdk/script"
	"github.com/bitcoin-sv/go-sdk/spv"
	"github.com/stretchr/testify/require"
)

func givenTXSpec(t *testing.T) GivenTXSpec {
	return GivenTX(t).
		WithSingleSourceInputs(2, 3, 4).
		WithP2PKHOutput(1).
		WithOPReturn("hello world").
		WithOutputScriptParts(OpCode(script.OpRETURN), PushData("hello world"))
}

/*
NOTE: Since the GivenTX is only a helper for other tests, and helpers usually don't need to be tested,
But in this case we have to make an exception, because we want to make sure that the mechanism
for given specification returns always the same hexes (BEEF, EF, raw).
*/

func TestMockTXGeneration(t *testing.T) {
	tests := map[string]struct {
		spec           GivenTXSpec
		shouldBeSigned bool
		beef           string
		rawTX          string
		ef             string
	}{
		"empty tx": {
			spec:  GivenTX(t),
			beef:  "0100beef00010100000000000000000000",
			rawTX: "01000000000000000000",
			ef:    "010000000000000000ef000000000000",
		},
		"signed complex tx": {
			spec:  givenTXSpec(t),
			beef:  "0100beef01fde8030101000018aec743415be3827937f37e1bd6b1930e8d7cfbb22bcd72a373438399f1dcb70201000000012e3f4683e173b40a20527fe5719633ba070df649983614886e90e45aecf2ac56000000006a47304402205c05a6c9eadda5da97eddb55a37178e946d4be0b151670233bc76cbebbee011c022018845d47e6fa8dd0100258bf855dc4990ad49c72de65f04b1b85ea146d1f117a4121034d2d6d23fbcb6eefe3e80c47044e36797dcb80d0ac5e96e732ef03c3c550a116ffffffff0302000000000000001976a91494677c56fa2968644c90a517214338b4139899ce88ac03000000000000001976a91494677c56fa2968644c90a517214338b4139899ce88ac04000000000000001976a91494677c56fa2968644c90a517214338b4139899ce88ac000000000100010000000318aec743415be3827937f37e1bd6b1930e8d7cfbb22bcd72a373438399f1dcb7000000006b483045022100f5fb8c86bb12a2cc45c2f62cf590725209362c8da0618010a3009625462cddf302205856c21d1daf1f2758c1d03f6acb2f27e361f998d48d732cb8ab8d494eac8e894121034d2d6d23fbcb6eefe3e80c47044e36797dcb80d0ac5e96e732ef03c3c550a116ffffffff18aec743415be3827937f37e1bd6b1930e8d7cfbb22bcd72a373438399f1dcb7010000006a47304402207877451733ab726fc23288084dd99a6450490a70db06483c1709f7b482610cbc02207b7c84cca5e345a2a3a56fc59ba8ee12cadb856d3032a7aff7c3719c84048e144121034d2d6d23fbcb6eefe3e80c47044e36797dcb80d0ac5e96e732ef03c3c550a116ffffffff18aec743415be3827937f37e1bd6b1930e8d7cfbb22bcd72a373438399f1dcb7020000006b483045022100e500eba499e76490c6936c3ac13e1ed808edcd671d113f5fa76f23951116cbd40220576ff4993ce3ec757949aa0fc2bfe1b541660b520f0d1ed745cf79a210d9e3ff4121034d2d6d23fbcb6eefe3e80c47044e36797dcb80d0ac5e96e732ef03c3c550a116ffffffff0301000000000000001976a9143cf53c49c322d9d811728182939aee2dca087f9888ac00000000000000000e006a0b68656c6c6f20776f726c6400000000000000000d6a0b68656c6c6f20776f726c640000000000",
			rawTX: "010000000318aec743415be3827937f37e1bd6b1930e8d7cfbb22bcd72a373438399f1dcb7000000006b483045022100f5fb8c86bb12a2cc45c2f62cf590725209362c8da0618010a3009625462cddf302205856c21d1daf1f2758c1d03f6acb2f27e361f998d48d732cb8ab8d494eac8e894121034d2d6d23fbcb6eefe3e80c47044e36797dcb80d0ac5e96e732ef03c3c550a116ffffffff18aec743415be3827937f37e1bd6b1930e8d7cfbb22bcd72a373438399f1dcb7010000006a47304402207877451733ab726fc23288084dd99a6450490a70db06483c1709f7b482610cbc02207b7c84cca5e345a2a3a56fc59ba8ee12cadb856d3032a7aff7c3719c84048e144121034d2d6d23fbcb6eefe3e80c47044e36797dcb80d0ac5e96e732ef03c3c550a116ffffffff18aec743415be3827937f37e1bd6b1930e8d7cfbb22bcd72a373438399f1dcb7020000006b483045022100e500eba499e76490c6936c3ac13e1ed808edcd671d113f5fa76f23951116cbd40220576ff4993ce3ec757949aa0fc2bfe1b541660b520f0d1ed745cf79a210d9e3ff4121034d2d6d23fbcb6eefe3e80c47044e36797dcb80d0ac5e96e732ef03c3c550a116ffffffff0301000000000000001976a9143cf53c49c322d9d811728182939aee2dca087f9888ac00000000000000000e006a0b68656c6c6f20776f726c6400000000000000000d6a0b68656c6c6f20776f726c6400000000",
			ef:    "010000000000000000ef0318aec743415be3827937f37e1bd6b1930e8d7cfbb22bcd72a373438399f1dcb7000000006b483045022100f5fb8c86bb12a2cc45c2f62cf590725209362c8da0618010a3009625462cddf302205856c21d1daf1f2758c1d03f6acb2f27e361f998d48d732cb8ab8d494eac8e894121034d2d6d23fbcb6eefe3e80c47044e36797dcb80d0ac5e96e732ef03c3c550a116ffffffff02000000000000001976a91494677c56fa2968644c90a517214338b4139899ce88ac18aec743415be3827937f37e1bd6b1930e8d7cfbb22bcd72a373438399f1dcb7010000006a47304402207877451733ab726fc23288084dd99a6450490a70db06483c1709f7b482610cbc02207b7c84cca5e345a2a3a56fc59ba8ee12cadb856d3032a7aff7c3719c84048e144121034d2d6d23fbcb6eefe3e80c47044e36797dcb80d0ac5e96e732ef03c3c550a116ffffffff03000000000000001976a91494677c56fa2968644c90a517214338b4139899ce88ac18aec743415be3827937f37e1bd6b1930e8d7cfbb22bcd72a373438399f1dcb7020000006b483045022100e500eba499e76490c6936c3ac13e1ed808edcd671d113f5fa76f23951116cbd40220576ff4993ce3ec757949aa0fc2bfe1b541660b520f0d1ed745cf79a210d9e3ff4121034d2d6d23fbcb6eefe3e80c47044e36797dcb80d0ac5e96e732ef03c3c550a116ffffffff04000000000000001976a91494677c56fa2968644c90a517214338b4139899ce88ac0301000000000000001976a9143cf53c49c322d9d811728182939aee2dca087f9888ac00000000000000000e006a0b68656c6c6f20776f726c6400000000000000000d6a0b68656c6c6f20776f726c6400000000",
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			spec := test.spec

			// when
			tx := spec.TX()
			ok, err := spv.VerifyScripts(tx)

			// then:
			require.NoError(t, err)
			require.True(t, ok)

			require.Equal(t, test.beef, spec.BEEF())
			require.Equal(t, test.rawTX, spec.RawTX())
			require.Equal(t, test.ef, spec.EF())
		})
	}
}

func TestSpike(t *testing.T) {

	tx4 := GivenTX(t).
		WithInput(1000).
		WithP2PKHOutput(100).TX()

	// tx6 := tx4.InputIdx(0).SourceTransaction

	tx1 := GivenTX(t).
		WithInputFromUTXO(tx4, 0).
		WithP2PKHOutput(1).TX()

	for _, input := range tx1.Inputs {
		input.UnlockingScriptTemplate = RecipientInternal.P2PKHUnlockingScriptTemplate()
	}
	err := tx1.Sign()
	require.NoError(t, err)

	tx5 := GivenTX(t).
		WithInput(20).
		WithP2PKHOutput(2).
		WithP2PKHOutput(2).TX()

	// tx7 := tx5.InputIdx(0).SourceTransaction

	tx3 := GivenTX(t).
		WithInputFromUTXO(tx5, 0).
		WithP2PKHOutput(1).TX()
	for _, input := range tx3.Inputs {
		input.UnlockingScriptTemplate = RecipientInternal.P2PKHUnlockingScriptTemplate()
	}
	err = tx3.Sign()
	require.NoError(t, err)

	txn := GivenTX(t).
		WithInputFromUTXO(tx5, 1).
		WithP2PKHOutput(1).TX()
	for _, input := range txn.Inputs {
		input.UnlockingScriptTemplate = RecipientInternal.P2PKHUnlockingScriptTemplate()
	}
	err = txn.Sign()
	require.NoError(t, err)

	tx2 := GivenTX(t).
		WithInput(1000).
		WithP2PKHOutput(100).TX()

	tx0 := GivenTX(t).
		WithInputFromUTXO(tx1, 0).
		WithInputFromUTXO(tx2, 0).
		WithInputFromUTXO(tx3, 0).
		TX()

	for i, input := range tx0.Inputs {
		verified, err := spv.VerifyScripts(input.SourceTransaction)
		require.NoError(t, err, "input %s", i)
		require.True(t, verified, "input %s", i)
	}
}
