package gltf

import (
	"reflect"
	"testing"
)

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
		{"jpg", &Image{URI: "data:image/png;base64,dsjdsaGGUDXGA"}, true},
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

func TestImage_Data(t *testing.T) {
	tests := []struct {
		name    string
		im      *Image
		want    []uint8
		wantErr bool
	}{
		{"error", &Image{URI: "data:image/png;base64,_"}, []uint8{}, true},
		{"empty", &Image{URI: "data:image/png;base64,"}, []uint8{}, false},
		{"test", &Image{URI: "data:image/png;base64,TEST"}, []uint8{76, 68, 147}, false},
		{"complex", &Image{URI: "data:image/png;base64,YW55IGNhcm5hbCBwbGVhcw=="}, []uint8{97, 110, 121, 32, 99, 97, 114, 110, 97, 108, 32, 112, 108, 101, 97, 115}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.im.Data()
			if (err != nil) != tt.wantErr {
				t.Errorf("Image.Data() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Image.Data() = %v, want %v", got, tt.want)
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
			Matrix:   [16]float32{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1},
			Rotation: [4]float32{0, 0, 0, 1},
			Scale:    [3]float32{1, 1, 1},
			Camera:   -1,
			Mesh:     -1,
			Skin:     -1,
		}, false},
		{"nodefault", new(Node), args{[]byte(`{"matrix":[1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1],"rotation":[0,0,0,1],"scale":[1,1,1],"camera":0,"mesh":0,"skin":0}`)}, &Node{
			Matrix:   [16]float32{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1},
			Rotation: [4]float32{0, 0, 0, 1},
			Scale:    [3]float32{1, 1, 1},
			Camera:   0,
			Mesh:     0,
			Skin:     0,
		}, false},
		{"nodefault", new(Node), args{[]byte(`{"matrix":[1,2,2,0,0,1,3,4,0,0,1,0,5,0,0,5],"rotation":[1,2,3,4],"scale":[2,3,4],"camera":1,"mesh":2,"skin":3}`)}, &Node{
			Matrix:   [16]float32{1, 2, 2, 0, 0, 1, 3, 4, 0, 0, 1, 0, 5, 0, 0, 5},
			Rotation: [4]float32{1, 2, 3, 4},
			Scale:    [3]float32{2, 3, 4},
			Camera:   1,
			Mesh:     2,
			Skin:     3,
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

func TestPrimitive_UnmarshalJSON(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		p       *Primitive
		args    args
		want    *Primitive
		wantErr bool
	}{
		{"default", new(Primitive), args{[]byte("{}")}, &Primitive{Mode: Triangles, Indices: -1, Material: -1}, false},
		{"empty", new(Primitive), args{[]byte(`{"mode": 0, "indices": 0, "material": 0}`)}, &Primitive{Mode: Points, Indices: 0, Material: 0}, false},
		{"nodefault", new(Primitive), args{[]byte(`{"mode": 1, "indices": 2, "material": 3}`)}, &Primitive{Mode: Lines, Indices: 2, Material: 3}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.p.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Primitive.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(tt.p, tt.want) {
				t.Errorf("Primitive.UnmarshalJSON() = %v, want %v", tt.p, tt.want)
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
		{"default", new(Material), args{[]byte("{}")}, &Material{AlphaCutoff: 0.5, AlphaMode: Opaque}, false},
		{"nodefault", new(Material), args{[]byte(`{"alphaCutoff": 0.2, "alphaMode": "MASK"}`)}, &Material{AlphaCutoff: 0.2, AlphaMode: Mask}, false},
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
		{"default", new(NormalTexture), args{[]byte("{}")}, &NormalTexture{Scale: 1, Index: -1}, false},
		{"empty", new(NormalTexture), args{[]byte(`{"scale": 0, "index": 0}`)}, &NormalTexture{Scale: 0, Index: 0}, false},
		{"nodefault", new(NormalTexture), args{[]byte(`{"scale": 0.5, "index":2}`)}, &NormalTexture{Scale: 0.5, Index: 2}, false},
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
		{"default", new(OcclusionTexture), args{[]byte("{}")}, &OcclusionTexture{Strength: 1, Index: -1}, false},
		{"empty", new(OcclusionTexture), args{[]byte(`{"strength": 0, "index": 0}`)}, &OcclusionTexture{Strength: 0, Index: 0}, false},
		{"nodefault", new(OcclusionTexture), args{[]byte(`{"strength": 0.5, "index":2}`)}, &OcclusionTexture{Strength: 0.5, Index: 2}, false},
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
		{"default", new(PBRMetallicRoughness), args{[]byte("{}")}, &PBRMetallicRoughness{BaseColorFactor: [4]float32{1, 1, 1, 1}, MetallicFactor: 1, RoughnessFactor: 1}, false},
		{"nodefault", new(PBRMetallicRoughness), args{[]byte(`{"baseColorFactor": [0.1,0.2,0.6,0.7],"metallicFactor":0.5,"roughnessFactor":0.1}`)}, &PBRMetallicRoughness{
			BaseColorFactor: [4]float32{0.1, 0.2, 0.6, 0.7}, MetallicFactor: 0.5, RoughnessFactor: 0.1,
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

func TestSampler_UnmarshalJSON(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		s       *Sampler
		args    args
		want    *Sampler
		wantErr bool
	}{
		{"default", new(Sampler), args{[]byte("{}")}, &Sampler{WrapS: Repeat, WrapT: Repeat}, false},
		{"nodefault", new(Sampler), args{[]byte(`{"wrapS": 33648, "wrapT": 33071}`)}, &Sampler{WrapS: MirroredRepeat, WrapT: ClampToEdge}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Sampler.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(tt.s, tt.want) {
				t.Errorf("Sampler.UnmarshalJSON() = %v, want %v", tt.s, tt.want)
			}
		})
	}
}

func TestAnimationSampler_UnmarshalJSON(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		as      *AnimationSampler
		args    args
		want    *AnimationSampler
		wantErr bool
	}{
		{"default", new(AnimationSampler), args{[]byte("{}")}, &AnimationSampler{Interpolation: Linear, Input: -1, Output: -1}, false},
		{"empty", new(AnimationSampler), args{[]byte(`{"interpolation": "LINEAR", "input": 0, "output":0}`)}, &AnimationSampler{Interpolation: Linear, Input: 0, Output: 0}, false},
		{"nodefault", new(AnimationSampler), args{[]byte(`{"interpolation": "STEP", "input": 1, "output":2}`)}, &AnimationSampler{Interpolation: Step, Input: 1, Output: 2}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.as.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("AnimationSampler.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(tt.as, tt.want) {
				t.Errorf("AnimationSampler.UnmarshalJSON() = %v, want %v", tt.as, tt.want)
			}
		})
	}
}

func TestAccessor_UnmarshalJSON(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		a       *Accessor
		args    args
		want    *Accessor
		wantErr bool
	}{
		{"default", new(Accessor), args{[]byte("{}")}, &Accessor{BufferView: -1}, false},
		{"empty", new(Accessor), args{[]byte(`{"bufferView": 0}`)}, &Accessor{BufferView: 0}, false},
		{"nodefault", new(Accessor), args{[]byte(`{"bufferView": 1}`)}, &Accessor{BufferView: 1}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.a.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Accessor.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(tt.a, tt.want) {
				t.Errorf("Accessor.UnmarshalJSON() = %v, want %v", tt.a, tt.want)
			}
		})
	}
}

func TestChannelTarget_UnmarshalJSON(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		ch      *ChannelTarget
		args    args
		want    *ChannelTarget
		wantErr bool
	}{
		{"default", new(ChannelTarget), args{[]byte("{}")}, &ChannelTarget{Node: -1}, false},
		{"empty", new(ChannelTarget), args{[]byte(`{"node": 0}`)}, &ChannelTarget{Node: 0}, false},
		{"nodefault", new(ChannelTarget), args{[]byte(`{"node": 1}`)}, &ChannelTarget{Node: 1}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.ch.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("ChannelTarget.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(tt.ch, tt.want) {
				t.Errorf("ChannelTarget.UnmarshalJSON() = %v, want %v", tt.ch, tt.want)
			}
		})
	}
}

func TestChannel_UnmarshalJSON(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		ch      *Channel
		args    args
		want    *Channel
		wantErr bool
	}{
		{"default", new(Channel), args{[]byte("{}")}, &Channel{Sampler: -1}, false},
		{"empty", new(Channel), args{[]byte(`{"sampler": 0}`)}, &Channel{Sampler: 0}, false},
		{"nodefault", new(Channel), args{[]byte(`{"sampler": 1}`)}, &Channel{Sampler: 1}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.ch.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Channel.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(tt.ch, tt.want) {
				t.Errorf("Channel.UnmarshalJSON() = %v, want %v", tt.ch, tt.want)
			}
		})
	}
}

