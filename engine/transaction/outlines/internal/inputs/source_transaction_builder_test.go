package inputs

import (
	"testing"

	"github.com/bitcoin-sv/go-sdk/spv"
	sdk "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
	"github.com/stretchr/testify/require"
)

func TestSourceTransactionBuilder_Case2(t *testing.T) {
	// given:

	// beef node:
	const tx6BEEFHex = "0100beef01fde803010100008cce48ad84ba7ab6578b8bcdf2ed68aec49b098c43fdd8b47918afb94aeac8ab0101000000012e3f4683e173b40a20527fe5719633ba070df649983614886e90e45aecf2ac56000000006b483045022100fc3d9faa7c983d4a490e9e3ad13da6cb6b8f8da967a6585775ffffb307349dfa02202a7bea4ba5c27cf37de0234a0ab5d533a9715e15183868092e74fa900eecd9f64121034d2d6d23fbcb6eefe3e80c47044e36797dcb80d0ac5e96e732ef03c3c550a116ffffffff010a000000000000001976a91494677c56fa2968644c90a517214338b4139899ce88ac000000000100"

	// raw tx node:
	const tx1Hex = "01000000018cce48ad84ba7ab6578b8bcdf2ed68aec49b098c43fdd8b47918afb94aeac8ab000000006a473044022041a1354250793efbf331dee36e20e6d919db6416ff443087a1a49133d0c4a5e30220715fca40622835bcd6216e6fd95f22a902f64212d52146e1d9798bdfb7cfd4414121034d2d6d23fbcb6eefe3e80c47044e36797dcb80d0ac5e96e732ef03c3c550a116ffffffff0102000000000000001976a9143cf53c49c322d9d811728182939aee2dca087f9888ac00000000"

	// root tx node:
	const tx0Hex = "0100000001917b2f6a7523a439ecaae9e00df81d744241d34edf02d6c1959f76760d7bf78b0000000000ffffffff0000000000"

	// node edges:
	const tx0InputIdx0SourceID = "8bf77b0d76769f95c1d602df4ed34142741df80de0e9aaec39a423756a2f7b91"
	const tx1InputIdx0SourceID = "abc8ea4ab9af1879b4d8fd438c099bc4ae68edf2cd8b8b57b67aba84ad48ce8c"

	tx0, err := sdk.NewTransactionFromHex(tx0Hex)
	require.NoError(t, err)
	require.NotEmpty(t, tx0)

	builder := SourceTransactionBuilder{Tx: tx0}

	// when:
	err = builder.Build(TxQueryResultSlice{
		{
			SourceTXID: tx0InputIdx0SourceID,
			RawHex:     Ptr(tx1Hex),
		},
		{
			SourceTXID: tx1InputIdx0SourceID,
			BeefHex:    Ptr(tx6BEEFHex),
		},
	})

	// then:
	require.NoError(t, err)
}

