package notifications

/*func TestNewClient(t *testing.T) {
	type args struct {
		opts []ClientOps
	}
	tests := []struct {
		name    string
		args    args
		want    ClientInterface
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "empty",
			args:    args{opts: []ClientOps{}},
			want:    &Client{options: defaultClientOptions()},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewClient(tt.args.opts...)
			if !tt.wantErr(t, err, fmt.Sprintf("NewClient(%v)", tt.args.opts)) {
				return
			}
			assert.Equalf(t, tt.want, got, "NewClient(%v)", tt.args.opts)
		})
	}
}*/
