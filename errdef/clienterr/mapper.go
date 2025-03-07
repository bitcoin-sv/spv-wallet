package clienterr

import (
	"github.com/gin-gonic/gin"
	"github.com/joomcode/errorx"
	"github.com/rs/zerolog"
)

// Mapper is a fluent API for mapping errors to client errors.
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
	return &mapperBuilder{
		baseErr: err,
	}
}

type mapperBuilder struct {
	baseErr  error
	matching bool
	final    *Builder
}

func (r *mapperBuilder) IfOfType(typeToMatch *errorx.Type) OnMatch {
	if r.final != nil {
		return r
	}
	r.matching = false
	if errorx.IsOfType(r.baseErr, typeToMatch) {
		r.matching = true
	}
	return r
}

func (r *mapperBuilder) Then(errToReturn *Builder) Mapper {
	if r.final == nil && r.matching {
		r.final = errToReturn
	}
	return r
}

func (r *mapperBuilder) Finalize() error {
	if r.final == nil {
		return r.baseErr
	}

	return r.final.Wrap(r.baseErr).Err()
}

func (r *mapperBuilder) Response(c *gin.Context, log *zerolog.Logger) {
	err := r.Finalize()
	Response(c, err, log)
}
