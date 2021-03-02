package modeler

import (
	"errors"
	"fmt"
	"io"
	"reflect"

	"github.com/qmuntal/gltf"
	"github.com/qmuntal/gltf/binary"
)

// ReadAccessor returns the data references by acr
// as an slice whose element types are the ones associated with
// acr.ComponentType and acr.Type.
//
// If data is an slice whose elements type matches the accessor type
// then data will be used as backing slice, else a new slice will be allocated.
//
// ReadAccessor supports all types of accessors: non-interleaved, interleaved, sparse,
// without buffer views, ..., and any combinations of them.
//
// ReadAccessor is safe to use even with malformed documents.
// If that happens it will return an error instead of panic.
func ReadAccessor(doc *gltf.Document, acr *gltf.Accessor, data interface{}) (interface{}, error) {
	if acr.BufferView == nil && acr.Sparse == nil {
		return nil, nil
	}
	if data != nil {
		c, t, count := binary.Type(data)
		if count > 0 && c == acr.ComponentType && t == acr.Type {
			if uint32(count) < acr.Count {
				tmpSlice := binary.MakeSlice(acr.ComponentType, acr.Type, acr.Count-uint32(count))
				data = reflect.AppendSlice(reflect.ValueOf(data), reflect.ValueOf(tmpSlice)).Interface()
			} else if uint32(count) > acr.Count {
				data = reflect.ValueOf(data).Slice(0, int(acr.Count)).Interface()
			}
		} else {
			data = binary.MakeSlice(acr.ComponentType, acr.Type, acr.Count)
		}
	} else {
		data = binary.MakeSlice(acr.ComponentType, acr.Type, acr.Count)
	}
	if acr.BufferView != nil {
		buffer, err := readBufferView(doc, *acr.BufferView)
		if err != nil {
			return nil, err
		}
		byteStride := doc.BufferViews[*acr.BufferView].ByteStride
		err = binary.Read(buffer[acr.ByteOffset:], byteStride, data)
		if err != nil {
			return nil, err
		}
	}

	if acr.Sparse != nil {
		indicesBuffer, err := readBufferView(doc, acr.Sparse.Indices.BufferView)
		if err != nil {
			return nil, err
		}

		byteStride := doc.BufferViews[acr.Sparse.Indices.BufferView].ByteStride
		indices := binary.MakeSlice(acr.Sparse.Indices.ComponentType, gltf.AccessorScalar, acr.Sparse.Count)
		err = binary.Read(indicesBuffer[acr.Sparse.Indices.ByteOffset:], byteStride, indices)
		if err != nil {
			return nil, err
		}

		valuesBuffer, err := readBufferView(doc, acr.Sparse.Values.BufferView)
		if err != nil {
			return nil, err
		}
		byteStride = doc.BufferViews[acr.Sparse.Values.ByteOffset].ByteStride
		values := binary.MakeSlice(acr.ComponentType, acr.Type, acr.Sparse.Count)
		err = binary.Read(valuesBuffer[acr.Sparse.Values.ByteOffset:], byteStride, values)
		if err != nil {
			return nil, err
		}

		s := reflect.ValueOf(data)
		ind := reflect.ValueOf(indices)
		vals := reflect.ValueOf(values)
		tp := reflect.TypeOf((*int)(nil)).Elem()
		for i := 0; i < int(acr.Sparse.Count); i++ {
			id := ind.Index(i).Convert(tp).Interface().(int)
			s.Index(id).Set(vals.Index(i))
		}
	}
	return data, nil
}

func readBufferView(doc *gltf.Document, bufferViewIndex uint32) ([]byte, error) {
	if uint32(len(doc.BufferViews)) <= bufferViewIndex {
		return nil, errors.New("gltf: bufferview index overflows")
	}
	return ReadBufferView(doc, doc.BufferViews[bufferViewIndex])
}

