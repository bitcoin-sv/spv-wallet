package notifications

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_Notify(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	ctx := context.Background()
	webhookURL := "https://test.example.com/v1/api-endpoint"

	type args struct {
		modelType string
		eventType EventType
		model     interface{}
		id        string
	}

	useArgs := args{
		modelType: "transaction",
		eventType: EventTypeCreate,
		model:     map[string]interface{}{},
		id:        "test-id",
	}

	var tests = []struct {
		name      string
		options   []ClientOps
		args      args
		wantErr   assert.ErrorAssertionFunc
		httpCalls int
		httpMock  func()
	}{
		{
			name:      "empty notification",
			options:   nil,
			args:      useArgs,
			wantErr:   assert.NoError,
			httpCalls: 0,
			httpMock:  func() {},
		},
		{
			name: "http call done",
			options: []ClientOps{
				WithNotifications(webhookURL),
			},
			args:      useArgs,
			wantErr:   assert.NoError,
			httpCalls: 1,
			httpMock: func() {
				httpmock.RegisterResponder(http.MethodPost, webhookURL,
					httpmock.NewStringResponder(
						http.StatusOK,
						`OK`,
					),
				)
			},
		},
		{
			name: "http error",
			options: []ClientOps{
				WithNotifications(webhookURL),
			},
			args:      useArgs,
			wantErr:   assert.Error,
			httpCalls: 1,
			httpMock: func() {
				httpmock.RegisterResponder(http.MethodPost, webhookURL,
					httpmock.NewErrorResponder(errors.New("error")),
				)
			},
		},
	}
	for _, tt := range tests {
		httpmock.Reset()
		tt.httpMock()
		t.Run(tt.name, func(t *testing.T) {
			c, err := NewClient(tt.options...)
			require.NoError(t, err)
			tt.wantErr(t, c.Notify(ctx, tt.args.modelType, tt.args.eventType, tt.args.model, tt.args.id), fmt.Sprintf("Notify(%v, %v, %v, %v, %v)", ctx, tt.args.modelType, tt.args.eventType, tt.args.model, tt.args.id))
			assert.Equal(t, tt.httpCalls, httpmock.GetTotalCallCount())
		})
	}
}
