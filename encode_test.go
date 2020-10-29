package gltf

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/go-test/deep"
)

type mockChunkReadHandler struct {
	Chunks map[string][]byte
}

func (m mockChunkReadHandler) ReadFullResource(uri string, data []byte) error {
	copy(data, m.Chunks[uri])
	return nil
}

func (m mockChunkReadHandler) WriteResource(uri string, data []byte) error {
	m.Chunks[uri] = data
	return nil
}

func saveMemory(doc *Document, asBinary bool) (*Decoder, error) {
	buff := new(bytes.Buffer)
	m := &mockChunkReadHandler{Chunks: make(map[string][]byte)}
	e := NewEncoder(buff).WithWriteHandler(m)
	e.AsBinary = asBinary
	if err := e.Encode(doc); err != nil {
		return nil, err
	}

	return NewDecoder(buff).WithReadHandler(m), nil
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
		{"withInvalidBuffer", args{&Document{Buffers: []*Buffer{
			{Extras: 8.0, Name: "binary", ByteLength: 3, URI: "a.bin", Data: []byte{1, 2, 3}},
			{Extras: 8.0, Name: "binary2", ByteLength: 3, URI: "/../a.bin", Data: []byte{1, 2, 3}},
		}}}, true},
		{"empty", args{&Document{}}, false},
		{"withExtensions", args{&Document{Extras: 8.0, ExtensionsUsed: []string{"c"}, ExtensionsRequired: []string{"d", "e"}}}, false},
		{"withAsset", args{&Document{Asset: Asset{Extras: 8.0, Copyright: "@2019", Generator: "qmuntal/gltf", Version: "2.0", MinVersion: "1.0"}}}, false},
		{"withAccessors", args{&Document{Accessors: []*Accessor{
			{Extras: 8.0, Name: "acc_1", BufferView: Index(0), ByteOffset: 50, ComponentType: ComponentByte, Normalized: true, Count: 5, Type: AccessorVec3, Max: []float64{1, 2}, Min: []float64{2.4}},
			{BufferView: Index(0), Normalized: false, Count: 50, Type: AccessorVec4, Sparse: &Sparse{Extras: 8.0, Count: 2,
				Values:  SparseValues{Extras: 8.0, BufferView: 1, ByteOffset: 2},
				Indices: SparseIndices{Extras: 8.0, BufferView: 1, ByteOffset: 2, ComponentType: ComponentFloat}},
			},
		}}}, false},
		{"withAnimations", args{&Document{Animations: []*Animation{
			{Extras: 8.0, Name: "an_1", Channels: []*Channel{
				{Extras: 8.0, Sampler: Index(1), Target: ChannelTarget{Extras: 8.0, Node: Index(10), Path: TRSRotation}},
				{Extras: 8.0, Sampler: Index(2), Target: ChannelTarget{Extras: 8.0, Node: Index(10), Path: TRSScale}},
			}},
			{Extras: 8.0, Name: "an_2", Channels: []*Channel{
				{Extras: 8.0, Sampler: Index(1), Target: ChannelTarget{Extras: 8.0, Node: Index(3), Path: TRSWeights}},
				{Extras: 8.0, Sampler: Index(2), Target: ChannelTarget{Extras: 8.0, Node: Index(5), Path: TRSTranslation}},
			}},
			{Extras: 8.0, Name: "an_3", Samplers: []*AnimationSampler{
				{Extras: 8.0, Input: Index(1), Output: Index(1), Interpolation: InterpolationCubicSpline},
			}},
		}}}, false},
		{"withBuffer", args{&Document{Buffers: []*Buffer{
			{Extras: 8.0, Name: "binary", ByteLength: 3, URI: "a.bin", Data: []byte{1, 2, 3}},
			{Extras: 8.0, Name: "embedded", ByteLength: 2, URI: "data:application/octet-stream;base64,YW55ICsgb2xkICYgZGF0YQ==", Data: []byte("any + old & data")},
			{Extras: 8.0, Name: "external", ByteLength: 4, URI: "b.bin", Data: []byte{4, 5, 6, 7}},
			{Extras: 8.0, Name: "external", ByteLength: 4, URI: "a.drc", Data: []byte{0, 0, 0, 0}},
		}}}, false},
		{"withBufView", args{&Document{BufferViews: []*BufferView{
			{Extras: 8.0, Buffer: 0, ByteOffset: 1, ByteLength: 2, ByteStride: 5, Target: TargetArrayBuffer},
			{Buffer: 10, ByteOffset: 10, ByteLength: 20, ByteStride: 50, Target: TargetElementArrayBuffer},
		}}}, false},
		{"withCameras", args{&Document{Cameras: []*Camera{
			{Extras: 8.0, Name: "cam_1", Orthographic: &Orthographic{Extras: 8.0, Xmag: 1, Ymag: 2, Zfar: 3, Znear: 4}},
			{Extras: 8.0, Name: "cam_2", Perspective: &Perspective{Extras: 8.0, AspectRatio: Float64(1), Yfov: 2, Zfar: Float64(3), Znear: 4}},
		}}}, false},
		{"withImages", args{&Document{Images: []*Image{
			{Extras: 8.0, Name: "binary", BufferView: Index(1), MimeType: "data:image/png"},
			{Extras: 8.0, Name: "embedded", URI: "data:image/png;base64,dsjdsaGGUDXGA", MimeType: "data:image/png"},
			{Extras: 8.0, Name: "external", URI: "https://web.com/a", MimeType: "data:image/png"},
		}}}, false},
		{"withMaterials", args{&Document{Materials: []*Material{
			{Extras: 8.0, Name: "base", EmissiveFactor: [3]float64{1.0, 1.0, 1.0}, DoubleSided: true, AlphaCutoff: Float64(0.5), AlphaMode: AlphaOpaque},
			{Extras: 8.0, Name: "pbr", AlphaCutoff: Float64(0.5), AlphaMode: AlphaOpaque,
				PBRMetallicRoughness: &PBRMetallicRoughness{
					Extras: 8.0, MetallicFactor: Float64(1), RoughnessFactor: Float64(2), BaseColorFactor: &[4]float32{0.8, 0.8, 0.8, 1},
					BaseColorTexture:         &TextureInfo{Extras: 8.0, Index: 1, TexCoord: 3},
					MetallicRoughnessTexture: &TextureInfo{Extras: 8.0, Index: 6, TexCoord: 5},
				},
			},
			{Extras: 8.0, Name: "normal", AlphaCutoff: Float64(0.7), AlphaMode: AlphaBlend,
				NormalTexture: &NormalTexture{Extras: 8.0, Index: Index(1), TexCoord: 2, Scale: Float64(2.0)},
			},
			{Extras: 8.0, Name: "occlusion", AlphaCutoff: Float64(0.5), AlphaMode: AlphaMask,
				OcclusionTexture: &OcclusionTexture{Extras: 8.0, Index: Index(1), TexCoord: 2, Strength: Float64(2.0)},
			},
			{Extras: 8.0, Name: "emmisice", AlphaCutoff: Float64(0.5), AlphaMode: AlphaMask, EmissiveTexture: &TextureInfo{Extras: 8.0, Index: 4, TexCoord: 50}},
		}}}, false},
		{"withMeshes", args{&Document{Meshes: []*Mesh{
			{Extras: 8.0, Name: "mesh_1", Weights: []float64{1.2, 2}},
			{Extras: 8.0, Name: "mesh_2", Primitives: []*Primitive{
				{Extras: 8.0, Attributes: Attribute{"POSITION": 1}, Indices: Index(2), Material: Index(1), Mode: PrimitiveLines},
				{Extras: 8.0, Targets: []Attribute{{"POSITION": 1, "THEN": 4}, {"OTHER": 2}}, Indices: Index(2), Material: Index(1), Mode: PrimitiveLines},
			}},
		}}}, false},
		{"withNodes", args{&Document{Nodes: []*Node{
			{Extras: 8.0, Name: "n-1", Camera: Index(1), Children: []uint32{1, 2}, Skin: Index(3),
				Matrix: [16]float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}, Mesh: Index(15), Rotation: [4]float64{1.5, 1.3, 12, 0}, Scale: [3]float64{1, 3, 4}, Translation: [3]float64{0, 7.8, 9}, Weights: []float64{1, 3}},
			{Extras: 8.0, Name: "n-2", Camera: Index(1), Children: []uint32{1, 2}, Skin: Index(3),
				Matrix: [16]float64{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1}, Mesh: Index(15), Rotation: [4]float64{0, 0, 0, 1}, Scale: [3]float64{1, 1, 1}},
		}}}, false},
		{"withSampler", args{&Document{Samplers: []*Sampler{
			{Extras: 8.0, Name: "s_1", MagFilter: MagLinear, MinFilter: MinLinearMipMapLinear, WrapS: WrapClampToEdge, WrapT: WrapMirroredRepeat},
			{Extras: 8.0, Name: "s_2", MagFilter: MagNearest, MinFilter: MinLinearMipMapLinear, WrapS: WrapMirroredRepeat, WrapT: WrapRepeat},
		}}}, false},
		{"withScene", args{&Document{Scene: Index(1)}}, false},
		{"withScenes", args{&Document{Scenes: []*Scene{
			{Extras: 8.0, Name: "s_1", Nodes: []uint32{1, 2}},
			{Extras: 8.0, Name: "s_2", Nodes: []uint32{2, 3}},
		}}}, false},
		{"withSkins", args{&Document{Skins: []*Skin{
			{Extras: 8.0, Name: "skin_1", InverseBindMatrices: Index(2), Skeleton: Index(4), Joints: []uint32{5, 6}},
			{Extras: 8.0, Name: "skin_2", InverseBindMatrices: Index(3), Skeleton: Index(4), Joints: []uint32{7, 8}},
		}}}, false},
		{"withTextures", args{&Document{Textures: []*Texture{
			{Extras: 8.0, Name: "t_1", Sampler: Index(2), Source: Index(3)},
			{Extras: 8.0, Name: "t_2", Sampler: Index(3), Source: Index(4)},
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
