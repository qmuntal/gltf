package modeler

import (
	"errors"
	"image/color"
	"reflect"

	"github.com/qmuntal/gltf"
	"github.com/qmuntal/gltf/binary"
)

// Modeler wraps a Document and add usefull methods to build it.
type Modeler struct {
	*gltf.Document
}

// AddColor adds a new color accessor to the Document
// and fills the buffer with the color data.
// If success it returns the index of the new accessor.
func (m *Modeler) AddColor(bufferIndex uint32, data interface{}) (uint32, error) {
	componentType, accesorType, length := binary.Type(data)
	var (
		normalized bool
		ok         bool
	)
	switch data.(type) {
	case []color.RGBA, []color.RGBA64, [][4]uint8, [][3]uint8, [][4]uint16, [][3]uint16:
		normalized = true
		ok = true
	case []gltf.RGBA, []gltf.RGB, [][3]float32, [][4]float32:
		ok = true
	}
	if !ok || length <= 0 {
		return 0, errors.New("modeler.AddColor: invalid type " + reflect.TypeOf(data).String())
	}
	buffer := &m.Buffers[bufferIndex]
	offset := uint32(len(buffer.Data))
	size := length * binary.SizeOfElement(componentType, accesorType)
	buffer.ByteLength += uint32(size)
	buffer.Data = append(buffer.Data, make([]byte, size)...)
	if err := binary.Write(buffer.Data[offset:], data); err != nil {
		return 0, err
	}
	index := m.addAccessor(bufferIndex, uint32(length), offset, componentType, accesorType, false)
	m.Document.Accessors[index].Normalized = normalized
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

func (m *Modeler) addAccessor(buffer, count, offset uint32, componentType gltf.ComponentType, accessorType gltf.AccessorType, isIndex bool) uint32 {
	size := count * uint32(binary.SizeOfElement(componentType, accessorType))
	index := m.addBufferView(buffer, size, offset, isIndex)
	m.Accessors = append(m.Accessors, gltf.Accessor{
		BufferView:    gltf.Index(index),
		ByteOffset:    0,
		ComponentType: componentType,
		Type:          accessorType,
		Count:         count,
	})
	return uint32(len(m.Accessors) - 1)
}

func (m *Modeler) addBufferView(buffer, size, offset uint32, isIndices bool) uint32 {
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
	m.BufferViews = append(m.BufferViews, bufferView)
	return uint32(len(m.BufferViews)) - 1
}
