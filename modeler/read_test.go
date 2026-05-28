package modeler_test

import (
	"encoding/binary"
	"math"
	"reflect"
	"testing"

	"github.com/qmuntal/gltf"
	"github.com/qmuntal/gltf/modeler"
)

func float32Bytes(values ...float32) []byte {
	data := make([]byte, 4*len(values))
	for i, value := range values {
		binary.LittleEndian.PutUint32(data[i*4:], math.Float32bits(value))
	}
	return data
}

func TestReadColor64_ReuseBuffer(t *testing.T) {
	data := []byte{1, 2, 3, 4, 1, 2, 3, 4}
	acc1 := &gltf.Accessor{
		BufferView: gltf.Index(0), Count: 1, Type: gltf.AccessorVec4, ComponentType: gltf.ComponentUbyte, Normalized: true,
	}
	acc2 := &gltf.Accessor{
		BufferView: gltf.Index(0), Count: 2, Type: gltf.AccessorVec4, ComponentType: gltf.ComponentUbyte, Normalized: true,
	}
	doc := &gltf.Document{
		BufferViews: []*gltf.BufferView{
			{Buffer: 0, ByteLength: len(data)},
		},
		Buffers: []*gltf.Buffer{
			{Data: data, ByteLength: len(data)},
		},
	}
	var buf [][4]uint16
	var err error
	buf, err = modeler.ReadColor64(doc, acc2, buf)
	if err != nil {
		t.Error(err)
	}
	if len(buf) != int(acc2.Count) {
		t.Errorf("ReadColor() len = %d, want %d", len(buf), acc2.Count)
	}
	buf, err = modeler.ReadColor64(doc, acc1, buf)
	if err != nil {
		t.Error(err)
	}
	if len(buf) != int(acc1.Count) {
		t.Errorf("ReadColor() len = %d, want %d", len(buf), acc1.Count)
	}
}

func TestReadBufferView(t *testing.T) {
	type args struct {
		doc *gltf.Document
		bv  *gltf.BufferView
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{"base", args{&gltf.Document{Buffers: []*gltf.Buffer{
			{ByteLength: 9, Data: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9}},
		}}, &gltf.BufferView{
			Buffer: 0, ByteLength: 3, ByteOffset: 6,
		}}, []byte{7, 8, 9}, false},
		{"errbuffer", args{&gltf.Document{Buffers: []*gltf.Buffer{
			{ByteLength: 9, Data: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9}},
		}}, &gltf.BufferView{
			Buffer: 1, ByteLength: 3, ByteOffset: 6,
		}}, nil, true},
		{"shortbuffer", args{&gltf.Document{Buffers: []*gltf.Buffer{
			{ByteLength: 9, Data: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9}},
		}}, &gltf.BufferView{
			Buffer: 0, ByteLength: 10, ByteOffset: 6,
		}}, nil, true},
		{"nil-doc", args{nil, &gltf.BufferView{}}, nil, true},
		{"nil-bufferview", args{&gltf.Document{}, nil}, nil, true},
		{"negative-buffer", args{&gltf.Document{Buffers: []*gltf.Buffer{
			{ByteLength: 9, Data: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9}},
		}}, &gltf.BufferView{
			Buffer: -1, ByteLength: 3,
		}}, nil, true},
		{"negative-offset", args{&gltf.Document{Buffers: []*gltf.Buffer{
			{ByteLength: 9, Data: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9}},
		}}, &gltf.BufferView{
			Buffer: 0, ByteLength: 3, ByteOffset: -1,
		}}, nil, true},
		{"negative-length", args{&gltf.Document{Buffers: []*gltf.Buffer{
			{ByteLength: 9, Data: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9}},
		}}, &gltf.BufferView{
			Buffer: 0, ByteLength: -1,
		}}, nil, true},
		{"declared-buffer-short", args{&gltf.Document{Buffers: []*gltf.Buffer{
			{ByteLength: 2, Data: []byte{1, 2, 3, 4}},
		}}, &gltf.BufferView{
			Buffer: 0, ByteLength: 3,
		}}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := modeler.ReadBufferView(tt.args.doc, tt.args.bv)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadBufferView() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadBufferView() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadAccessor_Buffered(t *testing.T) {
	doc := &gltf.Document{Buffers: []*gltf.Buffer{
		{ByteLength: 9, Data: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9}},
	}, BufferViews: []*gltf.BufferView{{
		Buffer: 0, ByteLength: 6, ByteOffset: 3,
	}}}
	acr := &gltf.Accessor{
		BufferView: gltf.Index(0), ByteOffset: 3, ComponentType: gltf.ComponentUbyte, Type: gltf.AccessorScalar, Count: 3,
	}
	data, err := modeler.ReadAccessor(doc, acr, make([]byte, 100))
	if err != nil {
		t.Fatal(err)
	}
	if len(data.([]byte)) != int(acr.Count) {
		t.Errorf("ReadAccessor expecting length %v, got %v", len(data.([]byte)), acr.Count)
	}
}

