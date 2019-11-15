package binary

import (
	"unsafe"
)

type byteComponent struct{}
type shortComponent struct{}
type floatComponent struct{}

func (byteComponent) Scalar(b []byte) int8 {
	v := UnsignedByte.Scalar(b)
	return *(*int8)(unsafe.Pointer(&v))
}

func (byteComponent) PutScalar(b []byte, v int8) {
	UnsignedByte.PutScalar(b, *(*uint8)(unsafe.Pointer(&v)))
}

func (byteComponent) Vec2(b []byte) [2]int8 {
	v := UnsignedByte.Vec2(b)
	return *(*[2]int8)(unsafe.Pointer(&v))
}

func (byteComponent) PutVec2(b []byte, v [2]int8) {
	UnsignedByte.PutVec2(b, *(*[2]uint8)(unsafe.Pointer(&v)))
}

func (byteComponent) Vec3(b []byte) [3]int8 {
	v := UnsignedByte.Vec3(b)
	return *(*[3]int8)(unsafe.Pointer(&v))
}

func (byteComponent) PutVec3(b []byte, v [3]int8) {
	UnsignedByte.PutVec3(b, *(*[3]uint8)(unsafe.Pointer(&v)))
}

func (byteComponent) Vec4(b []byte) [4]int8 {
	v := UnsignedByte.Vec4(b)
	return *(*[4]int8)(unsafe.Pointer(&v))
}

func (byteComponent) PutVec4(b []byte, v [4]int8) {
	UnsignedByte.PutVec4(b, *(*[4]uint8)(unsafe.Pointer(&v)))
}

func (byteComponent) Mat2(b []byte) [2][2]int8 {
	v := UnsignedByte.Mat2(b)
	return *(*[2][2]int8)(unsafe.Pointer(&v))
}

func (byteComponent) PutMat2(b []byte, v [2][2]int8) {
	UnsignedByte.PutMat2(b, *(*[2][2]uint8)(unsafe.Pointer(&v)))
}

func (byteComponent) Mat3(b []byte) [3][3]int8 {
	v := UnsignedByte.Mat3(b)
	return *(*[3][3]int8)(unsafe.Pointer(&v))
}

func (byteComponent) PutMat3(b []byte, v [3][3]int8) {
	UnsignedByte.PutMat3(b, *(*[3][3]uint8)(unsafe.Pointer(&v)))
}

func (byteComponent) Mat4(b []byte) [4][4]int8 {
	v := UnsignedByte.Mat3(b)
	return *(*[4][4]int8)(unsafe.Pointer(&v))
}

func (byteComponent) PutMat4(b []byte, v [4][4]int8) {
	UnsignedByte.PutMat4(b, *(*[4][4]uint8)(unsafe.Pointer(&v)))
}

func (shortComponent) Scalar(b []byte) int16 {
	v := UnsignedShort.Scalar(b)
	return *(*int16)(unsafe.Pointer(&v))
}

func (shortComponent) PutScalar(b []byte, v int16) {
	UnsignedShort.PutScalar(b, *(*uint16)(unsafe.Pointer(&v)))
}

func (shortComponent) Vec2(b []byte) [2]int16 {
	v := UnsignedShort.Vec2(b)
	return *(*[2]int16)(unsafe.Pointer(&v))
}

func (shortComponent) PutVec2(b []byte, v [2]int16) {
	UnsignedShort.PutVec2(b, *(*[2]uint16)(unsafe.Pointer(&v)))
}

func (shortComponent) Vec3(b []byte) [3]int16 {
	v := UnsignedShort.Vec3(b)
	return *(*[3]int16)(unsafe.Pointer(&v))
}

func (shortComponent) PutVec3(b []byte, v [3]int16) {
	UnsignedShort.PutVec3(b, *(*[3]uint16)(unsafe.Pointer(&v)))
}

func (shortComponent) Vec4(b []byte) [4]int16 {
	v := UnsignedShort.Vec4(b)
	return *(*[4]int16)(unsafe.Pointer(&v))
}

