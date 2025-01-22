package users

import (
	"net/http"

	"github.com/bitcoin-sv/go-paymail"
	adminerrors "github.com/bitcoin-sv/spv-wallet/actions/v2/admin/errors"
	"github.com/bitcoin-sv/spv-wallet/actions/v2/admin/internal/mapping"
	"github.com/bitcoin-sv/spv-wallet/engine/domainmodels"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
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

	if err = reqctx.AppConfig(c).Paymail.CheckDomain(domain); err != nil {
		spverrors.ErrorResponse(c, err, logger)
		return
	}

	newPaymail := &domainmodels.NewPaymail{
		Alias:      alias,
		Domain:     domain,
		PublicName: requestBody.PublicName,
		Avatar:     requestBody.Avatar,
		UserID:     userID,
	}
	createdPaymail, err := reqctx.Engine(c).UserService().AppendPaymail(c, newPaymail)
	if err != nil {
		spverrors.ErrorResponse(c, adminerrors.ErrAddingPaymail.Wrap(err), logger)
		return
	}

	c.JSON(http.StatusCreated, mapping.CreatedPaymailResponse(createdPaymail))
}

func parsePaymail(request *adminrequest.AddPaymail) (string, string, error) {
	if request.Address != "" &&
		(request.Alias != "" || request.Domain != "") &&
		request.Address != request.Alias+"@"+request.Domain {
		return "", "", adminerrors.ErrPaymailInconsistent
	}
	pm := request.Address
	if pm == "" {
		pm = request.Alias + "@" + request.Domain
	}
	alias, domain, sanitized := paymail.SanitizePaymail(pm)
	if sanitized == "" {
		return "", "", adminerrors.ErrInvalidPaymail
	}
	return alias, domain, nil
}
