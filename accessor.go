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

// AccessorType specifies if the attribute is a scalar, vector, or matrix.
type AccessorType string

const (
	Scalar AccessorType = "SCALAR"
	Vec2                = "VEC2"
	Vec3                = "VEC3"
	Vec4                = "VEC4"
	Mat2                = "MAT2"
	Mat3                = "MAT3"
	Mat4                = "MAT4"
)

// An accessor is a typed view into a bufferView. A bufferView contains raw binary data.
// An accessor provides a typed view into a bufferView or a subset of a bufferView
// similar to how WebGL's vertexAttribPointer() defines an attribute in a buffer.
type Accessor struct {
	Named
	Extensible
	BufferView    uint32        `json:"bufferView,omitempty"`
	ByteOffset    uint32        `json:"byteOffset,omitempty"`
	ComponentType ComponentType `json:"componentType" validator:"oneof=5120 5121 5122 5123 5125 5126"`
	Normalized    bool          `json:"normalized,omitempty"`       // Specifies whether integer data values should be normalized.
	Count         uint32        `json:"count" validator:"required"` // The number of attributes referenced by this accessor.
	Type          AccessorType  `json:"type" validator:"oneof=SCALAR VEC2 VEC3 VEC4 MAT2 MAT3 MAT4"`
	Max           []float32     `json:"max,omitempty" validator:"omitempty,lte=16"` // Maximum value of each component in this attribute.
	Min           []float32     `json:"min,omitempty" validator:"omitempty,lte=16"` // Minimum value of each component in this attribute.
	Sparce        Sparse        `json:"sparce,omitempty"`                           // Sparse storage of attributes that deviate from their initialization value.
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
	ComponentType ComponentType `json:"componentType" validator:"oneof=5121 5123 5125"`
}
