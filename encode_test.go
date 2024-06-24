package gltf

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/go-test/deep"
)

type mockFile struct {
	data *[]byte
}

func (f mockFile) Write(p []byte) (n int, err error) {
	*f.data = append(*f.data, p...)
	return len(p), nil
}

func (f mockFile) Close() error {
	return nil
}

type mockChunkReadHandler struct {
	fstest.MapFS
}

func (m mockChunkReadHandler) Create(uri string) (io.WriteCloser, error) {
	m.MapFS[uri] = &fstest.MapFile{}
	return mockFile{&m.MapFS[uri].Data}, nil
}

func saveMemory(doc *Document, asBinary bool) (*Decoder, error) {
	buff := new(bytes.Buffer)
	m := mockChunkReadHandler{fstest.MapFS{}}
	e := NewEncoderFS(buff, m)
	e.AsBinary = asBinary
	if err := e.Encode(doc); err != nil {
		return nil, err
	}

	return NewDecoderFS(buff, m), nil
}

func saveMemoryIndent(doc *Document, asBinary bool) (*Decoder, error) {
	buff := new(bytes.Buffer)
	m := mockChunkReadHandler{fstest.MapFS{}}
	e := NewEncoderFS(buff, m)
	e.AsBinary = asBinary
	e.SetJSONIndent("", "\t")
	if err := e.Encode(doc); err != nil {
		return nil, err
	}

	return NewDecoderFS(buff, m), nil
}

func TestEncoder_Encode_AsBinary_WithoutBuffer(t *testing.T) {
	doc := &Document{}
	buff := new(bytes.Buffer)
	e := NewEncoder(buff)
	e.AsBinary = true
	if err := e.Encode(doc); err != nil {
		t.Errorf("Encoder.Encode() error = %v", err)
	}
	if strings.Contains(buff.String(), "BIN") {
		t.Error("Encoder.Encode() as binary without bin buffer should not contain bin chunk")
	}
}

func TestEncoder_Encode_AsBinary_WithoutBinChunk(t *testing.T) {
	doc := &Document{Buffers: []*Buffer{
		{Extras: 8.0, Name: "embedded", ByteLength: 2, URI: "data:application/octet-stream;base64,YW55ICsgb2xkICYgZGF0YQ==", Data: []byte("any + old & data")},
		{Extras: 8.0, Name: "external", ByteLength: 4, URI: "b.bin", Data: []byte{4, 5, 6, 7}},
		{Extras: 8.0, Name: "external", ByteLength: 4, URI: "a.drc", Data: []byte{0, 0, 0, 0}},
	}}
	buff := new(bytes.Buffer)
	m := mockChunkReadHandler{fstest.MapFS{}}
	e := NewEncoderFS(buff, m)
	e.AsBinary = true
	if err := e.Encode(doc); err != nil {
		t.Errorf("Encoder.Encode() error = %v", err)
	}
	if strings.Contains(buff.String(), "BIN") {
		t.Error("Encoder.Encode() as binary without bin buffer should not contain bin chunk")
	}
}

func TestEncoder_Encode_AsBinary_WithBinChunk(t *testing.T) {
	doc := &Document{Buffers: []*Buffer{
		{Extras: 8.0, Name: "binary", ByteLength: 3, Data: []byte{1, 2, 3}},
	}}
	buff := new(bytes.Buffer)
	e := NewEncoder(buff)
	e.AsBinary = true
	if err := e.Encode(doc); err != nil {
		t.Errorf("Encoder.Encode() error = %v", err)
	}
	if !strings.Contains(buff.String(), "BIN") {
		t.Error("Encoder.Encode() as binary with bin buffer should contain bin chunk")
	}
	var header glbHeader
	if err := binary.Read(buff, binary.LittleEndian, &header); err != nil {
		t.Fatal(err)
	}
	if got := header.Length; got != 116 {
		t.Errorf("Encoder.Encode() incorrect length. want = %v, got = %v", 116, got)
	}
}

