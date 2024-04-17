package notifications

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithNotifications(t *testing.T) {
	type args struct {
		webhookEndpoint string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "empty",
			args: args{webhookEndpoint: ""},
			want: "",
		},
		{
			name: "empty",
			args: args{webhookEndpoint: "https://example.com/v1"},
			want: "https://example.com/v1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := []ClientOps{WithNotifications(tt.args.webhookEndpoint)}
			client, err := NewClient(opts...)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, client.GetWebhookEndpoint())
		})
	}
}

func Test_defaultClientOptions(t *testing.T) {
	tests := []struct {
		name string
		want *clientOptions
	}{
		{
			name: "options",
			want: &clientOptions{
				config: &notificationsConfig{
					webhookEndpoint: "",
				},
				httpClient: &http.Client{
					Timeout: defaultHTTPTimeout,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, defaultClientOptions(), "defaultClientOptions()")
		})
	}
}