func TestBufferView_UnmarshalJSON(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		b       *BufferView
		args    args
		want    *BufferView
		wantErr bool
	}{
		{"default", new(BufferView), args{[]byte("{}")}, &BufferView{Buffer: -1}, false},
		{"empty", new(BufferView), args{[]byte(`{"Buffer": 0}`)}, &BufferView{Buffer: 0}, false},
		{"nodefault", new(BufferView), args{[]byte(`{"Buffer": 1}`)}, &BufferView{Buffer: 1}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.b.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("BufferView.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(tt.b, tt.want) {
				t.Errorf("BufferView.UnmarshalJSON() = %v, want %v", tt.b, tt.want)
			}
		})
	}
}

func TestTextureInfo_UnmarshalJSON(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		t       *TextureInfo
		args    args
		want    *TextureInfo
		wantErr bool
	}{
		{"default", new(TextureInfo), args{[]byte("{}")}, &TextureInfo{Index: -1}, false},
		{"empty", new(TextureInfo), args{[]byte(`{"index": 0}`)}, &TextureInfo{Index: 0}, false},
		{"nodefault", new(TextureInfo), args{[]byte(`{"index": 1}`)}, &TextureInfo{Index: 1}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.t.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("TextureInfo.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(tt.t, tt.want) {
				t.Errorf("TextureInfo.UnmarshalJSON() = %v, want %v", tt.t, tt.want)
			}
		})
	}
}

