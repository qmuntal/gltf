package modeler

import (
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
		{"base", &Modeler{Document: &gltf.Document{}, Compress: CompressionSafe}},
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
