package errdef

import (
	"fmt"
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
	p.Detail = fmt.Sprintf("%v%s%v", detail, separator, p.Detail)
	return p
}

func (p *ProblemDetails) FromInternalError(err error) *ProblemDetails {
	ex := errorx.Cast(err)
	if ex == nil {
		return p
	}

	for _, trait := range globalTraits {
		if ex.HasTrait(trait.Trait) {
			p.PushDetail(trait.Title)
		}
	}

	if hint, ok := ex.Property(PropPublicHint); ok {
		p.PushDetail(fmt.Sprintf("Hint: %v", hint))
	}

	if instance, ok := ex.Property(PropSpecificProblemOccurrence); ok {
		p.Instance = fmt.Sprintf("%v", instance)
	}

	return p
}
