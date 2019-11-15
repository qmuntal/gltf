package binary

import "github.com/qmuntal/gltf"

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

// ElementSize returns the size, in bytes, of an element.
func ElementSize(c gltf.ComponentType, t gltf.AccessorType) int {
	return SizeOf(c) * ComponentsOf(t)
}

// intDataSize returns the size of the data required to represent the data when encoded.
// It returns zero if the type cannot be implemented by the fast path in Read or Write.
func intDataSize(data interface{}) (element int, total int) {
	var length int
	switch data := data.(type) {
	case []int8:
		element = ElementSize(gltf.Byte, gltf.Scalar)
		length = len(data)
	case [][2]int8:
		element = ElementSize(gltf.Byte, gltf.Vec2)
		length = len(data)
	case [][3]int8:
		element = ElementSize(gltf.Byte, gltf.Vec3)
		length = len(data)
	case [][4]int8:
		element = ElementSize(gltf.Byte, gltf.Vec4)
		length = len(data)
	case [][2][2]int8:
		element = ElementSize(gltf.Byte, gltf.Mat2)
		length = len(data)
	case [][3][3]int8:
		element = ElementSize(gltf.Byte, gltf.Mat3)
		length = len(data)
	case [][4][4]int8:
		element = ElementSize(gltf.Byte, gltf.Mat4)
		length = len(data)
	case []uint8:
		element = ElementSize(gltf.UnsignedByte, gltf.Scalar)
		length = len(data)
	case [][2]uint8:
		element = ElementSize(gltf.UnsignedByte, gltf.Vec2)
		length = len(data)
	case [][3]uint8:
		element = ElementSize(gltf.UnsignedByte, gltf.Vec3)
		length = len(data)
	case [][4]uint8:
		element = ElementSize(gltf.UnsignedByte, gltf.Vec4)
		length = len(data)
	case [][2][2]uint8:
		element = ElementSize(gltf.UnsignedByte, gltf.Mat2)
		length = len(data)
	case [][3][3]uint8:
		element = ElementSize(gltf.UnsignedByte, gltf.Mat3)
		length = len(data)
	case [][4][4]uint8:
		element = ElementSize(gltf.UnsignedByte, gltf.Mat4)
		length = len(data)
	case []int16:
		element = ElementSize(gltf.Short, gltf.Scalar)
		length = len(data)
	case [][2]int16:
		element = ElementSize(gltf.Short, gltf.Vec2)
		length = len(data)
	case [][3]int16:
		element = ElementSize(gltf.Short, gltf.Vec3)
		length = len(data)
	case [][4]int16:
		element = ElementSize(gltf.Short, gltf.Vec4)
		length = len(data)
	case [][2][2]int16:
		element = ElementSize(gltf.Short, gltf.Mat2)
		length = len(data)
	case [][3][3]int16:
		element = ElementSize(gltf.Short, gltf.Mat3)
		length = len(data)
	case [][4][4]int16:
		element = ElementSize(gltf.Short, gltf.Mat4)
		length = len(data)
	case []uint16:
		element = ElementSize(gltf.UnsignedShort, gltf.Scalar)
		length = len(data)
	case [][2]uint16:
		element = ElementSize(gltf.UnsignedShort, gltf.Vec2)
		length = len(data)
	case [][3]uint16:
		element = ElementSize(gltf.UnsignedShort, gltf.Vec3)
		length = len(data)
	case [][4]uint16:
		element = ElementSize(gltf.UnsignedShort, gltf.Vec4)
		length = len(data)
	case [][2][2]uint16:
		element = ElementSize(gltf.UnsignedShort, gltf.Mat2)
		length = len(data)
	case [][3][3]uint16:
		element = ElementSize(gltf.UnsignedShort, gltf.Mat3)
		length = len(data)
	case [][4][4]uint16:
		element = ElementSize(gltf.UnsignedShort, gltf.Mat4)
		length = len(data)
	case []uint32:
		element = ElementSize(gltf.UnsignedInt, gltf.Scalar)
		length = len(data)
	case [][2]uint32:
		element = ElementSize(gltf.UnsignedInt, gltf.Vec2)
		length = len(data)
	case [][3]uint32:
		element = ElementSize(gltf.UnsignedInt, gltf.Vec3)
		length = len(data)
	case [][4]uint32:
		element = ElementSize(gltf.UnsignedInt, gltf.Vec4)
		length = len(data)
	case [][2][2]uint32:
		element = ElementSize(gltf.UnsignedInt, gltf.Mat2)
		length = len(data)
	case [][3][3]uint32:
		element = ElementSize(gltf.UnsignedInt, gltf.Mat3)
		length = len(data)
	case [][4][4]uint32:
		element = ElementSize(gltf.UnsignedInt, gltf.Mat4)
		length = len(data)
	case []float32:
		element = ElementSize(gltf.Float, gltf.Scalar)
		length = len(data)
	case [][2]float32:
		element = ElementSize(gltf.Float, gltf.Vec2)
		length = len(data)
	case [][3]float32:
		element = ElementSize(gltf.Float, gltf.Vec3)
		length = len(data)
	case [][4]float32:
		element = ElementSize(gltf.Float, gltf.Vec4)
		length = len(data)
	case [][2][2]float32:
		element = ElementSize(gltf.Float, gltf.Mat2)
		length = len(data)
	case [][3][3]float32:
		element = ElementSize(gltf.Float, gltf.Mat3)
		length = len(data)
	case [][4][4]float32:
		element = ElementSize(gltf.Float, gltf.Mat4)
		length = len(data)
	}
	return element, length
}
