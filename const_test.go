package gltf

import (
	"encoding/json"
	"testing"
)

func TestAccessorType_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		expType AccessorType
		typeStr []byte
		wantErr bool
	}{
		{"base", AccessorVec3, []byte(`"VEC3"`), false},
		{"incorrect-type", 100, []byte(`"CUSTOM_TYPE"`), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var accType AccessorType = 100
			err := json.Unmarshal(tt.typeStr, &accType)
			if (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if accType != tt.expType {
				t.Errorf("Expected: %d, got: %d", AccessorVec3, accType)
			}
		})
	}
}
