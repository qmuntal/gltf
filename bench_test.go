package gltf_test

import (
	"testing"

	"github.com/qmuntal/gltf"
)

func BenchmarkOpenASCII(b *testing.B) {
	benchs := []struct {
		name string
	}{
		{"testdata/AnimatedCube/glTF/AnimatedCube.gltf"},
		{"testdata/BoxVertexColors/glTF/BoxVertexColors.gltf"},
		{"testdata/Cameras/glTF/Cameras.gltf"},
		{"testdata/Cube/glTF/Cube.gltf"},
		{"testdata/EnvironmentTest/glTF/EnvironmentTest.gltf"},
		{"testdata/OrientationTest/glTF/OrientationTest.gltf"},
		{"testdata/Triangle/glTF/Triangle.gltf"},
		{"testdata/TriangleWithoutIndices/glTF/TriangleWithoutIndices.gltf"},
	}
	for _, bb := range benchs {
		b.Run(bb.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := gltf.Open(bb.name)
				if err != nil {
					b.Errorf("Open() error = %v", err)
					return
				}
			}
		})
	}
}

func BenchmarkOpenEmbedded(b *testing.B) {
	benchs := []struct {
		name string
	}{
		{"testdata/BoxVertexColors/glTF-Embedded/BoxVertexColors.gltf"},
		{"testdata/Cameras/glTF-Embedded/Cameras.gltf"},
		{"testdata/OrientationTest/glTF-Embedded/OrientationTest.gltf"},
		{"testdata/Triangle/glTF-Embedded/Triangle.gltf"},
		{"testdata/TriangleWithoutIndices/glTF-Embedded/TriangleWithoutIndices.gltf"},
	}
	for _, bb := range benchs {
		b.Run(bb.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := gltf.Open(bb.name)
				if err != nil {
					b.Errorf("Open() error = %v", err)
					return
				}
			}
		})
	}
}

func BenchmarkOpenBinary(b *testing.B) {
	benchs := []struct {
		name string
	}{
		{"testdata/BoxVertexColors/glTF-Binary/BoxVertexColors.glb"},
		{"testdata/OrientationTest/glTF-Binary/OrientationTest.glb"},
	}
	for _, bb := range benchs {
		b.Run(bb.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := gltf.Open(bb.name)
				if err != nil {
					b.Errorf("Open() error = %v", err)
					return
				}
			}
		})
	}
}
