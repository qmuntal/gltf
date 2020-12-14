package modeler_test

import (
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
					"POSITION": positionAccessor,
					"COLOR_0":  colorIndices,
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
					"POSITION": indices[0],
					"COLOR_0":  indices[1],
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
	attrs, _ := modeler.WriteAttributesInterleaved(doc, modeler.Attributes{
		Position:       [][3]float32{{1, 2, 3}, {0, 0, -1}},
		Normal:         [][3]float32{{1, 2, 3}, {0, 0, -1}},
		Tangent:        [][4]float32{{1, 2, 3, 4}, {1, 2, 3, 4}},
		TextureCoord_0: [][2]uint8{{0, 255}, {255, 0}},
		TextureCoord_1: [][2]float32{{1, 2}, {1, 2}},
		Joints:         [][4]uint8{{1, 2, 3, 4}, {1, 2, 3, 4}},
		Weights:        [][4]uint8{{1, 2, 3, 4}, {1, 2, 3, 4}},
		Color:          [][3]uint8{{255, 255, 255}, {0, 255, 0}},
		CustomAttributes: []modeler.CustomAttribute{
			{Name: "COLOR_1", Data: [][3]uint8{{0, 0, 255}, {100, 200, 0}}},
			{Name: "COLOR_2", Data: [][4]uint8{{23, 58, 188, 1}, {0, 155, 0, 0}}},
		},
	})
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
