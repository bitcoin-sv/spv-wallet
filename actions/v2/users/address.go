package users

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/api"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/addresses/addressesmodels"
	"github.com/bitcoin-sv/spv-wallet/server/reqctx"
	"github.com/gin-gonic/gin"
)

// CreateAddress creates a new paymail address for the user.
func (s *APIUsers) CreateAddress(c *gin.Context) {
	userContext := reqctx.GetUserContext(c)
	userID, err := userContext.ShouldGetUserID()
	if err != nil {
		spverrors.ErrorResponse(c, err, reqctx.Logger(c))
		return
	}

	var req api.RequestsCreatePaymailAddress
	if err := c.ShouldBindJSON(&req); err != nil {
		spverrors.ErrorResponse(c, err, reqctx.Logger(c))
		return
	}

	fullAddress := req.Paymail.Alias + "@" + req.Paymail.Domain

	newAddress := &addressesmodels.NewAddress{
		UserID:             userID,
		Address:            fullAddress,
		CustomInstructions: nil,
	}

	err = reqctx.Engine(c).AddressesService().Create(c.Request.Context(), newAddress)
	if err != nil {
		spverrors.ErrorResponse(c, err, reqctx.Logger(c))
		return
	}

	c.JSON(http.StatusOK, &api.ResponsesPaymailAddress{
		Alias:      req.Paymail.Alias,
		Domain:     req.Paymail.Domain,
		Paymail:    fullAddress,
		PublicName: valueOrDefault(req.Paymail.PublicName, req.Paymail.Alias),
		Avatar:     valueOrDefault(req.Paymail.AvatarURL, ""),
	})
}

func valueOrDefault(ptr *string, defaultValue string) string {
	if ptr != nil {
		return *ptr
	}
	return defaultValue
}
