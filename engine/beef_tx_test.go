package engine

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type beefTestCase struct {
	testID               int
	name                 string
	hexForProcessedTx    string
	outputValue          uint64
	receiverAddress      string
	ancestors            []*beefTestCaseAncestor
	expectedErrorMessage string
}

type beefTestCaseAncestor struct {
	// reverse condition to not set this value every time
	doNotAddToStore bool
	hex             string
	isMined         bool
	bumpJSON        string
	blockHeight     int
	parents         []*beefTestCaseAncestor
}

func Test_ToBeef_HappyPaths(t *testing.T) {
	testCases := []beefTestCase{
		{
			testID:            1,
			hexForProcessedTx: "0100000002cb3553424ffc94b59a60fb358b6cb6dfb694aee894dcd1effc0ed0a9052464e3000000006a4730440220515c3bf93d38fa7cc164746fae4bec8b66c60a82509eb553751afa5971c3e41d0220321517fd5c997ab5f8ef0e59048ce9157de46f92b10d882bf898e62f3ee7343d4121038f1273fcb299405d8d140b4de9a2111ecb39291b2846660ebecd864d13bee575ffffffff624fbcb4e68d162361f456b8b4fef6b9e7943013088b32b6bca7f5ced41ff004010000006a47304402203fb24f6e00a6487cf88a3b39d8454786db63d649142ea76374c2f55990777e6302207fbb903d038cf43e13ffb496a64f36637ec7323e5ac48bb96bdb4a885100abca4121024b003d3cf49a8f48c1fe79b711b1d08e306c42a0ab8da004d97fccc4ced3343affffffff026f000000000000001976a914f232d38cd4c2f87c117af06542b04a7061b6640188aca62a0000000000001976a9146058e52d00e3b94211939f68cc2d9a3fc1e3db0f88ac00000000",
			name:              "all inputs are already mined - 2 inputs on the same level",
			ancestors: []*beefTestCaseAncestor{
				{
					hex:         "0100000001cfc39e3adcd58ed58cf590079dc61c3eb6ec739abb7d22b592fb969d427f33ee000000006a4730440220253e674e64028459457d55b444f5f3dc15c658425e3184c628016739e4921fd502207c8fe20eb34e55e4115fbd82c23878b4e54f01f6c6ad0811282dd0b1df863b5e41210310a4366fd997127ad972b14c56ca2e18f39ca631ac9e3e4ad3d9827865d0cc70ffffffff0264000000000000001976a914668a92ff9cb5785eb8fc044771837a0818b028b588acdc4e0000000000001976a914b073264927a61cf84327dea77414df6c28b11e5988ac00000000",
					isMined:     true,
					bumpJSON:    `{"blockHeight":"817574","path":[[{"offset":"11432","hash":"3b535e0f8e266124bce9868420052d5a7585c67e82c1edc2c7fe05fd5e140307"},{"offset":"11433","hash":"e3642405a9d00efcefd1dc94e8ae94b6dfb66c8b35fb609ab594fc4f425335cb","txid":true}],[{"offset":"5717","hash":"6ef9c6dde7fff82fa893754109f12378c8453b47dc896596b5531433093ab5b7"}],[{"offset":"2859","hash":"daa67e00ad2aef787998b66cbb3417033fbec136da1e230a5f5df3186f5c0880"}],[{"offset":"1428","hash":"bc777a80d951fbf2b7bd3a8048a9bb78fbf1d23d4127290c3fed9740b4246dd2"}],[{"offset":"715","hash":"762b57f88e7258f5757b48cda96d075cbe767c0a39a83e7109574555fd2dd8ba"}],[{"offset":"356","hash":"bbaab745bcca4f8a4be39c06c7e9be3aa1994f32271e3c6b4f768897153e5522"}],[{"offset":"179","hash":"817694ccbde5dbf88f290c30e8735991708a3d406740f7dd31434ff516a5bfde"}],[{"offset":"88","hash":"ed5b52ba4af9198d398e934a84e18405f49e7abde91cafb6dfe5aeaedb33a979"}],[{"offset":"45","hash":"0e51ec9dd5319ceb32d2d20f620c0ca3e0d918260803c1005d49e686c9b18752"}],[{"offset":"23","hash":"08ab694ef1af4019e2999a543a632cf4a662ae04d5fee879c6aadaeb749f6374"}],[{"offset":"10","hash":"4223f47597b14ee0fa7ade08e611ec80948b5fa9da267ce6c8e5d952e7fdb38e"}],[{"offset":"4","hash":"b6dace0d2294fd6e0c11f74376b7f0a1fc8ee415b350caf90c3ae92749e2a8ee"}],[{"offset":"3","hash":"795e7514ebf6d63b454d3f04854e1e0db0ac3a549f61135d5e9ef8d5785f2c68"}],[{"offset":"0","hash":"3f458f2c06493c31cbc3a035ba131913b274ac7915b9b9bc79128001a75cf76d"}],[{"offset":"1","hash":"b9b9f80cc72a674e37b54a9fdee72a9bff761f8cbcb94146afc2bffef33be89f"}]]}`,
					blockHeight: 817574,
				},
				{
					hex:         "0100000001a114c7deb8deba851d87755aa10aa18c97bd77afee4e1bad01d1c50e07a644eb010000006a473044022041abd4f93bd1db1d0097f2d467ae183801d7842d23d0605fa9568040d245167402201be66c96bef4d6d051304f6df2aecbdfe23a8a05af0908ef2117ab5388d8903c412103c08545a40c819f6e50892e31e792d221b6df6da96ebdba9b6fe39305cc6cc768ffffffff0263040000000000001976a91454097d9d921f9a1f55084a943571d868552e924f88acb22a0000000000001976a914c36b3fca5159231033f3fbdca1cde942096d379f88ac00000000",
					isMined:     true,
					bumpJSON:    `{"blockHeight":"819138","path":[[{"offset":"648","hash":"121bca23ca64d9b925c055f89340802e51f0949ab5edf36ad6dffe050d28400e"},{"offset":"649","hash":"04f01fd4cef5a7bcb6328b08133094e7b9f6feb4b856f46123168de6b4bc4f62","txid":true}],[{"offset":"325","hash":"1e5b72effab8fb56da368f25bab8d8fae7891a5dc70c5da1a8dac4f81e75f990"}],[{"offset":"163","hash":"c97efaa344c57f5e0e46676cbf8629fad9f69a7b1f71d6fda8a1e03f2b546328"}],[{"offset":"80","hash":"5069f334f680952ee9abf37ca5f1cf327e5114920f79ec26b108ae7a491e0b3a"}],[{"offset":"41","hash":"a40cf9eb878b35f853198ebf23ac85061253ffb6e20c4c3eb1ac546b2a376f6d"}],[{"offset":"21","hash":"b5f91b76bf448529368e9421a89e6c756d6b92679ce06479557bf8a5dedb10c3"}],[{"offset":"11","hash":"dbf6acf27c7df7bf4a9de100fd6ad7f73db9b0d659e38235cef4eee26fe367a2"}],[{"offset":"4","hash":"e08fbe8bbdda28ff48478b90c909e4dad7acffc6ff5b3e46f8b4c597d76fc180"}],[{"offset":"3","hash":"0ae8ff1e623834f0624d78703498bb986e1d3d2c5f9d172f05b8e839a09ce0b7"}],[{"offset":"0","hash":"1f55fb14746170226dd929e47fedaac59c65d7b4b9b5502c758bca19a76d5bcf"}],[{"offset":"1","hash":"40ab6623661b0a927bbf231c8a1c0bbf1b64b0b0afb665449cac9ac70e8601dc"}]]}`,
					blockHeight: 819138,
				},
			},
			receiverAddress:      "1A1PjKqjWMNBzTVdcBru27EV1PHcXWc63W",
			outputValue:          3000,
			expectedErrorMessage: "",
		},
		{
			testID:            2,
			hexForProcessedTx: "0100000003e01e74cab9a0571ab5a7d86794826f756a9c65dd0dea3bb3720c4051c488cf50000000006b483045022100bc7fc6ace1a5b1ab8601599d56b3adad4a11b7f11757f3225e96b46ca1ab7f7c0220324d6074aa987a7c63c404ac5b03c26e55d3c4209e298b4ca9df0e90aca43ef3412103ee05b34332b5662830c600b73f9c908bb8bff1813bc9b2690e9cad00fad23d3cffffffff08c461a39a8877db46472f5cc59e5a108e417b1c9ea3091b71b65346d218f471000000006b483045022100a936c496423ec03b1ad0f3bfe2348572d7b29ab14e4435c0c8e2ee093d930fde02203d9e86647ea18043c150289f74c6cf2ceb9ca3b228ae31c7b19c4eef813fb68d412103a19014bcc672ccdf18abb6972dd699367baed89c29b704385253ce2ae0eddad5ffffffff2256c94d07451664749e440f55cec8a37da1c46cf30a97579e2f9696b84ad484000000006b48304502210091b0bcf2e84d9ee65de437e8396b379941345e4cffac331af2ae29b8a16968a602205a00eed18a7ffe36f59ae6eb477d9002324cfc249c875260e6ade5bce852692d4121021446bd1df2b61952088a22a516550e43cd95e47ca2a778822d21268bd8b1cebeffffffff02c4090000000000001976a91497ebeffef6d9dd88ffbce922f1df97cbcd7f88d388ac42000000000000001976a91449457f2c101859d1c8ff90096385d3cc30e5488388ac00000000",
			name:              "all inputs are already mined - 3 inputs on the same level",
			ancestors: []*beefTestCaseAncestor{
				{
					hex:         "01000000019ed68f94dfa952554d777dbaa9e5c01acb3df767e40cabad7b6fb7547bfa871a010000006a4730440220287534d6ff51166e014ad91a2b677be4bd88cf08785624006cdb66553eafc8cf02204862f38e9d2982a5ee95a7850222f2208bff38637349ecfe41abe185498e4ead4121035ca1a2c6d2b46c61fd29e7697018f5ce2bae1ae735e23627046a2dd17ca8fb24ffffffff02de000000000000001976a914f5c9505bf02a4a2fb591e3568183f9c53cf157be88aca62b0000000000001976a91489b5e639bce3209e0888ea8b7eb4203de1c6148888ac00000000",
					isMined:     true,
					bumpJSON:    `{"blockHeight":"819265","path":[[{"offset":"1484","hash":"84d44ab896962f9e57970af36cc4a17da3c8ce550f449e74641645074dc95622","txid":true},{"offset":"1485","hash":"657410f5ed93c9eb79d5de3f369ab112dcc0cd693a8c75d8ba83d419a5e0cbae"}],[{"offset":"743","hash":"1c610e4e420187f4625bf8bb4bd97d9d8c0189b09016a2653339c9eedfa19586"}],[{"offset":"370","hash":"9d84d9625ba6b72bd0cf4d4a4206abc528bf6f37e0e3c50d8af2740f7e268712"}],[{"offset":"184","hash":"f93ddfa2e2a485be23b527f00b6d2d12c0cac9a1d9ef031e3d49fc8d97265a68"}],[{"offset":"93","hash":"3cb99ac1431ce21b4e4c89dc2561b304ae0825a2514227eb2405990be9259d55"}],[{"offset":"47","hash":"0f8fbdc287329d0dd88aa672353452c41cb200b4034325caca35f4a19f1e4637"}],[{"offset":"22","hash":"296826f3a55a1dbc5ff550657be9a24dc97fe65ce352e9ac056d5341cef80910"}],[{"offset":"10","hash":"3d766a898268ba2ea8cb40baabb806582e06addc9f65c4e7f9a441a52207eaeb"}],[{"offset":"4","hash":"7712fc3a6f8653ad62dd035e1f8a66c8190a031a1783dd3d273ffce7364dde36"}],[{"offset":"3","hash":"8dfdcd43d0b3ea6cdb50cf2cc12280462e4b92bfe123dd97a044e12bd619998e"}],[{"offset":"0","hash":"3c9d541c4bd24c6e59e2e90a8b1c9fe271d591f5491d032bd95202463a4898ae"}]]}`,
					blockHeight: 819265,
				},
				{
					hex:         "010000000154aa46f1b3b7bde36c02e293b74d53e6c6eaed7411d286183b1dca766f42879a010000006b483045022100cd21d346073b4a0788018ff6938c44395d14cf5759fcc35a0899a8fe35a3c2a0022064eb9a005c3d0be03b61ab0e1c8757ed566dd935dacac37fcd1452adba4994b541210272d67492c31d0e6bead28c934fb1c9bb50ba9b46f886209fe95fb6a3e43bb27bffffffff0257040000000000001976a9140501308b6409cca5a7b5768c18ff2de8da4c1fa388ac39420000000000001976a91417e3d89f4aeacd5b4929fe04edc32c79b6182e1988ac00000000",
					isMined:     true,
					bumpJSON:    `{"blockHeight":"819265","path":[[{"offset":"1168","hash":"71f418d24653b6711b09a39e1c7b418e105a9ec55c2f4746db77889aa361c408","txid":true},{"offset":"1169","hash":"f487a07dc7f72556292f0bd85e5f22a05d9dab29772daefe48870895e3a18fff"}],[{"offset":"585","hash":"89ad2f3acd22cd3313054427d16bf2e0e412a95cdba2b77b97f7e338a911330f"}],[{"offset":"293","hash":"68a8bf0c4b6c378f5120a9f716e0452c7434fabc5b6e5bb5f2cde1286a985f5e"}],[{"offset":"147","hash":"d3c4698449afbd39403ddd35272f0c5f45726acc87fce700369aa17970a00aa6"}],[{"offset":"72","hash":"a3b2da485a2d0331969dc5e9b38550e4cbfb69a127418ef32252780ce4baa9b4"}],[{"offset":"37","hash":"2e11be0a25536b375e91257a85631a77f4194db005c44a106f90cc61be88af6d"}],[{"offset":"19","hash":"3b4be5245bf5e1e18362ebf3498447a7aaf7c5893b6f344cfe0a0a44a0517b59"}],[{"offset":"8","hash":"14f1cf466a0cfcebad570548d09bb180b9964bcbdee04da0807b63127e7668b6"}],[{"offset":"5","hash":"f02af5d04771132d1b6f08ed70f0d7a3fa868edd6006c8d90ddda6c967c99d5c"}],[{"offset":"3","hash":"8dfdcd43d0b3ea6cdb50cf2cc12280462e4b92bfe123dd97a044e12bd619998e"}],[{"offset":"0","hash":"3c9d541c4bd24c6e59e2e90a8b1c9fe271d591f5491d032bd95202463a4898ae"}]]}`,
					blockHeight: 819265,
				},
				{
					hex:         "0100000001e230ab1b300ac3ce334590fc308fee93ddbb252f6e4645e0a20f7e30dd541289010000006b483045022100a611fdf01eca42289d80e1265584e5bd487faa72e6142ebbc140a676f7c5037c0220409282aaadf580f458d97d61db43c94ac343e0b40674a80fd3ac47f43fd0c66c4121020a87e70cc26f7d5fe775f622d2705f27cfd6f5d2b574fea75401d6412a58b91affffffff02d2040000000000001976a9145d2117c4f66bdb335ce2707a74c46fa46d02cdb388acf23b0000000000001976a914effd80ee9df812990a8d7834fa8610491cbeb91688ac00000000",
					isMined:     true,
					bumpJSON:    `{"blockHeight":"819265","path":[[{"offset":"1068","hash":"50cf88c451400c72b33bea0ddd659c6a756f829467d8a7b51a57a0b9ca741ee0","txid":true},{"offset":"1069","hash":"b42c92309cabe31f84a04281331c4d4d3288d77e6696fbaa0f1cfb3a540c06d9"}],[{"offset":"535","hash":"847c23a1d815cbe09c0876402de53564945061945324a8f77e77f8635e795ddd"}],[{"offset":"266","hash":"d2fd8db438af79b228e7edac40eaa204c83191bf96c78365a5ab6418023d3722"}],[{"offset":"132","hash":"6520f0990428421e5632666380c27af04fc3d720d1c915ef06a1032bb99e10a8"}],[{"offset":"67","hash":"9272d24f6f60ed3c80e2eae3657fc4b83d89bd466454d0ca652c37f4d6c5786f"}],[{"offset":"32","hash":"53df03b923f4eac53729382a7fa611a0106cbb90e69310a37daad2792f3fbc46"}],[{"offset":"17","hash":"ab516a8baa05c2dc560cd929cbcf933a0e5574d073ab9d3dbb0e6720aaeb4ef1"}],[{"offset":"9","hash":"87cf2caba2df8fd0ca00f509ccf58acf5fecf6df26bce8f79fddbe6736ec4cd2"}],[{"offset":"5","hash":"f02af5d04771132d1b6f08ed70f0d7a3fa868edd6006c8d90ddda6c967c99d5c"}],[{"offset":"3","hash":"8dfdcd43d0b3ea6cdb50cf2cc12280462e4b92bfe123dd97a044e12bd619998e"}],[{"offset":"0","hash":"3c9d541c4bd24c6e59e2e90a8b1c9fe271d591f5491d032bd95202463a4898ae"}]]}`,
					blockHeight: 819265,
				},
			},
			receiverAddress:      "1A1PjKqjWMNBzTVdcBru27EV1PHcXWc63W",
			outputValue:          4000,
			expectedErrorMessage: "",
		},
		{
			testID:            3,
			hexForProcessedTx: "0100000002cb3553424ffc94b59a60fb358b6cb6dfb694aee894dcd1effc0ed0a9052464e3000000006a4730440220515c3bf93d38fa7cc164746fae4bec8b66c60a82509eb553751afa5971c3e41d0220321517fd5c997ab5f8ef0e59048ce9157de46f92b10d882bf898e62f3ee7343d4121038f1273fcb299405d8d140b4de9a2111ecb39291b2846660ebecd864d13bee575ffffffff624fbcb4e68d162361f456b8b4fef6b9e7943013088b32b6bca7f5ced41ff004010000006a47304402203fb24f6e00a6487cf88a3b39d8454786db63d649142ea76374c2f55990777e6302207fbb903d038cf43e13ffb496a64f36637ec7323e5ac48bb96bdb4a885100abca4121024b003d3cf49a8f48c1fe79b711b1d08e306c42a0ab8da004d97fccc4ced3343affffffff026f000000000000001976a914f232d38cd4c2f87c117af06542b04a7061b6640188aca62a0000000000001976a9146058e52d00e3b94211939f68cc2d9a3fc1e3db0f88ac00000000",
			name:              "not all inputs are mined but all required ancestors are mined - one level below inputs",
			ancestors: []*beefTestCaseAncestor{
				{
					hex:         "0100000001cfc39e3adcd58ed58cf590079dc61c3eb6ec739abb7d22b592fb969d427f33ee000000006a4730440220253e674e64028459457d55b444f5f3dc15c658425e3184c628016739e4921fd502207c8fe20eb34e55e4115fbd82c23878b4e54f01f6c6ad0811282dd0b1df863b5e41210310a4366fd997127ad972b14c56ca2e18f39ca631ac9e3e4ad3d9827865d0cc70ffffffff0264000000000000001976a914668a92ff9cb5785eb8fc044771837a0818b028b588acdc4e0000000000001976a914b073264927a61cf84327dea77414df6c28b11e5988ac00000000",
					isMined:     true,
					bumpJSON:    `{"blockHeight":"817574","path":[[{"offset":"11432","hash":"3b535e0f8e266124bce9868420052d5a7585c67e82c1edc2c7fe05fd5e140307"},{"offset":"11433","hash":"e3642405a9d00efcefd1dc94e8ae94b6dfb66c8b35fb609ab594fc4f425335cb","txid":true}],[{"offset":"5717","hash":"6ef9c6dde7fff82fa893754109f12378c8453b47dc896596b5531433093ab5b7"}],[{"offset":"2859","hash":"daa67e00ad2aef787998b66cbb3417033fbec136da1e230a5f5df3186f5c0880"}],[{"offset":"1428","hash":"bc777a80d951fbf2b7bd3a8048a9bb78fbf1d23d4127290c3fed9740b4246dd2"}],[{"offset":"715","hash":"762b57f88e7258f5757b48cda96d075cbe767c0a39a83e7109574555fd2dd8ba"}],[{"offset":"356","hash":"bbaab745bcca4f8a4be39c06c7e9be3aa1994f32271e3c6b4f768897153e5522"}],[{"offset":"179","hash":"817694ccbde5dbf88f290c30e8735991708a3d406740f7dd31434ff516a5bfde"}],[{"offset":"88","hash":"ed5b52ba4af9198d398e934a84e18405f49e7abde91cafb6dfe5aeaedb33a979"}],[{"offset":"45","hash":"0e51ec9dd5319ceb32d2d20f620c0ca3e0d918260803c1005d49e686c9b18752"}],[{"offset":"23","hash":"08ab694ef1af4019e2999a543a632cf4a662ae04d5fee879c6aadaeb749f6374"}],[{"offset":"10","hash":"4223f47597b14ee0fa7ade08e611ec80948b5fa9da267ce6c8e5d952e7fdb38e"}],[{"offset":"4","hash":"b6dace0d2294fd6e0c11f74376b7f0a1fc8ee415b350caf90c3ae92749e2a8ee"}],[{"offset":"3","hash":"795e7514ebf6d63b454d3f04854e1e0db0ac3a549f61135d5e9ef8d5785f2c68"}],[{"offset":"0","hash":"3f458f2c06493c31cbc3a035ba131913b274ac7915b9b9bc79128001a75cf76d"}],[{"offset":"1","hash":"b9b9f80cc72a674e37b54a9fdee72a9bff761f8cbcb94146afc2bffef33be89f"}]]}`,
					blockHeight: 817574,
				},
				{
					hex:         "0100000001a114c7deb8deba851d87755aa10aa18c97bd77afee4e1bad01d1c50e07a644eb010000006a473044022041abd4f93bd1db1d0097f2d467ae183801d7842d23d0605fa9568040d245167402201be66c96bef4d6d051304f6df2aecbdfe23a8a05af0908ef2117ab5388d8903c412103c08545a40c819f6e50892e31e792d221b6df6da96ebdba9b6fe39305cc6cc768ffffffff0263040000000000001976a91454097d9d921f9a1f55084a943571d868552e924f88acb22a0000000000001976a914c36b3fca5159231033f3fbdca1cde942096d379f88ac00000000",
					isMined:     false,
					bumpJSON:    ``,
					blockHeight: -1,
					parents: []*beefTestCaseAncestor{
						{
							hex:         "010000000150965003ea3d2c08bc79b116c9ffe7e730c9f9cf0a61e3df07868b24eac6f8d3000000006b4830450221009d3489f9e76ff3b043708972c52f85519e50a5fc35563d405e04b668780bf2ba0220024188508fc9c6870b2fc4f40b9484ae4163481199a5b4a7a338b86ec8952fee4121036a8b9d796ce2dee820d1f6d7a6ba07037dab4758f16028654fe4bc3a5c430b40ffffffff022a200000000000001976a91484c73348a8fbbc44cfa34f8f5441fc104f3bc78588ac162f0000000000001976a914590b1df63948c2c4e7a12a6e52012b36e25daa9888ac00000000",
							isMined:     true,
							bumpJSON:    `{"blockHeight":"817267","path":[[{"offset":"204","hash":"cc688d86f3ceb67d53bbdd4b140b10a1ff0cff919e2e1e45c5d82b2629a52f5d"},{"offset":"205","hash":"eb44a6070ec5d101ad1b4eeeaf77bd978ca10aa15a75871d85badeb8dec714a1","txid":true}],[{"offset":"103","hash":"0b0bae3a36eaf8fb10f8d52d61a04434eea745206bff04943a2d0361ebf52b67"}],[{"offset":"50","hash":"43af8bb7133505db727bf626f81081ab841ca90def48eba8bf70f5998d9a2af6"}],[{"offset":"24","hash":"18bf5d884d88570eea3006f46ca62d2ab13d41f4b2b7741c0fb72a7d3c19ec0f"}],[{"offset":"13","hash":"fccf4bf884dda23d6639b0768513be141c329ce8674eda6184eb28887859f227"}],[{"offset":"7","hash":"3353b9bd9666d87045c9661803dbe63dad7270b12580ef25cf3db45e6dddf4eb"}],[{"offset":"2","hash":"38dbe8ee0853c4adb254a2fd73040b02f5f98c6db65db49256515e5401cbd368"}],[{"offset":"0","hash":"c0542ab8ca895ea798840c8b90eaa27112c3e3ab7acb8500749f482cb73d9172"}],[{"offset":"1","hash":"c22b61b365e9ff54dd317122d1316a8244e826604650c934b0f957ad87771b49"}],[{"offset":"1","hash":"0de03905ab373fcc44443c17d54148c101bc6e5076b74f5fbcc52e37a6fd831f"}],[{"offset":"1","hash":"c1f2dbb43a68c3893784efef0c45a998fccceb63d873c8c7e3fc7a7ced33bfce"}],[{"offset":"1","hash":"2a7f4a0e8008299ac6e258229aedc67ba1d4069ebb78f653501e78991d5f36dc"}]]}`,
							blockHeight: 817267,
						},
					},
				},
			},
			receiverAddress:      "1A1PjKqjWMNBzTVdcBru27EV1PHcXWc63W",
			outputValue:          3500,
			expectedErrorMessage: "",
		},
		{
			testID:            4,
			hexForProcessedTx: "0100000002cb3553424ffc94b59a60fb358b6cb6dfb694aee894dcd1effc0ed0a9052464e3000000006a4730440220515c3bf93d38fa7cc164746fae4bec8b66c60a82509eb553751afa5971c3e41d0220321517fd5c997ab5f8ef0e59048ce9157de46f92b10d882bf898e62f3ee7343d4121038f1273fcb299405d8d140b4de9a2111ecb39291b2846660ebecd864d13bee575ffffffff624fbcb4e68d162361f456b8b4fef6b9e7943013088b32b6bca7f5ced41ff004010000006a47304402203fb24f6e00a6487cf88a3b39d8454786db63d649142ea76374c2f55990777e6302207fbb903d038cf43e13ffb496a64f36637ec7323e5ac48bb96bdb4a885100abca4121024b003d3cf49a8f48c1fe79b711b1d08e306c42a0ab8da004d97fccc4ced3343affffffff026f000000000000001976a914f232d38cd4c2f87c117af06542b04a7061b6640188aca62a0000000000001976a9146058e52d00e3b94211939f68cc2d9a3fc1e3db0f88ac00000000",
			name:              "not all inputs are mined but all required ancestors are mined - two levels below inputs",
			ancestors: []*beefTestCaseAncestor{
				{
					hex:         "0100000001cfc39e3adcd58ed58cf590079dc61c3eb6ec739abb7d22b592fb969d427f33ee000000006a4730440220253e674e64028459457d55b444f5f3dc15c658425e3184c628016739e4921fd502207c8fe20eb34e55e4115fbd82c23878b4e54f01f6c6ad0811282dd0b1df863b5e41210310a4366fd997127ad972b14c56ca2e18f39ca631ac9e3e4ad3d9827865d0cc70ffffffff0264000000000000001976a914668a92ff9cb5785eb8fc044771837a0818b028b588acdc4e0000000000001976a914b073264927a61cf84327dea77414df6c28b11e5988ac00000000",
					isMined:     true,
					bumpJSON:    `{"blockHeight":"817574","path":[[{"offset":"11432","hash":"3b535e0f8e266124bce9868420052d5a7585c67e82c1edc2c7fe05fd5e140307"},{"offset":"11433","hash":"e3642405a9d00efcefd1dc94e8ae94b6dfb66c8b35fb609ab594fc4f425335cb","txid":true}],[{"offset":"5717","hash":"6ef9c6dde7fff82fa893754109f12378c8453b47dc896596b5531433093ab5b7"}],[{"offset":"2859","hash":"daa67e00ad2aef787998b66cbb3417033fbec136da1e230a5f5df3186f5c0880"}],[{"offset":"1428","hash":"bc777a80d951fbf2b7bd3a8048a9bb78fbf1d23d4127290c3fed9740b4246dd2"}],[{"offset":"715","hash":"762b57f88e7258f5757b48cda96d075cbe767c0a39a83e7109574555fd2dd8ba"}],[{"offset":"356","hash":"bbaab745bcca4f8a4be39c06c7e9be3aa1994f32271e3c6b4f768897153e5522"}],[{"offset":"179","hash":"817694ccbde5dbf88f290c30e8735991708a3d406740f7dd31434ff516a5bfde"}],[{"offset":"88","hash":"ed5b52ba4af9198d398e934a84e18405f49e7abde91cafb6dfe5aeaedb33a979"}],[{"offset":"45","hash":"0e51ec9dd5319ceb32d2d20f620c0ca3e0d918260803c1005d49e686c9b18752"}],[{"offset":"23","hash":"08ab694ef1af4019e2999a543a632cf4a662ae04d5fee879c6aadaeb749f6374"}],[{"offset":"10","hash":"4223f47597b14ee0fa7ade08e611ec80948b5fa9da267ce6c8e5d952e7fdb38e"}],[{"offset":"4","hash":"b6dace0d2294fd6e0c11f74376b7f0a1fc8ee415b350caf90c3ae92749e2a8ee"}],[{"offset":"3","hash":"795e7514ebf6d63b454d3f04854e1e0db0ac3a549f61135d5e9ef8d5785f2c68"}],[{"offset":"0","hash":"3f458f2c06493c31cbc3a035ba131913b274ac7915b9b9bc79128001a75cf76d"}],[{"offset":"1","hash":"b9b9f80cc72a674e37b54a9fdee72a9bff761f8cbcb94146afc2bffef33be89f"}]]}`,
					blockHeight: 817574,
				},
				{
					hex:         "0100000001a114c7deb8deba851d87755aa10aa18c97bd77afee4e1bad01d1c50e07a644eb010000006a473044022041abd4f93bd1db1d0097f2d467ae183801d7842d23d0605fa9568040d245167402201be66c96bef4d6d051304f6df2aecbdfe23a8a05af0908ef2117ab5388d8903c412103c08545a40c819f6e50892e31e792d221b6df6da96ebdba9b6fe39305cc6cc768ffffffff0263040000000000001976a91454097d9d921f9a1f55084a943571d868552e924f88acb22a0000000000001976a914c36b3fca5159231033f3fbdca1cde942096d379f88ac00000000",
					isMined:     false,
					bumpJSON:    ``,
					blockHeight: -1,
					parents: []*beefTestCaseAncestor{
						{
							hex:         "010000000150965003ea3d2c08bc79b116c9ffe7e730c9f9cf0a61e3df07868b24eac6f8d3000000006b4830450221009d3489f9e76ff3b043708972c52f85519e50a5fc35563d405e04b668780bf2ba0220024188508fc9c6870b2fc4f40b9484ae4163481199a5b4a7a338b86ec8952fee4121036a8b9d796ce2dee820d1f6d7a6ba07037dab4758f16028654fe4bc3a5c430b40ffffffff022a200000000000001976a91484c73348a8fbbc44cfa34f8f5441fc104f3bc78588ac162f0000000000001976a914590b1df63948c2c4e7a12a6e52012b36e25daa9888ac00000000",
							isMined:     false,
							bumpJSON:    ``,
							blockHeight: -1,
							parents: []*beefTestCaseAncestor{
								{
									hex:         "0100000002787a565270ec00b1bf6ed20100223176656705dc0cfe5ef9d1810ca6569f12d1020000006a47304402203cfe36be7ff5c2ac939bb6a625e4a1226be242f1f9950672b5f696ec58a3358902202a48d6c6e81e5950dc49d0dd1a35b46fa8f919b109b0e7c05deaef3db6051890412102fb130326dbd7c43841cde467196e5f289b9d8596e237725df84f768468426d8bffffffff008d9db2a5c8c310e6394c24c1f3c23b3adbdd6ab4a719e917a4a0ed78768773020000006a473044022049c80385f7f69e8ba6039ebe84fe5e6578f4c3c83eb622442a96219c59ac1a750220317fe2b47838dff11f88d909732d0846eba20acff57cb357a3ff39b5a7b61b3741210322b79b40a759c485eac318eabba60a73a49ec3307ded79ba8c47204405bb2f3fffffffff05414f0000000000001976a91400414bcf2602f309171901d837b4a155adbfb5ce88ac50c30000000000001976a91489ef778cc07c77cce1ad3ff6274615afe15f20c088ac204e0000000000001976a914971b76df1dc6acf01e8e7d2f8bfb3c86e69bc64c88acef250000000000001976a9144b4a836b444d5ed8d245ddb1aa878908e36cd6b588ac9d860100000000001976a9144405da67e318e9cfd9d6ce9dffce27af5f60522888ac00000000",
									isMined:     true,
									bumpJSON:    `{"blockHeight":"817117","path":[[{"offset":"90","hash":"d3f8c6ea248b8607dfe3610acff9c930e7e7ffc916b179bc082c3dea03509650","txid":true},{"offset":"91","hash":"5b52ad65ab613867da9a710d60898a6e5da62dea97dac25da40a0dc385253ad2"}],[{"offset":"44","hash":"84c338bea7f65ccaf7a27ca9ae6d4b11372339cf6aa6021523de3ce6f5fe4f0c"}],[{"offset":"23","hash":"5860f292e051c0a5d9d8d69a451311a009c9cde8da6522df915587913a5180dd"}],[{"offset":"10","hash":"633fb08a689363af6a8245d3482fff232b27a62b94a4d119e67700fb9608ef78"}],[{"offset":"4","hash":"4b80bb130cebb1b8c313eb4088d098178ae122fd490a255218ceada19ab9eb52"}],[{"offset":"3","hash":"cf3d0335dda3223c8b4cf28ca2c03c7e025e3088525d51981d0ee1bd2ea210cf"}],[{"offset":"0","hash":"99c7462c2530abd1be779b170b7c2afbf7b883c07175871c971734d2bd38d35b"}],[{"offset":"1","hash":"9a66c7e35426281b1be6f43ecad44a3b65a9d2234d69a55b87d535f5903d677f"}],[{"offset":"1","hash":"40e34161018499a3ad5d1ef0d74a2e557733b6c7b5c07c1d8b872ffd504e545b"}],[{"offset":"1","hash":"dce6b1d7924d3ea2b1d7e3f2c4ae15abeea3b63e1362cc0896a82bca57a21387"}],[{"offset":"1","hash":"e1fe4e2b97189cb9da328bc59430a51595ee248079bd01a3f658af810e14cb7c"}]]}`,
									blockHeight: 817117,
								},
							},
						},
					},
				},
			},
			receiverAddress:      "1A1PjKqjWMNBzTVdcBru27EV1PHcXWc63W",
			outputValue:          4500,
			expectedErrorMessage: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			ctx, client, deferMe := initSimpleTestCase(t)
			store := NewMockTransactionStore()
			defer deferMe()

			var ancestors []*Transaction
			for _, ancestor := range tc.ancestors {
				ancestors = append(ancestors, addAncestor(ctx, ancestor, client, store, t))
			}
			newTx := createProcessedTx(ctx, t, client, &tc, ancestors)

			// when
			result, err := ToBeef(ctx, newTx, store)

			// then
			assert.Equal(t, expectedBeefHex[tc.testID], result)
			assert.NoError(t, nil, err)
		})
	}
}

