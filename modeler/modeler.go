package modeler

import (
	"errors"
	"image/color"
	"math"
	"reflect"

	"github.com/qmuntal/gltf"
	"github.com/qmuntal/gltf/binary"
)

// Modeler wraps a Document and add usefull methods to build it.
type Modeler struct {
	*gltf.Document
}

// AddIndices adds a new INDICES accessor to the Document
// and fills the buffer with the indices data.
// If success it returns the index of the new accessor.
func (m *Modeler) AddIndices(bufferIndex uint32, data interface{}) (uint32, error) {
	var ok bool
	switch data.(type) {
	case []uint8, []uint16, []uint32:
		ok = true
	}
	componentType, accessorType, length := binary.Type(data)
	if !ok || length <= 0 {
		return 0, errors.New("modeler.AddIndices: invalid type " + reflect.TypeOf(data).String())
	}
	index, err := m.addAccessor(bufferIndex, length, data, componentType, accessorType, true)
	if err != nil {
		return 0, err
	}
	return uint32(index), nil
}

// AddNormal adds a new NORMAL accessor to the Document
// and fills the buffer with the indices data.
// If success it returns the index of the new accessor.
func (m *Modeler) AddNormal(bufferIndex uint32, data interface{}) (uint32, error) {
	var ok bool
	switch data.(type) {
	case [][3]float32:
		ok = true
	}
	componentType, accessorType, length := binary.Type(data)
	if !ok || length <= 0 {
		return 0, errors.New("modeler.AddNormal: invalid type " + reflect.TypeOf(data).String())
	}
	index, err := m.addAccessor(bufferIndex, length, data, componentType, accessorType, true)
	if err != nil {
		return 0, err
	}
	return uint32(index), nil
}

// AddTangent adds a new TANGENT accessor to the Document
// and fills the buffer with the indices data.
// If success it returns the index of the new accessor.
func (m *Modeler) AddTangent(bufferIndex uint32, data interface{}) (uint32, error) {
	var ok bool
	switch data.(type) {
	case [][4]float32:
		ok = true
	}
	componentType, accessorType, length := binary.Type(data)
	if !ok || length <= 0 {
		return 0, errors.New("modeler.AddTangent: invalid type " + reflect.TypeOf(data).String())
	}
	index, err := m.addAccessor(bufferIndex, length, data, componentType, accessorType, true)
	if err != nil {
		return 0, err
	}
	return uint32(index), nil
}

// AddTextureCoord adds a new TEXTURECOORD accessor to the Document
// and fills the buffer with the texturecoord data.
// If success it returns the index of the new accessor.
func (m *Modeler) AddTextureCoord(bufferIndex uint32, data interface{}) (uint32, error) {
	var normalized, ok bool
	switch data.(type) {
	case [][2]uint8, [][2]uint16:
		ok = true
		normalized = true
	case [][2]float32:
		ok = true
	}
	componentType, accessorType, length := binary.Type(data)
	if !ok || length <= 0 {
		return 0, errors.New("modeler.AddTextureCoord: invalid type " + reflect.TypeOf(data).String())
	}
	index, err := m.addAccessor(bufferIndex, length, data, componentType, accessorType, true)
	if err != nil {
		return 0, err
	}
	m.Document.Accessors[index].Normalized = normalized
	return uint32(index), nil
}

// AddWeights adds a new WEIGTHS accessor to the Document
// and fills the buffer with the weights data.
// If success it returns the index of the new accessor.
func (m *Modeler) AddWeights(bufferIndex uint32, data interface{}) (uint32, error) {
	var normalized, ok bool
	switch data.(type) {
	case [][4]uint8, [][4]uint16:
		ok = true
		normalized = true
	case [][4]float32:
		ok = true
	}
	componentType, accessorType, length := binary.Type(data)
	if !ok || length <= 0 {
		return 0, errors.New("modeler.AddWeights: invalid type " + reflect.TypeOf(data).String())
	}
	index, err := m.addAccessor(bufferIndex, length, data, componentType, accessorType, true)
	if err != nil {
		return 0, err
	}
	m.Document.Accessors[index].Normalized = normalized
	return uint32(index), nil
}

