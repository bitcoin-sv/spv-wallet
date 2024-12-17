package filter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdminContactFilter(t *testing.T) {
	tests := []struct {
		name    string
		filter  AdminContactFilter
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name:   "Empty filter",
			filter: AdminContactFilter{},
			want: map[string]interface{}{
				"deleted_at": nil,
			},
			wantErr: false,
		},
		{
			name: "With XPubID",
			filter: AdminContactFilter{
				XPubID: ptrString("623bc25ce1c0fc510dea72b5ee27b2e70384c099f1f3dce9e73dd987198c3486"),
			},
			want: map[string]interface{}{
				"xpub_id":    "623bc25ce1c0fc510dea72b5ee27b2e70384c099f1f3dce9e73dd987198c3486",
				"deleted_at": nil,
			},
			wantErr: false,
		},
		{
			name: "With ContactFilter conditions",
			filter: AdminContactFilter{ContactFilter: ContactFilter{
				Paymail: ptrString("test@example.com"),
				ModelFilter: ModelFilter{
					IncludeDeleted: ptrBool(true),
				},
			},
				XPubID: ptrString("623bc25ce1c0fc510dea72b5ee27b2e70384c099f1f3dce9e73dd987198c3486"),
			},
			want: map[string]interface{}{
				"paymail": "test@example.com",
				"xpub_id": "623bc25ce1c0fc510dea72b5ee27b2e70384c099f1f3dce9e73dd987198c3486",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.filter.ToDbConditions()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func ptrString(s string) *string {
	return &s
}

func ptrBool(b bool) *bool {
	return &b
}