func Test_ToBeef_ErrorPaths(t *testing.T) {
	testCases := []beefTestCase{
		{
			hexForProcessedTx: "0100000002cb3553424ffc94b59a60fb358b6cb6dfb694aee894dcd1effc0ed0a9052464e3000000006a4730440220515c3bf93d38fa7cc164746fae4bec8b66c60a82509eb553751afa5971c3e41d0220321517fd5c997ab5f8ef0e59048ce9157de46f92b10d882bf898e62f3ee7343d4121038f1273fcb299405d8d140b4de9a2111ecb39291b2846660ebecd864d13bee575ffffffff624fbcb4e68d162361f456b8b4fef6b9e7943013088b32b6bca7f5ced41ff004010000006a47304402203fb24f6e00a6487cf88a3b39d8454786db63d649142ea76374c2f55990777e6302207fbb903d038cf43e13ffb496a64f36637ec7323e5ac48bb96bdb4a885100abca4121024b003d3cf49a8f48c1fe79b711b1d08e306c42a0ab8da004d97fccc4ced3343affffffff026f000000000000001976a914f232d38cd4c2f87c117af06542b04a7061b6640188aca62a0000000000001976a9146058e52d00e3b94211939f68cc2d9a3fc1e3db0f88ac00000000",
			name:              "one input is mined and properly stored, the second one is missing in the store - should return error",
			ancestors: []*beefTestCaseAncestor{
				{
					hex:         "0100000001cfc39e3adcd58ed58cf590079dc61c3eb6ec739abb7d22b592fb969d427f33ee000000006a4730440220253e674e64028459457d55b444f5f3dc15c658425e3184c628016739e4921fd502207c8fe20eb34e55e4115fbd82c23878b4e54f01f6c6ad0811282dd0b1df863b5e41210310a4366fd997127ad972b14c56ca2e18f39ca631ac9e3e4ad3d9827865d0cc70ffffffff0264000000000000001976a914668a92ff9cb5785eb8fc044771837a0818b028b588acdc4e0000000000001976a914b073264927a61cf84327dea77414df6c28b11e5988ac00000000",
					isMined:     false,
					bumpJSON:    ``,
					blockHeight: -1,
				},
				{
					hex:         "0100000001a114c7deb8deba851d87755aa10aa18c97bd77afee4e1bad01d1c50e07a644eb010000006a473044022041abd4f93bd1db1d0097f2d467ae183801d7842d23d0605fa9568040d245167402201be66c96bef4d6d051304f6df2aecbdfe23a8a05af0908ef2117ab5388d8903c412103c08545a40c819f6e50892e31e792d221b6df6da96ebdba9b6fe39305cc6cc768ffffffff0263040000000000001976a91454097d9d921f9a1f55084a943571d868552e924f88acb22a0000000000001976a914c36b3fca5159231033f3fbdca1cde942096d379f88ac00000000",
					isMined:     false,
					bumpJSON:    ``,
					blockHeight: -1,
				},
			},
			receiverAddress:      "1A1PjKqjWMNBzTVdcBru27EV1PHcXWc63W",
			outputValue:          3000,
			expectedErrorMessage: "prepareBUMPFactors() error: required transactions not found in database: [ee337f429d96fb92b5227dbb9a73ecb63e1cc69d0790f58cd58ed5dc3a9ec3cf]",
		},
		{
			hexForProcessedTx: "0100000002cb3553424ffc94b59a60fb358b6cb6dfb694aee894dcd1effc0ed0a9052464e3000000006a4730440220515c3bf93d38fa7cc164746fae4bec8b66c60a82509eb553751afa5971c3e41d0220321517fd5c997ab5f8ef0e59048ce9157de46f92b10d882bf898e62f3ee7343d4121038f1273fcb299405d8d140b4de9a2111ecb39291b2846660ebecd864d13bee575ffffffff624fbcb4e68d162361f456b8b4fef6b9e7943013088b32b6bca7f5ced41ff004010000006a47304402203fb24f6e00a6487cf88a3b39d8454786db63d649142ea76374c2f55990777e6302207fbb903d038cf43e13ffb496a64f36637ec7323e5ac48bb96bdb4a885100abca4121024b003d3cf49a8f48c1fe79b711b1d08e306c42a0ab8da004d97fccc4ced3343affffffff026f000000000000001976a914f232d38cd4c2f87c117af06542b04a7061b6640188aca62a0000000000001976a9146058e52d00e3b94211939f68cc2d9a3fc1e3db0f88ac00000000",
			name:              "inputs not mined - no parents in store - should return error",
			ancestors: []*beefTestCaseAncestor{
				{
					hex:         "0100000001cfc39e3adcd58ed58cf590079dc61c3eb6ec739abb7d22b592fb969d427f33ee000000006a4730440220253e674e64028459457d55b444f5f3dc15c658425e3184c628016739e4921fd502207c8fe20eb34e55e4115fbd82c23878b4e54f01f6c6ad0811282dd0b1df863b5e41210310a4366fd997127ad972b14c56ca2e18f39ca631ac9e3e4ad3d9827865d0cc70ffffffff0264000000000000001976a914668a92ff9cb5785eb8fc044771837a0818b028b588acdc4e0000000000001976a914b073264927a61cf84327dea77414df6c28b11e5988ac00000000",
					isMined:     true,
					bumpJSON:    `{"blockHeight":"817574","path":[[{"offset":"11432","hash":"3b535e0f8e266124bce9868420052d5a7585c67e82c1edc2c7fe05fd5e140307"},{"offset":"11433","hash":"e3642405a9d00efcefd1dc94e8ae94b6dfb66c8b35fb609ab594fc4f425335cb","txid":true}],[{"offset":"5717","hash":"6ef9c6dde7fff82fa893754109f12378c8453b47dc896596b5531433093ab5b7"}],[{"offset":"2859","hash":"daa67e00ad2aef787998b66cbb3417033fbec136da1e230a5f5df3186f5c0880"}],[{"offset":"1428","hash":"bc777a80d951fbf2b7bd3a8048a9bb78fbf1d23d4127290c3fed9740b4246dd2"}],[{"offset":"715","hash":"762b57f88e7258f5757b48cda96d075cbe767c0a39a83e7109574555fd2dd8ba"}],[{"offset":"356","hash":"bbaab745bcca4f8a4be39c06c7e9be3aa1994f32271e3c6b4f768897153e5522"}],[{"offset":"179","hash":"817694ccbde5dbf88f290c30e8735991708a3d406740f7dd31434ff516a5bfde"}],[{"offset":"88","hash":"ed5b52ba4af9198d398e934a84e18405f49e7abde91cafb6dfe5aeaedb33a979"}],[{"offset":"45","hash":"0e51ec9dd5319ceb32d2d20f620c0ca3e0d918260803c1005d49e686c9b18752"}],[{"offset":"23","hash":"08ab694ef1af4019e2999a543a632cf4a662ae04d5fee879c6aadaeb749f6374"}],[{"offset":"10","hash":"4223f47597b14ee0fa7ade08e611ec80948b5fa9da267ce6c8e5d952e7fdb38e"}],[{"offset":"4","hash":"b6dace0d2294fd6e0c11f74376b7f0a1fc8ee415b350caf90c3ae92749e2a8ee"}],[{"offset":"3","hash":"795e7514ebf6d63b454d3f04854e1e0db0ac3a549f61135d5e9ef8d5785f2c68"}],[{"offset":"0","hash":"3f458f2c06493c31cbc3a035ba131913b274ac7915b9b9bc79128001a75cf76d"}],[{"offset":"1","hash":"b9b9f80cc72a674e37b54a9fdee72a9bff761f8cbcb94146afc2bffef33be89f"}]]}`,
					blockHeight: 817574,
				},
				{
					hex:             "0100000001a114c7deb8deba851d87755aa10aa18c97bd77afee4e1bad01d1c50e07a644eb010000006a473044022041abd4f93bd1db1d0097f2d467ae183801d7842d23d0605fa9568040d245167402201be66c96bef4d6d051304f6df2aecbdfe23a8a05af0908ef2117ab5388d8903c412103c08545a40c819f6e50892e31e792d221b6df6da96ebdba9b6fe39305cc6cc768ffffffff0263040000000000001976a91454097d9d921f9a1f55084a943571d868552e924f88acb22a0000000000001976a914c36b3fca5159231033f3fbdca1cde942096d379f88ac00000000",
					isMined:         false,
					bumpJSON:        ``,
					blockHeight:     -1,
					doNotAddToStore: true,
				},
			},
			receiverAddress:      "1A1PjKqjWMNBzTVdcBru27EV1PHcXWc63W",
			outputValue:          3000,
			expectedErrorMessage: "prepareBUMPFactors() error: required transactions not found in database: [04f01fd4cef5a7bcb6328b08133094e7b9f6feb4b856f46123168de6b4bc4f62]",
		},
		{
			hexForProcessedTx: "0100000002cb3553424ffc94b59a60fb358b6cb6dfb694aee894dcd1effc0ed0a9052464e3000000006a4730440220515c3bf93d38fa7cc164746fae4bec8b66c60a82509eb553751afa5971c3e41d0220321517fd5c997ab5f8ef0e59048ce9157de46f92b10d882bf898e62f3ee7343d4121038f1273fcb299405d8d140b4de9a2111ecb39291b2846660ebecd864d13bee575ffffffff624fbcb4e68d162361f456b8b4fef6b9e7943013088b32b6bca7f5ced41ff004010000006a47304402203fb24f6e00a6487cf88a3b39d8454786db63d649142ea76374c2f55990777e6302207fbb903d038cf43e13ffb496a64f36637ec7323e5ac48bb96bdb4a885100abca4121024b003d3cf49a8f48c1fe79b711b1d08e306c42a0ab8da004d97fccc4ced3343affffffff026f000000000000001976a914f232d38cd4c2f87c117af06542b04a7061b6640188aca62a0000000000001976a9146058e52d00e3b94211939f68cc2d9a3fc1e3db0f88ac00000000",
			name:              "last ancestor has corrupted hex in database - should return error",
			ancestors: []*beefTestCaseAncestor{
				{
					hex:         "0100000001cfc39e3adcd58ed58cf590079dc61c3eb6ec739abb7d22b592fb969d427f33ee000000006a4730440220253e674e64028459457d55b444f5f3dc15c658425e3184c628016739e4921fd502207c8fe20eb34e55e4115fbd82c23878b4e54f01f6c6ad0811282dd0b1df863b5e41210310a4366fd997127ad972b14c56ca2e18f39ca631ac9e3e4ad3d9827865d0cc70ffffffff0264000000000000001976a914668a92ff9cb5785eb8fc044771837a0818b028b588acdc4e0000000000001976a914b073264927a61cf84327dea77414df6c28b11e5988ac00000000",
					isMined:     true,
					bumpJSON:    `{"blockHeight":"817574","path":[[{"offset":"11432","hash":"3b535e0f8e266124bce9868420052d5a7585c67e82c1edc2c7fe05fd5e140307"},{"offset":"11433","hash":"e3642405a9d00efcefd1dc94e8ae94b6dfb66c8b35fb609ab594fc4f425335cb","txid":true}],[{"offset":"5717","hash":"6ef9c6dde7fff82fa893754109f12378c8453b47dc896596b5531433093ab5b7"}],[{"offset":"2859","hash":"daa67e00ad2aef787998b66cbb3417033fbec136da1e230a5f5df3186f5c0880"}],[{"offset":"1428","hash":"bc777a80d951fbf2b7bd3a8048a9bb78fbf1d23d4127290c3fed9740b4246dd2"}],[{"offset":"715","hash":"762b57f88e7258f5757b48cda96d075cbe767c0a39a83e7109574555fd2dd8ba"}],[{"offset":"356","hash":"bbaab745bcca4f8a4be39c06c7e9be3aa1994f32271e3c6b4f768897153e5522"}],[{"offset":"179","hash":"817694ccbde5dbf88f290c30e8735991708a3d406740f7dd31434ff516a5bfde"}],[{"offset":"88","hash":"ed5b52ba4af9198d398e934a84e18405f49e7abde91cafb6dfe5aeaedb33a979"}],[{"offset":"45","hash":"0e51ec9dd5319ceb32d2d20f620c0ca3e0d918260803c1005d49e686c9b18752"}],[{"offset":"23","hash":"08ab694ef1af4019e2999a543a632cf4a662ae04d5fee879c6aadaeb749f6374"}],[{"offset":"10","hash":"4223f47597b14ee0fa7ade08e611ec80948b5fa9da267ce6c8e5d952e7fdb38e"}],[{"offset":"4","hash":"b6dace0d2294fd6e0c11f74376b7f0a1fc8ee415b350caf90c3ae92749e2a8ee"}],[{"offset":"3","hash":"795e7514ebf6d63b454d3f04854e1e0db0ac3a549f61135d5e9ef8d5785f2c68"}],[{"offset":"0","hash":"3f458f2c06493c31cbc3a035ba131913b274ac7915b9b9bc79128001a75cf76d"}],[{"offset":"1","hash":"b9b9f80cc72a674e37b54a9fdee72a9bff761f8cbcb94146afc2bffef33be89f"}]]}`,
					blockHeight: 817574,
				},
				{
					hex:         "0100000001a114c7deb8deba851d87755aa10aa18c97bd77afee4e1bad01d1c50e07a644eb010000006a473044022041abd4f93bd1db1d0097f2d467ae183801d7842d23d0605fa9568040d245167402201be66c96bef4d6d051304f6df2aecbdfe23a8a05af0908ef2117ab5388d8903c412103c08545a40c819f6e50892e31e792d221b6df6da96ebdba9b6fe39305cc6cc768ffffffff0263040000000000001976a91454097d9d921f9a1f55084a943571d868552e924f88acb22a0000000000001976a914c36b3fca5159231033f3fbdca1cde942096d379f88ac00000000",
					isMined:     false,
					bumpJSON:    ``,
					blockHeight: -1,
					parents: []*beefTestCaseAncestor{
						{
							hex:         "010000000150965003ea3d2c08bc79b116c9ffe7e730c9f9cf0a61e3df07868b24eac6f8d3000000006b4830450221009d3489f9e76ff3b043708972c52f85519e50a5fc35563d405e04b668780bf2ba0220024188508fc9c6870b2fc4f40b9484ae4163481199a5b4a7a338b86ec8952fee4121036a8b9d796ce2dee820d1f6d7a6ba07037dab4758f16028654fe4bc3a5c430b40ffffffff022a200000000000001976a91484c73348a8fbbc44cfa34f8f5441fc104f3bc78588ac162f0000000000001976a914590b1df63948c2c4e7a12a6e52012b36e25daa9888ac00000000",
							isMined:     false,
							bumpJSON:    ``,
							blockHeight: -1,
							parents: []*beefTestCaseAncestor{
								{
									hex:         "0100000002787a565270ec00b1bf6ed20100223176656705dc0cfe5ef9d1810ca6569f12d1020000006a47304402203cfe36be7ff5c2ac939bb6a625e4a1226be242f1f9950672b5f696ec58a3358902202a48d6c6e81e5950dc49d0dd1a35b46fa8f919b109b0e7c05deaef3db6051890412102fb130326dbd7c43841cde467196e5f289b9d8596e237725df84f768468426d8bffffffff008d9db2a5c8c310e6394c24c1f3c23b3adbdd6ab4a719e917a4a0ed78768773020000006a473044022049c80385f7f69e8ba6039ebe84fe5e6578f4c3c83eb622442a96219c59ac1a750220317fe2b47838dff11f88d909732d0846eba20acff57cb357a3ff39b5a7b61b3741210322b79b40a759c485eac318eabba60a73a49ec3307ded79ba8c47204405bb2f3fffffffff05414f0000000000001976a91400414bcf2602f309171901d837b4a155adbfb5ce88ac50c30000000000001976a91489ef778cc07c77cce1ad3ff6274615afe15f20c088ac204e0000000000001976a914971b76df1dc6acf01e8e7d2f8bfb3c86e69bc64c88acef250000000000001976a9144b4a836b444d5ed8d245ddb1aa87890",
									isMined:     true,
									bumpJSON:    `{"blockHeight":"817117","path":[[{"offset":"90","hash":"d3f8c6ea248b8607dfe3610acff9c930e7e7ffc916b179bc082c3dea03509650","txid":true},{"offset":"91","hash":"5b52ad65ab613867da9a710d60898a6e5da62dea97dac25da40a0dc385253ad2"}],[{"offset":"44","hash":"84c338bea7f65ccaf7a27ca9ae6d4b11372339cf6aa6021523de3ce6f5fe4f0c"}],[{"offset":"23","hash":"5860f292e051c0a5d9d8d69a451311a009c9cde8da6522df915587913a5180dd"}],[{"offset":"10","hash":"633fb08a689363af6a8245d3482fff232b27a62b94a4d119e67700fb9608ef78"}],[{"offset":"4","hash":"4b80bb130cebb1b8c313eb4088d098178ae122fd490a255218ceada19ab9eb52"}],[{"offset":"3","hash":"cf3d0335dda3223c8b4cf28ca2c03c7e025e3088525d51981d0ee1bd2ea210cf"}],[{"offset":"0","hash":"99c7462c2530abd1be779b170b7c2afbf7b883c07175871c971734d2bd38d35b"}],[{"offset":"1","hash":"9a66c7e35426281b1be6f43ecad44a3b65a9d2234d69a55b87d535f5903d677f"}],[{"offset":"1","hash":"40e34161018499a3ad5d1ef0d74a2e557733b6c7b5c07c1d8b872ffd504e545b"}],[{"offset":"1","hash":"dce6b1d7924d3ea2b1d7e3f2c4ae15abeea3b63e1362cc0896a82bca57a21387"}],[{"offset":"1","hash":"e1fe4e2b97189cb9da328bc59430a51595ee248079bd01a3f658af810e14cb7c"}]]}`,
									blockHeight: 817117,
								},
							},
						},
					},
				},
			},
			receiverAddress:      "1A1PjKqjWMNBzTVdcBru27EV1PHcXWc63W",
			outputValue:          4500,
			expectedErrorMessage: "prepareBUMPFactors() error: required transactions not found in database: [d3f8c6ea248b8607dfe3610acff9c930e7e7ffc916b179bc082c3dea03509650]",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			ctx, client, deferMe := initSimpleTestCase(t)
			store := NewMockTransactionStore()
			defer deferMe()

			var ancestors []*Transaction
			for _, ancestor := range tc.ancestors {
				ancestors = append(ancestors, addAncestor(ctx, ancestor, client, store, t))
			}
			newTx := createProcessedTx(ctx, t, client, &tc, ancestors)

			// when
			result, err := ToBeef(ctx, newTx, store)

			// then
			assert.Equal(t, "", result)
			assert.NotNil(t, err)
			assert.Equal(t, tc.expectedErrorMessage, err.Error())
		})
	}
}