func (shortComponent) PutVec4(b []byte, v [4]int16) {
	UnsignedShort.PutVec4(b, *(*[4]uint16)(unsafe.Pointer(&v)))
}

func (shortComponent) Mat2(b []byte) [2][2]int16 {
	v := UnsignedShort.Mat2(b)
	return *(*[2][2]int16)(unsafe.Pointer(&v))
}

func (shortComponent) PutMat2(b []byte, v [2][2]int16) {
	UnsignedShort.PutMat2(b, *(*[2][2]uint16)(unsafe.Pointer(&v)))
}

func (shortComponent) Mat3(b []byte) [3][3]int16 {
	v := UnsignedShort.Mat3(b)
	return *(*[3][3]int16)(unsafe.Pointer(&v))
}

func (shortComponent) PutMat3(b []byte, v [3][3]int16) {
	UnsignedShort.PutMat3(b, *(*[3][3]uint16)(unsafe.Pointer(&v)))
}

func (shortComponent) Mat4(b []byte) [4][4]int16 {
	v := UnsignedShort.Mat4(b)
	return *(*[4][4]int16)(unsafe.Pointer(&v))
}

func (shortComponent) PutMat4(b []byte, v [4][4]int16) {
	UnsignedShort.PutMat4(b, *(*[4][4]uint16)(unsafe.Pointer(&v)))
}

func (floatComponent) Scalar(b []byte) float32 {
	v := UnsignedInt.Scalar(b)
	return *(*float32)(unsafe.Pointer(&v))
}

func (floatComponent) PutScalar(b []byte, v float32) {
	UnsignedInt.PutScalar(b, *(*uint32)(unsafe.Pointer(&v)))
}

func (floatComponent) Vec2(b []byte) [2]float32 {
	v := UnsignedInt.Vec2(b)
	return *(*[2]float32)(unsafe.Pointer(&v))
}

func (floatComponent) PutVec2(b []byte, v [2]float32) {
	UnsignedInt.PutVec2(b, *(*[2]uint32)(unsafe.Pointer(&v)))
}

func (floatComponent) Vec3(b []byte) [3]float32 {
	v := UnsignedInt.Vec3(b)
	return *(*[3]float32)(unsafe.Pointer(&v))
}

func (floatComponent) PutVec3(b []byte, v [3]float32) {
	UnsignedInt.PutVec3(b, *(*[3]uint32)(unsafe.Pointer(&v)))
}

func (floatComponent) Vec4(b []byte) [4]float32 {
	v := UnsignedInt.Vec4(b)
	return *(*[4]float32)(unsafe.Pointer(&v))
}

func (floatComponent) PutVec4(b []byte, v [4]float32) {
	UnsignedInt.PutVec4(b, *(*[4]uint32)(unsafe.Pointer(&v)))
}

func (floatComponent) Mat2(b []byte) [2][2]float32 {
	v := UnsignedInt.Mat2(b)
	return *(*[2][2]float32)(unsafe.Pointer(&v))
}

func (floatComponent) PutMat2(b []byte, v [2][2]float32) {
	UnsignedInt.PutMat2(b, *(*[2][2]uint32)(unsafe.Pointer(&v)))
}

func (floatComponent) Mat3(b []byte) [3][3]float32 {
	v := UnsignedInt.Mat3(b)
	return *(*[3][3]float32)(unsafe.Pointer(&v))
}

func (floatComponent) PutMat3(b []byte, v [3][3]float32) {
	UnsignedInt.PutMat3(b, *(*[3][3]uint32)(unsafe.Pointer(&v)))
}

func (floatComponent) Mat4(b []byte) [4][4]float32 {
	v := UnsignedInt.Mat4(b)
	return *(*[4][4]float32)(unsafe.Pointer(&v))
}

func (floatComponent) PutMat4(b []byte, v [4][4]float32) {
	UnsignedInt.PutMat4(b, *(*[4][4]uint32)(unsafe.Pointer(&v)))
}
