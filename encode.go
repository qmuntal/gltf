package gltf

import (
	"encoding/binary"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"unsafe"
)

// WriteHandler is the interface that wraps the Write method.
//
// WriteResource should behaves as io.Write in terms of reading the writing resource.
type WriteHandler interface {
	WriteResource(uri string, data []byte) error
}

// Save will save a document as a glTF with the specified by name.
func Save(doc *Document, name string) error {
	return save(doc, name, false)
}

// SaveBinary will save a document as a GLB file with the specified by name.
func SaveBinary(doc *Document, name string) error {
	return save(doc, name, true)
}

func save(doc *Document, name string, asBinary bool) error {
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	e := NewEncoder(f).WithWriteHandler(&RelativeFileHandler{Dir: filepath.Dir(name)})
	e.AsBinary = asBinary
	if err := e.Encode(doc); err != nil {
		f.Close()
		return err
	}
	return f.Close()
}

// An Encoder writes a GLTF to an output stream
// with relative external buffers support.
type Encoder struct {
	AsBinary     bool
	WriteHandler WriteHandler
	w            io.Writer
}

// NewEncoder returns a new encoder that writes to w as a normal glTF file.
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		AsBinary:     true,
		WriteHandler: new(RelativeFileHandler),
		w:            w,
	}
}

// WithWriteHandler sets the WriteHandler.
func (e *Encoder) WithWriteHandler(h WriteHandler) *Encoder {
	e.WriteHandler = h
	return e
}

// Encode writes the encoding of doc to the stream.
func (e *Encoder) Encode(doc *Document) error {
	if doc.Asset.Version == "" {
		doc.Asset.Version = "2.0"
	}
	var err error
	var externalBufferIndex = 0
	if e.AsBinary {
		err = e.encodeBinary(doc)
		externalBufferIndex = 1
	} else {
		err = json.NewEncoder(e.w).Encode(doc)
	}
	if err != nil {
		return err
	}

	for i := externalBufferIndex; i < len(doc.Buffers); i++ {
		buffer := doc.Buffers[i]
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
	if err := validateBufferURI(buffer.URI); err != nil {
		return err
	}

	return e.WriteHandler.WriteResource(buffer.URI, buffer.Data)
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
		binBuffer = doc.Buffers[0]
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
