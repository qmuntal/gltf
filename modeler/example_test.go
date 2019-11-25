package modeler_test

import (
	"github.com/qmuntal/gltf"
	"github.com/qmuntal/gltf/modeler"
)

func Example() {
	m := modeler.NewModeler()
	positionAccessor := m.AddPosition(0, [][3]float32{{43, 43, 0}, {83, 43, 0}, {63, 63, 40}, {43, 83, 0}, {83, 83, 0}})
	indicesAccessor := m.AddIndices(0, []uint8{0, 1, 2, 3, 1, 0, 0, 2, 3, 1, 4, 2, 4, 3, 2, 4, 1, 3})
	colorIndices := m.AddColor(0, [][3]uint8{{50, 155, 255}, {0, 100, 200}, {255, 155, 50}, {155, 155, 155}, {25, 25, 25}})
	m.Document.Meshes = []gltf.Mesh{{
		Name: "Pyramid",
		Primitives: []gltf.Primitive{
			{
				Indices: gltf.Index(indicesAccessor),
				Attributes: map[string]uint32{
					"POSITION": positionAccessor,
					"COLOR_0":  colorIndices,
				},
			},
		},
	}}
	m.Nodes = []gltf.Node{{Name: "Root", Mesh: gltf.Index(0)}}
	m.Scenes[0].Nodes = append(m.Scenes[0].Nodes, 0)
	if err := gltf.SaveBinary(m.Document, "./a.gltf"); err != nil {
		panic(err)
	}
	// Output:
	//
}
