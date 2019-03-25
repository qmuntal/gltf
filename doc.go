// Package gltf implements a glTF 2.0 file decoder, encoder and validator.
/*
The glTF 2.0 specification is at https://github.com/KhronosGroup/glTF/tree/master/specification/2.0/.

Optional Properties

All optional properties whose default value does not match with the golang type zero value are defines as pointers or as zero arrays.

Guidelines when working with optional values:
It is safe to not define them when saving the glTF if the desired value is the default one.
It is safe to expect that the optional values are not nil when opening a glTF.
When quering values of optional properties that are not indices it is recommended to use the utility functions that returns the property as value
if not nil or the default value if nil.
When assigning values to optional properties one can use the utility functions that take the reference of basic types. Examples:
 gltf.Index(1)
 gltf.Float64(0.5)
*/
package gltf
