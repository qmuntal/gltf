package binary

import (
	"image/color"
	"reflect"

	"github.com/qmuntal/gltf"
)

// SizeOfElement returns the size, in bytes, of an element.
// The element size may not be (component size) * (number of components),
// as some of the elements are tightly packed in order to ensure
// that they are aligned to 4-byte boundaries.
func SizeOfElement(c gltf.ComponentType, t gltf.AccessorType) uint32 {
	// special cases
	switch {
	case (t == gltf.AccessorVec3 || t == gltf.AccessorVec2) && (c == gltf.ComponentByte || c == gltf.ComponentUbyte):
		return 4
	case t == gltf.AccessorVec3 && (c == gltf.ComponentShort || c == gltf.ComponentUshort):
		return 8
	case t == gltf.AccessorMat2 && (c == gltf.ComponentByte || c == gltf.ComponentUbyte):
		return 8
	case t == gltf.AccessorMat3 && (c == gltf.ComponentByte || c == gltf.ComponentUbyte):
		return 12
	case t == gltf.AccessorMat3 && (c == gltf.ComponentShort || c == gltf.ComponentUshort):
		return 24
	}
	return c.ByteSize() * t.Components()
}

// Type returns the associated glTF type data.
// If data is an slice, it also returns the length of the slice.
// If data does not have an associated glTF type length will be 0.
func Type(data interface{}) (c gltf.ComponentType, t gltf.AccessorType, length uint32) {
	v := reflect.ValueOf(data)
	switch v.Kind() {
	case reflect.Slice:
		length = uint32(v.Len())
	}
	switch data.(type) {
	case []int8, int8:
		c, t = gltf.ComponentByte, gltf.AccessorScalar
	case [][2]int8, [2]int8:
		c, t = gltf.ComponentByte, gltf.AccessorVec2
	case [][3]int8, [3]int8:
		c, t = gltf.ComponentByte, gltf.AccessorVec3
	case [][4]int8, [4]int8:
		c, t = gltf.ComponentByte, gltf.AccessorVec4
	case [][2][2]int8, [2][2]int8:
		c, t = gltf.ComponentByte, gltf.AccessorMat2
	case [][3][3]int8, [3][3]int8:
		c, t = gltf.ComponentByte, gltf.AccessorMat3
	case [][4][4]int8, [4][4]int8:
		c, t = gltf.ComponentByte, gltf.AccessorMat4
	case []uint8, uint8:
		c, t = gltf.ComponentUbyte, gltf.AccessorScalar
	case [][2]uint8, [2]uint8:
		c, t = gltf.ComponentUbyte, gltf.AccessorVec2
	case [][3]uint8, [3]uint8:
		c, t = gltf.ComponentUbyte, gltf.AccessorVec3
	case []color.RGBA, color.RGBA, [][4]uint8, [4]uint8:
		c, t = gltf.ComponentUbyte, gltf.AccessorVec4
	case [][2][2]uint8, [2][2]uint8:
		c, t = gltf.ComponentUbyte, gltf.AccessorMat2
	case [][3][3]uint8, [3][3]uint8:
		c, t = gltf.ComponentUbyte, gltf.AccessorMat3
	case [][4][4]uint8, [4][4]uint8:
		c, t = gltf.ComponentUbyte, gltf.AccessorMat4
	case []int16, int16:
		c, t = gltf.ComponentShort, gltf.AccessorScalar
	case [][2]int16, [2]int16:
		c, t = gltf.ComponentShort, gltf.AccessorVec2
	case [][3]int16, [3]int16:
		c, t = gltf.ComponentShort, gltf.AccessorVec3
	case [][4]int16, [4]int16:
		c, t = gltf.ComponentShort, gltf.AccessorVec4
	case [][2][2]int16, [2][2]int16:
		c, t = gltf.ComponentShort, gltf.AccessorMat2
	case [][3][3]int16, [3][3]int16:
		c, t = gltf.ComponentShort, gltf.AccessorMat3
	case [][4][4]int16, [4][4]int16:
		c, t = gltf.ComponentShort, gltf.AccessorMat4
	case []uint16, uint16:
		c, t = gltf.ComponentUshort, gltf.AccessorScalar
	case [][2]uint16, [2]uint16:
		c, t = gltf.ComponentUshort, gltf.AccessorVec2
	case [][3]uint16, [3]uint16:
		c, t = gltf.ComponentUshort, gltf.AccessorVec3
	case []color.RGBA64, color.RGBA64, [][4]uint16, [4]uint16:
		c, t = gltf.ComponentUshort, gltf.AccessorVec4
	case [][2][2]uint16, [2][2]uint16:
		c, t = gltf.ComponentUshort, gltf.AccessorMat2
	case [][3][3]uint16, [3][3]uint16:
		c, t = gltf.ComponentUshort, gltf.AccessorMat3
	case [][4][4]uint16, [4][4]uint16:
		c, t = gltf.ComponentUshort, gltf.AccessorMat4
	case []uint32, uint32:
		c, t = gltf.ComponentUint, gltf.AccessorScalar
	case [][2]uint32, [2]uint32:
		c, t = gltf.ComponentUint, gltf.AccessorVec2
	case [][3]uint32, [3]uint32:
		c, t = gltf.ComponentUint, gltf.AccessorVec3
	case [][4]uint32, [4]uint32:
		c, t = gltf.ComponentUint, gltf.AccessorVec4
	case [][2][2]uint32, [2][2]uint32:
		c, t = gltf.ComponentUint, gltf.AccessorMat2
	case [][3][3]uint32, [3][3]uint32:
		c, t = gltf.ComponentUint, gltf.AccessorMat3
	case [][4][4]uint32, [4][4]uint32:
		c, t = gltf.ComponentUint, gltf.AccessorMat4
	case []float32, float32:
		c, t = gltf.ComponentFloat, gltf.AccessorScalar
	case [][2]float32, [2]float32:
		c, t = gltf.ComponentFloat, gltf.AccessorVec2
	case []gltf.RGB, gltf.RGB, [][3]float32, [3]float32:
		c, t = gltf.ComponentFloat, gltf.AccessorVec3
	case []gltf.RGBA, gltf.RGBA, [][4]float32, [4]float32:
		c, t = gltf.ComponentFloat, gltf.AccessorVec4
	case [][2][2]float32, [2][2]float32:
		c, t = gltf.ComponentFloat, gltf.AccessorMat2
	case [][3][3]float32, [3][3]float32:
		c, t = gltf.ComponentFloat, gltf.AccessorMat3
	case [][4][4]float32, [4][4]float32:
		c, t = gltf.ComponentFloat, gltf.AccessorMat4
	default:
		length = 0
	}
	return
}
