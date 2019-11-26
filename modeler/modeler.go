package modeler

import (
	"bytes"
	"image/color"
	"io"
	"io/ioutil"
	"math"
	"reflect"

	"github.com/qmuntal/gltf"
	"github.com/qmuntal/gltf/binary"
)

// CompressionLevel defines the different levels of compression.
type CompressionLevel uint8

const (
	// CompressionNone will not apply any compression.
	CompressionNone CompressionLevel = iota
	// CompressionLossless will reduce the byte size without sacrificing quality.
	CompressionLossless
)

// Modeler wraps a Document and add usefull methods to build it.
// If Compress is true, all the data added to accessors that support different component types
// will be evaluated to see if it fits in a smaller component type.
type Modeler struct {
	*gltf.Document
	Compression CompressionLevel
}

// NewModeler returns a new Modeler instance.
func NewModeler() *Modeler {
	return &Modeler{
		Document: &gltf.Document{
			Scene:  gltf.Index(0),
			Scenes: []gltf.Scene{{Name: "Root Scene"}},
		},
		Compression: CompressionLossless,
	}
}

// AddIndices adds a new INDICES accessor to the Document
// and fills the buffer with the indices data.
// If success it returns the index of the new accessor.
func (m *Modeler) AddIndices(bufferIndex uint32, data interface{}) uint32 {
	var ok bool
	switch data.(type) {
	case []uint8:
		ok = true
	case []uint16:
		ok = true
		if m.Compression >= CompressionLossless {
			data = compressUint16(data.([]uint16))
		}
	case []uint32:
		ok = true
		if m.Compression >= CompressionLossless {
			data = compressUint32(data.([]uint32))
		}
	}
	componentType, accessorType, length := binary.Type(data)
	if !ok || length <= 0 {
		panic("modeler.AddIndices: invalid type " + reflect.TypeOf(data).String())
	}
	index := m.addAccessor(bufferIndex, length, data, componentType, accessorType, true)
	return uint32(index)
}

// AddNormal adds a new NORMAL accessor to the Document
// and fills the buffer with the indices data.
// If success it returns the index of the new accessor.
func (m *Modeler) AddNormal(bufferIndex uint32, data [][3]float32) uint32 {
	componentType, accessorType, length := binary.Type(data)
	index := m.addAccessor(bufferIndex, length, data, componentType, accessorType, false)
	return uint32(index)
}

// AddTangent adds a new TANGENT accessor to the Document
// and fills the buffer with the indices data.
// If success it returns the index of the new accessor.
func (m *Modeler) AddTangent(bufferIndex uint32, data [][4]float32) uint32 {
	componentType, accessorType, length := binary.Type(data)
	index := m.addAccessor(bufferIndex, length, data, componentType, accessorType, false)
	return uint32(index)
}

// AddTextureCoord adds a new TEXTURECOORD accessor to the Document
// and fills the buffer with the texturecoord data.
// If success it returns the index of the new accessor.
func (m *Modeler) AddTextureCoord(bufferIndex uint32, data interface{}) uint32 {
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
		panic("modeler.AddTextureCoord: invalid type " + reflect.TypeOf(data).String())
	}
	index := m.addAccessor(bufferIndex, length, data, componentType, accessorType, false)
	m.Document.Accessors[index].Normalized = normalized
	return uint32(index)
}

// AddWeights adds a new WEIGTHS accessor to the Document
// and fills the buffer with the weights data.
// If success it returns the index of the new accessor.
func (m *Modeler) AddWeights(bufferIndex uint32, data interface{}) uint32 {
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
		panic("modeler.AddWeights: invalid type " + reflect.TypeOf(data).String())
	}
	index := m.addAccessor(bufferIndex, length, data, componentType, accessorType, false)
	m.Document.Accessors[index].Normalized = normalized
	return uint32(index)
}

// AddJoints adds a new JOINTS accessor to the Document
// and fills the buffer with the joints data.
// If success it returns the index of the new accessor.
func (m *Modeler) AddJoints(bufferIndex uint32, data interface{}) uint32 {
	var ok bool
	switch data.(type) {
	case [][4]uint8, [][4]uint16:
		ok = true
	}
	componentType, accessorType, length := binary.Type(data)
	if !ok || length <= 0 {
		panic("modeler.AddJoints: invalid type " + reflect.TypeOf(data).String())
	}
	index := m.addAccessor(bufferIndex, length, data, componentType, accessorType, false)
	return uint32(index)
}

