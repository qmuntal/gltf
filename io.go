package gltf

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
)

// RelativeFileHandler implements a secure ReadHandler supporting relative paths.
// If Dir is empty the os.Getws will be used. It comes with directory traversal protection.
type RelativeFileHandler struct {
	Dir string
}

// ReadFull should as io.ReadFull in terms of reading the external resource.
func (h *RelativeFileHandler) ReadFull(uri string, data []byte) (err error) {
	dir := h.Dir
	if dir == "" {
		if dir, err = os.Getwd(); err != nil {
			return
		}
	}
	var f http.File
	f, err = http.Dir(dir).Open(uri)
	if err != nil {
		return
	}
	_, err = io.ReadFull(f, data)
	return
}

// ProtocolRegistry implements a secure ProtocolReadHandler as a map of supported schemes.
type ProtocolRegistry map[string]ReadHandler

// ReadFull should as io.ReadFull in terms of reading the external resource.
// An error is returned when the scheme is not supported.
func (reg ProtocolRegistry) ReadFull(uri string, data []byte) error {
	u, err := url.Parse(uri)
	if err != nil {
		return err
	}
	if f, ok := reg[u.Scheme]; ok {
		return f.ReadFull(uri, data)
	}
	return errors.New("gltf: not supported scheme")
}
