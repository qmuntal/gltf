package gltf

import (
	"io"
	"os"
)

// An FS provides access to a hierarchical file system.
// Must follow the same naming convetion as as io/fs.FS.
type FS interface {
	Open(name string) (io.ReadCloser, error)
}

// A CreateFS provides access to a hierarchical file system.
// Must follow the same naming convetion as as io/fs.FS.
type CreateFS interface {
	Create(name string) (io.WriteCloser, error)
}

// dirFS implements a file system (an fs.FS) for the tree of files rooted at the directory dir.
type dirFS string

// Open opens the named file for reading.
func (dir dirFS) Open(name string) (io.ReadCloser, error) {
	return os.Open(string(dir) + "/" + name)
}

// Create creates or truncates the named file.
func (dir dirFS) Create(name string) (io.WriteCloser, error) {
	return os.Create(string(dir) + "/" + name)
}
