//go:generate stringer -linecomment -type ComponentType,AccessorType,PrimitiveMode,AlphaMode,MagFilter,MinFilter,WrappingMode,Interpolation,Target,TRSProperty -output const_strings.go

package gltf

import (
	"encoding/json"
	"fmt"
)

var (
	// DefaultMatrix defines an identity matrix.
	DefaultMatrix = [16]float32{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1}
	// DefaultRotation defines a quaternion without rotation.
	DefaultRotation = [4]float32{0, 0, 0, 1}
	// DefaultScale defines a scaling that does not modify the size of the object.
	DefaultScale = [3]float32{1, 1, 1}
	// DefaultTranslation defines a translation that does not move the object.
	DefaultTranslation = [3]float32{0, 0, 0}
)

// Attribute names defined in the spec.
const (
	POSITION   = "POSITION"
	NORMAL     = "NORMAL"
	TANGENT    = "TANGENT"
	TEXCOORD_0 = "TEXCOORD_0"
	TEXCOORD_1 = "TEXCOORD_1"
	WEIGHTS_0  = "WEIGHTS_0"
	JOINTS_0   = "JOINTS_0"
	COLOR_0    = "COLOR_0"
)

var (
	emptyMatrix   = [16]float32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	emptyRotation = [4]float32{0, 0, 0, 0}
	emptyScale    = [3]float32{0, 0, 0}
)

// The ComponentType is the datatype of components in the attribute. All valid values correspond to WebGL enums.
// 5125 (UNSIGNED_INT) is only allowed when the accessor contains indices.
type ComponentType uint16

const (
	ComponentFloat  ComponentType = iota // FLOAT
	ComponentByte                        // BYTE
	ComponentUbyte                       // UNSIGNED_BYTE
	ComponentShort                       // SHORT
	ComponentUshort                      // UNSIGNED_SHORT
	ComponentUint                        // UNSIGNED_INT
)

// ByteSize returns the size of a component in bytes.
func (c ComponentType) ByteSize() uint32 {
	return map[ComponentType]uint32{
		ComponentByte:   1,
		ComponentUbyte:  1,
		ComponentShort:  2,
		ComponentUshort: 2,
		ComponentUint:   4,
		ComponentFloat:  4,
	}[c]
}

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
		ComponentByte:   5120,
		ComponentUbyte:  5121,
		ComponentShort:  5122,
		ComponentUshort: 5123,
		ComponentUint:   5125,
		ComponentFloat:  5126,
	}[*c])
}

// AccessorType specifies if the attribute is a scalar, vector, or matrix.
type AccessorType uint8

const (
	AccessorScalar AccessorType = iota // SCALAR
	AccessorVec2                       // VEC2
	AccessorVec3                       // VEC3
	AccessorVec4                       // VEC4
	AccessorMat2                       // MAT2
	AccessorMat3                       // MAT3
	AccessorMat4                       // MAT4
)

// Components returns the number of components of an accessor type.
func (a AccessorType) Components() uint32 {
	return map[AccessorType]uint32{
		AccessorScalar: 1,
		AccessorVec2:   2,
		AccessorVec3:   3,
		AccessorVec4:   4,
		AccessorMat2:   4,
		AccessorMat3:   9,
		AccessorMat4:   16,
	}[a]
}

// UnmarshalJSON unmarshal the accessor type with the correct default values.
func (a *AccessorType) UnmarshalJSON(data []byte) error {
	var tmp string
	err := json.Unmarshal(data, &tmp)
	if err == nil {
		accType, ok := map[string]AccessorType{
			"SCALAR": AccessorScalar,
			"VEC2":   AccessorVec2,
			"VEC3":   AccessorVec3,
			"VEC4":   AccessorVec4,
			"MAT2":   AccessorMat2,
			"MAT3":   AccessorMat3,
			"MAT4":   AccessorMat4,
		}[tmp]
		if !ok {
			return fmt.Errorf("gltf: unknown accessor's type: %s", tmp)
		}
		*a = accType
	}
	return err
}

// MarshalJSON marshal the accessor type with the correct default values.
func (a *AccessorType) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[AccessorType]string{
		AccessorScalar: "SCALAR",
		AccessorVec2:   "VEC2",
		AccessorVec3:   "VEC3",
		AccessorVec4:   "VEC4",
		AccessorMat2:   "MAT2",
		AccessorMat3:   "MAT3",
		AccessorMat4:   "MAT4",
	}[*a])
}

// The Target that the GPU buffer should be bound to.
type Target uint16

const (
	TargetNone               Target = 0     // NONE
	TargetArrayBuffer        Target = 34962 // ARRAY_BUFFER
	TargetElementArrayBuffer Target = 34963 // ELEMENT_ARRAY_BUFFER
)

// Attribute is a map that each key corresponds to mesh attribute semantic and each value is the index of the accessor containing attribute's data.
type Attribute = map[string]uint32

// PrimitiveMode defines the type of primitives to render. All valid values correspond to WebGL enums.
type PrimitiveMode uint8

const (
	PrimitiveTriangles     PrimitiveMode = iota // TRIANGLES
	PrimitivePoints                             // POINTS
	PrimitiveLines                              // LINES
	PrimitiveLineLoop                           // LINE_LOOP
	PrimitiveLineStrip                          // LINE_STRIP
	PrimitiveTriangleStrip                      // TRIANGLE_STRIP
	PrimitiveTriangleFan                        // TRIANGLE_FAN
)

