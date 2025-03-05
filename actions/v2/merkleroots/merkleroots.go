package merkleroots

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/actions/v2/internal/mapping"
	"github.com/bitcoin-sv/spv-wallet/api"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/server/reqctx"
	"github.com/gin-gonic/gin"
)

// MerkleRoots returns merkleroots from block headers service according to given query params
func (s *APIMerkleRoots) MerkleRoots(c *gin.Context, params api.MerkleRootsParams) {
	res, err := s.engine.Chain().GetMerkleRoots(c, c.Request.URL.Query())
	if err != nil {
		spverrors.ErrorResponse(c, err, reqctx.Logger(c))
		return
	}

	c.JSON(http.StatusOK, mapping.MerkleRootsPagedResponse(res))
}
