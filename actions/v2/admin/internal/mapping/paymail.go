package mapping

import (
	"github.com/bitcoin-sv/go-paymail"
	"github.com/bitcoin-sv/spv-wallet/api"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/paymails/paymailsmodels"
	"github.com/bitcoin-sv/spv-wallet/errdef/clienterr"
	"github.com/bitcoin-sv/spv-wallet/lox"
)

// PaymailToAdminResponse maps a paymail to a response
func PaymailToAdminResponse(p *paymailsmodels.Paymail) api.ModelsPaymail {
	return api.ModelsPaymail{
		Id:         p.ID,
		Alias:      p.Alias,
		Domain:     p.Domain,
		Paymail:    p.Alias + "@" + p.Domain,
		PublicName: p.PublicName,
		Avatar:     p.Avatar,
	}
}

// UsersPaymailToResponse maps a user's paymail to a response
func UsersPaymailToResponse(p *paymailsmodels.Paymail) api.ModelsPaymail {
	return api.ModelsPaymail{
		Id:         p.ID,
		Alias:      p.Alias,
		Domain:     p.Domain,
		Paymail:    p.Alias + "@" + p.Domain,
		PublicName: p.PublicName,
		Avatar:     p.Avatar,
	}
}

// RequestAddPaymailToNewPaymailModel maps a add paymail request to new paymail model
func RequestAddPaymailToNewPaymailModel(r *api.RequestsAddPaymail, userID string) (*paymailsmodels.NewPaymail, error) {
	alias, domain, err := parsePaymail(r)
	if err != nil {
		return nil, err
	}

	newPaymail := &paymailsmodels.NewPaymail{
		Alias:      alias,
		Domain:     domain,
		PublicName: lox.Unwrap(r.PublicName).Else(""),
		Avatar:     lox.Unwrap(r.AvatarURL).Else(""),
	}

	if userID != "" {
		newPaymail.UserID = userID
	}

	return newPaymail, nil
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

// parsePaymail parses the paymail address from the request body.
// Uses either Alias + Domain or the whole paymail Address field
// If both Alias + Domain and Address are set, and they are inconsistent, an error is returned.
func parsePaymail(r *api.RequestsAddPaymail) (string, string, error) {
	request := &addPaymailRequest{r}
	if request.HasAddress() &&
		(request.HasAlias() || request.HasDomain()) &&
		!request.AddressEqualsTo(request.Alias+"@"+request.Domain) {
		return "", "", clienterr.BadRequest.
			Detailed(
				"inconsistent_alias_domain_and_address",
				"Inconsistent alias@domain and address fields: %s, %s, %s. Hint: use either alias and domain or address (not both)",
				request.Alias, request.Domain, request.Address,
			).Err()
	}
	if !request.HasAddress() {
		request.Address = request.Alias + "@" + request.Domain
	}
	alias, domain, sanitized := paymail.SanitizePaymail(request.Address)
	if sanitized == "" {
		return "", "", clienterr.BadRequest.
			Detailed("invalid_paymail_address", "Invalid paymail address: %s", request.Address).Err()
	}
	return alias, domain, nil
}