func TestSkin_UnmarshalJSON(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		s       *Skin
		args    args
		want    *Skin
		wantErr bool
	}{
		{"default", new(Skin), args{[]byte("{}")}, &Skin{InverseBindMatrices: -1, Skeleton: -1}, false},
		{"empty", new(Skin), args{[]byte(`{"InverseBindMatrices": 0, "skeleton": 0}`)}, &Skin{InverseBindMatrices: 0, Skeleton: 0}, false},
		{"nodefault", new(Skin), args{[]byte(`{"InverseBindMatrices": 1, "skeleton": 2}`)}, &Skin{InverseBindMatrices: 1, Skeleton: 2}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Skin.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(tt.s, tt.want) {
				t.Errorf("Skin.UnmarshalJSON() = %v, want %v", tt.s, tt.want)
			}
		})
	}
}

func TestTexture_UnmarshalJSON(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		t       *Texture
		args    args
		want    *Texture
		wantErr bool
	}{
		{"default", new(Texture), args{[]byte("{}")}, &Texture{Sampler: -1, Source: -1}, false},
		{"empty", new(Texture), args{[]byte(`{"sampler": 0, "source": 0}`)}, &Texture{Sampler: 0, Source: 0}, false},
		{"nodefault", new(Texture), args{[]byte(`{"sampler": 1, "source": 2}`)}, &Texture{Sampler: 1, Source: 2}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.t.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Texture.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(tt.t, tt.want) {
				t.Errorf("Texture.UnmarshalJSON() = %v, want %v", tt.t, tt.want)
			}
		})
	}
}

func TestDocument_UnmarshalJSON(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		d       *Document
		args    args
		want    *Document
		wantErr bool
	}{
		{"default", new(Document), args{[]byte("{}")}, &Document{Scene: -1}, false},
		{"empty", new(Document), args{[]byte(`{"scene": 0}`)}, &Document{Scene: 0}, false},
		{"nodefault", new(Document), args{[]byte(`{"scene": 1}`)}, &Document{Scene: 1}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.d.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Document.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(tt.d, tt.want) {
				t.Errorf("Document.UnmarshalJSON() = %v, want %v", tt.d, tt.want)
			}
		})
	}
}
