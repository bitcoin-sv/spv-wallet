package engine

import (
	"context"
	"database/sql"
	"testing"

	"github.com/bitcoin-sv/go-sdk/script"
	trx "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/datastore"
	customTypes "github.com/bitcoin-sv/spv-wallet/engine/datastore/customtypes"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testTxHex             = "020000000165bb8d2733298b2d3b441a871868d6323c5392facf0d3eced3a6c6a17dc84c10000000006a473044022057b101e9a017cdcc333ef66a4a1e78720ae15adf7d1be9c33abec0fe56bc849d022013daa203095522039fadaba99e567ec3cf8615861d3b7258d5399c9f1f4ace8f412103b9c72aebee5636664b519e5f7264c78614f1e57fa4097ae83a3012a967b1c4b9ffffffff03e0930400000000001976a91413473d21dc9e1fb392f05a028b447b165a052d4d88acf9020000000000001976a91455decebedd9a6c2c2d32cf0ee77e2640c3955d3488ac00000000000000000c006a09446f7457616c6c657400000000"
	testTxID              = "1b52eac9d1eb0adf3ce6a56dee1c4768780b8126e288aca65dd1db32f173b853"
	testTxID2             = "104cc87da1c6a6d3ce3e0dcffa92533c32d66818871a443b2d8b2933278dbb65"
	testTx2Hex            = "020000000189fbccca3a5e2bfc8a161bf7f54e8cb5898e296ae8c23b620b89ed570711f931000000006a47304402204e94380ae4d27f8bb9b40dd9944b4fea532d5fe12cf62c1994a6a495c81490f202204aab42f8f1b15259a032e58a3810fbbfd691771b92317f8a12a0da84761a400641210382229c0295e4d63ee54c541eba40be2963f0e80489b7da34e022d513a723181fffffffff0259970400000000001976a914e069bd2e2fe3ea702c40d5e65b491b734c01686788ac00000000000000000c006a09446f7457616c6c657400000000"
	testTxInID            = "9b0495704e23e4b3bef3682c6a5c40abccc32a3e6b7b01ae3295e93a9d3a0482"
	testTxInScriptPubKey  = "76a914e069bd2e2fe3ea702c40d5e65b491b734c01686788ac"
	testTxScriptPubKey1   = "76a91413473d21dc9e1fb392f05a028b447b165a052d4d88ac"
	testTxScriptPubKey2   = "76a91455decebedd9a6c2c2d32cf0ee77e2640c3955d3488ac"
	testSTAStxHex         = "0100000002e0d43ab21510a337ca66a58744c11c1bc9519ef733d54cb4d1824c7e8ed3fde9000000006a47304402203342239754aac17471c1fc7bda7f60685d729a2fdf0db0fbef2185fa16f41956022043926c24a3a516727e910408f74f3cb4fc1b94a2e90d5927d7ae030ef5ff18c7412102f892d89cd0e522a0ff3bac1195a4eeba76c366be3b10f3ffea915fd4d4b2bdf1ffffffffd99b8883f6bf6faf2205488470173fc824cdf9b6445dbcc8490488554ee31c44020000006b483045022100bb2062177404040bceab92125dba929871314998ced25ed13ddf9f71bfc7c522022036429c265e0cb2c4d8b49728726ec7d62fbe81084000a1cde14cee193943f9dc4121032639cbd16258e6b0788a0e17eb899ea0f0c65e44a09d35a19f081769737d7525ffffffff04dc05000000000000fd800676a9140311c6e2114620d68ddfc71519c1a00e0bf9d10b88ac6976aa607f5f7f7c5e7f7c5d7f7c5c7f7c5b7f7c5a7f7c597f7c587f7c577f7c567f7c557f7c547f7c537f7c527f7c517f7c7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7c5f7f7c5e7f7c5d7f7c5c7f7c5b7f7c5a7f7c597f7c587f7c577f7c567f7c557f7c547f7c537f7c527f7c517f7c7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e01007e818b21414136d08c5ed2bf3ba048afe6dcaebafeffffffffffffffffffffffffffffff007d976e7c5296a06394677768827601249301307c7e23022079be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798027e7c7e7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e01417e21038ff83d8cf12121491609c4939dc11c4aa35503508fe432dc5a5c1905608b9218ad547f7701207f01207f7701247f517f7801007e8102fd00a063546752687f7801007e817f727e7b01177f777b557a766471567a577a786354807e7e676d68aa880067765158a569765187645294567a5379587a7e7e78637c8c7c53797e577a7e6878637c8c7c53797e577a7e6878637c8c7c53797e577a7e6878637c8c7c53797e577a7e6878637c8c7c53797e577a7e6867567a6876aa587a7d54807e577a597a5a7a786354807e6f7e7eaa727c7e676d6e7eaa7c687b7eaa587a7d877663516752687c72879b69537a647500687c7b547f77517f7853a0916901247f77517f7c01007e817602fc00a06302fd00a063546752687f7c01007e816854937f77788c6301247f77517f7c01007e817602fc00a06302fd00a063546752687f7c01007e816854937f777852946301247f77517f7c01007e817602fc00a06302fd00a063546752687f7c01007e816854937f77686877517f7c52797d8b9f7c53a09b91697c76638c7c587f77517f7c01007e817602fc00a06302fd00a063546752687f7c01007e81687f777c6876638c7c587f77517f7c01007e817602fc00a06302fd00a063546752687f7c01007e81687f777c6863587f77517f7c01007e817602fc00a06302fd00a063546752687f7c01007e81687f7768587f517f7801007e817602fc00a06302fd00a063546752687f7801007e81727e7b7b687f75537f7c0376a9148801147f775379645579887567726881766968789263556753687a76026c057f7701147f8263517f7c766301007e817f7c6775006877686b537992635379528763547a6b547a6b677c6b567a6b537a7c717c71716868547a587f7c81547a557964936755795187637c686b687c547f7701207f75748c7a7669765880748c7a76567a876457790376a9147e7c7e557967041976a9147c7e0288ac687e7e5579636c766976748c7a9d58807e6c0376a9147e748c7a7e6c7e7e676c766b8263828c007c80517e846864745aa0637c748c7a76697d937b7b58807e56790376a9147e748c7a7e55797e7e6868686c567a5187637500678263828c007c80517e846868647459a0637c748c7a76697d937b7b58807e55790376a9147e748c7a7e55797e7e687459a0637c748c7a76697d937b7b58807e55790376a9147e748c7a7e55797e7e68687c537a9d547963557958807e041976a91455797e0288ac7e7e68aa87726d77776a1492a08b11fc38494853997469a4cb62bfe8aa3d990101066567344970754cde7c2056534e207c2065396664643338653765346338326431623434636435333366373965353163393162316363313434383761353636636133376133313031356232336164346530207c2068747470733a2f2f666972656261736573746f726167652e676f6f676c65617069732e636f6d2f76302f622f6d757369636172746465762f6f2f6e667441737365747325324664363861363331322d633138632d346634392d623332612d3535363933636239306265665f323530783235303f616c743d6d65646961207c2033207c2031207c2073666161647366736466207cdc05000000000000fd800676a9140311c6e2114620d68ddfc71519c1a00e0bf9d10b88ac6976aa607f5f7f7c5e7f7c5d7f7c5c7f7c5b7f7c5a7f7c597f7c587f7c577f7c567f7c557f7c547f7c537f7c527f7c517f7c7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7c5f7f7c5e7f7c5d7f7c5c7f7c5b7f7c5a7f7c597f7c587f7c577f7c567f7c557f7c547f7c537f7c527f7c517f7c7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e01007e818b21414136d08c5ed2bf3ba048afe6dcaebafeffffffffffffffffffffffffffffff007d976e7c5296a06394677768827601249301307c7e23022079be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798027e7c7e7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e01417e21038ff83d8cf12121491609c4939dc11c4aa35503508fe432dc5a5c1905608b9218ad547f7701207f01207f7701247f517f7801007e8102fd00a063546752687f7801007e817f727e7b01177f777b557a766471567a577a786354807e7e676d68aa880067765158a569765187645294567a5379587a7e7e78637c8c7c53797e577a7e6878637c8c7c53797e577a7e6878637c8c7c53797e577a7e6878637c8c7c53797e577a7e6878637c8c7c53797e577a7e6867567a6876aa587a7d54807e577a597a5a7a786354807e6f7e7eaa727c7e676d6e7eaa7c687b7eaa587a7d877663516752687c72879b69537a647500687c7b547f77517f7853a0916901247f77517f7c01007e817602fc00a06302fd00a063546752687f7c01007e816854937f77788c6301247f77517f7c01007e817602fc00a06302fd00a063546752687f7c01007e816854937f777852946301247f77517f7c01007e817602fc00a06302fd00a063546752687f7c01007e816854937f77686877517f7c52797d8b9f7c53a09b91697c76638c7c587f77517f7c01007e817602fc00a06302fd00a063546752687f7c01007e81687f777c6876638c7c587f77517f7c01007e817602fc00a06302fd00a063546752687f7c01007e81687f777c6863587f77517f7c01007e817602fc00a06302fd00a063546752687f7c01007e81687f7768587f517f7801007e817602fc00a06302fd00a063546752687f7801007e81727e7b7b687f75537f7c0376a9148801147f775379645579887567726881766968789263556753687a76026c057f7701147f8263517f7c766301007e817f7c6775006877686b537992635379528763547a6b547a6b677c6b567a6b537a7c717c71716868547a587f7c81547a557964936755795187637c686b687c547f7701207f75748c7a7669765880748c7a76567a876457790376a9147e7c7e557967041976a9147c7e0288ac687e7e5579636c766976748c7a9d58807e6c0376a9147e748c7a7e6c7e7e676c766b8263828c007c80517e846864745aa0637c748c7a76697d937b7b58807e56790376a9147e748c7a7e55797e7e6868686c567a5187637500678263828c007c80517e846868647459a0637c748c7a76697d937b7b58807e55790376a9147e748c7a7e55797e7e687459a0637c748c7a76697d937b7b58807e55790376a9147e748c7a7e55797e7e68687c537a9d547963557958807e041976a91455797e0288ac7e7e68aa87726d77776a1492a08b11fc38494853997469a4cb62bfe8aa3d990101066567344970754cde7c2056534e207c2065396664643338653765346338326431623434636435333366373965353163393162316363313434383761353636636133376133313031356232336164346530207c2068747470733a2f2f666972656261736573746f726167652e676f6f676c65617069732e636f6d2f76302f622f6d757369636172746465762f6f2f6e667441737365747325324664363861363331322d633138632d346634392d623332612d3535363933636239306265665f323530783235303f616c743d6d65646961207c2033207c2032207c2073666161647366736466207cdc05000000000000fd800676a9140311c6e2114620d68ddfc71519c1a00e0bf9d10b88ac6976aa607f5f7f7c5e7f7c5d7f7c5c7f7c5b7f7c5a7f7c597f7c587f7c577f7c567f7c557f7c547f7c537f7c527f7c517f7c7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7c5f7f7c5e7f7c5d7f7c5c7f7c5b7f7c5a7f7c597f7c587f7c577f7c567f7c557f7c547f7c537f7c527f7c517f7c7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e01007e818b21414136d08c5ed2bf3ba048afe6dcaebafeffffffffffffffffffffffffffffff007d976e7c5296a06394677768827601249301307c7e23022079be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798027e7c7e7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e01417e21038ff83d8cf12121491609c4939dc11c4aa35503508fe432dc5a5c1905608b9218ad547f7701207f01207f7701247f517f7801007e8102fd00a063546752687f7801007e817f727e7b01177f777b557a766471567a577a786354807e7e676d68aa880067765158a569765187645294567a5379587a7e7e78637c8c7c53797e577a7e6878637c8c7c53797e577a7e6878637c8c7c53797e577a7e6878637c8c7c53797e577a7e6878637c8c7c53797e577a7e6867567a6876aa587a7d54807e577a597a5a7a786354807e6f7e7eaa727c7e676d6e7eaa7c687b7eaa587a7d877663516752687c72879b69537a647500687c7b547f77517f7853a0916901247f77517f7c01007e817602fc00a06302fd00a063546752687f7c01007e816854937f77788c6301247f77517f7c01007e817602fc00a06302fd00a063546752687f7c01007e816854937f777852946301247f77517f7c01007e817602fc00a06302fd00a063546752687f7c01007e816854937f77686877517f7c52797d8b9f7c53a09b91697c76638c7c587f77517f7c01007e817602fc00a06302fd00a063546752687f7c01007e81687f777c6876638c7c587f77517f7c01007e817602fc00a06302fd00a063546752687f7c01007e81687f777c6863587f77517f7c01007e817602fc00a06302fd00a063546752687f7c01007e81687f7768587f517f7801007e817602fc00a06302fd00a063546752687f7801007e81727e7b7b687f75537f7c0376a9148801147f775379645579887567726881766968789263556753687a76026c057f7701147f8263517f7c766301007e817f7c6775006877686b537992635379528763547a6b547a6b677c6b567a6b537a7c717c71716868547a587f7c81547a557964936755795187637c686b687c547f7701207f75748c7a7669765880748c7a76567a876457790376a9147e7c7e557967041976a9147c7e0288ac687e7e5579636c766976748c7a9d58807e6c0376a9147e748c7a7e6c7e7e676c766b8263828c007c80517e846864745aa0637c748c7a76697d937b7b58807e56790376a9147e748c7a7e55797e7e6868686c567a5187637500678263828c007c80517e846868647459a0637c748c7a76697d937b7b58807e55790376a9147e748c7a7e55797e7e687459a0637c748c7a76697d937b7b58807e55790376a9147e748c7a7e55797e7e68687c537a9d547963557958807e041976a91455797e0288ac7e7e68aa87726d77776a1492a08b11fc38494853997469a4cb62bfe8aa3d990101066567344970754cde7c2056534e207c2065396664643338653765346338326431623434636435333366373965353163393162316363313434383761353636636133376133313031356232336164346530207c2068747470733a2f2f666972656261736573746f726167652e676f6f676c65617069732e636f6d2f76302f622f6d757369636172746465762f6f2f6e667441737365747325324664363861363331322d633138632d346634392d623332612d3535363933636239306265665f323530783235303f616c743d6d65646961207c2033207c2033207c2073666161647366736466207cad180000000000001976a9147190407e487fd53c7504031956c9b995bb8dfd3988ac00000000"
	testSTAStxID          = "76a4f090140242a34c41fc2ac1936b140dc0efad65b8a61fed32227c13ff11f4"
	testSTASLockingScript = "76a9140311c6e2114620d68ddfc71519c1a00e0bf9d10b88ac6976aa607f5f7f7c5e7f7c5d7f7c5c7f7c5b7f7c5a7f7c597f7c587f7c577f7c567f7c557f7c547f7c537f7c527f7c517f7c7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7c5f7f7c5e7f7c5d7f7c5c7f7c5b7f7c5a7f7c597f7c587f7c577f7c567f7c557f7c547f7c537f7c527f7c517f7c7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e01007e818b21414136d08c5ed2bf3ba048afe6dcaebafeffffffffffffffffffffffffffffff007d976e7c5296a06394677768827601249301307c7e23022079be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798027e7c7e7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e01417e21038ff83d8cf12121491609c4939dc11c4aa35503508fe432dc5a5c1905608b9218ad547f7701207f01207f7701247f517f7801007e8102fd00a063546752687f7801007e817f727e7b01177f777b557a766471567a577a786354807e7e676d68aa880067765158a569765187645294567a5379587a7e7e78637c8c7c53797e577a7e6878637c8c7c53797e577a7e6878637c8c7c53797e577a7e6878637c8c7c53797e577a7e6878637c8c7c53797e577a7e6867567a6876aa587a7d54807e577a597a5a7a786354807e6f7e7eaa727c7e676d6e7eaa7c687b7eaa587a7d877663516752687c72879b69537a647500687c7b547f77517f7853a0916901247f77517f7c01007e817602fc00a06302fd00a063546752687f7c01007e816854937f77788c6301247f77517f7c01007e817602fc00a06302fd00a063546752687f7c01007e816854937f777852946301247f77517f7c01007e817602fc00a06302fd00a063546752687f7c01007e816854937f77686877517f7c52797d8b9f7c53a09b91697c76638c7c587f77517f7c01007e817602fc00a06302fd00a063546752687f7c01007e81687f777c6876638c7c587f77517f7c01007e817602fc00a06302fd00a063546752687f7c01007e81687f777c6863587f77517f7c01007e817602fc00a06302fd00a063546752687f7c01007e81687f7768587f517f7801007e817602fc00a06302fd00a063546752687f7801007e81727e7b7b687f75537f7c0376a9148801147f775379645579887567726881766968789263556753687a76026c057f7701147f8263517f7c766301007e817f7c6775006877686b537992635379528763547a6b547a6b677c6b567a6b537a7c717c71716868547a587f7c81547a557964936755795187637c686b687c547f7701207f75748c7a7669765880748c7a76567a876457790376a9147e7c7e557967041976a9147c7e0288ac687e7e5579636c766976748c7a9d58807e6c0376a9147e748c7a7e6c7e7e676c766b8263828c007c80517e846864745aa0637c748c7a76697d937b7b58807e56790376a9147e748c7a7e55797e7e6868686c567a5187637500678263828c007c80517e846868647459a0637c748c7a76697d937b7b58807e55790376a9147e748c7a7e55797e7e687459a0637c748c7a76697d937b7b58807e55790376a9147e748c7a7e55797e7e68687c537a9d547963557958807e041976a91455797e0288ac7e7e68aa87726d77776a1492a08b11fc38494853997469a4cb62bfe8aa3d990101066567344970754cde7c2056534e207c2065396664643338653765346338326431623434636435333366373965353163393162316363313434383761353636636133376133313031356232336164346530207c2068747470733a2f2f666972656261736573746f726167652e676f6f676c65617069732e636f6d2f76302f622f6d757369636172746465762f6f2f6e667441737365747325324664363861363331322d633138632d346634392d623332612d3535363933636239306265665f323530783235303f616c743d6d65646961207c2033207c2031207c2073666161647366736466207c"
	testSTASScriptPubKey  = "76a9140311c6e2114620d68ddfc71519c1a00e0bf9d10b88ac"
	testSTAStx2Hex        = "0100000002f411ff137c2232ed1fa6b865adefc00d146b93c12afc414ca342021490f0a47600000000fde00702dc0514df5d8fa1cb1f668212dc8ca438bff19997df3e0302eb01147190407e487fd53c7504031956c9b995bb8dfd390020805de4ba94c12e8620fa6deee1a5334fd1ee6131b116d4c9cdf4ba5f0106bee0004d1f0701000000fa143d763f2c79e84e885a36e4b961a7c6816541f3e060e7efdc88276d19d6c9752adad0a7b9ceca853768aebb6965eca126a62965f698a0c1bc43d83db632adf411ff137c2232ed1fa6b865adefc00d146b93c12afc414ca342021490f0a47600000000fd800676a9140311c6e2114620d68ddfc71519c1a00e0bf9d10b88ac6976aa607f5f7f7c5e7f7c5d7f7c5c7f7c5b7f7c5a7f7c597f7c587f7c577f7c567f7c557f7c547f7c537f7c527f7c517f7c7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7c5f7f7c5e7f7c5d7f7c5c7f7c5b7f7c5a7f7c597f7c587f7c577f7c567f7c557f7c547f7c537f7c527f7c517f7c7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e01007e818b21414136d08c5ed2bf3ba048afe6dcaebafeffffffffffffffffffffffffffffff007d976e7c5296a06394677768827601249301307c7e23022079be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798027e7c7e7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e01417e21038ff83d8cf12121491609c4939dc11c4aa35503508fe432dc5a5c1905608b9218ad547f7701207f01207f7701247f517f7801007e8102fd00a063546752687f7801007e817f727e7b01177f777b557a766471567a577a786354807e7e676d68aa880067765158a569765187645294567a5379587a7e7e78637c8c7c53797e577a7e6878637c8c7c53797e577a7e6878637c8c7c53797e577a7e6878637c8c7c53797e577a7e6878637c8c7c53797e577a7e6867567a6876aa587a7d54807e577a597a5a7a786354807e6f7e7eaa727c7e676d6e7eaa7c687b7eaa587a7d877663516752687c72879b69537a647500687c7b547f77517f7853a0916901247f77517f7c01007e817602fc00a06302fd00a063546752687f7c01007e816854937f77788c6301247f77517f7c01007e817602fc00a06302fd00a063546752687f7c01007e816854937f777852946301247f77517f7c01007e817602fc00a06302fd00a063546752687f7c01007e816854937f77686877517f7c52797d8b9f7c53a09b91697c76638c7c587f77517f7c01007e817602fc00a06302fd00a063546752687f7c01007e81687f777c6876638c7c587f77517f7c01007e817602fc00a06302fd00a063546752687f7c01007e81687f777c6863587f77517f7c01007e817602fc00a06302fd00a063546752687f7c01007e81687f7768587f517f7801007e817602fc00a06302fd00a063546752687f7801007e81727e7b7b687f75537f7c0376a9148801147f775379645579887567726881766968789263556753687a76026c057f7701147f8263517f7c766301007e817f7c6775006877686b537992635379528763547a6b547a6b677c6b567a6b537a7c717c71716868547a587f7c81547a557964936755795187637c686b687c547f7701207f75748c7a7669765880748c7a76567a876457790376a9147e7c7e557967041976a9147c7e0288ac687e7e5579636c766976748c7a9d58807e6c0376a9147e748c7a7e6c7e7e676c766b8263828c007c80517e846864745aa0637c748c7a76697d937b7b58807e56790376a9147e748c7a7e55797e7e6868686c567a5187637500678263828c007c80517e846868647459a0637c748c7a76697d937b7b58807e55790376a9147e748c7a7e55797e7e687459a0637c748c7a76697d937b7b58807e55790376a9147e748c7a7e55797e7e68687c537a9d547963557958807e041976a91455797e0288ac7e7e68aa87726d77776a1492a08b11fc38494853997469a4cb62bfe8aa3d990101066567344970754cde7c2056534e207c2065396664643338653765346338326431623434636435333366373965353163393162316363313434383761353636636133376133313031356232336164346530207c2068747470733a2f2f666972656261736573746f726167652e676f6f676c65617069732e636f6d2f76302f622f6d757369636172746465762f6f2f6e667441737365747325324664363861363331322d633138632d346634392d623332612d3535363933636239306265665f323530783235303f616c743d6d65646961207c2033207c2031207c2073666161647366736466207cdc05000000000000ffffffff62b663e4d26475d1a6271fcbc5c7739f7e48efd1c76a8156dd0f2894a6f1fd140000000041000000483045022100a38cb7dee256592cac8c9b11f4c597189114041f761e94fdb464be0137f1c97502207927fc1bbca3e4f8ecaad056b69be1705fdb5ef77d6339c9bd13357fddae1786412102a598951f485328b1eda162237b603f2af9d4d9127447f9121d662cba6b2025daffffffff805de4ba94c12e8620fa6deee1a5334fd1ee6131b116d4c9cdf4ba5f0106bee0000000006b483045022100bc8a366b0f41818805bd9b84f0fdc8d2c1642cb2a7c5b06e49e5c47265717351022041072c5ec390fa6aa81681c14ff69af4bb3f9fad78cac07a025931c9b373176a4121032639cbd16258e6b0788a0e17eb899ea0f0c65e44a09d35a19f081769737d7525ffffffff02dc05000000000000fd800676a914df5d8fa1cb1f668212dc8ca438bff19997df3e0388ac6976aa607f5f7f7c5e7f7c5d7f7c5c7f7c5b7f7c5a7f7c597f7c587f7c577f7c567f7c557f7c547f7c537f7c527f7c517f7c7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7c5f7f7c5e7f7c5d7f7c5c7f7c5b7f7c5a7f7c597f7c587f7c577f7c567f7c557f7c547f7c537f7c527f7c517f7c7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e01007e818b21414136d08c5ed2bf3ba048afe6dcaebafeffffffffffffffffffffffffffffff007d976e7c5296a06394677768827601249301307c7e23022079be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798027e7c7e7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e01417e21038ff83d8cf12121491609c4939dc11c4aa35503508fe432dc5a5c1905608b9218ad547f7701207f01207f7701247f517f7801007e8102fd00a063546752687f7801007e817f727e7b01177f777b557a766471567a577a786354807e7e676d68aa880067765158a569765187645294567a5379587a7e7e78637c8c7c53797e577a7e6878637c8c7c53797e577a7e6878637c8c7c53797e577a7e6878637c8c7c53797e577a7e6878637c8c7c53797e577a7e6867567a6876aa587a7d54807e577a597a5a7a786354807e6f7e7eaa727c7e676d6e7eaa7c687b7eaa587a7d877663516752687c72879b69537a647500687c7b547f77517f7853a0916901247f77517f7c01007e817602fc00a06302fd00a063546752687f7c01007e816854937f77788c6301247f77517f7c01007e817602fc00a06302fd00a063546752687f7c01007e816854937f777852946301247f77517f7c01007e817602fc00a06302fd00a063546752687f7c01007e816854937f77686877517f7c52797d8b9f7c53a09b91697c76638c7c587f77517f7c01007e817602fc00a06302fd00a063546752687f7c01007e81687f777c6876638c7c587f77517f7c01007e817602fc00a06302fd00a063546752687f7c01007e81687f777c6863587f77517f7c01007e817602fc00a06302fd00a063546752687f7c01007e81687f7768587f517f7801007e817602fc00a06302fd00a063546752687f7801007e81727e7b7b687f75537f7c0376a9148801147f775379645579887567726881766968789263556753687a76026c057f7701147f8263517f7c766301007e817f7c6775006877686b537992635379528763547a6b547a6b677c6b567a6b537a7c717c71716868547a587f7c81547a557964936755795187637c686b687c547f7701207f75748c7a7669765880748c7a76567a876457790376a9147e7c7e557967041976a9147c7e0288ac687e7e5579636c766976748c7a9d58807e6c0376a9147e748c7a7e6c7e7e676c766b8263828c007c80517e846864745aa0637c748c7a76697d937b7b58807e56790376a9147e748c7a7e55797e7e6868686c567a5187637500678263828c007c80517e846868647459a0637c748c7a76697d937b7b58807e55790376a9147e748c7a7e55797e7e687459a0637c748c7a76697d937b7b58807e55790376a9147e748c7a7e55797e7e68687c537a9d547963557958807e041976a91455797e0288ac7e7e68aa87726d77776a1492a08b11fc38494853997469a4cb62bfe8aa3d990101066567344970754cde7c2056534e207c2065396664643338653765346338326431623434636435333366373965353163393162316363313434383761353636636133376133313031356232336164346530207c2068747470733a2f2f666972656261736573746f726167652e676f6f676c65617069732e636f6d2f76302f622f6d757369636172746465762f6f2f6e667441737365747325324664363861363331322d633138632d346634392d623332612d3535363933636239306265665f323530783235303f616c743d6d65646961207c2033207c2031207c2073666161647366736466207ceb010000000000001976a9147190407e487fd53c7504031956c9b995bb8dfd3988ac00000000"
	testSTAStx2ID         = "c7abebcc6a49a28e509687afd7cda8d147cd58fa5a5f7e45fa3e04a64f39973a"
)

