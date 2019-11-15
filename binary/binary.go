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
package binary

import "encoding/binary"

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

type ubyteComponent struct{}
type ushortComponent struct{}
type uintComponent struct{}

func (ubyteComponent) Scalar(b []byte) uint8 {
	return b[0]
}

func (ubyteComponent) PutScalar(b []byte, v uint8) {
	b[0] = v
}

func (ubyteComponent) Vec2(b []byte) [2]uint8 {
	return [2]uint8{b[0], b[1]}
}

func (ubyteComponent) PutVec2(b []byte, v [2]uint8) {
	b[0] = v[0]
	b[1] = v[1]
}

func (ubyteComponent) Vec3(b []byte) [3]uint8 {
	return [3]uint8{b[0], b[1], b[2]}
}

func (ubyteComponent) PutVec3(b []byte, v [3]uint8) {
	b[0] = v[0]
	b[1] = v[1]
	b[2] = v[2]
}

func (ubyteComponent) Vec4(b []byte) [4]uint8 {
	return [4]uint8{b[0], b[1], b[2], b[3]}
}

func (ubyteComponent) PutVec4(b []byte, v [4]uint8) {
	b[0] = v[0]
	b[1] = v[1]
	b[2] = v[2]
	b[3] = v[3]
}

func (ubyteComponent) Mat2(b []byte) [2][2]uint8 {
	return [2][2]uint8{
		{b[0], b[4]},
		{b[1], b[5]},
	}
}

func (ubyteComponent) PutMat2(b []byte, v [2][2]uint8) {
	b[0] = v[0][0]
	b[1] = v[1][0]
	b[4] = v[0][1]
	b[5] = v[1][1]
}

func (ubyteComponent) Mat3(b []byte) [3][3]uint8 {
	return [3][3]uint8{
		{b[0], b[4], b[8]},
		{b[1], b[5], b[9]},
		{b[2], b[6], b[10]},
	}
}

func (ubyteComponent) PutMat3(b []byte, v [3][3]uint8) {
	b[0] = v[0][0]
	b[1] = v[1][0]
	b[2] = v[2][0]
	b[4] = v[0][1]
	b[5] = v[1][1]
	b[6] = v[2][1]
	b[8] = v[0][2]
	b[9] = v[1][2]
	b[10] = v[2][2]
}
func (ubyteComponent) Mat4(b []byte) [4][4]uint8 {
	return [4][4]uint8{
		{b[0], b[4], b[8], b[12]},
		{b[1], b[5], b[9], b[13]},
		{b[2], b[6], b[10], b[14]},
		{b[3], b[7], b[11], b[15]},
	}
}

func (ubyteComponent) PutMat4(b []byte, v [4][4]uint8) {
	b[0] = v[0][0]
	b[1] = v[1][0]
	b[2] = v[2][0]
	b[3] = v[3][0]
	b[4] = v[0][1]
	b[5] = v[1][1]
	b[6] = v[2][1]
	b[7] = v[3][1]
	b[8] = v[0][2]
	b[9] = v[1][2]
	b[10] = v[2][2]
	b[11] = v[3][2]
	b[12] = v[0][3]
	b[13] = v[1][3]
	b[14] = v[2][3]
	b[15] = v[3][3]
}

func getUint16(b []byte) uint16 {
	return binary.LittleEndian.Uint16(b)
}

func (ushortComponent) Scalar(b []byte) uint16 {
	return getUint16(b)
}

func (ushortComponent) PutScalar(b []byte, v uint16) {
	binary.LittleEndian.PutUint16(b, v)
}

func (ushortComponent) Vec2(b []byte) [2]uint16 {
	return [2]uint16{getUint16(b), getUint16(b[2:])}
}

func (ushortComponent) PutVec2(b []byte, v [2]uint16) {
	binary.LittleEndian.PutUint16(b, v[0])
	binary.LittleEndian.PutUint16(b[2:], v[1])
}

func (ushortComponent) Vec3(b []byte) [3]uint16 {
	return [3]uint16{getUint16(b), getUint16(b[2:]), getUint16(b[4:])}
}

func (ushortComponent) PutVec3(b []byte, v [3]uint16) {
	binary.LittleEndian.PutUint16(b, v[0])
	binary.LittleEndian.PutUint16(b[2:], v[1])
	binary.LittleEndian.PutUint16(b[4:], v[2])
}

func (ushortComponent) Vec4(b []byte) [4]uint16 {
	return [4]uint16{getUint16(b), getUint16(b[2:]), getUint16(b[4:]), getUint16(b[6:])}
}

