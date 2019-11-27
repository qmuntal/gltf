package gltf

import "encoding/json"

var (
	// DefaultMatrix defines an identity matrix.
	DefaultMatrix = [16]float64{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1}
	// DefaultRotation defines a quaternion without rotation.
	DefaultRotation = [4]float64{0, 0, 0, 1}
	// DefaultScale defines a scaling that does not modify the size of the object.
	DefaultScale = [3]float64{1, 1, 1}
	// DefaultTranslation defines a translation that does not move the object.
	DefaultTranslation = [3]float64{0, 0, 0}
)

var (
	emptyMatrix   = [16]float64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	emptyRotation = [4]float64{0, 0, 0, 0}
	emptyScale    = [3]float64{0, 0, 0}
)

// The ComponentType is the datatype of components in the attribute. All valid values correspond to WebGL enums.
// 5125 (UNSIGNED_INT) is only allowed when the accessor contains indices.
type ComponentType uint16

const (
	// ComponentFloat corresponds to a Float32Array.
	ComponentFloat ComponentType = iota
	// ComponentByte corresponds to a Int8Array.
	ComponentByte
	// ComponentUbyte corresponds to a Uint8Array.
	ComponentUbyte
	// ComponentShort corresponds to a Int16Array.
	ComponentShort
	// ComponentUshort corresponds to a Uint16Array.
	ComponentUshort
	// ComponentUint corresponds to a Uint32Array.
	ComponentUint
)

// UnmarshalJSON unmarshal the component type with the correct default values.
func (c *ComponentType) UnmarshalJSON(data []byte) error {
	var tmp uint16
	err := json.Unmarshal(data, &tmp)
	if err == nil {
		*c = map[uint16]ComponentType{
			5120: ComponentByte,
			5121: ComponentUbyte,
			5122: ComponentShort,
			5123: ComponentUshort,
			5125: ComponentUint,
			5126: ComponentFloat,
		}[tmp]
	}
	return err
}

// MarshalJSON marshal the component type with the correct default values.
func (c *ComponentType) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[ComponentType]uint16{
		ComponentByte:          5120,
		ComponentUbyte:  5121,
		ComponentShort:         5122,
		ComponentUshort: 5123,
		ComponentUint:   5125,
		ComponentFloat:         5126,
	}[*c])
}

// AccessorType specifies if the attribute is a scalar, vector, or matrix.
type AccessorType uint8

const (
	// Scalar corresponds to a single dimension value.
	Scalar AccessorType = iota
	// Vec2 corresponds to a two dimensions array.
	Vec2
	// Vec3 corresponds to a three dimensions array.
	Vec3
	// Vec4 corresponds to a four dimensions array.
	Vec4
	// Mat2 corresponds to a 2x2 matrix.
	Mat2
	// Mat3 corresponds to a 3x3 matrix.
	Mat3
	// Mat4 corresponds to a 4x4 matrix.
	Mat4
)

// UnmarshalJSON unmarshal the accessor type with the correct default values.
func (a *AccessorType) UnmarshalJSON(data []byte) error {
	var tmp string
	err := json.Unmarshal(data, &tmp)
	if err == nil {
		*a = map[string]AccessorType{
			"SCALAR": Scalar,
			"VEC2":   Vec2,
			"VEC3":   Vec3,
			"VEC4":   Vec4,
			"MAT2":   Mat2,
			"MAT3":   Mat3,
			"MAT4":   Mat4,
		}[tmp]
	}
	return err
}

// MarshalJSON marshal the accessor type with the correct default values.
func (a *AccessorType) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[AccessorType]string{
		Scalar: "SCALAR",
		Vec2:   "VEC2",
		Vec3:   "VEC3",
		Vec4:   "VEC4",
		Mat2:   "MAT2",
		Mat3:   "MAT3",
		Mat4:   "MAT4",
	}[*a])
}

// The Target that the GPU buffer should be bound to.
type Target uint16

const (
	// None is used when the buffer should not bound to a target, for example when referenced by an sparce indices.
	None = 0
	// ArrayBuffer corresponds to an array buffer.
	ArrayBuffer Target = 34962
	// ElementArrayBuffer corresponds to an element array buffer.
	ElementArrayBuffer = 34963
)

// Attribute is a map that each key corresponds to mesh attribute semantic and each value is the index of the accessor containing attribute's data.
type Attribute = map[string]uint32

