package gltf

import (
	"bytes"
	"encoding/json"
	"errors"
)

// UnmarshalJSON unmarshal the node with the correct default values.
func (n *Node) UnmarshalJSON(data []byte) error {
	type alias Node
	tmp := alias(Node{
		Matrix:   DefaultMatrix,
		Rotation: DefaultRotation,
		Scale:    DefaultScale,
	})
	err := json.Unmarshal(data, &tmp)
	if err == nil {
		*n = Node(tmp)
	}
	return err
}

// MarshalJSON marshal the node with the correct default values.
func (n *Node) MarshalJSON() ([]byte, error) {
	type alias Node
	out, err := json.Marshal(&struct{ *alias }{alias: (*alias)(n)})
	if err == nil {
		if n.Matrix == DefaultMatrix {
			out = removeProperty([]byte(`"matrix":[1,0,0,0,0,1,0,0,0,0,1,0,0,0,0,1]`), out)
		} else if n.Matrix == emptyMatrix {
			out = removeProperty([]byte(`"matrix":[0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0]`), out)
		}

		if n.Rotation == DefaultRotation {
			out = removeProperty([]byte(`"rotation":[0,0,0,1]`), out)
		} else if n.Rotation == emptyRotation {
			out = removeProperty([]byte(`"rotation":[0,0,0,0]`), out)
		}

		if n.Scale == DefaultScale {
			out = removeProperty([]byte(`"scale":[1,1,1]`), out)
		} else if n.Scale == emptyScale {
			out = removeProperty([]byte(`"scale":[0,0,0]`), out)
		}

		if n.Translation == DefaultTranslation {
			out = removeProperty([]byte(`"translation":[0,0,0]`), out)
		}
		out = sanitizeJSON(out)
	}
	return out, err
}

// MarshalJSON marshal the camera with the correct default values.
func (c *Camera) MarshalJSON() ([]byte, error) {
	type alias Camera
	if c.Perspective != nil {
		return json.Marshal(&struct {
			Type string `json:"type"`
			*alias
		}{
			Type:  "perspective",
			alias: (*alias)(c),
		})
	} else if c.Orthographic != nil {
		return json.Marshal(&struct {
			Type string `json:"type"`
			*alias
		}{
			Type:  "orthographic",
			alias: (*alias)(c),
		})
	}
	return nil, errors.New("gltf: camera must defined either the perspective or orthographic property")
}

// UnmarshalJSON unmarshal the material with the correct default values.
func (m *Material) UnmarshalJSON(data []byte) error {
	type alias Material
	tmp := alias(Material{AlphaCutoff: Float64(0.5)})
	err := json.Unmarshal(data, &tmp)
	if err == nil {
		*m = Material(tmp)
	}
	return err
}

// MarshalJSON marshal the material with the correct default values.
func (m *Material) MarshalJSON() ([]byte, error) {
	type alias Material
	out, err := json.Marshal(&struct{ *alias }{alias: (*alias)(m)})
	if err == nil {
		if m.AlphaCutoff != nil && *m.AlphaCutoff == 0.5 {
			out = removeProperty([]byte(`"alphaCutoff":0.5`), out)
		}
		if m.EmissiveFactor == [3]float64{0, 0, 0} {
			out = removeProperty([]byte(`"emissiveFactor":[0,0,0]`), out)
		}
		out = sanitizeJSON(out)
	}
	return out, err
}

// UnmarshalJSON unmarshal the texture info with the correct default values.
func (n *NormalTexture) UnmarshalJSON(data []byte) error {
	type alias NormalTexture
	tmp := alias(NormalTexture{Scale: Float64(1)})
	err := json.Unmarshal(data, &tmp)
	if err == nil {
		*n = NormalTexture(tmp)
	}
	return err
}

// MarshalJSON marshal the texture info with the correct default values.
func (n *NormalTexture) MarshalJSON() ([]byte, error) {
	type alias NormalTexture
	if n.Scale != nil && *n.Scale == 1 {
		return json.Marshal(&struct {
			Scale float64 `json:"scale,omitempty"`
			*alias
		}{
			Scale: 0,
			alias: (*alias)(n),
		})
	}
	return json.Marshal(&struct{ *alias }{alias: (*alias)(n)})
}

