package gltf

// An Asset is metadata about the glTF asset.
type Asset struct {
	Extensible
	Copyright  string `json:"copyright,omitempty"`          // A copyright message suitable for display to credit the content creator.
	Generator  string `json:"generator,omitempty"`          // Tool that generated this glTF model. Useful for debugging.
	Version    string `json:"version" validator:"required"` // The glTF version that this asset targets.
	MinVersion string `json:"minVersion,omitempty"`         // The minimum glTF version that this asset targets.
}

// Extensible is an object that has the Extension and Extras properties.
type Extensible struct {
	Extensions interface{}            `json:"extensions,omitempty"` // Dictionary object with extension-specific objects.
	Extras     map[string]interface{} `json:"extras,omitempty"`     // Application-specific data.
}

// Named is an object that has a Name property.
type Named struct {
	Name string `json:"name,omitempty"` // The user-defined name of this object.
}

// Document defines the root object for a glTF asset.
type Document struct {
	Extensible
	ExtensionsUsed     []string     `json:"extensionsUsed,omitempty"`
	ExtensionsRequired []string     `json:"extensionsRequired,omitempty"`
	Accessors          []Accessor   `json:"accessors,omitempty"`
	Animations         []Animation  `json:"animations,omitempty"`
	Asset              Asset        `json:"asset"`
	Buffers            []Buffer     `json:"buffers,omitempty"`
	BufferViews        []BufferView `json:"bufferViews,omitempty"`
	Cameras            []Camera     `json:"cameras,omitempty"`
	Images             []Image      `json:"images,omitempty"`
	Materials          []Material   `json:"materials,omitempty"`
	Meshes             []Mesh       `json:"meshes,omitempty"`
	Nodes              []Node       `json:"nodes,omitempty"`
	Samplers           []Sampler    `json:"samplers,omitempty"`
	Scene              uint32       `json:"scene,omitempty"`
	Scenes             []Scene      `json:"scenes,omitempty"`
	Skins              []Skin       `json:"skins,omitempty"`
	Textures           []Texture    `json:"textures,omitempty"`
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
	Named
	Extensible
	BufferView    uint32        `json:"bufferView,omitempty"`
	ByteOffset    uint32        `json:"byteOffset,omitempty"`
	ComponentType ComponentType `json:"componentType" validator:"oneof=5120 5121 5122 5123 5125 5126"`
	Normalized    bool          `json:"normalized,omitempty"`       // Specifies whether integer data values should be normalized.
	Count         uint32        `json:"count" validator:"required"` // The number of attributes referenced by this accessor.
	Type          AccessorType  `json:"type" validator:"oneof=SCALAR VEC2 VEC3 VEC4 MAT2 MAT3 MAT4"`
	Max           []float32     `json:"max,omitempty" validator:"omitempty,lte=16"` // Maximum value of each component in this attribute.
	Min           []float32     `json:"min,omitempty" validator:"omitempty,lte=16"` // Minimum value of each component in this attribute.
	Sparce        Sparse        `json:"sparce,omitempty"`                           // Sparse storage of attributes that deviate from their initialization value.
}

// Sparse storage of attributes that deviate from their initialization value.
type Sparse struct {
	Extensible
	Count   uint32        `json:"count" validator:"gte=1"` // Number of entries stored in the sparse array.
	Indices SparceIndices `json:"indices"`                 // Index array of size count that points to those accessor attributes that deviate from their initialization value.
	Values  SparseValues  `json:"values"`                  // Array of size count times number of components, storing the displaced accessor attributes pointed by indices.
}

// SparseValues stores the displaced accessor attributes pointed by accessor.sparse.indices.
type SparseValues struct {
	Extensible
	BufferView uint32 `json:"bufferView"`
	ByteOffset uint32 `json:"byteOffset,omitempty"`
}

// SparceIndices defines the indices of those attributes that deviate from their initialization value.
type SparceIndices struct {
	Extensible
	BufferView    uint32        `json:"bufferView"`
	ByteOffset    uint32        `json:"byteOffset,omitempty"`
	ComponentType ComponentType `json:"componentType" validator:"oneof=5121 5123 5125"`
}

// The target that the GPU buffer should be bound to.
type Target uint16

const (
	ArrayBuffer        Target = 34962
	ElementArrayBuffer        = 34963
)

// A buffer points to binary geometry, animation, or skins.
type Buffer struct {
	Named
	Extensible
	URI        string `json:"uri,omitempty" validator:"omitempty|uri|datauri"`
	ByteLength uint32 `json:"byteLength" validator:"required"`
}

