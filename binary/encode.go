package binary

import (
	"bytes"
	"encoding/binary"
	"image/color"
	"io"

	"github.com/qmuntal/gltf"
)

// Read reads structured binary data from r into data.
// Data should be a slice of glTF predefined fixed-size types,
// else it fallbacks to `encoding/binary.Read`.
//
// If data length is greater than the length of b, Read returns io.ErrShortBuffer.
func Read(b []byte, data interface{}) error {
	c, t, n := Type(data)
	if n <= 0 {
		return binary.Read(bytes.NewReader(b), binary.LittleEndian, data)
	}
	e := SizeOfElement(c, t)
	if len(b) < n*e {
		return io.ErrShortBuffer
	}
	switch data := data.(type) {
	case []color.RGBA:
		for i := range data {
			c := UnsignedByte.Vec4(b[e*i:])
			data[i] = color.RGBA{R: c[0], G: c[1], B: c[2], A: c[3]}
		}
	case []color.RGBA64:
		for i := range data {
			c := UnsignedShort.Vec4(b[e*i:])
			data[i] = color.RGBA64{R: c[0], G: c[1], B: c[2], A: c[3]}
		}
	case []gltf.RGBA:
		for i := range data {
			c := Float.Vec4(b[e*i:])
			data[i] = gltf.RGBA{R: float64(c[0]), G: float64(c[1]), B: float64(c[2]), A: float64(c[3])}
		}
	case []gltf.RGB:
		for i := range data {
			c := Float.Vec3(b[e*i:])
			data[i] = gltf.RGB{R: float64(c[0]), G: float64(c[1]), B: float64(c[2])}
		}
	case []int8:
		for i, x := range b {
			data[i] = int8(x)
		}
	case [][2]int8:
		for i := range data {
			data[i] = Byte.Vec2(b[e*i:])
		}
	case [][3]int8:
		for i := range data {
			data[i] = Byte.Vec3(b[e*i:])
		}
	case [][4]int8:
		for i := range data {
			data[i] = Byte.Vec4(b[e*i:])
		}
	case [][2][2]int8:
		for i := range data {
			data[i] = Byte.Mat2(b[e*i:])
		}
	case [][3][3]int8:
		for i := range data {
			data[i] = Byte.Mat3(b[e*i:])
		}
	case [][4][4]int8:
		for i := range data {
			data[i] = Byte.Mat4(b[e*i:])
		}
	case []uint8:
		copy(data, b)
	case [][2]uint8:
		for i := range data {
			data[i] = UnsignedByte.Vec2(b[e*i:])
		}
	case [][3]uint8:
		for i := range data {
			data[i] = UnsignedByte.Vec3(b[e*i:])
		}
	case [][4]uint8:
		for i := range data {
			data[i] = UnsignedByte.Vec4(b[e*i:])
		}
	case [][2][2]uint8:
		for i := range data {
			data[i] = UnsignedByte.Mat2(b[e*i:])
		}
	case [][3][3]uint8:
		for i := range data {
			data[i] = UnsignedByte.Mat3(b[e*i:])
		}
	case [][4][4]uint8:
		for i := range data {
			data[i] = UnsignedByte.Mat4(b[e*i:])
		}
	case []int16:
		for i := range data {
			data[i] = Short.Scalar(b[e*i:])
		}
	case [][2]int16:
		for i := range data {
			data[i] = Short.Vec2(b[e*i:])
		}
	case [][3]int16:
		for i := range data {
			data[i] = Short.Vec3(b[e*i:])
		}
	case [][4]int16:
		for i := range data {
			data[i] = Short.Vec4(b[e*i:])
		}
	case [][2][2]int16:
		for i := range data {
			data[i] = Short.Mat2(b[e*i:])
		}
	case [][3][3]int16:
		for i := range data {
			data[i] = Short.Mat3(b[e*i:])
		}
	case [][4][4]int16:
		for i := range data {
			data[i] = Short.Mat4(b[e*i:])
		}
	case []uint16:
		for i := range data {
			data[i] = UnsignedShort.Scalar(b[e*i:])
		}
	case [][2]uint16:
		for i := range data {
			data[i] = UnsignedShort.Vec2(b[e*i:])
		}
	case [][3]uint16:
		for i := range data {
			data[i] = UnsignedShort.Vec3(b[e*i:])
		}
	case [][4]uint16:
		for i := range data {
			data[i] = UnsignedShort.Vec4(b[e*i:])
		}
	case [][2][2]uint16:
		for i := range data {
			data[i] = UnsignedShort.Mat2(b[e*i:])
		}
	case [][3][3]uint16:
		for i := range data {
			data[i] = UnsignedShort.Mat3(b[e*i:])
		}
	case [][4][4]uint16:
		for i := range data {
			data[i] = UnsignedShort.Mat4(b[e*i:])
		}
	case []float32:
		for i := range data {
			data[i] = Float.Scalar(b[e*i:])
		}
	case [][2]float32:
		for i := range data {
			data[i] = Float.Vec2(b[e*i:])
		}
	case [][3]float32:
		for i := range data {
			data[i] = Float.Vec3(b[e*i:])
		}
	case [][4]float32:
		for i := range data {
			data[i] = Float.Vec4(b[e*i:])
		}
	case [][2][2]float32:
		for i := range data {
			data[i] = Float.Mat2(b[e*i:])
		}
	case [][3][3]float32:
		for i := range data {
			data[i] = Float.Mat3(b[e*i:])
		}
	case [][4][4]float32:
		for i := range data {
			data[i] = Float.Mat4(b[e*i:])
		}
	case []uint32:
		for i := range data {
			data[i] = UnsignedInt.Scalar(b[e*i:])
		}
	case [][2]uint32:
		for i := range data {
			data[i] = UnsignedInt.Vec2(b[e*i:])
		}
	case [][3]uint32:
		for i := range data {
			data[i] = UnsignedInt.Vec3(b[e*i:])
		}
	case [][4]uint32:
		for i := range data {
			data[i] = UnsignedInt.Vec4(b[e*i:])
		}
	case [][2][2]uint32:
		for i := range data {
			data[i] = UnsignedInt.Mat2(b[e*i:])
		}
	case [][3][3]uint32:
		for i := range data {
			data[i] = UnsignedInt.Mat3(b[e*i:])
		}
	case [][4][4]uint32:
		for i := range data {
			data[i] = UnsignedInt.Mat4(b[e*i:])
		}
	}
	return nil
}

