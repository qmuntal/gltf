package binary

import (
	"unsafe"
)

type byteComponent struct{}
type shortComponent struct{}
type floatComponent struct{}

func (byteComponent) Scalar(b []byte) int8 {
	v := Ubyte.Scalar(b)
	return *(*int8)(unsafe.Pointer(&v))
}

func (byteComponent) PutScalar(b []byte, v int8) {
	Ubyte.PutScalar(b, *(*uint8)(unsafe.Pointer(&v)))
}

func (byteComponent) Vec2(b []byte) [2]int8 {
	v := Ubyte.Vec2(b)
	return *(*[2]int8)(unsafe.Pointer(&v))
}

func (byteComponent) PutVec2(b []byte, v [2]int8) {
	Ubyte.PutVec2(b, *(*[2]uint8)(unsafe.Pointer(&v)))
}

func (byteComponent) Vec3(b []byte) [3]int8 {
	v := Ubyte.Vec3(b)
	return *(*[3]int8)(unsafe.Pointer(&v))
}

func (byteComponent) PutVec3(b []byte, v [3]int8) {
	Ubyte.PutVec3(b, *(*[3]uint8)(unsafe.Pointer(&v)))
}

func (byteComponent) Vec4(b []byte) [4]int8 {
	v := Ubyte.Vec4(b)
	return *(*[4]int8)(unsafe.Pointer(&v))
}

func (byteComponent) PutVec4(b []byte, v [4]int8) {
	Ubyte.PutVec4(b, *(*[4]uint8)(unsafe.Pointer(&v)))
}

func (byteComponent) Mat2(b []byte) [2][2]int8 {
	v := Ubyte.Mat2(b)
	return *(*[2][2]int8)(unsafe.Pointer(&v))
}

func (byteComponent) PutMat2(b []byte, v [2][2]int8) {
	Ubyte.PutMat2(b, *(*[2][2]uint8)(unsafe.Pointer(&v)))
}

func (byteComponent) Mat3(b []byte) [3][3]int8 {
	v := Ubyte.Mat3(b)
	return *(*[3][3]int8)(unsafe.Pointer(&v))
}

func (byteComponent) PutMat3(b []byte, v [3][3]int8) {
	Ubyte.PutMat3(b, *(*[3][3]uint8)(unsafe.Pointer(&v)))
}

func (byteComponent) Mat4(b []byte) [4][4]int8 {
	v := Ubyte.Mat4(b)
	return *(*[4][4]int8)(unsafe.Pointer(&v))
}

func (byteComponent) PutMat4(b []byte, v [4][4]int8) {
	Ubyte.PutMat4(b, *(*[4][4]uint8)(unsafe.Pointer(&v)))
}

func (shortComponent) Scalar(b []byte) int16 {
	v := Ushort.Scalar(b)
	return *(*int16)(unsafe.Pointer(&v))
}

func (shortComponent) PutScalar(b []byte, v int16) {
	Ushort.PutScalar(b, *(*uint16)(unsafe.Pointer(&v)))
}

func (shortComponent) Vec2(b []byte) [2]int16 {
	v := Ushort.Vec2(b)
	return *(*[2]int16)(unsafe.Pointer(&v))
}

func (shortComponent) PutVec2(b []byte, v [2]int16) {
	Ushort.PutVec2(b, *(*[2]uint16)(unsafe.Pointer(&v)))
}

func (shortComponent) Vec3(b []byte) [3]int16 {
	v := Ushort.Vec3(b)
	return *(*[3]int16)(unsafe.Pointer(&v))
}

func (shortComponent) PutVec3(b []byte, v [3]int16) {
	Ushort.PutVec3(b, *(*[3]uint16)(unsafe.Pointer(&v)))
}

func (shortComponent) Vec4(b []byte) [4]int16 {
	v := Ushort.Vec4(b)
	return *(*[4]int16)(unsafe.Pointer(&v))
}

func (shortComponent) PutVec4(b []byte, v [4]int16) {
	Ushort.PutVec4(b, *(*[4]uint16)(unsafe.Pointer(&v)))
}

func (shortComponent) Mat2(b []byte) [2][2]int16 {
	v := Ushort.Mat2(b)
	return *(*[2][2]int16)(unsafe.Pointer(&v))
}

func (shortComponent) PutMat2(b []byte, v [2][2]int16) {
	Ushort.PutMat2(b, *(*[2][2]uint16)(unsafe.Pointer(&v)))
}

func (shortComponent) Mat3(b []byte) [3][3]int16 {
	v := Ushort.Mat3(b)
	return *(*[3][3]int16)(unsafe.Pointer(&v))
}

func (shortComponent) PutMat3(b []byte, v [3][3]int16) {
	Ushort.PutMat3(b, *(*[3][3]uint16)(unsafe.Pointer(&v)))
}

func (shortComponent) Mat4(b []byte) [4][4]int16 {
	v := Ushort.Mat4(b)
	return *(*[4][4]int16)(unsafe.Pointer(&v))
}

func (shortComponent) PutMat4(b []byte, v [4][4]int16) {
	Ushort.PutMat4(b, *(*[4][4]uint16)(unsafe.Pointer(&v)))
}

func (floatComponent) Scalar(b []byte) float32 {
	v := UInt.Scalar(b)
	return *(*float32)(unsafe.Pointer(&v))
}

func (floatComponent) PutScalar(b []byte, v float32) {
	UInt.PutScalar(b, *(*uint32)(unsafe.Pointer(&v)))
}

func (floatComponent) Vec2(b []byte) [2]float32 {
	v := UInt.Vec2(b)
	return *(*[2]float32)(unsafe.Pointer(&v))
}

func (floatComponent) PutVec2(b []byte, v [2]float32) {
	UInt.PutVec2(b, *(*[2]uint32)(unsafe.Pointer(&v)))
}

func (floatComponent) Vec3(b []byte) [3]float32 {
	v := UInt.Vec3(b)
	return *(*[3]float32)(unsafe.Pointer(&v))
}

func (floatComponent) PutVec3(b []byte, v [3]float32) {
	UInt.PutVec3(b, *(*[3]uint32)(unsafe.Pointer(&v)))
}

func (floatComponent) Vec4(b []byte) [4]float32 {
	v := UInt.Vec4(b)
	return *(*[4]float32)(unsafe.Pointer(&v))
}

func (floatComponent) PutVec4(b []byte, v [4]float32) {
	UInt.PutVec4(b, *(*[4]uint32)(unsafe.Pointer(&v)))
}

func (floatComponent) Mat2(b []byte) [2][2]float32 {
	v := UInt.Mat2(b)
	return *(*[2][2]float32)(unsafe.Pointer(&v))
}

func (floatComponent) PutMat2(b []byte, v [2][2]float32) {
	UInt.PutMat2(b, *(*[2][2]uint32)(unsafe.Pointer(&v)))
}

func (floatComponent) Mat3(b []byte) [3][3]float32 {
	v := UInt.Mat3(b)
	return *(*[3][3]float32)(unsafe.Pointer(&v))
}

func (floatComponent) PutMat3(b []byte, v [3][3]float32) {
	UInt.PutMat3(b, *(*[3][3]uint32)(unsafe.Pointer(&v)))
}

func (floatComponent) Mat4(b []byte) [4][4]float32 {
	v := UInt.Mat4(b)
	return *(*[4][4]float32)(unsafe.Pointer(&v))
}

func (floatComponent) PutMat4(b []byte, v [4][4]float32) {
	UInt.PutMat4(b, *(*[4][4]uint32)(unsafe.Pointer(&v)))
}
