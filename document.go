package gltf

type extensible struct {
	Extensions interface{}            `json:"extensions,omitempty"` // Dictionary object with extension-specific objects.
	Extras     map[string]interface{} `json:"extras,omitempty"`     // Application-specific data.
}

type named struct {
	Name string `json:"name,omitempty"` // The user-defined name of this object.
}

// An Asset is metadata about the glTF asset.
type Asset struct {
	extensible
	Copyright  string `json:"copyright,omitempty"`         // A copyright message suitable for display to credit the content creator.
	Generator  string `json:"generator,omitempty"`         // Tool that generated this glTF model. Useful for debugging.
	Version    string `json:"version" validate:"required"` // The glTF version that this asset targets.
	MinVersion string `json:"minVersion,omitempty"`        // The minimum glTF version that this asset targets.
}

// Document defines the root object for a glTF asset.
type Document struct {
	extensible
	ExtensionsUsed     []string     `json:"extensionsUsed,omitempty"`
	ExtensionsRequired []string     `json:"extensionsRequired,omitempty"`
	Accessors          []Accessor   `json:"accessors,omitempty" validate:"dive"`
	Animations         []Animation  `json:"animations,omitempty" validate:"dive"`
	Asset              Asset        `json:"asset"`
	Buffers            []Buffer     `json:"buffers,omitempty" validate:"dive"`
	BufferViews        []BufferView `json:"bufferViews,omitempty" validate:"dive"`
	Cameras            []Camera     `json:"cameras,omitempty" validate:"dive"`
	Images             []Image      `json:"images,omitempty" validate:"dive"`
	Materials          []Material   `json:"materials,omitempty" validate:"dive"`
	Meshes             []Mesh       `json:"meshes,omitempty" validate:"dive"`
	Nodes              []Node       `json:"nodes,omitempty" validate:"dive"`
	Samplers           []Sampler    `json:"samplers,omitempty" validate:"dive"`
	Scene              uint32       `json:"scene,omitempty"`
	Scenes             []Scene      `json:"scenes,omitempty" validate:"dive"`
	Skins              []Skin       `json:"skins,omitempty" validate:"dive"`
	Textures           []Texture    `json:"textures,omitempty" validate:"dive"`
}

//The ComponentType is the datatype of components in the attribute. All valid values correspond to WebGL enums.
// The corresponding typed arrays are Int8Array, Uint8Array, Int16Array, Uint16Array, Uint32Array, and Float32Array, respectively.
// 5125 (UNSIGNED_INT) is only allowed when the accessor contains indices.
type ComponentType uint16

const (
	Byte          ComponentType = 5120
	UnsignedByte                = 5121
	Short                       = 5122
	UnsignedShort               = 5123
	UnsignedInt                 = 5125
	Float                       = 5126
)

// AccessorType specifies if the attribute is a scalar, vector, or matrix.
type AccessorType string

const (
	Scalar AccessorType = "SCALAR"
	Vec2                = "VEC2"
	Vec3                = "VEC3"
	Vec4                = "VEC4"
	Mat2                = "MAT2"
	Mat3                = "MAT3"
	Mat4                = "MAT4"
)

// An accessor is a typed view into a bufferView. A bufferView contains raw binary data.
// An accessor provides a typed view into a bufferView or a subset of a bufferView
// similar to how WebGL's vertexAttribPointer() defines an attribute in a buffer.
type Accessor struct {
	named
	extensible
	BufferView    uint32        `json:"bufferView,omitempty"`
	ByteOffset    uint32        `json:"byteOffset,omitempty"`
	ComponentType ComponentType `json:"componentType" validate:"oneof=5120 5121 5122 5123 5125 5126"`
	Normalized    bool          `json:"normalized,omitempty"`      // Specifies whether integer data values should be normalized.
	Count         uint32        `json:"count" validate:"required"` // The number of attributes referenced by this accessor.
	Type          AccessorType  `json:"type" validate:"oneof=SCALAR VEC2 VEC3 VEC4 MAT2 MAT3 MAT4"`
	Max           []float32     `json:"max,omitempty" validate:"omitempty,lte=16"` // Maximum value of each component in this attribute.
	Min           []float32     `json:"min,omitempty" validate:"omitempty,lte=16"` // Minimum value of each component in this attribute.
	Sparce        *Sparse       `json:"sparce,omitempty"`                          // Sparse storage of attributes that deviate from their initialization value.
}

