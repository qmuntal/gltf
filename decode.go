package gltf

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"unsafe"
)

// ReadQuotas defines maximum allocation sizes to prevent DOS's from malicious files.
type ReadQuotas struct {
	MaxBufferCount      int
	MaxMemoryAllocation int
}

type dataContext struct {
	Quotas       ReadQuotas
	ReadCallback ExternalResourceCallback
}

// ExternalResourceCallback defines a callback that will be called when an external resource should be loaded.
// The string parameter is the URI of the resource.
type ExternalResourceCallback = func(string) (io.Reader, error)

// A Decoder reads and decodes glTF and GLB values from an input stream.
type Decoder struct {
	r      *bufio.Reader
	cb     ExternalResourceCallback
	quotas ReadQuotas
}

// NewDecoder returns a new decoder that reads from r.
func NewDecoder(r io.Reader, cb ExternalResourceCallback) *Decoder {
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

func (d *Decoder) glbHeader() *GLBHeader {
	var header GLBHeader
	chunk, err := d.r.Peek(int(unsafe.Sizeof(header)))
	if err != nil {
		return nil
	}
	r := bytes.NewReader(chunk)
	if err := binary.Read(r, binary.LittleEndian, &header); err != nil {
		return nil
	}
	if header.Magic != glbHeaderMagic {
		return nil
	}
	d.r.Read(chunk)
	return &header
}

func (d *Decoder) chunkHeader() (*ChunkHeader, error) {
	var header ChunkHeader
	chunk := make([]byte, unsafe.Sizeof(header))
	_, err := d.r.Read(chunk)
	if err != nil {
		return nil, err
	}
	r := bytes.NewReader(chunk)
	if err := binary.Read(r, binary.LittleEndian, &header); err != nil {
		return nil, err
	}
	return &header, nil
}

func (d *Decoder) checkForBinary() (bool, error) {
	glbHeader := d.glbHeader()
	if glbHeader == nil {
		return false, nil
	}
	if int(glbHeader.Length) > d.quotas.MaxMemoryAllocation {
		return false, errors.New("gltf: Quota exceeded, bytes of glb buffer > MaxMemoryAllocation")
	}
	jsonHeader, err := d.chunkHeader()
	if err != nil || jsonHeader.Type != glbChunkJSON || (jsonHeader.Length+uint32(unsafe.Sizeof(glbHeader)+unsafe.Sizeof(jsonHeader))) > glbHeader.Length {
		return false, errors.New("gltf: Invalid GLB JSON header")
	}
	return true, nil
}

// Decode reads the next JSON-encoded value from its
// input and stores it in the value pointed to by doc.
func (d *Decoder) Decode(doc *Document) error {
	isBinary, err := d.checkForBinary()
	if err != nil {
		return err
	}
	jd := json.NewDecoder(d.r)
	jd.Decode(doc)
	d.r.Reset(jd.Buffered())
	if len(doc.Buffers) > d.quotas.MaxBufferCount {
		return errors.New("gltf: Quota exceeded, number of buffer > MaxBufferCount")
	}
	if isBinary {
		return d.decodeBinaryBuffer(&doc.Buffers[0])
	}
	for _, buffer := range doc.Buffers {
		if err = d.decodeBuffer(&buffer); err != nil {
			break
		}
	}
	return nil
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
	} else if err = d.validateBufferURI(buffer.URI); err != nil {
		r, err := d.cb(buffer.URI)
		if err == nil {
			buffer.Data = make([]uint8, buffer.ByteLength)
			_, err = r.Read(buffer.Data)
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

func (d *Decoder) validateBufferURI(uri string) error {
	if uri == "" || strings.IndexAny(uri, "..") != -1 || strings.HasPrefix(uri, "/") || strings.HasPrefix(uri, "\\") {
		return fmt.Errorf("gltf: Invalid buffer.uri value '%s'", uri)
	}
	return nil
}
