package gltf

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
	MagFilter MagFilter    `json:"magFilter"`
	MinFilter MinFilter    `json:"minFilter"`
	WrapS     WrappingMode `json:"wrapS"`
	WrapT     WrappingMode `json:"wrapT"`
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
