package binary

import (
	"encoding/binary"
	"math"
)

func (byteComponent) PutScalar(b []byte, v int8) {
	b[0] = int8bits(v)
}

func (byteComponent) PutVec2(b []byte, v [2]int8) {
	b[0] = int8bits(v[0])
	b[1] = int8bits(v[1])
}

func (byteComponent) PutVec3(b []byte, v [3]int8) {
	b[0] = int8bits(v[0])
	b[1] = int8bits(v[1])
	b[2] = int8bits(v[2])
}

func (byteComponent) PutVec4(b []byte, v [4]int8) {
	b[0] = int8bits(v[0])
	b[1] = int8bits(v[1])
	b[2] = int8bits(v[2])
	b[3] = int8bits(v[3])
}

func (byteComponent) PutMat2(b []byte, v [2][2]int8) {
	b[0] = int8bits(v[0][0])
	b[1] = int8bits(v[1][0])
	b[4] = int8bits(v[0][1])
	b[5] = int8bits(v[1][1])
}

func (byteComponent) PutMat3(b []byte, v [3][3]int8) {
	b[0] = int8bits(v[0][0])
	b[1] = int8bits(v[1][0])
	b[2] = int8bits(v[2][0])
	b[4] = int8bits(v[0][1])
	b[5] = int8bits(v[1][1])
	b[6] = int8bits(v[2][1])
	b[8] = int8bits(v[0][2])
	b[9] = int8bits(v[1][2])
	b[10] = int8bits(v[2][2])
}

func (byteComponent) PutMat4(b []byte, v [4][4]int8) {
	b[0] = int8bits(v[0][0])
	b[1] = int8bits(v[1][0])
	b[2] = int8bits(v[2][0])
	b[3] = int8bits(v[3][0])
	b[4] = int8bits(v[0][1])
	b[5] = int8bits(v[1][1])
	b[6] = int8bits(v[2][1])
	b[7] = int8bits(v[3][1])
	b[8] = int8bits(v[0][2])
	b[9] = int8bits(v[1][2])
	b[10] = int8bits(v[2][2])
	b[11] = int8bits(v[3][2])
	b[12] = int8bits(v[0][3])
	b[13] = int8bits(v[1][3])
	b[14] = int8bits(v[2][3])
	b[15] = int8bits(v[3][3])
}

func (ubyteComponent) PutScalar(b []byte, v uint8) {
	b[0] = v
}

func (ubyteComponent) PutVec2(b []byte, v [2]uint8) {
	b[0] = v[0]
	b[1] = v[1]
}

func (ubyteComponent) PutVec3(b []byte, v [3]uint8) {
	b[0] = v[0]
	b[1] = v[1]
	b[2] = v[2]
}

func (ubyteComponent) PutVec4(b []byte, v [4]uint8) {
	b[0] = v[0]
	b[1] = v[1]
	b[2] = v[2]
	b[3] = v[3]
}