// Sparse storage of attributes that deviate from their initialization value.
type Sparse struct {
	extensible
	Count   uint32        `json:"count" validate:"gte=1"` // Number of entries stored in the sparse array.
	Indices SparseIndices `json:"indices"`                // Index array of size count that points to those accessor attributes that deviate from their initialization value.
	Values  SparseValues  `json:"values"`                 // Array of size count times number of components, storing the displaced accessor attributes pointed by indices.
}

// SparseValues stores the displaced accessor attributes pointed by accessor.sparse.indices.
type SparseValues struct {
	extensible
	BufferView uint32 `json:"bufferView"`
	ByteOffset uint32 `json:"byteOffset,omitempty"`
}

// SparseIndices defines the indices of those attributes that deviate from their initialization value.
type SparseIndices struct {
	extensible
	BufferView    uint32        `json:"bufferView"`
	ByteOffset    uint32        `json:"byteOffset,omitempty"`
	ComponentType ComponentType `json:"componentType" validate:"oneof=5121 5123 5125"`
}

// The target that the GPU buffer should be bound to.
type Target uint16

const (
	ArrayBuffer        Target = 34962
	ElementArrayBuffer        = 34963
)

// A buffer points to binary geometry, animation, or skins.
type Buffer struct {
	named
	extensible
	URI        string `json:"uri,omitempty" validate:"omitempty,uri|datauri"`
	ByteLength uint32 `json:"byteLength" validate:"required"`
}

// BufferView is a view into a buffer generally representing a subset of the buffer.
type BufferView struct {
	extensible
	Buffer     uint32 `json:"buffer"`
	ByteOffset uint32 `json:"byteOffset,omitempty"`
	ByteLength uint32 `json:"byteLength" validate:"required"`
	ByteStride uint32 `json:"byteStride,omitempty" validate:"omitempty,gte=4,lte=252"`
	Target     Target `json:"target,omitempty" validate:"omitempty,oneof=34962 34963"`
}

// The Scene contains a list of root nodes.
type Scene struct {
	named
	extensible
	Nodes []uint32 `json:"nodes,omitempty" validate:"omitempty,unique"`
}

// A node in the node hierarchy.
// A node can have either a matrix or any combination of translation/rotation/scale (TRS) properties.
type Node struct {
	named
	extensible
	Camera      uint32      `json:"camera,omitempty"`
	Children    []uint32    `json:"children,omitempty" validate:"omitempty,unique"`
	Skin        uint32      `json:"skin,omitempty"`
	Matrix      [16]float32 `json:"matrix,omitempty"` // A 4x4 transformation matrix stored in column-major order.
	Mesh        uint32      `json:"mesh,omitempty"`
	Rotation    [4]float64  `json:"rotation" validate:"omitempty,dive,gte=-1,lte=1"` // The node's unit quaternion rotation in the order (x, y, z, w), where w is the scalar.
	Scale       [3]float32  `json:"scale,omitempty"`
	Translation [3]float32  `json:"translation,omitempty"`
	Weights     []float32   `json:"weights,omitempty"` // The weights of the instantiated Morph Target.
}

// Skin defines joints and matrices.
type Skin struct {
	named
	extensible
	InverseBindMatrices uint32   `json:"inverseBindMatrices,omitempty"`      // The index of the accessor containing the floating-point 4x4 inverse-bind matrices.
	Skeleton            uint32   `json:"skeleton,omitempty"`                 // The index of the node used as a skeleton root. When undefined, joints transforms resolve to scene root.
	Joints              []uint32 `json:"joints" validate:"omitempty,unique"` // Indices of skeleton nodes, used as joints in this skin.
}

// CameraType specifies if the camera uses a perspective or orthographic projection.
// Based on this, either the camera's perspective or orthographic property will be defined.
type CameraType string

const (
	PerspectiveType CameraType = "perspective"
	OrtographicType            = "orthographic"
)

// A Camera's projection. A node can reference a camera to apply a transform to place the camera in the scene.
type Camera struct {
	named
	extensible
	Ortographic *Ortographic `json:"orthographic,omitempty"`
	Perspective *Perspective `json:"perspective,omitempty"`
	Type        CameraType   `json:"type" validate:"oneof=perspective orthographic"`
}