func TestReadAccessor(t *testing.T) {
	type args struct {
		doc *gltf.Document
		acr *gltf.Accessor
	}
	tests := []struct {
		name    string
		args    args
		want    any
		wantErr bool
	}{
		{"nodata", args{&gltf.Document{}, &gltf.Accessor{}}, []float32{}, false},
		{"base", args{&gltf.Document{Buffers: []*gltf.Buffer{
			{ByteLength: 9, Data: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9}},
		}, BufferViews: []*gltf.BufferView{{
			Buffer: 0, ByteLength: 6, ByteOffset: 3,
		}}}, &gltf.Accessor{
			BufferView: gltf.Index(0), ByteOffset: 3, ComponentType: gltf.ComponentUbyte, Type: gltf.AccessorScalar, Count: 3,
		}}, []byte{7, 8, 9}, false},
		{"shortbuffer", args{&gltf.Document{Buffers: []*gltf.Buffer{
			{ByteLength: 9, Data: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9}},
		}, BufferViews: []*gltf.BufferView{{
			Buffer: 0, ByteLength: 3, ByteOffset: 3,
		}}}, &gltf.Accessor{
			BufferView: gltf.Index(0), ByteOffset: 3, ComponentType: gltf.ComponentUbyte, Type: gltf.AccessorScalar, Count: 3,
		}}, nil, true},
		{"viewoverflow", args{&gltf.Document{Buffers: []*gltf.Buffer{
			{ByteLength: 9, Data: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9}},
		}, BufferViews: []*gltf.BufferView{{
			Buffer: 0, ByteLength: 6, ByteOffset: 3,
		}}}, &gltf.Accessor{
			BufferView: gltf.Index(1), ByteOffset: 3, ComponentType: gltf.ComponentUbyte, Type: gltf.AccessorScalar, Count: 3,
		}}, nil, true},
		{"negative-view", args{&gltf.Document{Buffers: []*gltf.Buffer{
			{ByteLength: 9, Data: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9}},
		}, BufferViews: []*gltf.BufferView{{
			Buffer: 0, ByteLength: 6, ByteOffset: 3,
		}}}, &gltf.Accessor{
			BufferView: gltf.Index(-1), ComponentType: gltf.ComponentUbyte, Type: gltf.AccessorScalar, Count: 3,
		}}, nil, true},
		{"negative-offset", args{&gltf.Document{Buffers: []*gltf.Buffer{
			{ByteLength: 9, Data: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9}},
		}, BufferViews: []*gltf.BufferView{{
			Buffer: 0, ByteLength: 6, ByteOffset: 3,
		}}}, &gltf.Accessor{
			BufferView: gltf.Index(0), ByteOffset: -1, ComponentType: gltf.ComponentUbyte, Type: gltf.AccessorScalar, Count: 3,
		}}, nil, true},
		{"offset-overflow", args{&gltf.Document{Buffers: []*gltf.Buffer{
			{ByteLength: 9, Data: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9}},
		}, BufferViews: []*gltf.BufferView{{
			Buffer: 0, ByteLength: 6, ByteOffset: 3,
		}}}, &gltf.Accessor{
			BufferView: gltf.Index(0), ByteOffset: 7, ComponentType: gltf.ComponentUbyte, Type: gltf.AccessorScalar, Count: 0,
		}}, nil, true},
		{"negative-count", args{&gltf.Document{}, &gltf.Accessor{
			ComponentType: gltf.ComponentFloat, Type: gltf.AccessorScalar, Count: -1,
		}}, nil, true},
		{"invalid-stride", args{&gltf.Document{
			Buffers:     []*gltf.Buffer{{ByteLength: 24, Data: make([]byte, 24)}},
			BufferViews: []*gltf.BufferView{{Buffer: 0, ByteLength: 24, ByteStride: 4}},
		}, &gltf.Accessor{
			BufferView: gltf.Index(0), ComponentType: gltf.ComponentFloat, Type: gltf.AccessorVec3, Count: 2,
		}}, nil, true},
		{"nonconformant-stride", args{&gltf.Document{
			Buffers:     []*gltf.Buffer{{ByteLength: 3, Data: []byte{1, 0, 2}}},
			BufferViews: []*gltf.BufferView{{Buffer: 0, ByteLength: 3, ByteStride: 2}},
		}, &gltf.Accessor{
			BufferView: gltf.Index(0), ComponentType: gltf.ComponentUbyte, Type: gltf.AccessorScalar, Count: 2,
		}}, nil, true},
		{"interleaved", args{&gltf.Document{
			Buffers: []*gltf.Buffer{{ByteLength: 52, Data: []byte{
				0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 128, 63, 0, 0, 0, 64, 0, 0, 64, 64,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128, 191,
			}}}, BufferViews: []*gltf.BufferView{{Buffer: 0, ByteOffset: 4, ByteLength: 48, ByteStride: 24}},
		}, &gltf.Accessor{
			BufferView: gltf.Index(0), ByteOffset: 12, ComponentType: gltf.ComponentFloat, Type: gltf.AccessorVec3, Count: 2,
		}}, [][3]float32{{1, 2, 3}, {0, 0, -1}}, false},
		{"sparse", args{&gltf.Document{
			Buffers: []*gltf.Buffer{{ByteLength: 284, Data: []byte{
				0, 0, 8, 0, 7, 0, 0, 0, 1, 0, 8, 0, 1, 0, 9, 0, 8, 0, 1, 0, 2, 0, 9, 0,
				2, 0, 10, 0, 9, 0, 2, 0, 3, 0, 10, 0, 3, 0, 11, 0, 10, 0, 3, 0, 4, 0, 11,
				0, 4, 0, 12, 0, 11, 0, 4, 0, 5, 0, 12, 0, 5, 0, 13, 0, 12, 0, 5, 0, 6, 0,
				13, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128, 63, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64, 64, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 128, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 160, 64, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 192, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128, 63, 0,
				0, 0, 0, 0, 0, 128, 63, 0, 0, 128, 63, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 128, 63,
				0, 0, 0, 0, 0, 0, 64, 64, 0, 0, 128, 63, 0, 0, 0, 0, 0, 0, 128, 64, 0, 0, 128,
				63, 0, 0, 0, 0, 0, 0, 160, 64, 0, 0, 128, 63, 0, 0, 0, 0, 0, 0, 192, 64, 0, 0,
				128, 63, 0, 0, 0, 0, 8, 0, 10, 0, 12, 0, 0, 0, 0, 0, 128, 63, 0, 0, 0, 64, 0, 0,
				0, 0, 0, 0, 64, 64, 0, 0, 64, 64, 0, 0, 0, 0, 0, 0, 160, 64, 0, 0, 128, 64, 0, 0, 0, 0}}},
			BufferViews: []*gltf.BufferView{
				{Buffer: 0, ByteOffset: 72, ByteLength: 168},
				{Buffer: 0, ByteOffset: 240, ByteLength: 6},
				{Buffer: 0, ByteOffset: 248, ByteLength: 36},
			},
		}, &gltf.Accessor{
			BufferView: gltf.Index(0), ComponentType: gltf.ComponentFloat, Type: gltf.AccessorVec3, Count: 14,
			Sparse: &gltf.Sparse{
				Count:   3,
				Indices: gltf.SparseIndices{BufferView: 1, ComponentType: gltf.ComponentUshort},
				Values:  gltf.SparseValues{BufferView: 2},
			},
		}}, [][3]float32{
			{0, 0, 0}, {1, 0, 0}, {2, 0, 0}, {3, 0, 0}, {4, 0, 0}, {5, 0, 0}, {6, 0, 0},
			{0, 1, 0}, {1, 2, 0}, {2, 1, 0}, {3, 3, 0}, {4, 1, 0}, {5, 4, 0}, {6, 1, 0}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := modeler.ReadAccessor(tt.args.doc, tt.args.acr, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadAccessor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadAccessor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadAccessorSparseMalformed(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		acr  *gltf.Accessor
	}{
		{"index-overflow", append([]byte{2}, float32Bytes(1)...), &gltf.Accessor{
			ComponentType: gltf.ComponentFloat,
			Type:          gltf.AccessorScalar,
			Count:         2,
			Sparse: &gltf.Sparse{Count: 1,
				Indices: gltf.SparseIndices{BufferView: 0, ComponentType: gltf.ComponentUbyte},
				Values:  gltf.SparseValues{BufferView: 1},
			},
		}},
		{"duplicate-indices", append([]byte{1, 1}, float32Bytes(1, 2)...), &gltf.Accessor{
			ComponentType: gltf.ComponentFloat,
			Type:          gltf.AccessorScalar,
			Count:         3,
			Sparse: &gltf.Sparse{Count: 2,
				Indices: gltf.SparseIndices{BufferView: 0, ComponentType: gltf.ComponentUbyte},
				Values:  gltf.SparseValues{BufferView: 1},
			},
		}},
		{"out-of-order-indices", append([]byte{2, 1}, float32Bytes(1, 2)...), &gltf.Accessor{
			ComponentType: gltf.ComponentFloat,
			Type:          gltf.AccessorScalar,
			Count:         3,
			Sparse: &gltf.Sparse{Count: 2,
				Indices: gltf.SparseIndices{BufferView: 0, ComponentType: gltf.ComponentUbyte},
				Values:  gltf.SparseValues{BufferView: 1},
			},
		}},
		{"count-overflow", nil, &gltf.Accessor{
			ComponentType: gltf.ComponentFloat,
			Type:          gltf.AccessorScalar,
			Count:         1,
			Sparse: &gltf.Sparse{Count: 2,
				Indices: gltf.SparseIndices{BufferView: 0, ComponentType: gltf.ComponentUbyte},
				Values:  gltf.SparseValues{BufferView: 1},
			},
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := &gltf.Document{}
			if tt.data != nil {
				doc.Buffers = []*gltf.Buffer{{Data: tt.data, ByteLength: len(tt.data)}}
				doc.BufferViews = []*gltf.BufferView{
					{Buffer: 0, ByteLength: tt.acr.Sparse.Count},
					{Buffer: 0, ByteOffset: tt.acr.Sparse.Count, ByteLength: len(tt.data) - tt.acr.Sparse.Count},
				}
			}
			if _, err := modeler.ReadAccessor(doc, tt.acr, nil); err == nil {
				t.Fatal("ReadAccessor() expected an error")
			}
		})
	}
}

func TestReadAccessorAllocs(t *testing.T) {
	doc := &gltf.Document{
		Buffers: []*gltf.Buffer{{ByteLength: 52, Data: []byte{
			0, 0, 128, 63, 0, 0, 0, 64, 0, 0, 64, 64,
			0, 0, 128, 63, 0, 0, 0, 64, 0, 0, 64, 64,
			0, 0, 128, 63, 0, 0, 0, 64, 0, 0, 64, 64,
			0, 0, 128, 63, 0, 0, 0, 64, 0, 0, 64, 64,
		}}}, BufferViews: []*gltf.BufferView{{Buffer: 0, ByteLength: 48}},
	}
	acr := &gltf.Accessor{
		BufferView: gltf.Index(0), ComponentType: gltf.ComponentFloat, Type: gltf.AccessorVec3, Count: 4,
	}

	testFunc := func(t *testing.T, buf []byte, want float32) {
		allocs := testing.AllocsPerRun(50, func() {
			modeler.ReadAccessor(doc, acr, buf)
		})
		if allocs > float64(want) {
			t.Errorf("ReadAccessor expected %v allocs got %v", want, allocs)
		}

	}
	t.Run("nil", func(t *testing.T) {
		testFunc(t, nil, 2)
	})
	t.Run("2", func(t *testing.T) {
		buf := make([]byte, 24)
		testFunc(t, buf, 3)
		testFunc(t, buf, 3)
		testFunc(t, buf, 3)
		testFunc(t, buf, 3)
	})
	t.Run("4", func(t *testing.T) {
		buf := make([]byte, 48)
		testFunc(t, buf, 1)
		testFunc(t, buf, 1)
		testFunc(t, buf, 1)
		testFunc(t, buf, 1)
	})
}

func TestReadIndices(t *testing.T) {
	type args struct {
		data   []byte
		acr    *gltf.Accessor
		buffer []uint32
	}
	tests := []struct {
		name    string
		args    args
		want    []uint32
		wantErr bool
	}{
		{"uint8", args{[]byte{1, 2}, &gltf.Accessor{
			BufferView: gltf.Index(0), Count: 2, Type: gltf.AccessorScalar, ComponentType: gltf.ComponentUbyte,
		}, nil}, []uint32{1, 2}, false},
		{"uint16", args{[]byte{1, 0, 2, 0}, &gltf.Accessor{
			BufferView: gltf.Index(0), Count: 2, Type: gltf.AccessorScalar, ComponentType: gltf.ComponentUshort,
		}, nil}, []uint32{1, 2}, false},
		{"uint32", args{[]byte{1, 0, 0, 0, 2, 0, 0, 0}, &gltf.Accessor{
			BufferView: gltf.Index(0), Count: 2, Type: gltf.AccessorScalar, ComponentType: gltf.ComponentUint,
		}, nil}, []uint32{1, 2}, false},
		{"uint32-withbuffer", args{[]byte{1, 0, 0, 0, 2, 0, 0, 0}, &gltf.Accessor{
			BufferView: gltf.Index(0), Count: 2, Type: gltf.AccessorScalar, ComponentType: gltf.ComponentUint,
		}, make([]uint32, 1)}, []uint32{1, 2}, false},
		{"incorrect-type", args{[]byte{}, &gltf.Accessor{
			BufferView: gltf.Index(0), Type: gltf.AccessorMat2, ComponentType: gltf.ComponentUint,
		}, nil}, nil, true},
		{"incorrect-componenttype", args{[]byte{}, &gltf.Accessor{
			BufferView: gltf.Index(0), Type: gltf.AccessorScalar, ComponentType: gltf.ComponentFloat,
		}, nil}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := &gltf.Document{
				BufferViews: []*gltf.BufferView{
					{Buffer: 0, ByteLength: len(tt.args.data)},
				},
				Buffers: []*gltf.Buffer{
					{Data: tt.args.data, ByteLength: len(tt.args.data)},
				},
			}
			got, err := modeler.ReadIndices(doc, tt.args.acr, tt.args.buffer)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadIndices() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadIndices() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadNormal(t *testing.T) {
	type args struct {
		data   []byte
		acr    *gltf.Accessor
		buffer [][3]float32
	}
	tests := []struct {
		name    string
		args    args
		want    [][3]float32
		wantErr bool
	}{
		{"float32", args{[]byte{0, 0, 128, 63, 0, 0, 0, 64, 0, 0, 64, 64}, &gltf.Accessor{
			BufferView: gltf.Index(0), Count: 1, Type: gltf.AccessorVec3, ComponentType: gltf.ComponentFloat,
		}, nil}, [][3]float32{{1, 2, 3}}, false},
		{"float32-withbuffer", args{[]byte{0, 0, 128, 63, 0, 0, 0, 64, 0, 0, 64, 64}, &gltf.Accessor{
			BufferView: gltf.Index(0), Count: 1, Type: gltf.AccessorVec3, ComponentType: gltf.ComponentFloat,
		}, make([][3]float32, 1)}, [][3]float32{{1, 2, 3}}, false},
		{"incorrect-type", args{[]byte{}, &gltf.Accessor{
			BufferView: gltf.Index(0), Type: gltf.AccessorMat2, ComponentType: gltf.ComponentFloat,
		}, nil}, nil, true},
		{"incorrect-componenttype", args{[]byte{}, &gltf.Accessor{
			BufferView: gltf.Index(0), Type: gltf.AccessorVec3, ComponentType: gltf.ComponentByte,
		}, nil}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := &gltf.Document{
				BufferViews: []*gltf.BufferView{
					{Buffer: 0, ByteLength: len(tt.args.data)},
				},
				Buffers: []*gltf.Buffer{
					{Data: tt.args.data, ByteLength: len(tt.args.data)},
				},
			}
			got, err := modeler.ReadNormal(doc, tt.args.acr, tt.args.buffer)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadNormal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadNormal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadTangent(t *testing.T) {
	type args struct {
		data   []byte
		acr    *gltf.Accessor
		buffer [][4]float32
	}
	tests := []struct {
		name    string
		args    args
		want    [][4]float32
		wantErr bool
	}{
		{"float32", args{[]byte{0, 0, 128, 63, 0, 0, 0, 64, 0, 0, 64, 64, 0, 0, 128, 64, 0, 0, 0, 0, 0, 0}, &gltf.Accessor{
			BufferView: gltf.Index(0), Count: 1, Type: gltf.AccessorVec4, ComponentType: gltf.ComponentFloat,
		}, nil}, [][4]float32{{1, 2, 3, 4}}, false},
		{"float32-withbuffer", args{[]byte{0, 0, 128, 63, 0, 0, 0, 64, 0, 0, 64, 64, 0, 0, 128, 64, 0, 0, 0, 0, 0, 0}, &gltf.Accessor{
			BufferView: gltf.Index(0), Count: 1, Type: gltf.AccessorVec4, ComponentType: gltf.ComponentFloat,
		}, make([][4]float32, 1)}, [][4]float32{{1, 2, 3, 4}}, false},
		{"incorrect-type", args{[]byte{}, &gltf.Accessor{
			BufferView: gltf.Index(0), Type: gltf.AccessorMat2, ComponentType: gltf.ComponentFloat,
		}, nil}, nil, true},
		{"incorrect-componenttype", args{[]byte{}, &gltf.Accessor{
			BufferView: gltf.Index(0), Type: gltf.AccessorVec4, ComponentType: gltf.ComponentByte,
		}, nil}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := &gltf.Document{
				BufferViews: []*gltf.BufferView{
					{Buffer: 0, ByteLength: len(tt.args.data)},
				},
				Buffers: []*gltf.Buffer{
					{Data: tt.args.data, ByteLength: len(tt.args.data)},
				},
			}
			got, err := modeler.ReadTangent(doc, tt.args.acr, tt.args.buffer)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadTangent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadTangent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadTextureCoord(t *testing.T) {
	type args struct {
		data   []byte
		acr    *gltf.Accessor
		buffer [][2]float32
	}
	tests := []struct {
		name    string
		args    args
		want    [][2]float32
		wantErr bool
	}{
		{"uint8", args{[]byte{255, 0, 0, 0}, &gltf.Accessor{
			BufferView: gltf.Index(0), Count: 1, Type: gltf.AccessorVec2, ComponentType: gltf.ComponentUbyte, Normalized: true,
		}, nil}, [][2]float32{{1, 0}}, false},
		{"uint16", args{[]byte{255, 255, 0, 0}, &gltf.Accessor{
			BufferView: gltf.Index(0), Count: 1, Type: gltf.AccessorVec2, ComponentType: gltf.ComponentUshort, Normalized: true,
		}, nil}, [][2]float32{{1, 0}}, false},
		{"float32", args{[]byte{0, 0, 128, 63, 0, 0, 0, 64}, &gltf.Accessor{
			BufferView: gltf.Index(0), Count: 1, Type: gltf.AccessorVec2, ComponentType: gltf.ComponentFloat,
		}, nil}, [][2]float32{{1, 2}}, false},
		{"float32-withbuffer", args{[]byte{0, 0, 128, 63, 0, 0, 0, 64}, &gltf.Accessor{
			BufferView: gltf.Index(0), Count: 1, Type: gltf.AccessorVec2, ComponentType: gltf.ComponentFloat,
		}, make([][2]float32, 1)}, [][2]float32{{1, 2}}, false},
		{"incorrect-type", args{[]byte{}, &gltf.Accessor{
			BufferView: gltf.Index(0), Type: gltf.AccessorMat2, ComponentType: gltf.ComponentFloat,
		}, nil}, nil, true},
		{"incorrect-componenttype", args{[]byte{}, &gltf.Accessor{
			BufferView: gltf.Index(0), Type: gltf.AccessorVec2, ComponentType: gltf.ComponentByte,
		}, nil}, nil, true},
		{"integer-not-normalized", args{[]byte{255, 0}, &gltf.Accessor{
			BufferView: gltf.Index(0), Count: 1, Type: gltf.AccessorVec2, ComponentType: gltf.ComponentUbyte,
		}, nil}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := &gltf.Document{
				BufferViews: []*gltf.BufferView{
					{Buffer: 0, ByteLength: len(tt.args.data)},
				},
				Buffers: []*gltf.Buffer{
					{Data: tt.args.data, ByteLength: len(tt.args.data)},
				},
			}
			got, err := modeler.ReadTextureCoord(doc, tt.args.acr, tt.args.buffer)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadTextureCoord() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadTextureCoord() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadWeights(t *testing.T) {
	type args struct {
		data   []byte
		acr    *gltf.Accessor
		buffer [][4]float32
	}
	tests := []struct {
		name    string
		args    args
		want    [][4]float32
		wantErr bool
	}{
		{"uint8", args{[]byte{255, 0, 255, 0}, &gltf.Accessor{
			BufferView: gltf.Index(0), Count: 1, Type: gltf.AccessorVec4, ComponentType: gltf.ComponentUbyte, Normalized: true,
		}, nil}, [][4]float32{{1, 0, 1, 0}}, false},
		{"uint16", args{[]byte{0, 0, 255, 255, 0, 0, 255, 255}, &gltf.Accessor{
			BufferView: gltf.Index(0), Count: 1, Type: gltf.AccessorVec4, ComponentType: gltf.ComponentUshort, Normalized: true,
		}, nil}, [][4]float32{{0, 1, 0, 1}}, false},
		{"float32", args{[]byte{0, 0, 128, 63, 0, 0, 0, 64, 0, 0, 64, 64, 0, 0, 128, 64}, &gltf.Accessor{
			BufferView: gltf.Index(0), Count: 1, Type: gltf.AccessorVec4, ComponentType: gltf.ComponentFloat,
		}, nil}, [][4]float32{{1, 2, 3, 4}}, false},
		{"float32-withbuffer", args{[]byte{0, 0, 128, 63, 0, 0, 0, 64, 0, 0, 64, 64, 0, 0, 128, 64}, &gltf.Accessor{
			BufferView: gltf.Index(0), Count: 1, Type: gltf.AccessorVec4, ComponentType: gltf.ComponentFloat,
		}, make([][4]float32, 1)}, [][4]float32{{1, 2, 3, 4}}, false},
		{"incorrect-type", args{[]byte{}, &gltf.Accessor{
			BufferView: gltf.Index(0), Type: gltf.AccessorMat2, ComponentType: gltf.ComponentFloat,
		}, nil}, nil, true},
		{"incorrect-componenttype", args{[]byte{}, &gltf.Accessor{
			BufferView: gltf.Index(0), Type: gltf.AccessorVec4, ComponentType: gltf.ComponentByte,
		}, nil}, nil, true},
		{"integer-not-normalized", args{[]byte{255, 0, 255, 0}, &gltf.Accessor{
			BufferView: gltf.Index(0), Count: 1, Type: gltf.AccessorVec4, ComponentType: gltf.ComponentUbyte,
		}, nil}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := &gltf.Document{
				BufferViews: []*gltf.BufferView{
					{Buffer: 0, ByteLength: len(tt.args.data)},
				},
				Buffers: []*gltf.Buffer{
					{Data: tt.args.data, ByteLength: len(tt.args.data)},
				},
			}
			got, err := modeler.ReadWeights(doc, tt.args.acr, tt.args.buffer)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadWeights() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadWeights() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadJoints(t *testing.T) {
	type args struct {
		data   []byte
		acr    *gltf.Accessor
		buffer [][4]uint16
	}
	tests := []struct {
		name    string
		args    args
		want    [][4]uint16
		wantErr bool
	}{
		{"uint8", args{[]byte{255, 0, 255, 0}, &gltf.Accessor{
			BufferView: gltf.Index(0), Count: 1, Type: gltf.AccessorVec4, ComponentType: gltf.ComponentUbyte,
		}, nil}, [][4]uint16{{255, 0, 255, 0}}, false},
		{"uint16", args{[]byte{0, 0, 255, 255, 0, 0, 255, 255}, &gltf.Accessor{
			BufferView: gltf.Index(0), Count: 1, Type: gltf.AccessorVec4, ComponentType: gltf.ComponentUshort,
		}, nil}, [][4]uint16{{0, 65535, 0, 65535}}, false},
		{"incorrect-type", args{[]byte{}, &gltf.Accessor{
			BufferView: gltf.Index(0), Type: gltf.AccessorMat2, ComponentType: gltf.ComponentUshort,
		}, nil}, nil, true},
		{"incorrect-componenttype", args{[]byte{}, &gltf.Accessor{
			BufferView: gltf.Index(0), Type: gltf.AccessorVec4, ComponentType: gltf.ComponentByte,
		}, nil}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := &gltf.Document{
				BufferViews: []*gltf.BufferView{
					{Buffer: 0, ByteLength: len(tt.args.data)},
				},
				Buffers: []*gltf.Buffer{
					{Data: tt.args.data, ByteLength: len(tt.args.data)},
				},
			}
			got, err := modeler.ReadJoints(doc, tt.args.acr, tt.args.buffer)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadJoints() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadJoints() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadPosition(t *testing.T) {
	type args struct {
		data   []byte
		acr    *gltf.Accessor
		buffer [][3]float32
	}
	tests := []struct {
		name    string
		args    args
		want    [][3]float32
		wantErr bool
	}{
		{"nil-bufferView", args{nil, &gltf.Accessor{
			Type: gltf.AccessorVec3, ComponentType: gltf.ComponentFloat,
		}, nil}, [][3]float32{}, false},
		{"float32", args{[]byte{0, 0, 128, 63, 0, 0, 0, 64, 0, 0, 64, 64}, &gltf.Accessor{
			BufferView: gltf.Index(0), Count: 1, Type: gltf.AccessorVec3, ComponentType: gltf.ComponentFloat,
		}, nil}, [][3]float32{{1, 2, 3}}, false},
		{"float32-withbuffer", args{[]byte{0, 0, 128, 63, 0, 0, 0, 64, 0, 0, 64, 64}, &gltf.Accessor{
			BufferView: gltf.Index(0), Count: 1, Type: gltf.AccessorVec3, ComponentType: gltf.ComponentFloat,
		}, make([][3]float32, 1)}, [][3]float32{{1, 2, 3}}, false},
		{"incorrect-type", args{[]byte{}, &gltf.Accessor{
			BufferView: gltf.Index(0), Type: gltf.AccessorMat2, ComponentType: gltf.ComponentFloat,
		}, nil}, nil, true},
		{"incorrect-componenttype", args{[]byte{}, &gltf.Accessor{
			BufferView: gltf.Index(0), Type: gltf.AccessorVec3, ComponentType: gltf.ComponentByte,
		}, nil}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := new(gltf.Document)
			if tt.args.data != nil {
				doc.BufferViews = []*gltf.BufferView{
					{Buffer: 0, ByteLength: len(tt.args.data)},
				}
				doc.Buffers = []*gltf.Buffer{
					{Data: tt.args.data, ByteLength: len(tt.args.data)},
				}
			}
			got, err := modeler.ReadPosition(doc, tt.args.acr, tt.args.buffer)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadPosition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadPosition() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadPosition_AfterReadNormalSparseNoBufferView(t *testing.T) {
	positionValues := float32Bytes(1, 2, 3)
	normalValues := float32Bytes(4, 5, 6, 7, 8, 9)
	data := append([]byte{1, 0}, positionValues...)
	data = append(data, normalValues...)

	doc := &gltf.Document{
		BufferViews: []*gltf.BufferView{
			{Buffer: 0, ByteLength: 2},
			{Buffer: 0, ByteOffset: 2, ByteLength: len(positionValues)},
			{Buffer: 0, ByteOffset: 2 + len(positionValues), ByteLength: len(normalValues)},
		},
		Buffers: []*gltf.Buffer{
			{Data: data, ByteLength: len(data)},
		},
	}
	positionAccessor := &gltf.Accessor{
		ComponentType: gltf.ComponentFloat,
		Type:          gltf.AccessorVec3,
		Count:         2,
		Sparse: &gltf.Sparse{
			Count: 1,
			Indices: gltf.SparseIndices{
				BufferView:    0,
				ComponentType: gltf.ComponentUshort,
			},
			Values: gltf.SparseValues{BufferView: 1},
		},
	}
	normalAccessor := &gltf.Accessor{
		BufferView:    gltf.Index(2),
		ComponentType: gltf.ComponentFloat,
		Type:          gltf.AccessorVec3,
		Count:         2,
	}
	wantPosition := [][3]float32{{0, 0, 0}, {1, 2, 3}}

	gotPosition, err := modeler.ReadPosition(doc, positionAccessor, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(gotPosition, wantPosition) {
		t.Errorf("ReadPosition() = %v, want %v", gotPosition, wantPosition)
	}
	if _, err := modeler.ReadNormal(doc, normalAccessor, nil); err != nil {
		t.Fatal(err)
	}
	gotPosition, err = modeler.ReadPosition(doc, positionAccessor, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(gotPosition, wantPosition) {
		t.Errorf("ReadPosition() after ReadNormal() = %v, want %v", gotPosition, wantPosition)
	}
}

func TestReadColor(t *testing.T) {
	type args struct {
		data   []byte
		acr    *gltf.Accessor
		buffer [][4]uint8
	}
	tests := []struct {
		name    string
		args    args
		want    [][4]uint8
		wantErr bool
	}{
		{"[4]uint8", args{[]byte{1, 2, 3, 4}, &gltf.Accessor{
			BufferView: gltf.Index(0), Count: 1, Type: gltf.AccessorVec4, ComponentType: gltf.ComponentUbyte, Normalized: true,
		}, nil}, [][4]uint8{{1, 2, 3, 4}}, false},
		{"[4]uint16", args{[]byte{0, 0, 115, 33, 200, 0, 255, 255}, &gltf.Accessor{
			BufferView: gltf.Index(0), Count: 1, Type: gltf.AccessorVec4, ComponentType: gltf.ComponentUshort, Normalized: true,
		}, nil}, [][4]uint8{{0, 33, 1, 255}}, false},
		{"[4]uint16-normalized", args{[]byte{0, 1, 115, 33, 200, 0, 255, 255}, &gltf.Accessor{
			BufferView: gltf.Index(0), Count: 1, Type: gltf.AccessorVec4, ComponentType: gltf.ComponentUshort, Normalized: true,
		}, nil}, [][4]uint8{{1, 33, 1, 255}}, false},
		{"[4]uint16-normalized-linear", args{[]byte{37, 3, 180, 50, 255, 255, 255, 255}, &gltf.Accessor{
			BufferView: gltf.Index(0), Count: 1, Type: gltf.AccessorVec4, ComponentType: gltf.ComponentUshort, Normalized: true,
		}, nil}, [][4]uint8{{3, 51, 255, 255}}, false},
		{"[4]float32", args{float32Bytes(0, 0.5, 1, 1), &gltf.Accessor{
			BufferView: gltf.Index(0), Count: 1, Type: gltf.AccessorVec4, ComponentType: gltf.ComponentFloat,
		}, nil}, [][4]uint8{{0, 128, 255, 255}}, false},
		{"[3]uint8", args{[]byte{1, 2, 3, 0}, &gltf.Accessor{
			BufferView: gltf.Index(0), Count: 1, Type: gltf.AccessorVec3, ComponentType: gltf.ComponentUbyte, Normalized: true,
		}, nil}, [][4]uint8{{1, 2, 3, 255}}, false},
		{"[3]uint16", args{[]byte{0, 0, 255, 0, 255, 0, 0, 0}, &gltf.Accessor{
			BufferView: gltf.Index(0), Count: 1, Type: gltf.AccessorVec3, ComponentType: gltf.ComponentUshort, Normalized: true,
		}, nil}, [][4]uint8{{0, 1, 1, 255}}, false},
		{"[3]uint16-normalized", args{[]byte{0, 1, 255, 0, 255, 0, 0, 0}, &gltf.Accessor{
			BufferView: gltf.Index(0), Count: 1, Type: gltf.AccessorVec3, ComponentType: gltf.ComponentUshort, Normalized: true,
		}, nil}, [][4]uint8{{1, 1, 1, 255}}, false},
		{"[3]float32", args{float32Bytes(805.0/65535.0, 12980.0/65535.0, 1), &gltf.Accessor{
			BufferView: gltf.Index(0), Count: 1, Type: gltf.AccessorVec3, ComponentType: gltf.ComponentFloat,
		}, nil}, [][4]uint8{{3, 51, 255, 255}}, false},
		{"incorrect-type", args{[]byte{}, &gltf.Accessor{
			BufferView: gltf.Index(0), Type: gltf.AccessorMat2, ComponentType: gltf.ComponentUbyte,
		}, nil}, nil, true},
		{"incorrect-componenttype", args{[]byte{}, &gltf.Accessor{
			BufferView: gltf.Index(0), Type: gltf.AccessorVec4, ComponentType: gltf.ComponentByte,
		}, nil}, nil, true},
		{"integer-not-normalized", args{[]byte{1, 2, 3, 4}, &gltf.Accessor{
			BufferView: gltf.Index(0), Count: 1, Type: gltf.AccessorVec4, ComponentType: gltf.ComponentUbyte,
		}, nil}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := &gltf.Document{
				BufferViews: []*gltf.BufferView{
					{Buffer: 0, ByteLength: len(tt.args.data)},
				},
				Buffers: []*gltf.Buffer{
					{Data: tt.args.data, ByteLength: len(tt.args.data)},
				},
			}
			got, err := modeler.ReadColor(doc, tt.args.acr, tt.args.buffer)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadColor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadColor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadColor64(t *testing.T) {
	type args struct {
		data   []byte
		acr    *gltf.Accessor
		buffer [][4]uint16
	}
	tests := []struct {
		name    string
		args    args
		want    [][4]uint16
		wantErr bool
	}{
		{"[4]uint8", args{[]byte{0, 115, 200, 255}, &gltf.Accessor{
			BufferView: gltf.Index(0), Count: 1, Type: gltf.AccessorVec4, ComponentType: gltf.ComponentUbyte, Normalized: true,
		}, nil}, [][4]uint16{{0, 29555, 51400, 65535}}, false},
		{"[4]uint16", args{[]byte{0, 0, 255, 255, 0, 0, 255, 255}, &gltf.Accessor{
			BufferView: gltf.Index(0), Count: 1, Type: gltf.AccessorVec4, ComponentType: gltf.ComponentUshort, Normalized: true,
		}, nil}, [][4]uint16{{0, 65535, 0, 65535}}, false},
		{"[4]uint16-normalized", args{[]byte{37, 3, 180, 50, 255, 255, 255, 255}, &gltf.Accessor{
			BufferView: gltf.Index(0), Count: 1, Type: gltf.AccessorVec4, ComponentType: gltf.ComponentUshort, Normalized: true,
		}, nil}, [][4]uint16{{805, 12980, 65535, 65535}}, false},
		{"[4]float32", args{float32Bytes(0, 0.5, 1, 1), &gltf.Accessor{
			BufferView: gltf.Index(0), Count: 1, Type: gltf.AccessorVec4, ComponentType: gltf.ComponentFloat,
		}, nil}, [][4]uint16{{0, 32768, 65535, 65535}}, false},
		{"[3]uint8", args{[]byte{0, 100, 200, 0}, &gltf.Accessor{
			BufferView: gltf.Index(0), Count: 1, Type: gltf.AccessorVec3, ComponentType: gltf.ComponentUbyte, Normalized: true,
		}, nil}, [][4]uint16{{0, 25700, 51400, 65535}}, false},
		{"[3]uint16", args{[]byte{0, 0, 255, 0, 255, 0, 0, 0}, &gltf.Accessor{
			BufferView: gltf.Index(0), Count: 1, Type: gltf.AccessorVec3, ComponentType: gltf.ComponentUshort, Normalized: true,
		}, nil}, [][4]uint16{{0, 255, 255, 65535}}, false},
		{"[3]float32", args{float32Bytes(805.0/65535.0, 12980.0/65535.0, 1), &gltf.Accessor{
			BufferView: gltf.Index(0), Count: 1, Type: gltf.AccessorVec3, ComponentType: gltf.ComponentFloat,
		}, nil}, [][4]uint16{{805, 12980, 65535, 65535}}, false},
		{"incorrect-type", args{[]byte{}, &gltf.Accessor{
			BufferView: gltf.Index(0), Type: gltf.AccessorMat2, ComponentType: gltf.ComponentUbyte,
		}, nil}, nil, true},
		{"incorrect-componenttype", args{[]byte{}, &gltf.Accessor{
			BufferView: gltf.Index(0), Type: gltf.AccessorVec4, ComponentType: gltf.ComponentByte,
		}, nil}, nil, true},
		{"integer-not-normalized", args{[]byte{0, 115, 200, 255}, &gltf.Accessor{
			BufferView: gltf.Index(0), Count: 1, Type: gltf.AccessorVec4, ComponentType: gltf.ComponentUbyte,
		}, nil}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := &gltf.Document{
				BufferViews: []*gltf.BufferView{
					{Buffer: 0, ByteLength: len(tt.args.data)},
				},
				Buffers: []*gltf.Buffer{
					{Data: tt.args.data, ByteLength: len(tt.args.data)},
				},
			}
			got, err := modeler.ReadColor64(doc, tt.args.acr, tt.args.buffer)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadColor64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadColor64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadInverseBindMatrices(t *testing.T) {
	type args struct {
		data   []byte
		acr    *gltf.Accessor
		buffer [][4][4]float32
	}
	tests := []struct {
		name    string
		args    args
		want    [][4][4]float32
		wantErr bool
	}{
		{"base", args{[]byte{
			0, 0, 128, 63, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 64, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 128, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		}, &gltf.Accessor{BufferView: gltf.Index(0), Count: 1, Type: gltf.AccessorMat4, ComponentType: gltf.ComponentFloat}, nil},
			[][4][4]float32{{{1, 2, 3, 4}}}, false,
		},
		{"incorrect-type", args{[]byte{}, &gltf.Accessor{
			BufferView: gltf.Index(0), Type: gltf.AccessorMat2, ComponentType: gltf.ComponentFloat,
		}, nil}, nil, true},
		{"incorrect-componenttype", args{[]byte{}, &gltf.Accessor{
			BufferView: gltf.Index(0), Type: gltf.AccessorMat4, ComponentType: gltf.ComponentByte,
		}, nil}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := &gltf.Document{
				BufferViews: []*gltf.BufferView{
					{Buffer: 0, ByteLength: len(tt.args.data)},
				},
				Buffers: []*gltf.Buffer{
					{Data: tt.args.data, ByteLength: len(tt.args.data)},
				},
			}
			got, err := modeler.ReadInverseBindMatrices(doc, tt.args.acr, tt.args.buffer)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadInverseBindMatrices() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadInverseBindMatrices() = %v, want %v", got, tt.want)
			}
		})
	}

}
