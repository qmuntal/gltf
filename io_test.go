package gltf

import "testing"

func TestRelativeFileHandler_ReadFull(t *testing.T) {
	type args struct {
		uri  string
		data []byte
	}
	tests := []struct {
		name    string
		h       *RelativeFileHandler
		args    args
		wantErr bool
	}{
		{"no dir", new(RelativeFileHandler), args{"a.bin", []byte{}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.h.ReadFull(tt.args.uri, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("RelativeFileHandler.ReadFull() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestProtocolRegistry_ReadFull(t *testing.T) {
	type args struct {
		uri  string
		data []byte
	}
	tests := []struct {
		name    string
		reg     ProtocolRegistry
		args    args
		wantErr bool
	}{
		{"invalid url", make(ProtocolRegistry), args{"%$·$·23", []byte{}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.reg.ReadFull(tt.args.uri, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("ProtocolRegistry.ReadFull() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
