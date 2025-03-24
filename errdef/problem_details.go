package errdef

import (
	"fmt"
	"strings"

	"github.com/bitcoin-sv/spv-wallet/api"
	"github.com/joomcode/errorx"
)

// ProblemDetails is a struct that represents a problem details object as defined in RFC 7807.
// https://datatracker.ietf.org/doc/html/rfc7807
type ProblemDetails struct {
	api.ErrorsProblemDetails
}

// PushDetail appends a detail to the existing details, separated by a semicolon.
func (p *ProblemDetails) PushDetail(detail string) *ProblemDetails {
	separator := ""
	if p.Detail != "" {
		separator = "; "
	}
	p.Detail = fmt.Sprintf("%v%s%v", p.Detail, separator, detail)
	return p
}

// FromInternalError maps an internal error to a ProblemDetails object.
func (p *ProblemDetails) FromInternalError(err error) *ProblemDetails {
	ex := errorx.Cast(err)
	if ex == nil {
		return p
	}

	if hint, ok := ex.Property(PropPublicHint); ok {
		p.PushDetail(fmt.Sprintf("Hint: %v", hint))
	}

	p.Instance = strings.ReplaceAll(ex.Type().FullName(), ".", "/")

	return p
}
