package gltf

// The Scene contains a list of root nodes.
type Scene struct {
	Named
	Extensible
	Nodes []uint32 `json:"nodes,omitempty"`
}

// A node in the node hierarchy.
// A node can have either a matrix or any combination of translation/rotation/scale (TRS) properties.
type Node struct {
	Named
	Extensible
	Camera      uint32      `json:"camera,omitempty"`
	Children    []uint32    `json:"children,omitempty" validator:"omitempty,unique"`
	Skin        uint32      `json:"skin,omitempty"`
	Matrix      [16]float32 `json:"matrix,omitempty"` // A 4x4 transformation matrix stored in column-major order.
	Mesh        uint32      `json:"mesh,omitempty"`
	Rotation    [4]float64  `json:"rotation" validator:"omitempty,dive,gte=-1,lte=1"` // The node's unit quaternion rotation in the order (x, y, z, w), where w is the scalar.
	Scale       [3]float32  `json:"scale,omitempty"`
	Translation [3]float32  `json:"translation,omitempty"`
	Weights     []float32   `json:"weights,omitempty"` // The weights of the instantiated Morph Target.
}

// Skin defines joints and matrices.
type Skin struct {
	Named
	Extensible
	InverseBindMatrices uint32   `json:"inverseBindMatrices,omitempty"` // The index of the accessor containing the floating-point 4x4 inverse-bind matrices.
	Skeleton            uint32   `json:"skeleton,omitempty"`            // The index of the node used as a skeleton root. When undefined, joints transforms resolve to scene root.
	Joints              []uint32 `json:"joints"`                        // Indices of skeleton nodes, used as joints in this skin.
}