// PrimitiveMode defines the type of primitives to render. All valid values correspond to WebGL enums.
type PrimitiveMode uint8

const (
	// Triangles corresponds to a Triangle primitive.
	Triangles PrimitiveMode = iota
	// Points corresponds to a Point primitive.
	Points
	// Lines corresponds to a Line primitive.
	Lines
	// LineLoop corresponds to a Line Loop primitive.
	LineLoop
	// LineStrip corresponds to a Line Strip primitive.
	LineStrip
	// TriangleStrip corresponds to a Triangle Strip primitive.
	TriangleStrip
	// TriangleFan corresponds to a Triangle Fan primitive.
	TriangleFan
)

// UnmarshalJSON unmarshal the primitive mode with the correct default values.
func (p *PrimitiveMode) UnmarshalJSON(data []byte) error {
	var tmp uint8
	err := json.Unmarshal(data, &tmp)
	if err == nil {
		*p = map[uint8]PrimitiveMode{
			0: Points,
			1: Lines,
			2: LineLoop,
			3: LineStrip,
			4: Triangles,
			5: TriangleStrip,
			6: TriangleFan,
		}[tmp]
	}
	return err
}

// MarshalJSON marshal the primitive mode with the correct default values.
func (p *PrimitiveMode) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[PrimitiveMode]uint8{
		Points:        0,
		Lines:         1,
		LineLoop:      2,
		LineStrip:     3,
		Triangles:     4,
		TriangleStrip: 5,
		TriangleFan:   6,
	}[*p])
}

// The AlphaMode enumeration specifying the interpretation of the alpha value of the main factor and texture.
type AlphaMode uint8

const (
	// Opaque corresponds to an Opaque material.
	Opaque AlphaMode = iota
	// Mask corresponds to a masked material.
	Mask
	// Blend corresponds to a Blend material.
	Blend
)

// UnmarshalJSON unmarshal the alpha mode with the correct default values.
func (a *AlphaMode) UnmarshalJSON(data []byte) error {
	var tmp string
	err := json.Unmarshal(data, &tmp)
	if err == nil {
		*a = map[string]AlphaMode{
			"OPAQUE": Opaque,
			"MASK":   Mask,
			"BLEND":  Blend,
		}[tmp]
	}
	return err
}

// MarshalJSON marshal the alpha mode with the correct default values.
func (a *AlphaMode) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[AlphaMode]string{
		Opaque: "OPAQUE",
		Mask:   "MASK",
		Blend:  "BLEND",
	}[*a])
}

// MagFilter is the magnification filter.
type MagFilter uint16

const (
	// MagLinear corresponds to a linear magnification filter.
	MagLinear MagFilter = iota
	// MagNearest corresponds to a nearest magnification filter.
	MagNearest
)

// UnmarshalJSON unmarshal the mag filter with the correct default values.
func (m *MagFilter) UnmarshalJSON(data []byte) error {
	var tmp uint16
	err := json.Unmarshal(data, &tmp)
	if err == nil {
		*m = map[uint16]MagFilter{
			9728: MagNearest,
			9729: MagLinear,
		}[tmp]
	}
	return err
}

// MarshalJSON marshal the alpha mode with the correct default values.
func (m *MagFilter) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[MagFilter]uint16{
		MagNearest: 9728,
		MagLinear:  9729,
	}[*m])
}

// MinFilter is the minification filter.
type MinFilter uint16

const (
	// MinLinear corresponds to a linear minification filter.
	MinLinear MinFilter = iota
	// MinNearestMipMapLinear corresponds to a nearest mipmap linear minification filter.
	MinNearestMipMapLinear
	// MinNearest corresponds to a nearest minification filter.
	MinNearest
	// MinNearestMipMapNearest corresponds to a nearest mipmap nearest minification filter.
	MinNearestMipMapNearest
	// MinLinearMipMapNearest corresponds to a linear mipmap nearest minification filter.
	MinLinearMipMapNearest
	// MinLinearMipMapLinear corresponds to a linear mipmap linear minification filter.
	MinLinearMipMapLinear
)

// UnmarshalJSON unmarshal the min filter with the correct default values.
func (m *MinFilter) UnmarshalJSON(data []byte) error {
	var tmp uint16
	err := json.Unmarshal(data, &tmp)
	if err == nil {
		*m = map[uint16]MinFilter{
			9728: MinNearest,
			9729: MinLinear,
			9984: MinNearestMipMapNearest,
			9985: MinLinearMipMapNearest,
			9986: MinNearestMipMapLinear,
			9987: MinLinearMipMapLinear,
		}[tmp]
	}
	return err
}

