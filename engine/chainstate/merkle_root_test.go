package chainstate

import (
	"bytes"
	"context"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func initMockClient(ops ...ClientOps) (*Client, *buffLogger) {
	bLogger := newBuffLogger()
	ops = append(ops, WithLogger(bLogger.logger))
	c, _ := NewClient(
		context.Background(),
		ops...,
	)
	return c.(*Client), bLogger
}

func TestVerifyMerkleRoots(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mockURL := "http://block-headers-service.test/api/v1/chain/merkleroot/verify"

	t.Run("no block headers service client", func(t *testing.T) {
		c, _ := initMockClient()

		err := c.VerifyMerkleRoots(context.Background(), []MerkleRootConfirmationRequestItem{})

		assert.Error(t, err)
	})

	t.Run("block headers service is not online", func(t *testing.T) {
		httpmock.Reset()
		httpmock.RegisterResponder("POST", mockURL,
			httpmock.NewStringResponder(500, `{"error":"Internal Server Error"}`),
		)
		c, bLogger := initMockClient(WithConnectionToBlockHeaderService(mockURL, ""))

		err := c.VerifyMerkleRoots(context.Background(), []MerkleRootConfirmationRequestItem{})

		assert.Error(t, err)
		assert.Equal(t, 1, httpmock.GetTotalCallCount())
		assert.True(t, bLogger.contains("Block Headers Service client returned status code 500"))
	})

	t.Run("block headers service wrong auth", func(t *testing.T) {
		httpmock.Reset()
		httpmock.RegisterResponder("POST", mockURL,
			httpmock.NewStringResponder(401, `Unauthorized`),
		)
		c, bLogger := initMockClient(WithConnectionToBlockHeaderService(mockURL, "some-token"))

		err := c.VerifyMerkleRoots(context.Background(), []MerkleRootConfirmationRequestItem{})

		assert.Error(t, err)
		assert.Equal(t, 1, httpmock.GetTotalCallCount())
		assert.True(t, bLogger.contains("401"))
	})

	t.Run("block headers service invalid state", func(t *testing.T) {
		httpmock.Reset()
		httpmock.RegisterResponder("POST", mockURL,
			httpmock.NewJsonResponderOrPanic(200, MerkleRootsConfirmationsResponse{
				ConfirmationState: Invalid,
				Confirmations:     []MerkleRootConfirmation{},
			}),
		)
		c, bLogger := initMockClient(WithConnectionToBlockHeaderService(mockURL, "some-token"))

		err := c.VerifyMerkleRoots(context.Background(), []MerkleRootConfirmationRequestItem{})

		assert.Error(t, err)
		assert.Equal(t, 1, httpmock.GetTotalCallCount())
		assert.True(t, bLogger.contains("Not all merkle roots confirmed"))
	})

	t.Run("block headers service confirmedState", func(t *testing.T) {
		httpmock.Reset()
		httpmock.RegisterResponder("POST", mockURL,
			httpmock.NewJsonResponderOrPanic(200, MerkleRootsConfirmationsResponse{
				ConfirmationState: Confirmed,
				Confirmations: []MerkleRootConfirmation{
					{
						Hash:         "some-hash",
						BlockHeight:  1,
						MerkleRoot:   "some-merkle-root",
						Confirmation: Confirmed,
					},
				},
			}),
		)
		c, bLogger := initMockClient(WithConnectionToBlockHeaderService(mockURL, "some-token"))

		err := c.VerifyMerkleRoots(context.Background(), []MerkleRootConfirmationRequestItem{
			{
				MerkleRoot:  "some-merkle-root",
				BlockHeight: 1,
			},
		})

		assert.NoError(t, err)
		assert.Equal(t, 1, httpmock.GetTotalCallCount())
		assert.False(t, bLogger.contains("ERR"))
		assert.False(t, bLogger.contains("WARN"))
	})
}

// buffLogger allows to check if a certain string was logged
type buffLogger struct {
	logger *zerolog.Logger
	buf    *bytes.Buffer
}

func newBuffLogger() *buffLogger {
	var buf bytes.Buffer
	logger := zerolog.New(&buf).Level(zerolog.DebugLevel).With().Logger()
	return &buffLogger{logger: &logger, buf: &buf}
}

func (l *buffLogger) contains(expected string) bool {
	return bytes.Contains(l.buf.Bytes(), []byte(expected))
}
