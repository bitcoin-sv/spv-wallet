package mapping

import (
	"github.com/bitcoin-sv/spv-wallet/api"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/samber/lo"
)

// MerkleRootsPagedResponse maps a paged result of contacts to a response.
func MerkleRootsPagedResponse(merkleRoots *models.ExclusiveStartKeyPage[[]models.MerkleRoot]) api.ModelsGetMerkleRootResult {
	return api.ModelsGetMerkleRootResult{
		Page: api.ModelsExclusiveStartKeySearchPage{
			Size:             merkleRoots.Page.Size,
			LastEvaluatedKey: merkleRoots.Page.LastEvaluatedKey,
			TotalElements:    merkleRoots.Page.TotalElements,
		},
		Content: lo.Map(merkleRoots.Content, ModelsMerkleRoot),
	}
}

// ModelsMerkleRoot maps a merkle root model to a response.
func ModelsMerkleRoot(merkleRoot models.MerkleRoot, _ int) api.ModelsMerkleRoot {
	return api.ModelsMerkleRoot{
		MerkleRoot:  merkleRoot.MerkleRoot,
		BlockHeight: merkleRoot.BlockHeight,
	}
}
