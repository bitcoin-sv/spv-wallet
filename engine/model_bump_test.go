package engine

import (
	"testing"

	"github.com/libsv/go-bc"
	"github.com/stretchr/testify/assert"
)

// TestBUMPModel_CalculateBUMP will test the method CalculateMergedBUMP()
func TestBUMPModel_CalculateBUMP(t *testing.T) {
	t.Parallel()

	t.Run("Single BUMP", func(t *testing.T) {
		// given
		bumps := []BUMP{
			{
				BlockHeight: 0,
				Path: [][]BUMPLeaf{
					{
						{
							Offset: 0,
							Hash:   "123b00", // this has to be a valid hex now
						},
						{
							Offset: 1,
							Hash:   "123b",
							TxID:   true,
						},
					},
					{
						{
							Offset: 1,
							Hash:   "123b01",
						},
					},
					{
						{
							Offset: 1,
							Hash:   "123b02",
						},
					},
					{
						{
							Offset: 1,
							Hash:   "123b03",
						},
					},
				},
			},
		}
		expectedBUMP := &BUMP{
			BlockHeight: 0,
			Path: [][]BUMPLeaf{
				{
					{
						Offset: 0,
						Hash:   "123b00",
					},
					{
						Offset: 1,
						Hash:   "123b",
						TxID:   true,
					},
				},
				{
					{
						Offset: 1,
						Hash:   "123b01",
					},
				},
				{
					{
						Offset: 1,
						Hash:   "123b02",
					},
				},
				{
					{
						Offset: 1,
						Hash:   "123b03",
					},
				},
			},
			allNodes: []map[uint64]bool{
				{
					0: true,
					1: true,
				},
				{
					1: true,
				},
				{
					1: true,
				},
				{
					1: true,
				},
			},
		}

		// when
		bump, err := CalculateMergedBUMP(bumps)

		// then
		assert.NoError(t, err)
		assert.Equal(t, expectedBUMP, bump)
	})

	t.Run("Paired Transactions", func(t *testing.T) {
		// given
		bumps := []BUMP{
			{
				BlockHeight: 0,
				Path: [][]BUMPLeaf{
					{
						{
							Offset: 8,
							Hash:   "123b09",
							TxID:   true,
						},
						{
							Offset: 9,
							Hash:   "123b10",
						},
					},
					{
						{
							Offset: 5,
							Hash:   "123b1112",
						},
					},
					{
						{
							Offset: 3,
							Hash:   "123b13141516",
						},
					},
					{
						{
							Offset: 0,
							Hash:   "123b0102030405060708",
						},
					},
				},
			},
			{
				BlockHeight: 0,
				Path: [][]BUMPLeaf{
					{
						{
							Offset: 8,
							Hash:   "123b09",
						},
						{
							Offset: 9,
							Hash:   "123b10",
							TxID:   true,
						},
					},
					{
						{
							Offset: 5,
							Hash:   "123b1112",
						},
					},
					{
						{
							Offset: 3,
							Hash:   "123b13141516",
						},
					},
					{
						{
							Offset: 0,
							Hash:   "123b0102030405060708",
						},
					},
				},
			},
		}
		expectedBUMP := &BUMP{
			BlockHeight: 0,
			Path: [][]BUMPLeaf{
				{
					{
						Offset: 8,
						Hash:   "123b09",
						TxID:   true,
					},
					{
						Offset: 9,
						Hash:   "123b10",
						TxID:   true,
					},
				},
				{
					{
						Offset: 5,
						Hash:   "123b1112",
					},
				},
				{
					{
						Offset: 3,
						Hash:   "123b13141516",
					},
				},
				{
					{
						Offset: 0,
						Hash:   "123b0102030405060708",
					},
				},
			},
			allNodes: []map[uint64]bool{
				{
					8: true,
					9: true,
				},
				{
					5: true,
				},
				{
					3: true,
				},
				{
					0: true,
				},
			},
		}

		// when
		bump, err := CalculateMergedBUMP(bumps)

		// then
		assert.NoError(t, err)
		assert.Equal(t, expectedBUMP, bump)
	})

	t.Run("Different sizes of BUMPs", func(t *testing.T) {
		// given
		bumps := []BUMP{
			{
				BlockHeight: 0,
				Path: [][]BUMPLeaf{
					{
						{
							Offset: 8,
							Hash:   "123b09",
							TxID:   true,
						},
						{
							Offset: 9,
							Hash:   "123b10",
						},
					},
					{
						{
							Offset: 5,
							Hash:   "123b1112",
						},
					},
					{
						{
							Offset: 3,
							Hash:   "123b13141516",
						},
					},
					{
						{
							Offset: 0,
							Hash:   "123b0102030405060708",
						},
					},
				},
			},
			{
				BlockHeight: 0,
				Path: [][]BUMPLeaf{
					{
						{
							Offset: 8,
							Hash:   "123b09",
						},
						{
							Offset: 9,
							Hash:   "123b10",
							TxID:   true,
						},
					},
					{
						{
							Offset: 5,
							Hash:   "123b1112",
						},
					},
					{
						{
							Offset: 3,
							Hash:   "123b0102030405060708",
						},
					},
				},
			},
		}

		// when
		bump, err := CalculateMergedBUMP(bumps)

		// then
		assert.Error(t, err)
		assert.Nil(t, bump)
	})

	t.Run("BUMPs with different block heights", func(t *testing.T) {
		// given
		bumps := []BUMP{
			{
				BlockHeight: 0,
				Path: [][]BUMPLeaf{
					{
						{
							Offset: 8,
							Hash:   "123b09",
							TxID:   true,
						},
						{
							Offset: 9,
							Hash:   "123b10",
						},
					},
					{
						{
							Offset: 5,
							Hash:   "123b1112",
						},
					},
				},
			},
			{
				BlockHeight: 100,
				Path: [][]BUMPLeaf{
					{
						{
							Offset: 8,
							Hash:   "123b09",
						},
						{
							Offset: 9,
							Hash:   "123b10",
							TxID:   true,
						},
					},
					{
						{
							Offset: 5,
							Hash:   "123b1112",
						},
					},
				},
			},
		}

		// when
		bump, err := CalculateMergedBUMP(bumps)

		// then
		assert.Error(t, err)
		assert.Nil(t, bump)
	})

	t.Run("Empty slice of BUMPS", func(t *testing.T) {
		// given
		bumps := []BUMP{}

		// when
		bump, err := CalculateMergedBUMP(bumps)

		// then
		assert.NoError(t, err)
		assert.Nil(t, bump)
	})

	t.Run("Slice of empty BUMPS", func(t *testing.T) {
		// given
		bumps := []BUMP{
			{}, {}, {},
		}
		// when
		bump, err := CalculateMergedBUMP(bumps)

		// then
		assert.Error(t, err)
		assert.Nil(t, bump)
	})
}

