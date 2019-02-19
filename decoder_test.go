package gltf

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"testing"

	"github.com/go-test/deep"
)

func readFile(path string) []uint8 {
	cube, _ := ioutil.ReadFile(path)
	return cube
}

func TestOpen(t *testing.T) {
	type args struct {
		name     string
		embedded string
	}
	tests := []struct {
		args    args
		want    *Document
		wantErr bool
	}{
		{args{"openError", ""}, nil, true},
		{args{"testdata/Cube/glTF/Cube.gltf", ""}, &Document{
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
			Buffers:   []Buffer{{ByteLength: 1800, URI: "Cube.bin", Data: readFile("testdata/Cube/glTF/Cube.bin")}},
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
		{args{"testdata/Cameras/glTF/Cameras.gltf", "testdata/Cameras/glTF-Embedded/Cameras.gltf"}, &Document{
			Accessors: []Accessor{
				{BufferView: 0, ByteOffset: 0, ComponentType: UnsignedShort, Count: 6, Max: []float32{3}, Min: []float32{0}, Type: Scalar},
				{BufferView: 1, ByteOffset: 0, ComponentType: Float, Count: 4, Max: []float32{1, 1, 0}, Min: []float32{0, 0, 0}, Type: Vec3},
			},
			Asset: Asset{Version: "2.0"},
			BufferViews: []BufferView{
				{Buffer: 0, ByteLength: 12, ByteOffset: 0, Target: ElementArrayBuffer},
				{Buffer: 0, ByteLength: 48, ByteOffset: 12, Target: ArrayBuffer},
			},
			Buffers: []Buffer{{ByteLength: 60, URI: "simpleSquare.bin", Data: readFile("testdata/Cameras/glTF/simpleSquare.bin")}},
			Cameras: []Camera{
				{Type: PerspectiveType, Perspective: &Perspective{AspectRatio: 1.0, Yfov: 0.7, Zfar: 100, Znear: 0.01}},
				{Type: OrthographicType, Orthographic: &Orthographic{Xmag: 1.0, Ymag: 1.0, Zfar: 100, Znear: 0.01}},
			},
			Meshes: []Mesh{{Primitives: []Primitive{{Indices: 0, Material: -1, Mode: Triangles, Attributes: map[string]uint32{"POSITION": 1}}}}},
			Nodes: []Node{
				{Mesh: 0, Camera: -1, Skin: -1, Matrix: [16]float32{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1}, Rotation: [4]float32{-0.3, 0, 0, 0.9}, Scale: [3]float32{1, 1, 1}},
				{Mesh: -1, Camera: 0, Skin: -1, Matrix: [16]float32{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1}, Rotation: [4]float32{0, 0, 0, 1}, Scale: [3]float32{1, 1, 1}, Translation: [3]float32{0.5, 0.5, 3.0}},
				{Mesh: -1, Camera: 1, Skin: -1, Matrix: [16]float32{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1}, Rotation: [4]float32{0, 0, 0, 1}, Scale: [3]float32{1, 1, 1}, Translation: [3]float32{0.5, 0.5, 3.0}},
			},
			Scene:  -1,
			Scenes: []Scene{{Nodes: []uint32{0, 1, 2}}},
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
				return
			}
			if tt.args.embedded != "" {
				got, err = Open(tt.args.embedded)
				for i, b := range got.Buffers {
					if b.IsEmbeddedResource() {
						tt.want.Buffers[i].EmbeddedResource()
					}
				}
				if (err != nil) != tt.wantErr {
					t.Errorf("Open() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if diff := deep.Equal(got, tt.want); diff != nil {
					t.Errorf("Open() = %v", diff)
					return
				}
			}
		})
	}
}

func readCallback(name string) (io.ReadCloser, error) {
	return ioutil.NopCloser(bytes.NewBufferString("a")), nil
}

func TestDecoder_decodeBuffer(t *testing.T) {
	type args struct {
		buffer *Buffer
	}
	tests := []struct {
		name    string
		d       *Decoder
		args    args
		wantErr bool
	}{
		{"byteLength_0", &Decoder{quotas: ReadQuotas{MaxMemoryAllocation: 2}}, args{&Buffer{ByteLength: 0, URI: "a.bin"}}, true},
		{"noURI", &Decoder{quotas: ReadQuotas{MaxMemoryAllocation: 2}}, args{&Buffer{ByteLength: 1, URI: ""}}, true},
		{"invalidURI", &Decoder{quotas: ReadQuotas{MaxMemoryAllocation: 2}}, args{&Buffer{ByteLength: 1, URI: "../a.bin"}}, true},
		{"maxQuota", &Decoder{quotas: ReadQuotas{MaxMemoryAllocation: 2}}, args{&Buffer{ByteLength: 3, URI: "a.bin"}}, true},
		{"cbErr", NewDecoder(nil, func(name string) (io.ReadCloser, error) { return nil, errors.New("") }), args{&Buffer{ByteLength: 3, URI: "a.bin"}}, true},
		{"base", NewDecoder(nil, readCallback), args{&Buffer{ByteLength: 3, URI: "a.bin"}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.d.decodeBuffer(tt.args.buffer); (err != nil) != tt.wantErr {
				t.Errorf("Decoder.decodeBuffer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDecoder_decodeBinaryBuffer(t *testing.T) {
	type args struct {
		buffer *Buffer
	}
	tests := []struct {
		name    string
		d       *Decoder
		args    args
		wantErr bool
	}{
		{"invalidBuffer", new(Decoder), args{&Buffer{ByteLength: 0, URI: "a.bin"}}, true},
		{"readErr", NewDecoder(bytes.NewBufferString(""), nil), args{&Buffer{ByteLength: 1, URI: "a.bin"}}, true},
		{"invalidHeader", NewDecoder(bytes.NewBufferString("aaaaaaaa"), nil), args{&Buffer{ByteLength: 1, URI: "a.bin"}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.d.decodeBinaryBuffer(tt.args.buffer); (err != nil) != tt.wantErr {
				t.Errorf("Decoder.decodeBinaryBuffer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDecoder_Decode(t *testing.T) {
	type args struct {
		doc *Document
	}
	tests := []struct {
		name    string
		d       *Decoder
		args    args
		wantErr bool
	}{
		{"empty", NewDecoder(bytes.NewBufferString(""), nil), args{new(Document)}, true},
		{"invalidJSON", NewDecoder(bytes.NewBufferString("{asset: {}}"), nil), args{new(Document)}, true},
		{"invalidBuffer", NewDecoder(bytes.NewBufferString("{\"buffers\": [{\"byteLength\": 0}]}"), nil), args{new(Document)}, true},
		{"maxBuffers", NewDecoder(bytes.NewBufferString("{\"buffers\": [{\"byteLength\": 0}]}"), nil).SetQuotas(ReadQuotas{MaxBufferCount: 0}), args{new(Document)}, true},
		{"base", NewDecoder(bytes.NewBufferString("{\"buffers\": [{\"byteLength\": 1, \"URI\": \"a.bin\"}]}"), readCallback), args{new(Document)}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.d.Decode(tt.args.doc); (err != nil) != tt.wantErr {
				t.Errorf("Decoder.Decode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