// AddPosition adds a new POSITION accessor to the Document
// and fills the buffer with the vertices data.
// If success it returns the index of the new accessor.
func (m *Modeler) AddPosition(bufferIndex uint32, data [][3]float32) uint32 {
	min := [3]float64{math.MaxFloat64, math.MaxFloat64, math.MaxFloat64}
	max := [3]float64{-math.MaxFloat64, -math.MaxFloat64, -math.MaxFloat64}
	for _, v := range data {
		for i, x := range v {
			min[i] = math.Min(min[i], float64(x))
			max[i] = math.Max(max[i], float64(x))
		}
	}
	componentType, accessorType, length := binary.Type(data)
	index := m.addAccessor(bufferIndex, length, data, componentType, accessorType, false)
	m.Accessors[index].Min = min[:]
	m.Accessors[index].Max = max[:]
	return uint32(index)
}

// AddColor adds a new COLOR accessor to the Document
// and fills the buffer with the color data.
// If success it returns the index of the new accessor.
func (m *Modeler) AddColor(bufferIndex uint32, data interface{}) uint32 {
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
		panic("modeler.AddColor: invalid type " + reflect.TypeOf(data).String())
	}
	index := m.addAccessor(bufferIndex, length, data, componentType, accessorType, false)
	m.Document.Accessors[index].Normalized = normalized
	return uint32(index)
}

// AddImage adds a new image to the Document
// and fills the buffer with the image data.
// If success it returns the index of the new image.
func (m *Modeler) AddImage(bufferIndex uint32, name, mimeType string, r io.Reader) (uint32, error) {
	buffer := m.buffer(bufferIndex)
	offset := uint32(len(buffer.Data))
	switch r := r.(type) {
	case *bytes.Buffer:
		buffer.Data = append(buffer.Data, r.Bytes()...)
	default:
		data, err := ioutil.ReadAll(r)
		if err != nil {
			return 0, err
		}
		buffer.Data = append(buffer.Data, data...)
	}
	index := m.addBufferView(bufferIndex, uint32(len(buffer.Data))-offset, offset, 0, false)
	buffer.ByteLength += uint32(len(buffer.Data))
	m.BufferViews[index].Target = gltf.None
	m.Images = append(m.Images, gltf.Image{
		Name:       name,
		MimeType:   mimeType,
		BufferView: gltf.Index(index),
	})
	return uint32(len(m.Images) - 1), nil
}

func (m *Modeler) addAccessor(bufferIndex uint32, count int, data interface{}, componentType gltf.ComponentType, accessorType gltf.AccessorType, isIndex bool) uint32 {
	buffer := m.buffer(bufferIndex)
	offset := uint32(len(buffer.Data))
	padding := ((offset+3)/4)*4 - offset
	elementSize := binary.SizeOfElement(componentType, accessorType)
	size := uint32(count * elementSize)
	buffer.ByteLength += uint32(size + padding)
	buffer.Data = append(buffer.Data, make([]byte, size+padding)...)
	// Cannot return error as the buffer has enough size and the data type is controled.
	_ = binary.Write(buffer.Data[offset+padding:], data)
	index := m.addBufferView(bufferIndex, size, offset+padding, uint32(elementSize), isIndex)
	m.Accessors = append(m.Accessors, gltf.Accessor{
		BufferView:    gltf.Index(index),
		ByteOffset:    0,
		ComponentType: componentType,
		Type:          accessorType,
		Count:         uint32(count),
	})
	return uint32(len(m.Accessors) - 1)
}

func (m *Modeler) addBufferView(buffer, size, offset, stride uint32, isIndices bool) uint32 {
	bufferView := gltf.BufferView{
		Buffer:     buffer,
		ByteLength: size,
		ByteOffset: offset,
	}
	if isIndices {
		bufferView.Target = gltf.ElementArrayBuffer
	} else {
		bufferView.Target = gltf.ArrayBuffer
		bufferView.ByteStride = stride
	}
	m.BufferViews = append(m.BufferViews, bufferView)
	return uint32(len(m.BufferViews)) - 1
}

func (m *Modeler) buffer(bufferIndex uint32) *gltf.Buffer {
	if int(bufferIndex) >= len(m.Buffers) {
		m.Buffers = append(m.Buffers, make([]gltf.Buffer, int(bufferIndex)-len(m.Buffers)+1)...)
	}
	return &m.Buffers[bufferIndex]
}
