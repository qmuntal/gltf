package binary_test

import (
	"bytes"
	gobinary "encoding/binary"
	"testing"

	"github.com/qmuntal/gltf"
	"github.com/qmuntal/gltf/binary"
)

func BenchmarkNative(b *testing.B) {
	var s int = 1000
	bs := make([]byte, s*gltf.SizeOfElement(gltf.ComponentFloat, gltf.AccessorVec3))
	data := make([][3]float32, s)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for i := range data {
			binary.Float.PutVec3(bs[4*i:], data[i])
		}
	}
}

func BenchmarkWrite(b *testing.B) {
	var s int = 1000
	bs := make([]byte, s*gltf.SizeOfElement(gltf.ComponentFloat, gltf.AccessorVec3))
	data := make([][3]float32, s)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		binary.Write(bs, 0, data)
	}
}

func BenchmarkWrite_builtint(b *testing.B) {
	var s int = 1000
	bs := bytes.NewBuffer(make([]byte, s*gltf.SizeOfElement(gltf.ComponentFloat, gltf.AccessorVec3)))
	data := make([][3]float32, s)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		gobinary.Write(bs, gobinary.LittleEndian, data)
	}
}
