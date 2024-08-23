package query

import (
	"time"

	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
)

func ParseSearchParams[T any](c *gin.Context, _ T) (*filter.SearchParams[T], error) {
	var params filter.SearchParams[T]

	dicts, err := ShouldGetQueryNestedMap(c)
	if err != nil {
		return nil, err
	}

	config := mapstructure.DecoderConfig{
		DecodeHook:       mapstructure.StringToTimeHookFunc(time.RFC3339),
		WeaklyTypedInput: true,
		Result:           &params,
		TagName:          "json", // Small hax to reuse json tags which we have already defined
	}

	decoder, err := mapstructure.NewDecoder(&config)
	if err != nil {
		return nil, err
	}

	err = decoder.Decode(dicts)
	if err != nil {
		return nil, err
	}

	return &params, nil
}
