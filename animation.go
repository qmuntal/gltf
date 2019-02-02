package gltf

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
	Channels []Channel        `json:"channel"`
	Samplers AnimationSampler `json:"sampler"`
}

// AnimationSampler combines input and output accessors with an interpolation algorithm to define a keyframe graph (but not its target).
type AnimationSampler struct {
	Extensible
	Input         uint32        `json:"input"` // The index of an accessor containing keyframe input values.
	Interpolation Interpolation `json:"interpolation,omitempty"`
	Output        uint32        `json:"output"` // The index of an accessor containing keyframe output values.
}

// The channel targets an animation's sampler at a node's property.
type Channel struct {
	Extensible
	Sampler uint32        `json:"sampler"`
	Target  ChannelTarget `json:"target"`
}

// ChannelTarget describes the index of the node and TRS property that an animation channel targets.
// The Path represents the name of the node's TRS property to modify, or the \"weights\" of the Morph Targets it instantiates.
// For the \"translation\" property, the values that are provided by the sampler are the translation along the x, y, and z axes.
// For the \"rotation\" property, the values are a quaternion in the order (x, y, z, w), where w is the scalar.
// For the \"scale\" property, the values are the scaling factors along the x, y, and z axes.
type ChannelTarget struct {
	Extensible
	Node uint32 `json:"node,omitempty"`
	Path string `json:"path"`
}