type transactionServiceMock struct {
	destinations map[string]*Destination
	utxos        map[string]map[uint32]*Utxo
}

func (x transactionServiceMock) getDestinationByLockingScript(_ context.Context, lockingScript string, _ ...ModelOps) (*Destination, error) {
	return x.destinations[lockingScript], nil
}

func (x transactionServiceMock) getUtxo(_ context.Context, txID string, index uint32, _ ...ModelOps) (*Utxo, error) {
	return x.utxos[txID][index], nil
}

func TestTransaction_newTransaction(t *testing.T) {
	t.Parallel()

	t.Run("New transaction model", func(t *testing.T) {
		transaction, err := txFromHex(testTxHex, New())
		require.NoError(t, err)
		require.NotNil(t, transaction)

		assert.IsType(t, Transaction{}, *transaction)
		assert.Equal(t, ModelTransaction.String(), transaction.GetModelName())
		assert.Equal(t, testTxID, transaction.ID)
		assert.Equal(t, testTxID, transaction.GetID())
		assert.Equal(t, true, transaction.IsNew())
	})

	t.Run("New transaction model - no hex, no options", func(t *testing.T) {
		transaction := emptyTx()
		require.NotNil(t, transaction)
		assert.IsType(t, Transaction{}, *transaction)
		assert.Equal(t, ModelTransaction.String(), transaction.GetModelName())
		assert.Equal(t, "", transaction.ID)
		assert.Equal(t, "", transaction.GetID())
		assert.Equal(t, false, transaction.IsNew())
	})
}

