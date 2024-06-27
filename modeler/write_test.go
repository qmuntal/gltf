package modeler

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/go-test/deep"
	"github.com/qmuntal/gltf"
)

func TestAlignment(t *testing.T) {
	doc := gltf.NewDocument()
	WriteIndices(doc, []uint16{0, 1, 2})
	WritePosition(doc, [][3]float32{{0, 0, 0}, {1, 0, 0}, {0, 1, 0}})
	if len(doc.Buffers) != 1 {
		t.Errorf("Testalignment() buffer size = %v, want 1", len(doc.Buffers))
	}
	buffer := doc.Buffers[0]
	want := make([]byte, 44)
	want[2], want[4] = 1, 2
	want[22], want[23] = 0x80, 0x3f
	want[38], want[39] = 0x80, 0x3f
	if diff := deep.Equal(buffer.Data, want); diff != nil {
		t.Errorf("Testalignment() = %v", diff)
		return
	}
}

func TestWriteAttributesInterleaved(t *testing.T) {
	data := [][3]float32{{1, 2, 3}, {0, 0, -1}}
	doc := gltf.NewDocument()
	attrs, err := WriteAttributesInterleaved(doc,
		Attribute{Name: gltf.POSITION, Data: data},
		Attribute{Name: gltf.NORMAL, Data: data},
		Attribute{Name: gltf.TANGENT, Data: [][4]float32{{1, 2, 3, 4}, {1, 2, 3, 4}}},
		Attribute{Name: gltf.TEXCOORD_0, Data: [][2]float32{{1, 2}, {1, 2}}},
		Attribute{Name: gltf.TEXCOORD_1, Data: [][2]float32{{1, 2}, {1, 2}}},
		Attribute{Name: gltf.WEIGHTS_0, Data: [][4]uint8{{1, 2, 3, 4}, {1, 2, 3, 4}}},
		Attribute{Name: gltf.JOINTS_0, Data: [][4]uint8{{1, 2, 3, 4}, {1, 2, 3, 4}}},
		Attribute{Name: gltf.COLOR_0, Data: data},
		Attribute{Name: "COLOR_1", Data: data},
		Attribute{Name: "COLOR_2", Data: data},
	)
	if err != nil {
		t.Fatalf("TestWriteAttributesInterleaved() got error = %v", err)
	}
	if len(doc.Buffers) != 1 {
		t.Errorf("TestWriteAttributesInterleaved() buffer size = %v, want 1", len(doc.Buffers))
	}
	if len(doc.BufferViews) != 1 {
		t.Errorf("TestWriteAttributesInterleaved() buffer views size = %v, want 1", len(doc.BufferViews))
	}
	if len(doc.Accessors) != 10 {
		t.Errorf("TestWriteAttributesInterleaved() accessors size = %v, want 10", len(doc.Accessors))
	}
	want := map[string]uint32{
		gltf.POSITION:   0,
		gltf.NORMAL:     1,
		gltf.TANGENT:    2,
		gltf.TEXCOORD_0: 3,
		gltf.TEXCOORD_1: 4,
		gltf.WEIGHTS_0:  5,
		gltf.JOINTS_0:   6,
		gltf.COLOR_0:    7,
		"COLOR_1":       8,
		"COLOR_2":       9,
	}
	if diff := deep.Equal(attrs, want); diff != nil {
		t.Errorf("TestWriteAttributesInterleaved() = %v", diff)
	}
	accs := []*gltf.Accessor{
		{BufferView: gltf.Index(0), Count: 2, Type: gltf.AccessorVec3, Min: []float64{0, 0, -1}, Max: []float64{1, 2, 3}},
		{BufferView: gltf.Index(0), Count: 2, Type: gltf.AccessorVec3, ByteOffset: 12},
		{BufferView: gltf.Index(0), Count: 2, Type: gltf.AccessorVec4, ByteOffset: 24},
		{BufferView: gltf.Index(0), Count: 2, Type: gltf.AccessorVec2, ByteOffset: 40},
		{BufferView: gltf.Index(0), Count: 2, Type: gltf.AccessorVec2, ByteOffset: 48},
		{BufferView: gltf.Index(0), Count: 2, Type: gltf.AccessorVec4, ByteOffset: 56, ComponentType: gltf.ComponentUbyte, Normalized: true},
		{BufferView: gltf.Index(0), Count: 2, Type: gltf.AccessorVec4, ByteOffset: 60, ComponentType: gltf.ComponentUbyte},
		{BufferView: gltf.Index(0), Count: 2, Type: gltf.AccessorVec3, ByteOffset: 64},
		{BufferView: gltf.Index(0), Count: 2, Type: gltf.AccessorVec3, ByteOffset: 76},
		{BufferView: gltf.Index(0), Count: 2, Type: gltf.AccessorVec3, ByteOffset: 88},
	}
	if diff := deep.Equal(accs, doc.Accessors); diff != nil {
		t.Errorf("TestWriteAttributesInterleaved() = %v", diff)
	}
}

