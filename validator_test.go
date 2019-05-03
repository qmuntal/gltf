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
			Accessors: []Accessor{{ComponentType: 10, Count: 1}}}, true},
		{"Document.Accessors[0].Count", &Document{Asset: Asset{Version: "1.0"},
			Accessors: []Accessor{{ComponentType: Byte, Count: 0}}}, true},
		{"Document.Accessors[0].Type", &Document{Asset: Asset{Version: "1.0"},
			Accessors: []Accessor{{ComponentType: Byte, Count: 1, Type: 10}}}, true},
		{"Document.Accessors[0].Max", &Document{Asset: Asset{Version: "1.0"},
			Accessors: []Accessor{{ComponentType: Byte, Count: 1, Max: make([]float64, 17)}}}, true},
		{"Document.Accessors[0].Min", &Document{Asset: Asset{Version: "1.0"},
			Accessors: []Accessor{{ComponentType: Byte, Count: 1, Min: make([]float64, 17)}}}, true},
		{"Document.Accessors[0].Sparse.Count", &Document{Asset: Asset{Version: "1.0"},
			Accessors: []Accessor{{ComponentType: Byte, Count: 1,
				Sparse: &Sparse{Count: 0}}}}, true},
		{"Document.Accessors[0].Sparse.Indices.ComponentType", &Document{Asset: Asset{Version: "1.0"},
			Accessors: []Accessor{{ComponentType: Byte, Count: 1,
				Sparse: &Sparse{Count: 1, Indices: SparseIndices{ComponentType: 1}}}}}, true},
		{"Document.Buffers[0].URI", &Document{Asset: Asset{Version: "1.0"},
			Buffers: []Buffer{{ByteLength: 1, URI: "a.bin"}}}, false},
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
		{"Document.Cameras[0].Orthographic.Znear", &Document{Asset: Asset{Version: "1.0"},
			Cameras: []Camera{{Orthographic: &Orthographic{Znear: -1, Zfar: 1}}}}, true},
		{"Document.Cameras[0].Orthographic.Zfar", &Document{Asset: Asset{Version: "1.0"},
			Cameras: []Camera{{Orthographic: &Orthographic{Znear: 1, Zfar: 1}}}}, true},
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
			Materials: []Material{{AlphaMode: 5}}}, true},
		{"Document.Materials[0].PBRMetallicRoughness.BaseColorFactor.R", &Document{Asset: Asset{Version: "1.0"},
			Materials: []Material{{AlphaMode: Opaque,
				PBRMetallicRoughness: &PBRMetallicRoughness{BaseColorFactor: &RGBA{R: -0.1, G: 0.8, B: 0.8, A: 1}}}}}, true},
		{"Document.Materials[0].PBRMetallicRoughness.BaseColorFactor.G", &Document{Asset: Asset{Version: "1.0"},
			Materials: []Material{{AlphaMode: Opaque,
				PBRMetallicRoughness: &PBRMetallicRoughness{BaseColorFactor: &RGBA{R: 1, G: 1.1, B: 0, A: 0}}}}}, true},
		{"Document.Materials[0].PBRMetallicRoughness.MetallicFactor", &Document{Asset: Asset{Version: "1.0"},
			Materials: []Material{{AlphaMode: Opaque,
				PBRMetallicRoughness: &PBRMetallicRoughness{MetallicFactor: Float64(2)}}}}, true},
		{"Document.Materials[0].PBRMetallicRoughness.MetallicFactor", &Document{Asset: Asset{Version: "1.0"},
			Materials: []Material{{AlphaMode: Opaque,
				PBRMetallicRoughness: &PBRMetallicRoughness{MetallicFactor: Float64(-1)}}}}, true},
		{"Document.Materials[0].PBRMetallicRoughness.RoughnessFactor", &Document{Asset: Asset{Version: "1.0"},
			Materials: []Material{{AlphaMode: Opaque,
				PBRMetallicRoughness: &PBRMetallicRoughness{RoughnessFactor: Float64(2)}}}}, true},
		{"Document.Materials[0].PBRMetallicRoughness.RoughnessFactor", &Document{Asset: Asset{Version: "1.0"},
			Materials: []Material{{AlphaMode: Opaque,
				PBRMetallicRoughness: &PBRMetallicRoughness{RoughnessFactor: Float64(-1)}}}}, true},
		{"Document.Samplers[0].MagFilter", &Document{Asset: Asset{Version: "1.0"},
			Samplers: []Sampler{{MagFilter: 10, MinFilter: MinLinear, WrapS: ClampToEdge, WrapT: ClampToEdge}}}, true},
		{"Document.Samplers[0].MinFilter", &Document{Asset: Asset{Version: "1.0"},
			Samplers: []Sampler{{MagFilter: MagLinear, MinFilter: 10, WrapS: ClampToEdge, WrapT: ClampToEdge}}}, true},
		{"Document.Samplers[0].WrapS", &Document{Asset: Asset{Version: "1.0"},
			Samplers: []Sampler{{MagFilter: MagLinear, MinFilter: MinLinear, WrapS: 10, WrapT: ClampToEdge}}}, true},
		{"Document.Samplers[0].WrapT", &Document{Asset: Asset{Version: "1.0"},
			Samplers: []Sampler{{MagFilter: MagLinear, MinFilter: MinLinear, WrapS: ClampToEdge, WrapT: 10}}}, true},
		{"Document.Images[0].URI", &Document{Asset: Asset{Version: "1.0"},
			Images: []Image{{URI: "a.png"}}}, false},
		{"Document.Images[0].MimeType", &Document{Asset: Asset{Version: "1.0"},
			Images: []Image{{BufferView: Index(1)}}}, true},
		{"Document.Animations[0].Channels", &Document{Asset: Asset{Version: "1.0"},
			Animations: []Animation{{Samplers: []AnimationSampler{{}}}}}, true},
		{"Document.Animations[0].Channels[0].Target.Path", &Document{Asset: Asset{Version: "1.0"},
			Animations: []Animation{{Channels: []Channel{{Target: ChannelTarget{Path: 10}}}, Samplers: []AnimationSampler{{}}}}}, true},
		{"Document.Animations[0].Samplers[0].Interpolation", &Document{Asset: Asset{Version: "1.0"},
			Animations: []Animation{{Channels: []Channel{{Target: ChannelTarget{Path: Translation}}},
				Samplers: []AnimationSampler{{Interpolation: 10}}}}}, true},
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
