package gltf

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestImage_MarshalData(t *testing.T) {
	tests := []struct {
		name    string
		im      *Image
		want    []uint8
		wantErr bool
	}{
		{"error", &Image{URI: "data:image/png;base64,_"}, []uint8{}, true},
		{"external", &Image{URI: "http://web.com"}, []uint8{}, false},
		{"empty", &Image{URI: "data:image/png;base64,"}, []uint8{}, false},
		{"empty", &Image{URI: "data:image/jpeg;base64,"}, []uint8{}, false},
		{"test", &Image{URI: "data:image/png;base64,TEST"}, []uint8{76, 68, 147}, false},
		{"complex", &Image{URI: "data:image/png;base64,YW55IGNhcm5hbCBwbGVhcw=="}, []uint8{97, 110, 121, 32, 99, 97, 114, 110, 97, 108, 32, 112, 108, 101, 97, 115}, false},
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
		{"default", new(Material), args{[]byte("{}")}, &Material{AlphaCutoff: Float64(0.5), AlphaMode: Opaque}, false},
		{"nodefault", new(Material), args{[]byte(`{"alphaCutoff": 0.2, "alphaMode": "MASK"}`)}, &Material{AlphaCutoff: Float64(0.2), AlphaMode: Mask}, false},
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
		{"default", new(NormalTexture), args{[]byte("{}")}, &NormalTexture{Scale: Float64(1)}, false},
		{"empty", new(NormalTexture), args{[]byte(`{"scale": 0, "index": 0}`)}, &NormalTexture{Scale: Float64(0), Index: Index(0)}, false},
		{"nodefault", new(NormalTexture), args{[]byte(`{"scale": 0.5, "index":2}`)}, &NormalTexture{Scale: Float64(0.5), Index: Index(2)}, false},
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
		{"default", new(OcclusionTexture), args{[]byte("{}")}, &OcclusionTexture{Strength: Float64(1)}, false},
		{"empty", new(OcclusionTexture), args{[]byte(`{"strength": 0, "index": 0}`)}, &OcclusionTexture{Strength: Float64(0), Index: Index(0)}, false},
		{"nodefault", new(OcclusionTexture), args{[]byte(`{"strength": 0.5, "index":2}`)}, &OcclusionTexture{Strength: Float64(0.5), Index: Index(2)}, false},
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
		{"default", new(PBRMetallicRoughness), args{[]byte("{}")}, &PBRMetallicRoughness{BaseColorFactor: NewRGBA(), MetallicFactor: Float64(1), RoughnessFactor: Float64(1)}, false},
		{"nodefault", new(PBRMetallicRoughness), args{[]byte(`{"baseColorFactor": [0.1,0.2,0.6,0.7],"metallicFactor":0.5,"roughnessFactor":0.1}`)}, &PBRMetallicRoughness{
			BaseColorFactor: &RGBA{R: 0.1, G: 0.2, B: 0.6, A: 0.7}, MetallicFactor: Float64(0.5), RoughnessFactor: Float64(0.1),
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
		}, []byte(`{"camera":1,"skin":1,"matrix":[1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],"mesh":1,"rotation":[1,0,0,0],"scale":[1,0,0],"translation":[1,0,0]}`), false},
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
		{"default", &Material{AlphaCutoff: Float64(0.5), AlphaMode: Opaque}, []byte(`{}`), false},
		{"empty", &Material{AlphaMode: Blend}, []byte(`{"alphaMode":"BLEND"}`), false},
		{"nodefault", &Material{AlphaCutoff: Float64(1), AlphaMode: Blend}, []byte(`{"alphaMode":"BLEND","alphaCutoff":1}`), false},
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
		{"default", &NormalTexture{Scale: Float64(1)}, []byte(`{}`), false},
		{"empty", &NormalTexture{Index: Index(0)}, []byte(`{"index":0}`), false},
		{"nodefault", &NormalTexture{Index: Index(1), Scale: Float64(0.5)}, []byte(`{"index":1,"scale":0.5}`), false},
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
		{"default", &OcclusionTexture{Strength: Float64(1)}, []byte(`{}`), false},
		{"empty", &OcclusionTexture{Index: Index(0)}, []byte(`{"index":0}`), false},
		{"nodefault", &OcclusionTexture{Index: Index(1), Strength: Float64(0.5)}, []byte(`{"index":1,"strength":0.5}`), false},
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
		{"default", &PBRMetallicRoughness{MetallicFactor: Float64(1), RoughnessFactor: Float64(1), BaseColorFactor: NewRGBA()}, []byte(`{}`), false},
		{"empty", &PBRMetallicRoughness{MetallicFactor: Float64(0), RoughnessFactor: Float64(0)}, []byte(`{"metallicFactor":0,"roughnessFactor":0}`), false},
		{"nodefault", &PBRMetallicRoughness{MetallicFactor: Float64(0.5), RoughnessFactor: Float64(0.5), BaseColorFactor: &RGBA{R: 1, G: 0.5, B: 1, A: 1}}, []byte(`{"baseColorFactor":[1,0.5,1,1],"metallicFactor":0.5,"roughnessFactor":0.5}`), false},
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
	RegisterExtension("fake_ext", func(data []byte) (interface{}, error) {
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