// UnmarshalJSON unmarshal the primitive mode with the correct default values.
func (p *PrimitiveMode) UnmarshalJSON(data []byte) error {
	var tmp uint8
	err := json.Unmarshal(data, &tmp)
	if err == nil {
		*p = map[uint8]PrimitiveMode{
			0: PrimitivePoints,
			1: PrimitiveLines,
			2: PrimitiveLineLoop,
			3: PrimitiveLineStrip,
			4: PrimitiveTriangles,
			5: PrimitiveTriangleStrip,
			6: PrimitiveTriangleFan,
		}[tmp]
	}
	return err
}

// MarshalJSON marshal the primitive mode with the correct default values.
func (p *PrimitiveMode) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[PrimitiveMode]uint8{
		PrimitivePoints:        0,
		PrimitiveLines:         1,
		PrimitiveLineLoop:      2,
		PrimitiveLineStrip:     3,
		PrimitiveTriangles:     4,
		PrimitiveTriangleStrip: 5,
		PrimitiveTriangleFan:   6,
	}[*p])
}

// The AlphaMode enumeration specifying the interpretation of the alpha value of the main factor and texture.
type AlphaMode uint8

const (
	AlphaOpaque AlphaMode = iota // OPAQUE
	AlphaMask                    // MASK
	AlphaBlend                   // BLEND
)

// UnmarshalJSON unmarshal the alpha mode with the correct default values.
func (a *AlphaMode) UnmarshalJSON(data []byte) error {
	var tmp string
	err := json.Unmarshal(data, &tmp)
	if err == nil {
		*a = map[string]AlphaMode{
			"OPAQUE": AlphaOpaque,
			"MASK":   AlphaMask,
			"BLEND":  AlphaBlend,
		}[tmp]
	}
	return err
}

// MarshalJSON marshal the alpha mode with the correct default values.
func (a *AlphaMode) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[AlphaMode]string{
		AlphaOpaque: "OPAQUE",
		AlphaMask:   "MASK",
		AlphaBlend:  "BLEND",
	}[*a])
}

// MagFilter is the magnification filter.
type MagFilter uint16

const (
	MagUndefined MagFilter = iota // UNDEFINED
	MagLinear                     // LINEAR
	MagNearest                    // NEAREST
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
	MinUndefined            MinFilter = iota // UNDEFINED
	MinLinear                                // LINEAR
	MinNearestMipMapLinear                   // NEAREST_MIPMAP_LINEAR
	MinNearest                               // NEAREST
	MinNearestMipMapNearest                  // NEAREST_MIPMAP_NEAREST
	MinLinearMipMapNearest                   // LINEAR_MIPMAP_NEAREST
	MinLinearMipMapLinear                    // LINEAR_MIPMAP_LINEAR
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
	WrapRepeat         WrappingMode = iota // REPEAT
	WrapClampToEdge                        // CLAMP_TO_EDGE
	WrapMirroredRepeat                     // MIRRORED_REPEAT
)

// UnmarshalJSON unmarshal the wrapping mode with the correct default values.
func (w *WrappingMode) UnmarshalJSON(data []byte) error {
	var tmp uint16
	err := json.Unmarshal(data, &tmp)
	if err == nil {
		*w = map[uint16]WrappingMode{
			33071: WrapClampToEdge,
			33648: WrapMirroredRepeat,
			10497: WrapRepeat,
		}[tmp]
	}
	return err
}

// MarshalJSON marshal the wrapping mode with the correct default values.
func (w *WrappingMode) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[WrappingMode]uint16{
		WrapClampToEdge:    33071,
		WrapMirroredRepeat: 33648,
		WrapRepeat:         10497,
	}[*w])
}

// Interpolation algorithm.
type Interpolation uint8

const (
	InterpolationLinear      Interpolation = iota // LINEAR
	InterpolationStep                             // STEP
	InterpolationCubicSpline                      // CUBICSPLINE
)

// UnmarshalJSON unmarshal the interpolation with the correct default values.
func (i *Interpolation) UnmarshalJSON(data []byte) error {
	var tmp string
	err := json.Unmarshal(data, &tmp)
	if err == nil {
		*i = map[string]Interpolation{
			"LINEAR":      InterpolationLinear,
			"STEP":        InterpolationStep,
			"CUBICSPLINE": InterpolationCubicSpline,
		}[tmp]
	}
	return err
}

// MarshalJSON marshal the interpolation with the correct default values.
func (i *Interpolation) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[Interpolation]string{
		InterpolationLinear:      "LINEAR",
		InterpolationStep:        "STEP",
		InterpolationCubicSpline: "CUBICSPLINE",
	}[*i])
}

// TRSProperty defines a local space transformation.
// TRSproperties are converted to matrices and postmultiplied in the T * R * S order to compose the transformation matrix.
type TRSProperty uint8

const (
	TRSTranslation TRSProperty = iota // translation
	TRSRotation                       // rotation
	TRSScale                          // scale
	TRSWeights                        // weights
)

// UnmarshalJSON unmarshal the TRSProperty with the correct default values.
func (t *TRSProperty) UnmarshalJSON(data []byte) error {
	var tmp string
	err := json.Unmarshal(data, &tmp)
	if err == nil {
		*t = map[string]TRSProperty{
			"translation": TRSTranslation,
			"rotation":    TRSRotation,
			"scale":       TRSScale,
			"weights":     TRSWeights,
		}[tmp]
	}
	return err
}

// MarshalJSON marshal the TRSProperty with the correct default values.
func (t *TRSProperty) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[TRSProperty]string{
		TRSTranslation: "translation",
		TRSRotation:    "rotation",
		TRSScale:       "scale",
		TRSWeights:     "weights",
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
