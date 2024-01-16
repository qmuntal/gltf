package lightspuntual

import (
	"math"
	"reflect"
	"testing"

	"github.com/go-test/deep"
	"github.com/qmuntal/gltf"
)

func TestLight_IntensityOrDefault(t *testing.T) {
	tests := []struct {
		name string
		l    *Light
		want float64
	}{
		{"empty", &Light{}, 1},
		{"other", &Light{Intensity: gltf.Float(0.5)}, 0.5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.IntensityOrDefault(); got != tt.want {
				t.Errorf("Light.IntensityOrDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLight_ColorOrDefault(t *testing.T) {
	tests := []struct {
		name string
		l    *Light
		want [3]float64
	}{
		{"empty", &Light{}, [3]float64{1, 1, 1}},
		{"other", &Light{Color: &[3]float64{0.8, 0.8, 0.8}}, [3]float64{0.8, 0.8, 0.8}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.ColorOrDefault(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Light.ColorOrDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSpot_OuterConeAngleOrDefault(t *testing.T) {
	tests := []struct {
		name string
		s    *Spot
		want float64
	}{
		{"empty", &Spot{}, math.Pi / 4},
		{"other", &Spot{OuterConeAngle: gltf.Float(0.5)}, 0.5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.OuterConeAngleOrDefault(); got != tt.want {
				t.Errorf("Spot.OuterConeAngleOrDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLight_UnmarshalJSON(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		l       *Light
		args    args
		want    *Light
		wantErr bool
	}{
		{"default", new(Light), args{[]byte("{}")}, &Light{
			Color: &[3]float64{1, 1, 1}, Intensity: gltf.Float(1), Range: gltf.Float(math.Inf(0)),
		}, false},
		{"nodefault", new(Light), args{[]byte(`{
			"color": [0.3, 0.7, 1.0],
			"name": "AAA",
			"intensity": 40.0,
			"type": "spot",
			"range": 10.0,
			"spot": {
			  "innerConeAngle": 1.0,
			  "outerConeAngle": 2.0
			}
		  }`)}, &Light{
			Name: "AAA", Type: "spot", Color: &[3]float64{0.3, 0.7, 1}, Intensity: gltf.Float(40), Range: gltf.Float(10),
			Spot: &Spot{
				InnerConeAngle: 1.0,
				OuterConeAngle: gltf.Float(2.0),
			},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.l.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Light.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(tt.l, tt.want) {
				t.Errorf("Light.UnmarshalJSON() = %+v, want %+v", tt.l, tt.want)
			}
		})
	}
}

func TestUnmarshal(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{"error", args{[]byte(`{"light: 1}`)}, nil, true},
		{"index", args{[]byte(`{"light": 1}`)}, LightIndex(1), false},
		{"lights", args{[]byte(`{"lights": [
			{
			  "color": [1.0, 0.9, 0.7],
			  "name": "Directional",
			  "intensity": 3.0,
			  "type": "directional"
			},
			{
			  "color": [1.0, 0.0, 0.0],
			  "name": "Point",
			  "intensity": 20.0,
			  "type": "point"
			}
		  ]}`)}, Lights{
			{Color: &[3]float64{1, 0.9, 0.7}, Name: "Directional", Intensity: gltf.Float(3.0), Type: "directional", Range: gltf.Float(math.Inf(0))},
			{Color: &[3]float64{1, 0, 0}, Name: "Point", Intensity: gltf.Float(20.0), Type: "point", Range: gltf.Float(math.Inf(0))},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Unmarshal(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := deep.Equal(got, tt.want); diff != nil {
				t.Errorf("Unmarshal() = %v", diff)
			}
		})
	}
}
