// Package modeler implements helper methods to write common glTF entities
// (indices, positions, colors, ...) into buffers.
package modeler

import (
	"bytes"
	"errors"
	"fmt"
	"image/color"
	"io"
	"io/ioutil"
	"math"
	"reflect"

	"github.com/qmuntal/gltf"
	"github.com/qmuntal/gltf/binary"
)

// WriteIndices adds a new INDICES accessor to doc
// and fills the last buffer with data.
// If success it returns the index of the new accessor.
func WriteIndices(doc *gltf.Document, data interface{}) uint32 {
	switch data.(type) {
	case []uint16, []uint32:
	default:
		panic(fmt.Sprintf("modeler.WriteIndices: invalid type %T", data))
	}
	return WriteAccessor(doc, gltf.TargetElementArrayBuffer, data)
}

// WriteNormal adds a new NORMAL accessor to doc
// and fills the last buffer with data.
// If success it returns the index of the new accessor.
func WriteNormal(doc *gltf.Document, data [][3]float32) uint32 {
	return WriteAccessor(doc, gltf.TargetArrayBuffer, data)
}

// WriteTangent adds a new TANGENT accessor to doc
// and fills the last buffer with data.
// If success it returns the index of the new accessor.
func WriteTangent(doc *gltf.Document, data [][4]float32) uint32 {
	return WriteAccessor(doc, gltf.TargetArrayBuffer, data)
}

// WriteTextureCoord adds a new TEXTURECOORD accessor to doc
// and fills the last buffer with data.
// If success it returns the index of the new accessor.
func WriteTextureCoord(doc *gltf.Document, data interface{}) uint32 {
	normalized := checkTextureCoord(data)
	index := WriteAccessor(doc, gltf.TargetArrayBuffer, data)
	doc.Accessors[index].Normalized = normalized
	return index
}

func checkTextureCoord(data interface{}) bool {
	var normalized bool
	switch data.(type) {
	case [][2]uint8, [][2]uint16:
		normalized = true
	case [][2]float32:
	default:
		panic(fmt.Sprintf("modeler.WriteTextureCoord: invalid type %T", data))
	}
	return normalized
}

// WriteWeights adds a new WEIGHTS accessor to doc
// and fills the last buffer with data.
// If success it returns the index of the new accessor.
func WriteWeights(doc *gltf.Document, data interface{}) uint32 {
	normalized := checkWeights(data)
	index := WriteAccessor(doc, gltf.TargetArrayBuffer, data)
	doc.Accessors[index].Normalized = normalized
	return index
}

func checkWeights(data interface{}) bool {
	var normalized bool
	switch data.(type) {
	case [][4]uint8, [][4]uint16:
		normalized = true
	case [][4]float32:
	default:
		panic(fmt.Sprintf("modeler.WriteWeights: invalid type %T", data))
	}
	return normalized
}

// WriteJoints adds a new JOINTS accessor to doc
// and fills the last buffer with data.
// If success it returns the index of the new accessor.
func WriteJoints(doc *gltf.Document, data interface{}) uint32 {
	checkJoints(data)
	return WriteAccessor(doc, gltf.TargetArrayBuffer, data)
}

func checkJoints(data interface{}) {
	switch data.(type) {
	case [][4]uint8, [][4]uint16:
	default:
		panic(fmt.Sprintf("modeler.WriteJoints: invalid type %T", data))
	}
}

// WritePosition adds a new POSITION accessor to doc
// and fills the last buffer with data.
// If success it returns the index of the new accessor.
func WritePosition(doc *gltf.Document, data [][3]float32) uint32 {
	index := WriteAccessor(doc, gltf.TargetArrayBuffer, data)
	min, max := minMaxFloat32(data)
	doc.Accessors[index].Min = min[:]
	doc.Accessors[index].Max = max[:]
	return index
}

func minMaxFloat32(data [][3]float32) ([3]float64, [3]float64) {
	min := [3]float64{math.MaxFloat64, math.MaxFloat64, math.MaxFloat64}
	max := [3]float64{-math.MaxFloat64, -math.MaxFloat64, -math.MaxFloat64}
	for _, v := range data {
		for i, x := range v {
			min[i] = math.Min(float64(min[i]), float64(x))
			max[i] = math.Max(float64(max[i]), float64(x))
		}
	}
	return min, max
}

// WriteColor adds a new COLOR accessor to doc
// and fills the buffer with data.
// If success it returns the index of the new accessor.
func WriteColor(doc *gltf.Document, data interface{}) uint32 {
	normalized := checkColor(data)
	index := WriteAccessor(doc, gltf.TargetArrayBuffer, data)
	doc.Accessors[index].Normalized = normalized
	return index
}

