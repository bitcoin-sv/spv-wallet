package contacts

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/actions/v2/internal/mapping"
	"github.com/bitcoin-sv/spv-wallet/api"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/server/reqctx"
	"github.com/gin-gonic/gin"
)

func (s *APIContacts) GetContacts(c *gin.Context, params api.GetContactsParams) {
	userContext := reqctx.GetUserContext(c)
	userID, err := userContext.ShouldGetUserID()
	if err != nil {
		spverrors.ErrorResponse(c, err, s.logger)
		return
	}

	page := mapping.MapContactsParamToFilterPage(params)
	conditions := mapping.MapToDBConditions(params)

	pagedResult, err := s.engine.ContactService().PaginatedForUser(c.Request.Context(), userID, page, conditions)
	if err != nil {
		spverrors.ErrorResponse(c, err, s.logger)
		return
	}

	c.JSON(http.StatusOK, mapping.ContactsPagedResponse(pagedResult))
}
