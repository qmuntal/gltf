package specular_test

import (
	"reflect"
	"testing"

	"github.com/qmuntal/gltf"
	"github.com/qmuntal/gltf/ext/specular"
)

func TestPBRSpecularGlossiness_UnmarshalJSON(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		p       *specular.PBRSpecularGlossiness
		args    args
		want    *specular.PBRSpecularGlossiness
		wantErr bool
	}{
		{"default", new(specular.PBRSpecularGlossiness), args{[]byte("{}")}, &specular.PBRSpecularGlossiness{DiffuseFactor: &[4]float64{1, 1, 1, 1}, SpecularFactor: &[3]float64{1, 1, 1}, GlossinessFactor: gltf.Float(1)}, false},
		{"nodefault", new(specular.PBRSpecularGlossiness), args{[]byte(`{"diffuseFactor": [0.1,0.2,0.3,0.4],"specularFactor":[0.5,0.6,0.7],"glossinessFactor":0.5}`)}, &specular.PBRSpecularGlossiness{
			DiffuseFactor: &[4]float64{0.1, 0.2, 0.3, 0.4}, SpecularFactor: &[3]float64{0.5, 0.6, 0.7}, GlossinessFactor: gltf.Float(0.5),
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
		p       *specular.PBRSpecularGlossiness
		want    []byte
		wantErr bool
	}{
		{"default", &specular.PBRSpecularGlossiness{GlossinessFactor: gltf.Float(1), DiffuseFactor: &[4]float64{1, 1, 1, 1}, SpecularFactor: &[3]float64{1, 1, 1}}, []byte(`{}`), false},
		{"empty", &specular.PBRSpecularGlossiness{}, []byte(`{}`), false},
		{"nodefault", &specular.PBRSpecularGlossiness{GlossinessFactor: gltf.Float(0.5), DiffuseFactor: &[4]float64{1, 0.5, 1, 1}, SpecularFactor: &[3]float64{1, 1, 0.5}}, []byte(`{"diffuseFactor":[1,0.5,1,1],"specularFactor":[1,1,0.5],"glossinessFactor":0.5}`), false},
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

func TestUnmarshal(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    any
		wantErr bool
	}{
		{"base", args{[]byte("{}")}, &specular.PBRSpecularGlossiness{DiffuseFactor: &[4]float64{1, 1, 1, 1}, SpecularFactor: &[3]float64{1, 1, 1}, GlossinessFactor: gltf.Float(1)}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := specular.Unmarshal(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Unmarshal() = %v, want %v", got, tt.want)
			}
		})
	}
}
