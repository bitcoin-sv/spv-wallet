package contacts

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/bsv-blockchain/spv-wallet/engine/spverrors"
	"github.com/bsv-blockchain/spv-wallet/server/reqctx"
)

// UnconfirmContact unconfirms a contact.
func (s *APIContacts) UnconfirmContact(c *gin.Context, paymail string) {
	userContext := reqctx.GetUserContext(c)
	userID, err := userContext.ShouldGetUserID()
	if err != nil {
		spverrors.ErrorResponse(c, err, s.logger)
		return
	}

	err = s.contactsService.UnconfirmContact(c.Request.Context(), userID, paymail)
	if err != nil {
		spverrors.ErrorResponse(c, err, s.logger)
		return
	}

	c.Status(http.StatusOK)
}