func TestSourceTransactionBuilder_Case3(t *testing.T) {
	// given:

	// raw tx nodes:
	const tx0Hex = "010000000396ae2b2f9ed0a876bad6a2e0434149974e95ead5c9661d7051b6d16aba00a46d0000000000ffffffff4728499cf1517d3cba0bcf392e32ca3fd0e6ddc5a61196382fbe6328291e81300000000000ffffffffd651103d5af4d2acad3a154385e6ab5ff5bb938c099cb20c4a65f9b6c60d24040000000000ffffffff0000000000"
	const tx1Hex = "01000000018cce48ad84ba7ab6578b8bcdf2ed68aec49b098c43fdd8b47918afb94aeac8ab000000006a4730440220679101ae3f0006a4727149bb8e858d6bd48d9821bff21d6150ffe48740dae3de022008449c4d7f6f9908bd5d652445617a0c28bc63d2a9e7939fb450051fd795628c4121034d2d6d23fbcb6eefe3e80c47044e36797dcb80d0ac5e96e732ef03c3c550a116ffffffff0109000000000000001976a9143cf53c49c322d9d811728182939aee2dca087f9888ac00000000"
	const tx2Hex = "0100000001b73487ab61e8392624f2a84f78ac30c01872cef04bbfa82fb9363df9f825c6da000000006a4730440220182c473bfc8b820ac0b822cffe789202492662c66befc6bafa3dbd946c0ae09e022035039d4cdf48ee5fedbad422582ce8c512946157e96df17f4c0c366898dcfb604121034d2d6d23fbcb6eefe3e80c47044e36797dcb80d0ac5e96e732ef03c3c550a116ffffffff0107000000000000001976a9143cf53c49c322d9d811728182939aee2dca087f9888ac00000000"
	const tx3Hex = "0100000001467b1fc90961eed9e84ffe9e57099f41234c8225b8a0e2146fb00a5ba2e76419000000006a47304402205e016051c126abe69ad21d39ec595e408e2925d4dbcb1b4b12d9c0a7576e85a402204616dee62952548f9c1682115f94889b76ad0d12ea4cc1d6d3d26f6e26dec1e64121034d2d6d23fbcb6eefe3e80c47044e36797dcb80d0ac5e96e732ef03c3c550a116ffffffff0108000000000000001976a9143cf53c49c322d9d811728182939aee2dca087f9888ac00000000"

	// beef nodes:
	const tx4BEEFHex = "0100beef01fde803010100008cce48ad84ba7ab6578b8bcdf2ed68aec49b098c43fdd8b47918afb94aeac8ab0101000000012e3f4683e173b40a20527fe5719633ba070df649983614886e90e45aecf2ac56000000006b483045022100fc3d9faa7c983d4a490e9e3ad13da6cb6b8f8da967a6585775ffffb307349dfa02202a7bea4ba5c27cf37de0234a0ab5d533a9715e15183868092e74fa900eecd9f64121034d2d6d23fbcb6eefe3e80c47044e36797dcb80d0ac5e96e732ef03c3c550a116ffffffff010a000000000000001976a91494677c56fa2968644c90a517214338b4139899ce88ac000000000100"
	const tx5BEEFHex = "0100beef01fde80301010000467b1fc90961eed9e84ffe9e57099f41234c8225b8a0e2146fb00a5ba2e764190101000000012e3f4683e173b40a20527fe5719633ba070df649983614886e90e45aecf2ac56000000006b483045022100e66a8f4c94864466f3e8582cb2587fab05bfe02fc8b4127e952684104f8f289602204991685e305b490196be6ba021972dc7c71640c92bea689af8b94932e49db86a4121034d2d6d23fbcb6eefe3e80c47044e36797dcb80d0ac5e96e732ef03c3c550a116ffffffff0109000000000000001976a91494677c56fa2968644c90a517214338b4139899ce88ac000000000100"
	const tx6BEEFHex = "0100beef01fde80301010000b73487ab61e8392624f2a84f78ac30c01872cef04bbfa82fb9363df9f825c6da0101000000012e3f4683e173b40a20527fe5719633ba070df649983614886e90e45aecf2ac56000000006a473044022012ade0924fa8f8675874e2fe304720895c529fd187e1e6b2a150f2b0e93187e602206deeb26e88c5222c9f287c72c2c5124426023a184f0e6fe02c807cf13012c22d4121034d2d6d23fbcb6eefe3e80c47044e36797dcb80d0ac5e96e732ef03c3c550a116ffffffff0108000000000000001976a91494677c56fa2968644c90a517214338b4139899ce88ac000000000100"

	// node edges:
	const tx1InputIdx0SourceID = "abc8ea4ab9af1879b4d8fd438c099bc4ae68edf2cd8b8b57b67aba84ad48ce8c"
	const tx2InputIdx0SourceID = "dac625f8f93d36b92fa8bf4bf0ce7218c030ac784fa8f2242639e861ab8734b7"
	const tx3InputIdx0SourceID = "1964e7a25b0ab06f14e2a0b825824c23419f09579efe4fe8d9ee6109c91f7b46"

	const tx0InputIdx0SourceID = "6da400ba6ad1b651701d66c9d5ea954e97494143e0a2d6ba76a8d09e2f2bae96"
	const tx0InputIdx1SourceID = "30811e292863be2f389611a6c5dde6d03fca322e39cf0bba3c7d51f19c492847"
	const tx0InputIdx2SourceID = "04240dc6b6f9654a0cb29c098c93bbf55fabe68543153aadacd2f45a3d1051d6"

	tx0, err := sdk.NewTransactionFromHex(tx0Hex)
	require.NoError(t, err)
	require.NotEmpty(t, tx0)

	builder := SourceTransactionBuilder{Tx: tx0}

	// when:
	err = builder.Build(TxQueryResultSlice{
		{
			SourceTXID: tx0InputIdx0SourceID,
			RawHex:     Ptr(tx1Hex),
		},
		{
			SourceTXID: tx1InputIdx0SourceID,
			BeefHex:    Ptr(tx4BEEFHex),
		},
		{
			SourceTXID: tx0InputIdx1SourceID,
			RawHex:     Ptr(tx3Hex),
		},
		{
			SourceTXID: tx3InputIdx0SourceID,
			BeefHex:    Ptr(tx5BEEFHex),
		},
		{
			SourceTXID: tx0InputIdx2SourceID,
			RawHex:     Ptr(tx2Hex),
		},
		{
			SourceTXID: tx2InputIdx0SourceID,
			BeefHex:    Ptr(tx6BEEFHex),
		},
	})

	// then:
	require.NoError(t, err)
}

