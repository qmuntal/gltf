package texturetransform

import (
	"bytes"
	"encoding/json"

	"github.com/qmuntal/gltf"
)

var (
	// DefaultScale defines a scaling that does not modify the size of the object.
	DefaultScale = [2]float64{1, 1}
	emptyScale   = [2]float64{0, 0}
	emptyOffset  = [2]float64{0, 0}
)

const (
	// ExtensionName defines the ExtTextureTransform unique key.
	ExtensionName = "KHR_texture_transform"
)

// Unmarshal decodes the json data into the correct type.
func Unmarshal(data []byte) (any, error) {
	t := new(TextureTranform)
	err := json.Unmarshal(data, t)
	return t, err
}

func init() {
	gltf.RegisterExtension(ExtensionName, Unmarshal)
}

// TextureTranform can be used in textureInfo to pack many low-res texture into a single large texture atlas.
type TextureTranform struct {
	Offset   [2]float64 `json:"offset"`
	Rotation float64    `json:"rotation,omitempty"`
	Scale    [2]float64 `json:"scale"`
	TexCoord *uint32    `json:"texCoord,omitempty"`
}

// ScaleOrDefault returns the node scale if it represents a valid scale factor, else return the default one.
func (t *TextureTranform) ScaleOrDefault() [2]float64 {
	if t.Scale == emptyScale {
		return DefaultScale
	}
	return t.Scale
}

// UnmarshalJSON unmarshal the pbr with the correct default values.
func (t *TextureTranform) UnmarshalJSON(data []byte) error {
	type alias TextureTranform
	tmp := alias(TextureTranform{Scale: DefaultScale})
	err := json.Unmarshal(data, &tmp)
	if err == nil {
		*t = TextureTranform(tmp)
	}
	return err
}

// MarshalJSON marshal the pbr with the correct default values.
func (t *TextureTranform) MarshalJSON() ([]byte, error) {
	type alias TextureTranform
	out, err := json.Marshal(&struct{ *alias }{alias: (*alias)(t)})
	if err == nil {
		if t.Scale == DefaultScale {
			out = removeProperty([]byte(`"scale":[1,1]`), out)
		} else if t.Scale == emptyScale {
			out = removeProperty([]byte(`"scale":[0,0]`), out)
		}
		if t.Offset == emptyOffset {
			out = removeProperty([]byte(`"offset":[0,0]`), out)
		}
		out = sanitizeJSON(out)
	}
	return out, err
}

func removeProperty(str []byte, b []byte) []byte {
	b = bytes.Replace(b, str, []byte(""), 1)
	return bytes.Replace(b, []byte(`,,`), []byte(","), 1)
}

func sanitizeJSON(b []byte) []byte {
	b = bytes.Replace(b, []byte(`{,`), []byte("{"), 1)
	return bytes.Replace(b, []byte(`,}`), []byte("}"), 1)
}
