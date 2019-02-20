# gltf
[![Documentation](https://godoc.org/github.com/qmuntal/gltf?status.svg)](https://godoc.org/github.com/qmuntal/gltf)
[![Build Status](https://travis-ci.com/qmuntal/gltf.svg?branch=master)](https://travis-ci.com/qmuntal/gltf)
[![Go Report Card](https://goreportcard.com/badge/github.com/qmuntal/gltf)](https://goreportcard.com/report/github.com/qmuntal/gltf)
[![codecov](https://coveralls.io/repos/github/qmuntal/gltf/badge.svg)](https://coveralls.io/github/qmuntal/gltf?branch=master)
[![codeclimate](https://codeclimate.com/github/qmuntal/gltf/badges/gpa.svg)](https://codeclimate.com/github/qmuntal/gltf)
[![License](https://img.shields.io/badge/License-BSD%202--Clause-orange.svg)](https://opensource.org/licenses/BSD-2-Clause)

A Go package for simple, efficient, and robust serialization/deserialization of glTF 2.0 

## Features
* Moderate parsing time and memory consumption.
* glTF specification v2.0.0
  * [x] ASCII glTF.
  * [x] Binary glTF(GLB).
  * [x] PBR material description.
* glTF validaton
  * [x] Validate against schemas.
  * [ ] Validate coherence.
* Buffers
  * [x] Parse BASE64 encoded embedded buffer data(DataURI).
  * [x] Load .bin file.
* Read from io.Reader
  * [x] Boilerplate for disk loading.
  * [x] Custom callback handlers.
  * [x] Automatic ASCII / glTF detection.
* Write to io.Writer
  * [x] Boilerplate for disk saving.
  * [x] Custom callback handlers.
  * [x] ASCII / Binary
* Extensions
  * [ ] KHR_draco_mesh_compression
  * [ ] KHR_lights_punctual
  * [ ] KHR_materials_pbrSpecularGlossiness
  * [ ] KHR_materials_unlit
  * [ ] KHR_techniques_webgl
  * [ ] KHR_texture_transform
  
## Examples
### Read
```go
doc, err := gltf.Open("./a.gltf")
if err != nil {
  panic(err)
}
fmt.Print(doc.Asset)
```
### Write
```go
doc := &gltf.Document{
  Scene: 0, 
  Asset: gltf.Asset{Generator: "qmuntal/gltf"}, 
  Scenes: []gltf.Scene{{Extensions: 8.0, Extras: map[string]interface{}{"a": "b"}, Name: "s_1"}}
}
 
if err := gltf.Save(doc, "./a.gltf", true); err != nil {
  panic(err)
}
```
