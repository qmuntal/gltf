package gltf

// The ComponentType is the datatype of components in the attribute. All valid values correspond to WebGL enums.
// 5125 (UNSIGNED_INT) is only allowed when the accessor contains indices.
type ComponentType uint16

const (
	// Byte corresponds to a Int8Array.
	Byte ComponentType = 5120
	// UnsignedByte corresponds to a Uint8Array.
	UnsignedByte = 5121
	// Short corresponds to a Int16Array.
	Short = 5122
	// UnsignedShort corresponds to a Uint16Array.
	UnsignedShort = 5123
	// UnsignedInt corresponds to a Uint32Array.
	UnsignedInt = 5125
	// Float corresponds to a Float32Array.
	Float = 5126
)

// AccessorType specifies if the attribute is a scalar, vector, or matrix.
type AccessorType string

const (
	// Scalar corresponds to a single dimension value.
	Scalar AccessorType = "SCALAR"
	// Vec2 corresponds to a two dimensions array.
	Vec2 = "VEC2"
	// Vec3 corresponds to a three dimensions array.
	Vec3 = "VEC3"
	// Vec4 corresponds to a four dimensions array.
	Vec4 = "VEC4"
	// Mat2 corresponds to a 2x2 matrix.
	Mat2 = "MAT2"
	// Mat3 corresponds to a 3x3 matrix.
	Mat3 = "MAT3"
	// Mat4 corresponds to a 4x4 matrix.
	Mat4 = "MAT4"
)

// The Target that the GPU buffer should be bound to.
type Target uint16

const (
	// ArrayBuffer corresponds to an array buffer.
	ArrayBuffer Target = 34962
	// ElementArrayBuffer corresponds to an element array buffer.
	ElementArrayBuffer = 34963
)

// CameraType specifies if the camera uses a perspective or orthographic projection.
// Based on this, either the camera's perspective or orthographic property will be defined.
type CameraType string

const (
	// PerspectiveType corresponds to a perspective camera.
	PerspectiveType CameraType = "perspective"
	// OrthographicType corresponds to an orthographic camera.
	OrthographicType = "orthographic"
)

// Attribute is a map that each key corresponds to mesh attribute semantic and each value is the index of the accessor containing attribute's data.
type Attribute = map[string]uint32

// PrimitiveMode defines the type of primitives to render. All valid values correspond to WebGL enums.
type PrimitiveMode uint8

const (
	// Points corresponds to a Point primitive.
	Points PrimitiveMode = 0
	// Lines corresponds to a Line primitive.
	Lines = 1
	// LineLoop corresponds to a Line Loop primitive.
	LineLoop = 2
	// LineStrip corresponds to a Line Strip primitive.
	LineStrip = 3
	// Triangles corresponds to a Triangle primitive.
	Triangles = 4
	// TriangleStrip corresponds to a Triangle Strip primitive.
	TriangleStrip = 5
	// TriangleFan corresponds to a Triangle Fan primitive.
	TriangleFan = 6
)

// The AlphaMode enumeration specifying the interpretation of the alpha value of the main factor and texture.
type AlphaMode string

const (
	// Opaque corresponds to an Opaque material.
	Opaque AlphaMode = "OPAQUE"
	// Mask corresponds to a masked material.
	Mask = "MASK"
	// Blend corresponds to a Blend material.
	Blend = "BLEND"
)

// MagFilter is the magnification filter.
type MagFilter uint16

const (
	// MagNearest corresponds to a nearest magnification filter.
	MagNearest MagFilter = 9728
	// MagLinear corresponds to a linear magnification filter.
	MagLinear = 9729
)

// MinFilter is the minification filter.
type MinFilter uint16

const (
	// MinNearest corresponds to a nearest minification filter.
	MinNearest MinFilter = 9728
	// MinLinear corresponds to a linear minification filter.
	MinLinear = 9729
	// MinNearestMipMapNearest corresponds to a nearest mipmap nearest minification filter.
	MinNearestMipMapNearest = 9984
	// MinLinearMipMapNearest corresponds to a linear mipmap nearest minification filter.
	MinLinearMipMapNearest = 9985
	// MinNearestMipMapLinear corresponds to a nearest mipmap linear minification filter.
	MinNearestMipMapLinear = 9986
	// MinLinearMipMapLinear corresponds to a linear mipmap linear minification filter.
	MinLinearMipMapLinear = 9987
)

// WrappingMode is the wrapping mode of a texture.
type WrappingMode uint16

const (
	// ClampToEdge corresponds to a clamp to edge wrapping.
	ClampToEdge WrappingMode = 33071
	// MirroredRepeat corresponds to a mirrored repeat wrapping.
	MirroredRepeat = 33648
	// Repeat corresponds to a repeat wrapping.
	Repeat = 10497
)

// Interpolation algorithm.
type Interpolation string

const (
	// Linear corresponds to a linear interpolation.
	Linear Interpolation = "LINEAR"
	// Step corresponds to a step interpolation.
	Step = "STEP"
	// CubicSpline corresponds to a cubic spline interpolation.
	CubicSpline = "CUBICSPLINE"
)

// TRSProperty defines a local space transformation.
// TRSproperties are converted to matrices and postmultiplied in the T * R * S order to compose the transformation matrix.
type TRSProperty string

const (
	// Translation corresponds to a translation transform.
	Translation TRSProperty = "translation"
	// Rotation corresponds to a rotation transform.
	Rotation = "rotation"
	// Scale corresponds to a scale transform.
	Scale = "scale"
	// Weights corresponds to a weights transform.
	Weights = "weights"
)

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
