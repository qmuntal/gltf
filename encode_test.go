package gltf

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"testing"

	"github.com/go-test/deep"
)

type writeCloser struct {
	io.Writer
}

func (w *writeCloser) Close() error { return nil }

func saveMemory(doc *Document, asBinary bool) (*Decoder, error) {
	buff := new(bytes.Buffer)
	chunks := make(map[string]*bytes.Buffer)
	wcb := func(uri string, size int) (io.WriteCloser, error) {
		chunks[uri] = bytes.NewBuffer(make([]byte, 0, size))
		return &writeCloser{chunks[uri]}, nil
	}
	if err := NewEncoder(buff, wcb, asBinary).Encode(doc); err != nil {
		return nil, err
	}
	rcb := func(uri string) (io.ReadCloser, error) {
		return ioutil.NopCloser(chunks[uri]), nil
	}
	return NewDecoder(buff, rcb), nil
}

func TestEncoder_Encode(t *testing.T) {
	type args struct {
		doc *Document
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"withInvalidBuffer", args{&Document{Buffers: []Buffer{
			{Extras: 8.0, Extensions: Extensions{"a": "b"}, Name: "binary", ByteLength: 3, URI: "a.bin", Data: []uint8{1, 2, 3}},
			{Extras: 8.0, Extensions: Extensions{"a": "b"}, Name: "binary2", ByteLength: 3, URI: "/../a.bin", Data: []uint8{1, 2, 3}},
		}}}, true},
		{"empty", args{&Document{}}, false},
		{"withExtensions", args{&Document{Extras: 8.0, Extensions: Extensions{"a": "b"}, ExtensionsUsed: []string{"c"}, ExtensionsRequired: []string{"d", "e"}}}, false},
		{"withAsset", args{&Document{Asset: Asset{Extras: 8.0, Extensions: Extensions{"a": "b"}, Copyright: "@2019", Generator: "qmuntal/gltf", Version: "2.0", MinVersion: "1.0"}}}, false},
		{"withAccessors", args{&Document{Accessors: []Accessor{
			{Extras: 8.0, Extensions: Extensions{"a": "b"}, Name: "acc_1", BufferView: 0, ByteOffset: 50, ComponentType: Byte, Normalized: true, Count: 5, Type: Vec3, Max: []float64{1, 2}, Min: []float64{2.4}},
			{BufferView: 0, Normalized: false, Count: 50, Type: Vec4, Sparse: &Sparse{Extras: 8.0, Extensions: Extensions{"a": "b"}, Count: 2,
				Values:  SparseValues{Extras: 8.0, Extensions: Extensions{"a": "b"}, BufferView: 1, ByteOffset: 2},
				Indices: SparseIndices{Extras: 8.0, Extensions: Extensions{"a": "b"}, BufferView: 1, ByteOffset: 2, ComponentType: Float}},
			},
		}}}, false},
		{"withAnimations", args{&Document{Animations: []Animation{
			{Extras: 8.0, Extensions: Extensions{"a": "b"}, Name: "an_1", Channels: []Channel{
				{Extras: 8.0, Extensions: Extensions{"a": "b"}, Sampler: 1, Target: ChannelTarget{Extras: 8.0, Extensions: Extensions{"a": "b"}, Node: 10, Path: Rotation}},
				{Extras: 8.0, Extensions: Extensions{"a": "b"}, Sampler: 2, Target: ChannelTarget{Extras: 8.0, Extensions: Extensions{"a": "b"}, Node: 10, Path: Scale}},
			}},
			{Extras: 8.0, Extensions: Extensions{"a": "b"}, Name: "an_2", Channels: []Channel{
				{Extras: 8.0, Extensions: Extensions{"a": "b"}, Sampler: 1, Target: ChannelTarget{Extras: 8.0, Extensions: Extensions{"a": "b"}, Node: 3, Path: Weights}},
				{Extras: 8.0, Extensions: Extensions{"a": "b"}, Sampler: 2, Target: ChannelTarget{Extras: 8.0, Extensions: Extensions{"a": "b"}, Node: 5, Path: Translation}},
			}},
		}}}, false},
		{"withBuffer", args{&Document{Buffers: []Buffer{
			{Extras: 8.0, Extensions: Extensions{"a": "b"}, Name: "binary", ByteLength: 3, URI: "a.bin", Data: []uint8{1, 2, 3}},
			{Extras: 8.0, Extensions: Extensions{"a": "b"}, Name: "embedded", ByteLength: 2, URI: "data:application/octet-stream;base64,YW55ICsgb2xkICYgZGF0YQ==", Data: []byte("any + old & data")},
			{Extras: 8.0, Extensions: Extensions{"a": "b"}, Name: "external", ByteLength: 4, URI: "b.bin", Data: []uint8{4, 5, 6, 7}},
		}}}, false},
		{"withBufView", args{&Document{BufferViews: []BufferView{
			{Extras: 8.0, Extensions: Extensions{"a": "b"}, Buffer: 0, ByteOffset: 1, ByteLength: 2, ByteStride: 5, Target: ArrayBuffer},
			{Buffer: 10, ByteOffset: 10, ByteLength: 20, ByteStride: 50, Target: ElementArrayBuffer},
		}}}, false},
		{"withCameras", args{&Document{Cameras: []Camera{
			{Extras: 8.0, Extensions: Extensions{"a": "b"}, Name: "cam_1", Type: OrthographicType, Orthographic: &Orthographic{Extras: 8.0, Extensions: Extensions{"a": "b"}, Xmag: 1, Ymag: 2, Zfar: 3, Znear: 4}},
			{Extras: 8.0, Extensions: Extensions{"a": "b"}, Name: "cam_2", Type: PerspectiveType, Perspective: &Perspective{Extras: 8.0, Extensions: Extensions{"a": "b"}, AspectRatio: 1, Yfov: 2, Zfar: 3, Znear: 4}},
		}}}, false},
		{"withImages", args{&Document{Images: []Image{
			{Extras: 8.0, Extensions: Extensions{"a": "b"}, Name: "binary", BufferView: 1, MimeType: "data:image/png"},
			{Extras: 8.0, Extensions: Extensions{"a": "b"}, Name: "embedded", URI: "data:image/png;base64,dsjdsaGGUDXGA", MimeType: "data:image/png"},
			{Extras: 8.0, Extensions: Extensions{"a": "b"}, Name: "external", URI: "https://web.com/a", MimeType: "data:image/png"},
		}}}, false},
		{"withMaterials", args{&Document{Materials: []Material{
			{Extras: 8.0, Extensions: Extensions{"a": "b"}, Name: "base", EmissiveFactor: [3]float64{1.0, 1.0, 1.0}, DoubleSided: true, AlphaCutoff: 0.5, AlphaMode: Opaque},
			{Extras: 8.0, Extensions: Extensions{"a": "b"}, Name: "pbr", AlphaCutoff: 0.5, AlphaMode: Opaque,
				PBRMetallicRoughness: &PBRMetallicRoughness{
					Extras: 8.0, Extensions: Extensions{"a": "b"}, MetallicFactor: 1, RoughnessFactor: 2, BaseColorFactor: [4]float64{1, 2, 3, 4},
					BaseColorTexture:         &TextureInfo{Extras: 8.0, Extensions: Extensions{"a": "b"}, Index: 1, TexCoord: 3},
					MetallicRoughnessTexture: &TextureInfo{Extras: 8.0, Extensions: Extensions{"a": "b"}, Index: 6, TexCoord: 5},
				},
			},
			{Extras: 8.0, Extensions: Extensions{"a": "b"}, Name: "normal", AlphaCutoff: 0.7, AlphaMode: Blend,
				NormalTexture: &NormalTexture{Extras: 8.0, Extensions: Extensions{"a": "b"}, Index: 1, TexCoord: 2, Scale: 2.0},
			},
			{Extras: 8.0, Extensions: Extensions{"a": "b"}, Name: "occlusion", AlphaCutoff: 0.5, AlphaMode: Mask,
				OcclusionTexture: &OcclusionTexture{Extras: 8.0, Extensions: Extensions{"a": "b"}, Index: 1, TexCoord: 2, Strength: 2.0},
			},
			{Extras: 8.0, Extensions: Extensions{"a": "b"}, Name: "emmisice", AlphaCutoff: 0.5, AlphaMode: Mask, EmissiveTexture: &TextureInfo{Extras: 8.0, Extensions: Extensions{"a": "b"}, Index: 4, TexCoord: 50}},
		}}}, false},
		{"withMeshes", args{&Document{Meshes: []Mesh{
			{Extras: 8.0, Extensions: Extensions{"a": "b"}, Name: "mesh_1", Weights: []float64{1.2, 2}},
			{Extras: 8.0, Extensions: Extensions{"a": "b"}, Name: "mesh_2", Primitives: []Primitive{
				{Extras: 8.0, Extensions: Extensions{"a": "b"}, Attributes: Attribute{"POSITION": 1}, Indices: 2, Material: 1, Mode: Lines},
				{Extras: 8.0, Extensions: Extensions{"a": "b"}, Targets: []Attribute{{"POSITION": 1, "THEN": 4}, {"OTHER": 2}}, Indices: 2, Material: 1, Mode: Lines},
			}},
		}}}, false},
		{"withNodes", args{&Document{Nodes: []Node{
			{Extras: 8.0, Extensions: Extensions{"a": "b"}, Name: "n-1", Camera: 1, Children: []uint32{1, 2}, Skin: 3,
				Matrix: [16]float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}, Mesh: 15, Rotation: [4]float64{1.5, 1.3, 12, 0}, Scale: [3]float64{1, 3, 4}, Translation: [3]float64{0, 7.8, 9}, Weights: []float64{1, 3}},
			{Extras: 8.0, Extensions: Extensions{"a": "b"}, Name: "n-2", Camera: 1, Children: []uint32{1, 2}, Skin: 3,
				Matrix: [16]float64{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1}, Mesh: 15, Rotation: [4]float64{0, 0, 0, 1}, Scale: [3]float64{1, 1, 1}},
		}}}, false},
		{"withSampler", args{&Document{Samplers: []Sampler{
			{Extras: 8.0, Extensions: Extensions{"a": "b"}, Name: "s_1", MagFilter: MagLinear, MinFilter: MinLinearMipMapLinear, WrapS: ClampToEdge, WrapT: MirroredRepeat},
			{Extras: 8.0, Extensions: Extensions{"a": "b"}, Name: "s_2", MagFilter: MagNearest, MinFilter: MinLinearMipMapLinear, WrapS: MirroredRepeat, WrapT: Repeat},
		}}}, false},
		{"withScene", args{&Document{Scene: 1}}, false},
		{"withScenes", args{&Document{Scenes: []Scene{
			{Extras: 8.0, Extensions: Extensions{"a": "b"}, Name: "s_1", Nodes: []uint32{1, 2}},
			{Extras: 8.0, Extensions: Extensions{"a": "b"}, Name: "s_2", Nodes: []uint32{2, 3}},
		}}}, false},
		{"withSkins", args{&Document{Skins: []Skin{
			{Extras: 8.0, Extensions: Extensions{"a": "b"}, Name: "skin_1", InverseBindMatrices: 2, Skeleton: 4, Joints: []uint32{5, 6}},
			{Extras: 8.0, Extensions: Extensions{"a": "b"}, Name: "skin_2", InverseBindMatrices: 3, Skeleton: 4, Joints: []uint32{7, 8}},
		}}}, false},
		{"withTextures", args{&Document{Textures: []Texture{
			{Extras: 8.0, Extensions: Extensions{"a": "b"}, Name: "t_1", Sampler: 2, Source: 3},
			{Extras: 8.0, Extensions: Extensions{"a": "b"}, Name: "t_2", Sampler: 3, Source: 4},
		}}}, false},
	}
	for _, tt := range tests {
		for _, method := range []string{"json", "binary"} {
			t.Run(fmt.Sprintf("%s_%s", tt.name, method), func(t *testing.T) {
				var asBinary bool
				if method == "binary" && !tt.wantErr {
					asBinary = true
					for i := 1; i < len(tt.args.doc.Buffers); i++ {
						tt.args.doc.Buffers[i].EmbeddedResource()
					}
				}
				d, err := saveMemory(tt.args.doc, asBinary)
				if (err != nil) != tt.wantErr {
					t.Errorf("Encoder.Encode() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !tt.wantErr {
					doc := new(Document)
					d.Decode(doc)
					if diff := deep.Equal(doc, tt.args.doc); diff != nil {
						t.Errorf("Encoder.Encode() = %v", diff)
						return
					}
				}
			})
		}
	}
}
