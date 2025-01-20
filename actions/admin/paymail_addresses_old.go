package admin

import (
	"net/http"
	"slices"

	"github.com/bitcoin-sv/go-paymail"
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/bitcoin-sv/spv-wallet/server/reqctx"
	"github.com/gin-gonic/gin"
)

// paymailGetAddressOld will return a paymail address
// Get Paymail godoc
// @DeprecatedRouter /v1/admin/paymail/get [post]
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
// @Deprecated
func paymailGetAddressOld(c *gin.Context, _ *reqctx.AdminContext) {
	logger := reqctx.Logger(c)
	engine := reqctx.Engine(c)
	var requestBody PaymailAddress

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest, logger)
		return
	}

	if requestBody.Address == "" {
		spverrors.ErrorResponse(c, spverrors.ErrMissingAddress, logger)
		return
	}

	opts := engine.DefaultModelOptions()

	paymailAddress, err := engine.GetPaymailAddress(c.Request.Context(), requestBody.Address, opts...)
	if err != nil {
		spverrors.ErrorResponse(c, err, logger)
		return
	}

	paymailAddressContract := mappings.MapToOldPaymailContract(paymailAddress)

	c.JSON(http.StatusOK, paymailAddressContract)
}

// paymailAddressesSearchOld will fetch a list of paymail addresses filtered by metadata
// Paymail addresses search by metadata godoc
// @DeprecatedRouter /v1/admin/paymails/search [post]
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
// @Deprecated
func paymailAddressesSearchOld(c *gin.Context, _ *reqctx.AdminContext) {
	logger := reqctx.Logger(c)
	var reqParams filter.AdminSearchPaymails
	if err := c.Bind(&reqParams); err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest, logger)
		return
	}

	paymailAddresses, err := reqctx.Engine(c).GetPaymailAddresses(
		c.Request.Context(),
		mappings.MapToMetadata(reqParams.Metadata),
		reqParams.Conditions.ToDbConditions(),
		mappings.MapToQueryParams(reqParams.QueryParams),
	)
	if err != nil {
		spverrors.ErrorResponse(c, err, logger)
		return
	}

	paymailAddressContracts := make([]*models.PaymailAddress, 0)
	for _, paymailAddress := range paymailAddresses {
		paymailAddressContracts = append(paymailAddressContracts, mappings.MapToOldPaymailContract(paymailAddress))
	}

	c.JSON(http.StatusOK, paymailAddressContracts)
}

// paymailAddressesCount will count all paymail addresses filtered by metadata
// Paymail addresses count by metadata godoc
// @DeprecatedRouter /v1/admin/paymails/count [post]
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
// @Deprecated
func paymailAddressesCount(c *gin.Context, _ *reqctx.AdminContext) {
	logger := reqctx.Logger(c)
	var reqParams filter.AdminCountPaymails
	if err := c.Bind(&reqParams); err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest, logger)
		return
	}

	count, err := reqctx.Engine(c).GetPaymailAddressesCount(
		c.Request.Context(),
		mappings.MapToMetadata(reqParams.Metadata),
		reqParams.Conditions.ToDbConditions(),
	)
	if err != nil {
		spverrors.ErrorResponse(c, err, logger)
		return
	}

	c.JSON(http.StatusOK, count)
}

// paymailCreateAddressOld will create a new paymail address
// Create Paymail godoc
// @DeprecatedRouter /v1/admin/paymail/create [post]
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
// @Deprecated
func paymailCreateAddressOld(c *gin.Context, _ *reqctx.AdminContext) {
	logger := reqctx.Logger(c)
	var requestBody CreatePaymail
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest, logger)
		return
	}

	if requestBody.Key == "" {
		spverrors.ErrorResponse(c, spverrors.ErrMissingFieldXpub, logger)
		return
	}
	if requestBody.Address == "" {
		spverrors.ErrorResponse(c, spverrors.ErrMissingAddress, logger)
		return
	}

	opts := reqctx.Engine(c).DefaultModelOptions()

	if requestBody.Metadata != nil {
		opts = append(opts, engine.WithMetadatas(requestBody.Metadata))
	}

	config := reqctx.AppConfig(c)
	if config.Paymail.DomainValidationEnabled {
		_, actualDomain, _ := paymail.SanitizePaymail(requestBody.Address)
		if !slices.Contains(config.Paymail.Domains, actualDomain) {
			spverrors.ErrorResponse(c, spverrors.ErrInvalidDomain, logger)
			return
		}
	}

	var paymailAddress *engine.PaymailAddress
	paymailAddress, err := reqctx.Engine(c).NewPaymailAddress(
		c.Request.Context(), requestBody.Key, requestBody.Address, requestBody.PublicName, requestBody.Avatar, opts...)
	if err != nil {
		spverrors.ErrorResponse(c, err, logger)
		return
	}

	paymailAddressContract := mappings.MapToOldPaymailContract(paymailAddress)

	c.JSON(http.StatusCreated, paymailAddressContract)
}

// paymailDeleteAddressOld will delete a paymail address
// Delete Paymail godoc
// @DeprecatedRouter /v1/admin/paymail/delete [delete]
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
// @Deprecated
func paymailDeleteAddressOld(c *gin.Context, _ *reqctx.AdminContext) {
	logger := reqctx.Logger(c)
	engine := reqctx.Engine(c)
	var requestBody PaymailAddress
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest, logger)
		return
	}

	if requestBody.Address == "" {
		spverrors.ErrorResponse(c, spverrors.ErrMissingAddress, logger)
		return
	}

	opts := engine.DefaultModelOptions()

	// Delete a new paymail address
	err := engine.DeletePaymailAddress(c.Request.Context(), requestBody.Address, opts...)
	if err != nil {
		spverrors.ErrorResponse(c, err, logger)
		return
	}

	c.Status(http.StatusOK)
}
