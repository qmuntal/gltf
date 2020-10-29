package gltf

import (
	"math"
	"testing"

	"github.com/go-test/deep"
)

func TestNormalize(t *testing.T) {
	if got := NormalizeByte(0.0236); got != 3 {
		t.Errorf("NormalizeByte = %v, want %v", got, 3)
	}
	if got := NormalizeUbyte(0.2); got != 51 {
		t.Errorf("NormalizeUbyte = %v, want %v", got, 51)
	}
	if got := NormalizeShort(0.03053); got != 1000 {
		t.Errorf("NormalizeShort = %v, want %v", got, 1000)
	}
	if got := NormalizeUshort(0.2); got != 13107 {
		t.Errorf("NormalizeUshort = %v, want %v", got, 13107)
	}
}

func TestDenormalize(t *testing.T) {
	if got := DenormalizeByte(3); math.Abs(float64(got)-0.0236) > 1e-4 {
		t.Errorf("DenormalizeByte = %v, want %v", got, 0.0236)
	}
	if got := DenormalizeUbyte(51); got != 0.2 {
		t.Errorf("DenormalizeUbyte = %v, want %v", got, 0.2)
	}
	if got := DenormalizeShort(1000); math.Abs(float64(got)-0.03053) > 1e-4 {
		t.Errorf("DenormalizeShort = %v, want %v", got, 0.03053)
	}
	if got := DenormalizeUshort(13107); got != 0.2 {
		t.Errorf("DenormalizeUshort = %v, want %v", got, 0.2)
	}
}

func TestDenormalizeRGBA(t *testing.T) {
	deep.FloatPrecision = 6
	type args struct {
		v [4]uint8
	}
	tests := []struct {
		name string
		args args
		want [4]float32
	}{
		{"empty", args{[4]uint8{}}, [4]float32{}},
		{"base", args{[4]uint8{1, 1, 1, 1}}, [4]float32{0.0003035, 0.0003035, 0.0003035, 0.00392156}},
		{"max", args{[4]uint8{255, 255, 255, 255}}, [4]float32{1, 1, 1, 1}},
		{"other", args{[4]uint8{60, 120, 180, 220}}, [4]float32{0.045186, 0.1878207, 0.4564110, 0.86274509}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DenormalizeRGBA(tt.args.v)
			if diff := deep.Equal(got, tt.want); diff != nil {
				t.Errorf("DenormalizeRGBA() = %v", diff)
			}
		})
	}
}