func TestTransaction_newTransactionWithDraftID(t *testing.T) {
	t.Parallel()

	t.Run("New transaction model", func(t *testing.T) {
		transaction, err := newTransactionWithDraftID(testTxHex, testDraftID, New())
		require.NoError(t, err)
		require.NotNil(t, transaction)

		assert.IsType(t, Transaction{}, *transaction)
		assert.Equal(t, ModelTransaction.String(), transaction.GetModelName())
		assert.Equal(t, testTxID, transaction.ID)
		assert.Equal(t, testDraftID, transaction.DraftID)
		assert.Equal(t, testTxID, transaction.GetID())
		assert.Equal(t, true, transaction.IsNew())
	})

	t.Run("New transaction model - no hex - return error", func(t *testing.T) {
		transaction, err := newTransactionWithDraftID("", "")
		require.Nil(t, transaction)
		require.Error(t, err)
	})
}

func TestTransaction_getTransactionByID(t *testing.T) {
	t.Run("not found", func(t *testing.T) {
		ctx, client, deferMe := CreateTestSQLiteClient(t, false, false, withTaskManagerMockup())
		defer deferMe()
		transaction, err := getTransactionByID(ctx, testXPubID, testTxID, client.DefaultModelOptions()...)
		require.NoError(t, err)
		assert.Nil(t, transaction)
	})

	t.Run("found tx", func(t *testing.T) {
		ctx, client, deferMe := CreateTestSQLiteClient(t, false, true)
		defer deferMe()

		opts := client.DefaultModelOptions()
		tx, err := txFromHex(testTxHex, append(opts, New())...)
		require.NoError(t, err)
		txErr := tx.Save(ctx)
		require.NoError(t, txErr)

		transaction, err := getTransactionByID(ctx, testXPubID, testTxID, client.DefaultModelOptions()...)
		require.NoError(t, err)
		assert.NotNil(t, transaction)
		assert.Equal(t, testTxID, transaction.ID)
		assert.Equal(t, testTxHex, transaction.Hex)
		assert.Nil(t, transaction.XpubInIDs)
		assert.Nil(t, transaction.XpubOutIDs)
	})
}