// An Orthographic camera containing properties to create an orthographic projection matrix.
type Ortographic struct {
	extensible
	Xmag  float32 `json:"xmag"`                               // The horizontal magnification of the view.
	Ymag  float32 `json:"ymag"`                               // The vertical magnification of the view.
	Zfar  float32 `json:"zfar" validate:"gt=0,gtfield=Znear"` // The distance to the far clipping plane.
	Znear float32 `json:"znear" validate:"gte=0"`             // The distance to the near clipping plane.
}

// A perspective camera containing properties to create a perspective projection matrix.
type Perspective struct {
	extensible
	AspectRatio float32 `json:"aspectRatio,omitempty"`
	Yfov        float32 `json:"yfov"`           // The vertical field of view in radians.
	Zfar        float32 `json:"zfar,omitempty"` // The distance to the far clipping plane.
	Znear       float32 `json:"znear"`          // The distance to the near clipping plane.
}

// Each key corresponds to mesh attribute semantic and each value is the index of the accessor containing attribute's data.
type Attribute = map[string]uint32

// PrimitiveMode defines the type of primitives to render. All valid values correspond to WebGL enums.
type PrimitiveMode uint8

const (
	Points         PrimitiveMode = 0
	Lines                        = 1
	Line_Loop                    = 2
	Lin_Strip                    = 3
	Triangles                    = 4
	Triangle_Strip               = 5
	Triangle_Fan                 = 6
)

// A Mesh is a set of primitives to be rendered. A node can contain one mesh. A node's transform places the mesh in the scene.
type Mesh struct {
	named
	extensible
	Primitives []Primitive `json:"primitives" validate:"required,gt=0,dive"`
	Weights    []float32   `json:"weights,omitempty"`
}

// Geometry to be rendered with the given material.
type Primitive struct {
	extensible
	Attributes Attribute     `json:"attributes"`
	Indices    uint32        `json:"indices,omitempty"` // The index of the accessor that contains the indices.
	Material   uint32        `json:"material,omitempty"`
	Mode       PrimitiveMode `json:"mode" validate:"lte=6"`
	Targets    []Attribute   `json:"targets,omitempty" validate:"omitempty,dive,dive,keys,oneof=POSITION NORMAL TANGENT,endkeys"` // Only POSITION, NORMAL, and TANGENT supported.
}

// The AlphMode enumeration specifying the interpretation of the alpha value of the main factor and texture.
type AlphaMode string

const (
	Opaque AlphaMode = "OPAQUE"
	Mask             = "MASK"
	Blend            = "BLEND"
)

// MagFilter is the magnification filter.
type MagFilter uint16

const (
	MagNearest MagFilter = 9728
	MagLinear            = 9729
)

// MinFilter is the minification filter.
type MinFilter uint16

const (
	MinNearest              MinFilter = 9728
	MinLinear                         = 9729
	MinNearestMipMapNearest           = 9984
	MinLinearMipMapNearest            = 9985
	MinNearestMipMapLinear            = 9986
	MinLinearMipMapLinear             = 9987
)

// WrappingMode is the wrapping mode of a texture.
type WrappingMode uint16

const (
	ClampToEdge    WrappingMode = 33071
	MirroredRepeat              = 33648
	Repeat                      = 10497
)

// The material appearance of a primitive.
type Material struct {
	named
	extensible
	PBRMetallicRoughness *PBRMetallicRoughness `json:"pbrMetallicRoughness,omitempty"`
	NormalTexture        *NormalTexture        `json:"normalTexture,omitempty"`
	OcclusionTexture     *OcclusionTexture     `json:"occlusionTexture,omitempty"`
	EmissiveTexture      *TextureInfo          `json:"emissiveTexture,omitempty"`
	EmissiveFactor       [3]float32            `json:"emissiveFactor,omitempty" validate:"dive,gte=0,lte=1"`
	AlphaMode            AlphaMode             `json:"alphaMode,omitempty" validate:"oneof=OPAQUE MASK BLEND"`
	AlphaCutoff          string                `json:"alphaCutoff"`
	DoubleSided          bool                  `json:"doubleSided,omitempty"`
}

