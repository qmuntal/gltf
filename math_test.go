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

func TestNormalizeRGBA(t *testing.T) {
	deep.FloatPrecision = 6
	type args struct {
		v [4]float32
	}
	tests := []struct {
		name string
		args args
		want [4]uint8
	}{
		{"empty", args{[4]float32{}}, [4]uint8{}},
		{"base", args{[4]float32{0.0003035, 0.0003035, 0.0003035, 0.00392156}}, [4]uint8{0, 0, 0, 0}},
		{"max", args{[4]float32{1, 1, 1, 1}}, [4]uint8{255, 255, 255, 255}},
		{"other", args{[4]float32{0.045186, 0.1878207, 0.4564110, 0.86274509}}, [4]uint8{59, 119, 179, 220}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NormalizeRGBA(tt.args.v)
			if diff := deep.Equal(got, tt.want); diff != nil {
				t.Errorf("NormalizeRGBA() = %v", diff)
			}
		})
	}
}

func TestNormalizeRGB(t *testing.T) {
	deep.FloatPrecision = 6
	type args struct {
		v [3]float32
	}
	tests := []struct {
		name string
		args args
		want [3]uint8
	}{
		{"empty", args{[3]float32{}}, [3]uint8{}},
		{"base", args{[3]float32{0.0003035, 0.0003035, 0.0003035}}, [3]uint8{0, 0, 0}},
		{"max", args{[3]float32{1, 1, 1}}, [3]uint8{255, 255, 255}},
		{"other", args{[3]float32{0.045186, 0.1878207, 0.4564110}}, [3]uint8{59, 119, 179}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NormalizeRGB(tt.args.v)
			if diff := deep.Equal(got, tt.want); diff != nil {
				t.Errorf("NormalizeRGB() = %v", diff)
			}
		})
	}
}

func TestNormalizeRGBA64(t *testing.T) {
	deep.FloatPrecision = 6
	type args struct {
		v [4]float32
	}
	tests := []struct {
		name string
		args args
		want [4]uint16
	}{
		{"empty", args{[4]float32{}}, [4]uint16{}},
		{"base", args{[4]float32{0.001181, 0.001181, 0.001181, 0.0152590}}, [4]uint16{999, 999, 999, 999}},
		{"max", args{[4]float32{1, 1, 1, 1}}, [4]uint16{65535, 65535, 65535, 65535}},
		{"other", args{[4]float32{0.0087617, 0.0280836, 0.061314, 0.335698}}, [4]uint16{5999, 11999, 17999, 21999}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NormalizeRGBA64(tt.args.v)
			if diff := deep.Equal(got, tt.want); diff != nil {
				t.Errorf("NormalizeRGBA64() = %v", diff)
			}
		})
	}
}

func TestNormalizeRGB64(t *testing.T) {
	deep.FloatPrecision = 6
	type args struct {
		v [3]float32
	}
	tests := []struct {
		name string
		args args
		want [3]uint16
	}{
		{"empty", args{[3]float32{}}, [3]uint16{}},
		{"base", args{[3]float32{0.001181, 0.001181, 0.001181}}, [3]uint16{999, 999, 999}},
		{"max", args{[3]float32{1, 1, 1}}, [3]uint16{65535, 65535, 65535}},
		{"other", args{[3]float32{0.0087617, 0.0280836, 0.061314}}, [3]uint16{5999, 11999, 17999}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NormalizeRGB64(tt.args.v)
			if diff := deep.Equal(got, tt.want); diff != nil {
				t.Errorf("NormalizeRGB64() = %v", diff)
			}
		})
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

func TestDenormalizeRGB(t *testing.T) {
	deep.FloatPrecision = 6
	type args struct {
		v [3]uint8
	}
	tests := []struct {
		name string
		args args
		want [3]float32
	}{
		{"empty", args{[3]uint8{}}, [3]float32{}},
		{"base", args{[3]uint8{1, 1, 1}}, [3]float32{0.0003035, 0.0003035, 0.0003035}},
		{"max", args{[3]uint8{255, 255, 255}}, [3]float32{1, 1, 1}},
		{"other", args{[3]uint8{60, 120, 180}}, [3]float32{0.045186, 0.1878207, 0.4564110}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DenormalizeRGB(tt.args.v)
			if diff := deep.Equal(got, tt.want); diff != nil {
				t.Errorf("DenormalizeRGB() = %v", diff)
			}
		})
	}
}

func TestDenormalizeRGBA64(t *testing.T) {
	deep.FloatPrecision = 6
	type args struct {
		v [4]uint16
	}
	tests := []struct {
		name string
		args args
		want [4]float32
	}{
		{"empty", args{[4]uint16{}}, [4]float32{}},
		{"base", args{[4]uint16{1000, 1000, 1000, 1000}}, [4]float32{0.001181, 0.001181, 0.001181, 0.0152590}},
		{"max", args{[4]uint16{65535, 65535, 65535, 65535}}, [4]float32{1, 1, 1, 1}},
		{"other", args{[4]uint16{6000, 12000, 18000, 22000}}, [4]float32{0.0087617, 0.0280836, 0.061314, 0.335698}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DenormalizeRGBA64(tt.args.v)
			if diff := deep.Equal(got, tt.want); diff != nil {
				t.Errorf("DenormalizeRGBA64() = %v", diff)
			}
		})
	}
}

func TestDenormalizeRGB64(t *testing.T) {
	deep.FloatPrecision = 6
	type args struct {
		v [3]uint16
	}
	tests := []struct {
		name string
		args args
		want [3]float32
	}{
		{"empty", args{[3]uint16{}}, [3]float32{}},
		{"base", args{[3]uint16{1000, 1000, 1000}}, [3]float32{0.001181, 0.001181, 0.001181}},
		{"max", args{[3]uint16{65535, 65535, 65535}}, [3]float32{1, 1, 1}},
		{"other", args{[3]uint16{6000, 12000, 18000}}, [3]float32{0.0087617, 0.0280836, 0.061314}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DenormalizeRGB64(tt.args.v)
			if diff := deep.Equal(got, tt.want); diff != nil {
				t.Errorf("DenormalizeRGB64() = %v", diff)
			}
		})
	}
}
