// Hello World loads a glTF mesh in Go and displays it using WebGL in a browser.
package main

import (
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/qmuntal/gltf"
	"github.com/qmuntal/gltf/modeler"
)

//go:embed index.html
var page []byte

type mesh struct {
	Positions [][3]float32 `json:"positions"`
	Indices   []uint32     `json:"indices"`
}

func main() {
	path := flag.String("model", "testdata/BoxVertexColors/glTF/BoxVertexColors.gltf", "path to a glTF or GLB file")
	flag.Parse()
	m, err := loadMesh(*path)
	if err != nil {
		log.Fatal(err)
	}
	data, err := json.Marshal(m)
	if err != nil {
		log.Fatal(err)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(page)
	})
	mux.HandleFunc("/mesh", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	})
	server := &http.Server{
		Addr:              "127.0.0.1:8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}
	log.Print("Open http://127.0.0.1:8080 in a WebGL 2 browser. Press Ctrl+C to stop.")
	log.Fatal(server.ListenAndServe())
}

func loadMesh(path string) (*mesh, error) {
	doc, err := gltf.Open(path)
	if err != nil {
		return nil, err
	}
	// This example renders only the first primitive, in mesh-local coordinates.
	if len(doc.Meshes) == 0 || doc.Meshes[0] == nil || len(doc.Meshes[0].Primitives) == 0 || doc.Meshes[0].Primitives[0] == nil {
		return nil, fmt.Errorf("model has no mesh primitive")
	}
	primitive := doc.Meshes[0].Primitives[0]
	if primitive.Mode != gltf.PrimitiveTriangles {
		return nil, fmt.Errorf("example supports only triangle primitives")
	}
	position, ok := primitive.Attributes[gltf.POSITION]
	if !ok {
		return nil, fmt.Errorf("primitive has no POSITION attribute")
	}
	accessor := func(index int) (*gltf.Accessor, error) {
		if index < 0 || index >= len(doc.Accessors) || doc.Accessors[index] == nil {
			return nil, fmt.Errorf("invalid accessor index %d", index)
		}
		return doc.Accessors[index], nil
	}
	positions, err := accessor(position)
	if err != nil {
		return nil, err
	}
	m := new(mesh)
	m.Positions, err = modeler.ReadPosition(doc, positions, nil)
	if err != nil {
		return nil, fmt.Errorf("read positions: %w", err)
	}
	if len(m.Positions) == 0 {
		return nil, fmt.Errorf("primitive has no vertices")
	}
	if primitive.Indices != nil {
		indices, err := accessor(*primitive.Indices)
		if err != nil {
			return nil, err
		}
		m.Indices, err = modeler.ReadIndices(doc, indices, nil)
		if err != nil {
			return nil, fmt.Errorf("read indices: %w", err)
		}
	} else {
		// Non-indexed primitives draw their vertices in order.
		m.Indices = make([]uint32, len(m.Positions))
		for i := range m.Indices {
			m.Indices[i] = uint32(i)
		}
	}
	if len(m.Indices) == 0 || len(m.Indices)%3 != 0 {
		return nil, fmt.Errorf("triangle primitive must have a nonzero multiple of three indices")
	}
	for _, index := range m.Indices {
		if uint64(index) >= uint64(len(m.Positions)) {
			return nil, fmt.Errorf("vertex index %d out of range", index)
		}
	}
	return m, nil
}
