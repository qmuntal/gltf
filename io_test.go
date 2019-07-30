package gltf

import "testing"

func TestRelativeFileHandler_ReadFullResource(t *testing.T) {
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
			if err := tt.h.ReadFullResource(tt.args.uri, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("ReadFullResource.ReadFull() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
