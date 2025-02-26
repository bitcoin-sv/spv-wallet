package contacts

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/actions/v2/internal/mapping"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/server/reqctx"
	"github.com/gin-gonic/gin"
)

func (s *APIContacts) GetContact(c *gin.Context, paymail string) {
	userContext := reqctx.GetUserContext(c)
	userID, err := userContext.ShouldGetUserID()
	if err != nil {
		spverrors.ErrorResponse(c, err, s.logger)
		return
	}

	contact, err := s.engine.ContactService().Find(c.Request.Context(), userID, paymail)
	if err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrGetContact.WithTrace(err), s.logger)
		return
	} else if contact == nil {
		spverrors.ErrorResponse(c, spverrors.ErrContactNotFound, s.logger)
		return
	}

	res := mapping.MapToContactContract(contact)

	c.JSON(http.StatusOK, res)
}
