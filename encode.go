package gltf

import (
	"encoding/binary"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"unsafe"
)

// WriteResourceCallback defines a callback that will be called when an external resource should be writed.
// The string parameter is the URI of the resource.
type WriteResourceCallback = func(string, int) (io.WriteCloser, error)

// Save will save a document as a glTF or a GLB file specified by name.
func Save(doc *Document, name string, asBinary bool) error {
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	cb := func(uri string, size int) (io.WriteCloser, error) {
		return os.Create(filepath.Join(filepath.Dir(name), uri))
	}
	if err := NewEncoder(f, cb, asBinary).Encode(doc); err != nil {
		f.Close()
		return err
	}
	return f.Close()
}

// An Encoder writes a GLTF to an output stream.
type Encoder struct {
	w        io.Writer
	cb       WriteResourceCallback
	asBinary bool
}

// NewEncoder returns a new encoder that writes to w as a normal glTF file.
func NewEncoder(w io.Writer, cb WriteResourceCallback, asBinary bool) *Encoder {
	return &Encoder{
		w:        w,
		cb:       cb,
		asBinary: asBinary,
	}
}

// Encode writes the encoding of doc to the stream.
func (e *Encoder) Encode(doc *Document) error {
	if doc.Asset.Version == "" {
		doc.Asset.Version = "2.0"
	}
	var err error
	var externalBufferIndex = 0
	if e.asBinary {
		err = e.encodeBinary(doc)
		externalBufferIndex = 1
	} else {
		err = json.NewEncoder(e.w).Encode(doc)
	}
	if err != nil {
		return err
	}

	for i := externalBufferIndex; i < len(doc.Buffers); i++ {
		buffer := &doc.Buffers[i]
		if len(buffer.Data) == 0 || buffer.IsEmbeddedResource() {
			continue
		}
		if err = e.encodeBuffer(buffer); err != nil {
			return err
		}
	}

	return err
}

func (e *Encoder) encodeBuffer(buffer *Buffer) error {
	if len(buffer.Data) == 0 || buffer.IsEmbeddedResource() {
		return nil
	}

	if err := validateBufferURI(buffer.URI); err != nil {
		return err
	}

	r, err := e.cb(buffer.URI, int(buffer.ByteLength))
	if err != nil {
		return err
	}

	_, err = r.Write(buffer.Data)
	if err != nil {
		return err
	}
	return r.Close()
}

func (e *Encoder) encodeBinary(doc *Document) error {
	jsonText, err := json.Marshal(doc)
	if err != nil {
		return err
	}
	header := glbHeader{Magic: glbHeaderMagic, Version: 2, Length: 0, JSONHeader: chunkHeader{Length: 0, Type: glbChunkJSON}}
	binHeader := chunkHeader{Length: 0, Type: glbChunkBIN}
	var binBufferLength uint32
	var binBuffer *Buffer
	if len(doc.Buffers) > 0 {
		binBuffer = &doc.Buffers[0]
		binBufferLength = binBuffer.ByteLength
	}
	binPaddedLength := ((binBufferLength + 3) / 4) * 4
	binPadding := make([]byte, binPaddedLength-binBufferLength)
	binHeader.Length = binPaddedLength

	header.JSONHeader.Length = uint32(((len(jsonText) + 3) / 4) * 4)
	header.Length = uint32(unsafe.Sizeof(header)+unsafe.Sizeof(binHeader)) + header.JSONHeader.Length + binHeader.Length
	headerPadding := make([]byte, header.JSONHeader.Length-uint32(len(jsonText)))
	for i := range headerPadding {
		headerPadding[i] = ' '
	}
	for i := range binPadding {
		binPadding[i] = 0
	}
	err = binary.Write(e.w, binary.LittleEndian, &header)
	if err != nil {
		return err
	}
	binary.Write(e.w, binary.LittleEndian, jsonText)
	binary.Write(e.w, binary.LittleEndian, headerPadding)
	binary.Write(e.w, binary.LittleEndian, &binHeader)
	if binBuffer != nil {
		binary.Write(e.w, binary.LittleEndian, binBuffer.Data)
	}
	return binary.Write(e.w, binary.LittleEndian, binPadding)
}
