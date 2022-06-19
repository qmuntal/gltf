package gltf

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// Open will open a glTF or GLB file specified by name and return the Document.
func Open(name string) (*Document, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	dec := NewDecoderFS(f, os.DirFS(filepath.Dir(name)))
	doc := new(Document)
	if err = dec.Decode(doc); err != nil {
		doc = nil
	}
	return doc, err
}

// A Decoder reads and decodes glTF and GLB values from an input stream.
//
// Only buffers with relative URIs will be read from Fsys.
// Fsys is called to read external resources.
type Decoder struct {
	Fsys fs.FS
	r    *bufio.Reader
}

// NewDecoder returns a new decoder that reads from r.
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		r: bufio.NewReader(r),
	}
}

// NewDecoder returns a new decoder that reads from r.
func NewDecoderFS(r io.Reader, fsys fs.FS) *Decoder {
	return &Decoder{
		Fsys: fsys,
		r:    bufio.NewReader(r),
	}
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
		if err := d.decodeBinaryBuffer(doc.Buffers[0]); err != nil {
			return err
		}
	}
	for i := externalBufferIndex; i < len(doc.Buffers); i++ {
		if err := d.decodeBuffer(doc.Buffers[i]); err != nil {
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

	err = jd.Decode(doc)
	return isBinary, err
}

func (d *Decoder) readGLBHeader() (*glbHeader, error) {
	var header glbHeader
	chunk, err := d.r.Peek(binary.Size(header))
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
	if header.JSONHeader.Type != glbChunkJSON || (header.JSONHeader.Length+uint32(binary.Size(header))) > header.Length {
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
	} else {
		err = validateBufferURI(buffer.URI)
		if err == nil && d.Fsys != nil {
			uri, ok := sanitizeURI(buffer.URI)
			if ok {
				buffer.Data, err = fs.ReadFile(d.Fsys, uri)
				if len(buffer.Data) > int(buffer.ByteLength) {
					buffer.Data = buffer.Data[:buffer.ByteLength:buffer.ByteLength]
				}
			}
		}
	}
	if err != nil {
		buffer.Data = nil
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
	buffer.Data = make([]byte, buffer.ByteLength)
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

func sanitizeURI(uri string) (string, bool) {
	uri = strings.Replace(uri, "\\", "/", -1)
	uri = strings.Replace(uri, "/./", "/", -1)
	uri = strings.TrimPrefix(uri, "./")
	u, err := url.Parse(uri)
	if err != nil || u.IsAbs() {
		// Only relative paths supported.
		return "", false
	}
	return strings.TrimPrefix(u.RequestURI(), "/"), true
}
