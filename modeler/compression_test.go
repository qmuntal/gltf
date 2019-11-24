package modeler

import (
	"math"
	"reflect"
	"testing"
)

func Test_compressUint32(t *testing.T) {
	type args struct {
		data []uint32
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{"empty", args{[]uint32{}}, []uint8{}},
		{"32", args{[]uint32{1, math.MaxUint16}}, []uint32{1, math.MaxUint16}},
		{"32-2", args{[]uint32{1, math.MaxUint8, math.MaxUint16}}, []uint32{1, math.MaxUint8, math.MaxUint16}},
		{"32-3", args{[]uint32{1, math.MaxUint16, math.MaxUint8}}, []uint32{1, math.MaxUint16, math.MaxUint8}},
		{"32-4", args{[]uint32{1, math.MaxUint8}}, []uint16{1, math.MaxUint8}},
		{"16", args{[]uint32{1, math.MaxUint16 - 1}}, []uint16{1, math.MaxUint16 - 1}},
		{"8", args{[]uint32{1, math.MaxUint8 - 1}}, []uint8{1, math.MaxUint8 - 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := compressUint32(tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("compressUint32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_compressUint16(t *testing.T) {
	type args struct {
		data []uint16
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{"empty", args{[]uint16{}}, []uint8{}},
		{"16", args{[]uint16{1, math.MaxUint8}}, []uint16{1, math.MaxUint8}},
		{"8", args{[]uint16{1, math.MaxUint8 - 1}}, []uint8{1, math.MaxUint8 - 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := compressUint16(tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("compressUint16() = %v, want %v", got, tt.want)
			}
		})
	}
}
