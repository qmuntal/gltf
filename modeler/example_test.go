package modeler_test

import (
	"bytes"
	"io/ioutil"

	"github.com/qmuntal/gltf"
	"github.com/qmuntal/gltf/modeler"
)

func Example() {
	doc := gltf.NewDocument()
	positionAccessor := modeler.WritePosition(doc, [][3]float32{{43, 43, 0}, {83, 43, 0}, {63, 63, 40}, {43, 83, 0}, {83, 83, 0}})
	indicesAccessor := modeler.WriteIndices(doc, []uint16{0, 1, 2, 3, 1, 0, 0, 2, 3, 1, 4, 2, 4, 3, 2, 4, 1, 3})
	colorIndices := modeler.WriteColor(doc, [][3]uint8{{50, 155, 255}, {0, 100, 200}, {255, 155, 50}, {155, 155, 155}, {25, 25, 25}})
	doc.Meshes = []*gltf.Mesh{{
		Name: "Pyramid",
		Primitives: []*gltf.Primitive{
			{
				Indices: gltf.Index(indicesAccessor),
				Attributes: map[string]uint32{
					gltf.POSITION: positionAccessor,
					gltf.COLOR_0:  colorIndices,
				},
			},
		},
	}}
	doc.Nodes = []*gltf.Node{{Name: "Root", Mesh: gltf.Index(0)}}
	doc.Scenes[0].Nodes = append(doc.Scenes[0].Nodes, 0)
	if err := gltf.SaveBinary(doc, "./example.glb"); err != nil {
		panic(err)
	}
}

func ExampleWriteImage() {
	img, err := ioutil.ReadFile("../assets/gopher_high.png")
	if err != nil {
		panic(err)
	}
	doc := gltf.NewDocument()
	indicesAccessor := modeler.WriteIndices(doc, []uint16{0, 1, 2, 3, 1, 0, 0, 2, 3, 1, 4, 2, 4, 3, 2, 4, 1, 3})
	positionAccessor := modeler.WritePosition(doc, [][3]float32{{43, 43, 0}, {83, 43, 0}, {63, 63, 40}, {43, 83, 0}, {83, 83, 0}})
	textureAccessor := modeler.WriteTextureCoord(doc, [][2]float32{{0, 1}, {0.4, 1}, {0.4, 0}, {0.4, 1}, {0, 1}})
	imageIdx, err := modeler.WriteImage(doc, "gopher", "image/png", bytes.NewReader(img))
	if err != nil {
		panic(err)
	}
	doc.Textures = append(doc.Textures, &gltf.Texture{Source: gltf.Index(imageIdx)})
	doc.Materials = append(doc.Materials, &gltf.Material{
		Name: "Texture",
		PBRMetallicRoughness: &gltf.PBRMetallicRoughness{
			BaseColorTexture: &gltf.TextureInfo{
				Index: uint32(len(doc.Textures) - 1),
			},
			MetallicFactor: gltf.Float(0),
		},
	})
	doc.Meshes = []*gltf.Mesh{{
		Name: "Pyramid",
		Primitives: []*gltf.Primitive{
			{
				Indices: gltf.Index(indicesAccessor),
				Attributes: map[string]uint32{
					gltf.POSITION:   positionAccessor,
					gltf.TEXCOORD_0: textureAccessor,
				},
				Material: gltf.Index(uint32(len(doc.Materials) - 1)),
			},
		},
	}}
	doc.Nodes = []*gltf.Node{{Name: "Root", Mesh: gltf.Index(0)}}
	doc.Scenes[0].Nodes = append(doc.Scenes[0].Nodes, 0)
	if err := gltf.SaveBinary(doc, "./example.glb"); err != nil {
		panic(err)
	}
}

func ExampleWriteAccessorsInterleaved() {
	doc := gltf.NewDocument()
	indices, _ := modeler.WriteAccessorsInterleaved(doc,
		[][3]float32{{43, 43, 0}, {83, 43, 0}, {63, 63, 40}, {43, 83, 0}, {83, 83, 0}},
		[][3]uint8{{50, 155, 255}, {0, 100, 200}, {255, 155, 50}, {155, 155, 155}, {25, 25, 25}},
	)
	indicesAccessor := modeler.WriteIndices(doc, []uint16{0, 1, 2, 3, 1, 0, 0, 2, 3, 1, 4, 2, 4, 3, 2, 4, 1, 3})
	doc.Meshes = []*gltf.Mesh{{
		Name: "Pyramid",
		Primitives: []*gltf.Primitive{
			{
				Indices: gltf.Index(indicesAccessor),
				Attributes: map[string]uint32{
					gltf.POSITION: indices[0],
					gltf.COLOR_0:  indices[1],
				},
			},
		},
	}}
	doc.Nodes = []*gltf.Node{{Name: "Root", Mesh: gltf.Index(0)}}
	doc.Scenes[0].Nodes = append(doc.Scenes[0].Nodes, 0)
	if err := gltf.SaveBinary(doc, "./example.glb"); err != nil {
		panic(err)
	}
}

func ExampleWriteAttributesInterleaved() {
	doc := gltf.NewDocument()
	attrs, _ := modeler.WriteAttributesInterleaved(doc,
		modeler.Attribute{Name: gltf.POSITION, Data: [][3]float32{{1, 2, 3}, {0, 0, -1}}},
		modeler.Attribute{Name: gltf.NORMAL, Data: [][3]float32{{1, 2, 3}, {0, 0, -1}}},
		modeler.Attribute{Name: gltf.TANGENT, Data: [][4]float32{{1, 2, 3, 4}, {1, 2, 3, 4}}},
		modeler.Attribute{Name: gltf.TEXCOORD_0, Data: [][2]uint8{{0, 255}, {255, 0}}},
		modeler.Attribute{Name: gltf.TEXCOORD_1, Data: [][2]float32{{1, 2}, {1, 2}}},
		modeler.Attribute{Name: gltf.JOINTS_0, Data: [][4]uint8{{1, 2, 3, 4}, {1, 2, 3, 4}}},
		modeler.Attribute{Name: gltf.WEIGHTS_0, Data: [][4]uint8{{1, 2, 3, 4}, {1, 2, 3, 4}}},
		modeler.Attribute{Name: gltf.COLOR_0, Data: [][3]uint8{{255, 255, 255}, {0, 255, 0}}},
		modeler.Attribute{Name: "COLOR_1", Data: [][3]uint8{{0, 0, 255}, {100, 200, 0}}},
		modeler.Attribute{Name: "COLOR_2", Data: [][4]uint8{{23, 58, 188, 1}, {0, 155, 0, 0}}},
	)
	indicesAccessor := modeler.WriteIndices(doc, []uint16{0, 1, 2, 3, 1, 0, 0, 2, 3, 1, 4, 2, 4, 3, 2, 4, 1, 3})
	doc.Meshes = []*gltf.Mesh{{
		Name: "Pyramid",
		Primitives: []*gltf.Primitive{
			{
				Indices:    gltf.Index(indicesAccessor),
				Attributes: attrs,
			},
		},
	}}
	doc.Nodes = []*gltf.Node{{Name: "Root", Mesh: gltf.Index(0)}}
	doc.Scenes[0].Nodes = append(doc.Scenes[0].Nodes, 0)
	if err := gltf.SaveBinary(doc, "./example.glb"); err != nil {
		panic(err)
	}
}
