package modeler

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"testing"

	"github.com/go-test/deep"
	"github.com/qmuntal/gltf"
)

func TestNewModeler(t *testing.T) {
	tests := []struct {
		name string
		want *Modeler
	}{
		{"base", &Modeler{Document: &gltf.Document{}, Compression: CompressionSafe}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewModeler(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewModeler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestModeler_AddNormal(t *testing.T) {
	type args struct {
		bufferIndex uint32
		data        [][3]float32
	}
	tests := []struct {
		name    string
		m       *Modeler
		args    args
		want    uint32
		wantDoc *gltf.Document
	}{
		{"base", &Modeler{Document: &gltf.Document{
			Accessors: []gltf.Accessor{{}},
			Buffers:   []gltf.Buffer{{ByteLength: 10}},
		}}, args{0, [][3]float32{{1, 2, 3}}}, 1, &gltf.Document{
			Accessors: []gltf.Accessor{
				{},
				{BufferView: gltf.Index(0), Count: 1, Type: gltf.Vec3, ComponentType: gltf.Float},
			},
			BufferViews: []gltf.BufferView{
				{ByteLength: 12, Target: gltf.ArrayBuffer},
			},
			Buffers: []gltf.Buffer{
				{ByteLength: 22, Data: []byte{0, 0, 128, 63, 0, 0, 0, 64, 0, 0, 64, 64}},
			},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.m.AddNormal(tt.args.bufferIndex, tt.args.data)
			if tt.want != got {
				t.Errorf("Modeler.AddNormal() = %v, want %v", got, tt.want)
				return
			}
			if diff := deep.Equal(tt.m.Document, tt.wantDoc); diff != nil {
				t.Errorf("Modeler.AddNormal() = %v", diff)
				return
			}
		})
	}
}

func TestModeler_AddTangent(t *testing.T) {
	type args struct {
		bufferIndex uint32
		data        [][4]float32
	}
	tests := []struct {
		name    string
		m       *Modeler
		args    args
		want    uint32
		wantDoc *gltf.Document
	}{
		{"base", &Modeler{Document: &gltf.Document{
			Accessors: []gltf.Accessor{{}},
			Buffers:   []gltf.Buffer{{ByteLength: 10}},
		}}, args{0, [][4]float32{{1, 2, 3, 4}, {}}}, 1, &gltf.Document{
			Accessors: []gltf.Accessor{
				{},
				{BufferView: gltf.Index(0), Count: 2, Type: gltf.Vec4, ComponentType: gltf.Float},
			},
			BufferViews: []gltf.BufferView{
				{ByteLength: 32, Target: gltf.ArrayBuffer},
			},
			Buffers: []gltf.Buffer{
				{ByteLength: 42, Data: []byte{0, 0, 128, 63, 0, 0, 0, 64, 0, 0, 64, 64, 0, 0, 128, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}},
			},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.m.AddTangent(tt.args.bufferIndex, tt.args.data)
			if tt.want != got {
				t.Errorf("Modeler.AddTangent() = %v, want %v", got, tt.want)
				return
			}
			if diff := deep.Equal(tt.m.Document, tt.wantDoc); diff != nil {
				t.Errorf("Modeler.AddTangent() = %v", diff)
				return
			}
		})
	}
}

func TestModeler_AddPosition(t *testing.T) {
	type args struct {
		bufferIndex uint32
		data        [][3]float32
	}
	tests := []struct {
		name    string
		m       *Modeler
		args    args
		want    uint32
		wantDoc *gltf.Document
	}{
		{"base", &Modeler{Document: &gltf.Document{
			Accessors: []gltf.Accessor{{}},
			Buffers:   []gltf.Buffer{{ByteLength: 10}},
		}}, args{0, [][3]float32{{1, 2, 3}, {0, 0, -1}}}, 1, &gltf.Document{
			Accessors: []gltf.Accessor{
				{},
				{BufferView: gltf.Index(0), Count: 2, Type: gltf.Vec3, ComponentType: gltf.Float, Max: []float64{1, 2, 3}, Min: []float64{0, 0, -1}},
			},
			BufferViews: []gltf.BufferView{
				{ByteLength: 24, Target: gltf.ArrayBuffer},
			},
			Buffers: []gltf.Buffer{
				{ByteLength: 34, Data: []byte{0, 0, 128, 63, 0, 0, 0, 64, 0, 0, 64, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128, 191}},
			},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.m.AddPosition(tt.args.bufferIndex, tt.args.data)
			if tt.want != got {
				t.Errorf("Modeler.AddPosition() = %v, want %v", got, tt.want)
				return
			}
			if diff := deep.Equal(tt.m.Document, tt.wantDoc); diff != nil {
				t.Errorf("Modeler.AddPosition() = %v", diff)
				return
			}
		})
	}
}

func TestModeler_AddJoints(t *testing.T) {
	type args struct {
		bufferIndex uint32
		data        interface{}
	}
	tests := []struct {
		name    string
		m       *Modeler
		args    args
		want    uint32
		wantDoc *gltf.Document
	}{
		{"uint8", &Modeler{Document: &gltf.Document{
			Accessors: []gltf.Accessor{{}},
			Buffers:   []gltf.Buffer{{ByteLength: 10}},
		}}, args{0, [][4]uint8{{1, 2, 3, 4}}}, 1, &gltf.Document{
			Accessors: []gltf.Accessor{
				{},
				{BufferView: gltf.Index(0), Count: 1, Type: gltf.Vec4, ComponentType: gltf.UnsignedByte},
			},
			BufferViews: []gltf.BufferView{
				{ByteLength: 4, Target: gltf.ArrayBuffer},
			},
			Buffers: []gltf.Buffer{
				{ByteLength: 14, Data: []byte{1, 2, 3, 4}},
			},
		}},
		{"uint16", &Modeler{Document: &gltf.Document{
			Accessors: []gltf.Accessor{{}},
			Buffers:   []gltf.Buffer{{ByteLength: 10}},
		}}, args{0, [][4]uint16{{1, 2, 3, 4}}}, 1, &gltf.Document{
			Accessors: []gltf.Accessor{
				{},
				{BufferView: gltf.Index(0), Count: 1, Type: gltf.Vec4, ComponentType: gltf.UnsignedShort},
			},
			BufferViews: []gltf.BufferView{
				{ByteLength: 8, Target: gltf.ArrayBuffer},
			},
			Buffers: []gltf.Buffer{
				{ByteLength: 18, Data: []byte{1, 0, 2, 0, 3, 0, 4, 0}},
			},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.m.AddJoints(tt.args.bufferIndex, tt.args.data)
			if tt.want != got {
				t.Errorf("Modeler.AddJoints() = %v, want %v", got, tt.want)
				return
			}
			if diff := deep.Equal(tt.m.Document, tt.wantDoc); diff != nil {
				t.Errorf("Modeler.AddJoints() = %v", diff)
				return
			}
		})
	}
}

func TestModeler_AddWeights(t *testing.T) {
	type args struct {
		bufferIndex uint32
		data        interface{}
	}
	tests := []struct {
		name    string
		m       *Modeler
		args    args
		want    uint32
		wantDoc *gltf.Document
	}{
		{"uint8", &Modeler{Document: &gltf.Document{
			Accessors: []gltf.Accessor{{}},
			Buffers:   []gltf.Buffer{{ByteLength: 10}},
		}}, args{0, [][4]uint8{{1, 2, 3, 4}}}, 1, &gltf.Document{
			Accessors: []gltf.Accessor{
				{},
				{BufferView: gltf.Index(0), Count: 1, Type: gltf.Vec4, ComponentType: gltf.UnsignedByte, Normalized: true},
			},
			BufferViews: []gltf.BufferView{
				{ByteLength: 4, Target: gltf.ArrayBuffer},
			},
			Buffers: []gltf.Buffer{
				{ByteLength: 14, Data: []byte{1, 2, 3, 4}},
			},
		}},
		{"uint16", &Modeler{Document: &gltf.Document{
			Accessors: []gltf.Accessor{{}},
			Buffers:   []gltf.Buffer{{ByteLength: 10}},
		}}, args{0, [][4]uint16{{1, 2, 3, 4}}}, 1, &gltf.Document{
			Accessors: []gltf.Accessor{
				{},
				{BufferView: gltf.Index(0), Count: 1, Type: gltf.Vec4, ComponentType: gltf.UnsignedShort, Normalized: true},
			},
			BufferViews: []gltf.BufferView{
				{ByteLength: 8, Target: gltf.ArrayBuffer},
			},
			Buffers: []gltf.Buffer{
				{ByteLength: 18, Data: []byte{1, 0, 2, 0, 3, 0, 4, 0}},
			},
		}},
		{"float", &Modeler{Document: &gltf.Document{
			Accessors: []gltf.Accessor{{}},
			Buffers:   []gltf.Buffer{{ByteLength: 10}},
		}}, args{0, [][4]float32{{1, 2, 3, 4}, {}}}, 1, &gltf.Document{
			Accessors: []gltf.Accessor{
				{},
				{BufferView: gltf.Index(0), Count: 2, Type: gltf.Vec4, ComponentType: gltf.Float},
			},
			BufferViews: []gltf.BufferView{
				{ByteLength: 32, Target: gltf.ArrayBuffer},
			},
			Buffers: []gltf.Buffer{
				{ByteLength: 42, Data: []byte{0, 0, 128, 63, 0, 0, 0, 64, 0, 0, 64, 64, 0, 0, 128, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}},
			},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.m.AddWeights(tt.args.bufferIndex, tt.args.data)
			if tt.want != got {
				t.Errorf("Modeler.AddWeights() = %v, want %v", got, tt.want)
				return
			}
			if diff := deep.Equal(tt.m.Document, tt.wantDoc); diff != nil {
				t.Errorf("Modeler.AddWeights() = %v", diff)
				return
			}
		})
	}
}

func TestModeler_AddTextureCoord(t *testing.T) {
	type args struct {
		bufferIndex uint32
		data        interface{}
	}
	tests := []struct {
		name    string
		m       *Modeler
		args    args
		want    uint32
		wantDoc *gltf.Document
	}{
		{"uint8", &Modeler{Document: &gltf.Document{
			Accessors: []gltf.Accessor{{}},
			Buffers:   []gltf.Buffer{{ByteLength: 10}},
		}}, args{0, [][2]uint8{{1, 2}}}, 1, &gltf.Document{
			Accessors: []gltf.Accessor{
				{},
				{BufferView: gltf.Index(0), Count: 1, Type: gltf.Vec2, ComponentType: gltf.UnsignedByte, Normalized: true},
			},
			BufferViews: []gltf.BufferView{
				{ByteLength: 2, Target: gltf.ArrayBuffer},
			},
			Buffers: []gltf.Buffer{
				{ByteLength: 12, Data: []byte{1, 2}},
			},
		}},
		{"uint16", &Modeler{Document: &gltf.Document{
			Accessors: []gltf.Accessor{{}},
			Buffers:   []gltf.Buffer{{ByteLength: 10}},
		}}, args{0, [][2]uint16{{1, 2}}}, 1, &gltf.Document{
			Accessors: []gltf.Accessor{
				{},
				{BufferView: gltf.Index(0), Count: 1, Type: gltf.Vec2, ComponentType: gltf.UnsignedShort, Normalized: true},
			},
			BufferViews: []gltf.BufferView{
				{ByteLength: 4, Target: gltf.ArrayBuffer},
			},
			Buffers: []gltf.Buffer{
				{ByteLength: 14, Data: []byte{1, 0, 2, 0}},
			},
		}},
		{"float", &Modeler{Document: &gltf.Document{
			Accessors: []gltf.Accessor{{}},
			Buffers:   []gltf.Buffer{{ByteLength: 10}},
		}}, args{0, [][2]float32{{1, 2}, {}}}, 1, &gltf.Document{
			Accessors: []gltf.Accessor{
				{},
				{BufferView: gltf.Index(0), Count: 2, Type: gltf.Vec2, ComponentType: gltf.Float},
			},
			BufferViews: []gltf.BufferView{
				{ByteLength: 16, Target: gltf.ArrayBuffer},
			},
			Buffers: []gltf.Buffer{
				{ByteLength: 26, Data: []byte{0, 0, 128, 63, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0}},
			},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.m.AddTextureCoord(tt.args.bufferIndex, tt.args.data)
			if tt.want != got {
				t.Errorf("Modeler.AddTextureCoord() = %v, want %v", got, tt.want)
				return
			}
			if diff := deep.Equal(tt.m.Document, tt.wantDoc); diff != nil {
				t.Errorf("Modeler.AddTextureCoord() = %v", diff)
				return
			}
		})
	}
}

func TestModeler_AddIndices(t *testing.T) {
	type args struct {
		bufferIndex uint32
		data        interface{}
	}
	tests := []struct {
		name    string
		m       *Modeler
		args    args
		want    uint32
		wantDoc *gltf.Document
	}{
		{"uint8", &Modeler{Document: &gltf.Document{
			Accessors: []gltf.Accessor{{}},
			Buffers:   []gltf.Buffer{{ByteLength: 10}},
		}}, args{0, []uint8{1, 2}}, 1, &gltf.Document{
			Accessors: []gltf.Accessor{
				{},
				{BufferView: gltf.Index(0), Count: 2, Type: gltf.Scalar, ComponentType: gltf.UnsignedByte},
			},
			BufferViews: []gltf.BufferView{
				{ByteLength: 2, Target: gltf.ElementArrayBuffer},
			},
			Buffers: []gltf.Buffer{
				{ByteLength: 12, Data: []byte{1, 2}},
			},
		}},
		{"uint16", &Modeler{Document: &gltf.Document{
			Accessors: []gltf.Accessor{{}},
			Buffers:   []gltf.Buffer{{ByteLength: 10}},
		}}, args{0, []uint16{1, 2}}, 1, &gltf.Document{
			Accessors: []gltf.Accessor{
				{},
				{BufferView: gltf.Index(0), Count: 2, Type: gltf.Scalar, ComponentType: gltf.UnsignedShort},
			},
			BufferViews: []gltf.BufferView{
				{ByteLength: 4, Target: gltf.ElementArrayBuffer},
			},
			Buffers: []gltf.Buffer{
				{ByteLength: 14, Data: []byte{1, 0, 2, 0}},
			},
		}},
		{"uint16-compress", &Modeler{Document: &gltf.Document{
			Accessors: []gltf.Accessor{{}},
			Buffers:   []gltf.Buffer{{ByteLength: 10}},
		}, Compression: CompressionSafe}, args{0, []uint16{1, 2}}, 1, &gltf.Document{
			Accessors: []gltf.Accessor{
				{},
				{BufferView: gltf.Index(0), Count: 2, Type: gltf.Scalar, ComponentType: gltf.UnsignedByte},
			},
			BufferViews: []gltf.BufferView{
				{ByteLength: 2, Target: gltf.ElementArrayBuffer},
			},
			Buffers: []gltf.Buffer{
				{ByteLength: 12, Data: []byte{1, 2}},
			},
		}},
		{"uint32", &Modeler{Document: &gltf.Document{
			Accessors: []gltf.Accessor{{}},
			Buffers:   []gltf.Buffer{{ByteLength: 10}},
		}}, args{0, []uint32{1, 2}}, 1, &gltf.Document{
			Accessors: []gltf.Accessor{
				{},
				{BufferView: gltf.Index(0), Count: 2, Type: gltf.Scalar, ComponentType: gltf.UnsignedInt},
			},
			BufferViews: []gltf.BufferView{
				{ByteLength: 8, Target: gltf.ElementArrayBuffer},
			},
			Buffers: []gltf.Buffer{
				{ByteLength: 18, Data: []byte{1, 0, 0, 0, 2, 0, 0, 0}},
			},
		}},
		{"uint32-compress", &Modeler{Document: &gltf.Document{
			Accessors: []gltf.Accessor{{}},
			Buffers:   []gltf.Buffer{{ByteLength: 10}},
		}, Compression: CompressionSafe}, args{0, []uint32{1, 2}}, 1, &gltf.Document{
			Accessors: []gltf.Accessor{
				{},
				{BufferView: gltf.Index(0), Count: 2, Type: gltf.Scalar, ComponentType: gltf.UnsignedByte},
			},
			BufferViews: []gltf.BufferView{
				{ByteLength: 2, Target: gltf.ElementArrayBuffer},
			},
			Buffers: []gltf.Buffer{
				{ByteLength: 12, Data: []byte{1, 2}},
			},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.m.AddIndices(tt.args.bufferIndex, tt.args.data)
			if tt.want != got {
				t.Errorf("Modeler.AddIndices() = %v, want %v", got, tt.want)
				return
			}
			if diff := deep.Equal(tt.m.Document, tt.wantDoc); diff != nil {
				t.Errorf("Modeler.AddIndices() = %v", diff)
				return
			}
		})
	}
}

func TestModeler_AddColor(t *testing.T) {
	type args struct {
		bufferIndex uint32
		data        interface{}
	}
	tests := []struct {
		name    string
		m       *Modeler
		args    args
		want    uint32
		wantDoc *gltf.Document
	}{
		{"uint8", &Modeler{Document: &gltf.Document{
			Accessors: []gltf.Accessor{{}},
			Buffers:   []gltf.Buffer{{ByteLength: 10}},
		}}, args{0, [][4]uint8{{1, 2, 3, 4}}}, 1, &gltf.Document{
			Accessors: []gltf.Accessor{
				{},
				{BufferView: gltf.Index(0), Count: 1, Type: gltf.Vec4, ComponentType: gltf.UnsignedByte, Normalized: true},
			},
			BufferViews: []gltf.BufferView{
				{ByteLength: 4, Target: gltf.ArrayBuffer},
			},
			Buffers: []gltf.Buffer{
				{ByteLength: 14, Data: []byte{1, 2, 3, 4}},
			},
		}},
		{"uint16", &Modeler{Document: &gltf.Document{
			Accessors: []gltf.Accessor{{}},
			Buffers:   []gltf.Buffer{{ByteLength: 10}},
		}}, args{0, [][4]uint16{{1, 2, 3, 4}}}, 1, &gltf.Document{
			Accessors: []gltf.Accessor{
				{},
				{BufferView: gltf.Index(0), Count: 1, Type: gltf.Vec4, ComponentType: gltf.UnsignedShort, Normalized: true},
			},
			BufferViews: []gltf.BufferView{
				{ByteLength: 8, Target: gltf.ArrayBuffer},
			},
			Buffers: []gltf.Buffer{
				{ByteLength: 18, Data: []byte{1, 0, 2, 0, 3, 0, 4, 0}},
			},
		}},
		{"float", &Modeler{Document: &gltf.Document{
			Accessors: []gltf.Accessor{{}},
			Buffers:   []gltf.Buffer{{ByteLength: 10}},
		}}, args{0, [][4]float32{{1, 2, 3, 4}, {}}}, 1, &gltf.Document{
			Accessors: []gltf.Accessor{
				{},
				{BufferView: gltf.Index(0), Count: 2, Type: gltf.Vec4, ComponentType: gltf.Float},
			},
			BufferViews: []gltf.BufferView{
				{ByteLength: 32, Target: gltf.ArrayBuffer},
			},
			Buffers: []gltf.Buffer{
				{ByteLength: 42, Data: []byte{0, 0, 128, 63, 0, 0, 0, 64, 0, 0, 64, 64, 0, 0, 128, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}},
			},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.m.AddColor(tt.args.bufferIndex, tt.args.data)
			if tt.want != got {
				t.Errorf("Modeler.AddColor() = %v, want %v", got, tt.want)
				return
			}
			if diff := deep.Equal(tt.m.Document, tt.wantDoc); diff != nil {
				t.Errorf("Modeler.AddColor() = %v", diff)
				return
			}
		})
	}
}

func TestModeler_AddImage(t *testing.T) {
	type args struct {
		bufferIndex uint32
		name        string
		mimeType    string
		r           io.Reader
	}
	tests := []struct {
		name    string
		m       *Modeler
		args    args
		want    uint32
		wantDoc *gltf.Document
		wantErr bool
	}{
		{"base", &Modeler{Document: &gltf.Document{
			Images:  []gltf.Image{{}},
			Buffers: []gltf.Buffer{{ByteLength: 10}},
		}}, args{0, "fake", "fake/type", bytes.NewReader([]byte{1, 2})}, 1, &gltf.Document{
			Images: []gltf.Image{
				{},
				{BufferView: gltf.Index(0), Name: "fake", MimeType: "fake/type"},
			},
			BufferViews: []gltf.BufferView{
				{ByteLength: 2, Target: gltf.None},
			},
			Buffers: []gltf.Buffer{
				{ByteLength: 12, Data: []byte{1, 2}},
			},
		}, false},
		{"buffer", &Modeler{Document: &gltf.Document{
			Images:  []gltf.Image{{}},
			Buffers: []gltf.Buffer{{ByteLength: 10}},
		}}, args{0, "fake", "fake/type", bytes.NewBuffer([]byte{1, 2})}, 1, &gltf.Document{
			Images: []gltf.Image{
				{},
				{BufferView: gltf.Index(0), Name: "fake", MimeType: "fake/type"},
			},
			BufferViews: []gltf.BufferView{
				{ByteLength: 2, Target: gltf.None},
			},
			Buffers: []gltf.Buffer{
				{ByteLength: 12, Data: []byte{1, 2}},
			},
		}, false},
		{"err", &Modeler{Document: &gltf.Document{
			Images:  []gltf.Image{{}},
			Buffers: []gltf.Buffer{{ByteLength: 10}},
		}}, args{0, "fake", "fake/type", &errReader{}}, 0, &gltf.Document{
			Images: []gltf.Image{
				{},
			},
			Buffers: []gltf.Buffer{
				{ByteLength: 10},
			},
		}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.AddImage(tt.args.bufferIndex, tt.args.name, tt.args.mimeType, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("Modeler.AddImage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != got {
				t.Errorf("Modeler.AddImage() = %v, want %v", got, tt.want)
				return
			}
			if diff := deep.Equal(tt.m.Document, tt.wantDoc); diff != nil {
				t.Errorf("Modeler.AddImage() = %v", diff)
				return
			}
		})
	}
}

type errReader struct {
	r io.Reader
}

func (r *errReader) Read(p []byte) (int, error) {
	return 0, errors.New("")
}
