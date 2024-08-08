package accesskeys

import (
	"fmt"
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/actions/common"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/bitcoin-sv/spv-wallet/models/response"
	"github.com/bitcoin-sv/spv-wallet/server/auth"
	"github.com/gin-gonic/gin"
)

// oldSearch will fetch a list of access keys filtered by metadata
// Search access key godoc
// @Summary		Search access key - Use (GET) /api/v1/users/current/keys instead.
// @Description	This endpoint has been deprecated. Use (GET) /api/v1/users/current/keys instead.
// @Tags		Access-key
// @Produce		json
// @Param		SearchAccessKeys body filter.SearchAccessKeys false "Supports targeted resource searches with filters and metadata, plus options for pagination and sorting to streamline data exploration and analysis"
// @Success		200 {object} []models.AccessKey "List of access keys"
// @Failure		400	"Bad request - Error while SearchAccessKeys from request body"
// @Failure 	500	"Internal server error - Error while searching for access keys"
// @DeprecatedRouter  /v1/access-key/search [post]
// @Security	x-auth-xpub
func (a *Action) oldSearch(c *gin.Context) {
	reqXPubID := c.GetString(auth.ParamXPubHashKey)

	var reqParams filter.SearchAccessKeys
	if err := c.Bind(&reqParams); err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest, a.Services.Logger)
		return
	}

	accessKeys, err := a.Services.SpvWalletEngine.GetAccessKeysByXPubID(
		c.Request.Context(),
		reqXPubID,
		mappings.MapToMetadata(reqParams.Metadata),
		reqParams.Conditions.ToDbConditions(),
		mappings.MapToQueryParams(reqParams.QueryParams),
	)
	if err != nil {
		spverrors.ErrorResponse(c, err, a.Services.Logger)
		return
	}

	accessKeyContracts := make([]*models.AccessKey, 0)
	for _, accessKey := range accessKeys {
		accessKeyContracts = append(accessKeyContracts, mappings.MapToOldAccessKeyContract(accessKey))
	}

	c.JSON(http.StatusOK, accessKeyContracts)
}

// search will fetch a list of access keys filtered by metadata
// Search access key godoc
// @Summary		Search access key
// @Description	Search access key
// @Tags		Access-key
// @Produce		json
// @Param		SearchAccessKeys body filter.SearchAccessKeys false "Supports targeted resource searches with filters and metadata, plus options for pagination and sorting to streamline data exploration and analysis"
// @Success		200 {object} []response.AccessKey "List of access keys"
// @Failure		400	"Bad request - Error while SearchAccessKeys from request body"
// @Failure 	500	"Internal server error - Error while searching for access keys"
// @Router		/api/v1//users/current/keys [get]
// @Security	x-auth-xpub
func (a *Action) search(c *gin.Context) {
	reqXPubID := c.GetString(auth.ParamXPubHashKey)

	var reqParams filter.SearchAccessKeys
	if err := c.Bind(&reqParams); err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest, a.Services.Logger)
		return
	}

	conditions := reqParams.Conditions.ToDbConditions()
	reqParams.DefaultsIfNil()

	accessKeys, err := a.Services.SpvWalletEngine.GetAccessKeysByXPubID(
		c.Request.Context(),
		reqXPubID,
		mappings.MapToMetadata(reqParams.Metadata),
		conditions,
		mappings.MapToQueryParams(reqParams.QueryParams),
	)
	if err != nil {
		spverrors.ErrorResponse(c, err, a.Services.Logger)
		return
	}

	accessKeyContracts := make([]*response.AccessKey, 0)
	for _, accessKey := range accessKeys {
		accessKeyContracts = append(accessKeyContracts, mappings.MapToAccessKeyContract(accessKey))
	}

	count, err := a.Services.SpvWalletEngine.GetContactsByXPubIDCount(
		c.Request.Context(),
		reqXPubID,
		mappings.MapToMetadata(reqParams.Metadata),
		conditions,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	response := response.PageModel[response.AccessKey]{
		Content: accessKeyContracts,
		Page: response.PageDescription{
			Size:          len(accessKeyContracts),
			Number:        0,
			TotalElements: int(count),
			TotalPages:    len(accessKeyContracts) / int(count),
		},
	}

	c.JSON(http.StatusOK, response)
}

func (a *Action) searchTest(c *gin.Context) {
	pageable := common.ExtractPageableFromRequest(c)

	fmt.Printf("Pageable from request: %+v\n", pageable)

	reqXPubID := c.GetString(auth.ParamXPubHashKey)

	var reqParams filter.SearchAccessKeys
	if err := c.Bind(&reqParams); err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest, a.Services.Logger)
		return
	}

	conditions := reqParams.Conditions.ToDbConditions()
	reqParams.DefaultsIfNil()

	accessKeys, err := a.Services.SpvWalletEngine.GetAccessKeysByXPubID(
		c.Request.Context(),
		reqXPubID,
		mappings.MapToMetadata(reqParams.Metadata),
		conditions,
		mappings.MapToQueryParams(reqParams.QueryParams),
	)
	if err != nil {
		spverrors.ErrorResponse(c, err, a.Services.Logger)
		return
	}

	accessKeyContracts := make([]*response.AccessKey, 0)
	for _, accessKey := range accessKeys {
		accessKeyContracts = append(accessKeyContracts, mappings.MapToAccessKeyContract(accessKey))
	}

	count, err := a.Services.SpvWalletEngine.GetContactsByXPubIDCount(
		c.Request.Context(),
		reqXPubID,
		mappings.MapToMetadata(reqParams.Metadata),
		conditions,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	response := response.PageModel[response.AccessKey]{
		Content: accessKeyContracts,
		Page:    common.GetPageDescriptionFromQueryParams(pageable, count),
	}

	c.JSON(http.StatusOK, response)

	// c.JSON(http.StatusOK, gin.H{
	// 	"pageable": "OK",
	// })
}