func TestWriteAttributesInterleaved_OnlyPosition(t *testing.T) {
	doc := gltf.NewDocument()
	_, err := WriteAttributesInterleaved(doc,
		Attribute{Name: gltf.POSITION, Data: [][3]float32{{1, 2, 3}, {0, 0, -1}}},
		Attribute{Name: gltf.TANGENT, Data: make([][4]float32, 0)},
		Attribute{Name: "COLOR_1"})
	if err != nil {
		t.Fatalf("TestWriteAttributesInterleaved_OnlyPosition() got error = %v", err)
	}
	if len(doc.Accessors) != 1 {
		t.Errorf("TestWriteAttributesInterleaved_OnlyPosition() accessors size = %v, want 1", len(doc.Accessors))
	}
}

func TestWriteAttributesInterleaved_Error(t *testing.T) {
	doc := gltf.NewDocument()
	_, err := WriteAttributesInterleaved(doc,
		Attribute{Name: gltf.POSITION, Data: [][3]float32{{1, 2, 3}, {0, 0, -1}}},
		Attribute{Name: gltf.COLOR_0, Data: [][3]float32{{1, 2, 3}}},
	)
	if err == nil {
		t.Error("TestWriteAttributesInterleaved_Error() expected an error")
	}
}

func TestWriteAccessorsInterleaved(t *testing.T) {
	doc := gltf.NewDocument()
	indices, err := WriteAccessorsInterleaved(doc,
		[][3]float32{{1, 2, 3}, {0, 0, -1}},
		[][4]float32{{1, 2, 3, 4}, {1, 2, 3, 4}},
		[][3]float32{{3, 1, 2}, {4, 0, 1}},
	)
	if err != nil {
		t.Fatalf("TestWriteAccessorsInterleaved() got error = %v", err)
	}
	if len(doc.Buffers) != 1 {
		t.Errorf("TestWriteAccessorsInterleaved() buffer size = %v, want 1", len(doc.Buffers))
	}
	if len(doc.BufferViews) != 1 {
		t.Errorf("TestWriteAccessorsInterleaved() buffer views size = %v, want 1", len(doc.BufferViews))
	}
	if got := doc.Accessors[indices[0]].ByteOffset; got != 0 {
		t.Errorf("TestWriteAccessorsInterleaved() accessor 0 has ByteOffset = %v, want 0", got)
	}
	if got := doc.Accessors[indices[1]].ByteOffset; got != 12 {
		t.Errorf("TestWriteAccessorsInterleaved() accessor 1 has ByteOffset = %v, want 12", got)
	}
	if got := doc.Accessors[indices[2]].ByteOffset; got != 28 {
		t.Errorf("TestWriteAccessorsInterleaved() accessor 2 has ByteOffset = %v, want 28", got)
	}
	buffer := doc.Buffers[0]
	want := []byte{
		0, 0, 128, 63, 0, 0, 0, 64, 0, 0, 64, 64,
		0, 0, 128, 63, 0, 0, 0, 64, 0, 0, 64, 64, 0, 0, 128, 64,
		0, 0, 64, 64, 0, 0, 128, 63, 0, 0, 0, 64,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128, 191,
		0, 0, 128, 63, 0, 0, 0, 64, 0, 0, 64, 64, 0, 0, 128, 64,
		0, 0, 128, 64, 0, 0, 0, 0, 0, 0, 128, 63,
	}
	if diff := deep.Equal(buffer.Data, want); diff != nil {
		t.Errorf("TestWriteAccessorsInterleaved() = %v", diff)
		return
	}
}