func TestSourceTransactionBuilder_Case6(t *testing.T) {
	// given:

	// raw tx nodes:
	const tx0Hex = "010000000319e1523dc73f409cec5f8a0f566ab1936d8f5bc9fcdbae65703e1dc2146c5b520000000000fffffffff4fe86255a7f800f6329c540a561721111db17406d350fc003b68a01bc4986db0000000000ffffffff39d7b89f13e4ce3604232e61e36cdfb68cf287ee246ab4c4d1ea87134b5defd40000000000ffffffff0000000000"
	const tx1Hex = "0100000001917b2f6a7523a439ecaae9e00df81d744241d34edf02d6c1959f76760d7bf78b000000006a473044022075b67d42bcb2883aa6dda02f13b28864119bd31309bd563ac036d7d51cfb63d9022073da123fb0ce331163572219b611abcb6eeb1a3455fadd64a86f3b2c2625ba6d4121020edd3dad5a457e087a448ab7dcf2b23ff5206e3cbc34b8cf006680fa15d45f5effffffff0101000000000000001976a9143cf53c49c322d9d811728182939aee2dca087f9888ac00000000"
	const tx2Hex = "01000000017cd347a6a099f82cde68faec941e888ebc3489b25403e3ffedd3280f3fa4cc03000000006b4830450221008bc002baff1b1bc89c3af82cb1cdfc1bc5d081356a5e10eb4e6042027c2d3ee4022029da3ecaaa62f4213b124ac9ae4254ab2e2910b14451cf807b174d71fb8485a94121034d2d6d23fbcb6eefe3e80c47044e36797dcb80d0ac5e96e732ef03c3c550a116ffffffff01c8000000000000001976a9143cf53c49c322d9d811728182939aee2dca087f9888ac00000000"
	const tx3Hex = "010000000112a3c30238c9d65b3c21c152e950a20acaf0db04b4b4b92e0612313d7935b2d8000000006a473044022036d64e8a6c2ed963dcb76639d3ea46a3e2544d5b67110cc4eb5a942b1b6919a302205bcd1bcc4bbf833de4f726f762b93426e4d8aa645167fd1c964f3eca4130b3b14121020edd3dad5a457e087a448ab7dcf2b23ff5206e3cbc34b8cf006680fa15d45f5effffffff010a000000000000001976a9143cf53c49c322d9d811728182939aee2dca087f9888ac00000000"
	const tx4Hex = "01000000018cce48ad84ba7ab6578b8bcdf2ed68aec49b098c43fdd8b47918afb94aeac8ab000000006a473044022041a1354250793efbf331dee36e20e6d919db6416ff443087a1a49133d0c4a5e30220715fca40622835bcd6216e6fd95f22a902f64212d52146e1d9798bdfb7cfd4414121034d2d6d23fbcb6eefe3e80c47044e36797dcb80d0ac5e96e732ef03c3c550a116ffffffff0102000000000000001976a9143cf53c49c322d9d811728182939aee2dca087f9888ac00000000"
	const tx5Hex = "010000000157829e06b620a578ebf54bb3e1f4115bbdaeae1f5e5483bb80ef68db936608cf000000006b483045022100c8700b9912265fac01c23b2d64e9abfea21e31a6c353cdd1a95cdea6a694b48302207296968f8abe2fa9273da15b11f44f223695fcaec0da0670ca6b82a8eb7c88ce4121034d2d6d23fbcb6eefe3e80c47044e36797dcb80d0ac5e96e732ef03c3c550a116ffffffff0114000000000000001976a9143cf53c49c322d9d811728182939aee2dca087f9888ac00000000"

	// beef nodes:
	const tx6BEEFHex = "0100beef01fde803010100008cce48ad84ba7ab6578b8bcdf2ed68aec49b098c43fdd8b47918afb94aeac8ab0101000000012e3f4683e173b40a20527fe5719633ba070df649983614886e90e45aecf2ac56000000006b483045022100fc3d9faa7c983d4a490e9e3ad13da6cb6b8f8da967a6585775ffffb307349dfa02202a7bea4ba5c27cf37de0234a0ab5d533a9715e15183868092e74fa900eecd9f64121034d2d6d23fbcb6eefe3e80c47044e36797dcb80d0ac5e96e732ef03c3c550a116ffffffff010a000000000000001976a91494677c56fa2968644c90a517214338b4139899ce88ac000000000100"
	const tx7BEEFHex = "0100beef01fde8030101000057829e06b620a578ebf54bb3e1f4115bbdaeae1f5e5483bb80ef68db936608cf0101000000012e3f4683e173b40a20527fe5719633ba070df649983614886e90e45aecf2ac56000000006b483045022100e4ebed1698286072db434fdcf1fd08a3d999f257ddfaf6376b3417b2ff4805e002200d2b7582a8cda8ea61986b16922a51d5e2c0c548457bd0bfda3d24390f1fbd8f4121034d2d6d23fbcb6eefe3e80c47044e36797dcb80d0ac5e96e732ef03c3c550a116ffffffff0164000000000000001976a91494677c56fa2968644c90a517214338b4139899ce88ac000000000100"
	const tx8BEEFHex = "0100beef01fde803010100007cd347a6a099f82cde68faec941e888ebc3489b25403e3ffedd3280f3fa4cc030101000000012e3f4683e173b40a20527fe5719633ba070df649983614886e90e45aecf2ac56000000006a47304402200b661b459f6a8e61af42e4e68ba74254988a333d579f98e49ee3cc588237627402205315f09e40b3815c50eaaaf06fd0163b9d8a72221712cac1296223750e4421de4121034d2d6d23fbcb6eefe3e80c47044e36797dcb80d0ac5e96e732ef03c3c550a116ffffffff01e8030000000000001976a91494677c56fa2968644c90a517214338b4139899ce88ac000000000100"

	// node edges:
	const tx0InputIdx0SourceID = "525b6c14c21d3e7065aedbfcc95b8f6d93b16a560f8a5fec9c403fc73d52e119"
	const tx0InputIdx1SourceID = "db8649bc018ab603c00f356d4017db11117261a540c529630f807f5a2586fef4"
	const tx0InputIdx2SourceID = "d4ef5d4b1387ead1c4b46a24ee87f28cb6df6ce3612e230436cee4139fb8d739"

	const tx1InputIdx0SourceID = "8bf77b0d76769f95c1d602df4ed34142741df80de0e9aaec39a423756a2f7b91"
	const tx4InputIdx0SourceID = "abc8ea4ab9af1879b4d8fd438c099bc4ae68edf2cd8b8b57b67aba84ad48ce8c"

	const tx2InputIdx0SourceID = "03cca43f0f28d3edffe30354b28934bc8e881e94ecfa68de2cf899a0a647d37c"

	const tx3InputIdx0SourceID = "d8b235793d3112062eb9b4b404dbf0ca0aa250e952c1213c5bd6c93802c3a312"
	const tx5InputIdx0SourceID = "cf086693db68ef80bb83545e1faeaebd5b11f4e1b34bf5eb78a520b6069e8257"

	tx0, err := sdk.NewTransactionFromHex(tx0Hex)
	require.NoError(t, err)
	require.NotEmpty(t, tx0)

	builder := SourceTransactionBuilder{Tx: tx0}

	// when:
	err = builder.Build(TxQueryResultSlice{
		// left branch:
		{
			SourceTXID: tx0InputIdx0SourceID,
			RawHex:     Ptr(tx1Hex),
		},
		{
			SourceTXID: tx1InputIdx0SourceID,
			RawHex:     Ptr(tx4Hex),
		},
		{
			SourceTXID: tx4InputIdx0SourceID,
			BeefHex:    Ptr(tx6BEEFHex),
		},
		// right branch:
		{
			SourceTXID: tx0InputIdx1SourceID,
			RawHex:     Ptr(tx3Hex),
		},
		{
			SourceTXID: tx3InputIdx0SourceID,
			RawHex:     Ptr(tx5Hex),
		},
		{
			SourceTXID: tx5InputIdx0SourceID,
			BeefHex:    Ptr(tx7BEEFHex),
		},
		// bottom branch:
		{
			SourceTXID: tx0InputIdx2SourceID,
			RawHex:     Ptr(tx2Hex),
		},
		{
			SourceTXID: tx2InputIdx0SourceID,
			BeefHex:    Ptr(tx8BEEFHex),
		},
	})

	// then:
	require.NoError(t, err)
}

