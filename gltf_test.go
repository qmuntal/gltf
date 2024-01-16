package gltf

import (
	"path/filepath"
	"reflect"
	"testing"
)

func TestDocument(t *testing.T) {
	tests := []struct {
		name string
		want *Document
	}{
		{"base", &Document{
			Scene:  Index(0),
			Scenes: []*Scene{{Name: "Root Scene"}},
			Asset: Asset{
				Generator: "qmuntal/gltf",
				Version:   "2.0",
			},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDocument(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDocument() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuffer_IsEmbeddedResource(t *testing.T) {
	tests := []struct {
		name string
		b    *Buffer
		want bool
	}{
		{"embedded", &Buffer{URI: "data:application/octet-stream;base64,dsjdsaGGUDXGA"}, true},
		{"external", &Buffer{URI: "https://web.com/a"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.IsEmbeddedResource(); got != tt.want {
				t.Errorf("Buffer.IsEmbeddedResource() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuffer_EmbeddedResource(t *testing.T) {
	tests := []struct {
		name string
		b    *Buffer
		want string
	}{
		{"base", &Buffer{Data: []byte("any + old & data")}, "data:application/octet-stream;base64,YW55ICsgb2xkICYgZGF0YQ=="},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.EmbeddedResource()
			if got := tt.b.URI; got != tt.want {
				t.Errorf("Buffer.EmbeddedResource() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestImage_IsEmbeddedResource(t *testing.T) {
	tests := []struct {
		name string
		im   *Image
		want bool
	}{
		{"png", &Image{URI: "data:image/png;base64,dsjdsaGGUDXGA"}, true},
		{"jpg", &Image{URI: "data:image/jpeg;base64,dsjdsaGGUDXGA"}, true},
		{"external", &Image{URI: "https://web.com/a"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.im.IsEmbeddedResource(); got != tt.want {
				t.Errorf("Image.IsEmbeddedResource() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuffer_marshalData(t *testing.T) {
	tests := []struct {
		name    string
		b       *Buffer
		want    []byte
		wantErr bool
	}{
		{"error", &Buffer{URI: "data:application/octet-stream;base64,_"}, nil, true},
		{"external", &Buffer{URI: "http://web.com"}, nil, false},
		{"empty", &Buffer{URI: "data:application/octet-stream;base64,"}, nil, false},
		{"test", &Buffer{URI: "data:application/octet-stream;base64,TEST"}, []byte{76, 68, 147}, false},
		{"complex", &Buffer{URI: "data:application/octet-stream;base64,YW55IGNhcm5hbCBwbGVhcw=="}, []byte{97, 110, 121, 32, 99, 97, 114, 110, 97, 108, 32, 112, 108, 101, 97, 115}, false},
		{"invalid", &Buffer{URI: "data:application/octet-stream;base64"}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.b.marshalData()
			if (err != nil) != tt.wantErr {
				t.Errorf("Buffer.marshalData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Buffer.marshalData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNode_MatrixOrDefault(t *testing.T) {
	tests := []struct {
		name string
		n    *Node
		want [16]float64
	}{
		{"default", &Node{Matrix: DefaultMatrix}, DefaultMatrix},
		{"zeros", &Node{Matrix: emptyMatrix}, DefaultMatrix},
		{"other", &Node{Matrix: [16]float64{2, 0, 0, 0, 0, 2, 0, 0, 0, 0, 2, 0, 0, 0, 0, 2}}, [16]float64{2, 0, 0, 0, 0, 2, 0, 0, 0, 0, 2, 0, 0, 0, 0, 2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.n.MatrixOrDefault(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Node.MatrixOrDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNode_RotationOrDefault(t *testing.T) {
	tests := []struct {
		name string
		n    *Node
		want [4]float64
	}{
		{"default", &Node{Rotation: DefaultRotation}, DefaultRotation},
		{"zeros", &Node{Rotation: emptyRotation}, DefaultRotation},
		{"other", &Node{Rotation: [4]float64{1, 2, 3, 4}}, [4]float64{1, 2, 3, 4}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.n.RotationOrDefault(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Node.RotationOrDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNode_ScaleOrDefault(t *testing.T) {
	tests := []struct {
		name string
		n    *Node
		want [3]float64
	}{
		{"default", &Node{Scale: DefaultScale}, DefaultScale},
		{"zeros", &Node{Scale: emptyScale}, DefaultScale},
		{"other", &Node{Scale: [3]float64{1, 2, 3}}, [3]float64{1, 2, 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.n.ScaleOrDefault(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Node.ScaleOrDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNode_TranslationOrDefault(t *testing.T) {
	tests := []struct {
		name string
		n    *Node
		want [3]float64
	}{
		{"default", &Node{Translation: DefaultTranslation}, DefaultTranslation},
		{"other", &Node{Translation: [3]float64{1, 2, 3}}, [3]float64{1, 2, 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.n.TranslationOrDefault(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Node.TranslationOrDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPBRMetallicRoughness_BaseColorFactorOrDefault(t *testing.T) {
	tests := []struct {
		name string
		p    *PBRMetallicRoughness
		want [4]float64
	}{
		{"empty", &PBRMetallicRoughness{}, [4]float64{1, 1, 1, 1}},
		{"other", &PBRMetallicRoughness{BaseColorFactor: &[4]float64{0.8, 0.8, 0.8, 0.5}}, [4]float64{0.8, 0.8, 0.8, 0.5}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.BaseColorFactorOrDefault(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PBRMetallicRoughness.BaseColorFactorOrDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOcclusionTexture_StrengthOrDefault(t *testing.T) {
	tests := []struct {
		name string
		o    *OcclusionTexture
		want float64
	}{
		{"empty", &OcclusionTexture{}, 1},
		{"other", &OcclusionTexture{Strength: Float(2)}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.o.StrengthOrDefault(); got != tt.want {
				t.Errorf("OcclusionTexture.StrengthOrDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNormalTexture_ScaleOrDefault(t *testing.T) {
	tests := []struct {
		name string
		n    *NormalTexture
		want float64
	}{
		{"empty", &NormalTexture{}, 1},
		{"other", &NormalTexture{Scale: Float(2)}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.n.ScaleOrDefault(); got != tt.want {
				t.Errorf("NormalTexture.ScaleOrDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMaterial_AlphaCutoffOrDefault(t *testing.T) {
	tests := []struct {
		name string
		m    *Material
		want float64
	}{
		{"empty", &Material{}, 0.5},
		{"other", &Material{AlphaCutoff: Float(2)}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.AlphaCutoffOrDefault(); got != tt.want {
				t.Errorf("Material.AlphaCutoffOrDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPBRMetallicRoughness_MetallicFactorOrDefault(t *testing.T) {
	tests := []struct {
		name string
		p    *PBRMetallicRoughness
		want float64
	}{
		{"empty", &PBRMetallicRoughness{}, 1},
		{"other", &PBRMetallicRoughness{MetallicFactor: Float(2)}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.MetallicFactorOrDefault(); got != tt.want {
				t.Errorf("PBRMetallicRoughness.MetallicFactorOrDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPBRMetallicRoughness_RoughnessFactorOrDefault(t *testing.T) {
	tests := []struct {
		name string
		p    *PBRMetallicRoughness
		want float64
	}{
		{"empty", &PBRMetallicRoughness{}, 1},
		{"other", &PBRMetallicRoughness{RoughnessFactor: Float(2)}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.RoughnessFactorOrDefault(); got != tt.want {
				t.Errorf("PBRMetallicRoughness.RoughnessFactorOrDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSizeOfElement(t *testing.T) {
	type args struct {
		c ComponentType
		t AccessorType
	}
	tests := []struct {
		name string
		args args
		want uint32
	}{
		{"byte-vec2", args{ComponentByte, AccessorVec2}, 4},
		{"ubyte-vec2", args{ComponentUbyte, AccessorVec2}, 4},
		{"byte-vec3", args{ComponentByte, AccessorVec3}, 4},
		{"ubyte-vec3", args{ComponentUbyte, AccessorVec3}, 4},
		{"short-vec3", args{ComponentShort, AccessorVec3}, 8},
		{"ushort-vec3", args{ComponentUshort, AccessorVec3}, 8},
		{"byte-mat2", args{ComponentByte, AccessorMat2}, 8},
		{"ubyte-mat2", args{ComponentUbyte, AccessorMat2}, 8},
		{"byte-mat3", args{ComponentByte, AccessorMat3}, 12},
		{"ubyte-mat3", args{ComponentUbyte, AccessorMat3}, 12},
		{"short-mat3", args{ComponentShort, AccessorMat3}, 24},
		{"ushort-mat3", args{ComponentUshort, AccessorMat3}, 24},
		{"other", args{ComponentUshort, AccessorMat4}, 32},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SizeOfElement(tt.args.c, tt.args.t); got != tt.want {
				t.Errorf("SizeOfElement() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIssue63(t *testing.T) {
	srcDoc, err := Open("testdata/issue63/simple.gltf")
	if nil != err {
		t.Fatal(err)
	}
	outputGlb := filepath.Join(t.TempDir(), "out.glb")
	err = SaveBinary(srcDoc, outputGlb)
	if nil != err {
		t.Fatal(err)
	}
	_, err = Open(outputGlb)
	if err != nil {
		t.Fatal(err)
	}
}
