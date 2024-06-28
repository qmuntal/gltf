package binary

import (
	"reflect"
	"testing"

	"github.com/qmuntal/gltf"
)

func TestMakeSlice(t *testing.T) {
	type args struct {
		c     gltf.ComponentType
		t     gltf.AccessorType
		count int
	}
	tests := []struct {
		name string
		args args
		want any
	}{
		// Scalar
		{"[]uint8", args{gltf.ComponentUbyte, gltf.AccessorScalar, 5}, make([]uint8, 5)},
		{"[]int8", args{gltf.ComponentByte, gltf.AccessorScalar, 5}, make([]int8, 5)},
		{"[]uint16", args{gltf.ComponentUshort, gltf.AccessorScalar, 5}, make([]uint16, 5)},
		{"[]int16", args{gltf.ComponentShort, gltf.AccessorScalar, 5}, make([]int16, 5)},
		{"[]uint32", args{gltf.ComponentUint, gltf.AccessorScalar, 5}, make([]uint32, 5)},
		{"[]float32", args{gltf.ComponentFloat, gltf.AccessorScalar, 5}, make([]float32, 5)},
		// Vec2
		{"[][2]uint8", args{gltf.ComponentUbyte, gltf.AccessorVec2, 5}, make([][2]uint8, 5)},
		{"[][2]int8", args{gltf.ComponentByte, gltf.AccessorVec2, 5}, make([][2]int8, 5)},
		{"[][2]uint16", args{gltf.ComponentUshort, gltf.AccessorVec2, 5}, make([][2]uint16, 5)},
		{"[][2]int16", args{gltf.ComponentShort, gltf.AccessorVec2, 5}, make([][2]int16, 5)},
		{"[][2]uint32", args{gltf.ComponentUint, gltf.AccessorVec2, 5}, make([][2]uint32, 5)},
		{"[][2]float32", args{gltf.ComponentFloat, gltf.AccessorVec2, 5}, make([][2]float32, 5)},
		// Vec3
		{"[][3]uint8", args{gltf.ComponentUbyte, gltf.AccessorVec3, 5}, make([][3]uint8, 5)},
		{"[][3]int8", args{gltf.ComponentByte, gltf.AccessorVec3, 5}, make([][3]int8, 5)},
		{"[][3]uint16", args{gltf.ComponentUshort, gltf.AccessorVec3, 5}, make([][3]uint16, 5)},
		{"[][3]int16", args{gltf.ComponentShort, gltf.AccessorVec3, 5}, make([][3]int16, 5)},
		{"[][3]uint32", args{gltf.ComponentUint, gltf.AccessorVec3, 5}, make([][3]uint32, 5)},
		{"[][3]float32", args{gltf.ComponentFloat, gltf.AccessorVec3, 5}, make([][3]float32, 5)},
		// Vec4
		{"[][4]uint8", args{gltf.ComponentUbyte, gltf.AccessorVec4, 5}, make([][4]uint8, 5)},
		{"[][4]int8", args{gltf.ComponentByte, gltf.AccessorVec4, 5}, make([][4]int8, 5)},
		{"[][4]uint16", args{gltf.ComponentUshort, gltf.AccessorVec4, 5}, make([][4]uint16, 5)},
		{"[][4]int16", args{gltf.ComponentShort, gltf.AccessorVec4, 5}, make([][4]int16, 5)},
		{"[][4]uint32", args{gltf.ComponentUint, gltf.AccessorVec4, 5}, make([][4]uint32, 5)},
		{"[][4]float32", args{gltf.ComponentFloat, gltf.AccessorVec4, 5}, make([][4]float32, 5)},
		// Mat2
		{"[][2][2]uint8", args{gltf.ComponentUbyte, gltf.AccessorMat2, 5}, make([][2][2]uint8, 5)},
		{"[][2][2]int8", args{gltf.ComponentByte, gltf.AccessorMat2, 5}, make([][2][2]int8, 5)},
		{"[][2][2]uint16", args{gltf.ComponentUshort, gltf.AccessorMat2, 5}, make([][2][2]uint16, 5)},
		{"[][2][2]int16", args{gltf.ComponentShort, gltf.AccessorMat2, 5}, make([][2][2]int16, 5)},
		{"[][2][2]uint32", args{gltf.ComponentUint, gltf.AccessorMat2, 5}, make([][2][2]uint32, 5)},
		{"[][2][2]float32", args{gltf.ComponentFloat, gltf.AccessorMat2, 5}, make([][2][2]float32, 5)},
		// Mat3
		{"[][3][3]uint8", args{gltf.ComponentUbyte, gltf.AccessorMat3, 5}, make([][3][3]uint8, 5)},
		{"[][3][3]int8", args{gltf.ComponentByte, gltf.AccessorMat3, 5}, make([][3][3]int8, 5)},
		{"[][3][3]uint16", args{gltf.ComponentUshort, gltf.AccessorMat3, 5}, make([][3][3]uint16, 5)},
		{"[][3][3]int16", args{gltf.ComponentShort, gltf.AccessorMat3, 5}, make([][3][3]int16, 5)},
		{"[][3][3]uint32", args{gltf.ComponentUint, gltf.AccessorMat3, 5}, make([][3][3]uint32, 5)},
		{"[][3][3]float32", args{gltf.ComponentFloat, gltf.AccessorMat3, 5}, make([][3][3]float32, 5)},
		// Mat4
		{"[][4][4]uint8", args{gltf.ComponentUbyte, gltf.AccessorMat4, 5}, make([][4][4]uint8, 5)},
		{"[][4][4]int8", args{gltf.ComponentByte, gltf.AccessorMat4, 5}, make([][4][4]int8, 5)},
		{"[][4][4]uint16", args{gltf.ComponentUshort, gltf.AccessorMat4, 5}, make([][4][4]uint16, 5)},
		{"[][4][4]int16", args{gltf.ComponentShort, gltf.AccessorMat4, 5}, make([][4][4]int16, 5)},
		{"[][4][4]uint32", args{gltf.ComponentUint, gltf.AccessorMat4, 5}, make([][4][4]uint32, 5)},
		{"[][4][4]float32", args{gltf.ComponentFloat, gltf.AccessorMat4, 5}, make([][4][4]float32, 5)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MakeSlice(tt.args.c, tt.args.t, tt.args.count); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MakeSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMakeSliceBuffer(t *testing.T) {
	type args struct {
		c      gltf.ComponentType
		t      gltf.AccessorType
		count  int
		buffer any
	}
	tests := []struct {
		name string
		args args
		want any
	}{
		{"nil buffer", args{gltf.ComponentUbyte, gltf.AccessorVec2, 2, nil}, make([][2]uint8, 2)},
		{"empty buffer", args{gltf.ComponentUbyte, gltf.AccessorVec2, 2, make([][2]uint8, 0)}, make([][2]uint8, 2)},
		{"different buffer", args{gltf.ComponentUbyte, gltf.AccessorVec2, 2, make([][3]int8, 3)}, make([][2]uint8, 2)},
		{"small buffer", args{gltf.ComponentUbyte, gltf.AccessorVec2, 2, make([][2]uint8, 1)}, make([][2]uint8, 2)},
		{"large buffer", args{gltf.ComponentUbyte, gltf.AccessorVec2, 2, make([][2]uint8, 3)}, make([][2]uint8, 2)},
		{"same buffer", args{gltf.ComponentUbyte, gltf.AccessorVec2, 2, make([][2]uint8, 2)}, make([][2]uint8, 2)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MakeSliceBuffer(tt.args.c, tt.args.t, tt.args.count, tt.args.buffer); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MakeSliceBuffer() = %v, want %v", got, tt.want)
			}
		})
	}
}
