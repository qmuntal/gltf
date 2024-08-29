package modeler

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"sync"

	"github.com/qmuntal/gltf"
	"github.com/qmuntal/gltf/binary"
)

var uint32Pool = sync.Pool{
	New: func() any {
		buf := make([]uint32, 100)
		return &buf
	},
}

// ReadAccessor returns the data references by acr as an slice
// whose element types are the ones associated with acr.ComponentType and acr.Type.
//
// If buffer is not nil, it will be used as backing slice.
//
// ReadAccessor supports all types of accessors: non-interleaved, interleaved, sparse,
// without buffer views, ..., and any combinations of them.
//
// ReadAccessor is safe to use even with malformed documents.
// If that happens it will return an error instead of panic.
func ReadAccessor(doc *gltf.Document, acr *gltf.Accessor, buffer []byte) (any, error) {
	if acr.BufferView == nil && acr.Sparse == nil {
		return nil, nil
	}
	data, err := binary.MakeSliceBuffer(acr.ComponentType, acr.Type, acr.Count, buffer)
	if err != nil {
		return nil, err
	}
	if acr.BufferView != nil {
		buf, err := readBufferView(doc, *acr.BufferView)
		if err != nil {
			return nil, err
		}
		byteStride := doc.BufferViews[*acr.BufferView].ByteStride
		err = binary.Read(buf[acr.ByteOffset:], byteStride, data)
		if err != nil {
			return nil, err
		}
	}

	if acr.Sparse != nil {
		bufPtr := uint32Pool.Get().(*[]uint32)
		defer uint32Pool.Put(bufPtr)
		indices, err := ReadIndices(doc, &gltf.Accessor{
			ComponentType: acr.Sparse.Indices.ComponentType,
			Count:         acr.Sparse.Count,
			Type:          gltf.AccessorScalar,
			BufferView:    &acr.Sparse.Indices.BufferView,
			ByteOffset:    acr.Sparse.Indices.ByteOffset,
		}, *bufPtr)
		if err != nil {
			return nil, err
		}

		valuesBufPtr := bufPool.Get().(*[]byte)
		defer bufPool.Put(valuesBufPtr)
		values, err := ReadAccessor(doc, &gltf.Accessor{
			ComponentType: acr.ComponentType,
			Count:         acr.Sparse.Count,
			Type:          acr.Type,
			BufferView:    &acr.Sparse.Values.BufferView,
			ByteOffset:    acr.Sparse.Values.ByteOffset,
		}, *valuesBufPtr)
		if err != nil {
			return nil, err
		}

		s := reflect.ValueOf(data)
		vals := reflect.ValueOf(values)
		for i := 0; i < int(acr.Sparse.Count); i++ {
			s.Index(int(indices[i])).Set(vals.Index(i))
		}
	}
	return data, nil
}

func readBufferView(doc *gltf.Document, bufferViewIndex int) ([]byte, error) {
	if len(doc.BufferViews) <= bufferViewIndex {
		return nil, errors.New("gltf: bufferview index overflows")
	}
	return ReadBufferView(doc, doc.BufferViews[bufferViewIndex])
}

// ReadBufferView returns the slice of bytes associated with the BufferView.
// The slice is a view of the buffer data, so it is not safe to modify it.
//
// It is safe to use even with malformed documents.
// If that happens it will return an error instead of panic.
func ReadBufferView(doc *gltf.Document, bv *gltf.BufferView) ([]byte, error) {
	if len(doc.Buffers) <= bv.Buffer {
		return nil, errors.New("gltf: buffer index overflows")
	}
	buf := doc.Buffers[bv.Buffer].Data

	high := bv.ByteOffset + bv.ByteLength
	if len(buf) < high {
		return nil, io.ErrShortBuffer
	}
	return buf[bv.ByteOffset:high], nil
}

var bufPool = sync.Pool{
	New: func() any {
		buf := make([]byte, 1024)
		return &buf
	},
}

func makeBufferOf[T any](count int, buffer []T) []T {
	if len(buffer) < count {
		buffer = append(buffer, make([]T, count-len(buffer))...)
	} else {
		buffer = buffer[:count]
	}
	return buffer
}