func createProcessedTx(ctx context.Context, t *testing.T, client ClientInterface, testCase *beefTestCase, ancestors []*Transaction) *Transaction {
	draftTx, err := newDraftTransaction(
		testXPub, &TransactionConfig{
			Inputs: createInputsUsingAncestors(ancestors, client),
			Outputs: []*TransactionOutput{{
				To:       testCase.receiverAddress,
				Satoshis: testCase.outputValue,
			}},
			ChangeNumberOfDestinations: 1,
			Sync: &SyncConfig{
				Broadcast:        true,
				BroadcastInstant: false,
				PaymailP2P:       false,
				SyncOnChain:      false,
			},
		},
		append(client.DefaultModelOptions(), New())...,
	)
	require.NoError(t, err)

	transaction, err := txFromHex(testCase.hexForProcessedTx, append(client.DefaultModelOptions(), New())...)
	require.NoError(t, err)

	transaction.draftTransaction = draftTx
	transaction.DraftID = draftTx.ID

	require.NotEmpty(t, transaction)

	return transaction
}

func addAncestor(ctx context.Context, testCase *beefTestCaseAncestor, client ClientInterface, store *MockTransactionStore, t *testing.T) *Transaction {
	ancestor, err := txFromHex(testCase.hex, append(client.DefaultModelOptions(), New())...)
	if err != nil {
		ancestor = emptyTx(append(client.DefaultModelOptions(), New())...)
		ancestor.Hex = testCase.hex
	}

	if testCase.isMined {
		ancestor.BlockHeight = uint64(testCase.blockHeight)

		var bump BUMP
		err := json.Unmarshal([]byte(testCase.bumpJSON), &bump)
		require.NoError(t, err)
		ancestor.BUMP = bump
	} else {
		// if we marked transaction as not mined, we need to add it's parents
		for _, parent := range testCase.parents {
			// no need a result from this func - we just want to add ancestors from level 1 and above to database if required
			_ = addAncestor(ctx, parent, client, store, t)
		}
	}

	if !testCase.doNotAddToStore {
		store.AddToStore(ancestor)
	}

	return ancestor
}

func createInputsUsingAncestors(ancestors []*Transaction, client ClientInterface) []*TransactionInput {
	var inputs []*TransactionInput

	for _, input := range ancestors {
		inputs = append(inputs, &TransactionInput{Utxo: *newUtxoFromTxID(input.GetID(), 0, append(client.DefaultModelOptions(), New())...)})
	}

	return inputs
}