func TestWriteAccessorsInterleaved_Error(t *testing.T) {
	doc := gltf.NewDocument()
	_, err := WriteAccessorsInterleaved(doc,
		[][3]float32{{1, 2, 3}, {0, 0, -1}},
		[][3]float32{{3, 1, 2}},
	)
	if err == nil {
		t.Error("TestWriteAccessorsInterleaved_Error() expected an error")
	}
}

func TestWriteBufferViewInterleaved(t *testing.T) {
	doc := gltf.NewDocument()
	_, err := WriteBufferViewInterleaved(doc,
		[][3]float32{{1, 2, 3}, {0, 0, -1}},
		[][4]float32{{1, 2, 3, 4}, {1, 2, 3, 4}},
		[][3]float32{{3, 1, 2}, {4, 0, 1}},
	)
	if err != nil {
		t.Fatalf("TestWriteBufferViewInterleaved() got error = %v", err)
	}
	if len(doc.Buffers) != 1 {
		t.Errorf("TestWriteBufferViewInterleaved() buffer size = %v, want 1", len(doc.Buffers))
	}
	buffer := doc.Buffers[0]
	want := []byte{
		0, 0, 128, 63, 0, 0, 0, 64, 0, 0, 64, 64,
		0, 0, 128, 63, 0, 0, 0, 64, 0, 0, 64, 64, 0, 0, 128, 64,
		0, 0, 64, 64, 0, 0, 128, 63, 0, 0, 0, 64,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128, 191,
		0, 0, 128, 63, 0, 0, 0, 64, 0, 0, 64, 64, 0, 0, 128, 64,
		0, 0, 128, 64, 0, 0, 0, 0, 0, 0, 128, 63,
	}
	if diff := deep.Equal(buffer.Data, want); diff != nil {
		t.Errorf("TestWriteBufferViewInterleaved() = %v", diff)
		return
	}
}

func TestWriteBufferViewInterleaved_Error(t *testing.T) {
	doc := gltf.NewDocument()
	_, err := WriteBufferViewInterleaved(doc,
		[][3]float32{{1, 2, 3}, {0, 0, -1}},
		[][3]float32{{3, 1, 2}},
	)
	if err == nil {
		t.Error("TestWriteBufferViewInterleaved_Error() expected an error")
	}
}

func TestWriteNormal(t *testing.T) {
	type args struct {
		data [][3]float32
	}
	tests := []struct {
		name    string
		m       *gltf.Document
		args    args
		want    uint32
		wantDoc *gltf.Document
	}{
		{"base", &gltf.Document{
			Accessors: []*gltf.Accessor{{}},
			Buffers:   []*gltf.Buffer{{ByteLength: 10}},
		}, args{[][3]float32{{1, 2, 3}}}, 1, &gltf.Document{
			Accessors: []*gltf.Accessor{
				{},
				{BufferView: gltf.Index(0), Count: 1, Type: gltf.AccessorVec3, ComponentType: gltf.ComponentFloat},
			},
			BufferViews: []*gltf.BufferView{
				{ByteLength: 12, Target: gltf.TargetArrayBuffer},
			},
			Buffers: []*gltf.Buffer{
				{ByteLength: 22, Data: []byte{0, 0, 128, 63, 0, 0, 0, 64, 0, 0, 64, 64}},
			},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WriteNormal(tt.m, tt.args.data)
			if tt.want != got {
				t.Errorf("WriteNormal() = %v, want %v", got, tt.want)
				return
			}
			if diff := deep.Equal(tt.m, tt.wantDoc); diff != nil {
				t.Errorf("WriteNormal() = %v", diff)
				return
			}
		})
	}
}

func TestWriteTangent(t *testing.T) {
	type args struct {
		data [][4]float32
	}
	tests := []struct {
		name    string
		m       *gltf.Document
		args    args
		want    uint32
		wantDoc *gltf.Document
	}{
		{"base", &gltf.Document{
			Accessors: []*gltf.Accessor{{}},
			Buffers:   []*gltf.Buffer{{ByteLength: 10}},
		}, args{[][4]float32{{1, 2, 3, 4}, {}}}, 1, &gltf.Document{
			Accessors: []*gltf.Accessor{
				{},
				{BufferView: gltf.Index(0), Count: 2, Type: gltf.AccessorVec4, ComponentType: gltf.ComponentFloat},
			},
			BufferViews: []*gltf.BufferView{
				{ByteLength: 32, Target: gltf.TargetArrayBuffer},
			},
			Buffers: []*gltf.Buffer{
				{ByteLength: 42, Data: []byte{0, 0, 128, 63, 0, 0, 0, 64, 0, 0, 64, 64, 0, 0, 128, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}},
			},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WriteTangent(tt.m, tt.args.data)
			if tt.want != got {
				t.Errorf("WriteTangent() = %v, want %v", got, tt.want)
				return
			}
			if diff := deep.Equal(tt.m, tt.wantDoc); diff != nil {
				t.Errorf("WriteTangent() = %v", diff)
				return
			}
		})
	}
}

