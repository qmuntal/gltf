package main

import (
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/qmuntal/gltf"
	"github.com/qmuntal/gltf/modeler"
)

func TestLoadMesh(t *testing.T) {
	tests := []struct {
		path      string
		positions int
		indices   int
	}{
		{"BoxVertexColors/glTF/BoxVertexColors.gltf", 24, 36},
		{"BoxVertexColors/glTF-Binary/BoxVertexColors.glb", 24, 36},
		{"BoxVertexColors/glTF-Embedded/BoxVertexColors.gltf", 24, 36},
		{"TriangleWithoutIndices/glTF/TriangleWithoutIndices.gltf", 3, 3},
	}
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			m, err := loadMesh(filepath.Join("..", "..", "testdata", filepath.FromSlash(tt.path)))
			if err != nil {
				t.Fatal(err)
			}
			if len(m.Positions) != tt.positions || len(m.Indices) != tt.indices {
				t.Fatalf("got %d positions and %d indices, want %d and %d", len(m.Positions), len(m.Indices), tt.positions, tt.indices)
			}
			if tt.positions == 3 {
				want := [][3]float32{{0, 0, 0}, {1, 0, 0}, {0, 1, 0}}
				if !reflect.DeepEqual(m.Positions, want) || !reflect.DeepEqual(m.Indices, []uint32{0, 1, 2}) {
					t.Fatalf("unexpected triangle data: %+v", m)
				}
			}
		})
	}
}

func TestLoadMeshErrors(t *testing.T) {
	tests := []struct {
		name   string
		change func(*gltf.Document)
		want   string
	}{
		{"no mesh", func(d *gltf.Document) { d.Meshes = nil }, "no mesh primitive"},
		{"no primitive", func(d *gltf.Document) { d.Meshes[0].Primitives = nil }, "mesh.primitives is required"},
		{"lines", func(d *gltf.Document) { d.Meshes[0].Primitives[0].Mode = gltf.PrimitiveLines }, "only triangle"},
		{"no position", func(d *gltf.Document) {
			d.Meshes[0].Primitives[0].Attributes = gltf.PrimitiveAttributes{gltf.NORMAL: 0}
		}, "no POSITION"},
		{"bad position accessor", func(d *gltf.Document) { d.Meshes[0].Primitives[0].Attributes[gltf.POSITION] = 99 }, "invalid accessor"},
		{"bad index accessor", func(d *gltf.Document) { d.Meshes[0].Primitives[0].Indices = gltf.Index(99) }, "invalid accessor"},
		{"wrong position type", func(d *gltf.Document) { d.Accessors[0].Type = gltf.AccessorScalar }, "read positions"},
		{"wrong index type", func(d *gltf.Document) { d.Accessors[1].Type = gltf.AccessorVec3 }, "read indices"},
		{"empty positions", func(d *gltf.Document) { d.Accessors[0].Count = 0 }, "accessor.count must be greater than zero"},
		{"incomplete triangle", func(d *gltf.Document) { d.Accessors[1].Count = 2 }, "multiple of three"},
		{"out of range index", func(d *gltf.Document) { d.Accessors[0].Count = 2 }, "out of range"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := gltf.NewDocument()
			position := modeler.WritePosition(doc, [][3]float32{{0, 0, 0}, {1, 0, 0}, {0, 1, 0}})
			indices := modeler.WriteIndices(doc, []uint16{0, 1, 2})
			doc.Meshes = []*gltf.Mesh{{Primitives: []*gltf.Primitive{{
				Attributes: gltf.PrimitiveAttributes{gltf.POSITION: position},
				Indices:    gltf.Index(indices),
			}}}}
			tt.change(doc)
			path := filepath.Join(t.TempDir(), "mesh.glb")
			if err := gltf.SaveBinary(doc, path); err != nil {
				t.Fatal(err)
			}
			if _, err := loadMesh(path); err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("got error %v, want %q", err, tt.want)
			}
		})
	}
	t.Run("missing file", func(t *testing.T) {
		if _, err := loadMesh(filepath.Join(t.TempDir(), "missing.gltf")); err == nil {
			t.Fatal("expected an error for a missing file")
		}
	})
}
