package users

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/actions/v2/admin/internal/mapping"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/server/reqctx"
	"github.com/gin-gonic/gin"
)

func (s *APIAdminUsers) GetApiV2AdminUsersId(c *gin.Context, id string) {
	user, err := reqctx.Engine(c).UsersService().GetByID(c, id)
	if err != nil {
		spverrors.ErrorResponse(c, err, reqctx.Logger(c))
		return
	}

	c.JSON(http.StatusOK, mapping.UserToResponse(user))
}
