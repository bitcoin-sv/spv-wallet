package config

import (
	"testing"
)

func TestChainstateConfig_Validate(t *testing.T) {
	type fields struct {
		Monitoring *MonitoringConfig
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "valid - enabled mempool monitoring with filter",
			fields: fields{
				Monitoring: &MonitoringConfig{
					Mempool: &MempoolConfiguration{
						Enabled: true,
						Filter:  "0010abdefe",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid - enabled mempool monitoring with empty filter",
			fields: fields{
				Monitoring: &MonitoringConfig{
					Mempool: &MempoolConfiguration{
						Enabled: true,
						Filter:  "",
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ChainstateConfig{
				Monitoring: tt.fields.Monitoring,
			}
			err := c.Validate()
			if err != nil && !tt.wantErr {
				t.Fatalf("Validate() unexepectedly failed: %v", err)
			}
		})
	}
}
