package testabilities

import (
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

type testServer struct {
	handlers *gin.Engine
}

func (t testServer) RoundTrip(request *http.Request) (*http.Response, error) {
	r := httptest.NewRecorder()
	t.handlers.ServeHTTP(r, request)
	return r.Result(), nil
}
