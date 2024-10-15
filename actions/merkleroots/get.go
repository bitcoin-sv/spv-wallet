package merkleroots

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/server/reqctx"
	"github.com/gin-gonic/gin"
)

// get will fetch merkleroots from Block Header Service (BHS) according to given query params
// Get Merkleroots godoc
// @Summary		Get Merkleroots
// @Description	Get Merkleroots from Block Header Service
// @Tags		Merkleroots, Block Header Service, BHS
// @Produce		json
// @Param		batchSize query int false "batch size of merkleroots to be returned"
// @Param		lastEvaluatedKey query string false "last processed merkleroot in client's database"
// @Success		200 {object} models.MerkleRootsBHSResponse "Paged response with Merkle Roots array in content "
// @Failure		400	"Bad request - batchSize must be 0 or a positive integer"
// @Failure		404	"Not found - No block with provided merkleroot was found"
// @Failure		409	"Conflict - Provided merkleroot is not part of the longest chain"
// @Failure 	500	"Internal Server Error - cannot create Block Header Service url. Please check your configuration"
// @Failure 	500	"Internal Server Error - Block Header Service cannot be requested"
// @Failure 	500	"Internal Server Error - Error while fetching Merkle Roots"
// @Router  /api/v1/merkleroots [get]
// @Security	x-auth-xpub
func get(c *gin.Context, userContext *reqctx.UserContext) {
	res, err := reqctx.Engine(c).Chain().GetMerkleRoots(c.Request.Context(), c.Request.URL.Query())

	if err != nil {
		spverrors.ErrorResponse(c, err, reqctx.Logger(c))
		return
	}

	c.JSON(http.StatusOK, res)
}