func TestWritePosition(t *testing.T) {
	type args struct {
		data [][3]float32
	}
	tests := []struct {
		name    string
		m       *gltf.Document
		args    args
		want    uint32
		wantDoc *gltf.Document
	}{
		{"base", &gltf.Document{
			Accessors: []*gltf.Accessor{{}},
			Buffers:   []*gltf.Buffer{{ByteLength: 10}},
		}, args{[][3]float32{{1, 2, 3}, {0, 0, -1}}}, 1, &gltf.Document{
			Accessors: []*gltf.Accessor{
				{},
				{BufferView: gltf.Index(0), Count: 2, Type: gltf.AccessorVec3, ComponentType: gltf.ComponentFloat, Max: []float64{1, 2, 3}, Min: []float64{0, 0, -1}},
			},
			BufferViews: []*gltf.BufferView{
				{ByteLength: 24, Target: gltf.TargetArrayBuffer},
			},
			Buffers: []*gltf.Buffer{
				{ByteLength: 34, Data: []byte{0, 0, 128, 63, 0, 0, 0, 64, 0, 0, 64, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128, 191}},
			},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WritePosition(tt.m, tt.args.data)
			if tt.want != got {
				t.Errorf("WritePosition() = %v, want %v", got, tt.want)
				return
			}
			if diff := deep.Equal(tt.m, tt.wantDoc); diff != nil {
				t.Errorf("WritePosition() = %v", diff)
				return
			}
		})
	}
}

func TestWriteJoints(t *testing.T) {
	type args struct {
		data any
	}
	tests := []struct {
		name    string
		m       *gltf.Document
		args    args
		want    uint32
		wantDoc *gltf.Document
	}{
		{"uint8", &gltf.Document{
			Accessors: []*gltf.Accessor{{}},
			Buffers:   []*gltf.Buffer{{ByteLength: 10}},
		}, args{[][4]uint8{{1, 2, 3, 4}}}, 1, &gltf.Document{
			Accessors: []*gltf.Accessor{
				{},
				{BufferView: gltf.Index(0), Count: 1, Type: gltf.AccessorVec4, ComponentType: gltf.ComponentUbyte},
			},
			BufferViews: []*gltf.BufferView{
				{ByteLength: 4, Target: gltf.TargetArrayBuffer},
			},
			Buffers: []*gltf.Buffer{
				{ByteLength: 14, Data: []byte{1, 2, 3, 4}},
			},
		}},
		{"uint16", &gltf.Document{
			Accessors: []*gltf.Accessor{{}},
			Buffers:   []*gltf.Buffer{{ByteLength: 10}},
		}, args{[][4]uint16{{1, 2, 3, 4}}}, 1, &gltf.Document{
			Accessors: []*gltf.Accessor{
				{},
				{BufferView: gltf.Index(0), Count: 1, Type: gltf.AccessorVec4, ComponentType: gltf.ComponentUshort},
			},
			BufferViews: []*gltf.BufferView{
				{ByteLength: 8, Target: gltf.TargetArrayBuffer},
			},
			Buffers: []*gltf.Buffer{
				{ByteLength: 18, Data: []byte{1, 0, 2, 0, 3, 0, 4, 0}},
			},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WriteJoints(tt.m, tt.args.data)
			if tt.want != got {
				t.Errorf("WriteJoints() = %v, want %v", got, tt.want)
				return
			}
			if diff := deep.Equal(tt.m, tt.wantDoc); diff != nil {
				t.Errorf("WriteJoints() = %v", diff)
				return
			}
		})
	}
}

