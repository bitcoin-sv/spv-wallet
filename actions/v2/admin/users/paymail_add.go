package users

import (
	"net/http"

	"github.com/bitcoin-sv/go-paymail"
	adminerrors "github.com/bitcoin-sv/spv-wallet/actions/v2/admin/errors"
	"github.com/bitcoin-sv/spv-wallet/actions/v2/admin/internal/mapping"
	"github.com/bitcoin-sv/spv-wallet/api"
	configerrors "github.com/bitcoin-sv/spv-wallet/config/errors"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/paymails/paymailsmodels"
	"github.com/bitcoin-sv/spv-wallet/lox"
	"github.com/gin-gonic/gin"
)

// AddPaymailToUser add paymails to the user
func (s *APIAdminUsers) AddPaymailToUser(c *gin.Context, id string) {
	var request api.RequestsAddPaymail
	if err := c.Bind(&request); err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest.Wrap(err), s.logger)
		return
	}

	requestBody := addPaymailRequest{&request}

	alias, domain, err := parsePaymail(&requestBody)
	if err != nil {
		spverrors.ErrorResponse(c, err, s.logger)
		return
	}

	newPaymail := &paymailsmodels.NewPaymail{
		Alias:      alias,
		Domain:     domain,
		PublicName: lox.Unwrap(requestBody.PublicName).Else(""),
		Avatar:     lox.Unwrap(requestBody.AvatarURL).Else(""),
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
func parsePaymail(request *addPaymailRequest) (string, string, error) {
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

type addPaymailRequest struct {
	*api.RequestsAddPaymail
}

func (a addPaymailRequest) HasAddress() bool {
	return a.Address != ""
}

// HasAlias returns true if the paymail alias is set
func (a addPaymailRequest) HasAlias() bool {
	return a.Alias != ""
}

// HasDomain returns true if the paymail domain is set
func (a addPaymailRequest) HasDomain() bool {
	return a.Domain != ""
}

// AddressEqualsTo returns true if the paymail address is equal to the given string
func (a addPaymailRequest) AddressEqualsTo(s string) bool {
	return a.Address == s
}
