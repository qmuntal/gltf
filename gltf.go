package gltf

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

// Interpolation algorithm.
type Interpolation string

const (
	Linear      Interpolation = "LINEAR"
	Step                      = "STEP"
	CubicSpline               = "CUBICSPLINE"
)

// CameraType specifies if the camera uses a perspective or orthographic projection.
// Based on this, either the camera's perspective or orthographic property will be defined.
type CameraType string

const (
	PerspectiveType CameraType = "perspective"
	OrtographicType            = "orthographic"
)

// The AlphMode enumeration specifying the interpretation of the alpha value of the main factor and texture.
type AlphaMode string

const (
	Opaque AlphaMode = "OPAQUE"
	Mask             = "MASK"
	Blend            = "BLEND"
)

// An accessor is a typed view into a bufferView. A bufferView contains raw binary data.
// An accessor provides a typed view into a bufferView or a subset of a bufferView
// similar to how WebGL's vertexAttribPointer() defines an attribute in a buffer.
type Accessor struct {
	Named
	Extensible
	BufferView    uint32        `json:"bufferView,omitempty"` // The index of the bufferView.
	ByteOffset    uint32        `json:"byteOffset,omitempty"` // The offset relative to the start of the bufferView in bytes.
	ComponentType ComponentType `json:"componentType"`        // The datatype of components in the attribute.
	Normalized    bool          `json:"normalized,omitempty"` // Specifies whether integer data values should be normalized.
	Count         uint32        `json:"count"`                // The number of attributes referenced by this accessor.
	Type          string        `json:"type"`                 // Specifies if the attribute is a scalar, vector, or matrix.
	Max           []float32     `json:"max,omitempty"`        // Maximum value of each component in this attribute.
	Min           []float32     `json:"min,omitempty"`        // Minimum value of each component in this attribute.
	Sparce        interface{}   `json:"sparce,omitempty"`     // Sparse storage of attributes that deviate from their initialization value.
}

// An Animation keyframe.
type Animation struct {
	Named
	Extensible
	Channels []Channel `json:"channel"` // An array of channels, each of which targets an animation's sampler at a node's property. Different channels of the same animation can't have equal targets.
}

type AnimationSampler struct {
	Extensible
	Input         uint32        `json:"input"`                   // The index of an accessor containing keyframe input values.
	Interpolation Interpolation `json:"interpolation,omitempty"` // Interpolation algorithm.
	Output        uint32        `json:"output"`                  // The index of an accessor containing keyframe output values.
}

// An Asset is metadata about the glTF asset.
type Asset struct {
	Extensible
	Copyright  string `json:"copyright,omitempty"`  // A copyright message suitable for display to credit the content creator.
	Generator  string `json:"generator,omitempty"`  // Tool that generated this glTF model. Useful for debugging.
	Version    string `json:"version"`              // The glTF version that this asset targets.
	MinVersion string `json:"minVersion,omitempty"` // The minimum glTF version that this asset targets.
}

// A buffer points to binary geometry, animation, or skins.
type Buffer struct {
	Named
	Extensible
	URI        string `json:"uri,omitempty"` // The uri of the buffer.
	ByteLength uint32 `json:"byteLength"`    // The total byte length of the buffer view.
}

// BufferView is a view into a buffer generally representing a subset of the buffer.
type BufferView struct {
	Extensible
	Buffer     uint32 `json:"buffer"`               // The index of the buffer.
	ByteOffset uint32 `json:"byteOffset,omitempty"` // The offset into the buffer in bytes.
	ByteLength uint32 `json:"byteLength"`           // The length of the bufferView in bytes.
	ByteStride uint32 `json:"byteStride,omitempty"` // The stride, in bytes.
	Target     uint32 `json:"target,omitempty"`     // The target that the GPU buffer should be bound to.
}

// A Camera's projection. A node can reference a camera to apply a transform to place the camera in the scene.
type Camera struct {
	Named
	Extensible
	Ortographic Ortographic `json:"orthographic,omitempty"` // An orthographic camera containing properties to create an orthographic projection matrix.
	Perspective Perspective `json:"perspective,omitempty"`  // A perspective camera containing properties to create a perspective projection matrix.
	Type        CameraType  `json:"type"`                   // Specifies if the camera uses a perspective or orthographic projection.
}

// The channel targets an animation's sampler at a node's property.
type Channel struct {
	Extensible
	Sampler uint32        `json:"sampler"` // The index of a sampler in this animation used to compute the value for the target.
	Target  ChannelTarget `json:"target"`  // The index of the node and TRS property to target.
}

// ChannelTarget describes the index of the node and TRS property that an animation channel targets.
// The Path represents the name of the node's TRS property to modify, or the \"weights\" of the Morph Targets it instantiates.
// For the \"translation\" property, the values that are provided by the sampler are the translation along the x, y, and z axes.
// For the \"rotation\" property, the values are a quaternion in the order (x, y, z, w), where w is the scalar.
// For the \"scale\" property, the values are the scaling factors along the x, y, and z axes.
type ChannelTarget struct {
	Node uint32 `json:"node,omitempty"` // The index of the node to target.
	Path string `json:"path"`           // TRS property.
}

