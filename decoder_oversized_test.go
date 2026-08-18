package gltf

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"testing"
)

// TestDecodeBinaryBufferOversized ensures a GLB whose buffer byteLength is far
// larger than the data actually present does not trigger a huge allocation.
// Before the fix, decodeBinaryBuffer ran make([]byte, buffer.ByteLength) and
// attempted a multi-GiB allocation from this ~96 byte input.
func TestDecodeBinaryBufferOversized(t *testing.T) {
	const byteLength = 0xFFFFFFFC // ~4 GiB, small enough that byteLength+3 does not overflow uint32

	jsonChunk := []byte(`{"asset":{"version":"2.0"},"buffers":[{"byteLength":4294967292}]}`)
	for len(jsonChunk)%4 != 0 {
		jsonChunk = append(jsonChunk, ' ')
	}

	var b bytes.Buffer
	binary.Write(&b, binary.LittleEndian, uint32(glbHeaderMagic))
	binary.Write(&b, binary.LittleEndian, uint32(2))              // version
	binary.Write(&b, binary.LittleEndian, uint32(0xFFFFFFF0))     // overall length
	binary.Write(&b, binary.LittleEndian, uint32(len(jsonChunk))) // JSON chunk length
	binary.Write(&b, binary.LittleEndian, uint32(glbChunkJSON))   // JSON chunk type
	b.Write(jsonChunk)
	binary.Write(&b, binary.LittleEndian, uint32(byteLength))  // BIN chunk length
	binary.Write(&b, binary.LittleEndian, uint32(glbChunkBIN)) // BIN chunk type
	// no BIN data follows

	err := NewDecoder(bytes.NewReader(b.Bytes())).Decode(new(Document))
	if err == nil {
		t.Fatal("expected error for truncated oversized binary buffer, got nil")
	}
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Errorf("expected io.ErrUnexpectedEOF; got %v", err)
	}
}
