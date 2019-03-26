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
		Accessors: []gltf.Accessor{
			{BufferView: gltf.Index(0), ComponentType: gltf.UnsignedShort, Count: 36, Type: gltf.Scalar},
			{BufferView: gltf.Index(1), ComponentType: gltf.Float, Count: 24, Max: []float64{0.5, 0.5, 0.5}, Min: []float64{-0.5, -0.5, -0.5}, Type: gltf.Vec3},
			{BufferView: gltf.Index(2), ComponentType: gltf.Float, Count: 24, Type: gltf.Vec3},
			{BufferView: gltf.Index(3), ComponentType: gltf.Float, Count: 24, Type: gltf.Vec4},
			{BufferView: gltf.Index(4), ComponentType: gltf.Float, Count: 24, Type: gltf.Vec2},
		},
		Asset: gltf.Asset{Version: "2.0", Generator: "FBX2glTF"},
		BufferViews: []gltf.BufferView{
			{Buffer: 0, ByteLength: 72, ByteOffset: 0, Target: gltf.ElementArrayBuffer},
			{Buffer: 0, ByteLength: 288, ByteOffset: 72, Target: gltf.ArrayBuffer},
			{Buffer: 0, ByteLength: 288, ByteOffset: 360, Target: gltf.ArrayBuffer},
			{Buffer: 0, ByteLength: 384, ByteOffset: 648, Target: gltf.ArrayBuffer},
			{Buffer: 0, ByteLength: 192, ByteOffset: 1032, Target: gltf.ArrayBuffer},
		},
		Buffers: []gltf.Buffer{{ByteLength: 1224, Data: []uint8{97, 110, 121, 32, 99, 97, 114, 110, 97, 108, 32, 112, 108, 101, 97, 115}}},
		Materials: []gltf.Material{{
			Name: "Default", AlphaMode: gltf.Opaque, AlphaCutoff: gltf.Float64(0.5),
			PBRMetallicRoughness: &gltf.PBRMetallicRoughness{BaseColorFactor: &gltf.RGBA{R: 0.8, G: 0.8, B: 0.8, A: 1}, MetallicFactor: gltf.Float64(0.1), RoughnessFactor: gltf.Float64(0.99)},
		}},
		Meshes: []gltf.Mesh{{Name: "Cube", Primitives: []gltf.Primitive{{Indices: gltf.Index(0), Material: gltf.Index(0), Mode: gltf.Triangles, Attributes: map[string]uint32{"POSITION": 1, "COLOR_0": 3, "NORMAL": 2, "TEXCOORD_0": 4}}}}},
		Nodes: []gltf.Node{
			{Name: "RootNode", Children: []uint32{1, 2, 3}},
			{Name: "Mesh"},
			{Name: "Cube", Mesh: gltf.Index(0)},
			{Name: "Texture Group"},
		},
		Samplers: []gltf.Sampler{{WrapS: gltf.Repeat, WrapT: gltf.Repeat}},
		Scene:    gltf.Index(0),
		Scenes:   []gltf.Scene{{Name: "Root Scene", Nodes: []uint32{0}}},
	}
	if err := gltf.Save(doc, "./a.gltf", true); err != nil {
		panic(err)
	}
}