// TestBUMPModel_Hex will test the method Hex()
func TestBUMPModel_Hex(t *testing.T) {
	t.Run("BUMP to HEX - simple example", func(t *testing.T) {
		// given
		expectedHex := "01" + // block height
			"03" + // tree height
			// ---- LEVEL 0 -----
			"02" + // nLeafes on level 0
			"00" + // offset
			"00" + // flag - data follows, not a client txid
			"0a" + // hash
			"01" + // offset
			"02" + // flag - data follows, not a client txid
			"0b" + // hash
			// ---- LEVEL 1 -----
			"01" + // nLeafes on level 0
			"01" + // offset
			"00" + // flag - data follows, not a client txid
			"cd" + // hash
			// ---- LEVEL 2 -----
			"01" + // nLeafes on level 0
			"01" + // offset
			"00" + // flag - data follows, not a client txid
			"abef" // hash (little endian - reversed bytes)
		bump := BUMP{
			BlockHeight: 1,
			Path: [][]BUMPLeaf{
				{
					{
						Offset: 0,
						Hash:   "0a",
					},
					{
						Offset: 1,
						TxID:   true,
						Hash:   "0b",
					},
				},
				{
					{
						Offset: 1,
						Hash:   "cd",
					},
				},
				{
					{
						Offset: 1,
						Hash:   "efab",
					},
				},
			},
		}

		// when
		actualHex := bump.Hex()

		// then
		assert.Equal(t, expectedHex, actualHex)
	})

	t.Run("BUMP to HEX - standard example", func(t *testing.T) {
		// given
		expectedHex := "fe8a6a0c000c04fde80b0011774f01d26412f0d16ea3f0447be0b5ebec67b0782e321a7a01cbdf7f734e30fde90b02004e53753e3fe4667073063a17987292cfdea278824e9888e52180581d7188d8fdea0b025e441996fc53f0191d649e68a200e752fb5f39e0d5617083408fa179ddc5c998fdeb0b0102fdf405000671394f72237d08a4277f4435e5b6edf7adc272f25effef27cdfe805ce71a81fdf50500262bccabec6c4af3ed00cc7a7414edea9c5efa92fb8623dd6160a001450a528201fdfb020101fd7c010093b3efca9b77ddec914f8effac691ecb54e2c81d0ab81cbc4c4b93befe418e8501bf01015e005881826eb6973c54003a02118fe270f03d46d02681c8bc71cd44c613e86302f8012e00e07a2bb8bb75e5accff266022e1e5e6e7b4d6d943a04faadcf2ab4a22f796ff30116008120cafa17309c0bb0e0ffce835286b3a2dcae48e4497ae2d2b7ced4f051507d010a00502e59ac92f46543c23006bff855d96f5e648043f0fb87a7a5949e6a9bebae430104001ccd9f8f64f4d0489b30cc815351cf425e0e78ad79a589350e4341ac165dbe45010301010000af8764ce7e1cc132ab5ed2229a005c87201c9a5ee15c0f91dd53eff31ab30cd4"
		bump := BUMP{
			BlockHeight: 813706,
			Path: [][]BUMPLeaf{
				{
					{
						Offset: 3048,
						Hash:   "304e737fdfcb017a1a322e78b067ecebb5e07b44f0a36ed1f01264d2014f7711",
					},
					{
						Offset: 3049,
						TxID:   true,
						Hash:   "d888711d588021e588984e8278a2decf927298173a06737066e43f3e75534e00",
					},
					{
						Offset: 3050,
						TxID:   true,
						Hash:   "98c9c5dd79a18f40837061d5e0395ffb52e700a2689e641d19f053fc9619445e",
					},
					{
						Offset:    3051,
						Duplicate: true,
					},
				},
				{
					{
						Offset: 1524,
						Hash:   "811ae75c80fecd27efff5ef272c2adf7edb6e535447f27a4087d23724f397106",
					},
					{
						Offset: 1525,
						Hash:   "82520a4501a06061dd2386fb92fa5e9ceaed14747acc00edf34a6cecabcc2b26",
					},
				},
				{
					{
						Offset:    763,
						Duplicate: true,
					},
				},
				{
					{
						Offset: 380,
						Hash:   "858e41febe934b4cbc1cb80a1dc8e254cb1e69acff8e4f91ecdd779bcaefb393",
					},
				},
				{
					{
						Offset:    191,
						Duplicate: true,
					},
				},
				{
					{
						Offset: 94,
						Hash:   "f80263e813c644cd71bcc88126d0463df070e28f11023a00543c97b66e828158",
					},
				},
				{
					{
						Offset: 46,
						Hash:   "f36f792fa2b42acfadfa043a946d4d7b6e5e1e2e0266f2cface575bbb82b7ae0",
					},
				},
				{
					{
						Offset: 22,
						Hash:   "7d5051f0d4ceb7d2e27a49e448aedca2b3865283ceffe0b00b9c3017faca2081",
					},
				},
				{
					{
						Offset: 10,
						Hash:   "43aeeb9b6a9e94a5a787fbf04380645e6fd955f8bf0630c24365f492ac592e50",
					},
				},
				{
					{
						Offset: 4,
						Hash:   "45be5d16ac41430e3589a579ad780e5e42cf515381cc309b48d0f4648f9fcd1c",
					},
				},
				{
					{
						Offset:    3,
						Duplicate: true,
					},
				},
				{
					{
						Offset: 0,
						Hash:   "d40cb31af3ef53dd910f5ce15e9a1c20875c009a22d25eab32c11c7ece6487af",
					},
				},
			},
		}

		// when
		actualHex := bump.Hex()

		// then
		assert.Equal(t, expectedHex, actualHex)
	})
}