func TestWriteWeights(t *testing.T) {
	type args struct {
		data any
	}
	tests := []struct {
		name    string
		m       *gltf.Document
		args    args
		want    uint32
		wantDoc *gltf.Document
	}{
		{"uint8", &gltf.Document{
			Accessors: []*gltf.Accessor{{}},
			Buffers:   []*gltf.Buffer{{ByteLength: 10}},
		}, args{[][4]uint8{{1, 2, 3, 4}}}, 1, &gltf.Document{
			Accessors: []*gltf.Accessor{
				{},
				{BufferView: gltf.Index(0), Count: 1, Type: gltf.AccessorVec4, ComponentType: gltf.ComponentUbyte, Normalized: true},
			},
			BufferViews: []*gltf.BufferView{
				{ByteLength: 4, Target: gltf.TargetArrayBuffer},
			},
			Buffers: []*gltf.Buffer{
				{ByteLength: 14, Data: []byte{1, 2, 3, 4}},
			},
		}},
		{"uint16", &gltf.Document{
			Accessors: []*gltf.Accessor{{}},
			Buffers:   []*gltf.Buffer{{ByteLength: 10}},
		}, args{[][4]uint16{{1, 2, 3, 4}}}, 1, &gltf.Document{
			Accessors: []*gltf.Accessor{
				{},
				{BufferView: gltf.Index(0), Count: 1, Type: gltf.AccessorVec4, ComponentType: gltf.ComponentUshort, Normalized: true},
			},
			BufferViews: []*gltf.BufferView{
				{ByteLength: 8, Target: gltf.TargetArrayBuffer},
			},
			Buffers: []*gltf.Buffer{
				{ByteLength: 18, Data: []byte{1, 0, 2, 0, 3, 0, 4, 0}},
			},
		}},
		{"float", &gltf.Document{
			Accessors: []*gltf.Accessor{{}},
			Buffers:   []*gltf.Buffer{{ByteLength: 10}},
		}, args{[][4]float32{{1, 2, 3, 4}, {}}}, 1, &gltf.Document{
			Accessors: []*gltf.Accessor{
				{},
				{BufferView: gltf.Index(0), Count: 2, Type: gltf.AccessorVec4, ComponentType: gltf.ComponentFloat},
			},
			BufferViews: []*gltf.BufferView{
				{ByteLength: 32, Target: gltf.TargetArrayBuffer},
			},
			Buffers: []*gltf.Buffer{
				{ByteLength: 42, Data: []byte{0, 0, 128, 63, 0, 0, 0, 64, 0, 0, 64, 64, 0, 0, 128, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}},
			},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WriteWeights(tt.m, tt.args.data)
			if tt.want != got {
				t.Errorf("WriteWeights() = %v, want %v", got, tt.want)
				return
			}
			if diff := deep.Equal(tt.m, tt.wantDoc); diff != nil {
				t.Errorf("WriteWeights() = %v", diff)
				return
			}
		})
	}
}

func TestWriteTextureCoord(t *testing.T) {
	type args struct {
		data any
	}
	tests := []struct {
		name    string
		m       *gltf.Document
		args    args
		want    uint32
		wantDoc *gltf.Document
	}{
		{"uint8", &gltf.Document{
			Accessors: []*gltf.Accessor{{}},
			Buffers:   []*gltf.Buffer{{ByteLength: 10}},
		}, args{[][2]uint8{{1, 2}}}, 1, &gltf.Document{
			Accessors: []*gltf.Accessor{
				{},
				{BufferView: gltf.Index(0), Count: 1, Type: gltf.AccessorVec2, ComponentType: gltf.ComponentUbyte, Normalized: true},
			},
			BufferViews: []*gltf.BufferView{
				{ByteLength: 4, Target: gltf.TargetArrayBuffer, ByteStride: 4},
			},
			Buffers: []*gltf.Buffer{
				{ByteLength: 14, Data: []byte{1, 2, 0, 0}},
			},
		}},
		{"uint16", &gltf.Document{
			Accessors: []*gltf.Accessor{{}},
			Buffers:   []*gltf.Buffer{{ByteLength: 10}},
		}, args{[][2]uint16{{1, 2}}}, 1, &gltf.Document{
			Accessors: []*gltf.Accessor{
				{},
				{BufferView: gltf.Index(0), Count: 1, Type: gltf.AccessorVec2, ComponentType: gltf.ComponentUshort, Normalized: true},
			},
			BufferViews: []*gltf.BufferView{
				{ByteLength: 4, Target: gltf.TargetArrayBuffer},
			},
			Buffers: []*gltf.Buffer{
				{ByteLength: 14, Data: []byte{1, 0, 2, 0}},
			},
		}},
		{"float", &gltf.Document{
			Accessors: []*gltf.Accessor{{}},
			Buffers:   []*gltf.Buffer{{ByteLength: 10}},
		}, args{[][2]float32{{1, 2}, {}}}, 1, &gltf.Document{
			Accessors: []*gltf.Accessor{
				{},
				{BufferView: gltf.Index(0), Count: 2, Type: gltf.AccessorVec2, ComponentType: gltf.ComponentFloat},
			},
			BufferViews: []*gltf.BufferView{
				{ByteLength: 16, Target: gltf.TargetArrayBuffer, Buffer: 0},
			},
			Buffers: []*gltf.Buffer{
				{ByteLength: 26, Data: []byte{0, 0, 128, 63, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0}},
			},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WriteTextureCoord(tt.m, tt.args.data)
			if tt.want != got {
				t.Errorf("WriteTextureCoord() = %v, want %v", got, tt.want)
				return
			}
			if diff := deep.Equal(tt.m, tt.wantDoc); diff != nil {
				t.Errorf("WriteTextureCoord() = %v", diff)
				return
			}
		})
	}
}

