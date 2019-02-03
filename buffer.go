package gltf

// The target that the GPU buffer should be bound to.
type Target uint16

const (
	ArrayBuffer        Target = 34962
	ElementArrayBuffer        = 34963
)

// A buffer points to binary geometry, animation, or skins.
type Buffer struct {
	Named
	Extensible
	URI        string `json:"uri,omitempty" validator:"omitempty|uri|datauri"`
	ByteLength uint32 `json:"byteLength" validator:"required"`
}

// BufferView is a view into a buffer generally representing a subset of the buffer.
type BufferView struct {
	Extensible
	Buffer     uint32 `json:"buffer"`
	ByteOffset uint32 `json:"byteOffset,omitempty"`
	ByteLength uint32 `json:"byteLength" validator:"required"`
	ByteStride uint32 `json:"byteStride,omitempty" validator:"omitempty,gte=4,lte=252"`
	Target     Target `json:"target,omitempty" validator:"omitempty,oneof=34962 34963"`
}
