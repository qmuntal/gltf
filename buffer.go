package gltf

// A buffer points to binary geometry, animation, or skins.
type Buffer struct {
	Named
	Extensible
	URI        string `json:"uri,omitempty"`
	ByteLength uint32 `json:"byteLength"`
}

// BufferView is a view into a buffer generally representing a subset of the buffer.
type BufferView struct {
	Extensible
	Buffer     uint32 `json:"buffer"`
	ByteOffset uint32 `json:"byteOffset,omitempty"`
	ByteLength uint32 `json:"byteLength"`
	ByteStride uint32 `json:"byteStride,omitempty"`
	Target     uint32 `json:"target,omitempty"`
}