func Test_TxGraphBuilding_Case2(t *testing.T) {
	tx1 := fixtures.GivenTX(t).WithInput(10).WithP2PKHOutput(2).TX()
	tx0 := fixtures.GivenTX(t).WithInputFromUTXO(tx1, 0).TX()

	for i, input := range tx0.Inputs {
		verified, err := spv.VerifyScripts(input.SourceTransaction)
		require.NoError(t, err, "input %v", i)
		require.True(t, verified, "input %v", i)
	}
}

func Test_TxGraphBuilding_Case3(t *testing.T) {
	tx1 := fixtures.GivenTX(t).WithInput(10).WithP2PKHOutput(9).TX()
	tx3 := fixtures.GivenTX(t).WithInput(9).WithP2PKHOutput(8).TX()
	tx2 := fixtures.GivenTX(t).WithInput(8).WithP2PKHOutput(7).TX()

	tx0 := fixtures.GivenTX(t).WithInputFromUTXO(tx1, 0).WithInputFromUTXO(tx2, 0).WithInputFromUTXO(tx3, 0).TX()
	for i, input := range tx0.Inputs {
		verified, err := spv.VerifyScripts(input.SourceTransaction)
		require.NoError(t, err, "input %v", i)
		require.True(t, verified, "input %v", i)
	}
}

