package admin

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/actions"
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/gin-gonic/gin"
)

// paymailGetAddress will return a paymail address
// Get Paymail godoc
// @Summary		Get paymail
// @Description	Get paymail
// @Tags		Admin
// @Param		PaymailAddress body PaymailAddress true "CreateAccessKey model containing paymail address"
// @Produce		json
// @Success		200
// @Router		/v1/admin/paymail/get [post]
// @Security	x-auth-xpub
func (a *Action) paymailGetAddress(c *gin.Context) {
	var requestBody PaymailAddress

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if requestBody.Address == "" {
		c.JSON(http.StatusExpectationFailed, "address is required")
		return
	}

	opts := a.Services.SpvWalletEngine.DefaultModelOptions()

	paymailAddress, err := a.Services.SpvWalletEngine.GetPaymailAddress(c.Request.Context(), requestBody.Address, opts...)
	if err != nil {
		c.JSON(http.StatusExpectationFailed, err.Error())
		return
	}

	c.JSON(http.StatusOK, paymailAddress)
}

// paymailAddressesSearch will fetch a list of paymail addresses filtered by metadata
// Paymail addresses search by metadata godoc
// @Summary		Paymail addresses search
// @Description	Paymail addresses search
// @Tags		Admin
// @Produce		json
// @Param		SearchRequestParameters body actions.SearchRequestParameters false "SearchRequestParameters model containing metadata, conditions and query params"
// @Success		200
// @Router		/v1/admin/paymails/search [post]
// @Security	x-auth-xpub
func (a *Action) paymailAddressesSearch(c *gin.Context) {
	queryParams, metadata, conditions, err := actions.GetSearchQueryParameters(c)
	if err != nil {
		c.JSON(http.StatusExpectationFailed, err.Error())
		return
	}

	var paymailAddresses []*engine.PaymailAddress
	if paymailAddresses, err = a.Services.SpvWalletEngine.GetPaymailAddresses(
		c.Request.Context(),
		metadata,
		conditions,
		queryParams,
	); err != nil {
		c.JSON(http.StatusExpectationFailed, err.Error())
		return
	}

	c.JSON(http.StatusOK, paymailAddresses)
}

// paymailAddressesCount will count all paymail addresses filtered by metadata
// Paymail addresses count by metadata godoc
// @Summary		Paymail addresses count
// @Description	Paymail addresses count
// @Tags		Admin
// @Produce		json
// @Param		CountRequestParameters body actions.CountRequestParameters false "CountRequestParameters model containing metadata and conditions"
// @Success		200
// @Router		/v1/admin/paymails/count [post]
// @Security	x-auth-xpub
func (a *Action) paymailAddressesCount(c *gin.Context) {
	metadata, conditions, err := actions.GetCountQueryParameters(c)
	if err != nil {
		c.JSON(http.StatusExpectationFailed, err.Error())
		return
	}

	var count int64
	if count, err = a.Services.SpvWalletEngine.GetPaymailAddressesCount(
		c.Request.Context(),
		metadata,
		conditions,
	); err != nil {
		c.JSON(http.StatusExpectationFailed, err.Error())
		return
	}

	c.JSON(http.StatusOK, count)
}

// paymailCreateAddress will create a new paymail address
// Create Paymail godoc
// @Summary		Create paymail
// @Description	Create paymail
// @Tags		Admin
// @Param		CreatePaymail body CreatePaymail true "CreatePaymail model containing paymail information"
// @Produce		json
// @Success		201
// @Router		/v1/admin/paymail/create [post]
// @Security	x-auth-xpub
func (a *Action) paymailCreateAddress(c *gin.Context) {
	var requestBody CreatePaymail
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if requestBody.XpubID == "" {
		c.JSON(http.StatusExpectationFailed, "xpub is required")
		return
	}
	if requestBody.Address == "" {
		c.JSON(http.StatusExpectationFailed, "address is required")
		return
	}

	opts := a.Services.SpvWalletEngine.DefaultModelOptions()

	if requestBody.Metadata != nil {
		opts = append(opts, engine.WithMetadatas(requestBody.Metadata))
	}

	var paymailAddress *engine.PaymailAddress
	paymailAddress, err := a.Services.SpvWalletEngine.NewPaymailAddress(
		c.Request.Context(), requestBody.XpubID, requestBody.Address, requestBody.PublicName, requestBody.Avatar, opts...)
	if err != nil {
		c.JSON(http.StatusExpectationFailed, err.Error())
		return
	}

	c.JSON(http.StatusCreated, paymailAddress)
}

// paymailDeleteAddress will delete a paymail address
// Delete Paymail godoc
// @Summary		Delete paymail
// @Description	Delete paymail
// @Tags		Admin
// @Param		PaymailAddress body string true "PaymailAddress model containing paymail address to delete"
// @Produce		json
// @Success		200
// @Router		/v1/admin/paymail/delete [delete]
// @Security	x-auth-xpub
func (a *Action) paymailDeleteAddress(c *gin.Context) {
	var requestBody PaymailAddress
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if requestBody.Address == "" {
		c.JSON(http.StatusExpectationFailed, "address is required")
		return
	}

	opts := a.Services.SpvWalletEngine.DefaultModelOptions()

	// Delete a new paymail address
	err := a.Services.SpvWalletEngine.DeletePaymailAddress(c.Request.Context(), requestBody.Address, opts...)
	if err != nil {
		c.JSON(http.StatusExpectationFailed, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
