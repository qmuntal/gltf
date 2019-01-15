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
	PBRMetallicRoughness PBRMetallicRoughness `json:"pbrMetallicRoughness,omitempty"` //Metallic-roughness material model parameters.
	NormalTexture        NormalTexture        `json:"normalTexture,omitempty"`        // A tangent space normal map.
	OcclusionTexture     OcclusionTexture     `json:"occlusionTexture,omitempty"`     // The occlusion map texture.
	EmissiveTexture      TextureInfo          `json:"emissiveTexture,omitempty"`      // The emissive map controls the color and intensity of the light being emitted by the material.
	EmissiveFactor       [3]float32           `json:"emissiveFactor,omitempty"`       // The RGB components of the emissive color of the material. These values are linear
	AlphaMode            AlphaMode            `json:"alphaMode,omitempty"`            // The material's alpha rendering mode enumeration specifying the interpretation of the alpha value of the main factor and texture.
	AlphaCutoff          string               `json:"alphaCutoff"`                    // Specifies the cutoff threshold when in MASK mode.
	DoubleSided          bool                 `json:"doubleSided,omitempty"`          // Specifies whether the material is double sided.
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

// TextureSampler properties for filtering and wrapping modes.
type TextureSampler struct {
	Named
	Extensible
	MagFilter MagFilter    `json:"magFilter"` // Magnification filter.
	MinFilter MinFilter    `json:"minFilter"` // Minification filter.
	WrapS     WrappingMode `json:"wrapS"`     // S wrapping mode.
	WrapT     WrappingMode `json:"wrapT"`     // T wrapping mode.
}
