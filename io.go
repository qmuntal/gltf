package gltf

import (
	"io"
	"io/fs"
	"os"
)

// A CreateFS provides access to a hierarchical file system.
// Must follow the same naming convention as io/fs.FS.
type CreateFS interface {
	fs.FS
	Create(name string) (io.WriteCloser, error)
}

// dirFS implements a file system (an fs.FS) for the tree of files rooted at the directory dir.
type dirFS struct {
	fs.FS
	dir string
}

// Create creates or truncates the named file.
func (d dirFS) Create(name string) (io.WriteCloser, error) {
	return os.Create(d.dir + "/" + name)
}