// TestBUMPModel_CalculateMergedBUMPAndHex will test both the CalculateMergedBUMP() and Hex() methods.
func TestBUMPModel_CalculateMergedBUMPAndHex(t *testing.T) {
	t.Parallel()

	t.Run("Real Merkle Proof", func(t *testing.T) {
		// given
		merkleProof := []bc.MerkleProof{
			{
				Index:  1153,
				TxOrID: "2130b63dcbfe1356a30137fe9578691f59c6cf42d5e8928a800619de7f8e14da",
				Nodes: []string{
					"4d4bde1dc35c87bba992944ec0379e0bb009916108113dc3de1c4aecda6457a3",
					"168595f83accfcec66d0e0df06df89e6a9a2eaa3aa69427fb86cb54d8ea5b1e9",
					"c2edd41b237844a45a0e6248a9e7c520af303a5c91cc8a443ad0075d6a3dec79",
					"bdd0fddf45fee49324e55dfc6fdb9044c86dc5be3dbf941a80b395838495ac09",
					"3e5ec052b86621b5691d15ad54fab2551c27a36d9ab84f428a304b607aa33d33",
					"9feb9b1aaa2cd8486edcacb60b9d477a89aec5867d292608c3c59a18324d608a",
					"22e1db219f8d874315845b7cee84832dc0865b5f9e18221a011043a4d6704e7d",
					"7f118890abd8df3f8a51c344da0f9235609f5fd380e38cfe519e81262aedb2a7",
					"20dcf60bbcecd2f587e8d3344fb68c71f2f2f7a6cc85589b9031c2312a433fe6",
					"0be65c1f3b53b937608f8426e43cb41c1db31227d0d9933e8b0ce3b8cc30d67f",
					"a8036cf77d8de296f60607862b228174733a30486a37962a56465f5e8c214d87",
					"b8e4d7975537bb775e320f01f874c06cf38dd2ce7bb836a1afe0337aeb9fb06f",
					"88e6b0bd93e02b057ea43a80a5bb8cf9673f143340af3f569fe0c55c085e5efb",
					"15f731176e17f4402802d5be3893419e690225e732d69dfd27f6e614f188233d",
				},
			},
		}
		expectedBUMP := &BUMP{
			BlockHeight: 0,
			Path: [][]BUMPLeaf{
				{
					{
						Offset: 1152,
						Hash:   "4d4bde1dc35c87bba992944ec0379e0bb009916108113dc3de1c4aecda6457a3",
					},
					{
						Offset: 1153,
						Hash:   "2130b63dcbfe1356a30137fe9578691f59c6cf42d5e8928a800619de7f8e14da",
						TxID:   true,
					},
				},
				{
					{
						Offset: 577,
						Hash:   "168595f83accfcec66d0e0df06df89e6a9a2eaa3aa69427fb86cb54d8ea5b1e9",
					},
				},
				{
					{
						Offset: 289,
						Hash:   "c2edd41b237844a45a0e6248a9e7c520af303a5c91cc8a443ad0075d6a3dec79",
					},
				},
				{
					{
						Offset: 145,
						Hash:   "bdd0fddf45fee49324e55dfc6fdb9044c86dc5be3dbf941a80b395838495ac09",
					},
				},
				{
					{
						Offset: 73,
						Hash:   "3e5ec052b86621b5691d15ad54fab2551c27a36d9ab84f428a304b607aa33d33",
					},
				},
				{
					{
						Offset: 37,
						Hash:   "9feb9b1aaa2cd8486edcacb60b9d477a89aec5867d292608c3c59a18324d608a",
					},
				},
				{
					{
						Offset: 19,
						Hash:   "22e1db219f8d874315845b7cee84832dc0865b5f9e18221a011043a4d6704e7d",
					},
				},
				{
					{
						Offset: 8,
						Hash:   "7f118890abd8df3f8a51c344da0f9235609f5fd380e38cfe519e81262aedb2a7",
					},
				},
				{
					{
						Offset: 5,
						Hash:   "20dcf60bbcecd2f587e8d3344fb68c71f2f2f7a6cc85589b9031c2312a433fe6",
					},
				},
				{
					{
						Offset: 3,
						Hash:   "0be65c1f3b53b937608f8426e43cb41c1db31227d0d9933e8b0ce3b8cc30d67f",
					},
				},
				{
					{
						Offset: 0,
						Hash:   "a8036cf77d8de296f60607862b228174733a30486a37962a56465f5e8c214d87",
					},
				},
				{
					{
						Offset: 1,
						Hash:   "b8e4d7975537bb775e320f01f874c06cf38dd2ce7bb836a1afe0337aeb9fb06f",
					},
				},
				{
					{
						Offset: 1,
						Hash:   "88e6b0bd93e02b057ea43a80a5bb8cf9673f143340af3f569fe0c55c085e5efb",
					},
				},
				{
					{
						Offset: 1,
						Hash:   "15f731176e17f4402802d5be3893419e690225e732d69dfd27f6e614f188233d",
					},
				},
			},
			allNodes: []map[uint64]bool{
				{
					1152: true,
					1153: true,
				},
				{
					577: true,
				},
				{
					289: true,
				},
				{
					145: true,
				},
				{
					73: true,
				},
				{
					37: true,
				},
				{
					19: true,
				},
				{
					8: true,
				},
				{
					5: true,
				},
				{
					3: true,
				},
				{
					0: true,
				},
				{
					1: true,
				},
				{
					1: true,
				},
				{
					1: true,
				},
			},
		}
		expectedHex := "00" + // block height (dummy value)
			"0e" + // 13 - tree height
			"02" + // nLeafs at this level
			"fd8004" + // offset - 1152
			"00" + // flags - data follows, not a client txid
			"a35764daec4a1cdec33d1108619109b00b9e37c04e9492a9bb875cc31dde4b4d" + // hash
			"fd8104" + // offset - 1153
			"02" + // flags - data follows, client txid
			"da148e7fde1906808a92e8d542cfc6591f697895fe3701a35613fecb3db63021" + // hash
			// ----------------------
			// implied end of leaves at this height
			// height of next leaves is therefore 12
			"01" +
			"fd4102" +
			"00" +
			"e9b1a58e4db56cb87f4269aaa3eaa2a9e689df06dfe0d066ecfccc3af8958516" +
			"01" +
			"fd2101" +
			"00" +
			"79ec3d6a5d07d03a448acc915c3a30af20c5e7a948620e5aa44478231bd4edc2" +
			"01" +
			"91" +
			"00" +
			"09ac95848395b3801a94bf3dbec56dc84490db6ffc5de52493e4fe45dffdd0bd" +
			"01" +
			"49" +
			"00" +
			"333da37a604b308a424fb89a6da3271c55b2fa54ad151d69b52166b852c05e3e" +
			"01" +
			"25" +
			"00" +
			"8a604d32189ac5c30826297d86c5ae897a479d0bb6acdc6e48d82caa1a9beb9f" +
			"01" +
			"13" +
			"00" +
			"7d4e70d6a44310011a22189e5f5b86c02d8384ee7c5b841543878d9f21dbe122" +
			"01" +
			"08" +
			"00" +
			"a7b2ed2a26819e51fe8ce380d35f9f6035920fda44c3518a3fdfd8ab9088117f" +
			"01" +
			"05" +
			"00" +
			"e63f432a31c231909b5885cca6f7f2f2718cb64f34d3e887f5d2ecbc0bf6dc20" +
			"01" +
			"03" +
			"00" +
			"7fd630ccb8e30c8b3e93d9d02712b31d1cb43ce426848f6037b9533b1f5ce60b" +
			"01" +
			"00" +
			"00" +
			"874d218c5e5f46562a96376a48303a737481222b860706f696e28d7df76c03a8" +
			"01" +
			"01" +
			"00" +
			"6fb09feb7a33e0afa136b87bced28df36cc074f8010f325e77bb375597d7e4b8" +
			"01" +
			"01" +
			"00" +
			"fb5e5e085cc5e09f563faf4033143f67f98cbba5803aa47e052be093bdb0e688" +
			"01" +
			"01" +
			"00" +
			"3d2388f114e6f627fd9dd632e72502699e419338bed5022840f4176e1731f715"

		// when
		bumps := make([]BUMP, 0)
		for _, mp := range merkleProof {
			bumps = append(bumps, merkleProofToBUMP(&mp, 0))
		}
		bump, err := CalculateMergedBUMP(bumps)
		actualHex := bump.Hex()

		// then
		assert.NoError(t, err)
		assert.Equal(t, expectedBUMP, bump)
		assert.Equal(t, expectedHex, actualHex)
	})
}

