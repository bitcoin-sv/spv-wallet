package contacts

import (
	"errors"
	"net/http"

	"github.com/bitcoin-sv/go-paymail"
	"github.com/bitcoin-sv/spv-wallet/actions/v2/internal/mapping"
	"github.com/bitcoin-sv/spv-wallet/api"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/contacts/contactsmodels"
	"github.com/bitcoin-sv/spv-wallet/server/reqctx"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// UpsertContact creates new contact or updates existing one.
func (s *APIContacts) UpsertContact(c *gin.Context, paymail string) {
	var requestBody api.RequestsUpsertContact
	err := c.ShouldBindWith(&requestBody, binding.JSON)
	if err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest.Wrap(err), s.logger)
		return
	}

	userContext := reqctx.GetUserContext(c)
	userID, err := userContext.ShouldGetUserID()
	if err != nil {
		spverrors.ErrorResponse(c, err, s.logger)
		return
	}

	err = validatePaymail(paymail)
	if err != nil {
		spverrors.ErrorResponse(c, err, s.logger)
		return
	}

	if requestBody.FullName == "" {
		spverrors.ErrorResponse(c, spverrors.ErrContactFullNameRequired, s.logger)
		return
	}

	newContact := contactsmodels.NewContact{
		FullName:          requestBody.FullName,
		NewContactPaymail: paymail,
		RequesterPaymail:  requestBody.RequesterPaymail,
		UserID:            userID,
	}

	contact, err := s.contactsService.UpsertContact(c, newContact)
	if err != nil && !errors.Is(err, spverrors.ErrAddingContactRequest) {
		spverrors.ErrorResponse(c, err, s.logger)
		return
	}

	res := mapping.MapToContactContract(contact)
	c.JSON(http.StatusOK, res)
}

func validatePaymail(paymailAddress string) error {
	_, _, sanitized := paymail.SanitizePaymail(paymailAddress)
	if sanitized == "" {
		return spverrors.ErrContactInvalidPaymail
	}
	return nil
}
