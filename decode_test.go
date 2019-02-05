package gltf

import (
	"testing"
)

func TestOpen(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		args    args
		want    *Document
		wantErr bool
	}{
		{args{"./test/data/2.0/2CylinderEngine/glTF/2CylinderEngine.gltf"}, nil, false},
		{args{"./test/data/2.0/2CylinderEngine/glTF-Binary/2CylinderEngine.glb"}, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.args.name, func(t *testing.T) {
			_, err := Open(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Open() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// if !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("Open() = %v, want %v", got, tt.want)
			// }
		})
	}
}
