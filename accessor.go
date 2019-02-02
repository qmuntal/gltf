package gltf

//The ComponentType is the datatype of components in the attribute. All valid values correspond to WebGL enums.
// The corresponding typed arrays are Int8Array, Uint8Array, Int16Array, Uint16Array, Uint32Array, and Float32Array, respectively.
// 5125 (UNSIGNED_INT) is only allowed when the accessor contains indices.
type ComponentType uint16

const (
	Byte          ComponentType = 5120
	UnsignedByte                = 5121
	Short                       = 5122
	UnsignedShort               = 5123
	UnsignedInt                 = 5125
	Float                       = 5126
)

// An accessor is a typed view into a bufferView. A bufferView contains raw binary data.
// An accessor provides a typed view into a bufferView or a subset of a bufferView
// similar to how WebGL's vertexAttribPointer() defines an attribute in a buffer.
type Accessor struct {
	Named
	Extensible
	BufferView    uint32        `json:"bufferView,omitempty"`
	ByteOffset    uint32        `json:"byteOffset,omitempty"`
	ComponentType ComponentType `json:"componentType"`
	Normalized    bool          `json:"normalized,omitempty"` // Specifies whether integer data values should be normalized.
	Count         uint32        `json:"count"`                // The number of attributes referenced by this accessor.
	Type          string        `json:"type"`                 // Specifies if the attribute is a scalar, vector, or matrix.
	Max           []float32     `json:"max,omitempty"`        // Maximum value of each component in this attribute.
	Min           []float32     `json:"min,omitempty"`        // Minimum value of each component in this attribute.
	Sparce        Sparse        `json:"sparce,omitempty"`     // Sparse storage of attributes that deviate from their initialization value.
}

// Sparse storage of attributes that deviate from their initialization value.
type Sparse struct {
	Extensible
	Count   uint32        `json:"count"`   // Number of entries stored in the sparse array.
	Indices SparceIndices `json:"indices"` // Index array of size count that points to those accessor attributes that deviate from their initialization value.
	Values  SparseValues  `json:"values"`  // Array of size count times number of components, storing the displaced accessor attributes pointed by indices.
}

// SparseValues stores the displaced accessor attributes pointed by accessor.sparse.indices.
type SparseValues struct {
	Extensible
	BufferView uint32 `json:"bufferView"`
	ByteOffset uint32 `json:"byteOffset,omitempty"`
}

// SparceIndices defines the indices of those attributes that deviate from their initialization value.
type SparceIndices struct {
	Extensible
	BufferView    uint32        `json:"bufferView"`
	ByteOffset    uint32        `json:"byteOffset,omitempty"`
	ComponentType ComponentType `json:"componentType"`
}
