package binary

import (
	"fmt"
	"image/color"
	"reflect"

	"github.com/qmuntal/gltf"
)

// MakeSliceBuffer returns the slice type associated with c and t and with the given element count.
// If the buffer is an slice which type matches with the expected by the acr then it will
// be used as backing slice.
func MakeSliceBuffer(c gltf.ComponentType, t gltf.AccessorType, count int, buffer any) any {
	if buffer == nil {
		return MakeSlice(c, t, count)
	}
	c1, t1, count1 := Type(buffer)
	if count1 == 0 || c1 != c || t1 != t {
		return MakeSlice(c, t, count)
	}
	if count1 < count {
		tmpSlice := MakeSlice(c, t, count-count1)
		return reflect.AppendSlice(reflect.ValueOf(buffer), reflect.ValueOf(tmpSlice)).Interface()
	}
	if count1 > count {
		return reflect.ValueOf(buffer).Slice(0, int(count)).Interface()
	}
	return buffer
}

// MakeSlice returns the slice type associated with c and t and with the given element count.
// For example, if c is gltf.ComponentFloat and t is gltf.AccessorVec3
// then MakeSlice(c, t, 5) is equivalent to make([][3]float32, 5).
func MakeSlice(c gltf.ComponentType, t gltf.AccessorType, count int) any {
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
	return reflect.MakeSlice(reflect.SliceOf(tp), count, count).Interface()
}

// Type returns the associated glTF type data.
// It panics if data is not an slice.
func Type(data any) (c gltf.ComponentType, t gltf.AccessorType, count int) {
	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Slice {
		panic(fmt.Sprintf("gltf: binary.Type expecting a slice but got %s", v.Kind()))
	}
	count = v.Len()
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
		panic(fmt.Sprintf("gltf: binary.Type expecting a glTF supported type but got %s", v.Kind()))
	}
	return
}
