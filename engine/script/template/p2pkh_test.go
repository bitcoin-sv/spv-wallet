package template

import (
	"testing"

	assert "github.com/stretchr/testify/require"
)

func TestCreatePikeOutput(t *testing.T) {
	tests := []struct {
		name     string
		satoshis uint64
		wantErr  bool
		expected []P2PKHTemplate
	}{
		{
			name:     "valid input",
			satoshis: 1000,
			wantErr:  false,
			expected: []P2PKHTemplate{
				{
					Script:   "76a9fd88ac",
					Satoshis: 1000,
				},
			},
		},
		{
			name:     "zero satoshis",
			satoshis: 0,
			wantErr:  false,
			expected: []P2PKHTemplate{
				{
					Script:   "76a9fd88ac",
					Satoshis: 0,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := P2PKH(tt.satoshis)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, got)
			}
		})
	}
}
