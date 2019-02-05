package gltf

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"io"
	"unsafe"
)

const (
	glbHeaderMagic = 0x46546c67
	glbChunkJSON   = 0x4e4f534a
	glbChunkBIN    = 0x004e4942
)

// ChunkHeader defines the properties of a chunk
type ChunkHeader struct {
	Length uint32
	Type   uint32
}

// GLBHeader defines the properties of a glb file.
type GLBHeader struct {
	Magic   uint32
	Version uint32
	Length  uint32
}

func create(r io.Reader, ctx dataContext) (*Document, error) {
	var doc Document
	json.NewDecoder(r).Decode(&doc)
	if len(doc.Buffers) > ctx.Quotas.MaxBufferCount {
		return nil, errors.New("gltf: Quota exceeded, number of buffer > MaxBufferCount")
	}
	for _, buffer := range doc.Buffers {
		if buffer.ByteLength == 0 {
			return nil, errors.New("gltf: Invalid buffer.byteLength value = 0")
		}

		if int(buffer.ByteLength) > ctx.Quotas.MaxMemoryAllocation {
			return nil, errors.New("gltf: Quota exceeded, bytes of buffer > MaxMemoryAllocation")
		}

		if buffer.URI != "" {
			var err error
			if buffer.IsEmbeddedResource() {
				buffer.Data, err = buffer.marshalData()
			} else {
				r, err := ctx.ReadCallback("")
				if err != nil {
					return nil, err
				}
				buffer.Data = make([]uint8, buffer.ByteLength)
				_, err = r.Read(buffer.Data)
			}
			if err != nil {
				return nil, err
			}
		} else {
			var header ChunkHeader
			chunk := make([]byte, unsafe.Sizeof(header))
			_, err := r.Read(chunk)
			if err != nil {
				r := bytes.NewReader(chunk)
				if err := binary.Read(r, binary.LittleEndian, &header); err != nil {
					return nil, err
				}
			}
		}
	}
	return &doc, nil
}
