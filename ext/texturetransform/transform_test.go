package texturetransform

import (
	"reflect"
	"testing"

	"github.com/qmuntal/gltf"
)

func TestTextureTranform_ScaleOrDefault(t *testing.T) {
	tests := []struct {
		name string
		t    *TextureTranform
		want [2]float32
	}{
		{"default", &TextureTranform{Scale: DefaultScale}, DefaultScale},
		{"zeros", &TextureTranform{Scale: emptyScale}, DefaultScale},
		{"other", &TextureTranform{Scale: [2]float32{1, 2}}, [2]float32{1, 2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.ScaleOrDefault(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TextureTranform.ScaleOrDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTextureTranform_UnmarshalJSON(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		t       *TextureTranform
		args    args
		want    *TextureTranform
		wantErr bool
	}{
		{"default", new(TextureTranform), args{[]byte("{}")}, &TextureTranform{Scale: DefaultScale}, false},
		{"nodefault", new(TextureTranform), args{[]byte(`{"offset": [0.1,0.2],"rotation":1.57,"scale":[1, -1],"texCoord":2}`)}, &TextureTranform{
			Offset: [2]float32{0.1, 0.2}, Rotation: 1.57, Scale: [2]float32{1, -1}, TexCoord: gltf.Index(2),
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.t.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("TextureTranform.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(tt.t, tt.want) {
				t.Errorf("TextureTranform.UnmarshalJSON() = %v, want %v", tt.t, tt.want)
			}
		})
	}
}

func TestTextureTranform_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		t       *TextureTranform
		want    []byte
		wantErr bool
	}{
		{"default", &TextureTranform{Scale: DefaultScale}, []byte(`{}`), false},
		{"empty", &TextureTranform{}, []byte(`{}`), false},
		{"nodefault", &TextureTranform{Offset: [2]float32{0.1, 0.2}, Rotation: 1.57, Scale: [2]float32{1, -1}, TexCoord: gltf.Index(2)}, []byte(`{"offset":[0.1,0.2],"rotation":1.57,"scale":[1,-1],"texCoord":2}`), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.t.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("TextureTranform.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TextureTranform.MarshalJSON() = %v, want %v", got, tt.want)
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
		want    interface{}
		wantErr bool
	}{
		{"base", args{[]byte("{}")}, &TextureTranform{Scale: DefaultScale}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Unmarshal(tt.args.data)
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
