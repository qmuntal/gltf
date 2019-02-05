package gltf

import (
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

// ReadQuotas defines maximum allocation sizes to prevent DOS's from malicious files.
type ReadQuotas struct {
	MaxBufferCount      int
	MaxMemoryAllocation int
}

// ReadResourceCallback defines a callback that will be called when an external resource should be loaded.
// The string parameter is the URI of the resource.
type ReadResourceCallback = func(string) (io.ReadCloser, error)

// Open will open a glTF or GLB file specified by name and return the Document.
func Open(name string) (*Document, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	cb := func(uri string) (io.ReadCloser, error) {
		return os.Open(filepath.Join(filepath.Dir(name), uri))
	}
	doc := new(Document)
	err = NewDecoder(f, cb).Decode(doc)
	return doc, err
}

// A Decoder reads and decodes glTF and GLB values from an input stream.
type Decoder struct {
	r      *bufio.Reader
	cb     ReadResourceCallback
	quotas ReadQuotas
}

// NewDecoder returns a new decoder that reads from r.
func NewDecoder(r io.Reader, cb ReadResourceCallback) *Decoder {
	return &Decoder{
		r:  bufio.NewReader(r),
		cb: cb,
		quotas: ReadQuotas{
			MaxBufferCount:      8,
			MaxMemoryAllocation: 32 * 1024 * 1024,
		}}
}

// SetQuotas sets the read memory limits.
func (d *Decoder) SetQuotas(quotas ReadQuotas) {
	d.quotas = quotas
}

// Decode reads the next JSON-encoded value from its
// input and stores it in the value pointed to by doc.
func (d *Decoder) Decode(doc *Document) error {
	isBinary, err := d.decodeDocument(doc)
	if err != nil {
		return nil
	}
	if len(doc.Buffers) > d.quotas.MaxBufferCount {
		return errors.New("gltf: Quota exceeded, number of buffer > MaxBufferCount")
	}
	if isBinary {
		return d.decodeBinaryBuffer(&doc.Buffers[0])
	}
	for _, buffer := range doc.Buffers {
		if err := d.decodeBuffer(&buffer); err != nil {
			break
		}
	}
	return nil
}

func (d *Decoder) decodeDocument(doc *Document) (isBinary bool, err error) {
	glbHeader, err := d.readGLBHeader()
	if err != nil {
		return
	}
	var jd *json.Decoder
	if glbHeader != nil {
		jd = json.NewDecoder(&io.LimitedReader{R: d.r, N: int64(glbHeader.JSONHeader.Length)})
		isBinary = true
	} else {
		jd = json.NewDecoder(d.r)
		isBinary = false
	}

	err = jd.Decode(doc)
	return
}

func (d *Decoder) readGLBHeader() (*GLBHeader, error) {
	var header GLBHeader
	chunk, err := d.r.Peek(int(unsafe.Sizeof(header)))
	if err != nil {
		return nil, nil
	}
	r := bytes.NewReader(chunk)
	if err := binary.Read(r, binary.LittleEndian, &header); err != nil {
		return nil, nil
	}
	if header.Magic != glbHeaderMagic {
		return nil, nil
	}
	d.r.Read(chunk)
	if int(header.Length) > d.quotas.MaxMemoryAllocation {
		return nil, errors.New("gltf: Quota exceeded, bytes of glb buffer > MaxMemoryAllocation")
	}
	if header.JSONHeader.Type != glbChunkJSON || (header.JSONHeader.Length+uint32(unsafe.Sizeof(header))) > header.Length {
		return nil, errors.New("gltf: Invalid GLB JSON header")
	}
	return &header, nil
}

func (d *Decoder) chunkHeader() (*ChunkHeader, error) {
	var header ChunkHeader
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
		r, err := d.cb(buffer.URI)
		if err == nil {
			buffer.Data = make([]uint8, buffer.ByteLength)
			_, err = r.Read(buffer.Data)
		}
		r.Close()
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
	_, err = d.r.Read(buffer.Data)
	return err
}

func (d *Decoder) validateBuffer(buffer *Buffer) error {
	if buffer.ByteLength == 0 {
		return errors.New("gltf: Invalid buffer.byteLength value = 0")
	}

	if int(buffer.ByteLength) > d.quotas.MaxMemoryAllocation {
		return errors.New("gltf: Quota exceeded, bytes of buffer > MaxMemoryAllocation")
	}
	return nil
}

func validateBufferURI(uri string) error {
	if uri == "" || strings.Contains(uri, "..") || strings.HasPrefix(uri, "/") || strings.HasPrefix(uri, "\\") {
		return fmt.Errorf("gltf: Invalid buffer.uri value '%s'", uri)
	}
	return nil
}
