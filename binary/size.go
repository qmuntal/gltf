package binary

import (
	"fmt"
	"image/color"
	"reflect"

	"github.com/qmuntal/gltf"
)

// SizeOfElement returns the size, in bytes, of an element.
// The element size may not be (component size) * (number of components),
// as some of the elements are tightly packed in order to ensure
// that they are aligned to 4-byte boundaries.
func SizeOfElement(c gltf.ComponentType, t gltf.AccessorType) uint32 {
	// special cases
	switch {
	case (t == gltf.AccessorVec3 || t == gltf.AccessorVec2) && (c == gltf.ComponentByte || c == gltf.ComponentUbyte):
		return 4
	case t == gltf.AccessorVec3 && (c == gltf.ComponentShort || c == gltf.ComponentUshort):
		return 8
	case t == gltf.AccessorMat2 && (c == gltf.ComponentByte || c == gltf.ComponentUbyte):
		return 8
	case t == gltf.AccessorMat3 && (c == gltf.ComponentByte || c == gltf.ComponentUbyte):
		return 12
	case t == gltf.AccessorMat3 && (c == gltf.ComponentShort || c == gltf.ComponentUshort):
		return 24
	}
	return c.ByteSize() * t.Components()
}

// MakeSlice returns the slice type associated with c and t and with the given element count.
// For example, if c is gltf.ComponentFloat and t is gltf.AccessorVec3
// then MakeSlice(c, t, 5) is equivalent to make([][3]float32, 5).
func MakeSlice(c gltf.ComponentType, t gltf.AccessorType, count uint32) interface{} {
	var tp reflect.Type
	switch c {
	case gltf.ComponentUbyte:
		tp = reflect.TypeOf((*uint8)(nil))
	case gltf.ComponentByte:
		tp = reflect.TypeOf((*int8)(nil))
	case gltf.ComponentUshort:
		tp = reflect.TypeOf((*uint16)(nil))
	case gltf.ComponentShort:
		tp = reflect.TypeOf((*int16)(nil))
	case gltf.ComponentUint:
		tp = reflect.TypeOf((*uint32)(nil))
	case gltf.ComponentFloat:
		tp = reflect.TypeOf((*float32)(nil))
	}
	tp = tp.Elem()
	switch t {
	case gltf.AccessorVec2:
		tp = reflect.ArrayOf(2, tp)
	case gltf.AccessorVec3:
		tp = reflect.ArrayOf(3, tp)
	case gltf.AccessorVec4:
		tp = reflect.ArrayOf(4, tp)
	case gltf.AccessorMat2:
		tp = reflect.ArrayOf(2, reflect.ArrayOf(2, tp))
	case gltf.AccessorMat3:
		tp = reflect.ArrayOf(3, reflect.ArrayOf(3, tp))
	case gltf.AccessorMat4:
		tp = reflect.ArrayOf(4, reflect.ArrayOf(4, tp))
	}
	return reflect.MakeSlice(reflect.SliceOf(tp), int(count), int(count)).Interface()
}

