package testabilities

import (
	"encoding/json"
	"slices"

	"github.com/bitcoin-sv/spv-wallet/models"
)

// MockedBHSMerkleRootsData is mocked  merkle roots data on Block Header Service (BHS) side
var MockedBHSMerkleRootsData = []models.MerkleRoot{
	{
		BlockHeight: 0,
		MerkleRoot:  "4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b",
	},
	{
		BlockHeight: 1,
		MerkleRoot:  "0e3e2357e806b6cdb1f70b54c3a3a17b6714ee1f0e68bebb44a74b1efd512098",
	},
	{
		BlockHeight: 2,
		MerkleRoot:  "9b0fc92260312ce44e74ef369f5c66bbb85848f2eddd5a7a1cde251e54ccfdd5",
	},
	{
		BlockHeight: 3,
		MerkleRoot:  "999e1c837c76a1b7fbb7e57baf87b309960f5ffefbf2a9b95dd890602272f644",
	},
	{
		BlockHeight: 4,
		MerkleRoot:  "df2b060fa2e5e9c8ed5eaf6a45c13753ec8c63282b2688322eba40cd98ea067a",
	},
	{
		BlockHeight: 5,
		MerkleRoot:  "63522845d294ee9b0188ae5cac91bf389a0c3723f084ca1025e7d9cdfe481ce1",
	},
	{
		BlockHeight: 6,
		MerkleRoot:  "20251a76e64e920e58291a30d4b212939aae976baca40e70818ceaa596fb9d37",
	},
	{
		BlockHeight: 7,
		MerkleRoot:  "8aa673bc752f2851fd645d6a0a92917e967083007d9c1684f9423b100540673f",
	},
	{
		BlockHeight: 8,
		MerkleRoot:  "a6f7f1c0dad0f2eb6b13c4f33de664b1b0e9f22efad5994a6d5b6086d85e85e3",
	},
	{
		BlockHeight: 9,
		MerkleRoot:  "0437cd7f8525ceed2324359c2d0ba26006d92d856a9c20fa0241106ee5a597c9",
	},
	{
		BlockHeight: 10,
		MerkleRoot:  "d3ad39fa52a89997ac7381c95eeffeaf40b66af7a57e9eba144be0a175a12b11",
	},
	{
		BlockHeight: 11,
		MerkleRoot:  "f8325d8f7fa5d658ea143629288d0530d2710dc9193ddc067439de803c37066e",
	},
	{
		BlockHeight: 12,
		MerkleRoot:  "3b96bb7e197ef276b85131afd4a09c059cc368133a26ca04ebffb0ab4f75c8b8",
	},
	{
		BlockHeight: 13,
		MerkleRoot:  "9962d5c704ec27243364cbe9d384808feeac1c15c35ac790dffd1e929829b271",
	},
	{
		BlockHeight: 14,
		MerkleRoot:  "e1afd89295b68bc5247fe0ca2885dd4b8818d7ce430faa615067d7bab8640156",
	},
}

func simulateBHSMerkleRootsAPI(lastMerkleRoot string) (string, error) {
	var response models.MerkleRootsBHSResponse
	marshallResponseError := models.SPVError{StatusCode: 500, Message: "Error during marshalling BHS response", Code: "err-marchall-bhs-res"}

	if lastMerkleRoot == "" {
		response.Content = MockedBHSMerkleRootsData
		response.Page = models.ExclusiveStartKeyPageInfo{
			LastEvaluatedKey: "",
			TotalElements:    len(MockedBHSMerkleRootsData),
			Size:             len(MockedBHSMerkleRootsData),
		}

		resString, err := json.Marshal(response)
		if err != nil {
			return "", marshallResponseError.Wrap(err)
		}

		return string(resString), nil
	}

	lastMerkleRootIdx := slices.IndexFunc(MockedBHSMerkleRootsData, func(mr models.MerkleRoot) bool {
		return mr.MerkleRoot == lastMerkleRoot
	})

	// handle case when lastMerkleRoot is already highest in the servers database
	if lastMerkleRootIdx == len(MockedBHSMerkleRootsData)-1 {
		response.Content = []models.MerkleRoot{}
		response.Page = models.ExclusiveStartKeyPageInfo{
			LastEvaluatedKey: "",
			TotalElements:    len(MockedBHSMerkleRootsData),
			Size:             0,
		}

		resString, err := json.Marshal(response)
		if err != nil {
			return "", marshallResponseError.Wrap(err)
		}

		return string(resString), nil
	}

	content := MockedBHSMerkleRootsData[lastMerkleRootIdx+1:]
	lastEvaluatedKey := content[len(content)-1].MerkleRoot

	if lastEvaluatedKey == MockedBHSMerkleRootsData[len(MockedBHSMerkleRootsData)-1].MerkleRoot {
		lastEvaluatedKey = ""
	}

	response.Content = content
	response.Page = models.ExclusiveStartKeyPageInfo{
		LastEvaluatedKey: lastEvaluatedKey,
		TotalElements:    len(MockedBHSMerkleRootsData),
		Size:             len(content),
	}

	resString, err := json.Marshal(response)
	if err != nil {
		return "", marshallResponseError.Wrap(err)
	}

	return string(resString), nil
}
