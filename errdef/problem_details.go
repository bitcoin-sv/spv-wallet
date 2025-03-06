package errdef

import (
	"fmt"
	"strings"

	"github.com/joomcode/errorx"
)

type ProblemDetails struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Status   int    `json:"status"`
	Detail   string `json:"detail"`
	Instance string `json:"instance"`
}

func (p *ProblemDetails) PushDetail(detail string) *ProblemDetails {
	separator := ""
	if p.Detail != "" {
		separator = "; "
	}
	p.Detail = fmt.Sprintf("%v%s%v", p.Detail, separator, detail)
	return p
}

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
