package gltf

import (
	"testing"

	val "github.com/go-playground/validator"
)

func TestValidateDocument(t *testing.T) {
	tests := []struct {
		name    string
		doc     *Document
		wantErr bool
	}{
		{"Document.Asset.Version", new(Document), true},
		{"Document.Accessors[0].ComponentType", &Document{Asset: Asset{Version: "1.0"},
			Accessors: []Accessor{{ComponentType: 1, Count: 1, Type: "SCALAR"}}}, true},
		{"Document.Accessors[0].Count", &Document{Asset: Asset{Version: "1.0"},
			Accessors: []Accessor{{ComponentType: Byte, Count: 0, Type: "SCALAR"}}}, true},
		{"Document.Accessors[0].Type", &Document{Asset: Asset{Version: "1.0"},
			Accessors: []Accessor{{ComponentType: Byte, Count: 1, Type: "OTHER"}}}, true},
		{"Document.Accessors[0].Max", &Document{Asset: Asset{Version: "1.0"},
			Accessors: []Accessor{{ComponentType: Byte, Count: 1, Type: "SCALAR", Max: make([]float64, 17)}}}, true},
		{"Document.Accessors[0].Min", &Document{Asset: Asset{Version: "1.0"},
			Accessors: []Accessor{{ComponentType: Byte, Count: 1, Type: "SCALAR", Min: make([]float64, 17)}}}, true},
		{"Document.Accessors[0].Sparse.Count", &Document{Asset: Asset{Version: "1.0"},
			Accessors: []Accessor{{ComponentType: Byte, Count: 1, Type: "SCALAR",
				Sparse: &Sparse{Count: 0}}}}, true},
		{"Document.Accessors[0].Sparse.Indices.ComponentType", &Document{Asset: Asset{Version: "1.0"},
			Accessors: []Accessor{{ComponentType: Byte, Count: 1, Type: "SCALAR",
				Sparse: &Sparse{Count: 1, Indices: SparseIndices{ComponentType: 1}}}}}, true},
		{"Document.Buffers[0].URI", &Document{Asset: Asset{Version: "1.0"},
			Buffers: []Buffer{{ByteLength: 1, URI: "http://[web].com"}}}, true},
		{"Document.Buffers[0].ByteLength", &Document{Asset: Asset{Version: "1.0"},
			Buffers: []Buffer{{ByteLength: 0, URI: "http://web.com"}}}, true},
		{"Document.BufferViews[0].ByteLength", &Document{Asset: Asset{Version: "1.0"},
			BufferViews: []BufferView{{ByteLength: 0}}}, true},
		{"Document.BufferViews[0].ByteStride", &Document{Asset: Asset{Version: "1.0"},
			BufferViews: []BufferView{{ByteLength: 1, ByteStride: 3}}}, true},
		{"Document.BufferViews[0].ByteStride", &Document{Asset: Asset{Version: "1.0"},
			BufferViews: []BufferView{{ByteLength: 1, ByteStride: 253}}}, true},
		{"Document.BufferViews[0].Target", &Document{Asset: Asset{Version: "1.0"},
			BufferViews: []BufferView{{ByteLength: 1, ByteStride: 4, Target: 2}}}, true},
		{"Document.Scenes[0].Nodes", &Document{Asset: Asset{Version: "1.0"},
			Scenes: []Scene{{Nodes: []uint32{1, 1}}}}, true},
		{"Document.Nodes[0].Children", &Document{Asset: Asset{Version: "1.0"},
			Nodes: []Node{{Children: []uint32{1, 1}}}}, true},
		{"Document.Nodes[0].Rotation[0]", &Document{Asset: Asset{Version: "1.0"},
			Nodes: []Node{{Rotation: [4]float64{2, 1, 1, 1}}}}, true},
		{"Document.Nodes[0].Rotation[1]", &Document{Asset: Asset{Version: "1.0"},
			Nodes: []Node{{Rotation: [4]float64{1, -2, 1, 1}}}}, true},
		{"Document.Skins[0].Joints", &Document{Asset: Asset{Version: "1.0"},
			Skins: []Skin{{Joints: []uint32{1, 1}}}}, true},
		{"Document.Cameras[0].Type", &Document{Asset: Asset{Version: "1.0"},
			Cameras: []Camera{{Type: "incorrect"}}}, true},
		{"Document.Cameras[0].Orthographic.Znear", &Document{Asset: Asset{Version: "1.0"},
			Cameras: []Camera{{Type: OrthographicType, Orthographic: &Orthographic{Znear: -1, Zfar: 1}}}}, true},
		{"Document.Cameras[0].Orthographic.Zfar", &Document{Asset: Asset{Version: "1.0"},
			Cameras: []Camera{{Type: OrthographicType, Orthographic: &Orthographic{Znear: 1, Zfar: 1}}}}, true},
		{"Document.Meshes[0].Primitives", &Document{Asset: Asset{Version: "1.0"},
			Meshes: []Mesh{{Primitives: make([]Primitive, 0)}}}, true},
		{"Document.Meshes[0].Primitives[0].Mode", &Document{Asset: Asset{Version: "1.0"},
			Meshes: []Mesh{{Primitives: []Primitive{{Mode: 7}}}}}, true},
		{"Document.Meshes[0].Primitives[0].Targets[0][OTHER]", &Document{Asset: Asset{Version: "1.0"},
			Meshes: []Mesh{{Primitives: []Primitive{{Targets: []Attribute{{"OTHER": 1}}}}}}}, true},
		{"Document.Materials[0].EmissiveFactor[0]", &Document{Asset: Asset{Version: "1.0"},
			Materials: []Material{{EmissiveFactor: [3]float64{-1, 1, 1}}}}, true},
		{"Document.Materials[0].EmissiveFactor[1]", &Document{Asset: Asset{Version: "1.0"},
			Materials: []Material{{AlphaMode: Opaque, EmissiveFactor: [3]float64{1, 2, 1}}}}, true},
		{"Document.Materials[0].AlphaMode", &Document{Asset: Asset{Version: "1.0"},
			Materials: []Material{{AlphaMode: "OTHER"}}}, true},
		{"Document.Materials[0].PBRMetallicRoughness.BaseColorFactor[0]", &Document{Asset: Asset{Version: "1.0"},
			Materials: []Material{{AlphaMode: Opaque,
				PBRMetallicRoughness: &PBRMetallicRoughness{BaseColorFactor: [4]float64{-1, 0, 0, 0}}}}}, true},
		{"Document.Materials[0].PBRMetallicRoughness.BaseColorFactor[1]", &Document{Asset: Asset{Version: "1.0"},
			Materials: []Material{{AlphaMode: Opaque,
				PBRMetallicRoughness: &PBRMetallicRoughness{BaseColorFactor: [4]float64{1, 2, 0, 0}}}}}, true},
		{"Document.Materials[0].PBRMetallicRoughness.MetallicFactor", &Document{Asset: Asset{Version: "1.0"},
			Materials: []Material{{AlphaMode: Opaque,
				PBRMetallicRoughness: &PBRMetallicRoughness{MetallicFactor: 2}}}}, true},
		{"Document.Materials[0].PBRMetallicRoughness.MetallicFactor", &Document{Asset: Asset{Version: "1.0"},
			Materials: []Material{{AlphaMode: Opaque,
				PBRMetallicRoughness: &PBRMetallicRoughness{MetallicFactor: -1}}}}, true},
		{"Document.Materials[0].PBRMetallicRoughness.RoughnessFactor", &Document{Asset: Asset{Version: "1.0"},
			Materials: []Material{{AlphaMode: Opaque,
				PBRMetallicRoughness: &PBRMetallicRoughness{RoughnessFactor: 2}}}}, true},
		{"Document.Materials[0].PBRMetallicRoughness.RoughnessFactor", &Document{Asset: Asset{Version: "1.0"},
			Materials: []Material{{AlphaMode: Opaque,
				PBRMetallicRoughness: &PBRMetallicRoughness{RoughnessFactor: -1}}}}, true},
		{"Document.Samplers[0].MagFilter", &Document{Asset: Asset{Version: "1.0"},
			Samplers: []Sampler{{MagFilter: 1, MinFilter: MinLinear, WrapS: ClampToEdge, WrapT: ClampToEdge}}}, true},
		{"Document.Samplers[0].MinFilter", &Document{Asset: Asset{Version: "1.0"},
			Samplers: []Sampler{{MagFilter: MagLinear, MinFilter: 1, WrapS: ClampToEdge, WrapT: ClampToEdge}}}, true},
		{"Document.Samplers[0].WrapS", &Document{Asset: Asset{Version: "1.0"},
			Samplers: []Sampler{{MagFilter: MagLinear, MinFilter: MinLinear, WrapS: 1, WrapT: ClampToEdge}}}, true},
		{"Document.Samplers[0].WrapT", &Document{Asset: Asset{Version: "1.0"},
			Samplers: []Sampler{{MagFilter: MagLinear, MinFilter: MinLinear, WrapS: ClampToEdge, WrapT: 1}}}, true},
		{"Document.Images[0].URI", &Document{Asset: Asset{Version: "1.0"},
			Images: []Image{{URI: "http://[web].com"}}}, true},
		{"Document.Images[0].MimeType", &Document{Asset: Asset{Version: "1.0"},
			Images: []Image{{BufferView: 1}}}, true},
		{"Document.Animations[0].Channels", &Document{Asset: Asset{Version: "1.0"},
			Animations: []Animation{{Samplers: []AnimationSampler{{}}}}}, true},
		{"Document.Animations[0].Channels[0].Target.Path", &Document{Asset: Asset{Version: "1.0"},
			Animations: []Animation{{Channels: []Channel{{Target: ChannelTarget{Path: "other"}}}, Samplers: []AnimationSampler{{}}}}}, true},
		{"Document.Animations[0].Samplers[0].Interpolation", &Document{Asset: Asset{Version: "1.0"},
			Animations: []Animation{{Channels: []Channel{{Target: ChannelTarget{Path: "translation"}}},
				Samplers: []AnimationSampler{{Interpolation: "OTHER"}}}}}, true},
		{"ok", &Document{
			ExtensionsUsed:     []string{"one", "another"},
			ExtensionsRequired: []string{"that", "this"},
			Asset:              Asset{Copyright: "glTF", Generator: "qmuntal", Version: "1.0", MinVersion: "0.5"},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.doc.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Document.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				errVal := err.(val.ValidationErrors)[0].Namespace()
				if errVal != tt.name {
					t.Errorf("Document.Validate() error = %v, wantErr %v", errVal, tt.name)
				}
			}
		})
	}
}
