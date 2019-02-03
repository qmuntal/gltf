package gltf

// Each key corresponds to mesh attribute semantic and each value is the index of the accessor containing attribute's data.
type Attribute = map[string]uint32

// PrimitiveMode defines the type of primitives to render. All valid values correspond to WebGL enums.
type PrimitiveMode uint8

const (
	Points         PrimitiveMode = 0
	Lines                        = 1
	Line_Loop                    = 2
	Lin_Strip                    = 3
	Triangles                    = 4
	Triangle_Strip               = 5
	Triangle_Fan                 = 6
)

// A Mesh is a set of primitives to be rendered. A node can contain one mesh. A node's transform places the mesh in the scene.
type Mesh struct {
	Named
	Extensible
	Primitives []Primitive `json:"primitives" validator:"required"`
	Weights    []float32   `json:"weights,omitempty"`
}

// Geometry to be rendered with the given material.
type Primitive struct {
	Extensible
	Attributes Attribute   `json:"attributes"`
	Indices    uint32      `json:"indices,omitempty"` // The index of the accessor that contains the indices.
	Material   uint32      `json:"material,omitempty"`
	Mode       uint32      `json:"mode" validator:"lte=6"`
	Targets    []Attribute `json:"targets,omitempty" validator:"omitempty,dive,keys,oneof=POSITION NORMAL TANGENT,endkeys"` // Only POSITION, NORMAL, and TANGENT supported.
}
