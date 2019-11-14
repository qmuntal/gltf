// Package binary implements simple translation between numbers and byte sequences.
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
// 	 Scalar([]byte) interface{}
//   Vec2([]byte) [2]interface{}
//   Vec3([]byte) [2]interface{}
//   Vec4([]byte) [2]interface{}
//   Mat2([]byte) [2][2]interface{}
//   Mat3([]byte) [3][3]interface{}
//   Mat4([]byte) [4][4]interface{}
// 	 PutScalar([]byte, interface{})
// 	 PutVec2([]byte, [2]interface{})
// 	 PutVec3([]byte, [3]interface{})
// 	 PutVec4([]byte, [4]interface{})
// 	 PutMat2([]byte, [2][2]interface{})
// 	 PutMat3([]byte, [3][3]interface{})
// 	 PutMat4([]byte, [4][4]interface{})
// }
//
// This package favors simplicity over efficiency. Clients that require
// high-performance serialization, especially for large data structures,
// should look at more advanced solutions such as custom encoders.
package binary

// Byte is the byte implementation of Component.
var Byte byteComponent

// UnsignedByte is the usnigned byte implementation of Component.
var UnsignedByte ubyteComponent

// Short is the byte short of Component.
var Short shortComponent

// UnsignedShort is the usnigned short implementation of Component.
var UnsignedShort ushortComponent

// UnsignedInt is the unsigned int implementation of Component.
var UnsignedInt uintComponent

// Float is the float implementation of Component.
var Float floatComponent

type byteComponent struct{}
type ubyteComponent struct{}
type shortComponent struct{}
type ushortComponent struct{}
type uintComponent struct{}
type floatComponent struct{}
