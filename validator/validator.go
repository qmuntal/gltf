package validator

import (
	"gltf"

	val "github.com/go-playground/validator"
)

// ValidateDocument ensures that a document follows the glTF 2.0 specs.
func ValidateDocument(doc *gltf.Document) error {
	validate := val.New()
	validate.RegisterStructValidation(imageValidation, gltf.Image{})
	return validate.Struct(doc)
}

func imageValidation(sl val.StructLevel) {
	image := sl.Current().Interface().(gltf.Image)

	if image.URI == "" && image.MimeType == "" {
		sl.ReportError(image.MimeType, "MimeType", "mimeType", "", "")
	}
}