func TestTransaction_getTransactionsByXpubID(t *testing.T) {
	t.Run("tx not found", func(t *testing.T) {
		ctx, client, deferMe := CreateTestSQLiteClient(t, false, false, withTaskManagerMockup())
		defer deferMe()
		transactions, err := getTransactionsByXpubID(ctx, testXPub, nil, nil, nil, client.DefaultModelOptions()...)
		require.NoError(t, err)
		assert.Nil(t, transactions)
	})

	t.Run("tx found", func(t *testing.T) {
		ctx, client, _ := CreateTestSQLiteClient(t, true, true)
		opts := client.DefaultModelOptions()
		tx, err := txFromHex(testTxHex, append(opts, New())...)
		require.NoError(t, err)

		tx.XpubInIDs = append(tx.XpubInIDs, testXPubID)
		txErr := tx.Save(ctx)
		require.NoError(t, txErr)

		transactions, err := getTransactionsByXpubID(ctx, testXPubID, nil, nil, nil, opts...)
		require.NoError(t, err)
		require.NotNil(t, transactions)
		require.Len(t, transactions, 1)
		assert.Equal(t, testTxID, transactions[0].ID)
		assert.Equal(t, testTxHex, transactions[0].Hex)
		assert.Equal(t, testXPubID, transactions[0].XpubInIDs[0])
		assert.Nil(t, transactions[0].XpubOutIDs)
	})
}

