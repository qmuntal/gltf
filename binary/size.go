package binary

import "github.com/qmuntal/gltf"

const (
	sizeByte   = 1
	sizeUbyte  = 1
	sizeShort  = 2
	sizeUshort = 2
	sizeUint   = 4
	sizeFloat  = 4
)

const (
	componentsScalar = 1
	componentsVec2   = 2
	componentsVec3   = 3
	componentsVec4   = 4
	componentsMat2   = 4
	componentsMat3   = 9
	componentsMat4   = 16
)

// SizeOf returns the size, in bytes, of a component type.
func SizeOf(c gltf.ComponentType) int {
	return map[gltf.ComponentType]int{
		gltf.Byte:          sizeByte,
		gltf.UnsignedByte:  sizeUbyte,
		gltf.Short:         sizeShort,
		gltf.UnsignedShort: sizeUshort,
		gltf.UnsignedInt:   sizeUint,
		gltf.Float:         sizeFloat,
	}[c]
}

// ComponentsOf returns the number of components of an accessor type.
func ComponentsOf(t gltf.AccessorType) int {
	return map[gltf.AccessorType]int{
		gltf.Scalar: componentsScalar,
		gltf.Vec2:   componentsVec2,
		gltf.Vec3:   componentsVec3,
		gltf.Vec4:   componentsVec4,
		gltf.Mat2:   componentsMat2,
		gltf.Mat3:   componentsMat3,
		gltf.Mat4:   componentsMat4,
	}[t]
}

// ElementSize returns the size, in bytes, of an element.
func ElementSize(c gltf.ComponentType, t gltf.AccessorType) int {
	return SizeOf(c) * ComponentsOf(t)
}