// ReadBufferView returns the slice of bytes associated with the BufferView.
//
// It is safe to use even with malformed documents.
// If that happens it will return an error instead of panic.
func ReadBufferView(doc *gltf.Document, bv *gltf.BufferView) ([]byte, error) {
	if uint32(len(doc.Buffers)) <= bv.Buffer {
		return nil, errors.New("gltf: buffer index overflows")
	}
	buf := doc.Buffers[bv.Buffer].Data

	high := bv.ByteOffset + bv.ByteLength
	if uint32(len(buf)) < high {
		return nil, io.ErrShortBuffer
	}
	return buf[bv.ByteOffset:high], nil
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
	data, err := ReadAccessor(doc, acr, buffer)
	if err != nil {
		return nil, err
	}
	if uint32(len(buffer)) < acr.Count {
		buffer = append(buffer, make([]uint32, acr.Count-uint32(len(buffer)))...)
	}
	switch acr.ComponentType {
	case gltf.ComponentUbyte:
		for i, e := range data.([]uint8) {
			buffer[i] = uint32(e)
		}
	case gltf.ComponentUshort:
		for i, e := range data.([]uint16) {
			buffer[i] = uint32(e)
		}
	case gltf.ComponentUint:
		buffer = data.([]uint32)
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
	data, err := ReadAccessor(doc, acr, buffer)
	if err != nil {
		return nil, err
	}
	return data.([][3]float32), nil
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
	data, err := ReadAccessor(doc, acr, buffer)
	if err != nil {
		return nil, err
	}
	return data.([][4]float32), nil
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
	data, err := ReadAccessor(doc, acr, buffer)
	if err != nil {
		return nil, err
	}
	if uint32(len(buffer)) < acr.Count {
		buffer = append(buffer, make([][2]float32, acr.Count-uint32(len(buffer)))...)
	}
	switch acr.ComponentType {
	case gltf.ComponentUbyte:
		for i, e := range data.([][2]uint8) {
			buffer[i] = [2]float32{
				gltf.DenormalizeUbyte(e[0]), gltf.DenormalizeUbyte(e[1]),
			}
		}
	case gltf.ComponentUshort:
		for i, e := range data.([][2]uint16) {
			buffer[i] = [2]float32{
				gltf.DenormalizeUshort(e[0]), gltf.DenormalizeUshort(e[1]),
			}
		}
	case gltf.ComponentFloat:
		buffer = data.([][2]float32)
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
	data, err := ReadAccessor(doc, acr, buffer)
	if err != nil {
		return nil, err
	}
	if uint32(len(buffer)) < acr.Count {
		buffer = append(buffer, make([][4]float32, acr.Count-uint32(len(buffer)))...)
	}
	switch acr.ComponentType {
	case gltf.ComponentUbyte:
		for i, e := range data.([][4]uint8) {
			buffer[i] = [4]float32{
				gltf.DenormalizeUbyte(e[0]), gltf.DenormalizeUbyte(e[1]),
				gltf.DenormalizeUbyte(e[2]), gltf.DenormalizeUbyte(e[3]),
			}
		}
	case gltf.ComponentUshort:
		for i, e := range data.([][4]uint16) {
			buffer[i] = [4]float32{
				gltf.DenormalizeUshort(e[0]), gltf.DenormalizeUshort(e[1]),
				gltf.DenormalizeUshort(e[2]), gltf.DenormalizeUshort(e[3]),
			}
		}
	case gltf.ComponentFloat:
		buffer = data.([][4]float32)
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
	data, err := ReadAccessor(doc, acr, buffer)
	if err != nil {
		return nil, err
	}
	if uint32(len(buffer)) < acr.Count {
		buffer = append(buffer, make([][4]uint16, acr.Count-uint32(len(buffer)))...)
	}
	switch acr.ComponentType {
	case gltf.ComponentUbyte:
		for i, e := range data.([][4]uint8) {
			buffer[i] = [4]uint16{
				uint16(e[0]), uint16(e[1]),
				uint16(e[2]), uint16(e[3]),
			}
		}
	case gltf.ComponentUshort:
		buffer = data.([][4]uint16)
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
	data, err := ReadAccessor(doc, acr, buffer)
	if err != nil {
		return nil, err
	}
	return data.([][3]float32), nil
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
	data, err := ReadAccessor(doc, acr, buffer)
	if err != nil {
		return nil, err
	}
	if uint32(len(buffer)) < acr.Count {
		buffer = append(buffer, make([][4]uint8, acr.Count-uint32(len(buffer)))...)
	}
	switch acr.ComponentType {
	case gltf.ComponentUbyte:
		if acr.Type == gltf.AccessorVec3 {
			for i, e := range data.([][3]uint8) {
				buffer[i] = [4]uint8{e[0], e[1], e[2], 255}
			}
		} else {
			buffer = data.([][4]uint8)
		}
	case gltf.ComponentUshort:
		if acr.Type == gltf.AccessorVec3 {
			for i, e := range data.([][3]uint16) {
				buffer[i] = [4]uint8{uint8(e[0]), uint8(e[1]), uint8(e[2]), 255}
			}
		} else {
			for i, e := range data.([][4]uint16) {
				buffer[i] = [4]uint8{uint8(e[0]), uint8(e[1]), uint8(e[2]), uint8(e[3])}
			}
		}
	case gltf.ComponentFloat:
		if acr.Type == gltf.AccessorVec3 {
			for i, e := range data.([][3]float32) {
				tmp := gltf.NormalizeRGB(e)
				buffer[i] = [4]uint8{tmp[0], tmp[1], tmp[2], 255}
			}
		} else {
			for i, e := range data.([][4]float32) {
				buffer[i] = gltf.NormalizeRGBA(e)
			}
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
	data, err := ReadAccessor(doc, acr, buffer)
	if err != nil {
		return nil, err
	}
	if uint32(len(buffer)) < acr.Count {
		buffer = append(buffer, make([][4]uint16, acr.Count-uint32(len(buffer)))...)
	}
	switch acr.ComponentType {
	case gltf.ComponentUbyte:
		if acr.Type == gltf.AccessorVec3 {
			for i, e := range data.([][3]uint8) {
				buffer[i] = [4]uint16{
					uint16(e[0]) | uint16(e[0])<<8,
					uint16(e[1]) | uint16(e[1])<<8,
					uint16(e[2]) | uint16(e[2])<<8,
					65535}
			}
		} else {
			for i, e := range data.([][4]uint8) {
				buffer[i] = [4]uint16{
					uint16(e[0]) | uint16(e[0])<<8,
					uint16(e[1]) | uint16(e[1])<<8,
					uint16(e[2]) | uint16(e[2])<<8,
					uint16(e[3]) | uint16(e[3])<<8,
				}
			}
		}
	case gltf.ComponentUshort:
		if acr.Type == gltf.AccessorVec3 {
			for i, e := range data.([][3]uint16) {
				buffer[i] = [4]uint16{e[0], e[1], e[2], 65535}
			}
		} else {
			buffer = data.([][4]uint16)
		}
	case gltf.ComponentFloat:
		if acr.Type == gltf.AccessorVec3 {
			for i, e := range data.([][3]float32) {
				tmp := gltf.NormalizeRGB64(e)
				buffer[i] = [4]uint16{tmp[0], tmp[1], tmp[2], 65535}
			}
		} else {
			for i, e := range data.([][4]float32) {
				buffer[i] = gltf.NormalizeRGBA64(e)
			}
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
