package chainstate

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/libsv/go-bk/envelope"
	"github.com/libsv/go-bt/v2"
	"github.com/tonicpow/go-minercraft/v2"
	"github.com/tonicpow/go-minercraft/v2/apis/mapi"
)

var (
	minerTaal = &minercraft.Miner{
		MinerID: "030d1fe5c1b560efe196ba40540ce9017c20daa9504c4c4cec6184fc702d9f274e",
		Name:    "Taal",
	}

	minerMempool = &minercraft.Miner{
		MinerID: "03e92d3e5c3f7bd945dfbf48e7a99393b1bfb3f11f380ae30d286e7ff2aec5a270",
		Name:    "Mempool",
	}

	minerMatterPool = &minercraft.Miner{
		MinerID: "0253a9b2d017254b91704ba52aad0df5ca32b4fb5cb6b267ada6aefa2bc5833a93",
		Name:    "Matterpool",
	}

	minerGorillaPool = &minercraft.Miner{
		MinerID: "03ad780153c47df915b3d2e23af727c68facaca4facd5f155bf5018b979b9aeb83",
		Name:    "GorillaPool",
	}

	allMiners = []*minercraft.Miner{
		minerTaal,
		minerMempool,
		minerGorillaPool,
		minerMatterPool,
	}

	minerAPIs = []*minercraft.MinerAPIs{
		{
			MinerID: "03e92d3e5c3f7bd945dfbf48e7a99393b1bfb3f11f380ae30d286e7ff2aec5a270",
			APIs: []minercraft.API{
				{
					URL:  "https://merchantapi.taal.com",
					Type: minercraft.MAPI,
				},
				{
					URL:  "https://tapi.taal.com/arc",
					Type: minercraft.Arc,
				},
			},
		},
		{
			MinerID: "03e92d3e5c3f7bd945dfbf48e7a99393b1bfb3f11f380ae30d286e7ff2aec5a270",
			APIs: []minercraft.API{
				{
					Token: "561b756d12572020ea9a104c3441b71790acbbce95a6ddbf7e0630971af9424b",
					URL:   "https://www.ddpurse.com/openapi",
					Type:  minercraft.MAPI,
				},
			},
		},
		{
			MinerID: "0253a9b2d017254b91704ba52aad0df5ca32b4fb5cb6b267ada6aefa2bc5833a93",
			APIs: []minercraft.API{
				{
					URL:  "https://merchantapi.matterpool.io",
					Type: minercraft.MAPI,
				},
			},
		},
		{
			MinerID: "03ad780153c47df915b3d2e23af727c68facaca4facd5f155bf5018b979b9aeb83",
			APIs: []minercraft.API{
				{
					URL:  "https://merchantapi.gorillapool.io",
					Type: minercraft.MAPI,
				},
				{
					URL:  "https://arc.gorillapool.io",
					Type: minercraft.Arc,
				},
			},
		},
	}
)

// MinerCraftBase is a mock implementation of the minercraft.MinerCraft interface.
type MinerCraftBase struct{}

// AddMiner adds a new miner to the list of miners.
func (m *MinerCraftBase) AddMiner(miner minercraft.Miner, apis []minercraft.API) error {
	existingMiner := m.MinerByName(miner.Name)
	if existingMiner != nil {
		return fmt.Errorf("miner %s already exists", miner.Name)
	}
	// Append the new miner
	allMiners = append(allMiners, &miner)

	// Append the new miner APIs
	minerAPIs = append(minerAPIs, &minercraft.MinerAPIs{
		MinerID: miner.MinerID,
		APIs:    apis,
	})

	return nil
}

// BestQuote returns the best quote for the given fee type and amount.
func (m *MinerCraftBase) BestQuote(context.Context, string, string) (*minercraft.FeeQuoteResponse, error) {
	return nil, nil
}

// FastestQuote returns the fastest quote for the given fee type and amount.
func (m *MinerCraftBase) FastestQuote(context.Context, time.Duration) (*minercraft.FeeQuoteResponse, error) {
	return nil, nil
}

// FeeQuote returns a fee quote for the given miner.
func (m *MinerCraftBase) FeeQuote(context.Context, *minercraft.Miner) (*minercraft.FeeQuoteResponse, error) {
	return &minercraft.FeeQuoteResponse{
		Quote: &mapi.FeePayload{
			Fees: []*bt.Fee{
				{
					FeeType:   bt.FeeTypeData,
					MiningFee: bt.FeeUnit(*MockDefaultFee),
				},
			},
		},
	}, nil
}