func TestEncoder_Encode_Buffers_Without_URI(t *testing.T) {
	doc := &Document{Buffers: []*Buffer{
		{Name: "binary", ByteLength: 3, Data: []byte{1, 2, 3}},
		{Name: "binary2", ByteLength: 3, Data: []byte{4, 5, 6}},
	}}
	buf := new(bytes.Buffer)
	e := NewEncoder(buf)
	e.AsBinary = false
	if err := e.Encode(doc); err != nil {
		t.Errorf("Encoder.Encode() error = %v", err)
	}
	if !strings.Contains(buf.String(), mimetypeApplicationOctet+",AQID") ||
		!strings.Contains(buf.String(), mimetypeApplicationOctet+",BAUG") {
		t.Error("Encoder.Encode() should auto embed buffers without URI")
	}
	buf.Reset()
	e.AsBinary = true
	if err := e.Encode(doc); err != nil {
		t.Errorf("Encoder.Encode() error = %v", err)
	}
	if strings.Contains(buf.String(), mimetypeApplicationOctet+",AQID") {
		t.Error("Encoder.Encode() as binary should not embed fur buffer")
	}
	if !strings.Contains(buf.String(), mimetypeApplicationOctet+",BAUG") {
		t.Error("Encoder.Encode() should auto embed buffers without URI")
	}
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
			{Extras: 8.0, Name: "an_1", Channels: []*AnimationChannel{
				{Extras: 8.0, Sampler: 1, Target: AnimationChannelTarget{Extras: 8.0, Node: Index(10), Path: TRSRotation}},
				{Extras: 8.0, Sampler: 2, Target: AnimationChannelTarget{Extras: 8.0, Node: Index(10), Path: TRSScale}},
			}},
			{Extras: 8.0, Name: "an_2", Channels: []*AnimationChannel{
				{Extras: 8.0, Sampler: 1, Target: AnimationChannelTarget{Extras: 8.0, Node: Index(3), Path: TRSWeights}},
				{Extras: 8.0, Sampler: 2, Target: AnimationChannelTarget{Extras: 8.0, Node: Index(5), Path: TRSTranslation}},
			}},
			{Extras: 8.0, Name: "an_3", Samplers: []*AnimationSampler{
				{Extras: 8.0, Input: 1, Output: 1, Interpolation: InterpolationCubicSpline},
			}},
		}}}, false},
		{"withBufView", args{&Document{BufferViews: []*BufferView{
			{Extras: 8.0, Buffer: 0, ByteOffset: 1, ByteLength: 2, ByteStride: 5, Target: TargetArrayBuffer},
			{Buffer: 10, ByteOffset: 10, ByteLength: 20, ByteStride: 50, Target: TargetElementArrayBuffer},
		}}}, false},
		{"withCameras", args{&Document{Cameras: []*Camera{
			{Extras: 8.0, Name: "cam_1", Orthographic: &Orthographic{Extras: 8.0, Xmag: 1, Ymag: 2, Zfar: 3, Znear: 4}},
			{Extras: 8.0, Name: "cam_2", Perspective: &Perspective{Extras: 8.0, AspectRatio: Float(1), Yfov: 2, Zfar: Float(3), Znear: 4}},
		}}}, false},
		{"withImages", args{&Document{Images: []*Image{
			{Extras: 8.0, Name: "binary", BufferView: Index(1), MimeType: "data:image/png"},
			{Extras: 8.0, Name: "embedded", URI: "data:image/png;base64,dsjdsaGGUDXGA", MimeType: "data:image/png"},
			{Extras: 8.0, Name: "external", URI: "https://web.com/a", MimeType: "data:image/png"},
		}}}, false},
		{"withMaterials", args{&Document{Materials: []*Material{
			{Extras: 8.0, Name: "base", EmissiveFactor: [3]float64{1.0, 1.0, 1.0}, DoubleSided: true, AlphaCutoff: Float(0.5), AlphaMode: AlphaOpaque},
			{Extras: 8.0, Name: "pbr", AlphaCutoff: Float(0.5), AlphaMode: AlphaOpaque,
				PBRMetallicRoughness: &PBRMetallicRoughness{
					Extras: 8.0, MetallicFactor: Float(1), RoughnessFactor: Float(2), BaseColorFactor: &[4]float64{0.8, 0.8, 0.8, 1},
					BaseColorTexture:         &TextureInfo{Extras: 8.0, Index: 1, TexCoord: 3},
					MetallicRoughnessTexture: &TextureInfo{Extras: 8.0, Index: 6, TexCoord: 5},
				},
			},
			{Extras: 8.0, Name: "normal", AlphaCutoff: Float(0.7), AlphaMode: AlphaBlend,
				NormalTexture: &NormalTexture{Extras: 8.0, Index: Index(1), TexCoord: 2, Scale: Float(2.0)},
			},
			{Extras: 8.0, Name: "occlusion", AlphaCutoff: Float(0.5), AlphaMode: AlphaMask,
				OcclusionTexture: &OcclusionTexture{Extras: 8.0, Index: Index(1), TexCoord: 2, Strength: Float(2.0)},
			},
			{Extras: 8.0, Name: "emmisice", AlphaCutoff: Float(0.5), AlphaMode: AlphaMask, EmissiveTexture: &TextureInfo{Extras: 8.0, Index: 4, TexCoord: 50}},
		}}}, false},
		{"withMeshes", args{&Document{Meshes: []*Mesh{
			{Extras: 8.0, Name: "mesh_1", Weights: []float64{1.2, 2}},
			{Extras: 8.0, Name: "mesh_2", Primitives: []*Primitive{
				{Extras: 8.0, Attributes: Attribute{POSITION: 1}, Indices: Index(2), Material: Index(1), Mode: PrimitiveLines},
				{Extras: 8.0, Targets: []Attribute{{POSITION: 1, "THEN": 4}, {"OTHER": 2}}, Indices: Index(2), Material: Index(1), Mode: PrimitiveLines},
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
		tt.args.doc.Asset.Version = "2.0"
		for _, method := range []string{"json", "binary"} {
			t.Run(fmt.Sprintf("%s_%s", tt.name, method), func(t *testing.T) {
				d, err := saveMemory(tt.args.doc, method == "binary")
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

				d, err = saveMemoryIndent(tt.args.doc, method == "binary")
				if (err != nil) != tt.wantErr {
					t.Errorf("Encoder.Encode() withIndent error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !tt.wantErr {
					doc := new(Document)
					d.Decode(doc)
					if diff := deep.Equal(doc, tt.args.doc); diff != nil {
						t.Errorf("Encoder.Encode() withIndent = %v", diff)
						return
					}
				}
			})
		}
	}
}

func TestImage_MarshalData(t *testing.T) {
	tests := []struct {
		name    string
		im      *Image
		want    []byte
		wantErr bool
	}{
		{"error", &Image{URI: "data:image/png;base64,_"}, []byte{}, true},
		{"external", &Image{URI: "http://web.com"}, []byte{}, false},
		{"empty", &Image{URI: "data:image/png;base64,"}, []byte{}, false},
		{"empty", &Image{URI: "data:image/jpeg;base64,"}, []byte{}, false},
		{"test", &Image{URI: "data:image/png;base64,TEST"}, []byte{76, 68, 147}, false},
		{"complex", &Image{URI: "data:image/png;base64,YW55IGNhcm5hbCBwbGVhcw=="}, []byte{97, 110, 121, 32, 99, 97, 114, 110, 97, 108, 32, 112, 108, 101, 97, 115}, false},
		{"invalid", &Image{URI: "data:image/png;base64"}, []byte{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.im.MarshalData()
			if (err != nil) != tt.wantErr {
				t.Errorf("Image.MarshalData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Image.MarshalData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNode_UnmarshalJSON(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		n       *Node
		args    args
		want    *Node
		wantErr bool
	}{
		{"default", new(Node), args{[]byte("{}")}, &Node{
			Matrix:   [16]float64{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1},
			Rotation: [4]float64{0, 0, 0, 1},
			Scale:    [3]float64{1, 1, 1},
		}, false},
		{"nodefault", new(Node), args{[]byte(`{"matrix":[1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1],"rotation":[0,0,0,1],"scale":[1,1,1],"camera":0,"mesh":0,"skin":0}`)}, &Node{
			Matrix:   [16]float64{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1},
			Rotation: [4]float64{0, 0, 0, 1},
			Scale:    [3]float64{1, 1, 1},
			Camera:   Index(0),
			Mesh:     Index(0),
			Skin:     Index(0),
		}, false},
		{"nodefault", new(Node), args{[]byte(`{"matrix":[1,2,2,0,0,1,3,4,0,0,1,0,5,0,0,5],"rotation":[1,2,3,4],"scale":[2,3,4],"camera":1,"mesh":2,"skin":3}`)}, &Node{
			Matrix:   [16]float64{1, 2, 2, 0, 0, 1, 3, 4, 0, 0, 1, 0, 5, 0, 0, 5},
			Rotation: [4]float64{1, 2, 3, 4},
			Scale:    [3]float64{2, 3, 4},
			Camera:   Index(1),
			Mesh:     Index(2),
			Skin:     Index(3),
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.n.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Node.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(tt.n, tt.want) {
				t.Errorf("Node.UnmarshalJSON() = %v, want %v", tt.n, tt.want)
			}
		})
	}
}

func TestMaterial_UnmarshalJSON(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		m       *Material
		args    args
		want    *Material
		wantErr bool
	}{
		{"default", new(Material), args{[]byte("{}")}, &Material{AlphaCutoff: Float(0.5), AlphaMode: AlphaOpaque}, false},
		{"nodefault", new(Material), args{[]byte(`{"alphaCutoff": 0.2, "alphaMode": "MASK"}`)}, &Material{AlphaCutoff: Float(0.2), AlphaMode: AlphaMask}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Material.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(tt.m, tt.want) {
				t.Errorf("Material.UnmarshalJSON() = %v, want %v", tt.m, tt.want)
			}
		})
	}
}

func TestNormalTexture_UnmarshalJSON(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		n       *NormalTexture
		args    args
		want    *NormalTexture
		wantErr bool
	}{
		{"default", new(NormalTexture), args{[]byte("{}")}, &NormalTexture{Scale: Float(1)}, false},
		{"empty", new(NormalTexture), args{[]byte(`{"scale": 0, "index": 0}`)}, &NormalTexture{Scale: Float(0), Index: Index(0)}, false},
		{"nodefault", new(NormalTexture), args{[]byte(`{"scale": 0.5, "index":2}`)}, &NormalTexture{Scale: Float(0.5), Index: Index(2)}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.n.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("NormalTexture.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(tt.n, tt.want) {
				t.Errorf("NormalTexture.UnmarshalJSON() = %v, want %v", tt.n, tt.want)
			}
		})
	}
}

func TestOcclusionTexture_UnmarshalJSON(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		o       *OcclusionTexture
		args    args
		want    *OcclusionTexture
		wantErr bool
	}{
		{"default", new(OcclusionTexture), args{[]byte("{}")}, &OcclusionTexture{Strength: Float(1)}, false},
		{"empty", new(OcclusionTexture), args{[]byte(`{"strength": 0, "index": 0}`)}, &OcclusionTexture{Strength: Float(0), Index: Index(0)}, false},
		{"nodefault", new(OcclusionTexture), args{[]byte(`{"strength": 0.5, "index":2}`)}, &OcclusionTexture{Strength: Float(0.5), Index: Index(2)}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.o.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("OcclusionTexture.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(tt.o, tt.want) {
				t.Errorf("OcclusionTexture.UnmarshalJSON() = %v, want %v", tt.o, tt.want)
			}
		})
	}
}

func TestPBRMetallicRoughness_UnmarshalJSON(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		p       *PBRMetallicRoughness
		args    args
		want    *PBRMetallicRoughness
		wantErr bool
	}{
		{"default", new(PBRMetallicRoughness), args{[]byte("{}")}, &PBRMetallicRoughness{BaseColorFactor: &[4]float64{1, 1, 1, 1}, MetallicFactor: Float(1), RoughnessFactor: Float(1)}, false},
		{"nodefault", new(PBRMetallicRoughness), args{[]byte(`{"baseColorFactor": [0.1,0.2,0.6,0.7],"metallicFactor":0.5,"roughnessFactor":0.1}`)}, &PBRMetallicRoughness{
			BaseColorFactor: &[4]float64{0.1, 0.2, 0.6, 0.7}, MetallicFactor: Float(0.5), RoughnessFactor: Float(0.1),
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.p.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("PBRMetallicRoughness.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(tt.p, tt.want) {
				t.Errorf("PBRMetallicRoughness.UnmarshalJSON() = %v, want %v", tt.p, tt.want)
			}
		})
	}
}

func TestNode_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		n       *Node
		want    []byte
		wantErr bool
	}{
		{"default", &Node{
			Matrix:   [16]float64{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1},
			Rotation: [4]float64{0, 0, 0, 1},
			Scale:    [3]float64{1, 1, 1},
		}, []byte("{}"), false},
		{"default2", &Node{
			Matrix:   [16]float64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			Rotation: [4]float64{0, 0, 0, 0},
			Scale:    [3]float64{0, 0, 0},
		}, []byte("{}"), false},
		{"empty", &Node{
			Matrix:   [16]float64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			Rotation: [4]float64{0, 0, 0, 0},
			Scale:    [3]float64{0, 0, 0},
			Camera:   Index(0),
			Skin:     Index(0),
			Mesh:     Index(0),
		}, []byte(`{"camera":0,"skin":0,"mesh":0}`), false},
		{"nodefault", &Node{
			Matrix:      [16]float64{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			Rotation:    [4]float64{1, 0, 0, 0},
			Scale:       [3]float64{1, 0, 0},
			Translation: [3]float64{1, 0, 0},
			Camera:      Index(1),
			Skin:        Index(1),
			Mesh:        Index(1),
		}, []byte(`{"matrix":[1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],"rotation":[1,0,0,0],"scale":[1,0,0],"translation":[1,0,0],"camera":1,"skin":1,"mesh":1}`), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.n.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("Node.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Node.MarshalJSON() = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}

func TestMaterial_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		m       *Material
		want    []byte
		wantErr bool
	}{
		{"default", &Material{AlphaCutoff: Float(0.5), AlphaMode: AlphaOpaque}, []byte(`{}`), false},
		{"empty", &Material{AlphaMode: AlphaBlend}, []byte(`{"alphaMode":"BLEND"}`), false},
		{"nodefault", &Material{AlphaCutoff: Float(1), AlphaMode: AlphaBlend}, []byte(`{"alphaCutoff":1,"alphaMode":"BLEND"}`), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("Material.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Material.MarshalJSON() = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}

func TestNormalTexture_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		n       *NormalTexture
		want    []byte
		wantErr bool
	}{
		{"default", &NormalTexture{Scale: Float(1)}, []byte(`{}`), false},
		{"empty", &NormalTexture{Index: Index(0)}, []byte(`{"index":0}`), false},
		{"nodefault", &NormalTexture{Index: Index(1), Scale: Float(0.5)}, []byte(`{"index":1,"scale":0.5}`), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.n.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("NormalTexture.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NormalTexture.MarshalJSON() = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}

func TestOcclusionTexture_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		o       *OcclusionTexture
		want    []byte
		wantErr bool
	}{
		{"default", &OcclusionTexture{Strength: Float(1)}, []byte(`{}`), false},
		{"empty", &OcclusionTexture{Index: Index(0)}, []byte(`{"index":0}`), false},
		{"nodefault", &OcclusionTexture{Index: Index(1), Strength: Float(0.5)}, []byte(`{"index":1,"strength":0.5}`), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.o.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("OcclusionTexture.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OcclusionTexture.MarshalJSON() = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}

func TestPBRMetallicRoughness_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		p       *PBRMetallicRoughness
		want    []byte
		wantErr bool
	}{
		{"default", &PBRMetallicRoughness{MetallicFactor: Float(1), RoughnessFactor: Float(1), BaseColorFactor: &[4]float64{1, 1, 1, 1}}, []byte(`{}`), false},
		{"empty", &PBRMetallicRoughness{MetallicFactor: Float(0), RoughnessFactor: Float(0)}, []byte(`{"metallicFactor":0,"roughnessFactor":0}`), false},
		{"nodefault", &PBRMetallicRoughness{MetallicFactor: Float(0.5), RoughnessFactor: Float(0.5), BaseColorFactor: &[4]float64{1, 0.5, 1, 1}}, []byte(`{"baseColorFactor":[1,0.5,1,1],"metallicFactor":0.5,"roughnessFactor":0.5}`), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.p.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("PBRMetallicRoughness.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PBRMetallicRoughness.MarshalJSON() = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}

func TestCamera_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		c       *Camera
		want    []byte
		wantErr bool
	}{
		{"empty", &Camera{}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("Camera.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Camera.MarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSampler_Encode(t *testing.T) {
	tests := []struct {
		name    string
		s       *Sampler
		want    []byte
		wantErr bool
	}{
		{"default", &Sampler{MagFilter: 0, MinFilter: 0, WrapS: 0, WrapT: 0}, []byte(`{}`), false},
		{"empty", &Sampler{}, []byte(`{}`), false},
		{"nondefault",
			&Sampler{MagFilter: MagLinear, MinFilter: MinNearest, WrapS: WrapRepeat, WrapT: WrapClampToEdge},
			[]byte(`{"magFilter":9729,"minFilter":9728,"wrapT":33071}`), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("Material.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Material.MarshalJSON() = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}

type fakeExt struct {
	A int `json:"a"`
}

func (f *fakeExt) UnmarshalJSON(data []byte) error {
	type alias fakeExt
	tmp := alias(fakeExt{})
	err := json.Unmarshal(data, &tmp)
	if err == nil {
		*f = fakeExt(tmp)
	}
	return err
}

func TestExtensions_UnmarshalJSON(t *testing.T) {
	RegisterExtension("fake_ext", func(data []byte) (any, error) {
		e := new(fakeExt)
		err := json.Unmarshal(data, e)
		return e, err
	})
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *Extensions
		wantErr bool
	}{
		{"fake_ext", args{[]byte(`{"fake_ext": {"a":2}}`)}, &Extensions{
			"fake_ext": &fakeExt{A: 2},
		}, false},
		{"err", args{[]byte(`{"fake_ext": {{"a":2}}`)}, &Extensions{}, true},
		{"errext", args{[]byte(`{"fake_ext": {"a":"incorrect"}}`)}, &Extensions{
			"fake_ext": json.RawMessage([]byte(`{"a":"incorrect"}`)),
		}, false},
		{"noregistered", args{[]byte(`{"fake_ext_1": {"a":"incorrect"}}`)}, &Extensions{
			"fake_ext_1": json.RawMessage([]byte(`{"a":"incorrect"}`)),
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ext := Extensions{}
			if err := ext.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Extensions.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(&ext, tt.want) {
				t.Errorf("Extensions.UnmarshalJSON() = %v, want %v", &ext, tt.want)
			}
		})
	}
}
