// Package binary implements simple translation between numbers and byte
// sequences and encoding and decoding of varints.
//
// Numbers are translated by reading and writing fixed-size values.
// A fixed-size value is either a fixed-size arithmetic
// type (bool, int8, uint8, int16, float32, complex64, ...)
// or an array or struct containing only fixed-size values.
//
// Even though the defined structs don't implement a common interface,
// they all follow the same `generic` interface, and the only difference is the
// type of the value:
//
// type Component interface {
// 	 PutScalar([]byte, interface{})
// 	 PutVec2([]byte, [2]interface{})
// 	 PutVec3([]byte, [3]interface{})
// 	 PutVec4([]byte, [4]interface{})
// 	 PutMat2([]byte, [2][2]interface{})
// 	 PutMat3([]byte, [3][3]interface{})
// 	 PutMat4([]byte, [4][4]interface{})
// }
package binary

// Byte is the byte implementation of Component.
type Byte byteComponent

// UnsignedByte is the usnigned byte implementation of Component.
type UnsignedByte ubyteComponent

// Short is the byte short of Component.
type Short shortComponent

// UnsignedShort is the usnigned short implementation of Component.
type UnsignedShort ushortComponent

// UnsignedInt is the unsigned int implementation of Component.
type UnsignedInt uintComponent

// Float is the float implementation of Component.
type Float floatComponent

type byteComponent struct{}
type ubyteComponent struct{}
type shortComponent struct{}
type ushortComponent struct{}
type uintComponent struct{}
type floatComponent struct{}
