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
type WriteResourceCallback = func(string) (io.WriteCloser, error)

// Save will save a document as a glTF or a GLB file specified by name.
func Save(doc *Document, name string, asBinary bool) error {
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	cb := func(uri string) (io.WriteCloser, error) {
		return os.Create(filepath.Join(filepath.Dir(name), uri))
	}

	if asBinary {
		err = NewEncoderBinary(f, cb).Encode(doc)
	} else {
		err = NewEncoder(f, cb).Encode(doc)
	}
	if err != nil {
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
func NewEncoder(w io.Writer, cb WriteResourceCallback) *Encoder {
	return &Encoder{
		w:        w,
		cb:       cb,
		asBinary: false,
	}
}

// NewEncoderBinary returns a new encoder that writes to w as a binary glTF file.
func NewEncoderBinary(w io.Writer, cb WriteResourceCallback) *Encoder {
	return &Encoder{
		w:        w,
		cb:       cb,
		asBinary: true,
	}
}

// Encode writes the encoding of doc to the stream.
func (e *Encoder) Encode(doc *Document) error {
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
		if err = e.encodeBuffer(buffer); err != nil {
			return err
		}
	}

	return err
}

func (e *Encoder) encodeBuffer(buffer *Buffer) error {
	if buffer.IsEmbeddedResource() {
		return nil
	}
	if err := validateBufferURI(buffer.URI); err != nil {
		return err
	}
	r, err := e.cb(buffer.URI)
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
	header := GLBHeader{Magic: glbHeaderMagic, Version: 2, Length: 0, JSONHeader: ChunkHeader{Length: 0, Type: glbChunkJSON}}
	binHeader := ChunkHeader{Length: 0, Type: glbChunkBIN}
	binBuffer := &doc.Buffers[0]
	binPaddedLength := ((binBuffer.ByteLength + 3) / 4) * 4
	binPadding := make([]byte, binPaddedLength-binBuffer.ByteLength)
	binHeader.Length = uint32(len(binPadding))

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
	binary.Write(e.w, binary.LittleEndian, binBuffer.Data)
	return binary.Write(e.w, binary.LittleEndian, binPadding)
}
