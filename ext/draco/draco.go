package draco

import (
	"encoding/json"

	"github.com/qmuntal/gltf"
)

const (
	// ExtensionName defines the KHR_draco_mesh_compression unique key.
	ExtensionName = "KHR_draco_mesh_compression"
)

func init() {
	gltf.RegisterExtension(ExtensionName, Unmarshal)
}

// Unmarshal decodes the json data into the correct type.
func Unmarshal(data []byte) (interface{}, error) {
	drc := new(Primitive)
	err := json.Unmarshal(data, drc)
	return drc, err
}

// Primitive extension points to the bufferView that contains the compressed data.
type Primitive struct {
	Extensions gltf.Extensions `json:"extensions,omitempty"`
	Extras     interface{}     `json:"extras,omitempty"`
	BufferView uint32          `json:"bufferView"`
	Attributes gltf.Attribute  `json:"attributes"`
}
