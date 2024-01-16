package unlit

import (
	"encoding/json"

	"github.com/qmuntal/gltf"
)

const (
	// ExtensionName defines the Unlit unique key.
	ExtensionName = "KHR_materials_unlit"
)

// Unmarshal decodes the json data into the correct type.
func Unmarshal(data []byte) (any, error) {
	u := new(Unlit)
	err := json.Unmarshal(data, u)
	return u, err
}

func init() {
	gltf.RegisterExtension(ExtensionName, Unmarshal)
}

// Unlit defines an unlit shading model.
// When present, the extension indicates that a material should be unlit.
// Additional properties on the extension object are allowed, but may lead to undefined behaviour in conforming viewers.
type Unlit map[string]any

// UnmarshalJSON unmarshal the pbr with the correct default values.
func (u *Unlit) UnmarshalJSON(data []byte) error {
	if len(*u) == 0 {
		*u = make(Unlit)
	}
	var raw map[string]json.RawMessage
	err := json.Unmarshal(data, &raw)
	if err == nil {
		for key, value := range raw {
			(*u)[key] = value
		}
	}
	return err
}
