package gltf

import "encoding/json"

const (
	ExtPBRSpecularGlossiness = "KHR_materials_pbrSpecularGlossiness"
)

// Extension is map where the keys are the extension identifiers and the values are the extensions payloads.
// If a key matches with one of the supported extensions the value will be marshalled as a pointer to the extension struct.
// If a key does not match with any of the supported extensions the value will be a json.RawMessage so its decoding can be delayed.
type Extensions map[string]interface{}

type envelope map[string]json.RawMessage

// UnmarshalJSON unmarshal the extensions with the supported extensions initialized.
func (ext *Extensions) UnmarshalJSON(data []byte) error {
	if len(*ext) == 0 {
		*ext = make(Extensions)
	}
	var raw envelope
	err := json.Unmarshal(data, &raw)
	if err == nil {
		for key, value := range raw {
			switch key {
			case ExtPBRSpecularGlossiness:
				n := &PBRSpecularGlossiness{}
				if err := json.Unmarshal(value, n); err != nil {
					return err
				}
				(*ext)[ExtPBRSpecularGlossiness] = n
			default:
				(*ext)[key] = value
			}
		}
	}

	return err
}

type PBRSpecularGlossiness struct {
	DiffuseFactor             [4]float64   `json:"diffuseFactor" validate:"dive,gte=0,lte=1"`
	DiffuseTexture            *TextureInfo `json:"diffuseTexture,omitempty"`
	SpecularFactor            [3]float64   `json:"specularFactor" validate:"dive,gte=0,lte=1"`
	GlossinessFactor          float64      `json:"glossinessFactor" validate:"gte=0,lte=1"`
	SpecularGlossinessTexture *TextureInfo `json:"specularGlossinessTexture,omitempty"`
}

// PBRSpecularGlossiness returns a default PBRSpecularGlossiness.
func NewPBRSpecularGlossiness() *PBRSpecularGlossiness {
	return &PBRSpecularGlossiness{DiffuseFactor: [4]float64{1, 1, 1, 1}, SpecularFactor: [3]float64{1, 1, 1}, GlossinessFactor: 1}
}

// UnmarshalJSON unmarshal the pbr with the correct default values.
func (p *PBRSpecularGlossiness) UnmarshalJSON(data []byte) error {
	type alias PBRSpecularGlossiness
	tmp := alias(*NewPBRSpecularGlossiness())
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
		if p.GlossinessFactor == 1 {
			out = removeProperty([]byte(`"glossinessFactor":1`), out)
		}
		if p.DiffuseFactor == [4]float64{1, 1, 1, 1} {
			out = removeProperty([]byte(`"diffuseFactor":[1,1,1,1]`), out)
		}
		if p.SpecularFactor == [3]float64{1, 1, 1} {
			out = removeProperty([]byte(`"specularFactor":[1,1,1]`), out)
		}
		out = sanitizeJSON(out)
	}
	return out, err
}
