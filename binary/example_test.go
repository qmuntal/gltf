package binary_test

import (
	"fmt"
	"github.com/qmuntal/gltf"
	"github.com/qmuntal/gltf/binary"
)

func ExampleWrite() {
	// Define data
	indices := []uint8{0, 1, 2, 3, 1, 0, 0, 2, 3, 1, 4, 2, 4, 3, 2, 4, 1, 3}
	vertices := [][3]float32{{43, 43, 0}, {83, 43, 0}, {63, 63, 40}, {43, 83, 0}, {83, 83, 0}}

	// Allocate buffer
	sizeIndices := len(indices) * binary.SizeOfElement(gltf.UnsignedByte, gltf.Scalar)
	sizeVertices := len(vertices) * binary.SizeOfElement(gltf.Float, gltf.Vec3)
	b := make([]byte, sizeIndices+sizeVertices)

	// Write
	binary.Write(b, indices)
	binary.Write(b[sizeIndices:], vertices)

	fmt.Print(b)
	// Output:
	// [0 1 2 3 1 0 0 2 3 1 4 2 4 3 2 4 1 3 0 0 44 66 0 0 44 66 0 0 0 0 0 0 166 66 0 0 44 66 0 0 0 0 0 0 124 66 0 0 124 66 0 0 32 66 0 0 44 66 0 0 166 66 0 0 0 0 0 0 166 66 0 0 166 66 0 0 0 0]
}

func ExampleRead() {
	// Allocate data
	indices := make([]uint8, 18)
	vertices := make([][3]float32, 5)

	// Define buffer
	b := []byte{0, 1, 2, 3, 1, 0, 0, 2, 3, 1, 4, 2, 4, 3, 2, 4, 1, 3, 0, 0, 44, 66, 0, 0, 44, 66, 0, 0, 0, 0, 0, 0, 166, 66, 0, 0, 44, 66, 0, 0, 0, 0, 0, 0, 124, 66, 0, 0, 124, 66, 0, 0, 32, 66, 0, 0, 44, 66, 0, 0, 166, 66, 0, 0, 0, 0, 0, 0, 166, 66, 0, 0, 166, 66, 0, 0, 0, 0}
	sizeIndices := len(indices) * binary.SizeOfElement(gltf.UnsignedByte, gltf.Scalar)

	// Write
	binary.Read(b, indices)
	binary.Read(b[sizeIndices:], vertices)

	fmt.Println(indices)
	fmt.Println(vertices)
	// Output:
	// [0 1 2 3 1 0 0 2 3 1 4 2 4 3 2 4 1 3]
	// [[43 43 0] [83 43 0] [63 63 40] [43 83 0] [83 83 0]]
}