func (ushortComponent) PutVec4(b []byte, v [4]uint16) {
	binary.LittleEndian.PutUint16(b, v[0])
	binary.LittleEndian.PutUint16(b[2:], v[1])
	binary.LittleEndian.PutUint16(b[4:], v[2])
	binary.LittleEndian.PutUint16(b[6:], v[3])
}

func (ushortComponent) Mat2(b []byte) [2][2]uint16 {
	return [2][2]uint16{
		{getUint16(b), getUint16(b[4:])},
		{getUint16(b[2:]), getUint16(b[6:])},
	}
}

func (ushortComponent) PutMat2(b []byte, v [2][2]uint16) {
	binary.LittleEndian.PutUint16(b, v[0][0])
	binary.LittleEndian.PutUint16(b[2:], v[1][0])
	binary.LittleEndian.PutUint16(b[4:], v[0][1])
	binary.LittleEndian.PutUint16(b[6:], v[1][1])
}

func (ushortComponent) Mat3(b []byte) [3][3]uint16 {
	return [3][3]uint16{
		{getUint16(b), getUint16(b[8:]), getUint16(b[16:])},
		{getUint16(b[2:]), getUint16(b[10:]), getUint16(b[18:])},
		{getUint16(b[4:]), getUint16(b[12:]), getUint16(b[20:])},
	}
}

func (ushortComponent) PutMat3(b []byte, v [3][3]uint16) {
	binary.LittleEndian.PutUint16(b, v[0][0])
	binary.LittleEndian.PutUint16(b[2:], v[1][0])
	binary.LittleEndian.PutUint16(b[4:], v[2][0])
	binary.LittleEndian.PutUint16(b[8:], v[0][1])
	binary.LittleEndian.PutUint16(b[10:], v[1][1])
	binary.LittleEndian.PutUint16(b[12:], v[2][1])
	binary.LittleEndian.PutUint16(b[16:], v[0][2])
	binary.LittleEndian.PutUint16(b[18:], v[1][2])
	binary.LittleEndian.PutUint16(b[20:], v[2][2])
}

func (ushortComponent) Mat4(b []byte) [4][4]uint16 {
	return [4][4]uint16{
		{getUint16(b), getUint16(b[8:]), getUint16(b[16:]), getUint16(b[24:])},
		{getUint16(b[2:]), getUint16(b[10:]), getUint16(b[18:]), getUint16(b[26:])},
		{getUint16(b[4:]), getUint16(b[12:]), getUint16(b[20:]), getUint16(b[28:])},
		{getUint16(b[6:]), getUint16(b[14:]), getUint16(b[22:]), getUint16(b[30:])},
	}
}

func (ushortComponent) PutMat4(b []byte, v [4][4]uint16) {
	binary.LittleEndian.PutUint16(b, v[0][0])
	binary.LittleEndian.PutUint16(b[2:], v[1][0])
	binary.LittleEndian.PutUint16(b[4:], v[2][0])
	binary.LittleEndian.PutUint16(b[6:], v[3][0])
	binary.LittleEndian.PutUint16(b[8:], v[0][1])
	binary.LittleEndian.PutUint16(b[10:], v[1][1])
	binary.LittleEndian.PutUint16(b[12:], v[2][1])
	binary.LittleEndian.PutUint16(b[14:], v[3][1])
	binary.LittleEndian.PutUint16(b[16:], v[0][2])
	binary.LittleEndian.PutUint16(b[18:], v[1][2])
	binary.LittleEndian.PutUint16(b[20:], v[2][2])
	binary.LittleEndian.PutUint16(b[22:], v[3][2])
	binary.LittleEndian.PutUint16(b[24:], v[0][3])
	binary.LittleEndian.PutUint16(b[26:], v[1][3])
	binary.LittleEndian.PutUint16(b[28:], v[2][3])
	binary.LittleEndian.PutUint16(b[30:], v[3][3])
}

func getUint32(b []byte) uint32 {
	return binary.LittleEndian.Uint32(b)
}

func (uintComponent) Scalar(b []byte) uint32 {
	return getUint32(b)
}

func (uintComponent) PutScalar(b []byte, v uint32) {
	binary.LittleEndian.PutUint32(b, v)
}

func (uintComponent) Vec2(b []byte) [2]uint32 {
	return [2]uint32{getUint32(b), getUint32(b[4:])}
}

func (uintComponent) PutVec2(b []byte, v [2]uint32) {
	binary.LittleEndian.PutUint32(b, v[0])
	binary.LittleEndian.PutUint32(b[4:], v[1])
}