// MinerByID returns a miner by its ID.
func (m *MinerCraftBase) MinerByID(minerID string) *minercraft.Miner {
	for index, miner := range allMiners {
		if strings.EqualFold(minerID, miner.MinerID) {
			return allMiners[index]
		}
	}
	return nil
}

// MinerByName returns a miner by its name.
func (m *MinerCraftBase) MinerByName(name string) *minercraft.Miner {
	for index, miner := range allMiners {
		if strings.EqualFold(name, miner.Name) {
			return allMiners[index]
		}
	}
	return nil
}

// Miners returns all miners.
func (m *MinerCraftBase) Miners() []*minercraft.Miner {
	return allMiners
}

// MinerUpdateToken updates the token for the given miner.
func (m *MinerCraftBase) MinerUpdateToken(name, token string, apiType minercraft.APIType) {
	if miner := m.MinerByName(name); miner != nil {
		api, _ := m.MinerAPIByMinerID(miner.MinerID, apiType)
		api.Token = token
	}
}

// PolicyQuote returns a policy quote for the given miner.
func (m *MinerCraftBase) PolicyQuote(context.Context, *minercraft.Miner) (*minercraft.PolicyQuoteResponse, error) {
	return nil, nil
}

// QueryTransaction returns a transaction for the given miner.
func (m *MinerCraftBase) QueryTransaction(context.Context, *minercraft.Miner, string, ...minercraft.QueryTransactionOptFunc) (*minercraft.QueryTransactionResponse, error) {
	return nil, nil
}

// RemoveMiner removes a miner from the list of miners.
func (m *MinerCraftBase) RemoveMiner(miner *minercraft.Miner) bool {
	for i, minerFound := range allMiners {
		if miner.Name == minerFound.Name || miner.MinerID == minerFound.MinerID {
			allMiners[i] = allMiners[len(allMiners)-1]
			allMiners = allMiners[:len(allMiners)-1]
			return true
		}
	}
	// Miner not found
	return false
}

// SubmitTransaction submits a transaction to the given miner.
func (m *MinerCraftBase) SubmitTransaction(context.Context, *minercraft.Miner, *minercraft.Transaction) (*minercraft.SubmitTransactionResponse, error) {
	return nil, nil
}

// SubmitTransactions submits transactions to the given miner.
func (m *MinerCraftBase) SubmitTransactions(context.Context, *minercraft.Miner, []minercraft.Transaction) (*minercraft.SubmitTransactionsResponse, error) {
	return nil, nil
}

// APIType will return the API type
func (m *MinerCraftBase) APIType() minercraft.APIType {
	return minercraft.MAPI
}

// MinerAPIByMinerID will return a miner's API given a miner id and API type
func (m *MinerCraftBase) MinerAPIByMinerID(minerID string, apiType minercraft.APIType) (*minercraft.API, error) {
	for _, minerAPI := range minerAPIs {
		if minerAPI.MinerID == minerID {
			for i := range minerAPI.APIs {
				if minerAPI.APIs[i].Type == apiType {
					return &minerAPI.APIs[i], nil
				}
			}
		}
	}
	return nil, &minercraft.APINotFoundError{MinerID: minerID, APIType: apiType}
}

// MinerAPIsByMinerID will return a miner's APIs given a miner id
func (m *MinerCraftBase) MinerAPIsByMinerID(minerID string) *minercraft.MinerAPIs {
	for _, minerAPIs := range minerAPIs {
		if minerAPIs.MinerID == minerID {
			return minerAPIs
		}
	}
	return nil
}

// UserAgent returns the user agent.
func (m *MinerCraftBase) UserAgent() string {
	return "default-user-agent"
}

type minerCraftTxOnChain struct {
	MinerCraftBase
}