func TestWriteIndices(t *testing.T) {
	type args struct {
		data any
	}
	tests := []struct {
		name    string
		m       *gltf.Document
		args    args
		want    uint32
		wantDoc *gltf.Document
	}{
		{"uint16", &gltf.Document{
			Accessors: []*gltf.Accessor{{}},
			Buffers:   []*gltf.Buffer{{ByteLength: 10}},
		}, args{[]uint16{1, 2}}, 1, &gltf.Document{
			Accessors: []*gltf.Accessor{
				{},
				{BufferView: gltf.Index(0), Count: 2, Type: gltf.AccessorScalar, ComponentType: gltf.ComponentUshort},
			},
			BufferViews: []*gltf.BufferView{
				{ByteLength: 4, Target: gltf.TargetElementArrayBuffer},
			},
			Buffers: []*gltf.Buffer{
				{ByteLength: 14, Data: []byte{1, 0, 2, 0}},
			},
		}},
		{"uint32", &gltf.Document{
			Accessors: []*gltf.Accessor{{}},
			Buffers:   []*gltf.Buffer{{ByteLength: 10}},
		}, args{[]uint32{1, 2}}, 1, &gltf.Document{
			Accessors: []*gltf.Accessor{
				{},
				{BufferView: gltf.Index(0), Count: 2, Type: gltf.AccessorScalar, ComponentType: gltf.ComponentUint},
			},
			BufferViews: []*gltf.BufferView{
				{ByteLength: 8, Target: gltf.TargetElementArrayBuffer},
			},
			Buffers: []*gltf.Buffer{
				{ByteLength: 18, Data: []byte{1, 0, 0, 0, 2, 0, 0, 0}},
			},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WriteIndices(tt.m, tt.args.data)
			if tt.want != got {
				t.Errorf("WriteIndices() = %v, want %v", got, tt.want)
				return
			}
			if diff := deep.Equal(tt.m, tt.wantDoc); diff != nil {
				t.Errorf("WriteIndices() = %v", diff)
				return
			}
		})
	}
}

