<p align="center"><img width="125" src="./assets/gopher.png" alt="Gopher glTF"></p>

[![PkgGoDev](https://pkg.go.dev/badge/qmuntal/gltf)](https://pkg.go.dev/github.com/qmuntal/gltf)
[![Build Status](https://travis-ci.com/qmuntal/gltf.svg?branch=master)](https://travis-ci.com/qmuntal/gltf)
[![Go Report Card](https://goreportcard.com/badge/github.com/qmuntal/gltf)](https://goreportcard.com/report/github.com/qmuntal/gltf)
[![codecov](https://coveralls.io/repos/github/qmuntal/gltf/badge.svg)](https://coveralls.io/github/qmuntal/gltf?branch=master)
[![codeclimate](https://codeclimate.com/github/qmuntal/gltf/badges/gpa.svg)](https://codeclimate.com/github/qmuntal/gltf)
[![License](https://img.shields.io/badge/License-BSD%202--Clause-orange.svg)](https://opensource.org/licenses/BSD-2-Clause)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go)  

# gltf

A Go package for simple, efficient, and robust serialization/deserialization of [glTF 2.0](https://www.khronos.org/gltf/) (GL Transmission Format), a royalty-free specification for the efficient transmission and loading of 3D scenes and models by applications.

## Features

* High parsing speed and moderate memory consumption
* glTF specification v2.0.0
  * [x] ASCII glTF
  * [x] Binary glTF(GLB)
  * [x] PBR material description
  * [x] Modeler package
* glTF validaton
  * [ ] Validate against schemas
  * [ ] Validate coherence
* Buffers
  * [x] Parse BASE64 encoded embedded buffer data(DataURI)
  * [x] Load .bin file
  * [x] Binary package
* Read from io.Reader
  * [x] Boilerplate for disk loading
  * [x] Custom callback handlers
  * [x] Automatic ASCII / glTF detection
* Write to io.Writer
  * [x] Boilerplate for disk saving
  * [x] Custom callback handlers
  * [x] ASCII / Binary
* Extensions
  * [x] KHR_draco_mesh_compression
  * [x] KHR_lights_punctual
  * [x] KHR_materials_pbrSpecularGlossiness
  * [x] KHR_materials_unlit
  * [ ] KHR_techniques_webgl
  * [x] KHR_texture_transform
  * [x] Support custom extensions

## Extensions

This module is designed to support dynamic extensions. By default only the core specification is decoded and the data inside the extensions objects are stored as `json.RawMessage` so they can be decoded outside this package or automatically encoded when saving the document.

To decode one of the supported extensions the only required action is to import the associated package, this way the extension will not be stored as `json.RawMessage` but as the type defined in the extension package:

```go
import (
  "github.com/qmuntal/gltf"
  "github.com/qmuntal/gltf/ext/lightspuntual"
)

func ExampleExension() {
  doc, _ := gltf.Open("...")
  if v, ok := doc.Extensions[lightspuntual.ExtensionName]; ok {
      for _, l := range v.(lightspuntual.Lights) {
          fmt.Print(l.Type)
      }
  }
}
```

It is not necessary to call `gltf.RegisterExtension` for built-in extensions, as these auto-register themselves on `init()`.

## Perfomance

All the functionality is benchmarked and tested using the official [glTF Samples](https://github.com/KhronosGroup/glTF-Sample-Models) in the utility package [qmuntal/gltf-bench](https://github.com/qmuntal/gltf-bench/).
The results show that the perfomance of this package is equivalent to [fx-gltf](https://github.com/jessey-git/fx-gltf), a reference perfomance-driven glTF implementation for C++, .

## Examples

### Read

```go
doc, err := gltf.Open("./a.gltf")
if err != nil {
    panic(err)
}
fmt.Print(doc.Asset)
```

### Create a glb using gltf/modeler

The following example generates a single triangle with colors per vertex.

![screenshot](./assets/color-triangle.png)

```go
package main

import (
    "github.com/qmuntal/gltf"
    "github.com/qmuntal/gltf/modeler"
)

func main() {
    doc := gltf.NewDocument()
    positionAccessor := modeler.WritePosition(doc, [][3]float32{{0, 0, 0}, {0, 10, 0}, {0, 0, 10}})
    indicesAccessor := modeler.WriteIndices(doc, []uint8{0, 1, 2})
    colorAccessor := modeler.WriteColor(doc, [][3]uint8{{255, 0, 0}, {0, 255, 0}, {0, 0, 255}})
    doc.Meshes = []*gltf.Mesh{{
        Name: "Pyramid",
        Primitives: []*gltf.Primitive{{
            Indices: gltf.Index(indicesAccessor),
            Attributes: map[string]uint32{
              "POSITION": positionAccessor,
              "COLOR_0":  colorAccessor,
            },
        }},
    }}
    doc.Nodes = []*gltf.Node{{Name: "Root", Mesh: gltf.Index(0)}}
    doc.Scenes[0].Nodes = append(doc.Scenes[0].Nodes, 0)
    if err := gltf.SaveBinary(doc, "./example.glb"); err != nil {
        panic(err)
    }
}
```

### Create a glb using raw data

The following example generates a 3D box with colors per vertex.

![screenshot](./assets/color-cube.png)

```go

package main

import "github.com/qmuntal/gltf"

func main() {
    doc := &gltf.Document{
        Accessors: []*gltf.Accessor{
            {BufferView: gltf.Index(0), ComponentType: gltf.ComponentUshort, Count: 36, Type: gltf.AccessorScalar},
            {BufferView: gltf.Index(1), ComponentType: gltf.ComponentFloat, Count: 24, Max: []float32{0.5, 0.5, 0.5}, Min: []float32{-0.5, -0.5, -0.5}, Type: gltf.AccessorVec3},
            {BufferView: gltf.Index(2), ComponentType: gltf.ComponentFloat, Count: 24, Type: gltf.AccessorVec3},
            {BufferView: gltf.Index(3), ComponentType: gltf.ComponentFloat, Count: 24, Type: gltf.AccessorVec4},
        },
        Asset: gltf.Asset{Version: "2.0", Generator: "FBX2glTF"},
        BufferViews: []*gltf.BufferView{
            {Buffer: 0, ByteLength: 72, ByteOffset: 0, Target: gltf.TargetElementArrayBuffer},
            {Buffer: 0, ByteLength: 288, ByteOffset: 72, Target: gltf.TargetArrayBuffer},
            {Buffer: 0, ByteLength: 288, ByteOffset: 360, Target: gltf.TargetArrayBuffer},
            {Buffer: 0, ByteLength: 384, ByteOffset: 648, Target: gltf.TargetArrayBuffer},
        },
        Buffers: []*gltf.Buffer{{ByteLength: 1224, URI: bufferData}},
        Materials: []*gltf.Material{{
            Name: "Default",
            AlphaMode: gltf.AlphaOpaque,
            AlphaCutoff: gltf.Float(0.5),
            PBRMetallicRoughness: &gltf.PBRMetallicRoughness{
              BaseColorFactor: &[4]float32{0.8, 0.8, 0.8, 0.5},
              MetallicFactor: gltf.Float(0.1),
              RoughnessFactor: gltf.Float(0.99),
            },
        }},
        Meshes: []*gltf.Mesh{{
          Name: "Cube",
          Primitives: []*gltf.Primitive{{
            Indices: gltf.Index(0),
            Material: gltf.Index(0),
            Mode: gltf.PrimitiveTriangles,
            Attributes: map[string]uint32{"POSITION": 1, "COLOR_0": 3, "NORMAL": 2},
          }},
        }},
        Nodes: []*gltf.Node{
            {Name: "RootNode", Children: []uint32{1}},
            {Name: "Cube", Mesh: gltf.Index(0)},
        },
        Scene:    gltf.Index(0),
        Scenes:   []*gltf.Scene{{Name: "Root Scene", Nodes: []uint32{0}}},
    }
    if err := gltf.Save(doc, "./cube.gltf"); err != nil {
        panic(err)
    }
}

const bufferData = "data:application/octet-stream;base64,AAABAAIAAQAAAAMABAAFAAYABwAEAAYACAAJAAoACwAJAAgADAANAA4ADQAMAA8AEAARABIAEAASABMAFAAVABYAFAAWABcAAAAAvwAAAD8AAAC/AAAAPwAAAD8AAAA/AAAAPwAAAD8AAAC/AAAAvwAAAD8AAAA/AAAAvwAAAD8AAAA/AAAAvwAAAD8AAAC/AAAAvwAAAL8AAAC/AAAAvwAAAL8AAAA/AAAAPwAAAD8AAAA/AAAAvwAAAL8AAAA/AAAAPwAAAL8AAAA/AAAAvwAAAD8AAAA/AAAAPwAAAD8AAAC/AAAAPwAAAL8AAAA/AAAAPwAAAL8AAAC/AAAAPwAAAD8AAAA/AAAAPwAAAL8AAAC/AAAAvwAAAL8AAAC/AAAAvwAAAD8AAAC/AAAAPwAAAD8AAAC/AAAAPwAAAL8AAAA/AAAAvwAAAL8AAAA/AAAAvwAAAL8AAAC/AAAAPwAAAL8AAAC/AAAAAAAAgD8AAAAAAAAAAAAAgD8AAAAAAAAAAAAAgD8AAAAAAAAAAAAAgD8AAAAAAACAvwAAAAAAAAAAAACAvwAAAAAAAAAAAACAvwAAAAAAAAAAAACAvwAAAAAAAAAAAAAAAAAAAAAAAIA/AAAAAAAAAAAAAIA/AAAAAAAAAAAAAIA/AAAAAAAAAAAAAIA/AACAPwAAAAAAAAAAAACAPwAAAAAAAAAAAACAPwAAAAAAAAAAAACAPwAAAAAAAAAAAAAAAAAAAAAAAIC/AAAAAAAAAAAAAIC/AAAAAAAAAAAAAIC/AAAAAAAAAAAAAIC/AAAAAAAAgL8AAAAAAAAAAAAAgL8AAAAAAAAAAAAAgL8AAAAAAAAAAAAAgL8AAAAAHekGMHp6JTz+/38/UdInMBHtfz9SNZs6rTHKPJsTwTrvy00/XMpNP0HLTT8AAIA/3ssCP4z/fz+gF+43whYUON7LAj+M/38/oBfuN8IWFDgd6QYwenolPP7/fz9R0icw6vR/PwjDNTqNSsY80xtiOu/LTT9cyk0/QctNPwAAgD8R7X8/UjWbOq0xyjybE8E678tNP1zKTT9By00/AACAP6bmIzurlgY+BNh/P0ziSzveywI/jP9/P6AX7jfCFhQ478tNP1zKTT9By00/AACAP6bmIzurlgY+BNh/P0ziSzsSAI885Oh+PxdguT574rE8Ee1/P1I1mzqtMco8mxPBOhIAjzzk6H4/F2C5PnvisTzq9H8/CMM1Oo1KxjzTG2I6HekGMHp6JTz+/38/UdInMO/LTT9cyk0/QctNPwAAgD+m5iM7q5YGPgTYfz9M4ks778tNP1zKTT9By00/AACAP+r0fz8IwzU6jUrGPNMbYjoSAI885Oh+PxdguT574rE8AACAPgAAAAAAAAA/qqqqPgAAAD8AAAAAAACAPqqqqj4AAIA+qqqqPgAAAACqqqo+AAAAAKqqKj8AAIA+qqoqPwAAAD+qqqo+AACAPqqqKj8AAAA/qqoqPwAAgD6qqqo+AABAP6qqqj4AAAA/qqoqPwAAQD+qqio/AAAAP6qqqj4AAEA/qqoqPwAAgD+qqio/AACAP6qqqj4AAEA/qqqqPgAAAD+qqio/AACAPqqqKj8AAIA+AACAPwAAAD8AAIA/"


```
