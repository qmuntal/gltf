package gltf

import (
	"encoding/json"
	"testing"
)

func TestAccessorType_UnmarshalJSON(t *testing.T) {
	type args struct {
		defaultType AccessorType
		expType     AccessorType
		typeStr     []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"base", args{50, AccessorVec3, []byte(`"VEC3"`)}, false},
		{"incorrect-type", args{100, 100, []byte(`"CUSTOM_TYPE"`)}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := json.Unmarshal(tt.args.typeStr, &tt.args.defaultType)
			if (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.args.defaultType != tt.args.expType {
				t.Errorf("Expected: %d, got: %d", tt.args.expType, tt.args.defaultType)
			}
		})
	}
}