// Type returns the associated glTF type data.
// It panics if data is not an slice.
func Type(data interface{}) (c gltf.ComponentType, t gltf.AccessorType, count uint32) {
	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Slice {
		panic(fmt.Sprintf("go3mf: binary.Type expecting a slice but got %s", v.Kind()))
	}
	count = uint32(v.Len())
	switch data.(type) {
	case []int8:
		c, t = gltf.ComponentByte, gltf.AccessorScalar
	case [][2]int8:
		c, t = gltf.ComponentByte, gltf.AccessorVec2
	case [][3]int8:
		c, t = gltf.ComponentByte, gltf.AccessorVec3
	case [][4]int8:
		c, t = gltf.ComponentByte, gltf.AccessorVec4
	case [][2][2]int8:
		c, t = gltf.ComponentByte, gltf.AccessorMat2
	case [][3][3]int8:
		c, t = gltf.ComponentByte, gltf.AccessorMat3
	case [][4][4]int8:
		c, t = gltf.ComponentByte, gltf.AccessorMat4
	case []uint8:
		c, t = gltf.ComponentUbyte, gltf.AccessorScalar
	case [][2]uint8:
		c, t = gltf.ComponentUbyte, gltf.AccessorVec2
	case [][3]uint8:
		c, t = gltf.ComponentUbyte, gltf.AccessorVec3
	case []color.RGBA, [][4]uint8:
		c, t = gltf.ComponentUbyte, gltf.AccessorVec4
	case [][2][2]uint8:
		c, t = gltf.ComponentUbyte, gltf.AccessorMat2
	case [][3][3]uint8:
		c, t = gltf.ComponentUbyte, gltf.AccessorMat3
	case [][4][4]uint8:
		c, t = gltf.ComponentUbyte, gltf.AccessorMat4
	case []int16:
		c, t = gltf.ComponentShort, gltf.AccessorScalar
	case [][2]int16:
		c, t = gltf.ComponentShort, gltf.AccessorVec2
	case [][3]int16:
		c, t = gltf.ComponentShort, gltf.AccessorVec3
	case [][4]int16:
		c, t = gltf.ComponentShort, gltf.AccessorVec4
	case [][2][2]int16:
		c, t = gltf.ComponentShort, gltf.AccessorMat2
	case [][3][3]int16:
		c, t = gltf.ComponentShort, gltf.AccessorMat3
	case [][4][4]int16:
		c, t = gltf.ComponentShort, gltf.AccessorMat4
	case []uint16:
		c, t = gltf.ComponentUshort, gltf.AccessorScalar
	case [][2]uint16:
		c, t = gltf.ComponentUshort, gltf.AccessorVec2
	case [][3]uint16:
		c, t = gltf.ComponentUshort, gltf.AccessorVec3
	case []color.RGBA64, [][4]uint16:
		c, t = gltf.ComponentUshort, gltf.AccessorVec4
	case [][2][2]uint16:
		c, t = gltf.ComponentUshort, gltf.AccessorMat2
	case [][3][3]uint16:
		c, t = gltf.ComponentUshort, gltf.AccessorMat3
	case [][4][4]uint16:
		c, t = gltf.ComponentUshort, gltf.AccessorMat4
	case []uint32:
		c, t = gltf.ComponentUint, gltf.AccessorScalar
	case [][2]uint32:
		c, t = gltf.ComponentUint, gltf.AccessorVec2
	case [][3]uint32:
		c, t = gltf.ComponentUint, gltf.AccessorVec3
	case [][4]uint32:
		c, t = gltf.ComponentUint, gltf.AccessorVec4
	case [][2][2]uint32:
		c, t = gltf.ComponentUint, gltf.AccessorMat2
	case [][3][3]uint32:
		c, t = gltf.ComponentUint, gltf.AccessorMat3
	case [][4][4]uint32:
		c, t = gltf.ComponentUint, gltf.AccessorMat4
	case []float32:
		c, t = gltf.ComponentFloat, gltf.AccessorScalar
	case [][2]float32:
		c, t = gltf.ComponentFloat, gltf.AccessorVec2
	case [][3]float32:
		c, t = gltf.ComponentFloat, gltf.AccessorVec3
	case [][4]float32:
		c, t = gltf.ComponentFloat, gltf.AccessorVec4
	case [][2][2]float32:
		c, t = gltf.ComponentFloat, gltf.AccessorMat2
	case [][3][3]float32:
		c, t = gltf.ComponentFloat, gltf.AccessorMat3
	case [][4][4]float32:
		c, t = gltf.ComponentFloat, gltf.AccessorMat4
	default:
		panic(fmt.Sprintf("go3mf: binary.Type expecting a glTF supported type but got %s", v.Kind()))
	}
	return
}
