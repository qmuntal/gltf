package binary

import (
	"image/color"
	"reflect"

	"github.com/qmuntal/gltf"
)

// SizeOf returns the size, in bytes, of a component type.
func SizeOf(c gltf.ComponentType) int {
	return map[gltf.ComponentType]int{
		gltf.Byte:          1,
		gltf.UnsignedByte:  1,
		gltf.Short:         2,
		gltf.UnsignedShort: 2,
		gltf.UnsignedInt:   4,
		gltf.Float:         4,
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
func SizeOfElement(c gltf.ComponentType, t gltf.AccessorType) int {
	// special cases
	switch {
	case t == gltf.Mat2 && (c == gltf.Byte || c == gltf.UnsignedByte):
		return 8
	case t == gltf.Mat3 && (c == gltf.Byte || c == gltf.UnsignedByte):
		return 12
	case t == gltf.Mat3 && (c == gltf.Short || c == gltf.UnsignedShort):
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
		c, t = gltf.Byte, gltf.Scalar
	case [][2]int8, [2]int8:
		c, t = gltf.Byte, gltf.Vec2
	case [][3]int8, [3]int8:
		c, t = gltf.Byte, gltf.Vec3
	case [][4]int8, [4]int8:
		c, t = gltf.Byte, gltf.Vec4
	case [][2][2]int8, [2][2]int8:
		c, t = gltf.Byte, gltf.Mat2
	case [][3][3]int8, [3][3]int8:
		c, t = gltf.Byte, gltf.Mat3
	case [][4][4]int8, [4][4]int8:
		c, t = gltf.Byte, gltf.Mat4
	case []uint8, uint8:
		c, t = gltf.UnsignedByte, gltf.Scalar
	case [][2]uint8, [2]uint8:
		c, t = gltf.UnsignedByte, gltf.Vec2
	case [][3]uint8, [3]uint8:
		c, t = gltf.UnsignedByte, gltf.Vec3
	case []color.RGBA, color.RGBA, [][4]uint8, [4]uint8:
		c, t = gltf.UnsignedByte, gltf.Vec4
	case [][2][2]uint8, [2][2]uint8:
		c, t = gltf.UnsignedByte, gltf.Mat2
	case [][3][3]uint8, [3][3]uint8:
		c, t = gltf.UnsignedByte, gltf.Mat3
	case [][4][4]uint8, [4][4]uint8:
		c, t = gltf.UnsignedByte, gltf.Mat4
	case []int16, int16:
		c, t = gltf.Short, gltf.Scalar
	case [][2]int16, [2]int16:
		c, t = gltf.Short, gltf.Vec2
	case [][3]int16, [3]int16:
		c, t = gltf.Short, gltf.Vec3
	case [][4]int16, [4]int16:
		c, t = gltf.Short, gltf.Vec4
	case [][2][2]int16, [2][2]int16:
		c, t = gltf.Short, gltf.Mat2
	case [][3][3]int16, [3][3]int16:
		c, t = gltf.Short, gltf.Mat3
	case [][4][4]int16, [4][4]int16:
		c, t = gltf.Short, gltf.Mat4
	case []uint16, uint16:
		c, t = gltf.UnsignedShort, gltf.Scalar
	case [][2]uint16, [2]uint16:
		c, t = gltf.UnsignedShort, gltf.Vec2
	case [][3]uint16, [3]uint16:
		c, t = gltf.UnsignedShort, gltf.Vec3
	case []color.RGBA64, color.RGBA64, [][4]uint16, [4]uint16:
		c, t = gltf.UnsignedShort, gltf.Vec4
	case [][2][2]uint16, [2][2]uint16:
		c, t = gltf.UnsignedShort, gltf.Mat2
	case [][3][3]uint16, [3][3]uint16:
		c, t = gltf.UnsignedShort, gltf.Mat3
	case [][4][4]uint16, [4][4]uint16:
		c, t = gltf.UnsignedShort, gltf.Mat4
	case []uint32, uint32:
		c, t = gltf.UnsignedInt, gltf.Scalar
	case [][2]uint32, [2]uint32:
		c, t = gltf.UnsignedInt, gltf.Vec2
	case [][3]uint32, [3]uint32:
		c, t = gltf.UnsignedInt, gltf.Vec3
	case [][4]uint32, [4]uint32:
		c, t = gltf.UnsignedInt, gltf.Vec4
	case [][2][2]uint32, [2][2]uint32:
		c, t = gltf.UnsignedInt, gltf.Mat2
	case [][3][3]uint32, [3][3]uint32:
		c, t = gltf.UnsignedInt, gltf.Mat3
	case [][4][4]uint32, [4][4]uint32:
		c, t = gltf.UnsignedInt, gltf.Mat4
	case []float32, float32:
		c, t = gltf.Float, gltf.Scalar
	case [][2]float32, [2]float32:
		c, t = gltf.Float, gltf.Vec2
	case []gltf.RGB, gltf.RGB, [][3]float32, [3]float32:
		c, t = gltf.Float, gltf.Vec3
	case []gltf.RGBA, gltf.RGBA, [][4]float32, [4]float32:
		c, t = gltf.Float, gltf.Vec4
	case [][2][2]float32, [2][2]float32:
		c, t = gltf.Float, gltf.Mat2
	case [][3][3]float32, [3][3]float32:
		c, t = gltf.Float, gltf.Mat3
	case [][4][4]float32, [4][4]float32:
		c, t = gltf.Float, gltf.Mat4
	default:
		length = -1
	}
	return
}