// ReadIndices returns the data referenced by acr.
// If acr.ComponentType is other than Uint the data
// will be converted appropriately.
//
// See ReadAccessor for more info.
func ReadIndices(doc *gltf.Document, acr *gltf.Accessor, buffer []uint32) ([]uint32, error) {
	switch acr.ComponentType {
	case gltf.ComponentUbyte, gltf.ComponentUshort, gltf.ComponentUint:
	default:
		return nil, errComponentType(acr.ComponentType)
	}
	if acr.Type != gltf.AccessorScalar {
		return nil, errAccessorType(acr.Type)
	}
	bufPtr := bufPool.Get().(*[]byte)
	defer bufPool.Put(bufPtr)
	data, err := ReadAccessor(doc, acr, *bufPtr)
	if err != nil {
		return nil, err
	}
	buffer = makeBufferOf(acr.Count, buffer)
	switch data := data.(type) {
	case []uint8:
		for i, e := range data {
			buffer[i] = uint32(e)
		}
	case []uint16:
		for i, e := range data {
			buffer[i] = uint32(e)
		}
	case []uint32:
		copy(buffer, data)
	}
	return buffer, nil
}

// ReadNormal returns the data referenced by acr.
//
// See ReadAccessor for more info.
func ReadNormal(doc *gltf.Document, acr *gltf.Accessor, buffer [][3]float32) ([][3]float32, error) {
	if acr.ComponentType != gltf.ComponentFloat {
		return nil, errComponentType(acr.ComponentType)
	}
	if acr.Type != gltf.AccessorVec3 {
		return nil, errAccessorType(acr.Type)
	}
	bufPtr := bufPool.Get().(*[]byte)
	defer bufPool.Put(bufPtr)
	data, err := ReadAccessor(doc, acr, *bufPtr)
	if err != nil {
		return nil, err
	}
	buffer = makeBufferOf(acr.Count, buffer)
	copy(buffer, data.([][3]float32))
	return buffer, nil
}

// ReadTangent returns the data referenced by acr.
//
// See ReadAccessor for more info.
func ReadTangent(doc *gltf.Document, acr *gltf.Accessor, buffer [][4]float32) ([][4]float32, error) {
	if acr.ComponentType != gltf.ComponentFloat {
		return nil, errComponentType(acr.ComponentType)
	}
	if acr.Type != gltf.AccessorVec4 {
		return nil, errAccessorType(acr.Type)
	}
	bufPtr := bufPool.Get().(*[]byte)
	defer bufPool.Put(bufPtr)
	data, err := ReadAccessor(doc, acr, *bufPtr)
	if err != nil {
		return nil, err
	}
	buffer = makeBufferOf(acr.Count, buffer)
	copy(buffer, data.([][4]float32))
	return buffer, nil
}

// ReadTextureCoord returns the data referenced by acr.
// If acr.ComponentType is other than Float the data
// will be converted and denormalized appropriately.
//
// See ReadAccessor for more info.
func ReadTextureCoord(doc *gltf.Document, acr *gltf.Accessor, buffer [][2]float32) ([][2]float32, error) {
	switch acr.ComponentType {
	case gltf.ComponentUbyte, gltf.ComponentUshort, gltf.ComponentFloat:
	default:
		return nil, errComponentType(acr.ComponentType)
	}
	if acr.Type != gltf.AccessorVec2 {
		return nil, errAccessorType(acr.Type)
	}
	bufPtr := bufPool.Get().(*[]byte)
	defer bufPool.Put(bufPtr)
	data, err := ReadAccessor(doc, acr, *bufPtr)
	if err != nil {
		return nil, err
	}
	buffer = makeBufferOf(acr.Count, buffer)
	switch data := data.(type) {
	case [][2]uint8:
		for i, e := range data {
			buffer[i] = [2]float32{
				gltf.DenormalizeUbyte(e[0]), gltf.DenormalizeUbyte(e[1]),
			}
		}
	case [][2]uint16:
		for i, e := range data {
			buffer[i] = [2]float32{
				gltf.DenormalizeUshort(e[0]), gltf.DenormalizeUshort(e[1]),
			}
		}
	case [][2]float32:
		copy(buffer, data)
	}
	return buffer, nil
}