// A NormalTexture references to a normal texture.
type NormalTexture struct {
	TextureInfo
	Scale float32 `json:"scale"`
}

// An OcclusionTexture references to an occlusion texture
type OcclusionTexture struct {
	TextureInfo
	Strength float32 `json:"strength" validate:"gte=0,lte=1"`
}

// A set of parameter values that are used to define the metallic-roughness material model from Physically-Based Rendering (PBR) methodology.
type PBRMetallicRoughness struct {
	extensible
	BaseColorFactor          [4]float32   `json:"baseColorFactor" validate:"dive,gte=0,lte=1"`
	BaseColorTexture         *TextureInfo `json:"baseColorTexture,omitempty"`
	MetallicFactor           float32      `json:"metallicFactor" validate:"gte=0,lte=1"`
	RoughnessFactor          float32      `json:"roughnessFactor" validate:"gte=0,lte=1"`
	MetallicRoughnessTexture *TextureInfo `json:"metallicRoughnessTexture,omitempty"`
}

// TextureInfo references to a texture.
type TextureInfo struct {
	extensible
	Index    uint32 `json:"index"`
	TexCoord uint32 `json:"texCoord,omitempty"` // The index of texture's TEXCOORD attribute used for texture coordinate mapping.
}

// A texture and its sampler.
type Texture struct {
	named
	extensible
	Sampler uint32 `json:"sampler,omitempty"`
	Source  uint32 `json:"source,omitempty"`
}

// Sampler of a texture for filtering and wrapping modes.
type Sampler struct {
	named
	extensible
	MagFilter MagFilter    `json:"magFilter" validate:"omitempty,oneof=9728 9729"`
	MinFilter MinFilter    `json:"minFilter" validate:"omitempty,oneof=9728 9729 9984 9985 9986 9987"`
	WrapS     WrappingMode `json:"wrapS" validate:"omitempty,oneof=33071 33648 10497"`
	WrapT     WrappingMode `json:"wrapT" validate:"omitempty,oneof=33071 33648 10497"`
}

// Image data used to create a texture. Image can be referenced by URI or bufferView index.
// mimeType is required in the latter case.
type Image struct {
	named
	extensible
	URI        string `json:"uri,omitempty" validate:"omitempty,uri|datauri"`
	MimeType   string `json:"mimeType,omitempty" validate:"omitempty,oneof=image/jpeg image/png"` // Manadatory if BufferView is defined.
	BufferView uint32 `json:"bufferView,omitempty"`                                               // Use this instead of the image's uri property.
}

// Interpolation algorithm.
type Interpolation string

const (
	Linear      Interpolation = "LINEAR"
	Step                      = "STEP"
	CubicSpline               = "CUBICSPLINE"
)

// An Animation keyframe.
type Animation struct {
	named
	extensible
	Channels []Channel          `json:"channels" validate:"required,gt=0,dive"`
	Samplers []AnimationSampler `json:"samplers" validate:"required,gt=0,dive"`
}

// AnimationSampler combines input and output accessors with an interpolation algorithm to define a keyframe graph (but not its target).
type AnimationSampler struct {
	extensible
	Input         uint32        `json:"input"` // The index of an accessor containing keyframe input values.
	Interpolation Interpolation `json:"interpolation,omitempty" validate:"omitempty,oneof=LINEAR STEP CUBICSPLINE"`
	Output        uint32        `json:"output"` // The index of an accessor containing keyframe output values.
}

// The channel targets an animation's sampler at a node's property.
type Channel struct {
	extensible
	Sampler uint32        `json:"sampler"`
	Target  ChannelTarget `json:"target"`
}

// ChannelTarget describes the index of the node and TRS property that an animation channel targets.
// The Path represents the name of the node's TRS property to modify, or the "weights" of the Morph Targets it instantiates.
// For the "translation" property, the values that are provided by the sampler are the translation along the x, y, and z axes.
// For the "rotation" property, the values are a quaternion in the order (x, y, z, w), where w is the scalar.
// For the "scale" property, the values are the scaling factors along the x, y, and z axes.
type ChannelTarget struct {
	extensible
	Node uint32 `json:"node,omitempty"`
	Path string `json:"path" validate:"oneof=translation rotation scale weights"`
}
