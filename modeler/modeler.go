package modeler

import (
	"errors"
	"image/color"
	"reflect"

	"github.com/qmuntal/gltf"
	"github.com/qmuntal/gltf/binary"
)

func AddColor(doc *gltf.Document, bufferIndex uint32, data interface{}) (uint32, error) {
	var (
		normalized    bool
		length        int
		componentType gltf.ComponentType
		accesorType   gltf.AccessorType
	)
	switch data := data.(type) {
	case []color.RGBA:
		length = len(data)
		componentType, accesorType = gltf.UnsignedByte, gltf.Vec4
		normalized = true
	case []color.RGBA64:
		length = len(data)
		componentType, accesorType = gltf.UnsignedShort, gltf.Vec4
		normalized = true
	case []gltf.RGB:
		length = len(data)
		componentType, accesorType = gltf.Float, gltf.Vec3
	case []gltf.RGBA:
		length = len(data)
		componentType, accesorType = gltf.Float, gltf.Vec4
	case [][4]uint8:
		length = len(data)
		componentType, accesorType = gltf.UnsignedByte, gltf.Vec4
		normalized = true
	case [][3]uint8:
		length = len(data)
		componentType, accesorType = gltf.UnsignedByte, gltf.Vec3
		normalized = true
	case [][3]uint16:
		length = len(data)
		componentType, accesorType = gltf.UnsignedShort, gltf.Vec3
		normalized = true
	case [][4]uint16:
		length = len(data)
		componentType, accesorType = gltf.UnsignedShort, gltf.Vec4
		normalized = true
	case [][3]float64:
		length = len(data)
		componentType, accesorType = gltf.Float, gltf.Vec3
	case [][4]float64:
		length = len(data)
		componentType, accesorType = gltf.Float, gltf.Vec4
	default:
		return 0, errors.New("modeler.AddColor: invalid type " + reflect.TypeOf(data).String())
	}
	buffer := &doc.Buffers[bufferIndex]
	offset := uint32(len(buffer.Data))
	if err := binary.Write(buffer.Data[offset:], data); err != nil {
		return 0, err
	}
	size := length * binary.SizeOfElement(componentType, accesorType)
	buffer.ByteLength += uint32(size)
	buffer.Data = append(buffer.Data, make([]byte, size)...)
	index := addAccessor(doc, bufferIndex, uint32(length), offset, componentType, accesorType, false)
	doc.Accessors[index].Normalized = normalized
	return uint32(index), nil

}

func padBuffer(buff []uint8) []uint8 {
	paddedLength := getPaddedBufferSize(len(buff))
	if l := paddedLength - len(buff); l > 0 {
		return append(buff, make([]uint8, l)...)
	}
	return buff
}

func getPaddedBufferSize(size int) int {
	return ((size + 4 - 1) / 4) * 4
}

func addAccessor(doc *gltf.Document, buffer, count, offset uint32, componentType gltf.ComponentType, accessorType gltf.AccessorType, isIndex bool) uint32 {
	size := count * uint32(binary.SizeOfElement(componentType, accessorType))
	index := addBufferView(doc, buffer, size, offset, isIndex)
	doc.Accessors = append(doc.Accessors, gltf.Accessor{
		BufferView:    gltf.Index(index),
		ByteOffset:    0,
		ComponentType: componentType,
		Type:          accessorType,
		Count:         count,
	})
	return uint32(len(doc.Accessors) - 1)
}

func addBufferView(doc *gltf.Document, buffer, size, offset uint32, isIndices bool) uint32 {
	bufferView := gltf.BufferView{
		Buffer:     buffer,
		ByteLength: size,
		ByteOffset: offset,
	}
	if isIndices {
		bufferView.Target = gltf.ElementArrayBuffer
	} else {
		bufferView.Target = gltf.ArrayBuffer
	}
	doc.BufferViews = append(doc.BufferViews, bufferView)
	return uint32(len(doc.BufferViews)) - 1
}
