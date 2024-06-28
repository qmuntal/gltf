package gltf_test

import (
	"fmt"

	"github.com/qmuntal/gltf"
)

func ExampleOpen() {
	doc, err := gltf.Open("fake")
	if err != nil {
		panic(err)
	}
	fmt.Print(doc.Asset)
}

func ExampleSave() {
	doc := &gltf.Document{
		Accessors: []*gltf.Accessor{
			{BufferView: gltf.Index(0), ComponentType: gltf.ComponentUshort, Count: 36, Type: gltf.AccessorScalar},
			{BufferView: gltf.Index(1), ComponentType: gltf.ComponentFloat, Count: 24, Max: []float64{0.5, 0.5, 0.5}, Min: []float64{-0.5, -0.5, -0.5}, Type: gltf.AccessorVec3},
			{BufferView: gltf.Index(2), ComponentType: gltf.ComponentFloat, Count: 24, Type: gltf.AccessorVec3},
			{BufferView: gltf.Index(3), ComponentType: gltf.ComponentFloat, Count: 24, Type: gltf.AccessorVec4},
			{BufferView: gltf.Index(4), ComponentType: gltf.ComponentFloat, Count: 24, Type: gltf.AccessorVec2},
		},
		Asset: gltf.Asset{Version: "2.0", Generator: "FBX2glTF"},
		BufferViews: []*gltf.BufferView{
			{Buffer: 0, ByteLength: 72, ByteOffset: 0, Target: gltf.TargetElementArrayBuffer},
			{Buffer: 0, ByteLength: 288, ByteOffset: 72, Target: gltf.TargetArrayBuffer},
			{Buffer: 0, ByteLength: 288, ByteOffset: 360, Target: gltf.TargetArrayBuffer},
			{Buffer: 0, ByteLength: 384, ByteOffset: 648, Target: gltf.TargetArrayBuffer},
			{Buffer: 0, ByteLength: 192, ByteOffset: 1032, Target: gltf.TargetArrayBuffer},
		},
		Buffers: []*gltf.Buffer{{ByteLength: 1224, Data: []byte{97, 110, 121, 32, 99, 97, 114, 110, 97, 108, 32, 112, 108, 101, 97, 115}}},
		Materials: []*gltf.Material{{
			Name: "Default", AlphaMode: gltf.AlphaOpaque, AlphaCutoff: gltf.Float(0.5),
			PBRMetallicRoughness: &gltf.PBRMetallicRoughness{BaseColorFactor: &[4]float64{0.8, 0.8, 0.8, 1}, MetallicFactor: gltf.Float(0.1), RoughnessFactor: gltf.Float(0.99)},
		}},
		Meshes: []*gltf.Mesh{{Name: "Cube", Primitives: []*gltf.Primitive{{Indices: gltf.Index(0), Material: gltf.Index(0), Mode: gltf.PrimitiveTriangles, Attributes: gltf.PrimitiveAttributes{gltf.POSITION: 1, gltf.COLOR_0: 3, gltf.NORMAL: 2, gltf.TEXCOORD_0: 4}}}}},
		Nodes: []*gltf.Node{
			{Name: "RootNode", Children: []uint32{1, 2, 3}},
			{Name: "Mesh"},
			{Name: "Cube", Mesh: gltf.Index(0)},
			{Name: "Texture Group"},
		},
		Samplers: []*gltf.Sampler{{WrapS: gltf.WrapRepeat, WrapT: gltf.WrapRepeat}},
		Scene:    gltf.Index(0),
		Scenes:   []*gltf.Scene{{Name: "Root Scene", Nodes: []uint32{0}}},
	}
	if err := gltf.SaveBinary(doc, "./example.glb"); err != nil {
		panic(err)
	}
}