func TestTransaction_UpdateTransactionMetadata(t *testing.T) {
	t.Run("tx without meta data", func(t *testing.T) {
		_, client, _ := CreateTestSQLiteClient(t, true, true)
		opts := client.DefaultModelOptions()
		tx, err := txFromHex(testTxHex, append(opts, New())...)
		require.NoError(t, err)

		assert.Nil(t, tx.XpubMetadata)

		metadata := Metadata{
			"test-key": "test-value",
		}
		err = tx.UpdateTransactionMetadata(testXPubID, metadata)
		require.NoError(t, err)
		assert.Equal(t, XpubMetadata{testXPubID: metadata}, tx.XpubMetadata)

		addMetadata := Metadata{
			"test-key-2": "test-value-2",
			"test-key-3": "test-value-3",
		}
		err = tx.UpdateTransactionMetadata(testXPubID, addMetadata)
		require.NoError(t, err)
		assert.Equal(t, XpubMetadata{testXPubID: Metadata{
			"test-key":   "test-value",
			"test-key-2": "test-value-2",
			"test-key-3": "test-value-3",
		}}, tx.XpubMetadata)

		editMetadata := Metadata{
			"test-key-2": nil,
			"test-key-3": "test-value-3333",
			"test-key-4": "test-value-4",
		}
		err = tx.UpdateTransactionMetadata(testXPubID, editMetadata)
		require.NoError(t, err)
		assert.Equal(t, XpubMetadata{testXPubID: Metadata{
			"test-key":   "test-value",
			"test-key-3": "test-value-3333",
			"test-key-4": "test-value-4",
		}}, tx.XpubMetadata)
	})
}

func TestTransaction_BeforeCreating(t *testing.T) {
	// t.Parallel()

	t.Run("incorrect transaction hex", func(t *testing.T) {
		_, err := txFromHex("test")
		assert.Error(t, err)
	})

	t.Run("no transaction hex", func(t *testing.T) {
		transaction := emptyTx()

		opts := DefaultClientOpts()
		client, _ := NewClient(context.Background(), opts...)
		transaction.client = client

		err := transaction.BeforeCreating(context.Background())
		assert.Error(t, err)
		assert.ErrorIs(t, ErrMissingFieldHex, err)
	})
}

func (ts *EmbeddedDBTestSuite) TestTransaction_BeforeCreating() {
	ts.T().Run("[sqlite] [in-memory] - valid transaction", func(t *testing.T) {
		tc := ts.genericDBClient(t, datastore.SQLite, true)
		defer tc.Close(tc.ctx)

		transaction, err := txFromHex(testTxHex, append(tc.client.DefaultModelOptions(), New())...)
		require.NoError(t, err)
		require.NotNil(t, transaction)

		err = transaction.BeforeCreating(tc.ctx)
		require.NoError(t, err)
	})
}

func TestTransaction_GetID(t *testing.T) {
	t.Parallel()

	t.Run("no id", func(t *testing.T) {
		transaction := emptyTx()
		require.NotNil(t, transaction)
		assert.Equal(t, "", transaction.GetID())
	})

	t.Run("valid id", func(t *testing.T) {
		transaction, err := txFromHex(testTxHex)
		require.NoError(t, err)
		require.NotNil(t, transaction)

		assert.Equal(t, testTxID, transaction.GetID())
	})
}

