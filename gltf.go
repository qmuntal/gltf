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
	BufferView    uint32        `json:"bufferView,omitempty"` // The index of the bufferView.
	ByteOffset    uint32        `json:"byteOffset,omitempty"` // The offset relative to the start of the bufferView in bytes.
	ComponentType ComponentType `json:"componentType"`        // The datatype of components in the attribute.
	Normalized    bool          `json:"normalized,omitempty"` // Specifies whether integer data values should be normalized.
	Count         uint32        `json:"count"`                // The number of attributes referenced by this accessor.
	Type          string        `json:"type"`                 // Specifies if the attribute is a scalar, vector, or matrix.
	Max           []float32     `json:"max,omitempty"`        // Maximum value of each component in this attribute.
	Min           []float32     `json:"min,omitempty"`        // Minimum value of each component in this attribute.
	Sparce        interface{}   `json:"sparce,omitempty"`     // Sparse storage of attributes that deviate from their initialization value.
}

// An Asset is metadata about the glTF asset.
type Asset struct {
	Extensible
	Copyright  string `json:"copyright,omitempty"`  // A copyright message suitable for display to credit the content creator.
	Generator  string `json:"generator,omitempty"`  // Tool that generated this glTF model. Useful for debugging.
	Version    string `json:"version"`              // The glTF version that this asset targets.
	MinVersion string `json:"minVersion,omitempty"` // The minimum glTF version that this asset targets.
}

// Image data used to create a texture. Image can be referenced by URI or bufferView index.
// mimeType is required in the latter case.
type Image struct {
	Named
	Extensible
	URI        string `json:"uri,omitempty"`        // The uri of the image.
	MimeType   string `json:"mimeType,omitempty"`   // The image's MIME type.
	BufferView uint32 `json:"bufferView,omitempty"` // The index of the bufferView that contains the image. Use this instead of the image's uri property.
}

// Indices of those attributes that deviate from their initialization value.
type Indices struct {
	Extensible
	BufferView    uint32        `json:"bufferView"`           // The index of the bufferView with sparse indices.
	ByteOffset    uint32        `json:"byteOffset,omitempty"` // The offset relative to the start of the bufferView in bytes. Must be aligned.
	ComponentType ComponentType `json:"componentType"`        // The indices data type. Valid values correspond to WebGL enums: 5121 (UNSIGNED_BYTE), 5123 (UNSIGNED_SHORT), 5125 (UNSIGNED_INT)
}

// A Mesh is a set of primitives to be rendered. A node can contain one mesh. A node's transform places the mesh in the scene.
type Mesh struct {
	Named
	Extensible
	Primitives []Primitive `json:"primitives"`        // An array of primitives, each defining geometry to be rendered with a material.
	Weights    []float32   `json:"weights,omitempty"` // Array of weights to be applied to the Morph Targets.
}

// Geometry to be rendered with the given material.
type Primitive struct {
	Extensible
	Attributes map[string]uint32 `json:"attributes"`         // A dictionary object, where each key corresponds to mesh attribute semantic and each value is the index of the accessor containing attribute's data.
	Indices    uint32            `json:"indices,omitempty"`  // The index of the accessor that contains the indices.
	Material   uint32            `json:"material,omitempty"` // The index of the material to apply to this primitive when rendering.
	Mode       uint32            `json:"mode"`               // The type of primitives to render.
	Targets    []interface{}     `json:"targets,omitempty"`  // An array of Morph Targets
}

// Extensible is an object that has the Extension and Extras properties.
type Extensible struct {
	Extensions interface{}            `json:"extensions,omitempty"` // Dictionary object with extension-specific objects.
	Extras     map[string]interface{} `json:"extras,omitempty"`     // Application-specific data.
}

// Named is an object that has a Name property.
type Named struct {
	Name string `json:"name,omitempty"` // The user-defined name of this object.
}
