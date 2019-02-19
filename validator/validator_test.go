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
		{"Document.Accessors[0].Sparse.Count", args{&gltf.Document{Asset: gltf.Asset{Version: "1.0"},
			Accessors: []gltf.Accessor{{ComponentType: gltf.Byte, Count: 1, Type: "SCALAR",
				Sparse: &gltf.Sparse{Count: 0}}}}}, true},
		{"Document.Accessors[0].Sparse.Indices.ComponentType", args{&gltf.Document{Asset: gltf.Asset{Version: "1.0"},
			Accessors: []gltf.Accessor{{ComponentType: gltf.Byte, Count: 1, Type: "SCALAR",
				Sparse: &gltf.Sparse{Count: 1, Indices: gltf.SparseIndices{ComponentType: 1}}}}}}, true},
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
			Nodes: []gltf.Node{{Rotation: [4]float32{2, 1, 1, 1}}}}}, true},
		{"Document.Nodes[0].Rotation[1]", args{&gltf.Document{Asset: gltf.Asset{Version: "1.0"},
			Nodes: []gltf.Node{{Rotation: [4]float32{1, -2, 1, 1}}}}}, true},
		{"Document.Skins[0].Joints", args{&gltf.Document{Asset: gltf.Asset{Version: "1.0"},
			Skins: []gltf.Skin{{Joints: []uint32{1, 1}}}}}, true},
		{"Document.Cameras[0].Type", args{&gltf.Document{Asset: gltf.Asset{Version: "1.0"},
			Cameras: []gltf.Camera{{Type: "incorrect"}}}}, true},
		{"Document.Cameras[0].Orthographic.Znear", args{&gltf.Document{Asset: gltf.Asset{Version: "1.0"},
			Cameras: []gltf.Camera{{Type: gltf.OrthographicType, Orthographic: &gltf.Orthographic{Znear: -1, Zfar: 1}}}}}, true},
		{"Document.Cameras[0].Orthographic.Zfar", args{&gltf.Document{Asset: gltf.Asset{Version: "1.0"},
			Cameras: []gltf.Camera{{Type: gltf.OrthographicType, Orthographic: &gltf.Orthographic{Znear: 1, Zfar: 1}}}}}, true},
		{"Document.Meshes[0].Primitives", args{&gltf.Document{Asset: gltf.Asset{Version: "1.0"},
			Meshes: []gltf.Mesh{{Primitives: make([]gltf.Primitive, 0)}}}}, true},
		{"Document.Meshes[0].Primitives[0].Mode", args{&gltf.Document{Asset: gltf.Asset{Version: "1.0"},
			Meshes: []gltf.Mesh{{Primitives: []gltf.Primitive{{Mode: 7}}}}}}, true},
		{"Document.Meshes[0].Primitives[0].Targets[0][OTHER]", args{&gltf.Document{Asset: gltf.Asset{Version: "1.0"},
			Meshes: []gltf.Mesh{{Primitives: []gltf.Primitive{{Targets: []gltf.Attribute{{"OTHER": 1}}}}}}}}, true},
		{"Document.Materials[0].EmissiveFactor[0]", args{&gltf.Document{Asset: gltf.Asset{Version: "1.0"},
			Materials: []gltf.Material{{EmissiveFactor: [3]float32{-1, 1, 1}}}}}, true},
		{"Document.Materials[0].EmissiveFactor[1]", args{&gltf.Document{Asset: gltf.Asset{Version: "1.0"},
			Materials: []gltf.Material{{AlphaMode: gltf.Opaque, EmissiveFactor: [3]float32{1, 2, 1}}}}}, true},
		{"Document.Materials[0].AlphaMode", args{&gltf.Document{Asset: gltf.Asset{Version: "1.0"},
			Materials: []gltf.Material{{AlphaMode: "OTHER"}}}}, true},
		{"Document.Materials[0].PBRMetallicRoughness.BaseColorFactor[0]", args{&gltf.Document{Asset: gltf.Asset{Version: "1.0"},
			Materials: []gltf.Material{{AlphaMode: gltf.Opaque,
				PBRMetallicRoughness: &gltf.PBRMetallicRoughness{BaseColorFactor: [4]float32{-1, 0, 0, 0}}}}}}, true},
		{"Document.Materials[0].PBRMetallicRoughness.BaseColorFactor[1]", args{&gltf.Document{Asset: gltf.Asset{Version: "1.0"},
			Materials: []gltf.Material{{AlphaMode: gltf.Opaque,
				PBRMetallicRoughness: &gltf.PBRMetallicRoughness{BaseColorFactor: [4]float32{1, 2, 0, 0}}}}}}, true},
		{"Document.Materials[0].PBRMetallicRoughness.MetallicFactor", args{&gltf.Document{Asset: gltf.Asset{Version: "1.0"},
			Materials: []gltf.Material{{AlphaMode: gltf.Opaque,
				PBRMetallicRoughness: &gltf.PBRMetallicRoughness{MetallicFactor: 2}}}}}, true},
		{"Document.Materials[0].PBRMetallicRoughness.MetallicFactor", args{&gltf.Document{Asset: gltf.Asset{Version: "1.0"},
			Materials: []gltf.Material{{AlphaMode: gltf.Opaque,
				PBRMetallicRoughness: &gltf.PBRMetallicRoughness{MetallicFactor: -1}}}}}, true},
		{"Document.Materials[0].PBRMetallicRoughness.RoughnessFactor", args{&gltf.Document{Asset: gltf.Asset{Version: "1.0"},
			Materials: []gltf.Material{{AlphaMode: gltf.Opaque,
				PBRMetallicRoughness: &gltf.PBRMetallicRoughness{RoughnessFactor: 2}}}}}, true},
		{"Document.Materials[0].PBRMetallicRoughness.RoughnessFactor", args{&gltf.Document{Asset: gltf.Asset{Version: "1.0"},
			Materials: []gltf.Material{{AlphaMode: gltf.Opaque,
				PBRMetallicRoughness: &gltf.PBRMetallicRoughness{RoughnessFactor: -1}}}}}, true},
		{"Document.Samplers[0].MagFilter", args{&gltf.Document{Asset: gltf.Asset{Version: "1.0"},
			Samplers: []gltf.Sampler{{MagFilter: 1, MinFilter: gltf.MinLinear, WrapS: gltf.ClampToEdge, WrapT: gltf.ClampToEdge}}}}, true},
		{"Document.Samplers[0].MinFilter", args{&gltf.Document{Asset: gltf.Asset{Version: "1.0"},
			Samplers: []gltf.Sampler{{MagFilter: gltf.MagLinear, MinFilter: 1, WrapS: gltf.ClampToEdge, WrapT: gltf.ClampToEdge}}}}, true},
		{"Document.Samplers[0].WrapS", args{&gltf.Document{Asset: gltf.Asset{Version: "1.0"},
			Samplers: []gltf.Sampler{{MagFilter: gltf.MagLinear, MinFilter: gltf.MinLinear, WrapS: 1, WrapT: gltf.ClampToEdge}}}}, true},
		{"Document.Samplers[0].WrapT", args{&gltf.Document{Asset: gltf.Asset{Version: "1.0"},
			Samplers: []gltf.Sampler{{MagFilter: gltf.MagLinear, MinFilter: gltf.MinLinear, WrapS: gltf.ClampToEdge, WrapT: 1}}}}, true},
		{"Document.Images[0].URI", args{&gltf.Document{Asset: gltf.Asset{Version: "1.0"},
			Images: []gltf.Image{{URI: "http://[web].com"}}}}, true},
		{"Document.Images[0].MimeType", args{&gltf.Document{Asset: gltf.Asset{Version: "1.0"},
			Images: []gltf.Image{{BufferView: 1}}}}, true},
		{"Document.Animations[0].Channels", args{&gltf.Document{Asset: gltf.Asset{Version: "1.0"},
			Animations: []gltf.Animation{{Samplers: []gltf.AnimationSampler{{}}}}}}, true},
		{"Document.Animations[0].Channels[0].Target.Path", args{&gltf.Document{Asset: gltf.Asset{Version: "1.0"},
			Animations: []gltf.Animation{{Channels: []gltf.Channel{{Target: gltf.ChannelTarget{Path: "other"}}}, Samplers: []gltf.AnimationSampler{{}}}}}}, true},
		{"Document.Animations[0].Samplers[0].Interpolation", args{&gltf.Document{Asset: gltf.Asset{Version: "1.0"},
			Animations: []gltf.Animation{{Channels: []gltf.Channel{{Target: gltf.ChannelTarget{Path: "translation"}}},
				Samplers: []gltf.AnimationSampler{{Interpolation: "OTHER"}}}}}}, true},
		{"ok", args{&gltf.Document{
			ExtensionsUsed:     []string{"one", "another"},
			ExtensionsRequired: []string{"that", "this"},
			Asset:              gltf.Asset{Copyright: "glTF", Generator: "qmuntal", Version: "1.0", MinVersion: "0.5"}},
		}, false},
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
