package gltf

import (
	"encoding/binary"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"unsafe"
)

// WriteResourceCallback defines a callback that will be called when an external resource should be writed.
// The string parameter is the URI of the resource.
type WriteResourceCallback = func(string, []byte) error

func discardWriteData(uri string, data []byte) error {
	return nil
}

// Save will save a document as a glTF with the specified by name.
func Save(doc *Document, name string) error {
	return save(doc, name, false)
}

// Save will save a document as a GLB file with the specified by name.
func SaveBinary(doc *Document, name string) error {
	return save(doc, name, true)
}

func save(doc *Document, name string, asBinary bool) error {
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	e := NewEncoder(f)
	e.AsBinary = asBinary
	e.WithCallback(func(uri string, data []byte) error {
		return ioutil.WriteFile(uri, data, 0664)
	})
	if err := e.Encode(doc); err != nil {
		f.Close()
		return err
	}
	return f.Close()
}

// An Encoder writes a GLTF to an output stream.
// The callback is called when an external resource shouldbe writed.
type Encoder struct {
	AsBinary bool
	Callback WriteResourceCallback
	w        io.Writer
}

// NewEncoder returns a new encoder that writes to w as a normal glTF file.
// By default the file is writed as binary and external buffers data is discarded.
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		AsBinary: true,
		Callback: discardWriteData,
		w:        w,
	}
}

// WithCallback sets the ReadResourceCallback.
func (e *Encoder) WithCallback(c WriteResourceCallback) *Encoder {
	e.Callback = c
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
	if err := validateBufferURI(buffer.URI); err != nil {
		return err
	}

	return e.Callback(buffer.URI, buffer.Data)
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
