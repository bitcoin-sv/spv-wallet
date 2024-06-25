package admin

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/gin-gonic/gin"
)

// paymailGetAddress will return a paymail address
// Get Paymail godoc
// @Summary		Get paymail
// @Description	Get paymail
// @Tags		Admin
// @Produce		json
// @Param		PaymailAddress body PaymailAddress false "PaymailAddress model containing paymail address to get"
// @Success		200	{object} models.PaymailAddress "PaymailAddress with given address"
// @Failure		400	"Bad request - Error while parsing PaymailAddress from request body"
// @Failure 	500	"Internal Server Error - Error while getting paymail address"
// @Router		/v1/admin/paymail/get [post]
// @Security	x-auth-xpub
func (a *Action) paymailGetAddress(c *gin.Context) {
	var requestBody PaymailAddress

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest, a.Services.Logger)
		return
	}

	if requestBody.Address == "" {
		c.JSON(http.StatusBadRequest, "address is required")
		return
	}

	opts := a.Services.SpvWalletEngine.DefaultModelOptions()

	paymailAddress, err := a.Services.SpvWalletEngine.GetPaymailAddress(c.Request.Context(), requestBody.Address, opts...)
	if err != nil {
		spverrors.ErrorResponse(c, err, a.Services.Logger)
		return
	}

	paymailAddressContract := mappings.MapToPaymailContract(paymailAddress)

	c.JSON(http.StatusOK, paymailAddressContract)
}

// paymailAddressesSearch will fetch a list of paymail addresses filtered by metadata
// Paymail addresses search by metadata godoc
// @Summary		Paymail addresses search
// @Description	Paymail addresses search
// @Tags		Admin
// @Produce		json
// @Param		SearchPaymails body filter.AdminSearchPaymails false "Supports targeted resource searches with filters and metadata, plus options for pagination and sorting to streamline data exploration and analysis"
// @Success		200 {object} []models.PaymailAddress "List of paymail addresses
// @Failure		400	"Bad request - Error while parsing SearchPaymails from request body"
// @Failure 	500	"Internal server error - Error while searching for paymail addresses"
// @Router		/v1/admin/paymails/search [post]
// @Security	x-auth-xpub
func (a *Action) paymailAddressesSearch(c *gin.Context) {
	var reqParams filter.AdminSearchPaymails
	if err := c.Bind(&reqParams); err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest, a.Services.Logger)
		return
	}

	paymailAddresses, err := a.Services.SpvWalletEngine.GetPaymailAddresses(
		c.Request.Context(),
		mappings.MapToMetadata(reqParams.Metadata),
		reqParams.Conditions.ToDbConditions(),
		mappings.MapToQueryParams(reqParams.QueryParams),
	)
	if err != nil {
		spverrors.ErrorResponse(c, err, a.Services.Logger)
		return
	}

	paymailAddressContracts := make([]*models.PaymailAddress, 0)
	for _, paymailAddress := range paymailAddresses {
		paymailAddressContracts = append(paymailAddressContracts, mappings.MapToPaymailContract(paymailAddress))
	}

	c.JSON(http.StatusOK, paymailAddressContracts)
}

// paymailAddressesCount will count all paymail addresses filtered by metadata
// Paymail addresses count by metadata godoc
// @Summary		Paymail addresses count
// @Description	Paymail addresses count
// @Tags		Admin
// @Produce		json
// @Param		CountPaymails body filter.AdminCountPaymails false "Enables filtering of elements to be counted"
// @Success		200	{number} int64 "Count of paymail addresses"
// @Failure		400	"Bad request - Error while parsing CountPaymails from request body"
// @Failure 	500	"Internal Server Error - Error while fetching count of paymail addresses"
// @Router		/v1/admin/paymails/count [post]
// @Security	x-auth-xpub
func (a *Action) paymailAddressesCount(c *gin.Context) {
	var reqParams filter.AdminCountPaymails
	if err := c.Bind(&reqParams); err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest, a.Services.Logger)
		return
	}

	count, err := a.Services.SpvWalletEngine.GetPaymailAddressesCount(
		c.Request.Context(),
		mappings.MapToMetadata(reqParams.Metadata),
		reqParams.Conditions.ToDbConditions(),
	)
	if err != nil {
		spverrors.ErrorResponse(c, err, a.Services.Logger)
		return
	}

	c.JSON(http.StatusOK, count)
}

// paymailCreateAddress will create a new paymail address
// Create Paymail godoc
// @Summary		Create paymail
// @Description	Create paymail
// @Tags		Admin
// @Produce		json
// @Param		CreatePaymail body CreatePaymail false " "
// @Success		201	{object} models.PaymailAddress "Created PaymailAddress"
// @Failure		400	"Bad request - Error while parsing CreatePaymail from request body or if xpub or address are missing"
// @Failure 	500	"Internal Server Error - Error while creating new paymail address"
// @Router		/v1/admin/paymail/create [post]
// @Security	x-auth-xpub
func (a *Action) paymailCreateAddress(c *gin.Context) {
	var requestBody CreatePaymail
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest, a.Services.Logger)
		return
	}

	if requestBody.Key == "" {
		spverrors.ErrorResponse(c, spverrors.ErrMissingFieldXpub, a.Services.Logger)
		return
	}
	if requestBody.Address == "" {
		spverrors.ErrorResponse(c, spverrors.ErrMissingAddress, a.Services.Logger)
		return
	}

	opts := a.Services.SpvWalletEngine.DefaultModelOptions()

	if requestBody.Metadata != nil {
		opts = append(opts, engine.WithMetadatas(requestBody.Metadata))
	}

	var paymailAddress *engine.PaymailAddress
	paymailAddress, err := a.Services.SpvWalletEngine.NewPaymailAddress(
		c.Request.Context(), requestBody.Key, requestBody.Address, requestBody.PublicName, requestBody.Avatar, opts...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	paymailAddressContract := mappings.MapToPaymailContract(paymailAddress)

	c.JSON(http.StatusCreated, paymailAddressContract)
}

// paymailDeleteAddress will delete a paymail address
// Delete Paymail godoc
// @Summary		Delete paymail
// @Description	Delete paymail
// @Tags		Admin
// @Produce		json
// @Param		PaymailAddress body PaymailAddress false "PaymailAddress model containing paymail address to delete"
// @Success		200
// @Failure		400	"Bad request - Error while parsing PaymailAddress from request body or if address is missing"
// @Failure 	500	"Internal Server Error - Error while deleting paymail address"
// @Router		/v1/admin/paymail/delete [delete]
// @Security	x-auth-xpub
func (a *Action) paymailDeleteAddress(c *gin.Context) {
	var requestBody PaymailAddress
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest, a.Services.Logger)
		return
	}

	if requestBody.Address == "" {
		spverrors.ErrorResponse(c, spverrors.ErrMissingAddress, a.Services.Logger)
		return
	}

	opts := a.Services.SpvWalletEngine.DefaultModelOptions()

	// Delete a new paymail address
	err := a.Services.SpvWalletEngine.DeletePaymailAddress(c.Request.Context(), requestBody.Address, opts...)
	if err != nil {
		spverrors.ErrorResponse(c, err, a.Services.Logger)
		return
	}

	c.Status(http.StatusOK)
}
