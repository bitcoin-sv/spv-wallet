package filters

import (
	"testing"

	"github.com/bitcoin-sv/spv-wallet/engine/chainstate"
)

func TestMetanet(t *testing.T) {
	type args struct {
		tx *chainstate.TxInfo
	}
	tests := []struct {
		name       string
		args       args
		passFilter bool
		wantErr    bool
	}{
		{
			name: "non-metanet transaction shouldn't pass filter",
			args: args{
				tx: &chainstate.TxInfo{
					Hex: "01000000010000000000000000000000000000000000000000000000000000000000000000ffffffff1c03d7c6082f7376706f6f6c2e636f6d2f3edff034600055b8467f0040ffffffff01247e814a000000001976a914492558fb8ca71a3591316d095afc0f20ef7d42f788ac00000000",
				},
			},
			passFilter: false,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Metanet(tt.args.tx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Metanet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.passFilter && got == nil {
				t.Errorf("Metanet() expected transaction to pass filter and it didn't")
			}
		})
	}
}
