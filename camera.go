package gltf

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
	Ortographic Ortographic `json:"orthographic,omitempty"` // An orthographic camera containing properties to create an orthographic projection matrix.
	Perspective Perspective `json:"perspective,omitempty"`  // A perspective camera containing properties to create a perspective projection matrix.
	Type        CameraType  `json:"type"`                   // Specifies if the camera uses a perspective or orthographic projection.
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