func Test_TxGraphBuilding_Case6(t *testing.T) {
	// left branch
	tx4 := fixtures.GivenTX(t).WithInput(10).WithP2PKHOutput(2).TX()
	tx1 := fixtures.GivenTX(t).WithInputFromUTXO(tx4, 0).WithP2PKHOutput(1).TX()

	for _, input := range tx1.Inputs {
		input.UnlockingScriptTemplate = fixtures.RecipientInternal.P2PKHUnlockingScriptTemplate()
	}
	require.NoError(t, tx1.Sign())

	// right branch
	tx5 := fixtures.GivenTX(t).WithInput(100).WithP2PKHOutput(20).TX()
	tx3 := fixtures.GivenTX(t).WithInputFromUTXO(tx5, 0).WithP2PKHOutput(10).TX()

	for _, input := range tx3.Inputs {
		input.UnlockingScriptTemplate = fixtures.RecipientInternal.P2PKHUnlockingScriptTemplate()
	}
	require.NoError(t, tx3.Sign())

	// bottom branch
	tx2 := fixtures.GivenTX(t).WithInput(1000).WithP2PKHOutput(200).TX()

	// root
	tx0 := fixtures.GivenTX(t).WithInputFromUTXO(tx1, 0).WithInputFromUTXO(tx2, 0).WithInputFromUTXO(tx3, 0).TX()

	for i, input := range tx0.Inputs {
		verified, err := spv.VerifyScripts(input.SourceTransaction)
		require.NoError(t, err, "input %v", i)
		require.True(t, verified, "input %v", i)
	}
}

func Ptr[T any](value T) *T {
	return &value
}
