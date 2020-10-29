package modeler

import (
	"reflect"
	"testing"

	"github.com/qmuntal/gltf"
)

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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadBufferView(tt.args.doc, tt.args.bv)
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

func TestReadAccessor(t *testing.T) {
	type args struct {
		doc *gltf.Document
		acr *gltf.Accessor
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{"nodata", args{&gltf.Document{}, &gltf.Accessor{}}, nil, false},
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
		}}, []byte{7, 8, 9}, false},
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
			got, err := ReadAccessor(tt.args.doc, tt.args.acr, nil)
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

	testFunc := func(t *testing.T, buf [][3]float32, want float64) {
		allocs := testing.AllocsPerRun(10, func() {
			ReadAccessor(doc, acr, buf)
		})
		if allocs != want {
			t.Errorf("ReadAccessor expected %v allocs got %v", want, allocs)
		}

	}
	t.Run("nil", func(t *testing.T) {
		testFunc(t, nil, 2)
	})
	t.Run("2", func(t *testing.T) {
		buf := make([][3]float32, 2)
		testFunc(t, buf, 6)
		testFunc(t, buf, 6)
		testFunc(t, buf, 6)
		testFunc(t, buf, 6)
	})
	t.Run("4", func(t *testing.T) {
		buf := make([][3]float32, 4)
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
					{Buffer: 0, ByteLength: uint32(len(tt.args.data))},
				},
				Buffers: []*gltf.Buffer{
					{Data: tt.args.data, ByteLength: uint32(len(tt.args.data))},
				},
			}
			got, err := ReadIndices(doc, tt.args.acr, tt.args.buffer)
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
					{Buffer: 0, ByteLength: uint32(len(tt.args.data))},
				},
				Buffers: []*gltf.Buffer{
					{Data: tt.args.data, ByteLength: uint32(len(tt.args.data))},
				},
			}
			got, err := ReadNormal(doc, tt.args.acr, tt.args.buffer)
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
					{Buffer: 0, ByteLength: uint32(len(tt.args.data))},
				},
				Buffers: []*gltf.Buffer{
					{Data: tt.args.data, ByteLength: uint32(len(tt.args.data))},
				},
			}
			got, err := ReadTangent(doc, tt.args.acr, tt.args.buffer)
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
			BufferView: gltf.Index(0), Count: 1, Type: gltf.AccessorVec2, ComponentType: gltf.ComponentUbyte,
		}, nil}, [][2]float32{{1, 0}}, false},
		{"uint16", args{[]byte{255, 255, 0, 0}, &gltf.Accessor{
			BufferView: gltf.Index(0), Count: 1, Type: gltf.AccessorVec2, ComponentType: gltf.ComponentUshort,
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := &gltf.Document{
				BufferViews: []*gltf.BufferView{
					{Buffer: 0, ByteLength: uint32(len(tt.args.data))},
				},
				Buffers: []*gltf.Buffer{
					{Data: tt.args.data, ByteLength: uint32(len(tt.args.data))},
				},
			}
			got, err := ReadTextureCoord(doc, tt.args.acr, tt.args.buffer)
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
			BufferView: gltf.Index(0), Count: 1, Type: gltf.AccessorVec4, ComponentType: gltf.ComponentUbyte,
		}, nil}, [][4]float32{{1, 0, 1, 0}}, false},
		{"uint16", args{[]byte{0, 0, 255, 255, 0, 0, 255, 255}, &gltf.Accessor{
			BufferView: gltf.Index(0), Count: 1, Type: gltf.AccessorVec4, ComponentType: gltf.ComponentUshort,
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := &gltf.Document{
				BufferViews: []*gltf.BufferView{
					{Buffer: 0, ByteLength: uint32(len(tt.args.data))},
				},
				Buffers: []*gltf.Buffer{
					{Data: tt.args.data, ByteLength: uint32(len(tt.args.data))},
				},
			}
			got, err := ReadWeights(doc, tt.args.acr, tt.args.buffer)
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
					{Buffer: 0, ByteLength: uint32(len(tt.args.data))},
				},
				Buffers: []*gltf.Buffer{
					{Data: tt.args.data, ByteLength: uint32(len(tt.args.data))},
				},
			}
			got, err := ReadJoints(doc, tt.args.acr, tt.args.buffer)
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
					{Buffer: 0, ByteLength: uint32(len(tt.args.data))},
				},
				Buffers: []*gltf.Buffer{
					{Data: tt.args.data, ByteLength: uint32(len(tt.args.data))},
				},
			}
			got, err := ReadPosition(doc, tt.args.acr, tt.args.buffer)
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
			BufferView: gltf.Index(0), Count: 1, Type: gltf.AccessorVec4, ComponentType: gltf.ComponentUbyte,
		}, nil}, [][4]uint8{{1, 2, 3, 4}}, false},
		{"[4]uint16", args{[]byte{0, 0, 255, 255, 0, 0, 255, 255}, &gltf.Accessor{
			BufferView: gltf.Index(0), Count: 1, Type: gltf.AccessorVec4, ComponentType: gltf.ComponentUshort,
		}, nil}, [][4]uint8{{0, 255, 0, 255}}, false},
		{"[4]float32", args{[]byte{0, 0, 128, 63, 0, 0, 0, 64, 0, 0, 64, 64, 0, 0, 128, 64}, &gltf.Accessor{
			BufferView: gltf.Index(0), Count: 1, Type: gltf.AccessorVec4, ComponentType: gltf.ComponentFloat,
		}, nil}, [][4]uint8{{255, 89, 155, 252}}, false},
		{"[3]uint8", args{[]byte{1, 2, 3, 0}, &gltf.Accessor{
			BufferView: gltf.Index(0), Count: 1, Type: gltf.AccessorVec3, ComponentType: gltf.ComponentUbyte,
		}, nil}, [][4]uint8{{1, 2, 3, 255}}, false},
		{"[3]uint16", args{[]byte{0, 0, 255, 0, 255, 0, 0, 0}, &gltf.Accessor{
			BufferView: gltf.Index(0), Count: 1, Type: gltf.AccessorVec3, ComponentType: gltf.ComponentUshort,
		}, nil}, [][4]uint8{{0, 255, 255, 255}}, false},
		{"[3]float32", args{[]byte{0, 0, 128, 63, 0, 0, 0, 64, 0, 0, 64, 64}, &gltf.Accessor{
			BufferView: gltf.Index(0), Count: 1, Type: gltf.AccessorVec3, ComponentType: gltf.ComponentFloat,
		}, nil}, [][4]uint8{{255, 89, 155, 255}}, false},
		{"incorrect-type", args{[]byte{}, &gltf.Accessor{
			BufferView: gltf.Index(0), Type: gltf.AccessorMat2, ComponentType: gltf.ComponentUbyte,
		}, nil}, nil, true},
		{"incorrect-componenttype", args{[]byte{}, &gltf.Accessor{
			BufferView: gltf.Index(0), Type: gltf.AccessorVec4, ComponentType: gltf.ComponentByte,
		}, nil}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := &gltf.Document{
				BufferViews: []*gltf.BufferView{
					{Buffer: 0, ByteLength: uint32(len(tt.args.data))},
				},
				Buffers: []*gltf.Buffer{
					{Data: tt.args.data, ByteLength: uint32(len(tt.args.data))},
				},
			}
			got, err := ReadColor(doc, tt.args.acr, tt.args.buffer)
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
		{"[4]uint8", args{[]byte{1, 2, 3, 4}, &gltf.Accessor{
			BufferView: gltf.Index(0), Count: 1, Type: gltf.AccessorVec4, ComponentType: gltf.ComponentUbyte,
		}, nil}, [][4]uint16{{1, 2, 3, 4}}, false},
		{"[4]uint16", args{[]byte{0, 0, 255, 255, 0, 0, 255, 255}, &gltf.Accessor{
			BufferView: gltf.Index(0), Count: 1, Type: gltf.AccessorVec4, ComponentType: gltf.ComponentUshort,
		}, nil}, [][4]uint16{{0, 65535, 0, 65535}}, false},
		{"[4]float32", args{[]byte{0, 0, 128, 63, 0, 0, 0, 64, 0, 0, 64, 64, 0, 0, 128, 64}, &gltf.Accessor{
			BufferView: gltf.Index(0), Count: 1, Type: gltf.AccessorVec4, ComponentType: gltf.ComponentFloat,
		}, nil}, [][4]uint16{{65535, 23149, 40135, 65532}}, false},
		{"[3]uint8", args{[]byte{1, 2, 3, 0}, &gltf.Accessor{
			BufferView: gltf.Index(0), Count: 1, Type: gltf.AccessorVec3, ComponentType: gltf.ComponentUbyte,
		}, nil}, [][4]uint16{{1, 2, 3, 65535}}, false},
		{"[3]uint16", args{[]byte{0, 0, 255, 0, 255, 0, 0, 0}, &gltf.Accessor{
			BufferView: gltf.Index(0), Count: 1, Type: gltf.AccessorVec3, ComponentType: gltf.ComponentUshort,
		}, nil}, [][4]uint16{{0, 255, 255, 65535}}, false},
		{"[3]float32", args{[]byte{0, 0, 128, 63, 0, 0, 0, 64, 0, 0, 64, 64}, &gltf.Accessor{
			BufferView: gltf.Index(0), Count: 1, Type: gltf.AccessorVec3, ComponentType: gltf.ComponentFloat,
		}, nil}, [][4]uint16{{65535, 23149, 40135, 65535}}, false},
		{"incorrect-type", args{[]byte{}, &gltf.Accessor{
			BufferView: gltf.Index(0), Type: gltf.AccessorMat2, ComponentType: gltf.ComponentUbyte,
		}, nil}, nil, true},
		{"incorrect-componenttype", args{[]byte{}, &gltf.Accessor{
			BufferView: gltf.Index(0), Type: gltf.AccessorVec4, ComponentType: gltf.ComponentByte,
		}, nil}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := &gltf.Document{
				BufferViews: []*gltf.BufferView{
					{Buffer: 0, ByteLength: uint32(len(tt.args.data))},
				},
				Buffers: []*gltf.Buffer{
					{Data: tt.args.data, ByteLength: uint32(len(tt.args.data))},
				},
			}
			got, err := ReadColor64(doc, tt.args.acr, tt.args.buffer)
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
