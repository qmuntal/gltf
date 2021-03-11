package gltf

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/go-test/deep"
)

func readFile(path string) []byte {
	r, _ := ioutil.ReadFile(path)
	return r
}

func TestOpen(t *testing.T) {
	deep.FloatPrecision = 5
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
			Accessors: []*Accessor{
				{BufferView: Index(0), ByteOffset: 0, ComponentType: ComponentUshort, Count: 36, Max: []float32{35}, Min: []float32{0}, Type: AccessorScalar},
				{BufferView: Index(1), ByteOffset: 0, ComponentType: ComponentFloat, Count: 36, Max: []float32{1, 1, 1}, Min: []float32{-1, -1, -1}, Type: AccessorVec3},
				{BufferView: Index(2), ByteOffset: 0, ComponentType: ComponentFloat, Count: 36, Max: []float32{1, 1, 1}, Min: []float32{-1, -1, -1}, Type: AccessorVec3},
				{BufferView: Index(3), ByteOffset: 0, ComponentType: ComponentFloat, Count: 36, Max: []float32{1, 0, 0, 1}, Min: []float32{0, 0, -1, -1}, Type: AccessorVec4},
				{BufferView: Index(4), ByteOffset: 0, ComponentType: ComponentFloat, Count: 36, Max: []float32{1, 1}, Min: []float32{-1, -1}, Type: AccessorVec2}},
			Asset: Asset{Generator: "VKTS glTF 2.0 exporter", Version: "2.0"},
			BufferViews: []*BufferView{
				{Buffer: 0, ByteLength: 72, ByteOffset: 0, Target: TargetElementArrayBuffer},
				{Buffer: 0, ByteLength: 432, ByteOffset: 72, Target: TargetArrayBuffer},
				{Buffer: 0, ByteLength: 432, ByteOffset: 504, Target: TargetArrayBuffer},
				{Buffer: 0, ByteLength: 576, ByteOffset: 936, Target: TargetArrayBuffer},
				{Buffer: 0, ByteLength: 288, ByteOffset: 1512, Target: TargetArrayBuffer},
			},
			Buffers:   []*Buffer{{ByteLength: 1800, URI: "Cube.bin", Data: readFile("testdata/Cube/glTF/Cube.bin")}},
			Images:    []*Image{{URI: "Cube_BaseColor.png"}, {URI: "Cube_MetallicRoughness.png"}},
			Materials: []*Material{{Name: "Cube", AlphaMode: AlphaOpaque, AlphaCutoff: Float(0.5), PBRMetallicRoughness: &PBRMetallicRoughness{BaseColorFactor: &[4]float32{1, 1, 1, 1}, MetallicFactor: Float(1), RoughnessFactor: Float(1), BaseColorTexture: &TextureInfo{Index: 0}, MetallicRoughnessTexture: &TextureInfo{Index: 1}}}},
			Meshes:    []*Mesh{{Name: "Cube", Primitives: []*Primitive{{Indices: Index(0), Material: Index(0), Mode: PrimitiveTriangles, Attributes: map[string]uint32{NORMAL: 2, POSITION: 1, TANGENT: 3, TEXCOORD_0: 4}}}}},
			Nodes:     []*Node{{Mesh: Index(0), Name: "Cube", Matrix: [16]float32{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1}, Rotation: [4]float32{0, 0, 0, 1}, Scale: [3]float32{1, 1, 1}}},
			Samplers:  []*Sampler{{WrapS: WrapRepeat, WrapT: WrapRepeat}},
			Scene:     Index(0),
			Scenes:    []*Scene{{Nodes: []uint32{0}}},
			Textures: []*Texture{
				{Sampler: Index(0), Source: Index(0)}, {Sampler: Index(0), Source: Index(1)},
			},
		}, false},
		{args{"testdata/Cameras/glTF/Cameras.gltf", "testdata/Cameras/glTF-Embedded/Cameras.gltf"}, &Document{
			Accessors: []*Accessor{
				{BufferView: Index(0), ByteOffset: 0, ComponentType: ComponentUshort, Count: 6, Max: []float32{3}, Min: []float32{0}, Type: AccessorScalar},
				{BufferView: Index(1), ByteOffset: 0, ComponentType: ComponentFloat, Count: 4, Max: []float32{1, 1, 0}, Min: []float32{0, 0, 0}, Type: AccessorVec3},
			},
			Asset: Asset{Version: "2.0"},
			BufferViews: []*BufferView{
				{Buffer: 0, ByteLength: 12, ByteOffset: 0, Target: TargetElementArrayBuffer},
				{Buffer: 0, ByteLength: 48, ByteOffset: 12, Target: TargetArrayBuffer},
			},
			Buffers: []*Buffer{{ByteLength: 60, URI: "simpleSquare.bin", Data: readFile("testdata/Cameras/glTF/simpleSquare.bin")}},
			Cameras: []*Camera{
				{Perspective: &Perspective{AspectRatio: Float(1.0), Yfov: 0.7, Zfar: Float(100), Znear: 0.01}},
				{Orthographic: &Orthographic{Xmag: 1.0, Ymag: 1.0, Zfar: 100, Znear: 0.01}},
			},
			Meshes: []*Mesh{{Primitives: []*Primitive{{Indices: Index(0), Mode: PrimitiveTriangles, Attributes: map[string]uint32{POSITION: 1}}}}},
			Nodes: []*Node{
				{Mesh: Index(0), Matrix: [16]float32{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1}, Rotation: [4]float32{-0.3, 0, 0, 0.9}, Scale: [3]float32{1, 1, 1}},
				{Camera: Index(0), Matrix: [16]float32{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1}, Rotation: [4]float32{0, 0, 0, 1}, Scale: [3]float32{1, 1, 1}, Translation: [3]float32{0.5, 0.5, 3.0}},
				{Camera: Index(1), Matrix: [16]float32{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1}, Rotation: [4]float32{0, 0, 0, 1}, Scale: [3]float32{1, 1, 1}, Translation: [3]float32{0.5, 0.5, 3.0}},
			},
			Scene:  nil,
			Scenes: []*Scene{{Nodes: []uint32{0, 1, 2}}},
		}, false},
		{args{"testdata/BoxVertexColors/glTF-Binary/BoxVertexColors.glb", ""}, &Document{
			Accessors: []*Accessor{
				{BufferView: Index(0), ByteOffset: 0, ComponentType: ComponentUshort, Count: 36, Type: AccessorScalar},
				{BufferView: Index(1), ByteOffset: 0, ComponentType: ComponentFloat, Count: 24, Max: []float32{0.5, 0.5, 0.5}, Min: []float32{-0.5, -0.5, -0.5}, Type: AccessorVec3},
				{BufferView: Index(2), ByteOffset: 0, ComponentType: ComponentFloat, Count: 24, Type: AccessorVec3},
				{BufferView: Index(3), ByteOffset: 0, ComponentType: ComponentFloat, Count: 24, Type: AccessorVec4},
				{BufferView: Index(4), ByteOffset: 0, ComponentType: ComponentFloat, Count: 24, Type: AccessorVec2},
			},
			Asset: Asset{Version: "2.0", Generator: "FBX2glTF"},
			BufferViews: []*BufferView{
				{Buffer: 0, ByteLength: 72, ByteOffset: 0, Target: TargetElementArrayBuffer},
				{Buffer: 0, ByteLength: 288, ByteOffset: 72, Target: TargetArrayBuffer},
				{Buffer: 0, ByteLength: 288, ByteOffset: 360, Target: TargetArrayBuffer},
				{Buffer: 0, ByteLength: 384, ByteOffset: 648, Target: TargetArrayBuffer},
				{Buffer: 0, ByteLength: 192, ByteOffset: 1032, Target: TargetArrayBuffer},
			},
			Buffers:   []*Buffer{{ByteLength: 1224, Data: readFile("testdata/BoxVertexColors/glTF-Binary/BoxVertexColors.glb")[1628+20+8:]}},
			Materials: []*Material{{Name: "Default", AlphaMode: AlphaOpaque, AlphaCutoff: Float(0.5), PBRMetallicRoughness: &PBRMetallicRoughness{BaseColorFactor: &[4]float32{0.8, 0.8, 0.8, 1}, MetallicFactor: Float(0.1), RoughnessFactor: Float(0.99)}}},
			Meshes:    []*Mesh{{Name: "Cube", Primitives: []*Primitive{{Indices: Index(0), Material: Index(0), Mode: PrimitiveTriangles, Attributes: map[string]uint32{POSITION: 1, COLOR_0: 3, NORMAL: 2, TEXCOORD_0: 4}}}}},
			Nodes: []*Node{
				{Name: "RootNode", Children: []uint32{1, 2, 3}, Matrix: [16]float32{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1}, Rotation: [4]float32{0, 0, 0, 1}, Scale: [3]float32{1, 1, 1}},
				{Name: "Mesh", Matrix: [16]float32{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1}, Rotation: [4]float32{0, 0, 0, 1}, Scale: [3]float32{1, 1, 1}},
				{Name: "Cube", Mesh: Index(0), Matrix: [16]float32{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1}, Rotation: [4]float32{0, 0, 0, 1}, Scale: [3]float32{1, 1, 1}},
				{Name: "Texture Group", Matrix: [16]float32{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1}, Rotation: [4]float32{0, 0, 0, 1}, Scale: [3]float32{1, 1, 1}},
			},
			Samplers: []*Sampler{{WrapS: WrapRepeat, WrapT: WrapRepeat}},
			Scene:    Index(0),
			Scenes:   []*Scene{{Name: "Root Scene", Nodes: []uint32{0}}},
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

type chunkedReader struct {
	s []byte
	n int
}

func (c *chunkedReader) Read(p []byte) (n int, err error) {
	c.n++
	if c.n == len(c.s)+1 {
		return 0, io.EOF
	}
	p[0] = c.s[c.n-1 : c.n][0]
	return 1, nil
}

type mockReadHandler struct {
	Payload string
}

func (m mockReadHandler) ReadFullResource(uri string, data []byte) error {
	copy(data, []byte(m.Payload))
	return nil
}

func TestDecoder_decodeBuffer(t *testing.T) {
	type args struct {
		buffer *Buffer
	}
	tests := []struct {
		name    string
		d       *Decoder
		args    args
		want    []byte
		wantErr bool
	}{
		{"byteLength_0", &Decoder{}, args{&Buffer{ByteLength: 0, URI: "a.bin"}}, nil, true},
		{"noURI", &Decoder{}, args{&Buffer{ByteLength: 1, URI: ""}}, nil, true},
		{"invalidURI", &Decoder{}, args{&Buffer{ByteLength: 1, URI: "../a.bin"}}, nil, true},
		{"noSchemeErr", NewDecoder(nil), args{&Buffer{ByteLength: 3, URI: "ftp://a.bin"}}, nil, true},
		{"base", NewDecoder(nil).WithReadHandler(&mockReadHandler{"abcdfg"}), args{&Buffer{ByteLength: 6, URI: "a.bin"}}, []byte("abcdfg"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.d.decodeBuffer(tt.args.buffer); (err != nil) != tt.wantErr {
				t.Errorf("Decoder.decodeBuffer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(tt.args.buffer.Data, tt.want) {
				t.Errorf("Decoder.decodeBuffer() buffer = %v, want %v", tt.args.buffer.Data, tt.want)
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
		want    []byte
		wantErr bool
	}{
		{"base", NewDecoder(bytes.NewBuffer([]byte{0x06, 0x00, 0x00, 0x00, 0x42, 0x49, 0x4e, 0x00, 1, 2, 3, 4, 5, 6})),
			args{&Buffer{ByteLength: 6}}, []byte{1, 2, 3, 4, 5, 6}, false},
		{"smallbuffer", NewDecoder(bytes.NewBuffer([]byte{0x6, 0x00, 0x00, 0x00, 0x42, 0x49, 0x4e, 0x00, 1, 2, 3, 4, 5, 6})),
			args{&Buffer{ByteLength: 5}}, []byte{1, 2, 3, 4, 5}, false},
		{"bigbuffer", NewDecoder(bytes.NewBuffer([]byte{0x6, 0x00, 0x00, 0x00, 0x42, 0x49, 0x4e, 0x00, 1, 2, 3, 4, 5, 6})),
			args{&Buffer{ByteLength: 7}}, nil, true},
		{"invalidBuffer", new(Decoder), args{&Buffer{ByteLength: 0}}, nil, true},
		{"readErr", NewDecoder(bytes.NewBufferString("")), args{&Buffer{ByteLength: 1}}, nil, true},
		{"invalidHeader", NewDecoder(bytes.NewBufferString("aaaaaaaa")), args{&Buffer{ByteLength: 1}}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.d.decodeBinaryBuffer(tt.args.buffer); (err != nil) != tt.wantErr {
				t.Errorf("Decoder.decodeBinaryBuffer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(tt.args.buffer.Data, tt.want) {
				t.Errorf("Decoder.decodeBinaryBuffer() buffer = %v, want %v", tt.args.buffer.Data, tt.want)
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
		{"baseJSON", NewDecoder(bytes.NewBufferString("{\"buffers\": [{\"byteLength\": 1, \"URI\": \"a.bin\"}]}")).WithReadHandler(&mockReadHandler{"abcdfg"}), args{new(Document)}, false},
		{"onlyGLBHeader", NewDecoder(bytes.NewBuffer([]byte{0x67, 0x6c, 0x54, 0x46, 0x02, 0x00, 0x00, 0x00, 0x40, 0x0b, 0x00, 0x00, 0x5c, 0x06, 0x00, 0x00, 0x4a, 0x53, 0x4f, 0x4e})).WithReadHandler(&mockReadHandler{"abcdfg"}), args{new(Document)}, true},
		{"glbNoJSONChunk", NewDecoder(bytes.NewBuffer([]byte{0x67, 0x6c, 0x54, 0x46, 0x02, 0x00, 0x00, 0x00, 0x40, 0x0b, 0x00, 0x00, 0x5c, 0x06, 0x00, 0x00, 0x4a, 0x52, 0x4f, 0x4e})).WithReadHandler(&mockReadHandler{"abcdfg"}), args{new(Document)}, true},
		{"empty", NewDecoder(bytes.NewBufferString("")), args{new(Document)}, true},
		{"invalidJSON", NewDecoder(bytes.NewBufferString("{asset: {}}")), args{new(Document)}, true},
		{"invalidBuffer", NewDecoder(bytes.NewBufferString("{\"buffers\": [{\"byteLength\": 0}]}")), args{new(Document)}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.d.Decode(tt.args.doc); (err != nil) != tt.wantErr {
				t.Errorf("Decoder.Decode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDecoder_validateDocumentQuotas(t *testing.T) {
	type args struct {
		doc      *Document
		isBinary bool
	}
	tests := []struct {
		name    string
		d       *Decoder
		args    args
		wantErr bool
	}{
		{
			"exceedBuffers", &Decoder{MaxMemoryAllocation: 100000, MaxExternalBufferCount: 1},
			args{&Document{Buffers: []*Buffer{{}, {}}}, false}, true,
		}, {
			"noExceedBuffers", &Decoder{MaxMemoryAllocation: 100000, MaxExternalBufferCount: 1},
			args{&Document{Buffers: []*Buffer{{}, {}}}, true}, false,
		}, {
			"exceedAllocs", &Decoder{MaxMemoryAllocation: 10, MaxExternalBufferCount: 100},
			args{&Document{Buffers: []*Buffer{{ByteLength: 11}}}, false}, true,
		}, {
			"noExceedAllocs", &Decoder{MaxMemoryAllocation: 11, MaxExternalBufferCount: 100},
			args{&Document{Buffers: []*Buffer{{ByteLength: 11}}}, true}, false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.d.validateDocumentQuotas(tt.args.doc, tt.args.isBinary); (err != nil) != tt.wantErr {
				t.Errorf("Decoder.validateDocumentQuotas() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSampler_Decode(t *testing.T) {

	tests := []struct {
		name    string
		s       []byte
		want    *Sampler
		wantErr bool
	}{
		{"empty", []byte(`{}`), &Sampler{}, false},
		{"nondefault",
			[]byte(`{"minFilter":9728,"wrapT":33071}`),
			&Sampler{MagFilter: MagLinear, MinFilter: MinNearest, WrapS: WrapRepeat, WrapT: WrapClampToEdge},
			false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got Sampler
			err := json.Unmarshal(tt.s, &got)
			if (err != nil) != tt.wantErr {
				t.Errorf("Unmarshaling Sampler error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(&got, tt.want) {
				t.Errorf("Unmarshaling Sampler = %v, want %v", string(tt.s), tt.want)
			}
		})
	}
}