// UnmarshalJSON unmarshal the texture info with the correct default values.
func (o *OcclusionTexture) UnmarshalJSON(data []byte) error {
	type alias OcclusionTexture
	tmp := alias(OcclusionTexture{Strength: Float64(1)})
	err := json.Unmarshal(data, &tmp)
	if err == nil {
		*o = OcclusionTexture(tmp)
	}
	return err
}

// MarshalJSON marshal the texture info with the correct default values.
func (o *OcclusionTexture) MarshalJSON() ([]byte, error) {
	type alias OcclusionTexture
	if o.Strength != nil && *o.Strength == 1 {
		return json.Marshal(&struct {
			Strength float64 `json:"strength,omitempty"`
			*alias
		}{
			Strength: 0,
			alias:    (*alias)(o),
		})
	}
	return json.Marshal(&struct{ *alias }{alias: (*alias)(o)})
}

// UnmarshalJSON unmarshal the color with the correct default values.
func (c *RGBA) UnmarshalJSON(data []byte) error {
	tmp := [4]float64{1, 1, 1, 1}
	err := json.Unmarshal(data, &tmp)
	if err == nil {
		c.R, c.G, c.B, c.A = tmp[0], tmp[1], tmp[2], tmp[3]
	}
	return err
}

// MarshalJSON marshal the color with the correct default values.
func (c *RGBA) MarshalJSON() ([]byte, error) {
	return json.Marshal([4]float64{c.R, c.G, c.B, c.A})
}

// UnmarshalJSON unmarshal the color with the correct default values.
func (c *RGB) UnmarshalJSON(data []byte) error {
	tmp := [3]float64{1, 1, 1}
	err := json.Unmarshal(data, &tmp)
	if err == nil {
		c.R, c.G, c.B = tmp[0], tmp[1], tmp[2]
	}
	return err
}

// MarshalJSON marshal the color with the correct default values.
func (c *RGB) MarshalJSON() ([]byte, error) {
	return json.Marshal([3]float64{c.R, c.G, c.B})
}

// UnmarshalJSON unmarshal the pbr with the correct default values.
func (p *PBRMetallicRoughness) UnmarshalJSON(data []byte) error {
	type alias PBRMetallicRoughness
	tmp := alias(PBRMetallicRoughness{BaseColorFactor: NewRGBA(), MetallicFactor: Float64(1), RoughnessFactor: Float64(1)})
	err := json.Unmarshal(data, &tmp)
	if err == nil {
		*p = PBRMetallicRoughness(tmp)
	}
	return err
}

// MarshalJSON marshal the pbr with the correct default values.
func (p *PBRMetallicRoughness) MarshalJSON() ([]byte, error) {
	type alias PBRMetallicRoughness
	out, err := json.Marshal(&struct{ *alias }{alias: (*alias)(p)})
	if err == nil {
		if p.MetallicFactor != nil && *p.MetallicFactor == 1 {
			out = removeProperty([]byte(`"metallicFactor":1`), out)
		}
		if p.RoughnessFactor != nil && *p.RoughnessFactor == 1 {
			out = removeProperty([]byte(`"roughnessFactor":1`), out)
		}
		if p.BaseColorFactor != nil && *p.BaseColorFactor == *NewRGBA() {
			out = removeProperty([]byte(`"baseColorFactor":[1,1,1,1]`), out)
		}
		out = sanitizeJSON(out)
	}
	return out, err
}

// UnmarshalJSON unmarshal the extensions with the supported extensions initialized.
func (ext *Extensions) UnmarshalJSON(data []byte) error {
	if len(*ext) == 0 {
		*ext = make(Extensions)
	}
	var raw map[string]json.RawMessage
	err := json.Unmarshal(data, &raw)
	if err == nil {
		for key, value := range raw {
			if extFactory, ok := extensions[key]; ok {
				n, err := extFactory(value)
				if err != nil {
					(*ext)[key] = value
				} else {
					(*ext)[key] = n
				}
			} else {
				(*ext)[key] = value
			}
		}
	}

	return err
}

func removeProperty(str []byte, b []byte) []byte {
	b = bytes.Replace(b, str, []byte(""), 1)
	return bytes.Replace(b, []byte(`,,`), []byte(","), 1)
}

func sanitizeJSON(b []byte) []byte {
	b = bytes.Replace(b, []byte(`{,`), []byte("{"), 1)
	return bytes.Replace(b, []byte(`,}`), []byte("}"), 1)
}