// BufferView is a view into a buffer generally representing a subset of the buffer.
type BufferView struct {
	Extensible
	Buffer     uint32 `json:"buffer"`
	ByteOffset uint32 `json:"byteOffset,omitempty"`
	ByteLength uint32 `json:"byteLength" validator:"required"`
	ByteStride uint32 `json:"byteStride,omitempty" validator:"omitempty,gte=4,lte=252"`
	Target     Target `json:"target,omitempty" validator:"omitempty,oneof=34962 34963"`
}

// The Scene contains a list of root nodes.
type Scene struct {
	Named
	Extensible
	Nodes []uint32 `json:"nodes,omitempty" validator:"omitempty,unique"`
}

// A node in the node hierarchy.
// A node can have either a matrix or any combination of translation/rotation/scale (TRS) properties.
type Node struct {
	Named
	Extensible
	Camera      uint32      `json:"camera,omitempty"`
	Children    []uint32    `json:"children,omitempty" validator:"omitempty,unique"`
	Skin        uint32      `json:"skin,omitempty"`
	Matrix      [16]float32 `json:"matrix,omitempty"` // A 4x4 transformation matrix stored in column-major order.
	Mesh        uint32      `json:"mesh,omitempty"`
	Rotation    [4]float64  `json:"rotation" validator:"omitempty,dive,gte=-1,lte=1"` // The node's unit quaternion rotation in the order (x, y, z, w), where w is the scalar.
	Scale       [3]float32  `json:"scale,omitempty"`
	Translation [3]float32  `json:"translation,omitempty"`
	Weights     []float32   `json:"weights,omitempty"` // The weights of the instantiated Morph Target.
}

// Skin defines joints and matrices.
type Skin struct {
	Named
	Extensible
	InverseBindMatrices uint32   `json:"inverseBindMatrices,omitempty"`       // The index of the accessor containing the floating-point 4x4 inverse-bind matrices.
	Skeleton            uint32   `json:"skeleton,omitempty"`                  // The index of the node used as a skeleton root. When undefined, joints transforms resolve to scene root.
	Joints              []uint32 `json:"joints" validator:"omitempty,unique"` // Indices of skeleton nodes, used as joints in this skin.
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
	Named
	Extensible
	Ortographic Ortographic `json:"orthographic,omitempty"`
	Perspective Perspective `json:"perspective,omitempty"`
	Type        CameraType  `json:"type" validator:"oneof=perspective orthographic"`
}

// An Orthographic camera containing properties to create an orthographic projection matrix.
type Ortographic struct {
	Extensible
	Xmag  float32 `json:"xmag"`                                // The horizontal magnification of the view.
	Ymag  float32 `json:"ymag"`                                // The vertical magnification of the view.
	Zfar  float32 `json:"zfar" validator:"gt=0,gtfield=Znear"` // The distance to the far clipping plane.
	Znear float32 `json:"znear" validator:"gte=0"`             // The distance to the near clipping plane.
}

// A perspective camera containing properties to create a perspective projection matrix.
type Perspective struct {
	Extensible
	AspectRatio float32 `json:"aspectRatio,omitempty" validator:"omitempty,gt=0"`
	Yfov        float32 `json:"yfov" validator:"gt=0"`                     // The vertical field of view in radians.
	Zfar        float32 `json:"zfar,omitempty" validator:"omitempty,gt=0"` // The distance to the far clipping plane.
	Znear       float32 `json:"znear" validator:"omitempty,gt=0"`          // The distance to the near clipping plane.
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
	Named
	Extensible
	Primitives []Primitive `json:"primitives" validator:"required"`
	Weights    []float32   `json:"weights,omitempty"`
}

// Geometry to be rendered with the given material.
type Primitive struct {
	Extensible
	Attributes Attribute   `json:"attributes"`
	Indices    uint32      `json:"indices,omitempty"` // The index of the accessor that contains the indices.
	Material   uint32      `json:"material,omitempty"`
	Mode       uint32      `json:"mode" validator:"lte=6"`
	Targets    []Attribute `json:"targets,omitempty" validator:"omitempty,dive,keys,oneof=POSITION NORMAL TANGENT,endkeys"` // Only POSITION, NORMAL, and TANGENT supported.
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
	Named
	Extensible
	PBRMetallicRoughness PBRMetallicRoughness `json:"pbrMetallicRoughness,omitempty"`
	NormalTexture        NormalTexture        `json:"normalTexture,omitempty"`
	OcclusionTexture     OcclusionTexture     `json:"occlusionTexture,omitempty"`
	EmissiveTexture      TextureInfo          `json:"emissiveTexture,omitempty"`
	EmissiveFactor       [3]float32           `json:"emissiveFactor,omitempty" validator:"dive,gte=0,lte=1"`
	AlphaMode            AlphaMode            `json:"alphaMode,omitempty" validator:"oneof=OPAQUE MASK BLEND"`
	AlphaCutoff          string               `json:"alphaCutoff"`
	DoubleSided          bool                 `json:"doubleSided,omitempty"`
}

