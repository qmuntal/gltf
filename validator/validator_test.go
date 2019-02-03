package validator

import (
	"gltf"
	"testing"

	val "github.com/go-playground/validator"
)

func TestValidateDocument(t *testing.T) {
	type args struct {
		doc *gltf.Document
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Document.Asset.Version", args{new(gltf.Document)}, true},
		{"Document.Accessors[0].ComponentType", args{&gltf.Document{Asset: gltf.Asset{Version: "1.0"},
			Accessors: []gltf.Accessor{{ComponentType: 1, Count: 1, Type: "SCALAR"}}}}, true},
		{"Document.Accessors[0].Count", args{&gltf.Document{Asset: gltf.Asset{Version: "1.0"},
			Accessors: []gltf.Accessor{{ComponentType: gltf.Byte, Count: 0, Type: "SCALAR"}}}}, true},
		{"Document.Accessors[0].Type", args{&gltf.Document{Asset: gltf.Asset{Version: "1.0"},
			Accessors: []gltf.Accessor{{ComponentType: gltf.Byte, Count: 1, Type: "OTHER"}}}}, true},
		{"Document.Accessors[0].Max", args{&gltf.Document{Asset: gltf.Asset{Version: "1.0"},
			Accessors: []gltf.Accessor{{ComponentType: gltf.Byte, Count: 1, Type: "SCALAR", Max: make([]float32, 17)}}}}, true},
		{"Document.Accessors[0].Min", args{&gltf.Document{Asset: gltf.Asset{Version: "1.0"},
			Accessors: []gltf.Accessor{{ComponentType: gltf.Byte, Count: 1, Type: "SCALAR", Min: make([]float32, 17)}}}}, true},
		{"Document.Accessors[0].Sparce.Count", args{&gltf.Document{Asset: gltf.Asset{Version: "1.0"},
			Accessors: []gltf.Accessor{{ComponentType: gltf.Byte, Count: 1, Type: "SCALAR",
				Sparce: &gltf.Sparse{Count: 0}}}}}, true},
		{"Document.Accessors[0].Sparce.Indices.ComponentType", args{&gltf.Document{Asset: gltf.Asset{Version: "1.0"},
			Accessors: []gltf.Accessor{{ComponentType: gltf.Byte, Count: 1, Type: "SCALAR",
				Sparce: &gltf.Sparse{Count: 1, Indices: gltf.SparseIndices{ComponentType: 1}}}}}}, true},
		{"Document.Buffers[0].URI", args{&gltf.Document{Asset: gltf.Asset{Version: "1.0"},
			Buffers: []gltf.Buffer{{ByteLength: 1, URI: "http://[web].com"}}}}, true},
		{"Document.Buffers[0].ByteLength", args{&gltf.Document{Asset: gltf.Asset{Version: "1.0"},
			Buffers: []gltf.Buffer{{ByteLength: 0, URI: "http://web.com"}}}}, true},
		{"Document.BufferViews[0].ByteLength", args{&gltf.Document{Asset: gltf.Asset{Version: "1.0"},
			BufferViews: []gltf.BufferView{{ByteLength: 0}}}}, true},
		{"Document.BufferViews[0].ByteStride", args{&gltf.Document{Asset: gltf.Asset{Version: "1.0"},
			BufferViews: []gltf.BufferView{{ByteLength: 1, ByteStride: 3}}}}, true},
		{"Document.BufferViews[0].ByteStride", args{&gltf.Document{Asset: gltf.Asset{Version: "1.0"},
			BufferViews: []gltf.BufferView{{ByteLength: 1, ByteStride: 253}}}}, true},
		{"Document.BufferViews[0].Target", args{&gltf.Document{Asset: gltf.Asset{Version: "1.0"},
			BufferViews: []gltf.BufferView{{ByteLength: 1, ByteStride: 4, Target: 2}}}}, true},
		{"Document.Scenes[0].Nodes", args{&gltf.Document{Asset: gltf.Asset{Version: "1.0"},
			Scenes: []gltf.Scene{{Nodes: []uint32{1, 1}}}}}, true},
		{"Document.Nodes[0].Children", args{&gltf.Document{Asset: gltf.Asset{Version: "1.0"},
			Nodes: []gltf.Node{{Children: []uint32{1, 1}}}}}, true},
		{"Document.Nodes[0].Rotation[0]", args{&gltf.Document{Asset: gltf.Asset{Version: "1.0"},
			Nodes: []gltf.Node{{Rotation: [4]float64{2, 1, 1, 1}}}}}, true},
		{"Document.Nodes[0].Rotation[1]", args{&gltf.Document{Asset: gltf.Asset{Version: "1.0"},
			Nodes: []gltf.Node{{Rotation: [4]float64{1, -2, 1, 1}}}}}, true},
		{"Document.Skins[0].Joints", args{&gltf.Document{Asset: gltf.Asset{Version: "1.0"},
			Skins: []gltf.Skin{{Joints: []uint32{1, 1}}}}}, true},
		{"Document.Cameras[0].Type", args{&gltf.Document{Asset: gltf.Asset{Version: "1.0"},
			Cameras: []gltf.Camera{{Type: "incorrect"}}}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDocument(tt.args.doc)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateDocument() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				errVal := err.(val.ValidationErrors)[0].Namespace()
				if errVal != tt.name {
					t.Errorf("ValidateDocument() error = %v, wantErr %v", errVal, tt.name)
				}
			}
		})
	}
}
