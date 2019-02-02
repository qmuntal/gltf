package gltf

// An Asset is metadata about the glTF asset.
type Asset struct {
	Extensible
	Copyright  string `json:"copyright,omitempty"`  // A copyright message suitable for display to credit the content creator.
	Generator  string `json:"generator,omitempty"`  // Tool that generated this glTF model. Useful for debugging.
	Version    string `json:"version"`              // The glTF version that this asset targets.
	MinVersion string `json:"minVersion,omitempty"` // The minimum glTF version that this asset targets.
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

// The root object for a glTF asset.
type GLTF struct {
	Extensible
	ExtensionsUsed     []string     `json:"extensionsUsed,omitempty"`     // Names of glTF extensions used somewhere in this asset.
	ExtensionsRequired []string     `json:"extensionsRequired,omitempty"` // Names of glTF extensions required to properly load this asset.
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