// A NormalTexture references to a normal texture.
type NormalTexture struct {
	TextureInfo
	Scale float32 `json:"scale"`
}

// An OcclusionTexture references to an occlusion texture
type OcclusionTexture struct {
	TextureInfo
	Strength float32 `json:"strength" validator:"gte=0,lte=1"`
}

// A set of parameter values that are used to define the metallic-roughness material model from Physically-Based Rendering (PBR) methodology.
type PBRMetallicRoughness struct {
	Extensible
	BaseColorFactor          [4]float32  `json:"baseColorFactor" validator:"dive,gte=0,lte=1"`
	BaseColorTexture         TextureInfo `json:"baseColorTexture,omitempty"`
	MetallicFactor           float32     `json:"metallicFactor" validator:"gte=0,lte=1"`
	RoughnessFactor          float32     `json:"roughnessFactor" validator:"gte=0,lte=1"`
	MetallicRoughnessTexture TextureInfo `json:"metallicRoughnessTexture,omitempty"`
}

// TextureInfo references to a texture.
type TextureInfo struct {
	Extensible
	Index    uint32 `json:"index"`
	TexCoord uint32 `json:"texCoord,omitempty"` // The index of texture's TEXCOORD attribute used for texture coordinate mapping.
}

// A texture and its sampler.
type Texture struct {
	Named
	Extensible
	Sampler uint32 `json:"sampler,omitempty"`
	Source  uint32 `json:"source,omitempty"`
}

// Sampler of a texture for filtering and wrapping modes.
type Sampler struct {
	Named
	Extensible
	MagFilter MagFilter    `json:"magFilter" validator:"omitempty,oneof=9728 9729"`
	MinFilter MinFilter    `json:"minFilter" validator:"omitempty,oneof=9728 9729 9984 9985 9986 9987"`
	WrapS     WrappingMode `json:"wrapS" validator:"omitempty,oneof=33071 33648 10497"`
	WrapT     WrappingMode `json:"wrapT" validator:"omitempty,oneof=33071 33648 10497"`
}

// Image data used to create a texture. Image can be referenced by URI or bufferView index.
// mimeType is required in the latter case.
type Image struct {
	Named
	Extensible
	URI        string `json:"uri,omitempty" validator:"omitempty,uri|datauri"`
	MimeType   string `json:"mimeType,omitempty" validator:"omitempty,oneof=image/jpeg image/png"`
	BufferView uint32 `json:"bufferView,omitempty"` // Use this instead of the image's uri property.
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
	Named
	Extensible
	Channels []Channel        `json:"channel" validator:"required"`
	Samplers AnimationSampler `json:"sampler" validator:"required"`
}

// AnimationSampler combines input and output accessors with an interpolation algorithm to define a keyframe graph (but not its target).
type AnimationSampler struct {
	Extensible
	Input         uint32        `json:"input"` // The index of an accessor containing keyframe input values.
	Interpolation Interpolation `json:"interpolation,omitempty" validator:"omitempty,oneof=LINEAR STEP CUBICSPLINE"`
	Output        uint32        `json:"output"` // The index of an accessor containing keyframe output values.
}

// The channel targets an animation's sampler at a node's property.
type Channel struct {
	Extensible
	Sampler uint32        `json:"sampler"`
	Target  ChannelTarget `json:"target"`
}

// ChannelTarget describes the index of the node and TRS property that an animation channel targets.
// The Path represents the name of the node's TRS property to modify, or the "weights" of the Morph Targets it instantiates.
// For the "translation" property, the values that are provided by the sampler are the translation along the x, y, and z axes.
// For the "rotation" property, the values are a quaternion in the order (x, y, z, w), where w is the scalar.
// For the "scale" property, the values are the scaling factors along the x, y, and z axes.
type ChannelTarget struct {
	Extensible
	Node uint32 `json:"node,omitempty"`
	Path string `json:"path" validator:"oneof=translation rotation scale weights"`
}