func TestWriteColor(t *testing.T) {
	type args struct {
		data any
	}
	tests := []struct {
		name    string
		m       *gltf.Document
		args    args
		want    uint32
		wantDoc *gltf.Document
	}{
		{"uint8", &gltf.Document{
			Accessors: []*gltf.Accessor{{}},
			Buffers:   []*gltf.Buffer{{ByteLength: 10}},
		}, args{[][4]uint8{{1, 2, 3, 4}}}, 1, &gltf.Document{
			Accessors: []*gltf.Accessor{
				{},
				{BufferView: gltf.Index(0), Count: 1, Type: gltf.AccessorVec4, ComponentType: gltf.ComponentUbyte, Normalized: true},
			},
			BufferViews: []*gltf.BufferView{
				{ByteLength: 4, Target: gltf.TargetArrayBuffer},
			},
			Buffers: []*gltf.Buffer{
				{ByteLength: 14, Data: []byte{1, 2, 3, 4}},
			},
		}},
		{"uint16", &gltf.Document{
			Accessors: []*gltf.Accessor{{}},
			Buffers:   []*gltf.Buffer{{ByteLength: 10}},
		}, args{[][4]uint16{{1, 2, 3, 4}}}, 1, &gltf.Document{
			Accessors: []*gltf.Accessor{
				{},
				{BufferView: gltf.Index(0), Count: 1, Type: gltf.AccessorVec4, ComponentType: gltf.ComponentUshort, Normalized: true},
			},
			BufferViews: []*gltf.BufferView{
				{ByteLength: 8, Target: gltf.TargetArrayBuffer},
			},
			Buffers: []*gltf.Buffer{
				{ByteLength: 18, Data: []byte{1, 0, 2, 0, 3, 0, 4, 0}},
			},
		}},
		{"float", &gltf.Document{
			Accessors: []*gltf.Accessor{{}},
		}, args{[][4]float32{{1, 2, 3, 4}, {}}}, 1, &gltf.Document{
			Accessors: []*gltf.Accessor{
				{},
				{BufferView: gltf.Index(0), Count: 2, Type: gltf.AccessorVec4, ComponentType: gltf.ComponentFloat},
			},
			BufferViews: []*gltf.BufferView{
				{ByteLength: 32, Target: gltf.TargetArrayBuffer},
			},
			Buffers: []*gltf.Buffer{
				{ByteLength: 32, Data: []byte{0, 0, 128, 63, 0, 0, 0, 64, 0, 0, 64, 64, 0, 0, 128, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}},
			},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WriteColor(tt.m, tt.args.data)
			if tt.want != got {
				t.Errorf("WriteColor() = %v, want %v", got, tt.want)
				return
			}
			if diff := deep.Equal(tt.m, tt.wantDoc); diff != nil {
				t.Errorf("WriteColor() = %v", diff)
				return
			}
		})
	}
}

func TestWriteImage(t *testing.T) {
	type args struct {
		name     string
		mimeType string
		r        io.Reader
	}
	tests := []struct {
		name    string
		m       *gltf.Document
		args    args
		want    uint32
		wantDoc *gltf.Document
		wantErr bool
	}{
		{"base", &gltf.Document{
			Images:  []*gltf.Image{{}},
			Buffers: []*gltf.Buffer{{ByteLength: 10, Data: make([]byte, 10)}},
		}, args{"fake", "fake/type", bytes.NewReader([]byte{1, 2})}, 1, &gltf.Document{
			Images: []*gltf.Image{
				{},
				{BufferView: gltf.Index(0), Name: "fake", MimeType: "fake/type"},
			},
			BufferViews: []*gltf.BufferView{
				{ByteOffset: 10, ByteLength: 2, Target: gltf.TargetNone},
			},
			Buffers: []*gltf.Buffer{
				{ByteLength: 12, Data: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 2}},
			},
		}, false},
		{"buffer", &gltf.Document{
			Images:  []*gltf.Image{{}},
			Buffers: []*gltf.Buffer{{ByteLength: 10, Data: make([]byte, 10)}},
		}, args{"fake", "fake/type", bytes.NewBuffer([]byte{1, 2})}, 1, &gltf.Document{
			Images: []*gltf.Image{
				{},
				{BufferView: gltf.Index(0), Name: "fake", MimeType: "fake/type"},
			},
			BufferViews: []*gltf.BufferView{
				{ByteOffset: 10, ByteLength: 2, Target: gltf.TargetNone},
			},
			Buffers: []*gltf.Buffer{
				{ByteLength: 12, Data: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 2}},
			},
		}, false},
		{"err", &gltf.Document{
			Images:  []*gltf.Image{{}},
			Buffers: []*gltf.Buffer{{ByteLength: 10}},
		}, args{"fake", "fake/type", &errReader{}}, 0, &gltf.Document{
			Images: []*gltf.Image{
				{},
			},
			Buffers: []*gltf.Buffer{
				{ByteLength: 10},
			},
		}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := WriteImage(tt.m, tt.args.name, tt.args.mimeType, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("WriteImage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != got {
				t.Errorf("WriteImage() = %v, want %v", got, tt.want)
				return
			}
			if diff := deep.Equal(tt.m, tt.wantDoc); diff != nil {
				t.Errorf("WriteImage() = %v", diff)
				return
			}
		})
	}
}

type errReader struct{}

func (r *errReader) Read(p []byte) (int, error) {
	return 0, errors.New("")
}
