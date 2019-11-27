package binary

import (
	"image/color"
	"reflect"

	"github.com/qmuntal/gltf"
)

// SizeOf returns the size, in bytes, of a component type.
func SizeOf(c gltf.ComponentType) int {
	return map[gltf.ComponentType]int{
		gltf.ComponentByte:   1,
		gltf.ComponentUbyte:  1,
		gltf.ComponentShort:  2,
		gltf.ComponentUshort: 2,
		gltf.ComponentUint:   4,
		gltf.ComponentFloat:  4,
	}[c]
}

// ComponentsOf returns the number of components of an accessor type.
func ComponentsOf(t gltf.AccessorType) int {
	return map[gltf.AccessorType]int{
		gltf.Scalar: 1,
		gltf.Vec2:   2,
		gltf.Vec3:   3,
		gltf.Vec4:   4,
		gltf.Mat2:   4,
		gltf.Mat3:   9,
		gltf.Mat4:   16,
	}[t]
}

// SizeOfElement returns the size, in bytes, of an element.
// The element size may not be (component size) * (number of components),
// as some of the elements are tightly packed in order to ensure
// that they are aligned to 4-byte boundaries.
func SizeOfElement(c gltf.ComponentType, t gltf.AccessorType) int {
	// special cases
	switch {
	case (t == gltf.Vec3 || t == gltf.Vec2) && (c == gltf.ComponentByte || c == gltf.ComponentUbyte):
		return 4
	case t == gltf.Vec3 && (c == gltf.ComponentShort || c == gltf.ComponentUshort):
		return 8
	case t == gltf.Mat2 && (c == gltf.ComponentByte || c == gltf.ComponentUbyte):
		return 8
	case t == gltf.Mat3 && (c == gltf.ComponentByte || c == gltf.ComponentUbyte):
		return 12
	case t == gltf.Mat3 && (c == gltf.ComponentShort || c == gltf.ComponentUshort):
		return 24
	}
	return SizeOf(c) * ComponentsOf(t)
}

// Type returns the associated glTF type data.
// If data is an slice, it also returns the length of the slice.
// If data does not have an associated glTF type, length will be -1.
func Type(data interface{}) (c gltf.ComponentType, t gltf.AccessorType, length int) {
	v := reflect.ValueOf(data)
	switch v.Kind() {
	case reflect.Slice:
		length = v.Len()
	}
	switch data.(type) {
	case []int8, int8:
		c, t = gltf.ComponentByte, gltf.Scalar
	case [][2]int8, [2]int8:
		c, t = gltf.ComponentByte, gltf.Vec2
	case [][3]int8, [3]int8:
		c, t = gltf.ComponentByte, gltf.Vec3
	case [][4]int8, [4]int8:
		c, t = gltf.ComponentByte, gltf.Vec4
	case [][2][2]int8, [2][2]int8:
		c, t = gltf.ComponentByte, gltf.Mat2
	case [][3][3]int8, [3][3]int8:
		c, t = gltf.ComponentByte, gltf.Mat3
	case [][4][4]int8, [4][4]int8:
		c, t = gltf.ComponentByte, gltf.Mat4
	case []uint8, uint8:
		c, t = gltf.ComponentUbyte, gltf.Scalar
	case [][2]uint8, [2]uint8:
		c, t = gltf.ComponentUbyte, gltf.Vec2
	case [][3]uint8, [3]uint8:
		c, t = gltf.ComponentUbyte, gltf.Vec3
	case []color.RGBA, color.RGBA, [][4]uint8, [4]uint8:
		c, t = gltf.ComponentUbyte, gltf.Vec4
	case [][2][2]uint8, [2][2]uint8:
		c, t = gltf.ComponentUbyte, gltf.Mat2
	case [][3][3]uint8, [3][3]uint8:
		c, t = gltf.ComponentUbyte, gltf.Mat3
	case [][4][4]uint8, [4][4]uint8:
		c, t = gltf.ComponentUbyte, gltf.Mat4
	case []int16, int16:
		c, t = gltf.ComponentShort, gltf.Scalar
	case [][2]int16, [2]int16:
		c, t = gltf.ComponentShort, gltf.Vec2
	case [][3]int16, [3]int16:
		c, t = gltf.ComponentShort, gltf.Vec3
	case [][4]int16, [4]int16:
		c, t = gltf.ComponentShort, gltf.Vec4
	case [][2][2]int16, [2][2]int16:
		c, t = gltf.ComponentShort, gltf.Mat2
	case [][3][3]int16, [3][3]int16:
		c, t = gltf.ComponentShort, gltf.Mat3
	case [][4][4]int16, [4][4]int16:
		c, t = gltf.ComponentShort, gltf.Mat4
	case []uint16, uint16:
		c, t = gltf.ComponentUshort, gltf.Scalar
	case [][2]uint16, [2]uint16:
		c, t = gltf.ComponentUshort, gltf.Vec2
	case [][3]uint16, [3]uint16:
		c, t = gltf.ComponentUshort, gltf.Vec3
	case []color.RGBA64, color.RGBA64, [][4]uint16, [4]uint16:
		c, t = gltf.ComponentUshort, gltf.Vec4
	case [][2][2]uint16, [2][2]uint16:
		c, t = gltf.ComponentUshort, gltf.Mat2
	case [][3][3]uint16, [3][3]uint16:
		c, t = gltf.ComponentUshort, gltf.Mat3
	case [][4][4]uint16, [4][4]uint16:
		c, t = gltf.ComponentUshort, gltf.Mat4
	case []uint32, uint32:
		c, t = gltf.ComponentUint, gltf.Scalar
	case [][2]uint32, [2]uint32:
		c, t = gltf.ComponentUint, gltf.Vec2
	case [][3]uint32, [3]uint32:
		c, t = gltf.ComponentUint, gltf.Vec3
	case [][4]uint32, [4]uint32:
		c, t = gltf.ComponentUint, gltf.Vec4
	case [][2][2]uint32, [2][2]uint32:
		c, t = gltf.ComponentUint, gltf.Mat2
	case [][3][3]uint32, [3][3]uint32:
		c, t = gltf.ComponentUint, gltf.Mat3
	case [][4][4]uint32, [4][4]uint32:
		c, t = gltf.ComponentUint, gltf.Mat4
	case []float32, float32:
		c, t = gltf.ComponentFloat, gltf.Scalar
	case [][2]float32, [2]float32:
		c, t = gltf.ComponentFloat, gltf.Vec2
	case []gltf.RGB, gltf.RGB, [][3]float32, [3]float32:
		c, t = gltf.ComponentFloat, gltf.Vec3
	case []gltf.RGBA, gltf.RGBA, [][4]float32, [4]float32:
		c, t = gltf.ComponentFloat, gltf.Vec4
	case [][2][2]float32, [2][2]float32:
		c, t = gltf.ComponentFloat, gltf.Mat2
	case [][3][3]float32, [3][3]float32:
		c, t = gltf.ComponentFloat, gltf.Mat3
	case [][4][4]float32, [4][4]float32:
		c, t = gltf.ComponentFloat, gltf.Mat4
	default:
		length = -1
	}
	return
}
