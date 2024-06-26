package modeler_test

import (
	"testing"

	"github.com/qmuntal/gltf"
	"github.com/qmuntal/gltf/modeler"
)

func BenchmarkReadAccessorSparse(b *testing.B) {
	doc := &gltf.Document{
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
	}
	acr := &gltf.Accessor{
		BufferView: gltf.Index(0), ComponentType: gltf.ComponentFloat, Type: gltf.AccessorVec3, Count: 14,
		Sparse: &gltf.Sparse{
			Count:   3,
			Indices: gltf.SparseIndices{BufferView: 1, ComponentType: gltf.ComponentUshort},
			Values:  gltf.SparseValues{BufferView: 2},
		},
	}
	for i := 0; i < b.N; i++ {
		_, err := modeler.ReadAccessor(doc, acr, nil)
		if err != nil {
			b.Fatal(err)
		}
	}
}
