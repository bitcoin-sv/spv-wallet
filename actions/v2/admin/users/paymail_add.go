package users

import (
	"errors"
	"net/http"

	"github.com/bitcoin-sv/go-paymail"
	adminerrors "github.com/bitcoin-sv/spv-wallet/actions/v2/admin/errors"
	"github.com/bitcoin-sv/spv-wallet/actions/v2/admin/internal/mapping"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/paymails/paymailsmodels"
	"github.com/bitcoin-sv/spv-wallet/models/request/adminrequest"
	"github.com/bitcoin-sv/spv-wallet/server/reqctx"
	"github.com/gin-gonic/gin"
)

func addPaymail(c *gin.Context, _ *reqctx.AdminContext) {
	logger := reqctx.Logger(c)

	var requestBody adminrequest.AddPaymail
	if err := c.Bind(&requestBody); err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest.Wrap(err), logger)
		return
	}

	userID := c.Param("id")

	alias, domain, err := parsePaymail(&requestBody)
	if err != nil {
		spverrors.ErrorResponse(c, err, logger)
		return
	}

	newPaymail := &paymailsmodels.NewPaymail{
		Alias:      alias,
		Domain:     domain,
		PublicName: requestBody.PublicName,
		Avatar:     requestBody.Avatar,
		UserID:     userID,
	}
	createdPaymail, err := reqctx.Engine(c).PaymailsService().Create(c, newPaymail)
	if err != nil {
		if !errors.Is(err, spverrors.ErrInvalidDomain) {
			err = adminerrors.ErrAddingPaymail.Wrap(err)
		}
		spverrors.ErrorResponse(c, err, logger)
		return
	}

	c.JSON(http.StatusCreated, mapping.CreatedPaymailResponse(createdPaymail))
}

// parsePaymail parses the paymail address from the request body.
// Uses either Alias + Domain or the whole paymail Address field
// If both Alias + Domain and Address are set, and they are inconsistent, an error is returned.
func parsePaymail(request *adminrequest.AddPaymail) (string, string, error) {
	if request.HasAddress() &&
		(request.HasAlias() || request.HasDomain()) &&
		!request.AddressEqualsTo(request.Alias+"@"+request.Domain) {
		return "", "", adminerrors.ErrPaymailInconsistent
	}
	if !request.HasAddress() {
		request.Address = request.Alias + "@" + request.Domain
	}
	alias, domain, sanitized := paymail.SanitizePaymail(request.Address)
	if sanitized == "" {
		return "", "", adminerrors.ErrInvalidPaymail
	}
	return alias, domain, nil
}
