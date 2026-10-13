package contacts

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/bsv-blockchain/spv-wallet/engine/spverrors"
	"github.com/bsv-blockchain/spv-wallet/server/reqctx"
)

// AcceptInvitation accepts a contact created by invitation.
func (s *APIContacts) AcceptInvitation(c *gin.Context, paymail string) {
	userContext := reqctx.GetUserContext(c)
	userID, err := userContext.ShouldGetUserID()
	if err != nil {
		spverrors.ErrorResponse(c, err, s.logger)
		return
	}

	err = s.contactsService.AcceptContact(c.Request.Context(), userID, paymail)
	if err != nil {
		spverrors.ErrorResponse(c, err, s.logger)
		return
	}

	c.Status(http.StatusOK)
}
