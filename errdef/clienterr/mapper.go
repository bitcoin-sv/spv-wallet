package clienterr

import (
	"github.com/gin-gonic/gin"
	"github.com/joomcode/errorx"
	"github.com/rs/zerolog"
)

type Mapper interface {
	IfOfType(typeToMatch *errorx.Type) OnMatch
	Finalize() error
	Response(c *gin.Context, log *zerolog.Logger)
}

// OnMatch provides a fluent API for defining the response to return when a match is found.
type OnMatch interface {
	Then(errToReturn *Builder) Mapper
}

// Map creates a new Mapper instance.
func Map(err error) Mapper {
	return &MapperBuilder{
		baseErr: err,
	}
}

type MapperBuilder struct {
	baseErr  error
	matching bool
	final    *Builder
}

func (r *MapperBuilder) IfOfType(typeToMatch *errorx.Type) OnMatch {
	if r.final != nil {
		return r
	}
	r.matching = false
	if errorx.IsOfType(r.baseErr, typeToMatch) {
		r.matching = true
	}
	return r
}

func (r *MapperBuilder) Then(errToReturn *Builder) Mapper {
	if r.final == nil && r.matching {
		r.final = errToReturn
	}
	return r
}

func (r *MapperBuilder) Finalize() error {
	if r.final == nil {
		return r.baseErr
	}

	return r.final.Wrap(r.baseErr).Err()
}

func (r *MapperBuilder) Response(c *gin.Context, log *zerolog.Logger) {
	err := r.Finalize()
	Response(c, err, log)
}