func (uintComponent) Vec3(b []byte) [3]uint32 {
	return [3]uint32{getUint32(b), getUint32(b[4:]), getUint32(b[8:])}
}

func (uintComponent) PutVec3(b []byte, v [3]uint32) {
	binary.LittleEndian.PutUint32(b, v[0])
	binary.LittleEndian.PutUint32(b[4:], v[1])
	binary.LittleEndian.PutUint32(b[8:], v[2])
}

func (uintComponent) Vec4(b []byte) [4]uint32 {
	return [4]uint32{getUint32(b), getUint32(b[4:]), getUint32(b[8:]), getUint32(b[12:])}
}

func (uintComponent) PutVec4(b []byte, v [4]uint32) {
	binary.LittleEndian.PutUint32(b, v[0])
	binary.LittleEndian.PutUint32(b[4:], v[1])
	binary.LittleEndian.PutUint32(b[8:], v[2])
	binary.LittleEndian.PutUint32(b[12:], v[3])
}

func (uintComponent) Mat2(b []byte) [2][2]uint32 {
	return [2][2]uint32{
		{getUint32(b), getUint32(b[8:])},
		{getUint32(b[4:]), getUint32(b[12:])},
	}
}

func (uintComponent) PutMat2(b []byte, v [2][2]uint32) {
	binary.LittleEndian.PutUint32(b, v[0][0])
	binary.LittleEndian.PutUint32(b[4:], v[1][0])
	binary.LittleEndian.PutUint32(b[8:], v[0][1])
	binary.LittleEndian.PutUint32(b[12:], v[1][1])
}

func (uintComponent) Mat3(b []byte) [3][3]uint32 {
	return [3][3]uint32{
		{getUint32(b), getUint32(b[12:]), getUint32(b[24:])},
		{getUint32(b[4:]), getUint32(b[16:]), getUint32(b[28:])},
		{getUint32(b[8:]), getUint32(b[20:]), getUint32(b[32:])},
	}
}

func (uintComponent) PutMat3(b []byte, v [3][3]uint32) {
	binary.LittleEndian.PutUint32(b, v[0][0])
	binary.LittleEndian.PutUint32(b[4:], v[1][0])
	binary.LittleEndian.PutUint32(b[8:], v[2][0])
	binary.LittleEndian.PutUint32(b[12:], v[0][1])
	binary.LittleEndian.PutUint32(b[16:], v[1][1])
	binary.LittleEndian.PutUint32(b[20:], v[2][1])
	binary.LittleEndian.PutUint32(b[24:], v[0][2])
	binary.LittleEndian.PutUint32(b[28:], v[1][2])
	binary.LittleEndian.PutUint32(b[32:], v[2][2])
}

func (uintComponent) Mat4(b []byte) [4][4]uint32 {
	return [4][4]uint32{
		{getUint32(b), getUint32(b[16:]), getUint32(b[32:]), getUint32(b[48:])},
		{getUint32(b[4:]), getUint32(b[20:]), getUint32(b[36:]), getUint32(b[52:])},
		{getUint32(b[8:]), getUint32(b[24:]), getUint32(b[40:]), getUint32(b[56:])},
		{getUint32(b[12:]), getUint32(b[28:]), getUint32(b[44:]), getUint32(b[60:])},
	}
}

func (uintComponent) PutMat4(b []byte, v [4][4]uint32) {
	binary.LittleEndian.PutUint32(b, v[0][0])
	binary.LittleEndian.PutUint32(b[4:], v[1][0])
	binary.LittleEndian.PutUint32(b[8:], v[2][0])
	binary.LittleEndian.PutUint32(b[12:], v[3][0])
	binary.LittleEndian.PutUint32(b[16:], v[0][1])
	binary.LittleEndian.PutUint32(b[20:], v[1][1])
	binary.LittleEndian.PutUint32(b[24:], v[2][1])
	binary.LittleEndian.PutUint32(b[28:], v[3][1])
	binary.LittleEndian.PutUint32(b[32:], v[0][2])
	binary.LittleEndian.PutUint32(b[36:], v[1][2])
	binary.LittleEndian.PutUint32(b[40:], v[2][2])
	binary.LittleEndian.PutUint32(b[44:], v[3][2])
	binary.LittleEndian.PutUint32(b[48:], v[0][3])
	binary.LittleEndian.PutUint32(b[52:], v[1][3])
	binary.LittleEndian.PutUint32(b[56:], v[2][3])
	binary.LittleEndian.PutUint32(b[60:], v[3][3])
}