func (ubyteComponent) PutMat2(b []byte, v [2][2]uint8) {
	b[0] = v[0][0]
	b[1] = v[1][0]
	b[4] = v[0][1]
	b[5] = v[1][1]
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

func putInt16(b []byte, v int16) {
	binary.LittleEndian.PutUint16(b, int16bits(v))
}

func (shortComponent) PutScalar(b []byte, v int16) {
	putInt16(b, v)
}

func (shortComponent) PutVec2(b []byte, v [2]int16) {
	putInt16(b, v[0])
	putInt16(b[2:], v[1])
}

func (shortComponent) PutVec3(b []byte, v [3]int16) {
	putInt16(b, v[0])
	putInt16(b[2:], v[1])
	putInt16(b[4:], v[2])
}

func (shortComponent) PutVec4(b []byte, v [4]int16) {
	putInt16(b, v[0])
	putInt16(b[2:], v[1])
	putInt16(b[4:], v[2])
	putInt16(b[6:], v[3])
}

func (shortComponent) PutMat2(b []byte, v [2][2]int16) {
	putInt16(b, v[0][0])
	putInt16(b[2:], v[1][0])
	putInt16(b[4:], v[0][1])
	putInt16(b[6:], v[1][1])
}

func (shortComponent) PutMat3(b []byte, v [3][3]int16) {
	putInt16(b, v[0][0])
	putInt16(b[2:], v[1][0])
	putInt16(b[4:], v[2][0])
	putInt16(b[8:], v[0][1])
	putInt16(b[10:], v[1][1])
	putInt16(b[12:], v[2][1])
	putInt16(b[16:], v[0][2])
	putInt16(b[18:], v[1][2])
	putInt16(b[20:], v[2][2])
}

func (shortComponent) PutMat4(b []byte, v [4][4]int16) {
	putInt16(b, v[0][0])
	putInt16(b[2:], v[1][0])
	putInt16(b[4:], v[2][0])
	putInt16(b[6:], v[3][0])
	putInt16(b[8:], v[0][1])
	putInt16(b[10:], v[1][1])
	putInt16(b[12:], v[2][1])
	putInt16(b[14:], v[3][1])
	putInt16(b[16:], v[0][2])
	putInt16(b[18:], v[1][2])
	putInt16(b[20:], v[2][2])
	putInt16(b[22:], v[3][2])
	putInt16(b[24:], v[0][3])
	putInt16(b[26:], v[1][3])
	putInt16(b[28:], v[2][3])
	putInt16(b[30:], v[3][3])
}

func (ushortComponent) PutScalar(b []byte, v uint16) {
	binary.LittleEndian.PutUint16(b, v)
}

func (ushortComponent) PutVec2(b []byte, v [2]uint16) {
	binary.LittleEndian.PutUint16(b, v[0])
	binary.LittleEndian.PutUint16(b[2:], v[1])
}

func (ushortComponent) PutVec3(b []byte, v [3]uint16) {
	binary.LittleEndian.PutUint16(b, v[0])
	binary.LittleEndian.PutUint16(b[2:], v[1])
	binary.LittleEndian.PutUint16(b[4:], v[2])
}

func (ushortComponent) PutVec4(b []byte, v [4]uint16) {
	binary.LittleEndian.PutUint16(b, v[0])
	binary.LittleEndian.PutUint16(b[2:], v[1])
	binary.LittleEndian.PutUint16(b[4:], v[2])
	binary.LittleEndian.PutUint16(b[6:], v[3])
}

func (ushortComponent) PutMat2(b []byte, v [2][2]uint16) {
	binary.LittleEndian.PutUint16(b, v[0][0])
	binary.LittleEndian.PutUint16(b[2:], v[1][0])
	binary.LittleEndian.PutUint16(b[4:], v[0][1])
	binary.LittleEndian.PutUint16(b[6:], v[1][1])
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

func (uintComponent) PutScalar(b []byte, v uint32) {
	binary.LittleEndian.PutUint32(b, v)
}

func (uintComponent) PutVec2(b []byte, v [2]uint32) {
	binary.LittleEndian.PutUint32(b, v[0])
	binary.LittleEndian.PutUint32(b[4:], v[1])
}

func (uintComponent) PutVec3(b []byte, v [3]uint32) {
	binary.LittleEndian.PutUint32(b, v[0])
	binary.LittleEndian.PutUint32(b[4:], v[1])
	binary.LittleEndian.PutUint32(b[8:], v[2])
}

func (uintComponent) PutVec4(b []byte, v [4]uint32) {
	binary.LittleEndian.PutUint32(b, v[0])
	binary.LittleEndian.PutUint32(b[4:], v[1])
	binary.LittleEndian.PutUint32(b[8:], v[2])
	binary.LittleEndian.PutUint32(b[12:], v[3])
}

func (uintComponent) PutMat2(b []byte, v [2][2]uint32) {
	binary.LittleEndian.PutUint32(b, v[0][0])
	binary.LittleEndian.PutUint32(b[4:], v[1][0])
	binary.LittleEndian.PutUint32(b[8:], v[0][1])
	binary.LittleEndian.PutUint32(b[12:], v[1][1])
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

func putFloat32(b []byte, v float32) {
	binary.LittleEndian.PutUint32(b, math.Float32bits(v))
}

func (floatComponent) PutScalar(b []byte, v float32) {
	putFloat32(b, v)
}

func (floatComponent) PutVec2(b []byte, v [2]float32) {
	putFloat32(b, v[0])
	putFloat32(b[4:], v[1])
}

func (floatComponent) PutVec3(b []byte, v [3]float32) {
	putFloat32(b, v[0])
	putFloat32(b[4:], v[1])
	putFloat32(b[8:], v[2])
}

func (floatComponent) PutVec4(b []byte, v [4]float32) {
	putFloat32(b, v[0])
	putFloat32(b[4:], v[1])
	putFloat32(b[8:], v[2])
	putFloat32(b[12:], v[3])
}

func (floatComponent) PutMat2(b []byte, v [2][2]float32) {
	putFloat32(b, v[0][0])
	putFloat32(b[4:], v[1][0])
	putFloat32(b[8:], v[0][1])
	putFloat32(b[12:], v[1][1])
}

func (floatComponent) PutMat3(b []byte, v [3][3]float32) {
	putFloat32(b, v[0][0])
	putFloat32(b[4:], v[1][0])
	putFloat32(b[8:], v[2][0])
	putFloat32(b[12:], v[0][1])
	putFloat32(b[16:], v[1][1])
	putFloat32(b[20:], v[2][1])
	putFloat32(b[24:], v[0][2])
	putFloat32(b[28:], v[1][2])
	putFloat32(b[32:], v[2][2])
}

func (floatComponent) PutMat4(b []byte, v [4][4]float32) {
	putFloat32(b, v[0][0])
	putFloat32(b[4:], v[1][0])
	putFloat32(b[8:], v[2][0])
	putFloat32(b[12:], v[3][0])
	putFloat32(b[16:], v[0][1])
	putFloat32(b[20:], v[1][1])
	putFloat32(b[24:], v[2][1])
	putFloat32(b[28:], v[3][1])
	putFloat32(b[32:], v[0][2])
	putFloat32(b[36:], v[1][2])
	putFloat32(b[40:], v[2][2])
	putFloat32(b[44:], v[3][2])
	putFloat32(b[48:], v[0][3])
	putFloat32(b[52:], v[1][3])
	putFloat32(b[56:], v[2][3])
	putFloat32(b[60:], v[3][3])
}
