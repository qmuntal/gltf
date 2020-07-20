package gltf

import (
	"math"
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"unsafe"
)

const (
	defaultMaxExternalBufferCount  = 10
	defaultMaxMemoryAllocation     = math.MaxUint32 // 4GB
)

// ReadHandler is the interface that wraps the ReadFullResource method.
//
// ReadFullResource should behaves as io.ReadFull in terms of reading the external resource.
// The data already has the correct size so it can be used directly to store the read output.
type ReadHandler interface {
	ReadFullResource(uri string, data []byte) error
}

// Open will open a glTF or GLB file specified by name and return the Document.
func Open(name string) (*Document, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	dec := NewDecoder(f).WithReadHandler(&RelativeFileHandler{Dir: filepath.Dir(name)})
	doc := new(Document)
	if err = dec.Decode(doc); err != nil {
		doc = nil
	}
	return doc, err
}

// A Decoder reads and decodes glTF and GLB values from an input stream.
// ReadHandler is called to read external resources.
type Decoder struct {
	ReadHandler            ReadHandler
	MaxExternalBufferCount int
	MaxMemoryAllocation    uint64
	r                      *bufio.Reader
}

// NewDecoder returns a new decoder that reads from r
// with relative external buffers support.
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		ReadHandler:            new(RelativeFileHandler),
		MaxExternalBufferCount: defaultMaxExternalBufferCount,
		MaxMemoryAllocation:    defaultMaxMemoryAllocation,
		r:                      bufio.NewReader(r),
	}
}

// WithReadHandler sets the ReadHandler.
func (d *Decoder) WithReadHandler(h ReadHandler) *Decoder {
	d.ReadHandler = h
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

func (d *Decoder) validateDocumentQuotas(doc *Document, isBinary bool) error {
	var externalCount int
	var allocs uint64
	for _, b := range doc.Buffers {
		allocs += uint64(b.ByteLength)
		if !b.IsEmbeddedResource() {
			externalCount++
		}
	}
	if isBinary {
		externalCount--
	}
	if externalCount > d.MaxExternalBufferCount {
		return errors.New("gltf: External buffer count quota exceeded")
	}
	if allocs > d.MaxMemoryAllocation {
		return errors.New("gltf: Memory allocation count quota exceeded")
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
	if err == nil {
		err = d.validateDocumentQuotas(doc, isBinary)
	}
	return isBinary, err
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
		buffer.Data = make([]byte, buffer.ByteLength)
		err = d.ReadHandler.ReadFullResource(buffer.URI, buffer.Data)
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
