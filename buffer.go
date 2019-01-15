package gltf

// A buffer points to binary geometry, animation, or skins.
type Buffer struct {
	Named
	Extensible
	URI        string `json:"uri,omitempty"` // The uri of the buffer.
	ByteLength uint32 `json:"byteLength"`    // The total byte length of the buffer view.
}

// BufferView is a view into a buffer generally representing a subset of the buffer.
type BufferView struct {
	Extensible
	Buffer     uint32 `json:"buffer"`               // The index of the buffer.
	ByteOffset uint32 `json:"byteOffset,omitempty"` // The offset into the buffer in bytes.
	ByteLength uint32 `json:"byteLength"`           // The length of the bufferView in bytes.
	ByteStride uint32 `json:"byteStride,omitempty"` // The stride, in bytes.
	Target     uint32 `json:"target,omitempty"`     // The target that the GPU buffer should be bound to.
}
