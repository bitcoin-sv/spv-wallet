package contacts

import (
	"net/http"

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

// getContacts will fetch a list of contacts
// @Summary		Get contacts
// @Description	Get contacts
// @Tags		Contacts
// @Produce		json
// @Param		SwaggerCommonParams query swagger.CommonFilteringQueryParams false "Supports options for pagination and sorting to streamline data exploration and analysis"
// @Param		ContactParams query filter.ContactFilter false "Supports targeted resource searches with filters"
// @Success		200 {object} response.PageModel[response.Contact] "Page of contacts"
// @Failure		400	"Bad request - Error while parsing SearchContacts from request body"
// @Failure 	500	"Internal server error - Error while searching for contacts"
// @Router		/api/v1/contacts [get]
// @Security	x-auth-xpub
func getContacts(c *gin.Context, userContext *reqctx.UserContext) {
	logger := reqctx.Logger(c)
	engine := reqctx.Engine(c)
	reqXPubID := userContext.GetXPubID()

	searchParams, err := query.ParseSearchParams[filter.ContactFilter](c)
	if err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotParseQueryParams, logger)
		return
	}

	conditions, err := searchParams.Conditions.ToDbConditions()
	if err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrInvalidConditions, logger)
		return
	}
	metadata := mappings.MapToMetadata(searchParams.Metadata)
	pageOptions := mappings.MapToDbQueryParams(&searchParams.Page)

	contacts, err := engine.GetContactsByXpubID(
		c.Request.Context(),
		reqXPubID,
		metadata,
		conditions,
		pageOptions,
	)
	if err != nil {
		spverrors.ErrorResponse(c, err, logger)
		return
	}

	contracts := make([]*response.Contact, 0)
	for _, contact := range contacts {
		contracts = append(contracts, mappings.MapToContactContract(contact))
	}

	count, err := engine.GetContactsByXPubIDCount(
		c.Request.Context(),
		reqXPubID,
		metadata,
		conditions,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	response := response.PageModel[response.Contact]{
		Content: contracts,
		Page:    common.GetPageDescriptionFromSearchParams(pageOptions, count),
	}

	c.JSON(http.StatusOK, response)
}

// getContactByPaymail will fetch a list of contacts
// @Summary		Get contact by paymail
// @Description	Get contact by paymail
// @Tags		Contacts
// @Produce		json
// @Success		200 {object} response.Contact "Contact"
// @Failure		400	"Bad request - Error while parsing SearchContacts from request body"
// @Failure 	500	"Internal server error - Error while searching for contacts"
// @Router		/api/v1/contacts/{paymail} [get]
// @Security	x-auth-xpub
func getContactByPaymail(c *gin.Context, userContext *reqctx.UserContext) {
	paymail := c.Param("paymail")

	contacts, _ := searchContacts(c, userContext.GetXPubID(), paymail)
	if contacts == nil {
		return
	}

	if contacts == nil || len(contacts) != 1 {
		spverrors.ErrorResponse(c, spverrors.ErrContactNotFound, reqctx.Logger(c))
		return
	}

	contracts := mappings.MapToContactContracts(contacts)

	c.JSON(http.StatusOK, contracts[0])
}

// searchContacts - a helper function for searching contacts
func searchContacts(c *gin.Context, reqXPubID string, paymail string) ([]*engine.Contact, int64) {
	logger := reqctx.Logger(c)
	engine := reqctx.Engine(c)
	var reqParams filter.SearchContacts

	if paymail != "" {
		preConditions := filter.ContactFilter{Paymail: &paymail}
		reqParams.ConditionsModel = filter.ConditionsModel[filter.ContactFilter]{
			Conditions: &preConditions,
		}
	}

	if err := c.Bind(&reqParams); err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest, logger)
		return nil, 0
	}

	conditions, err := reqParams.Conditions.ToDbConditions()
	if err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrInvalidConditions, logger)
		return nil, 0
	}

	contacts, err := engine.GetContactsByXpubID(
		c.Request.Context(),
		reqXPubID,
		mappings.MapToMetadata(reqParams.Metadata),
		conditions,
		mappings.DefaultDBQueryParams(),
	)
	if err != nil {
		spverrors.ErrorResponse(c, err, logger)
		return nil, 0
	}

	count, err := engine.GetContactsByXPubIDCount(
		c.Request.Context(),
		reqXPubID,
		mappings.MapToMetadata(reqParams.Metadata),
		conditions,
	)
	if err != nil {
		spverrors.ErrorResponse(c, err, logger)
		return nil, 0
	}

	return contacts, count
}
