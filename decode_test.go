package gltf

import (
	"io/ioutil"
	"testing"

	"github.com/go-test/deep"
)

func TestOpen(t *testing.T) {
	cube, _ := ioutil.ReadFile("testdata/Cube/gltf/Cube.bin")
	type args struct {
		name string
	}
	tests := []struct {
		args    args
		want    *Document
		wantErr bool
	}{
		{args{"testdata/Cube/gltf/Cube.gltf"}, &Document{
			Accessors: []Accessor{
				{BufferView: 0, ByteOffset: 0, ComponentType: UnsignedShort, Count: 36, Max: []float32{35}, Min: []float32{0}, Type: Scalar},
				{BufferView: 1, ByteOffset: 0, ComponentType: Float, Count: 36, Max: []float32{1, 1, 1}, Min: []float32{-1, -1, -1}, Type: Vec3},
				{BufferView: 2, ByteOffset: 0, ComponentType: Float, Count: 36, Max: []float32{1, 1, 1}, Min: []float32{-1, -1, -1}, Type: Vec3},
				{BufferView: 3, ByteOffset: 0, ComponentType: Float, Count: 36, Max: []float32{1, 0, 0, 1}, Min: []float32{0, 0, -1, -1}, Type: Vec4},
				{BufferView: 4, ByteOffset: 0, ComponentType: Float, Count: 36, Max: []float32{1, 1}, Min: []float32{-1, -1}, Type: Vec2}},
			Asset: Asset{Generator: "VKTS glTF 2.0 exporter", Version: "2.0"},
			BufferViews: []BufferView{
				{Buffer: 0, ByteLength: 72, ByteOffset: 0, Target: ElementArrayBuffer},
				{Buffer: 0, ByteLength: 432, ByteOffset: 72, Target: ArrayBuffer},
				{Buffer: 0, ByteLength: 432, ByteOffset: 504, Target: ArrayBuffer},
				{Buffer: 0, ByteLength: 576, ByteOffset: 936, Target: ArrayBuffer},
				{Buffer: 0, ByteLength: 288, ByteOffset: 1512, Target: ArrayBuffer},
			},
			Buffers:   []Buffer{{ByteLength: 1800, URI: "Cube.bin", Data: cube}},
			Images:    []Image{{URI: "Cube_BaseColor.png"}, {URI: "Cube_MetallicRoughness.png"}},
			Materials: []Material{{Name: "Cube", AlphaMode: Opaque, AlphaCutoff: 0.5, PBRMetallicRoughness: &PBRMetallicRoughness{BaseColorFactor: [4]float32{1, 1, 1, 1}, MetallicFactor: 1, RoughnessFactor: 1, BaseColorTexture: &TextureInfo{Index: 0}, MetallicRoughnessTexture: &TextureInfo{Index: 1}}}},
			Meshes:    []Mesh{{Name: "Cube", Primitives: []Primitive{{Indices: 0, Material: 0, Mode: Triangles, Attributes: map[string]uint32{"NORMAL": 2, "POSITION": 1, "TANGENT": 3, "TEXCOORD_0": 4}}}}},
			Nodes:     []Node{{Mesh: 0, Name: "Cube", Camera: -1, Skin: -1, Matrix: [16]float32{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1}, Rotation: [4]float32{0, 0, 0, 1}, Scale: [3]float32{1, 1, 1}}},
			Samplers:  []Sampler{{WrapS: Repeat, WrapT: Repeat}},
			Scene:     0,
			Scenes:    []Scene{{Nodes: []uint32{0}}},
			Textures: []Texture{
				{Sampler: 0, Source: 0}, {Sampler: 0, Source: 1},
			},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.args.name, func(t *testing.T) {
			got, err := Open(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Open() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := deep.Equal(got, tt.want); diff != nil {
				t.Errorf("Open() = %v", diff)
			}
		})
	}
}
