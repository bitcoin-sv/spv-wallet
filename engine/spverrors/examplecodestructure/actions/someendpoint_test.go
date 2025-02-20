//go:build errorx
// +build errorx

package actions

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestDBFailure(t *testing.T) {
	endpoint := NewSomeEndpoint()

	URL, err := url.Parse("")
	require.NoError(t, err)

	c := &gin.Context{
		Request: &http.Request{
			URL: URL,
		},
	}

	res := endpoint.SomeSearch(c)

	fmt.Printf("%+v", res)
}

func TestInvalidQuery(t *testing.T) {
	endpoint := NewSomeEndpoint()

	URL, err := url.Parse("?query=invalid")
	require.NoError(t, err)

	c := &gin.Context{
		Request: &http.Request{
			URL: URL,
		},
	}

	res := endpoint.SomeSearch(c)

	fmt.Printf("%+v", res)
}

func TestDBSuccess(t *testing.T) {
	endpoint := NewSomeEndpoint()

	URL, err := url.Parse("?query=success")
	require.NoError(t, err)

	c := &gin.Context{
		Request: &http.Request{
			URL: URL,
		},
	}

	res := endpoint.SomeSearch(c)

	fmt.Printf("%+v", res)
}
