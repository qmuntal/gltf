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

// Unmarshal decodes the json data into the correct type.
func Unmarshal(data []byte) (any, error) {
	pbr := new(PBRSpecularGlossiness)
	err := json.Unmarshal(data, pbr)
	return pbr, err
}

func init() {
	gltf.RegisterExtension(ExtensionName, Unmarshal)
}

// PBRSpecularGlossiness defines a specular-glossiness material model.
type PBRSpecularGlossiness struct {
	DiffuseFactor             *[4]float64       `json:"diffuseFactor,omitempty"`
	DiffuseTexture            *gltf.TextureInfo `json:"diffuseTexture,omitempty"`
	SpecularFactor            *[3]float64       `json:"specularFactor,omitempty"`
	GlossinessFactor          *float64          `json:"glossinessFactor,omitempty"`
	SpecularGlossinessTexture *gltf.TextureInfo `json:"specularGlossinessTexture,omitempty"`
}

// UnmarshalJSON unmarshal the pbr with the correct default values.
func (p *PBRSpecularGlossiness) UnmarshalJSON(data []byte) error {
	type alias PBRSpecularGlossiness
	tmp := alias(PBRSpecularGlossiness{DiffuseFactor: &[4]float64{1, 1, 1, 1}, SpecularFactor: &[3]float64{1, 1, 1}, GlossinessFactor: gltf.Float(1)})
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
		if p.DiffuseFactor != nil && *p.DiffuseFactor == [4]float64{1, 1, 1, 1} {
			out = removeProperty([]byte(`"diffuseFactor":[1,1,1,1]`), out)
		}
		if p.SpecularFactor != nil && *p.SpecularFactor == [3]float64{1, 1, 1} {
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