// Image data used to create a texture. Image can be referenced by URI or bufferView index.
// mimeType is required in the latter case.
type Image struct {
	Named
	Extensible
	URI        string `json:"uri,omitempty"`        // The uri of the image.
	MimeType   string `json:"mimeType,omitempty"`   // The image's MIME type.
	BufferView uint32 `json:"bufferView,omitempty"` // The index of the bufferView that contains the image. Use this instead of the image's uri property.
}

// Indices of those attributes that deviate from their initialization value.
type Indices struct {
	Extensible
	BufferView    uint32        `json:"bufferView"`           // The index of the bufferView with sparse indices.
	ByteOffset    uint32        `json:"byteOffset,omitempty"` // The offset relative to the start of the bufferView in bytes. Must be aligned.
	ComponentType ComponentType `json:"componentType"`        // The indices data type. Valid values correspond to WebGL enums: 5121 (UNSIGNED_BYTE), 5123 (UNSIGNED_SHORT), 5125 (UNSIGNED_INT)
}

// The material appearance of a primitive.
type Material struct {
	Named
	Extensible
	PBRMetallicRoughness PBRMetallicRoughness `json:"pbrMetallicRoughness,omitempty"` //Metallic-roughness material model parameters.
	NormalTexture        NormalTexture        `json:"normalTexture,omitempty"`        // A tangent space normal map.
	OcclusionTexture     OcclusionTexture     `json:"occlusionTexture,omitempty"`     // The occlusion map texture.
	EmissiveTexture      TextureInfo          `json:"emissiveTexture,omitempty"`      // The emissive map controls the color and intensity of the light being emitted by the material.
	EmissiveFactor       [3]float32           `json:"emissiveFactor,omitempty"`       // The RGB components of the emissive color of the material. These values are linear
	AlphaMode            AlphaMode            `json:"alphaMode,omitempty"`            // The material's alpha rendering mode enumeration specifying the interpretation of the alpha value of the main factor and texture.
	AlphaCutoff          string               `json:"alphaCutoff"`                    // Specifies the cutoff threshold when in MASK mode.
	DoubleSided          bool                 `json:"doubleSided,omitempty"`          // Specifies whether the material is double sided.
}

// A Mesh is a set of primitives to be rendered. A node can contain one mesh. A node's transform places the mesh in the scene.
type Mesh struct {
	Named
	Extensible
	Primitives []interface{} `json:"primitives"`        // An array of primitives, each defining geometry to be rendered with a material.
	Weights    []float32     `json:"weights,omitempty"` // Array of weights to be applied to the Morph Targets.
}

// A NormalTexture references to a normal texture.
type NormalTexture struct {
	TextureInfo
	Scale float32 `json:"scale"` // The scalar multiplier applied to each normal vector of the normal texture.
}

// An OcclusionTexture references to an occlusion texture
type OcclusionTexture struct {
	TextureInfo
	Strength float32 `json:"strength"` // A scalar multiplier controlling the amount of occlusion applied.
}

// An Orthographic camera containing properties to create an orthographic projection matrix.
type Ortographic struct {
	Extensible
	Xmag  float32 `json:"xmag"`  // The horizontal magnification of the view.
	Ymag  float32 `json:"ymag"`  // The vertical magnification of the view.
	Zfar  float32 `json:"zfar"`  // The distance to the far clipping plane. zfar must be greater than znear.
	Znear float32 `json:"znear"` // The distance to the near clipping plane.
}

// A perspective camera containing properties to create a perspective projection matrix.
type Perspective struct {
	Extensible
	AspectRatio float32 `json:"aspectRatio,omitempty"` // The aspect ratio of the field of view.
	Yfov        float32 `json:"yfov"`                  // The vertical field of view in radians.
	Zfar        float32 `json:"zfar,omitempty"`        // The distance to the far clipping plane.
	Znear       float32 `json:"znear"`                 // The distance to the near clipping plane.
}

// A set of parameter values that are used to define the metallic-roughness material model from Physically-Based Rendering (PBR) methodology.
type PBRMetallicRoughness struct {
	Extensible
	BaseColorFactor          [4]float32  `json:"baseColorFactor"`
	BaseColorTexture         TextureInfo `json:"baseColorTexture,omitempty"`
	MetallicFactor           float32     `json:"metallicFactor"`
	RoughnessFactor          float32     `json:"roughnessFactor"`
	MetallicRoughnessTexture TextureInfo `json:"metallicRoughnessTexture,omitempty"`
}

// TextureInfo references to a texture.
type TextureInfo struct {
	Extensible
	Index    uint32 `json:"index"`              // The index of the texture.
	TexCoord uint32 `json:"texCoord,omitempty"` // The set index of texture's TEXCOORD attribute used for texture coordinate mapping.
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