func TestTransaction_GetModelName(t *testing.T) {
	t.Parallel()

	t.Run("model name", func(t *testing.T) {
		transaction := emptyTx()
		assert.Equal(t, ModelTransaction.String(), transaction.GetModelName())
	})
}

func (ts *EmbeddedDBTestSuite) TestTransaction_processOutputs() {
	ts.T().Run("no outputs", func(t *testing.T) {
		tc := ts.genericDBClient(t, datastore.SQLite, true)
		defer tc.Close(tc.ctx)

		transaction, err := txFromHex(testTxHex, append(tc.client.DefaultModelOptions(), New())...)
		require.NoError(t, err)
		require.NotNil(t, transaction)

		transaction.transactionService = transactionServiceMock{}

		ctx := context.Background()
		err = transaction._processOutputs(ctx)
		require.NoError(t, err)
		assert.Nil(t, transaction.utxos)
		assert.Nil(t, transaction.XpubOutIDs)
	})

	ts.T().Run("no outputs", func(t *testing.T) {
		tc := ts.genericDBClient(t, datastore.SQLite, true)
		defer tc.Close(tc.ctx)

		transaction, err := txFromHex(testTxHex, append(tc.client.DefaultModelOptions(), New())...)
		require.NoError(t, err)
		require.NotNil(t, transaction)

		transaction.transactionService = transactionServiceMock{
			destinations: map[string]*Destination{
				"76a91413473d21dc9e1fb392f05a028b447b165a052d4d88ac": {
					Model:  Model{name: ModelDestination},
					XpubID: "test-xpub-id",
				},
			},
		}

		ctx := context.Background()
		err = transaction._processOutputs(ctx)
		require.NoError(t, err)
		require.NotNil(t, transaction.utxos)
		assert.IsType(t, Utxo{}, transaction.utxos[0])
		assert.Equal(t, testTxID, transaction.utxos[0].TransactionID)
		assert.Equal(t, "test-xpub-id", transaction.utxos[0].XpubID)
		assert.Equal(t, "test-xpub-id", transaction.XpubOutIDs[0])

		childModels := transaction.ChildModels()
		assert.Len(t, childModels, 1)
		assert.Equal(t, "utxo", childModels[0].Name())
	})

	ts.T().Run("STAS token", func(t *testing.T) {
		tc := ts.genericDBClient(t, datastore.SQLite, true)
		defer tc.Close(tc.ctx)

		transaction, err := txFromHex(testSTAStxHex, append(tc.client.DefaultModelOptions(), New())...)
		require.NoError(t, err)
		require.NotNil(t, transaction)

		transaction.transactionService = transactionServiceMock{
			destinations: map[string]*Destination{
				"76a9140311c6e2114620d68ddfc71519c1a00e0bf9d10b88ac": {
					Model:  Model{name: ModelDestination},
					XpubID: "test-xpub-id",
				},
			},
		}

		ctx := context.Background()
		err = transaction._processOutputs(ctx)
		require.NoError(t, err)
		require.NotNil(t, transaction.utxos)
		assert.IsType(t, Utxo{}, transaction.utxos[0])
		assert.Equal(t, testSTAStxID, transaction.utxos[0].TransactionID)
		assert.Equal(t, "test-xpub-id", transaction.utxos[0].XpubID)
		assert.Equal(t, "test-xpub-id", transaction.XpubOutIDs[0])

		childModels := transaction.ChildModels()
		assert.Len(t, childModels, 3)
		assert.Equal(t, "utxo", childModels[0].Name())
		for _, childModel := range childModels {
			utxo := childModel.(*Utxo)
			err = utxo.BeforeCreating(ctx)
			require.NoError(t, err)
			assert.Equal(t, utils.ScriptTypeTokenStas, utxo.Type)
		}
	})
}