func checkColor(data interface{}) bool {
	var normalized bool
	switch data.(type) {
	case []color.RGBA, []color.RGBA64, [][4]uint8, [][3]uint8, [][4]uint16, [][3]uint16:
		normalized = true
	case [][3]float32, [][4]float32:
	default:
		panic(fmt.Sprintf("modeler.WriteColor: invalid type %T", data))
	}
	return normalized
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
// and fills the buffer with the data.
// Returns the index of the new accessor.
func WriteAccessor(doc *gltf.Document, target gltf.Target, data interface{}) uint32 {
	ensurePadding(doc)
	index := WriteBufferView(doc, target, data)
	c, a, l := binary.Type(data)
	doc.Accessors = append(doc.Accessors, &gltf.Accessor{
		BufferView:    gltf.Index(index),
		ByteOffset:    0,
		ComponentType: c,
		Type:          a,
		Count:         l,
	})
	return uint32(len(doc.Accessors) - 1)
}

// WriteAccessorsInterleaved adds as many accessors as
// elements in data all pointing to the same interleaved buffer view
// and fills the buffer with the data.
// Returns an slice with the indices of the newly created accessors,
// with the same order as data or an error if the data elements
// don´t have all the same length.
func WriteAccessorsInterleaved(doc *gltf.Document, data ...interface{}) ([]uint32, error) {
	ensurePadding(doc)
	index, err := WriteBufferViewInterleaved(doc, data...)
	if err != nil {
		return nil, err
	}
	indices := make([]uint32, len(data))
	var byteOffset uint32
	for i, d := range data {
		c, t, l := binary.Type(d)
		doc.Accessors = append(doc.Accessors, &gltf.Accessor{
			BufferView:    gltf.Index(index),
			ByteOffset:    byteOffset,
			ComponentType: c,
			Type:          t,
			Count:         l,
		})
		byteOffset += gltf.SizeOfElement(c, t)
		indices[i] = uint32(len(doc.Accessors) - 1)
	}
	return indices, nil
}

// CustomAttribute defines an application-specific attribute
type CustomAttribute struct {
	Name string
	Data interface{}
}

// Attributes defines all the vertex attributes that can
// be associated to a primitive.
type Attributes struct {
	Position [][3]float32
	Normal   [][3]float32
	Tangent  [][4]float32
	// [][2]uint8, [][2]uint16 or [][2]float32
	TextureCoord_0, TextureCoord_1 interface{}
	// [][4]uint8, [][4]uint16 or [][4]float32
	Weights interface{}
	// [][4]uint8 or [][4]uint16
	Joints interface{}
	//[]color.RGBA, []color.RGBA64, [][4]uint8, [][3]uint8, [][4]uint16, [][3]uint16, [][3]float32 or [][4]float32
	Color            interface{}
	CustomAttributes []CustomAttribute
}

// WriteAttributesInterleaved write all the attributes in v
// which are not nil and have a non-zero length.
// Returns an attribute map that can be directly used
// as a primitive attributes.
func WriteAttributesInterleaved(doc *gltf.Document, v Attributes) (map[string]uint32, error) {
	type attrProps struct {
		Name       string
		Normalized bool
	}
	var (
		props []attrProps
		data  []interface{}
	)
	if len(v.Position) != 0 {
		props = append(props, attrProps{Name: gltf.POSITION})
		data = append(data, v.Position)
	}
	if len(v.Normal) != 0 {
		props = append(props, attrProps{Name: gltf.NORMAL})
		data = append(data, v.Normal)
	}
	if len(v.Tangent) != 0 {
		props = append(props, attrProps{Name: gltf.TANGENT})
		data = append(data, v.Tangent)
	}
	if sliceLength(v.TextureCoord_0) != 0 {
		normalized := checkTextureCoord(v.TextureCoord_0)
		props = append(props, attrProps{Name: gltf.TEXCOORD_0, Normalized: normalized})
		data = append(data, v.TextureCoord_0)
	}
	if sliceLength(v.TextureCoord_1) != 0 {
		normalized := checkTextureCoord(v.TextureCoord_1)
		props = append(props, attrProps{Name: gltf.TEXCOORD_1, Normalized: normalized})
		data = append(data, v.TextureCoord_1)
	}
	if sliceLength(v.Weights) != 0 {
		normalized := checkWeights(v.Weights)
		props = append(props, attrProps{Name: gltf.WEIGHTS_0, Normalized: normalized})
		data = append(data, v.Weights)
	}
	if sliceLength(v.Joints) != 0 {
		checkJoints(v.Joints)
		props = append(props, attrProps{Name: gltf.JOINTS_0})
		data = append(data, v.Joints)
	}
	if sliceLength(v.Color) != 0 {
		normalized := checkColor(v.Color)
		props = append(props, attrProps{Name: gltf.COLOR_0, Normalized: normalized})
		data = append(data, v.Color)
	}
	for _, c := range v.CustomAttributes {
		if sliceLength(c.Data) != 0 {
			props = append(props, attrProps{Name: c.Name})
			data = append(data, c.Data)
		}
	}
	indices, err := WriteAccessorsInterleaved(doc, data...)
	if err != nil {
		return nil, err
	}
	attrs := make(map[string]uint32, len(props))
	for i, index := range indices {
		prop := props[i]
		attrs[prop.Name] = index
		doc.Accessors[index].Normalized = prop.Normalized
	}
	if pos, ok := attrs[gltf.POSITION]; ok {
		min, max := minMaxFloat32(v.Position)
		doc.Accessors[pos].Min = min[:]
		doc.Accessors[pos].Max = max[:]
	}
	return attrs, nil
}

// WriteBufferViewInterleaved adds a new BufferView to doc
// and fills the buffer with one or more vertex attribute.
// If success it returns the index of the new buffer view.
// Returns the index of the new buffer view or an error if the data elements
// don´t have all the same length.
func WriteBufferViewInterleaved(doc *gltf.Document, data ...interface{}) (uint32, error) {
	return writeBufferViews(doc, gltf.TargetArrayBuffer, data...)
}

// WriteBufferView adds a new BufferView to doc
// and fills the buffer with the data.
// Returns the index of the new buffer view.
func WriteBufferView(doc *gltf.Document, target gltf.Target, data interface{}) uint32 {
	index, _ := writeBufferViews(doc, target, data)
	return index
}

func writeBufferViews(doc *gltf.Document, target gltf.Target, data ...interface{}) (uint32, error) {
	var refLength, stride, size uint32
	for i, d := range data {
		c, a, l := binary.Type(d)
		if i == 0 {
			refLength = l
		} else if refLength != l {
			return 0, errors.New("go3mf: interleaved data shall have the same number of elements in all chunks")
		}
		sizeOfElement := gltf.SizeOfElement(c, a)
		size += l * sizeOfElement
		if len(data) > 1 {
			stride += sizeOfElement
		} else if target == gltf.TargetArrayBuffer && c.ByteSize()*a.Components() != sizeOfElement {
			stride = sizeOfElement
		}
	}
	buffer := lastBuffer(doc)
	offset := uint32(len(buffer.Data))
	buffer.ByteLength += size
	buffer.Data = append(buffer.Data, make([]byte, size)...)
	dataOffset := offset
	for _, d := range data {
		// Cannot return error as the buffer has enough size and the data type is controlled.
		_ = binary.Write(buffer.Data[dataOffset:], stride, d)
		c, a, _ := binary.Type(d)
		dataOffset += gltf.SizeOfElement(c, a)
	}
	bufferView := &gltf.BufferView{
		Buffer:     uint32(len(doc.Buffers)) - 1,
		ByteLength: size,
		ByteOffset: offset,
		ByteStride: stride,
		Target:     target,
	}
	doc.BufferViews = append(doc.BufferViews, bufferView)
	return uint32(len(doc.BufferViews)) - 1, nil
}

func ensurePadding(doc *gltf.Document) {
	buffer := lastBuffer(doc)
	padding := getPadding(uint32(len(buffer.Data)))
	buffer.Data = append(buffer.Data, make([]byte, padding)...)
	buffer.ByteLength += padding
}

func lastBuffer(doc *gltf.Document) *gltf.Buffer {
	if len(doc.Buffers) == 0 {
		doc.Buffers = append(doc.Buffers, new(gltf.Buffer))
	}
	return doc.Buffers[len(doc.Buffers)-1]
}

func getPadding(offset uint32) uint32 {
	padAlign := offset % 4
	if padAlign == 0 {
		return 0
	}
	return 4 - padAlign
}

func sliceLength(data interface{}) int {
	if data == nil {
		return 0
	}
	v := reflect.ValueOf(data)
	if v.IsNil() {
		return 0
	}
	if v.Kind() != reflect.Slice {
		panic(fmt.Sprintf("go3mf: expecting a slice but got %s", v.Kind()))
	}
	return v.Len()
}
