package merkleroots

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/api"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/server/reqctx"
	"github.com/gin-gonic/gin"
)

func (s *APIMerkleRoots) GetMerkleRoots(c *gin.Context, params api.GetMerkleRootsParams) {
	res, err := s.engine.Chain().GetMerkleRoots(c.Request.Context(), c.Request.URL.Query())

	if err != nil {
		spverrors.ErrorResponse(c, err, reqctx.Logger(c))
		return
	}

	c.JSON(http.StatusOK, res)
}