func TestTransaction_processInputs(t *testing.T) {
	// t.Parallel()

	t.Run("no utxo", func(t *testing.T) {
		transaction, err := txFromHex(testTxHex, New())
		require.NoError(t, err)
		require.NotNil(t, transaction)

		transaction.transactionService = transactionServiceMock{}

		ctx := context.Background()
		err = transaction._processInputs(ctx)
		require.NoError(t, err)
		assert.Nil(t, transaction.utxos)
		assert.Nil(t, transaction.XpubInIDs)
	})

	t.Run("got utxo", func(t *testing.T) {
		transaction, err := txFromHex(testTxHex, New())
		require.NoError(t, err)
		require.NotNil(t, transaction)

		transaction.draftTransaction = &DraftTransaction{
			TransactionBase: TransactionBase{ID: testDraftID},
		}
		transaction.transactionService = transactionServiceMock{
			utxos: map[string]map[uint32]*Utxo{
				testTxID2: {
					uint32(0): {
						Model: Model{name: ModelUtxo},
						UtxoPointer: UtxoPointer{
							OutputIndex:   0,
							TransactionID: testTxID2,
						},
						XpubID: "test-xpub-id",
						DraftID: customTypes.NullString{NullString: sql.NullString{
							Valid:  true,
							String: testDraftID,
						}},
					},
				},
			},
		}

		ctx := context.Background()
		err = transaction._processInputs(ctx)
		require.NoError(t, err)
		require.NotNil(t, transaction.utxos)
		assert.IsType(t, Utxo{}, transaction.utxos[0])
		assert.Equal(t, testTxID2, transaction.utxos[0].TransactionID)
		assert.True(t, transaction.utxos[0].SpendingTxID.Valid)
		assert.Equal(t, testTxID, transaction.utxos[0].SpendingTxID.String)
		assert.Equal(t, "test-xpub-id", transaction.utxos[0].XpubID)
		assert.Equal(t, "test-xpub-id", transaction.XpubInIDs[0])

		childModels := transaction.ChildModels()
		assert.Len(t, childModels, 1)
		assert.Equal(t, "utxo", childModels[0].Name())
	})

	t.Run("spent utxo", func(t *testing.T) {
		transaction, err := txFromHex(testTxHex, New())
		require.NoError(t, err)
		require.NotNil(t, transaction)

		transaction.draftTransaction = &DraftTransaction{}
		transaction.transactionService = transactionServiceMock{
			utxos: map[string]map[uint32]*Utxo{
				testTxID2: {
					uint32(0): {
						Model: Model{name: ModelUtxo},
						UtxoPointer: UtxoPointer{
							OutputIndex:   0,
							TransactionID: testTxID2,
						},
						XpubID: "test-xpub-id",
						SpendingTxID: customTypes.NullString{NullString: sql.NullString{
							Valid:  true,
							String: testTxID,
						}},
						DraftID: customTypes.NullString{NullString: sql.NullString{
							Valid:  true,
							String: testDraftID2,
						}},
					},
				},
			},
		}

		ctx := context.Background()
		err = transaction._processInputs(ctx)
		require.ErrorIs(t, err, spverrors.ErrUtxoAlreadySpent)
	})

	t.Run("not reserved utxo", func(t *testing.T) {
		transaction, err := txFromHex(testTxHex, New())
		require.NoError(t, err)
		require.NotNil(t, transaction)

		transaction.draftTransaction = &DraftTransaction{
			TransactionBase: TransactionBase{ID: testDraftID},
		}
		transaction.transactionService = transactionServiceMock{
			utxos: map[string]map[uint32]*Utxo{
				testTxID2: {
					uint32(0): {
						Model: Model{name: ModelUtxo},
						UtxoPointer: UtxoPointer{
							OutputIndex:   0,
							TransactionID: testTxID2,
						},
						XpubID: "test-xpub-id",
					},
				},
			},
		}

		ctx := context.Background()
		err = transaction._processInputs(ctx)
		require.ErrorIs(t, err, spverrors.ErrUtxoNotReserved)
	})

	t.Run("incorrect reservation ID of utxo", func(t *testing.T) {
		transaction, err := txFromHex(testTxHex, New())
		require.NoError(t, err)
		require.NotNil(t, transaction)

		transaction.draftTransaction = &DraftTransaction{
			TransactionBase: TransactionBase{ID: testDraftID},
		}
		transaction.transactionService = transactionServiceMock{
			utxos: map[string]map[uint32]*Utxo{
				testTxID2: {
					uint32(0): {
						Model: Model{name: ModelUtxo},
						UtxoPointer: UtxoPointer{
							OutputIndex:   0,
							TransactionID: testTxID2,
						},
						XpubID: "test-xpub-id",
						DraftID: customTypes.NullString{NullString: sql.NullString{
							Valid:  true,
							String: testDraftID2,
						}},
					},
				},
			},
		}

		ctx := context.Background()
		err = transaction._processInputs(ctx)
		require.ErrorIs(t, err, spverrors.ErrDraftIDMismatch)
	})

	t.Run("inputUtxoChecksOff", func(t *testing.T) {
		ctx, client, deferMe := CreateTestSQLiteClient(t, false, false, withTaskManagerMockup(), WithIUCDisabled())
		defer deferMe()

		transaction, err := txFromHex(testTxHex, append(client.DefaultModelOptions(), New())...)
		require.NoError(t, err)
		require.NotNil(t, transaction)

		transaction.draftTransaction = &DraftTransaction{
			TransactionBase: TransactionBase{ID: testDraftID},
		}
		transaction.transactionService = transactionServiceMock{
			utxos: map[string]map[uint32]*Utxo{
				testTxID2: {
					uint32(0): {
						Model: Model{name: ModelUtxo},
						UtxoPointer: UtxoPointer{
							OutputIndex:   0,
							TransactionID: testTxID2,
						},
						XpubID: "test-xpub-id",
					},
				},
			},
		}

		err = transaction._processInputs(ctx)
		require.NoError(t, err)
	})

	t.Run("STAS token input", func(t *testing.T) {
		ctx := context.Background()

		transaction, err := txFromHex(testSTAStx2Hex, New())
		require.NoError(t, err)
		require.NotNil(t, transaction)

		transaction.draftTransaction = &DraftTransaction{
			TransactionBase: TransactionBase{ID: testDraftID},
		}
		transaction.transactionService = transactionServiceMock{
			utxos: map[string]map[uint32]*Utxo{
				testSTAStxID: {
					uint32(0): {
						Model: Model{name: ModelUtxo},
						UtxoPointer: UtxoPointer{
							OutputIndex:   0,
							TransactionID: testSTAStxID,
						},
						ScriptPubKey: testSTASLockingScript,
						XpubID:       "test-xpub-id",
						DraftID: customTypes.NullString{NullString: sql.NullString{
							Valid:  true,
							String: testDraftID,
						}},
					},
				},
			},
		}

		err = transaction._processInputs(ctx)
		require.NoError(t, err)
		require.NotNil(t, transaction.utxos)
		assert.IsType(t, Utxo{}, transaction.utxos[0])
		assert.Equal(t, testSTAStxID, transaction.utxos[0].TransactionID)
		assert.True(t, transaction.utxos[0].SpendingTxID.Valid)
		assert.Equal(t, testSTAStx2ID, transaction.utxos[0].SpendingTxID.String)
		assert.Equal(t, "test-xpub-id", transaction.utxos[0].XpubID)
		assert.Equal(t, "test-xpub-id", transaction.XpubInIDs[0])

		childModels := transaction.ChildModels()
		assert.Len(t, childModels, 1)
		assert.Equal(t, "utxo", childModels[0].Name())
	})
}

func TestTransaction_Display(t *testing.T) {
	t.Run("display without xpub data", func(t *testing.T) {
		tx := Transaction{
			Model:  Model{},
			XPubID: testXPubID,
		}

		displayTx := tx.Display().(*Transaction)
		assert.Nil(t, displayTx.Metadata)
		assert.Equal(t, int64(0), displayTx.OutputValue)
		assert.Nil(t, displayTx.XpubInIDs)
		assert.Nil(t, displayTx.XpubOutIDs)
		assert.Nil(t, displayTx.XpubMetadata)
		assert.Nil(t, displayTx.XpubOutputValue)
	})

	t.Run("display with xpub data", func(t *testing.T) {
		tx := Transaction{
			TransactionBase: TransactionBase{
				ID:  testTxID,
				Hex: "hex",
			},
			Model:           Model{},
			XpubInIDs:       IDs{testXPubID},
			XpubOutIDs:      nil,
			BlockHash:       "hash",
			BlockHeight:     123,
			Fee:             321,
			NumberOfInputs:  1,
			NumberOfOutputs: 2,
			DraftID:         testDraftID,
			TotalValue:      123499,
			XpubMetadata: XpubMetadata{
				testXPubID: Metadata{
					"test-key": "test-value",
				},
			},
			OutputValue: 12,
			XpubOutputValue: XpubOutputValue{
				testXPubID: 123499,
			},
			XPubID: testXPubID,
		}

		displayTx := tx.Display().(*Transaction)
		assert.Equal(t, Metadata{"test-key": "test-value"}, displayTx.Metadata)
		assert.Equal(t, int64(123499), displayTx.OutputValue)
		assert.Equal(t, TransactionDirectionIn, displayTx.Direction)
		assert.Nil(t, displayTx.XpubInIDs)
		assert.Nil(t, displayTx.XpubOutIDs)
		assert.Nil(t, displayTx.XpubMetadata)
		assert.Nil(t, displayTx.XpubOutputValue)
	})
}