// TestBUMPModel_merkleProofToBUMP will test the method merkleProofToBUMP()
func TestBUMPModel_merkleProofToBUMP(t *testing.T) {
	t.Parallel()

	t.Run("Valid Merkle Proof #1", func(t *testing.T) {
		// given
		blockHeight := uint64(0)
		mp := bc.MerkleProof{
			Index:  1,
			TxOrID: "txId",
			Nodes:  []string{"node0", "node1", "node2", "node3"},
		}
		expectedBUMP := BUMP{
			BlockHeight: blockHeight,
			Path: [][]BUMPLeaf{
				{
					{Offset: 0, Hash: "node0"},
					{Offset: 1, Hash: "txId", TxID: true},
				},
				{
					{Offset: 1, Hash: "node1"},
				},
				{
					{Offset: 1, Hash: "node2"},
				},
				{
					{Offset: 1, Hash: "node3"},
				},
			},
		}

		// when
		actualBUMP := merkleProofToBUMP(&mp, blockHeight)

		// then
		assert.Equal(t, expectedBUMP, actualBUMP)
	})

	t.Run("Valid Merkle Proof #2", func(t *testing.T) {
		// given
		blockHeight := uint64(0)
		mp := bc.MerkleProof{
			Index:  14,
			TxOrID: "txId",
			Nodes:  []string{"node0", "node1", "node2", "node3", "node4"},
		}
		expectedBUMP := BUMP{
			BlockHeight: blockHeight,
			Path: [][]BUMPLeaf{
				{
					{Offset: 14, Hash: "txId", TxID: true},
					{Offset: 15, Hash: "node0"},
				},
				{
					{Offset: 6, Hash: "node1"},
				},
				{
					{Offset: 2, Hash: "node2"},
				},
				{
					{Offset: 0, Hash: "node3"},
				},
				{
					{Offset: 1, Hash: "node4"},
				},
			},
		}

		// when
		actualBUMP := merkleProofToBUMP(&mp, blockHeight)

		// then
		assert.Equal(t, expectedBUMP, actualBUMP)
	})

	t.Run("Valid Merkle Proof #3 - with *", func(t *testing.T) {
		// given
		blockHeight := uint64(0)
		mp := bc.MerkleProof{
			Index:  14,
			TxOrID: "txId",
			Nodes:  []string{"*", "node1", "node2", "node3", "node4"},
		}
		expectedBUMP := BUMP{
			BlockHeight: blockHeight,
			Path: [][]BUMPLeaf{
				{
					{Offset: 14, Hash: "txId", TxID: true},
					{Offset: 15, Duplicate: true},
				},
				{
					{Offset: 6, Hash: "node1"},
				},
				{
					{Offset: 2, Hash: "node2"},
				},
				{
					{Offset: 0, Hash: "node3"},
				},
				{
					{Offset: 1, Hash: "node4"},
				},
			},
		}

		// when
		actualBUMP := merkleProofToBUMP(&mp, blockHeight)

		// then
		assert.Equal(t, expectedBUMP, actualBUMP)
	})

	t.Run("Empty Merkle Proof", func(t *testing.T) {
		blockHeight := uint64(0)
		mp := bc.MerkleProof{}
		actualBUMP := merkleProofToBUMP(&mp, blockHeight)
		assert.Equal(t, BUMP{BlockHeight: blockHeight}, actualBUMP)
	})
}
