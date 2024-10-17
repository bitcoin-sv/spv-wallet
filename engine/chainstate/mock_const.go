package chainstate

import (
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
)

const (
	// Dummy transaction data
	broadcastExample1TxID  = "15d31d00ed7533a83d7ab206115d7642812ec04a2cbae4248365febb82576ff3"
	broadcastExample1TxHex = "0100000001018d7ab1a0f0253120a0cb284e4170b47e5f83f70faaba5b0b55bbeeef624b45010000006b483045022100d5b0dddf76da9088e21cf1277f064dc7832c3da666732f003ee48f2458142e9a02201fe725a1c455b2bd964779391ae105b87730881f211cd299ca36d70d74d715ab412103673dffd80561b87825658f74076da805c238e8c47f25b5d804893c335514d074ffffffff02c4090000000000001976a914777242b335bc7781f43e1b05c60d8c2f2d08b44c88ac962e0000000000001976a91467d93a70ac575e15abb31bc8272a00ab1495d48388ac00000000"
	onChainExample1TxID    = "908c26f8227fa99f1b26f99a19648653a1382fb3b37b03870e9c138894d29b3b"
	onChainExampleArcTxID  = "a11b9e1ee08e264f9add02e4afa40dad3c00a23f250ac04449face095c68fab7"
)

// MockDefaultFee is a mock default fee used for assertions
var MockDefaultFee = &bsv.FeeUnit{
	Satoshis: 1,
	Bytes:    1000,
}