// ReadWeights returns the data referenced by acr.
// If acr.ComponentType is other than Float the data
// will be converted and denormalized appropriately.
//
// See ReadAccessor for more info.
func ReadWeights(doc *gltf.Document, acr *gltf.Accessor, buffer [][4]float32) ([][4]float32, error) {
	switch acr.ComponentType {
	case gltf.ComponentUbyte, gltf.ComponentUshort, gltf.ComponentFloat:
	default:
		return nil, errComponentType(acr.ComponentType)
	}
	if acr.Type != gltf.AccessorVec4 {
		return nil, errAccessorType(acr.Type)
	}
	bufPtr := bufPool.Get().(*[]byte)
	defer bufPool.Put(bufPtr)
	data, err := ReadAccessor(doc, acr, *bufPtr)
	if err != nil {
		return nil, err
	}
	buffer = makeBufferOf(acr.Count, buffer)
	switch data := data.(type) {
	case [][4]uint8:
		for i, e := range data {
			buffer[i] = [4]float32{
				gltf.DenormalizeUbyte(e[0]), gltf.DenormalizeUbyte(e[1]),
				gltf.DenormalizeUbyte(e[2]), gltf.DenormalizeUbyte(e[3]),
			}
		}
	case [][4]uint16:
		for i, e := range data {
			buffer[i] = [4]float32{
				gltf.DenormalizeUshort(e[0]), gltf.DenormalizeUshort(e[1]),
				gltf.DenormalizeUshort(e[2]), gltf.DenormalizeUshort(e[3]),
			}
		}
	case [][4]float32:
		copy(buffer, data)
	}
	return buffer, nil
}

// ReadJoints returns the data referenced by acr.
// If acr.ComponentType is other than Ushort the data
// will be converted and denormalized appropriately.
//
// See ReadAccessor for more info.
func ReadJoints(doc *gltf.Document, acr *gltf.Accessor, buffer [][4]uint16) ([][4]uint16, error) {
	switch acr.ComponentType {
	case gltf.ComponentUbyte, gltf.ComponentUshort:
	default:
		return nil, errComponentType(acr.ComponentType)
	}
	if acr.Type != gltf.AccessorVec4 {
		return nil, errAccessorType(acr.Type)
	}
	bufPtr := bufPool.Get().(*[]byte)
	defer bufPool.Put(bufPtr)
	data, err := ReadAccessor(doc, acr, *bufPtr)
	if err != nil {
		return nil, err
	}
	buffer = makeBufferOf(acr.Count, buffer)
	switch data := data.(type) {
	case [][4]uint8:
		for i, e := range data {
			buffer[i] = [4]uint16{
				uint16(e[0]), uint16(e[1]),
				uint16(e[2]), uint16(e[3]),
			}
		}
	case [][4]uint16:
		copy(buffer, data)
	}
	return buffer, nil
}

// ReadPosition returns the data referenced by acr.
//
// See ReadAccessor for more info.
func ReadPosition(doc *gltf.Document, acr *gltf.Accessor, buffer [][3]float32) ([][3]float32, error) {
	if acr.ComponentType != gltf.ComponentFloat {
		return nil, errComponentType(acr.ComponentType)
	}
	if acr.Type != gltf.AccessorVec3 {
		return nil, errAccessorType(acr.Type)
	}
	bufPtr := bufPool.Get().(*[]byte)
	defer bufPool.Put(bufPtr)
	data, err := ReadAccessor(doc, acr, *bufPtr)
	if err != nil {
		return nil, err
	}
	buffer = makeBufferOf(acr.Count, buffer)
	copy(buffer, data.([][3]float32))
	return buffer, nil
}