// SubmitTransaction submits a transaction to the given miner.
func (m *minerCraftTxOnChain) SubmitTransaction(_ context.Context, miner *minercraft.Miner,
	_ *minercraft.Transaction,
) (*minercraft.SubmitTransactionResponse, error) {
	if miner.Name == minercraft.MinerTaal {
		sig := "30440220008615778c5b8610c29b12925c8eb479f692ad6de9e62b7e622a3951baf9fbd8022014aaa27698cd3aba4144bfd707f3323e12ac20101d6e44f22eb8ed0856ef341a"
		pubKey := miner.MinerID
		return &minercraft.SubmitTransactionResponse{
			JSONEnvelope: minercraft.JSONEnvelope{
				Miner:     miner,
				Validated: true,
				JSONEnvelope: envelope.JSONEnvelope{
					Payload:   "{\"apiVersion\":\"1.4.0\",\"timestamp\":\"2022-02-01T15:19:40.889523Z\",\"txid\":\"683e11d4db8a776e293dc3bfe446edf66cf3b145a6ec13e1f5f1af6bb5855364\",\"returnResult\":\"failure\",\"resultDescription\":\"Missing inputs\",\"minerId\":\"030d1fe5c1b560efe196ba40540ce9017c20daa9504c4c4cec6184fc702d9f274e\",\"currentHighestBlockHash\":\"00000000000000000652def5827ad3de6380376f8fc8d3e835503095a761e0d2\",\"currentHighestBlockHeight\":724807,\"txSecondMempoolExpiry\":0}",
					Signature: &sig,
					PublicKey: &pubKey,
					Encoding:  utf8Type,
					MimeType:  applicationJSONType,
				},
			},
			Results: &minercraft.UnifiedSubmissionPayload{
				APIVersion:                "1.4.0",
				CurrentHighestBlockHash:   "00000000000000000652def5827ad3de6380376f8fc8d3e835503095a761e0d2",
				CurrentHighestBlockHeight: 724807,
				MinerID:                   miner.MinerID,
				ResultDescription:         "Missing inputs",
				ReturnResult:              mAPIFailure,
				Timestamp:                 "2022-02-01T15:19:40.889523Z",
				TxID:                      onChainExample1TxID,
			},
		}, nil
	} else if miner.Name == minercraft.MinerMempool {
		return &minercraft.SubmitTransactionResponse{
			JSONEnvelope: minercraft.JSONEnvelope{
				Miner:     miner,
				Validated: true,
				JSONEnvelope: envelope.JSONEnvelope{
					Payload:  "{\"apiVersion\":\"\",\"timestamp\":\"2022-02-01T17:47:52.518Z\",\"txid\":\"\",\"returnResult\":\"failure\",\"resultDescription\":\"ERROR: Missing inputs\",\"minerId\":null,\"currentHighestBlockHash\":\"0000000000000000064c900b1fceb316302426aedb2242852530b5e78144f2c1\",\"currentHighestBlockHeight\":724816,\"txSecondMempoolExpiry\":0}",
					Encoding: utf8Type,
					MimeType: applicationJSONType,
				},
			},
			Results: &minercraft.UnifiedSubmissionPayload{
				APIVersion:                "",
				CurrentHighestBlockHash:   "0000000000000000064c900b1fceb316302426aedb2242852530b5e78144f2c1",
				CurrentHighestBlockHeight: 724816,
				MinerID:                   miner.MinerID,
				ResultDescription:         "ERROR: Missing inputs",
				ReturnResult:              mAPIFailure,
				Timestamp:                 "2022-02-01T17:47:52.518Z",
				TxID:                      "",
			},
		}, nil
	} else if miner.Name == minercraft.MinerMatterpool {
		sig := matterCloudSig1
		pubKey := miner.MinerID
		return &minercraft.SubmitTransactionResponse{
			JSONEnvelope: minercraft.JSONEnvelope{
				Miner:     miner,
				Validated: true,
				JSONEnvelope: envelope.JSONEnvelope{
					Payload:   "{\"apiVersion\":\"1.1.0-1-g35ba2d3\",\"timestamp\":\"2022-02-01T17:50:15.130Z\",\"txid\":\"\",\"returnResult\":\"failure\",\"resultDescription\":\"ERROR: Missing inputs\",\"minerId\":\"0253a9b2d017254b91704ba52aad0df5ca32b4fb5cb6b267ada6aefa2bc5833a93\",\"currentHighestBlockHash\":\"0000000000000000064c900b1fceb316302426aedb2242852530b5e78144f2c1\",\"currentHighestBlockHeight\":724816,\"txSecondMempoolExpiry\":0}",
					Signature: &sig,
					PublicKey: &pubKey,
					Encoding:  utf8Type,
					MimeType:  applicationJSONType,
				},
			},
			Results: &minercraft.UnifiedSubmissionPayload{
				APIVersion:                "1.1.0-1-g35ba2d3",
				CurrentHighestBlockHash:   "0000000000000000064c900b1fceb316302426aedb2242852530b5e78144f2c1",
				CurrentHighestBlockHeight: 724816,
				MinerID:                   miner.MinerID,
				ResultDescription:         "ERROR: Missing inputs",
				ReturnResult:              mAPIFailure,
				Timestamp:                 "2022-02-01T17:50:15.130Z",
				TxID:                      "",
			},
		}, nil
	} else if miner.Name == minercraft.MinerGorillaPool {
		sig := gorillaPoolSig1
		pubKey := miner.MinerID
		return &minercraft.SubmitTransactionResponse{
			JSONEnvelope: minercraft.JSONEnvelope{
				Miner:     miner,
				Validated: true,
				JSONEnvelope: envelope.JSONEnvelope{
					Payload:   "{\"apiVersion\":\"\",\"timestamp\":\"2022-02-01T17:52:04.405Z\",\"txid\":\"\",\"returnResult\":\"failure\",\"resultDescription\":\"ERROR: Missing inputs\",\"minerId\":\"03ad780153c47df915b3d2e23af727c68facaca4facd5f155bf5018b979b9aeb83\",\"currentHighestBlockHash\":\"0000000000000000064c900b1fceb316302426aedb2242852530b5e78144f2c1\",\"currentHighestBlockHeight\":724816,\"txSecondMempoolExpiry\":0}",
					Signature: &sig,
					PublicKey: &pubKey,
					Encoding:  utf8Type,
					MimeType:  applicationJSONType,
				},
			},
			Results: &minercraft.UnifiedSubmissionPayload{
				APIVersion:                "",
				CurrentHighestBlockHash:   "0000000000000000064c900b1fceb316302426aedb2242852530b5e78144f2c1",
				CurrentHighestBlockHeight: 724816,
				MinerID:                   miner.MinerID,
				ResultDescription:         "ERROR: Missing inputs",
				ReturnResult:              mAPIFailure,
				Timestamp:                 "2022-02-01T17:52:04.405Z",
				TxID:                      "",
			},
		}, nil
	}

	return nil, errors.New("missing miner response")
}

