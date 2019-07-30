package gltf

import (
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

// RelativeFileHandler implements a secure ReadHandler supporting relative paths.
// If Dir is empty the os.Getws will be used. It comes with directory traversal protection.
type RelativeFileHandler struct {
	Dir string
}

func (h *RelativeFileHandler) fullName(uri string) string {
	dir := h.Dir
	if dir == "" {
		var err error
		if dir, err = os.Getwd(); err != nil {
			return ""
		}
	}
	return filepath.Join(dir, filepath.FromSlash(path.Clean("/"+uri)))
}

// WriteResource writes the resource using io.WriteFile.
func (h *RelativeFileHandler) WriteResource(uri string, data []byte) error {
	return ioutil.WriteFile(uri, data, 0664)
}

// ReadFullResource reads all the resource data using io.ReadFull.
func (h *RelativeFileHandler) ReadFullResource(uri string, data []byte) error {
	f, err := os.Open(h.fullName(uri))
	if err != nil {
		return err
	}
	_, err = io.ReadFull(f, data)
	f.Close()
	return err
}
