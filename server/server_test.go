package server

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/config"
	"github.com/bitcoin-sv/spv-wallet/engine/utils"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	testXpubAuth = "xpub661MyMwAqRbcGpZVrSHU7EZ5Zwx5cNZmD5iLHPcg8MPnVcPdsApRi4Z27Mg3Zy53XYMKuJC5GiwECCFVNkhNgrBrfcA22YoJhasH7GcArNX"
)

// TestSuite is for testing the entire package using real/mocked services
type TestSuite struct {
	tests.TestSuite
}

// SetupSuite runs at the start of the suite
func (ts *TestSuite) SetupSuite() {
	ts.BaseSetupSuite()
}

// TearDownSuite runs after the suite finishes
func (ts *TestSuite) TearDownSuite() {
	ts.BaseTearDownSuite()
}

// SetupTest runs before each test
func (ts *TestSuite) SetupTest() {
	ts.BaseSetupTest()

	setupServerRoutes(ts.AppConfig, ts.Services, ts.Router)
}

// TearDownTest runs after each test
func (ts *TestSuite) TearDownTest() {
	ts.BaseTearDownTest()
}

// TestTestSuite kick-starts all suite tests
func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

// TestAdminAuthentication will test admin authentication
func (ts *TestSuite) TestAdminAuthentication() {
	ts.T().Run("no value", func(t *testing.T) {
		w := httptest.NewRecorder()

		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/"+config.APIVersion+"/admin/status", nil)
		require.NoError(t, err)
		require.NotNil(t, req)

		ts.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	ts.T().Run("false value", func(t *testing.T) {
		w := httptest.NewRecorder()

		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/"+config.APIVersion+"/admin/status", nil)
		require.NoError(t, err)
		require.NotNil(t, req)

		req.Header.Set(models.AuthHeader, testXpubAuth)

		ts.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	ts.T().Run("admin key", func(t *testing.T) {
		w := httptest.NewRecorder()

		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/"+config.APIVersion+"/admin/status", bytes.NewReader([]byte("test")))
		require.NoError(t, err)
		require.NotNil(t, req)

		req.Header.Set(models.AuthHeader, ts.AppConfig.Authentication.AdminKey)

		ts.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

// TestApiAuthentication will api authentication
func (ts *TestSuite) TestApiAuthentication() {
	ts.T().Run("no value", func(t *testing.T) {
		w := httptest.NewRecorder()

		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/"+config.APIVersion+"/xpub", nil)
		require.NoError(t, err)
		require.NotNil(t, req)

		ts.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	ts.T().Run("false value", func(t *testing.T) {
		w := httptest.NewRecorder()

		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/"+config.APIVersion+"/xpub", nil)
		require.NoError(t, err)
		require.NotNil(t, req)

		req.Header.Set(models.AuthHeader, testXpubAuth)

		ts.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	ts.T().Run("valid value", func(t *testing.T) {
		w := httptest.NewRecorder()

		xpub, err := ts.Services.SpvWalletEngine.NewXpub(context.Background(), testXpubAuth)
		require.NoError(t, err)
		require.NotNil(t, xpub)

		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/"+config.APIVersion+"/xpub", bytes.NewReader([]byte("test")))
		require.NoError(t, err)
		require.NotNil(t, req)

		req.Header.Set(models.AuthHeader, xpub.RawXpub())

		ts.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

// TestBasicAuthentication will api authentication
func (ts *TestSuite) TestBasicAuthentication() {
	ts.T().Run("no value", func(t *testing.T) {
		w := httptest.NewRecorder()

		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/"+config.APIVersion+"/transaction", nil)
		require.NoError(t, err)
		require.NotNil(t, req)

		ts.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	ts.T().Run("non existing xpub", func(t *testing.T) {
		w := httptest.NewRecorder()

		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/"+config.APIVersion+"/transaction", nil)
		require.NoError(t, err)
		require.NotNil(t, req)

		req.Header.Set(models.AuthHeader, testXpubAuth)

		ts.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	ts.T().Run("valid value", func(t *testing.T) {
		w := httptest.NewRecorder()

		xpub, err := ts.Services.SpvWalletEngine.NewXpub(context.Background(), testXpubAuth)
		require.NoError(t, err)
		require.NotNil(t, xpub)

		destination, err := ts.Services.SpvWalletEngine.NewDestination(context.Background(), xpub.RawXpub(), 0, utils.ScriptTypePubKeyHash)
		require.NoError(t, err)
		require.NotNil(t, destination)

		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/"+config.APIVersion+"/destination?id="+destination.GetID(), bytes.NewReader([]byte("test")))
		require.NoError(t, err)
		require.NotNil(t, req)

		req.Header.Set(models.AuthHeader, xpub.RawXpub())

		ts.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
