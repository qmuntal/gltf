package specular

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/qmuntal/gltf"
)

func TestPBRSpecularGlossiness_UnmarshalJSON(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		p       *PBRSpecularGlossiness
		args    args
		want    *PBRSpecularGlossiness
		wantErr bool
	}{
		{"default", new(PBRSpecularGlossiness), args{[]byte("{}")}, &PBRSpecularGlossiness{DiffuseFactor: gltf.NewRGBA(), SpecularFactor: gltf.NewRGB(), GlossinessFactor: gltf.Float64(1)}, false},
		{"nodefault", new(PBRSpecularGlossiness), args{[]byte(`{"diffuseFactor": [0.1,0.2,0.3,0.4],"specularFactor":[0.5,0.6,0.7],"glossinessFactor":0.5}`)}, &PBRSpecularGlossiness{
			DiffuseFactor: &gltf.RGBA{R: 0.1, G: 0.2, B: 0.3, A: 0.4}, SpecularFactor: &gltf.RGB{R: 0.5, G: 0.6, B: 0.7}, GlossinessFactor: gltf.Float64(0.5),
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.p.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("PBRSpecularGlossiness.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(tt.p, tt.want) {
				t.Errorf("PBRSpecularGlossiness.UnmarshalJSON() = %v, want %v", tt.p, tt.want)
			}
		})
	}
}

func TestPBRSpecularGlossiness_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		p       *PBRSpecularGlossiness
		want    []byte
		wantErr bool
	}{
		{"default", &PBRSpecularGlossiness{GlossinessFactor: gltf.Float64(1), DiffuseFactor: gltf.NewRGBA(), SpecularFactor: gltf.NewRGB()}, []byte(`{}`), false},
		{"empty", &PBRSpecularGlossiness{}, []byte(`{}`), false},
		{"nodefault", &PBRSpecularGlossiness{GlossinessFactor: gltf.Float64(0.5), DiffuseFactor: &gltf.RGBA{R: 1, G: 0.5, B: 1, A: 1}, SpecularFactor: &gltf.RGB{R: 1, G: 1, B: 0.5}}, []byte(`{"diffuseFactor":[1,0.5,1,1],"specularFactor":[1,1,0.5],"glossinessFactor":0.5}`), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.p.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("PBRSpecularGlossiness.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PBRSpecularGlossiness.MarshalJSON() = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want json.Unmarshaler
	}{
		{"base", new(PBRSpecularGlossiness)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
