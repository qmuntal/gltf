package binary

import (
	"fmt"
	"image/color"
	"reflect"
	"unsafe"

	"github.com/qmuntal/gltf"
)

func reslice[T any](v []byte, length, rem int) []T {
	var ptr unsafe.Pointer
	if len(v) != 0 {
		ptr = unsafe.Pointer(&v[0])
	}
	s := unsafe.Slice((*T)(unsafe.Pointer(ptr)), length)
	if rem > 0 {
		return append(s, make([]T, rem)...)
	}
	return s
}

func castSlice(c gltf.ComponentType, t gltf.AccessorType, count int, v []byte) any {
	l := len(v) / (c.ByteSize() * t.Components())
	var rem int
	if count < l {
		l = count
	} else if count > l {
		rem = count - l
	}

	switch c {
	case gltf.ComponentUbyte:
		switch t {
		case gltf.AccessorScalar:
			return reslice[uint8](v, l, rem)
		case gltf.AccessorVec2:
			return reslice[[2]uint8](v, l, rem)
		case gltf.AccessorVec3:
			return reslice[[3]uint8](v, l, rem)
		case gltf.AccessorVec4:
			return reslice[[4]uint8](v, l, rem)
		case gltf.AccessorMat2:
			return reslice[[2][2]uint8](v, l, rem)
		case gltf.AccessorMat3:
			return reslice[[3][3]uint8](v, l, rem)
		case gltf.AccessorMat4:
			return reslice[[4][4]uint8](v, l, rem)
		}
	case gltf.ComponentByte:
		switch t {
		case gltf.AccessorScalar:
			return reslice[int8](v, l, rem)
		case gltf.AccessorVec2:
			return reslice[[2]int8](v, l, rem)
		case gltf.AccessorVec3:
			return reslice[[3]int8](v, l, rem)
		case gltf.AccessorVec4:
			return reslice[[4]int8](v, l, rem)
		case gltf.AccessorMat2:
			return reslice[[2][2]int8](v, l, rem)
		case gltf.AccessorMat3:
			return reslice[[3][3]int8](v, l, rem)
		case gltf.AccessorMat4:
			return reslice[[4][4]int8](v, l, rem)
		}
	case gltf.ComponentUshort:
		switch t {
		case gltf.AccessorScalar:
			return reslice[uint16](v, l, rem)
		case gltf.AccessorVec2:
			return reslice[[2]uint16](v, l, rem)
		case gltf.AccessorVec3:
			return reslice[[3]uint16](v, l, rem)
		case gltf.AccessorVec4:
			return reslice[[4]uint16](v, l, rem)
		case gltf.AccessorMat2:
			return reslice[[2][2]uint16](v, l, rem)
		case gltf.AccessorMat3:
			return reslice[[3][3]uint16](v, l, rem)
		case gltf.AccessorMat4:
			return reslice[[4][4]uint16](v, l, rem)
		}
	case gltf.ComponentShort:
		switch t {
		case gltf.AccessorScalar:
			return reslice[int16](v, l, rem)
		case gltf.AccessorVec2:
			return reslice[[2]int16](v, l, rem)
		case gltf.AccessorVec3:
			return reslice[[3]int16](v, l, rem)
		case gltf.AccessorVec4:
			return reslice[[4]int16](v, l, rem)
		case gltf.AccessorMat2:
			return reslice[[2][2]int16](v, l, rem)
		case gltf.AccessorMat3:
			return reslice[[3][3]int16](v, l, rem)
		case gltf.AccessorMat4:
			return reslice[[4][4]int16](v, l, rem)
		}
	case gltf.ComponentUint:
		switch t {
		case gltf.AccessorScalar:
			return reslice[uint32](v, l, rem)
		case gltf.AccessorVec2:
			return reslice[[2]uint32](v, l, rem)
		case gltf.AccessorVec3:
			return reslice[[3]uint32](v, l, rem)
		case gltf.AccessorVec4:
			return reslice[[4]uint32](v, l, rem)
		case gltf.AccessorMat2:
			return reslice[[2][2]uint32](v, l, rem)
		case gltf.AccessorMat3:

		case gltf.AccessorMat4:
			return reslice[[4][4]uint32](v, l, rem)
		}
	case gltf.ComponentFloat:
		switch t {
		case gltf.AccessorScalar:
			return reslice[float32](v, l, rem)
		case gltf.AccessorVec2:
			return reslice[[2]float32](v, l, rem)
		case gltf.AccessorVec3:
			return reslice[[3]float32](v, l, rem)
		case gltf.AccessorVec4:
			return reslice[[4]float32](v, l, rem)
		case gltf.AccessorMat2:
			return reslice[[2][2]float32](v, l, rem)
		case gltf.AccessorMat3:
			return reslice[[3][3]float32](v, l, rem)
		case gltf.AccessorMat4:
			return reslice[[4][4]float32](v, l, rem)
		}
	}
	panic(fmt.Errorf("gltf: unsupported component type %d or accessor type %d", c, t))
}

// MakeSliceBuffer returns the slice type associated with c and t and with the given element count.
// If the buffer is an slice which type matches with the expected by the acr then it will
// be used as backing slice.
func MakeSliceBuffer(c gltf.ComponentType, t gltf.AccessorType, count int, buffer []byte) (any, error) {
	if err := checkAccessorType(c, t); err != nil {
		return nil, err
	}
	buffer = buffer[:cap(buffer)] // Extend the slice to its capacity.
	if len(buffer) == 0 {
		return MakeSlice(c, t, count)
	}
	return castSlice(c, t, count, buffer), nil
}

// MakeSlice returns the slice type associated with c and t and with the given element count.
// For example, if c is gltf.ComponentFloat and t is gltf.AccessorVec3
// then MakeSlice(c, t, 5) is equivalent to make([][3]float32, 5).
func MakeSlice(c gltf.ComponentType, t gltf.AccessorType, count int) (any, error) {
	if err := checkAccessorType(c, t); err != nil {
		return nil, err
	}
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
	default:
		panic(fmt.Errorf("gltf: unsupported component type %d", c))
	}
	tp = tp.Elem()
	switch t {
	case gltf.AccessorScalar:
		// Nothing to do.
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
	default:
		panic(fmt.Errorf("gltf: unsupported accessor type %d", t))
	}
	return reflect.MakeSlice(reflect.SliceOf(tp), count, count).Interface(), nil
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

func checkAccessorType(c gltf.ComponentType, t gltf.AccessorType) error {
	switch c {
	case gltf.ComponentUbyte, gltf.ComponentByte, gltf.ComponentUshort,
		gltf.ComponentShort, gltf.ComponentUint, gltf.ComponentFloat:
	default:
		return fmt.Errorf("gltf: unsupported component type %d", c)
	}
	switch t {
	case gltf.AccessorScalar, gltf.AccessorVec2, gltf.AccessorVec3,
		gltf.AccessorVec4, gltf.AccessorMat2, gltf.AccessorMat3, gltf.AccessorMat4:
	default:
		return fmt.Errorf("gltf: unsupported accessor type %d", t)
	}
	return nil
}
