package admin

import (
	"net/http"

	"github.com/bitcoin-sv/go-paymail"
	"github.com/bitcoin-sv/spv-wallet/actions/common"
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/internal/query"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/bitcoin-sv/spv-wallet/models/response"
	"github.com/bitcoin-sv/spv-wallet/server/reqctx"
	"github.com/gin-gonic/gin"
)

// paymailGetAddress will return a paymail address
// Get Paymail godoc
// @Summary		Get paymail
// @Description	Get paymail
// @Tags		Admin
// @Produce		json
// @Param		PaymailAddress body PaymailAddress false "PaymailAddress model containing paymail address to get"
// @Success		200	{object} response.PaymailAddress "PaymailAddress with given address"
// @Failure		400	"Bad request - Error while parsing PaymailAddress from request body"
// @Failure 	500	"Internal Server Error - Error while getting paymail address"
// @Router		/api/v1/admin/paymails/{id} [get]
// @Security	x-auth-xpub
func paymailGetAddress(c *gin.Context, _ *reqctx.AdminContext) {
	logger := reqctx.Logger(c)
	engine := reqctx.Engine(c)
	id := c.Param("id")

	opts := engine.DefaultModelOptions()

	paymailAddress, err := engine.GetPaymailAddressByID(c.Request.Context(), id, opts...)
	if err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCouldNotFindPaymail.WithTrace(err), logger)
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
// @Param		SearchPaymails body filter.AdminPaymailFilter false "Supports targeted resource searches with filters and metadata, plus options for pagination and sorting to streamline data exploration and analysis"
// @Success		200 {object} []response.PaymailAddress "List of paymail addresses
// @Failure		400	"Bad request - Error while parsing SearchPaymails from request body"
// @Failure 	500	"Internal server error - Error while searching for paymail addresses"
// @Router		/api/v1/admin/paymails [get]
// @Security	x-auth-xpub
func paymailAddressesSearch(c *gin.Context, _ *reqctx.AdminContext) {
	logger := reqctx.Logger(c)

	searchParams, err := query.ParseSearchParams[filter.AdminPaymailFilter](c)
	if err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotParseQueryParams.WithTrace(err), logger)
		return
	}

	conditions := searchParams.Conditions.ToDbConditions()
	metadata := mappings.MapToMetadata(searchParams.Metadata)
	pageOptions := mappings.MapToDbQueryParams(&searchParams.Page)

	paymailAddresses, err := reqctx.Engine(c).GetPaymailAddresses(
		c.Request.Context(),
		metadata,
		conditions,
		pageOptions,
	)
	if err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCouldNotFindPaymail.WithTrace(err), logger)
		return
	}

	count, err := reqctx.Engine(c).GetPaymailAddressesCount(c.Request.Context(), metadata, conditions)
	if err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCouldNotFindPaymail.WithTrace(err), logger)
		return
	}

	paymailAddressContracts := common.MapToTypeContracts(paymailAddresses, mappings.MapToPaymailContract)

	result := response.PageModel[response.PaymailAddress]{
		Content: paymailAddressContracts,
		Page:    common.GetPageDescriptionFromSearchParams(pageOptions, count),
	}
	c.JSON(http.StatusOK, result)
}

// paymailCreateAddress will create a new paymail address
// Create Paymail godoc
// @Summary		Create paymail
// @Description	Create paymail
// @Tags		Admin
// @Produce		json
// @Param		CreatePaymail body CreatePaymail false " "
// @Success		201	{object} response.PaymailAddress "Created PaymailAddress"
// @Failure		400	"Bad request - Error while parsing CreatePaymail from request body or if xpub or address are missing"
// @Failure 	500	"Internal Server Error - Error while creating new paymail address"
// @Router		/api/v1/admin/paymails [post]
// @Security	x-auth-xpub
func paymailCreateAddress(c *gin.Context, _ *reqctx.AdminContext) {
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
		if !checkPaymailDomain(actualDomain, config.Paymail.Domains) {
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

	paymailAddressContract := mappings.MapToPaymailContract(paymailAddress)

	c.JSON(http.StatusCreated, paymailAddressContract)
}

func checkPaymailDomain(domain string, domains []string) bool {
	for _, d := range domains {
		if d == domain {
			return true
		}
	}
	return false
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
// @Router		/api/v1/admin/paymails/{id} [delete]
// @Security	x-auth-xpub
func paymailDeleteAddress(c *gin.Context, _ *reqctx.AdminContext) {
	logger := reqctx.Logger(c)
	engine := reqctx.Engine(c)
	var requestBody PaymailAddress
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest.WithTrace(err), logger)
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
		spverrors.ErrorResponse(c, spverrors.ErrDeletePaymailAddress.WithTrace(err), logger)
		return
	}

	c.Status(http.StatusOK)
}
