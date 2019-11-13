package binary

import "unsafe"

func int8bits(b int8) uint8 { return *(*uint8)(unsafe.Pointer(&b)) }

func int16bits(b int16) uint16 { return *(*uint16)(unsafe.Pointer(&b)) }
