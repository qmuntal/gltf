// Package lightspuntual defines three "punctual" light types: directional, point and spot.
// Punctual lights are defined as parameterized, infinitely small points
// that emit light in well-defined directions and intensities.
package lightspuntual

import (
	"encoding/json"
	"math"

	"github.com/qmuntal/gltf"
)

const (
	// ExtensionName defines the KHR_lights_punctual unique key.
	ExtensionName = "KHR_lights_punctual"
)

func init() {
	gltf.RegisterExtension(ExtensionName, Unmarshal)
}

type envelop struct {
	Lights Lights
	Light  *LightIndex
}

// Unmarshal decodes the json data into the correct type.
func Unmarshal(data []byte) (interface{}, error) {
	var env envelop
	if err := json.Unmarshal([]byte(data), &env); err != nil {
		return nil, err
	}
	if env.Light != nil {
		return *env.Light, nil
	}
	return env.Lights, nil
}

const (
	// TypeDirectional lights act as though they are infinitely far away and emit light in the direction of the local -z axis.
	TypeDirectional = "directional"
	// TypePoint lights emit light in all directions from their position in space.
	TypePoint = "point"
	// TypeSpot lights emit light in a cone in the direction of the local -z axis.
	TypeSpot = "spot"
)

// LightIndex is the id of the light referenced by this node.
type LightIndex uint32

// Spot defines the spot cone.
type Spot struct {
	InnerConeAngle float64  `json:"innerConeAngle,omitempty"`
	OuterConeAngle *float64 `json:"outerConeAngle,omitempty"`
}

// OuterConeAngleOrDefault returns the OuterConeAngle if it is not nil, else return the default one.
func (s *Spot) OuterConeAngleOrDefault() float64 {
	if s.OuterConeAngle == nil {
		return math.Pi / 4
	}
	return *s.OuterConeAngle
}

// UnmarshalJSON unmarshal the spot with the correct default values.
func (s *Spot) UnmarshalJSON(data []byte) error {
	type alias Spot
	tmp := alias(Spot{OuterConeAngle: gltf.Float(math.Pi / 4)})
	err := json.Unmarshal(data, &tmp)
	if err == nil {
		*s = Spot(tmp)
	}
	return err
}

// Lights defines a list of Lights.
type Lights []*Light

// Light defines a directional, point, or spot light.
// When a light's type is spot, the spot property on the light is required.
type Light struct {
	Type      string      `json:"type"`
	Name      string      `json:"name,omitempty"`
	Color     *[3]float64 `json:"color,omitempty" validate:"omitempty,dive,gte=0,lte=1"`
	Intensity *float64    `json:"intensity,omitempty"`
	Range     *float64    `json:"range,omitempty"`
	Spot      *Spot       `json:"spot,omitempty"`
}

// IntensityOrDefault returns the itensity if it is not nil, else return the default one.
func (l *Light) IntensityOrDefault() float64 {
	if l.Intensity == nil {
		return 1
	}
	return *l.Intensity
}

// ColorOrDefault returns the color if it is not nil, else return the default one.
func (l *Light) ColorOrDefault() [3]float64 {
	if l.Color == nil {
		return [3]float64{1, 1, 1}
	}
	return *l.Color
}

// UnmarshalJSON unmarshal the light with the correct default values.
func (l *Light) UnmarshalJSON(data []byte) error {
	type alias Light
	tmp := alias(Light{Color: &[3]float64{1, 1, 1}, Intensity: gltf.Float(1), Range: gltf.Float(math.Inf(0))})
	err := json.Unmarshal(data, &tmp)
	if err == nil {
		*l = Light(tmp)
	}
	return err
}
