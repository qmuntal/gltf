package draco

import (
	"reflect"
	"testing"

	"github.com/qmuntal/gltf"
)

func TestUnmarshal(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{"base", args{[]byte(`{
			"bufferView" : 5,
			"attributes" : {
				"POSITION" : 0,
				"NORMAL" : 1,
				"TEXCOORD_0" : 2,
				"WEIGHTS_0" : 3,
				"JOINTS_0" : 4
			}
		}`)}, &Primitive{BufferView: 5, Attributes: gltf.Attribute{
			"JOINTS_0":   4,
			"NORMAL":     1,
			"POSITION":   0,
			"TEXCOORD_0": 2,
			"WEIGHTS_0":  3,
		}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Unmarshal(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Unmarshal() = %v, want %v", got, tt.want)
			}
		})
	}
}