// MarshalJSON marshal the min filter with the correct default values.
func (m *MinFilter) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[MinFilter]uint16{
		MinNearest:              9728,
		MinLinear:               9729,
		MinNearestMipMapNearest: 9984,
		MinLinearMipMapNearest:  9985,
		MinNearestMipMapLinear:  9986,
		MinLinearMipMapLinear:   9987,
	}[*m])
}

// WrappingMode is the wrapping mode of a texture.
type WrappingMode uint16

const (
	// Repeat corresponds to a repeat wrapping.
	Repeat WrappingMode = iota
	// ClampToEdge corresponds to a clamp to edge wrapping.
	ClampToEdge
	// MirroredRepeat corresponds to a mirrored repeat wrapping.
	MirroredRepeat
)

// UnmarshalJSON unmarshal the wrapping mode with the correct default values.
func (w *WrappingMode) UnmarshalJSON(data []byte) error {
	var tmp uint16
	err := json.Unmarshal(data, &tmp)
	if err == nil {
		*w = map[uint16]WrappingMode{
			33071: ClampToEdge,
			33648: MirroredRepeat,
			10497: Repeat,
		}[tmp]
	}
	return err
}

// MarshalJSON marshal the wrapping mode with the correct default values.
func (w *WrappingMode) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[WrappingMode]uint16{
		ClampToEdge:    33071,
		MirroredRepeat: 33648,
		Repeat:         10497,
	}[*w])
}

// Interpolation algorithm.
type Interpolation uint8

const (
	// Linear corresponds to a linear interpolation.
	Linear Interpolation = iota
	// Step corresponds to a step interpolation.
	Step
	// CubicSpline corresponds to a cubic spline interpolation.
	CubicSpline
)

// UnmarshalJSON unmarshal the interpolation with the correct default values.
func (i *Interpolation) UnmarshalJSON(data []byte) error {
	var tmp string
	err := json.Unmarshal(data, &tmp)
	if err == nil {
		*i = map[string]Interpolation{
			"LINEAR":      Linear,
			"STEP":        Step,
			"CUBICSPLINE": CubicSpline,
		}[tmp]
	}
	return err
}

// MarshalJSON marshal the interpolation with the correct default values.
func (i *Interpolation) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[Interpolation]string{
		Linear:      "LINEAR",
		Step:        "STEP",
		CubicSpline: "CUBICSPLINE",
	}[*i])
}

// TRSProperty defines a local space transformation.
// TRSproperties are converted to matrices and postmultiplied in the T * R * S order to compose the transformation matrix.
type TRSProperty uint8

const (
	// Translation corresponds to a translation transform.
	Translation TRSProperty = iota
	// Rotation corresponds to a rotation transform.
	Rotation
	// Scale corresponds to a scale transform.
	Scale
	// Weights corresponds to a weights transform.
	Weights
)

// UnmarshalJSON unmarshal the TRSProperty with the correct default values.
func (t *TRSProperty) UnmarshalJSON(data []byte) error {
	var tmp string
	err := json.Unmarshal(data, &tmp)
	if err == nil {
		*t = map[string]TRSProperty{
			"translation": Translation,
			"rotation":    Rotation,
			"scale":       Scale,
			"weights":     Weights,
		}[tmp]
	}
	return err
}

// MarshalJSON marshal the TRSProperty with the correct default values.
func (t *TRSProperty) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[TRSProperty]string{
		Translation: "translation",
		Rotation:    "rotation",
		Scale:       "scale",
		Weights:     "weights",
	}[*t])
}

const (
	glbHeaderMagic = 0x46546c67
	glbChunkJSON   = 0x4e4f534a
	glbChunkBIN    = 0x004e4942
)

type chunkHeader struct {
	Length uint32
	Type   uint32
}

type glbHeader struct {
	Magic      uint32
	Version    uint32
	Length     uint32
	JSONHeader chunkHeader
}

const (
	mimetypeApplicationOctet = "data:application/octet-stream;base64"
	mimetypeImagePNG         = "data:image/png;base64"
	mimetypeImageJPG         = "data:image/jpeg;base64"
)
