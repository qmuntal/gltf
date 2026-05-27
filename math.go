package gltf

import "math"

// NormalizeByte normalize a float32 into a int8
func NormalizeByte(v float32) int8 {
	return int8(math.Round(float64(v) * 127))
}

// NormalizeByte denormalize a int8 into a float32
func DenormalizeByte(v int8) float32 {
	return float32(math.Max(float64(v)/127, -1))
}

// NormalizeUbyte normalize a float32 into a uint8
func NormalizeUbyte(v float32) uint8 {
	return uint8(math.Round(float64(v) * 255))
}

// DenormalizeUbyte denormalize a uint8 into a float32
func DenormalizeUbyte(v uint8) float32 {
	return float32(v) / 255
}

// NormalizeShort normalize a float32 into a int16
func NormalizeShort(v float32) int16 {
	return int16(math.Round(float64(v) * 32767))
}

// DenormalizeShort denormalize a int16 into a float32
func DenormalizeShort(v int16) float32 {
	return float32(math.Max(float64(v)/32767, -1))
}

// NormalizeUshort normalize a float32 into a uint16
func NormalizeUshort(v float32) uint16 {
	return uint16(math.Round(float64(v) * 65535))
}

// DenormalizeuShort denormalize a uint16 into a float32
func DenormalizeUshort(v uint16) float32 {
	return float32(v) / 65535
}

// NormalizeRGB transform a RGB float32 color (from 0 to 1)
// to its uint8 represtation (from 0 to 255).
func NormalizeRGB(v [3]float32) [3]uint8 {
	return [3]uint8{
		NormalizeUbyte(v[0]),
		NormalizeUbyte(v[1]),
		NormalizeUbyte(v[2]),
	}
}

// DenormalizeRGB transform a RGB uint8 color (from 0 to 255)
// to its float represtation (from 0 to 1).
func DenormalizeRGB(v [3]uint8) [3]float32 {
	return [3]float32{
		DenormalizeUbyte(v[0]),
		DenormalizeUbyte(v[1]),
		DenormalizeUbyte(v[2]),
	}
}

// NormalizeRGB transform a RGBA float32 color (from 0 to 1)
// to its uint8 represtation (from 0 to 255).
func NormalizeRGBA(v [4]float32) [4]uint8 {
	return [4]uint8{
		NormalizeUbyte(v[0]),
		NormalizeUbyte(v[1]),
		NormalizeUbyte(v[2]),
		NormalizeUbyte(v[3]),
	}
}

// DenormalizeRGBA transform a RGBA uint8 color (from 0 to 255)
// to its float represtation (from 0 to 1).
func DenormalizeRGBA(v [4]uint8) [4]float32 {
	return [4]float32{
		DenormalizeUbyte(v[0]),
		DenormalizeUbyte(v[1]),
		DenormalizeUbyte(v[2]),
		DenormalizeUbyte(v[3]),
	}
}

// NormalizeRGB64 transform a RGB float32 color (from 0 to 1)
// to its uint16 represtation (from 0 to 65535).
func NormalizeRGB64(v [3]float32) [3]uint16 {
	return [3]uint16{
		NormalizeUshort(v[0]),
		NormalizeUshort(v[1]),
		NormalizeUshort(v[2]),
	}
}

// DenormalizeRGB64 transform a RGB uint16 color (from 0 to 65535)
// to its float represtation (from 0 to 1).
func DenormalizeRGB64(v [3]uint16) [3]float32 {
	return [3]float32{
		DenormalizeUshort(v[0]),
		DenormalizeUshort(v[1]),
		DenormalizeUshort(v[2]),
	}
}

// NormalizeRGBA64 transform a RGBA float32 color (from 0 to 1)
// to its uint16 represtation (from 0 to 65535).
func NormalizeRGBA64(v [4]float32) [4]uint16 {
	return [4]uint16{
		NormalizeUshort(v[0]),
		NormalizeUshort(v[1]),
		NormalizeUshort(v[2]),
		NormalizeUshort(v[3]),
	}
}

// DenormalizeRGBA64 transform a RGBA uint16 color (from 0 to 65535)
// to its float represtation (from 0 to 1).
func DenormalizeRGBA64(v [4]uint16) [4]float32 {
	return [4]float32{
		DenormalizeUshort(v[0]),
		DenormalizeUshort(v[1]),
		DenormalizeUshort(v[2]),
		DenormalizeUshort(v[3]),
	}
}
