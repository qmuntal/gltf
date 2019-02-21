package gltf

import (
	"reflect"
	"testing"
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
		{"default", new(PBRSpecularGlossiness), args{[]byte("{}")}, &PBRSpecularGlossiness{DiffuseFactor: [4]float64{1, 1, 1, 1}, SpecularFactor: [3]float64{1, 1, 1}, GlossinessFactor: 1}, false},
		{"nodefault", new(PBRSpecularGlossiness), args{[]byte(`{"diffuseFactor": [0.1,0.2,0.3,0.4],"specularFactor":[0.5,0.6,0.7],"glossinessFactor":0.5}`)}, &PBRSpecularGlossiness{
			DiffuseFactor: [4]float64{0.1, 0.2, 0.3, 0.4}, SpecularFactor: [3]float64{0.5, 0.6, 0.7}, GlossinessFactor: 0.5,
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
		{"default", &PBRSpecularGlossiness{GlossinessFactor: 1, DiffuseFactor: [4]float64{1, 1, 1, 1}, SpecularFactor: [3]float64{1, 1, 1}}, []byte(`{}`), false},
		{"empty", &PBRSpecularGlossiness{GlossinessFactor: 0, DiffuseFactor: [4]float64{0, 0, 0, 0}, SpecularFactor: [3]float64{0, 0, 0}}, []byte(`{"diffuseFactor":[0,0,0,0],"specularFactor":[0,0,0],"glossinessFactor":0}`), false},
		{"nodefault", &PBRSpecularGlossiness{GlossinessFactor: 0.5, DiffuseFactor: [4]float64{1, 0.5, 1, 1}, SpecularFactor: [3]float64{1, 1, 0.5}}, []byte(`{"diffuseFactor":[1,0.5,1,1],"specularFactor":[1,1,0.5],"glossinessFactor":0.5}`), false},
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

func TestExtensions_UnmarshalJSON(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *Extensions
		wantErr bool
	}{
		{"specularGlossiness", args{[]byte(`{"KHR_materials_pbrSpecularGlossiness": {"diffuseFactor":[1,0.5,1,1],"specularFactor":[1,1,0.5],"glossinessFactor":0.5}}`)}, &Extensions{
			ExtPBRSpecularGlossiness: &PBRSpecularGlossiness{GlossinessFactor: 0.5, DiffuseFactor: [4]float64{1, 0.5, 1, 1}, SpecularFactor: [3]float64{1, 1, 0.5}},
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
				t.Errorf("PBRSpecularGlossiness.MarshalJSON() = %v, want %v", &ext, tt.want)
			}
		})
	}
}