// AddJoints adds a new JOINTS accessor to the Document
// and fills the buffer with the joints data.
// If success it returns the index of the new accessor.
func (m *Modeler) AddJoints(bufferIndex uint32, data interface{}) (uint32, error) {
	var ok bool
	switch data.(type) {
	case [][4]uint8, [][4]uint16:
		ok = true
	}
	componentType, accessorType, length := binary.Type(data)
	if !ok || length <= 0 {
		return 0, errors.New("modeler.AddJoints: invalid type " + reflect.TypeOf(data).String())
	}
	index, err := m.addAccessor(bufferIndex, length, data, componentType, accessorType, true)
	if err != nil {
		return 0, err
	}
	return uint32(index), nil
}

// AddPosition adds a new POSITION accessor to the Document
// and fills the buffer with the vertices data.
// If success it returns the index of the new accessor.
func (m *Modeler) AddPosition(bufferIndex uint32, data interface{}) (uint32, error) {
	var (
		ok       bool
		min, max []float64
	)
	switch data := data.(type) {
	case [][3]float32:
		ok = true
		for _, v := range data {
			for i, x := range v {
				min[i] = math.Min(min[i], float64(x))
				max[i] = math.Max(max[i], float64(x))
			}
		}
	}
	componentType, accessorType, length := binary.Type(data)
	if !ok || length <= 0 {
		return 0, errors.New("modeler.AddPosition: invalid type " + reflect.TypeOf(data).String())
	}
	index, err := m.addAccessor(bufferIndex, length, data, componentType, accessorType, false)
	if err != nil {
		return 0, err
	}
	m.Accessors[index].Min = min[:]
	m.Accessors[index].Max = max[:]
	return uint32(index), nil
}

// AddColor adds a new COLOR accessor to the Document
// and fills the buffer with the color data.
// If success it returns the index of the new accessor.
func (m *Modeler) AddColor(bufferIndex uint32, data interface{}) (uint32, error) {
	var normalized, ok bool
	switch data.(type) {
	case []color.RGBA, []color.RGBA64, [][4]uint8, [][3]uint8, [][4]uint16, [][3]uint16:
		normalized = true
		ok = true
	case []gltf.RGBA, []gltf.RGB, [][3]float32, [][4]float32:
		ok = true
	}
	componentType, accessorType, length := binary.Type(data)
	if !ok || length <= 0 {
		return 0, errors.New("modeler.AddColor: invalid type " + reflect.TypeOf(data).String())
	}
	index, err := m.addAccessor(bufferIndex, length, data, componentType, accessorType, false)
	if err != nil {
		return 0, err
	}
	m.Document.Accessors[index].Normalized = normalized
	return uint32(index), nil

}

func (m *Modeler) addAccessor(bufferIndex uint32, count int, data interface{}, componentType gltf.ComponentType, accessorType gltf.AccessorType, isIndex bool) (uint32, error) {
	buffer := &m.Buffers[bufferIndex]
	offset := uint32(len(buffer.Data))
	size := uint32(count * binary.SizeOfElement(componentType, accessorType))
	buffer.ByteLength += uint32(size)
	buffer.Data = append(buffer.Data, make([]byte, size)...)
	if err := binary.Write(buffer.Data[offset:], data); err != nil {
		return 0, err
	}
	index := m.addBufferView(bufferIndex, size, offset, isIndex)
	m.Accessors = append(m.Accessors, gltf.Accessor{
		BufferView:    gltf.Index(index),
		ByteOffset:    0,
		ComponentType: componentType,
		Type:          accessorType,
		Count:         uint32(count),
	})
	return uint32(len(m.Accessors) - 1), nil
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