// Write writes the binary representation of data into b.
// Data must be a slice of glTF predefined fixed-size types,
// else it fallbacks to `encoding/binary.Write`.
func Write(b []byte, data interface{}) error {
	c, t, n := Type(data)
	if n <= 0 {
		return binary.Write(bytes.NewBuffer(b), binary.LittleEndian, data)
	}
	e := SizeOfElement(c, t)
	if len(b) < e*n {
		return io.ErrShortBuffer
	}
	switch data := data.(type) {
	case []color.RGBA:
		for i, x := range data {
			UnsignedByte.PutVec4(b[e*i:], [4]uint8{x.R, x.G, x.B, x.A})
		}
	case []color.RGBA64:
		for i, x := range data {
			UnsignedShort.PutVec4(b[e*i:], [4]uint16{x.R, x.G, x.B, x.A})
		}
	case []gltf.RGBA:
		for i, x := range data {
			Float.PutVec4(b[e*i:], [4]float32{float32(x.R), float32(x.G), float32(x.B), float32(x.A)})
		}
	case []gltf.RGB:
		for i, x := range data {
			Float.PutVec3(b[e*i:], [3]float32{float32(x.R), float32(x.G), float32(x.B)})
		}
	case []int8:
		for i, x := range data {
			b[i] = byte(x)
		}
	case [][2]int8:
		for i := range data {
			Byte.PutVec2(b[e*i:], data[i])
		}
	case [][3]int8:
		for i := range data {
			Byte.PutVec3(b[e*i:], data[i])
		}
	case [][4]int8:
		for i := range data {
			Byte.PutVec4(b[e*i:], data[i])
		}
	case [][2][2]int8:
		for i := range data {
			Byte.PutMat2(b[e*i:], data[i])
		}
	case [][3][3]int8:
		for i := range data {
			Byte.PutMat3(b[e*i:], data[i])
		}
	case [][4][4]int8:
		for i := range data {
			Byte.PutMat4(b[e*i:], data[i])
		}
	case []uint8:
		copy(b, data)
	case [][2]uint8:
		for i := range data {
			UnsignedByte.PutVec2(b[e*i:], data[i])
		}
	case [][3]uint8:
		for i := range data {
			UnsignedByte.PutVec3(b[e*i:], data[i])
		}
	case [][4]uint8:
		for i := range data {
			UnsignedByte.PutVec4(b[e*i:], data[i])
		}
	case [][2][2]uint8:
		for i := range data {
			UnsignedByte.PutMat2(b[e*i:], data[i])
		}
	case [][3][3]uint8:
		for i := range data {
			UnsignedByte.PutMat3(b[e*i:], data[i])
		}
	case [][4][4]uint8:
		for i := range data {
			UnsignedByte.PutMat4(b[e*i:], data[i])
		}
	case []int16:
		for i := range data {
			Short.PutScalar(b[e*i:], data[i])
		}
	case [][2]int16:
		for i := range data {
			Short.PutVec2(b[e*i:], data[i])
		}
	case [][3]int16:
		for i := range data {
			Short.PutVec3(b[e*i:], data[i])
		}
	case [][4]int16:
		for i := range data {
			Short.PutVec4(b[e*i:], data[i])
		}
	case [][2][2]int16:
		for i := range data {
			Short.PutMat2(b[e*i:], data[i])
		}
	case [][3][3]int16:
		for i := range data {
			Short.PutMat3(b[e*i:], data[i])
		}
	case [][4][4]int16:
		for i := range data {
			Short.PutMat4(b[e*i:], data[i])
		}
	case []uint16:
		for i := range data {
			UnsignedShort.PutScalar(b[e*i:], data[i])
		}
	case [][2]uint16:
		for i := range data {
			UnsignedShort.PutVec2(b[e*i:], data[i])
		}
	case [][3]uint16:
		for i := range data {
			UnsignedShort.PutVec3(b[e*i:], data[i])
		}
	case [][4]uint16:
		for i := range data {
			UnsignedShort.PutVec4(b[e*i:], data[i])
		}
	case [][2][2]uint16:
		for i := range data {
			UnsignedShort.PutMat2(b[e*i:], data[i])
		}
	case [][3][3]uint16:
		for i := range data {
			UnsignedShort.PutMat3(b[e*i:], data[i])
		}
	case [][4][4]uint16:
		for i := range data {
			UnsignedShort.PutMat4(b[e*i:], data[i])
		}
	case []float32:
		for i := range data {
			Float.PutScalar(b[e*i:], data[i])
		}
	case [][2]float32:
		for i := range data {
			Float.PutVec2(b[e*i:], data[i])
		}
	case [][3]float32:
		for i := range data {
			Float.PutVec3(b[e*i:], data[i])
		}
	case [][4]float32:
		for i := range data {
			Float.PutVec4(b[e*i:], data[i])
		}
	case [][2][2]float32:
		for i := range data {
			Float.PutMat2(b[e*i:], data[i])
		}
	case [][3][3]float32:
		for i := range data {
			Float.PutMat3(b[e*i:], data[i])
		}
	case [][4][4]float32:
		for i := range data {
			Float.PutMat4(b[e*i:], data[i])
		}
	case []uint32:
		for i := range data {
			UnsignedInt.PutScalar(b[e*i:], data[i])
		}
	case [][2]uint32:
		for i := range data {
			UnsignedInt.PutVec2(b[e*i:], data[i])
		}
	case [][3]uint32:
		for i := range data {
			UnsignedInt.PutVec3(b[e*i:], data[i])
		}
	case [][4]uint32:
		for i := range data {
			UnsignedInt.PutVec4(b[e*i:], data[i])
		}
	case [][2][2]uint32:
		for i := range data {
			UnsignedInt.PutMat2(b[e*i:], data[i])
		}
	case [][3][3]uint32:
		for i := range data {
			UnsignedInt.PutMat3(b[e*i:], data[i])
		}
	case [][4][4]uint32:
		for i := range data {
			UnsignedInt.PutMat4(b[e*i:], data[i])
		}
	}
	return nil
}