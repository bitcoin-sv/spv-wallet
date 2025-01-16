package query

import (
	"time"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
)

// ParseSearchParams parses search params from the query string into a SearchParams struct with conditions of a given type.
func ParseSearchParams[T any](c *gin.Context) (*filter.SearchParams[T], error) {
	var params filter.SearchParams[T]

	dicts, err := ShouldGetQueryNestedMap(c)
	if err != nil {
		return nil, err
	}

	config := mapstructure.DecoderConfig{
		DecodeHook:       mapstructure.StringToTimeHookFunc(time.RFC3339),
		WeaklyTypedInput: true,
		Squash:           true,
		Result:           &params,
		TagName:          "json", // Small hax to reuse json tags which we have already defined
	}

	decoder, err := mapstructure.NewDecoder(&config)
	if err != nil {
		return nil, spverrors.Wrapf(err, spverrors.ErrCannotParseQueryParams.Error())
	}

	err = decoder.Decode(dicts)
	if err != nil {
		return nil, spverrors.Wrapf(err, spverrors.ErrCannotParseQueryParams.Error())
	}

	return &params, nil
}
