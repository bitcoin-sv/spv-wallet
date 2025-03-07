package users

import (
	"net/http"

	adminerrors "github.com/bitcoin-sv/spv-wallet/actions/v2/admin/errors"
	"github.com/bitcoin-sv/spv-wallet/actions/v2/admin/internal/mapping"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// UserById returns a user by ID
func (s *APIAdminUsers) UserById(c *gin.Context, id string) {
	user, err := s.engine.UsersService().GetByID(c, id)
	if err != nil {
		spverrors.MapResponse(c, err, s.logger).
			If(gorm.ErrRecordNotFound).Then(adminerrors.ErrUserNotFound).
			Else(adminerrors.ErrGetUserFailed)
		return
	}

	c.JSON(http.StatusOK, mapping.UserToResponse(user))
}
