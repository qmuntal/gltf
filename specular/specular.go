package specular

import (
	"bytes"
	"encoding/json"

	"github.com/qmuntal/gltf"
)

const (
	// ExtensionName defines the PBRSpecularGlossiness unique key.
	ExtensionName = "KHR_materials_pbrSpecularGlossiness"
)

// New returns a new specular.PBRSpecularGlossiness.
func New() json.Unmarshaler {
	return new(PBRSpecularGlossiness)
}

func init() {
	gltf.RegisterExtension(ExtensionName, New)
}

// PBRSpecularGlossiness defines a specular-glossiness material model.
type PBRSpecularGlossiness struct {
	DiffuseFactor             *gltf.RGBA        `json:"diffuseFactor,omitempty"`
	DiffuseTexture            *gltf.TextureInfo `json:"diffuseTexture,omitempty"`
	SpecularFactor            *gltf.RGB         `json:"specularFactor,omitempty"`
	GlossinessFactor          *float64          `json:"glossinessFactor,omitempty" validate:"omitempty,gte=0,lte=1"`
	SpecularGlossinessTexture *gltf.TextureInfo `json:"specularGlossinessTexture,omitempty"`
}

// UnmarshalJSON unmarshal the pbr with the correct default values.
func (p *PBRSpecularGlossiness) UnmarshalJSON(data []byte) error {
	type alias PBRSpecularGlossiness
	tmp := alias(PBRSpecularGlossiness{DiffuseFactor: gltf.NewRGBA(), SpecularFactor: gltf.NewRGB(), GlossinessFactor: gltf.Float64(1)})
	err := json.Unmarshal(data, &tmp)
	if err == nil {
		*p = PBRSpecularGlossiness(tmp)
	}
	return err
}

// MarshalJSON marshal the pbr with the correct default values.
func (p *PBRSpecularGlossiness) MarshalJSON() ([]byte, error) {
	type alias PBRSpecularGlossiness
	out, err := json.Marshal(&struct{ *alias }{alias: (*alias)(p)})
	if err == nil {
		if p.GlossinessFactor != nil && *p.GlossinessFactor == 1 {
			out = removeProperty([]byte(`"glossinessFactor":1`), out)
		}
		if p.DiffuseFactor != nil && *p.DiffuseFactor == *gltf.NewRGBA() {
			out = removeProperty([]byte(`"diffuseFactor":[1,1,1,1]`), out)
		}
		if p.SpecularFactor != nil && *p.SpecularFactor == *gltf.NewRGB() {
			out = removeProperty([]byte(`"specularFactor":[1,1,1]`), out)
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
