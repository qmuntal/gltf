package unlit

import (
	"encoding/json"

	"github.com/qmuntal/gltf"
)

const (
	// ExtUnlit defines the Unlit unique key.
	ExtUnlit = "KHR_materials_unlit"
)

// New returns a new unlit.Unlit.
func New() json.Unmarshaler {
	return new(Unlit)
}

func init() {
	gltf.RegisterExtension(ExtUnlit, New)
}

// Unlit defines an unlit shading model.
// When present, the extension indicates that a material should be unlit.
// Additional properties on the extension object are allowed, but may lead to undefined behaviour in conforming viewers.
type Unlit map[string]iterface{}

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
