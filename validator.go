package gltf

import (
	val "github.com/go-playground/validator"
)

// Validate ensures that a document follows the glTF 2.0 specs.
func (d *Document) Validate() error {
	validate := val.New()
	validate.RegisterStructValidation(imageValidation, Image{})
	return validate.Struct(d)
}

func imageValidation(sl val.StructLevel) {
	image := sl.Current().Interface().(Image)

	if image.URI == "" && image.MimeType == "" {
		sl.ReportError(image.MimeType, "MimeType", "mimeType", "", "")
	}
}
