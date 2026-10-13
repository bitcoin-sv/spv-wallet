package merkleroots

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/bsv-blockchain/spv-wallet/actions/v2/internal/mapping"
	"github.com/bsv-blockchain/spv-wallet/api"
	"github.com/bsv-blockchain/spv-wallet/engine/spverrors"
	"github.com/bsv-blockchain/spv-wallet/server/reqctx"
)

// MerkleRoots returns merkleroots from block headers service according to given query params
func (s *APIMerkleRoots) MerkleRoots(c *gin.Context, params api.MerkleRootsParams) {
	res, err := s.engine.Chain().GetMerkleRoots(c.Request.Context(), c.Request.URL.Query())
	if err != nil {
		spverrors.ErrorResponse(c, err, reqctx.Logger(c))
		return
	}

	c.JSON(http.StatusOK, mapping.MerkleRootsPagedResponse(res))
}