// QueryTransaction mocks the QueryTransaction method of the minercraft API.
func (m *minerCraftTxOnChain) QueryTransaction(_ context.Context, miner *minercraft.Miner,
	txID string, _ ...minercraft.QueryTransactionOptFunc,
) (*minercraft.QueryTransactionResponse, error) {
	if txID == onChainExample1TxID && miner.Name == minerTaal.Name {
		sig := "304402207ede387e82db1ac38e4286b0a967b4fe1c8446c413b3785ccf86b56009439b39022043931eae02d7337b039f109be41dbd44d0472abd10ed78d7e434824ea8ab01da"
		pubKey := minerTaal.MinerID
		return &minercraft.QueryTransactionResponse{
			JSONEnvelope: minercraft.JSONEnvelope{
				Miner:     minerTaal,
				Validated: true,
				JSONEnvelope: envelope.JSONEnvelope{
					Payload:   "{\"apiVersion\":\"1.4.0\",\"timestamp\":\"2022-01-23T19:42:18.6860061Z\",\"txid\":\"908c26f8227fa99f1b26f99a19648653a1382fb3b37b03870e9c138894d29b3b\",\"returnResult\":\"success\",\"blockHash\":\"0000000000000000015122781ab51d57b26a09518630b882f67f1b08d841979d\",\"blockHeight\":723229,\"confirmations\":319,\"minerId\":\"030d1fe5c1b560efe196ba40540ce9017c20daa9504c4c4cec6184fc702d9f274e\",\"txSecondMempoolExpiry\":0}",
					Signature: &sig,
					PublicKey: &pubKey,
					Encoding:  utf8Type,
					MimeType:  applicationJSONType,
				},
			},
			Query: &minercraft.QueryTxResponse{
				APIVersion:            "1.4.0",
				Timestamp:             "2022-01-23T19:42:18.6860061Z",
				TxID:                  onChainExample1TxID,
				ReturnResult:          mAPISuccess,
				ResultDescription:     "",
				BlockHash:             onChainExample1BlockHash,
				BlockHeight:           onChainExample1BlockHeight,
				MinerID:               minerTaal.MinerID,
				Confirmations:         onChainExample1Confirmations,
				TxSecondMempoolExpiry: 0,
			},
		}, nil
	} else if txID == onChainExample1TxID && miner.Name == minerMempool.Name {
		return &minercraft.QueryTransactionResponse{
			JSONEnvelope: minercraft.JSONEnvelope{
				Miner:     minerMempool,
				Validated: false,
				JSONEnvelope: envelope.JSONEnvelope{
					Payload:   "{\"apiVersion\":\"\",\"timestamp\":\"2022-01-23T19:51:10.046Z\",\"txid\":\"908c26f8227fa99f1b26f99a19648653a1382fb3b37b03870e9c138894d29b3b\",\"returnResult\":\"success\",\"resultDescription\":\"\",\"blockHash\":\"0000000000000000015122781ab51d57b26a09518630b882f67f1b08d841979d\",\"blockHeight\":723229,\"confirmations\":321,\"minerId\":null,\"txSecondMempoolExpiry\":0}",
					Signature: nil, // NOTE: missing from mempool response
					PublicKey: nil, // NOTE: missing from mempool response
					Encoding:  utf8Type,
					MimeType:  applicationJSONType,
				},
			},
			Query: &minercraft.QueryTxResponse{
				APIVersion:            "", // NOTE: missing from mempool response
				Timestamp:             "2022-01-23T19:51:10.046Z",
				TxID:                  onChainExample1TxID,
				ReturnResult:          mAPISuccess,
				ResultDescription:     "",
				BlockHash:             onChainExample1BlockHash,
				BlockHeight:           onChainExample1BlockHeight,
				MinerID:               "", // NOTE: missing from mempool response
				Confirmations:         onChainExample1Confirmations,
				TxSecondMempoolExpiry: 0,
			},
		}, nil
	}

	return nil, nil
}