func (ts *EmbeddedDBTestSuite) TestTransaction_Save() {
	parsedTx, errP := trx.NewTransactionFromHex(testTxHex)
	require.NoError(ts.T(), errP)
	require.NotNil(ts.T(), parsedTx)

	var parsedInTx *trx.Transaction
	parsedInTx, errP = trx.NewTransactionFromHex(testTx2Hex)
	require.NoError(ts.T(), errP)
	require.NotNil(ts.T(), parsedInTx)

	ts.T().Run("[sqlite] [in-memory] - Save transaction", func(t *testing.T) {
		tc := ts.genericDBClient(t, datastore.SQLite, false)
		defer tc.Close(tc.ctx)

		transaction, err := txFromHex(testTxHex, append(tc.client.DefaultModelOptions(), New())...)
		require.NoError(t, err)
		require.NotNil(t, transaction)

		err = transaction.Save(tc.ctx)
		require.NoError(t, err)

		var transaction2 *Transaction
		transaction2, err = tc.client.GetTransaction(tc.ctx, testXPubID, testTxID)
		require.NoError(t, err)
		require.NotNil(t, transaction2)
		assert.Equal(t, transaction2.ID, testTxID)

		// no utxos should have been saved, we don't recognize any of the destinations
		var utxo *Utxo
		utxo, err = getUtxo(tc.ctx, transaction.ID, 0, tc.client.DefaultModelOptions()...)
		require.NoError(t, err)
		assert.Nil(t, utxo)
	})

	ts.T().Run("[sqlite] [in-memory] - Save transaction - with utxos & outputs", func(t *testing.T) {
		tc := ts.genericDBClient(t, datastore.SQLite, false)
		defer tc.Close(tc.ctx)

		_, xPub, _ := CreateNewXPub(tc.ctx, t, tc.client)
		require.NotNil(t, xPub)

		_, xPub2, _ := CreateNewXPub(tc.ctx, t, tc.client)
		require.NotNil(t, xPub2)

		// NOTE: these are fake destinations, might want to replace with actual real data / methods

		// fake existing destinations, to generate utxos
		ls := parsedTx.Outputs[0].LockingScript
		destination := newDestination(xPub.GetID(), ls.String(), append(tc.client.DefaultModelOptions(), New())...)
		require.NotNil(t, destination)

		err := destination.Save(tc.ctx)
		require.NoError(t, err)

		ls2 := parsedTx.Outputs[1].LockingScript
		destination2 := newDestination(xPub2.GetID(), ls2.String(), append(tc.client.DefaultModelOptions(), New())...)
		require.NotNil(t, destination2)

		err = destination2.Save(tc.ctx)
		require.NoError(t, err)

		transaction, err := txFromHex(testTxHex, append(tc.client.DefaultModelOptions(), New())...)
		require.NoError(t, err)
		require.NotNil(t, transaction)

		err = transaction.processUtxos(tc.ctx)
		require.NoError(t, err)

		err = transaction.Save(tc.ctx)
		require.NoError(t, err)

		// check whether the XpubOutIDs were set properly
		var transaction2 *Transaction
		transaction2, err = tc.client.GetTransaction(tc.ctx, testXPubID, testTxID)
		require.NoError(t, err)
		require.NotNil(t, transaction2)
		assert.Equal(t, xPub.GetID(), transaction2.XpubOutIDs[0])
		assert.Equal(t, xPub2.GetID(), transaction2.XpubOutIDs[1])

		// utxos should have been saved for our fake destinations
		var utxo *Utxo
		utxo, err = getUtxo(tc.ctx, transaction.ID, 0, tc.client.DefaultModelOptions()...)
		require.NoError(t, err)
		require.NotNil(t, utxo)
		assert.Equal(t, xPub.GetID(), utxo.XpubID)
		assert.Equal(t, script.ScriptTypePubKeyHash, utxo.Type)
		assert.Equal(t, testTxScriptPubKey1, utxo.ScriptPubKey)
		assert.Empty(t, utxo.DraftID)
		assert.Empty(t, utxo.SpendingTxID)

		var utxo2 *Utxo
		utxo2, err = getUtxo(tc.ctx, transaction.ID, 1, tc.client.DefaultModelOptions()...)
		assert.Nil(t, err)
		assert.Equal(t, xPub2.GetID(), utxo2.XpubID)
		assert.Equal(t, script.ScriptTypePubKeyHash, utxo2.Type)
		assert.Equal(t, testTxScriptPubKey2, utxo2.ScriptPubKey)
		assert.Empty(t, utxo2.DraftID)
		assert.Empty(t, utxo2.SpendingTxID)
	})

	ts.T().Run("[sqlite] [in-memory] - Save transaction - with inputs", func(t *testing.T) {
		tc := ts.genericDBClient(t, datastore.SQLite, false)
		defer tc.Close(tc.ctx)

		_, xPub, _ := CreateNewXPub(tc.ctx, t, tc.client)
		require.NotNil(t, xPub)

		// NOTE: these are fake destinations, might want to replace with actual real data / methods

		// create a fake destination for our IN transaction
		ls := parsedInTx.Outputs[0].LockingScript
		destination := newDestination(xPub.GetID(), ls.String(), append(tc.client.DefaultModelOptions(), New())...)
		require.NotNil(t, destination)

		err := destination.Save(tc.ctx)
		require.NoError(t, err)

		// add the IN transaction
		transactionIn, err := txFromHex(testTx2Hex, append(tc.client.DefaultModelOptions(), New())...)
		require.NoError(t, err)
		require.NotNil(t, transactionIn)

		err = transactionIn.processUtxos(tc.ctx)
		require.NoError(t, err)

		err = transactionIn.Save(tc.ctx)
		require.NoError(t, err)

		var utxoIn *Utxo
		utxoIn, err = getUtxo(tc.ctx, transactionIn.ID, 0, tc.client.DefaultModelOptions()...)
		require.NotNil(t, utxoIn)
		require.NoError(t, err)
		assert.Equal(t, xPub.GetID(), utxoIn.XpubID)
		assert.Equal(t, script.ScriptTypePubKeyHash, utxoIn.Type)
		assert.Equal(t, testTxInScriptPubKey, utxoIn.ScriptPubKey)
		assert.Empty(t, utxoIn.SpendingTxID)

		draftConfig := &TransactionConfig{
			Outputs: []*TransactionOutput{{
				Satoshis: 202,
				To:       testExternalAddress,
			}},
		}
		draftTransaction, err := newDraftTransaction(
			xPub.rawXpubKey, draftConfig, append(tc.client.DefaultModelOptions(), New())...,
		)
		require.NoError(t, err)
		err = draftTransaction.Save(tc.ctx)
		require.NoError(t, err)

		// this transaction should spend the utxo of the IN transaction
		transaction, err := newTransactionWithDraftID(testTxHex, draftTransaction.ID,
			append(tc.client.DefaultModelOptions(), WithXPub(xPub.rawXpubKey), New())...)
		require.NoError(t, err)
		require.NotNil(t, transactionIn)

		transaction.draftTransaction = draftTransaction

		err = transaction.processUtxos(tc.ctx)
		require.NoError(t, err)

		err = transaction.Save(tc.ctx)
		require.NoError(t, err)

		// check whether the XpubInIDs were set properly
		var transaction2 *Transaction
		transaction2, err = tc.client.GetTransaction(tc.ctx, testXPubID, testTxID)
		require.NotNil(t, transaction2)
		require.NoError(t, err)
		assert.Equal(t, xPub.GetID(), transaction2.XpubInIDs[0])

		// Get the utxo for the IN transaction and make sure it is marked as spent
		var utxo *Utxo
		utxo, err = getUtxo(tc.ctx, transactionIn.ID, 0, tc.client.DefaultModelOptions()...)
		require.NotNil(t, transaction2)
		require.NoError(t, err)
		assert.Equal(t, testTxInID, utxo.ID)
		assert.True(t, utxo.SpendingTxID.Valid)
		assert.Equal(t, utxo.SpendingTxID.String, testTxID)
	})
}

// BenchmarkTransaction_newTransaction will benchmark the method newTransaction()
func BenchmarkTransaction_newTransaction(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = txFromHex(testTxHex, New())
	}
}
