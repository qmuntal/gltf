package modeler

import "math"

func compressUint32(data []uint32) interface{} {
	u8, u16 := true, true
	// Start at the end, as there are more possibilities to have a higher numbers.
	for i := len(data) - 1; i >= 0; i-- {
		x := data[i]
		if u16 && x >= math.MaxUint16 {
			u8, u16 = false, false
			break
		} else if u8 && x >= math.MaxUint8 {
			u8 = false
		}
	}
	if u8 {
		optData := make([]uint8, len(data))
		for i, x := range data {
			optData[i] = uint8(x)
		}
		return optData
	}
	if u16 {
		optData := make([]uint16, len(data))
		for i, x := range data {
			optData[i] = uint16(x)
		}
		return optData
	}
	return data
}

func compressUint16(data []uint16) interface{} {
	u8 := true
	// Start at the end, as there are more possibilities to have a higher numbers.
	for i := len(data) - 1; i >= 0; i-- {
		x := data[i]
		if u8 && x >= math.MaxUint8 {
			u8 = false
			break
		}
	}
	if u8 {
		optData := make([]uint8, len(data))
		for i, x := range data {
			optData[i] = uint8(x)
		}
		return optData
	}
	return data
}