// ReadColor returns the data referenced by acr.
// If acr.ComponentType is other than Ubyte the data
// will be converted and normalized appropriately.
//
// See ReadAccessor for more info.
func ReadColor(doc *gltf.Document, acr *gltf.Accessor, buffer [][4]uint8) ([][4]uint8, error) {
	switch acr.ComponentType {
	case gltf.ComponentUbyte, gltf.ComponentUshort, gltf.ComponentFloat:
	default:
		return nil, errComponentType(acr.ComponentType)
	}
	switch acr.Type {
	case gltf.AccessorVec3, gltf.AccessorVec4:
	default:
		return nil, errAccessorType(acr.Type)
	}
	bufPtr := bufPool.Get().(*[]byte)
	defer bufPool.Put(bufPtr)
	data, err := ReadAccessor(doc, acr, *bufPtr)
	if err != nil {
		return nil, err
	}
	buffer = makeBufferOf(acr.Count, buffer)
	switch data := data.(type) {
	case [][3]uint8:
		for i, e := range data {
			buffer[i] = [4]uint8{e[0], e[1], e[2], 255}
		}
	case [][4]uint8:
		copy(buffer, data)
	case [][3]uint16:
		for i, e := range data {
			buffer[i] = [4]uint8{uint8(e[0]), uint8(e[1]), uint8(e[2]), 255}
		}
	case [][4]uint16:
		for i, e := range data {
			buffer[i] = [4]uint8{uint8(e[0]), uint8(e[1]), uint8(e[2]), uint8(e[3])}
		}
	case [][3]float32:
		for i, e := range data {
			tmp := gltf.NormalizeRGB(e)
			buffer[i] = [4]uint8{tmp[0], tmp[1], tmp[2], 255}
		}
	case [][4]float32:
		for i, e := range data {
			buffer[i] = gltf.NormalizeRGBA(e)
		}
	}
	return buffer, nil
}

// ReadColor returns the data referenced by acr.
// If acr.ComponentType is other than Ushort the data
// will be converted and normalized appropriately.
//
// See ReadAccessor for more info.
func ReadColor64(doc *gltf.Document, acr *gltf.Accessor, buffer [][4]uint16) ([][4]uint16, error) {
	switch acr.ComponentType {
	case gltf.ComponentUbyte, gltf.ComponentUshort, gltf.ComponentFloat:
	default:
		return nil, errComponentType(acr.ComponentType)
	}
	switch acr.Type {
	case gltf.AccessorVec3, gltf.AccessorVec4:
	default:
		return nil, errAccessorType(acr.Type)
	}
	bufPtr := bufPool.Get().(*[]byte)
	defer bufPool.Put(bufPtr)
	data, err := ReadAccessor(doc, acr, *bufPtr)
	if err != nil {
		return nil, err
	}
	buffer = makeBufferOf(acr.Count, buffer)
	switch data := data.(type) {
	case [][3]uint8:
		for i, e := range data {
			buffer[i] = [4]uint16{
				uint16(e[0]) | uint16(e[0])<<8,
				uint16(e[1]) | uint16(e[1])<<8,
				uint16(e[2]) | uint16(e[2])<<8,
				65535,
			}
		}
	case [][4]uint8:
		for i, e := range data {
			buffer[i] = [4]uint16{
				uint16(e[0]) | uint16(e[0])<<8,
				uint16(e[1]) | uint16(e[1])<<8,
				uint16(e[2]) | uint16(e[2])<<8,
				uint16(e[3]) | uint16(e[3])<<8,
			}
		}
	case [][3]uint16:
		for i, e := range data {
			buffer[i] = [4]uint16{e[0], e[1], e[2], 65535}
		}
	case [][4]uint16:
		copy(buffer, data)
	case [][3]float32:
		for i, e := range data {
			tmp := gltf.NormalizeRGB64(e)
			buffer[i] = [4]uint16{tmp[0], tmp[1], tmp[2], 65535}
		}
	case [][4]float32:
		for i, e := range data {
			buffer[i] = gltf.NormalizeRGBA64(e)
		}
	}
	return buffer, nil
}

func errAccessorType(tp gltf.AccessorType) error {
	return fmt.Errorf("gltf: accessor type %v not allowed", tp)
}

func errComponentType(tp gltf.ComponentType) error {
	return fmt.Errorf("gltf: component type %v not allowed", tp)
}
