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
	"math"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"unsafe"
)

const (
	defaultMaxMemoryAllocation = math.MaxUint32 // 4GB
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
// FS is called to read external resources.
type Decoder struct {
	Fsys                fs.FS
	MaxMemoryAllocation uint64
	r                   *bufio.Reader
}

// NewDecoder returns a new decoder that reads from r.
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		MaxMemoryAllocation: defaultMaxMemoryAllocation,
		r:                   bufio.NewReader(r),
	}
}

// NewDecoder returns a new decoder that reads from r.
func NewDecoderFS(r io.Reader, fsys fs.FS) *Decoder {
	return &Decoder{
		MaxMemoryAllocation: defaultMaxMemoryAllocation,
		Fsys:                fsys,
		r:                   bufio.NewReader(r),
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

func (d *Decoder) validateDocumentQuotas(doc *Document) error {
	var allocs uint64
	for _, b := range doc.Buffers {
		allocs += uint64(b.ByteLength)
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
		err = d.validateDocumentQuotas(doc)
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
	} else if d.Fsys == nil {
		err = errors.New("gltf: external buffer requires Decoder.FS")
	} else if err = validateBufferURI(buffer.URI); err == nil {
		buffer.Data, err = fs.ReadFile(d.Fsys, sanitizeURI(buffer.URI))
		if len(buffer.Data) > int(buffer.ByteLength) {
			buffer.Data = buffer.Data[:buffer.ByteLength:buffer.ByteLength]
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

func sanitizeURI(uri string) string {
	u, err := url.Parse(uri)
	if err == nil {
		uri = strings.TrimPrefix(u.RequestURI(), "/")
	} else {
		uri = strings.Replace(uri, "\\", "/", -1)
		uri = strings.Replace(uri, "/./", "/", -1)
		uri = strings.TrimPrefix(uri, "./")
	}
	return uri
}
