package gltf

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"unsafe"
)

// ProtocolReadHandler defines a ReadFull interface.
//
// ReadFull should behaves as io.ReadFull in terms of reading the external resource.
// The data already has the correct size so it can be used directly to store the read output.
type ProtocolReadHandler interface {
	ReadFull(scheme, uri string, data []byte) error
}

// ProtocolRegistry implements a secure ProtocolReadHandler supporting http, https and relative paths.
// If Dir is empty the os.Getws will be used. It comes with directory traversal protection.
// If HTTPClient is nil http[s] will not be supported.
type ProtocolRegistry struct {
	Dir        string
	HTTPClient *http.Client
}

func (reg *ProtocolRegistry) readRelativeFile(uri string, data []byte) (err error) {
	dir := reg.Dir
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

// ReadFull should as io.ReadFull in terms of reading the external resource.
// An error is returned when the scheme is not supported.
func (reg *ProtocolRegistry) ReadFull(scheme, uri string, data []byte) error {
	switch scheme {
	case "": // probably relative path
		return reg.readRelativeFile(uri, data)
	case "http://", "https://":
		if reg.HTTPClient != nil {
			resp, err := reg.HTTPClient.Get(uri)
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				return fmt.Errorf("gltf: server responded with %d", resp.StatusCode)
			}
			_, err = io.ReadFull(resp.Body, data)
			return err
		}
	}
	return errors.New("gltf: scheme not supported")
}

// Open will open a glTF or GLB file specified by name and return the Document.
func Open(name string) (*Document, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	doc := new(Document)
	err = NewDecoder(f).WithProtocolReadHandler(&ProtocolRegistry{Dir: filepath.Dir(name)}).Decode(doc)
	f.Close()
	return doc, err
}

// A Decoder reads and decodes glTF and GLB values from an input stream.
// Callback is called to read external resources.
// If Callback is nil the external resource data in not loaded.
type Decoder struct {
	ProtocolReadHandler ProtocolReadHandler
	r                   *bufio.Reader
}

// NewDecoder returns a new decoder that reads from r.
// By default the external buffers are not read.
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		ProtocolReadHandler: new(ProtocolRegistry),
		r:                   bufio.NewReader(r),
	}
}

// WithProtocolReadHandler sets the ProtocolReadHandler.
func (d *Decoder) WithProtocolReadHandler(reg ProtocolReadHandler) *Decoder {
	d.ProtocolReadHandler = reg
	return d
}

// Decode reads the next JSON-encoded value from its
// input and stores it in the value pointed to by doc.
func (d *Decoder) Decode(doc *Document) error {
	isBinary, err := d.decodeDocument(doc)
	if err != nil {
		return err
	}

	var externalBufferIndex = 0
	if isBinary && len(doc.Buffers) > 0 {
		externalBufferIndex = 1
		if err := d.decodeBinaryBuffer(&doc.Buffers[0]); err != nil {
			return err
		}
	}
	for i := externalBufferIndex; i < len(doc.Buffers); i++ {
		if err := d.decodeBuffer(&doc.Buffers[i]); err != nil {
			return err
		}
	}
	return nil
}

func (d *Decoder) decodeDocument(doc *Document) (bool, error) {
	glbHeader, err := d.readGLBHeader()
	if err != nil {
		return false, err
	}
	var (
		jd       *json.Decoder
		isBinary bool
	)
	if glbHeader != nil {
		jd = json.NewDecoder(&io.LimitedReader{R: d.r, N: int64(glbHeader.JSONHeader.Length)})
		isBinary = true
	} else {
		jd = json.NewDecoder(d.r)
		isBinary = false
	}

	return isBinary, jd.Decode(doc)
}

func (d *Decoder) readGLBHeader() (*glbHeader, error) {
	var header glbHeader
	chunk, err := d.r.Peek(int(unsafe.Sizeof(header)))
	if err != nil {
		return nil, nil
	}
	r := bytes.NewReader(chunk)
	binary.Read(r, binary.LittleEndian, &header)
	if header.Magic != glbHeaderMagic {
		return nil, nil
	}
	d.r.Read(chunk)
	return &header, d.validateGLBHeader(&header)
}

func (d *Decoder) validateGLBHeader(header *glbHeader) error {
	if header.JSONHeader.Type != glbChunkJSON || (header.JSONHeader.Length+uint32(unsafe.Sizeof(header))) > header.Length {
		return errors.New("gltf: Invalid GLB JSON header")
	}
	return nil
}

func (d *Decoder) chunkHeader() (*chunkHeader, error) {
	var header chunkHeader
	if err := binary.Read(d.r, binary.LittleEndian, &header); err != nil {
		return nil, err
	}
	return &header, nil
}

func (d *Decoder) decodeBuffer(buffer *Buffer) error {
	if err := d.validateBuffer(buffer); err != nil {
		return err
	}
	if buffer.URI == "" {
		return errors.New("gltf: buffer without URI")
	}
	var err error
	if buffer.IsEmbeddedResource() {
		buffer.Data, err = buffer.marshalData()
	} else if err = validateBufferURI(buffer.URI); err == nil {
		buffer.Data = make([]uint8, buffer.ByteLength)
		var u *url.URL
		u, err = url.Parse(buffer.URI)
		if err == nil {
			err = d.ProtocolReadHandler.ReadFull(u.Scheme, buffer.URI, buffer.Data)
		}
	}
	return err
}

func (d *Decoder) decodeBinaryBuffer(buffer *Buffer) error {
	if err := d.validateBuffer(buffer); err != nil {
		return err
	}
	header, err := d.chunkHeader()
	if err != nil {
		return err
	}
	if header.Type != glbChunkBIN || header.Length < buffer.ByteLength {
		return errors.New("gltf: Invalid GLB BIN header")
	}
	buffer.Data = make([]uint8, buffer.ByteLength)
	_, err = io.ReadFull(d.r, buffer.Data)
	return err
}

func (d *Decoder) validateBuffer(buffer *Buffer) error {
	if buffer.ByteLength == 0 {
		return errors.New("gltf: Invalid buffer.byteLength value = 0")
	}
	return nil
}

func validateBufferURI(uri string) error {
	if uri == "" || strings.Contains(uri, "..") || strings.HasPrefix(uri, "/") || strings.HasPrefix(uri, "\\") {
		return fmt.Errorf("gltf: Invalid buffer.uri value '%s'", uri)
	}
	return nil
}
