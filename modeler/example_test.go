package modeler_test

import (
	"github.com/qmuntal/gltf"
	"github.com/qmuntal/gltf/modeler"
)

func Example() {
	doc := gltf.NewDocument()
	positionAccessor := modeler.WritePosition(doc, [][3]float32{{43, 43, 0}, {83, 43, 0}, {63, 63, 40}, {43, 83, 0}, {83, 83, 0}})
	indicesAccessor := modeler.WriteIndices(doc, []uint8{0, 1, 2, 3, 1, 0, 0, 2, 3, 1, 4, 2, 4, 3, 2, 4, 1, 3})
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
