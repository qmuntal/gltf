package gltf

type Attribute = map[string]uint32

// A Mesh is a set of primitives to be rendered. A node can contain one mesh. A node's transform places the mesh in the scene.
type Mesh struct {
	Named
	Extensible
	Primitives []Primitive `json:"primitives"`
	Weights    []float32   `json:"weights,omitempty"`
}

// Geometry to be rendered with the given material.
type Primitive struct {
	Extensible
	Attributes Attribute   `json:"attributes"`        // Each key corresponds to mesh attribute semantic and each value is the index of the accessor containing attribute's data.
	Indices    uint32      `json:"indices,omitempty"` // The index of the accessor that contains the indices.
	Material   uint32      `json:"material,omitempty"`
	Mode       uint32      `json:"mode"`
	Targets    []Attribute `json:"targets,omitempty"` // Only POSITION, NORMAL, and TANGENT supported.
}