type minerCraftBroadcastSuccess struct {
	MinerCraftBase
}

// SubmitTransaction mocks the SubmitTransaction method of the minercraft API.
func (m *minerCraftBroadcastSuccess) SubmitTransaction(_ context.Context, miner *minercraft.Miner,
	_ *minercraft.Transaction,
) (*minercraft.SubmitTransactionResponse, error) {
	if miner.Name == minercraft.MinerTaal {
		sig := "30440220268ad023bbe03c62a953f907f81c01754f34ffe4822bb9e89c5245613bda7b7602204c201e56b27fd044b3f8ad77ec2c24dc2b9571166a9a998c256d3cbf598fbbda"
		pubKey := miner.MinerID
		return &minercraft.SubmitTransactionResponse{
			JSONEnvelope: minercraft.JSONEnvelope{
				Miner:     miner,
				Validated: true,
				JSONEnvelope: envelope.JSONEnvelope{
					Payload:   "{\"apiVersion\":\"1.4.0\",\"timestamp\":\"2022-02-02T12:12:02.6089293Z\",\"txid\":\"15d31d00ed7533a83d7ab206115d7642812ec04a2cbae4248365febb82576ff3\",\"returnResult\":\"success\",\"resultDescription\":\"\",\"minerId\":\"030d1fe5c1b560efe196ba40540ce9017c20daa9504c4c4cec6184fc702d9f274e\",\"currentHighestBlockHash\":\"000000000000000006e6745f6a57a1da8096faf9f71dd59b2bab3f2b0219b7a0\",\"currentHighestBlockHeight\":724922,\"txSecondMempoolExpiry\":0}",
					Signature: &sig,
					PublicKey: &pubKey,
					Encoding:  utf8Type,
					MimeType:  applicationJSONType,
				},
			},
			Results: &minercraft.UnifiedSubmissionPayload{
				APIVersion:                "1.4.0",
				CurrentHighestBlockHash:   "000000000000000006e6745f6a57a1da8096faf9f71dd59b2bab3f2b0219b7a0",
				CurrentHighestBlockHeight: 724922,
				MinerID:                   miner.MinerID,
				ResultDescription:         "",
				ReturnResult:              mAPISuccess,
				Timestamp:                 "2022-02-02T12:12:02.6089293Z",
				TxID:                      broadcastExample1TxID,
			},
		}, nil
	} else if miner.Name == minercraft.MinerMempool {
		return &minercraft.SubmitTransactionResponse{
			JSONEnvelope: minercraft.JSONEnvelope{
				Miner:     miner,
				Validated: true,
				JSONEnvelope: envelope.JSONEnvelope{
					Payload:  "{\"apiVersion\":\"\",\"timestamp\":\"2022-02-02T12:12:02.6089293Z\",\"txid\":\"15d31d00ed7533a83d7ab206115d7642812ec04a2cbae4248365febb82576ff3\",\"returnResult\":\"success\",\"resultDescription\":\"\",\"minerId\":\"03e92d3e5c3f7bd945dfbf48e7a99393b1bfb3f11f380ae30d286e7ff2aec5a270\",\"currentHighestBlockHash\":\"000000000000000006e6745f6a57a1da8096faf9f71dd59b2bab3f2b0219b7a0\",\"currentHighestBlockHeight\":724922,\"txSecondMempoolExpiry\":0}",
					Encoding: utf8Type,
					MimeType: applicationJSONType,
				},
			},
			Results: &minercraft.UnifiedSubmissionPayload{
				APIVersion:                "",
				CurrentHighestBlockHash:   "000000000000000006e6745f6a57a1da8096faf9f71dd59b2bab3f2b0219b7a0",
				CurrentHighestBlockHeight: 724922,
				MinerID:                   miner.MinerID,
				ResultDescription:         "",
				ReturnResult:              mAPISuccess,
				Timestamp:                 "2022-02-01T17:47:52.518Z",
				TxID:                      broadcastExample1TxID,
			},
		}, nil
	} else if miner.Name == minercraft.MinerMatterpool {
		sig := matterCloudSig1
		pubKey := miner.MinerID
		return &minercraft.SubmitTransactionResponse{
			JSONEnvelope: minercraft.JSONEnvelope{
				Miner:     miner,
				Validated: true,
				JSONEnvelope: envelope.JSONEnvelope{
					Payload:   "{\"apiVersion\":\"1.1.0-1-g35ba2d3\",\"timestamp\":\"2022-02-02T12:12:02.6089293Z\",\"txid\":\"\",\"returnResult\":\"success\",\"resultDescription\":\"\",\"minerId\":\"0253a9b2d017254b91704ba52aad0df5ca32b4fb5cb6b267ada6aefa2bc5833a93\",\"currentHighestBlockHash\":\"000000000000000006e6745f6a57a1da8096faf9f71dd59b2bab3f2b0219b7a0\",\"currentHighestBlockHeight\":724922,\"txSecondMempoolExpiry\":0}",
					Signature: &sig,
					PublicKey: &pubKey,
					Encoding:  utf8Type,
					MimeType:  applicationJSONType,
				},
			},
			Results: &minercraft.UnifiedSubmissionPayload{
				APIVersion:                "1.1.0-1-g35ba2d3",
				CurrentHighestBlockHash:   "000000000000000006e6745f6a57a1da8096faf9f71dd59b2bab3f2b0219b7a0",
				CurrentHighestBlockHeight: 724922,
				MinerID:                   miner.MinerID,
				ResultDescription:         "",
				ReturnResult:              mAPISuccess,
				Timestamp:                 "2022-02-02T12:12:02.6089293Z",
				TxID:                      broadcastExample1TxID,
			},
		}, nil
	} else if miner.Name == minercraft.MinerGorillaPool {
		sig := gorillaPoolSig1
		pubKey := miner.MinerID
		return &minercraft.SubmitTransactionResponse{
			JSONEnvelope: minercraft.JSONEnvelope{
				Miner:     miner,
				Validated: true,
				JSONEnvelope: envelope.JSONEnvelope{
					Payload:   "{\"apiVersion\":\"\",\"timestamp\":\"2022-02-02T12:12:02.6089293Z\",\"txid\":\"\",\"returnResult\":\"success\",\"resultDescription\":\"\",\"minerId\":\"03ad780153c47df915b3d2e23af727c68facaca4facd5f155bf5018b979b9aeb83\",\"currentHighestBlockHash\":\"000000000000000006e6745f6a57a1da8096faf9f71dd59b2bab3f2b0219b7a0\",\"currentHighestBlockHeight\":724922,\"txSecondMempoolExpiry\":0}",
					Signature: &sig,
					PublicKey: &pubKey,
					Encoding:  utf8Type,
					MimeType:  applicationJSONType,
				},
			},
			Results: &minercraft.UnifiedSubmissionPayload{
				APIVersion:                "",
				CurrentHighestBlockHash:   "000000000000000006e6745f6a57a1da8096faf9f71dd59b2bab3f2b0219b7a0",
				CurrentHighestBlockHeight: 724922,
				MinerID:                   miner.MinerID,
				ResultDescription:         "",
				ReturnResult:              mAPISuccess,
				Timestamp:                 "2022-02-02T12:12:02.6089293Z",
				TxID:                      broadcastExample1TxID,
			},
		}, nil
	}

	return nil, errors.New("missing miner response")
}

type minerCraftUnreachable struct {
	MinerCraftBase
}

// FeeQuote returns an error.
func (m *minerCraftUnreachable) FeeQuote(context.Context, *minercraft.Miner) (*minercraft.FeeQuoteResponse, error) {
	return nil, errors.New("minercraft is unreachable")
}
