package binary

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/qmuntal/gltf"
)

func BenchmarkNative(b *testing.B) {
	s := 1000
	bs := make([]byte, s*SizeOfElement(gltf.ComponentFloat, gltf.AccessorVec3))
	data := make([][3]float32, s)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for i := range data {
			Float.PutVec3(bs[4*i:], data[i])
		}
	}
}

func BenchmarkWrite(b *testing.B) {
	s := 1000
	bs := make([]byte, s*SizeOfElement(gltf.ComponentFloat, gltf.AccessorVec3))
	data := make([][3]float32, s)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		Write(bs, data)
	}
}

func BenchmarkWrite_builtint(b *testing.B) {
	s := 1000
	bs := bytes.NewBuffer(make([]byte, s*SizeOfElement(gltf.ComponentFloat, gltf.AccessorVec3)))
	data := make([][3]float32, s)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		binary.Write(bs, binary.LittleEndian, data)
	}
}
