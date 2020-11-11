// Package modeler implements helper methods to write common glTF entities
// (indices, positions, colors, ...) into buffers.
package modeler

import (
	"bytes"
	"fmt"
	"image/color"
	"io"
	"io/ioutil"
	"math"

	"github.com/qmuntal/gltf"
	"github.com/qmuntal/gltf/binary"
)

// WriteIndices adds a new INDICES accessor to doc
// and fills the last buffer with the indices data.
// If success it returns the index of the new accessor.
func WriteIndices(doc *gltf.Document, data interface{}) uint32 {
	switch data.(type) {
	case []uint8, []uint16, []uint32:
	default:
		panic(fmt.Sprintf("modeler.WriteIndices: invalid type %T", data))
	}
	return WriteAccessor(doc, gltf.TargetElementArrayBuffer, data)
}

// WriteNormal adds a new NORMAL accessor to doc
// and fills the last buffer with the indices data.
// If success it returns the index of the new accessor.
func WriteNormal(doc *gltf.Document, data [][3]float32) uint32 {
	return WriteAccessor(doc, gltf.TargetArrayBuffer, data)
}

// WriteTangent adds a new TANGENT accessor to doc
// and fills the last buffer with the indices data.
// If success it returns the index of the new accessor.
func WriteTangent(doc *gltf.Document, data [][4]float32) uint32 {
	return WriteAccessor(doc, gltf.TargetArrayBuffer, data)
}

// WriteTextureCoord adds a new TEXTURECOORD accessor to doc
// and fills the last buffer with the texturecoord data.
// If success it returns the index of the new accessor.
func WriteTextureCoord(doc *gltf.Document, data interface{}) uint32 {
	var normalized bool
	switch data.(type) {
	case [][2]uint8, [][2]uint16:
		normalized = true
	case [][2]float32:
	default:
		panic(fmt.Sprintf("modeler.WriteTextureCoord: invalid type %T", data))
	}
	index := WriteAccessor(doc, gltf.TargetArrayBuffer, data)
	doc.Accessors[index].Normalized = normalized
	return index
}

// WriteWeights adds a new WEIGHTS accessor to doc
// and fills the last buffer with the weights data.
// If success it returns the index of the new accessor.
func WriteWeights(doc *gltf.Document, data interface{}) uint32 {
	var normalized bool
	switch data.(type) {
	case [][4]uint8, [][4]uint16:
		normalized = true
	case [][4]float32:
	default:
		panic(fmt.Sprintf("modeler.WriteWeights: invalid type %T", data))
	}
	index := WriteAccessor(doc, gltf.TargetArrayBuffer, data)
	doc.Accessors[index].Normalized = normalized
	return index
}

// WriteJoints adds a new JOINTS accessor to doc
// and fills the last buffer with the joints data.
// If success it returns the index of the new accessor.
func WriteJoints(doc *gltf.Document, data interface{}) uint32 {
	switch data.(type) {
	case [][4]uint8, [][4]uint16:
	default:
		panic(fmt.Sprintf("modeler.WriteJoints: invalid type %T", data))
	}
	return WriteAccessor(doc, gltf.TargetArrayBuffer, data)
}

// WritePosition adds a new POSITION accessor to doc
// and fills the last buffer with the vertices data.
// If success it returns the index of the new accessor.
func WritePosition(doc *gltf.Document, data [][3]float32) uint32 {
	min := [3]float32{math.MaxFloat32, math.MaxFloat32, math.MaxFloat32}
	max := [3]float32{-math.MaxFloat32, -math.MaxFloat32, -math.MaxFloat32}
	for _, v := range data {
		for i, x := range v {
			min[i] = float32(math.Min(float64(min[i]), float64(x)))
			max[i] = float32(math.Max(float64(max[i]), float64(x)))
		}
	}
	index := WriteAccessor(doc, gltf.TargetArrayBuffer, data)
	doc.Accessors[index].Min = min[:]
	doc.Accessors[index].Max = max[:]
	return index
}

// WriteColor adds a new COLOR accessor to doc
// and fills the buffer with the color data.
// If success it returns the index of the new accessor.
func WriteColor(doc *gltf.Document, data interface{}) uint32 {
	var normalized bool
	switch data.(type) {
	case []color.RGBA, []color.RGBA64, [][4]uint8, [][3]uint8, [][4]uint16, [][3]uint16:
		normalized = true
	case [][3]float32, [][4]float32:
	default:
		panic(fmt.Sprintf("modeler.WriteColor: invalid type %T", data))
	}
	index := WriteAccessor(doc, gltf.TargetArrayBuffer, data)
	doc.Accessors[index].Normalized = normalized
	return index
}

// WriteImage adds a new image to doc
// and fills the buffer with the image data.
// If success it returns the index of the new image.
func WriteImage(doc *gltf.Document, name string, mimeType string, r io.Reader) (uint32, error) {
	var data []byte
	switch r := r.(type) {
	case *bytes.Buffer:
		data = r.Bytes()
	default:
		var err error
		data, err = ioutil.ReadAll(r)
		if err != nil {
			return 0, err
		}
	}
	index := WriteBufferView(doc, gltf.TargetNone, data)
	doc.Images = append(doc.Images, &gltf.Image{
		Name:       name,
		MimeType:   mimeType,
		BufferView: gltf.Index(index),
	})
	return uint32(len(doc.Images) - 1), nil
}

// WriteAccessor adds a new Accessor to doc
// and fills the buffer with data.
// If success it returns the index of the new accessor.
func WriteAccessor(doc *gltf.Document, target gltf.Target, data interface{}) uint32 {
	c, a, l := binary.Type(data)
	index := WriteBufferView(doc, target, data)
	doc.Accessors = append(doc.Accessors, &gltf.Accessor{
		BufferView:    gltf.Index(index),
		ByteOffset:    0,
		ComponentType: c,
		Type:          a,
		Count:         l,
	})
	return uint32(len(doc.Accessors) - 1)
}

// WriteBufferView adds a new BufferView to doc
// and fills the buffer with the data.
// If success it returns the index of the new buffer view.
func WriteBufferView(doc *gltf.Document, target gltf.Target, data interface{}) uint32 {
	c, a, l := binary.Type(data)
	sizeElement := binary.SizeOfElement(c, a)
	size := l * sizeElement
	buffer := lastBuffer(doc)
	offset := uint32(len(buffer.Data))
	padding := getPadding(offset, c.ByteSize())
	buffer.ByteLength += size + padding
	buffer.Data = append(buffer.Data, make([]byte, size+padding)...)
	// Cannot return error as the buffer has enough size and the data type is controlled.
	_ = binary.Write(buffer.Data[offset+padding:], 0, data)
	var stride uint32
	if target == gltf.TargetArrayBuffer && c.ByteSize()*a.Components() != sizeElement {
		stride = sizeElement
	}
	bufferView := &gltf.BufferView{
		Buffer:     uint32(len(doc.Buffers)) - 1,
		ByteLength: size,
		ByteOffset: offset + padding,
		ByteStride: stride,
		Target:     target,
	}
	doc.BufferViews = append(doc.BufferViews, bufferView)
	return uint32(len(doc.BufferViews)) - 1
}

func lastBuffer(doc *gltf.Document) *gltf.Buffer {
	if len(doc.Buffers) == 0 {
		doc.Buffers = append(doc.Buffers, new(gltf.Buffer))
	}
	return doc.Buffers[len(doc.Buffers)-1]
}

func getPadding(offset uint32, alignment uint32) uint32 {
	padAlign := offset % alignment
	if padAlign == 0 {
		return 0
	}
	return alignment - padAlign
}
