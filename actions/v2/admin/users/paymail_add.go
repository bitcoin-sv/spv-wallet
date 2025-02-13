package users

import (
	"net/http"

	"github.com/bitcoin-sv/go-paymail"
	adminerrors "github.com/bitcoin-sv/spv-wallet/actions/v2/admin/errors"
	"github.com/bitcoin-sv/spv-wallet/actions/v2/admin/internal/mapping"
	configerrors "github.com/bitcoin-sv/spv-wallet/config/errors"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/paymails/paymailsmodels"
	"github.com/bitcoin-sv/spv-wallet/models/request/adminrequest"
	"github.com/gin-gonic/gin"
)

// AddPaymailToUser add paymails to the user
func (s *APIAdminUsers) AddPaymailToUser(c *gin.Context, id string) {
	var requestBody adminrequest.AddPaymail
	if err := c.Bind(&requestBody); err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest.Wrap(err), s.logger)
		return
	}

	alias, domain, err := parsePaymail(&requestBody)
	if err != nil {
		spverrors.ErrorResponse(c, err, s.logger)
		return
	}

	newPaymail := &paymailsmodels.NewPaymail{
		Alias:      alias,
		Domain:     domain,
		PublicName: requestBody.PublicName,
		Avatar:     requestBody.Avatar,
		UserID:     id,
	}
	createdPaymail, err := s.engine.PaymailsService().Create(c, newPaymail)
	if err != nil {
		spverrors.MapResponse(c, err, s.logger).
			If(configerrors.ErrUnsupportedDomain).Then(adminerrors.ErrInvalidDomain).
			Else(adminerrors.ErrAddingPaymail)
		return
	}

	c.JSON(http.StatusCreated, mapping.PaymailToAdminResponse(createdPaymail))
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
