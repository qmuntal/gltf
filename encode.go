package gltf

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
)

// WriteHandler is the interface that wraps the Write method.
//
// WriteResource should behaves as io.Write in terms of reading the writing resource.
type WriteHandler interface {
	WriteResource(uri string, data []byte) error
}

// Save will save a document as a glTF with the specified by name.
func Save(doc *Document, name string) error {
	return save(doc, name, false)
}

// SaveBinary will save a document as a GLB file with the specified by name.
func SaveBinary(doc *Document, name string) error {
	return save(doc, name, true)
}

func save(doc *Document, name string, asBinary bool) error {
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	e := NewEncoder(f).WithWriteHandler(&RelativeFileHandler{Dir: filepath.Dir(name)})
	e.AsBinary = asBinary
	if err := e.Encode(doc); err != nil {
		f.Close()
		return err
	}
	return f.Close()
}

// An Encoder writes a GLTF to an output stream
// with relative external buffers support.
type Encoder struct {
	AsBinary     bool
	WriteHandler WriteHandler
	w            io.Writer
}

// NewEncoder returns a new encoder that writes to w as a normal glTF file.
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		AsBinary:     true,
		WriteHandler: new(RelativeFileHandler),
		w:            w,
	}
}

// WithWriteHandler sets the WriteHandler.
func (e *Encoder) WithWriteHandler(h WriteHandler) *Encoder {
	e.WriteHandler = h
	return e
}

// Encode writes the encoding of doc to the stream.
func (e *Encoder) Encode(doc *Document) error {
	if doc.Asset.Version == "" {
		doc.Asset.Version = "2.0"
	}
	var err error
	var externalBufferIndex = 0
	if e.AsBinary {
		var hasBinChunk bool
		hasBinChunk, err = e.encodeBinary(doc)
		if hasBinChunk {
			externalBufferIndex = 1
		}
	} else {
		err = json.NewEncoder(e.w).Encode(doc)
	}
	if err != nil {
		return err
	}

	for i := externalBufferIndex; i < len(doc.Buffers); i++ {
		buffer := doc.Buffers[i]
		if len(buffer.Data) == 0 || buffer.IsEmbeddedResource() {
			continue
		}
		if err = e.encodeBuffer(buffer); err != nil {
			return err
		}
	}

	return err
}

func (e *Encoder) encodeBuffer(buffer *Buffer) error {
	if err := validateBufferURI(buffer.URI); err != nil {
		return err
	}

	return e.WriteHandler.WriteResource(buffer.URI, buffer.Data)
}

func (e *Encoder) encodeBinary(doc *Document) (bool, error) {
	jsonText, err := json.Marshal(doc)
	if err != nil {
		return false, err
	}
	jsonHeader := chunkHeader{
		Length: uint32(((len(jsonText) + 3) / 4) * 4),
		Type:   glbChunkJSON,
	}
	header := glbHeader{
		Magic:      glbHeaderMagic,
		Version:    2,
		Length:     12 + 8 + jsonHeader.Length, // 12-byte glb header + 8-byte json chunk header
		JSONHeader: jsonHeader,
	}
	headerPadding := make([]byte, header.JSONHeader.Length-uint32(len(jsonText)))
	for i := range headerPadding {
		headerPadding[i] = ' '
	}
	err = binary.Write(e.w, binary.LittleEndian, &header)
	if err != nil {
		return false, err
	}
	e.w.Write(jsonText)
	e.w.Write(headerPadding)

	hasBinChunk := len(doc.Buffers) > 0 && doc.Buffers[0].URI == ""
	if hasBinChunk {
		var binBufferLength uint32
		var binBuffer *Buffer
		if len(doc.Buffers) > 0 {
			binBuffer = doc.Buffers[0]
			binBufferLength = binBuffer.ByteLength
		}
		binPaddedLength := ((binBufferLength + 3) / 4) * 4
		binPadding := make([]byte, binPaddedLength-binBufferLength)
		binHeader := chunkHeader{Length: 0, Type: glbChunkBIN}
		binHeader.Length = binPaddedLength
		header.Length = uint32(8) + binHeader.Length
		for i := range binPadding {
			binPadding[i] = 0
		}
		binary.Write(e.w, binary.LittleEndian, &binHeader)
		e.w.Write(binBuffer.Data)
		_, err = e.w.Write(binPadding)
	}

	return hasBinChunk, err
}

// UnmarshalJSON unmarshal the node with the correct default values.
func (n *Node) UnmarshalJSON(data []byte) error {
	type alias Node
	tmp := alias(Node{
		Matrix:   DefaultMatrix,
		Rotation: DefaultRotation,
		Scale:    DefaultScale,
	})
	err := json.Unmarshal(data, &tmp)
	if err == nil {
		*n = Node(tmp)
	}
	return err
}

// MarshalJSON marshal the node with the correct default values.
func (n *Node) MarshalJSON() ([]byte, error) {
	type alias Node
	out, err := json.Marshal(&struct{ *alias }{alias: (*alias)(n)})
	if err == nil {
		if n.Matrix == DefaultMatrix {
			out = removeProperty([]byte(`"matrix":[1,0,0,0,0,1,0,0,0,0,1,0,0,0,0,1]`), out)
		} else if n.Matrix == emptyMatrix {
			out = removeProperty([]byte(`"matrix":[0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0]`), out)
		}

		if n.Rotation == DefaultRotation {
			out = removeProperty([]byte(`"rotation":[0,0,0,1]`), out)
		} else if n.Rotation == emptyRotation {
			out = removeProperty([]byte(`"rotation":[0,0,0,0]`), out)
		}

		if n.Scale == DefaultScale {
			out = removeProperty([]byte(`"scale":[1,1,1]`), out)
		} else if n.Scale == emptyScale {
			out = removeProperty([]byte(`"scale":[0,0,0]`), out)
		}

		if n.Translation == DefaultTranslation {
			out = removeProperty([]byte(`"translation":[0,0,0]`), out)
		}
		out = sanitizeJSON(out)
	}
	return out, err
}

// MarshalJSON marshal the camera with the correct default values.
func (c *Camera) MarshalJSON() ([]byte, error) {
	type alias Camera
	if c.Perspective != nil {
		return json.Marshal(&struct {
			Type string `json:"type"`
			*alias
		}{
			Type:  "perspective",
			alias: (*alias)(c),
		})
	} else if c.Orthographic != nil {
		return json.Marshal(&struct {
			Type string `json:"type"`
			*alias
		}{
			Type:  "orthographic",
			alias: (*alias)(c),
		})
	}
	return nil, errors.New("gltf: camera must defined either the perspective or orthographic property")
}

// UnmarshalJSON unmarshal the material with the correct default values.
func (m *Material) UnmarshalJSON(data []byte) error {
	type alias Material
	tmp := alias(Material{AlphaCutoff: Float(0.5)})
	err := json.Unmarshal(data, &tmp)
	if err == nil {
		*m = Material(tmp)
	}
	return err
}

// MarshalJSON marshal the material with the correct default values.
func (m *Material) MarshalJSON() ([]byte, error) {
	type alias Material
	out, err := json.Marshal(&struct{ *alias }{alias: (*alias)(m)})
	if err == nil {
		if m.AlphaCutoff != nil && *m.AlphaCutoff == 0.5 {
			out = removeProperty([]byte(`"alphaCutoff":0.5`), out)
		}
		if m.EmissiveFactor == [3]float32{0, 0, 0} {
			out = removeProperty([]byte(`"emissiveFactor":[0,0,0]`), out)
		}
		out = sanitizeJSON(out)
	}
	return out, err
}

// UnmarshalJSON unmarshal the texture info with the correct default values.
func (n *NormalTexture) UnmarshalJSON(data []byte) error {
	type alias NormalTexture
	tmp := alias(NormalTexture{Scale: Float(1)})
	err := json.Unmarshal(data, &tmp)
	if err == nil {
		*n = NormalTexture(tmp)
	}
	return err
}

// MarshalJSON marshal the texture info with the correct default values.
func (n *NormalTexture) MarshalJSON() ([]byte, error) {
	type alias NormalTexture
	if n.Scale != nil && *n.Scale == 1 {
		return json.Marshal(&struct {
			Scale float32 `json:"scale,omitempty"`
			*alias
		}{
			Scale: 0,
			alias: (*alias)(n),
		})
	}
	return json.Marshal(&struct{ *alias }{alias: (*alias)(n)})
}

// UnmarshalJSON unmarshal the texture info with the correct default values.
func (o *OcclusionTexture) UnmarshalJSON(data []byte) error {
	type alias OcclusionTexture
	tmp := alias(OcclusionTexture{Strength: Float(1)})
	err := json.Unmarshal(data, &tmp)
	if err == nil {
		*o = OcclusionTexture(tmp)
	}
	return err
}

// MarshalJSON marshal the texture info with the correct default values.
func (o *OcclusionTexture) MarshalJSON() ([]byte, error) {
	type alias OcclusionTexture
	if o.Strength != nil && *o.Strength == 1 {
		return json.Marshal(&struct {
			Strength float32 `json:"strength,omitempty"`
			*alias
		}{
			Strength: 0,
			alias:    (*alias)(o),
		})
	}
	return json.Marshal(&struct{ *alias }{alias: (*alias)(o)})
}

// UnmarshalJSON unmarshal the pbr with the correct default values.
func (p *PBRMetallicRoughness) UnmarshalJSON(data []byte) error {
	type alias PBRMetallicRoughness
	tmp := alias(PBRMetallicRoughness{BaseColorFactor: &[4]float32{1, 1, 1, 1}, MetallicFactor: Float(1), RoughnessFactor: Float(1)})
	err := json.Unmarshal(data, &tmp)
	if err == nil {
		*p = PBRMetallicRoughness(tmp)
	}
	return err
}

// MarshalJSON marshal the pbr with the correct default values.
func (p *PBRMetallicRoughness) MarshalJSON() ([]byte, error) {
	type alias PBRMetallicRoughness
	out, err := json.Marshal(&struct{ *alias }{alias: (*alias)(p)})
	if err == nil {
		if p.MetallicFactor != nil && *p.MetallicFactor == 1 {
			out = removeProperty([]byte(`"metallicFactor":1`), out)
		}
		if p.RoughnessFactor != nil && *p.RoughnessFactor == 1 {
			out = removeProperty([]byte(`"roughnessFactor":1`), out)
		}
		if p.BaseColorFactor != nil && *p.BaseColorFactor == [4]float32{1, 1, 1, 1} {
			out = removeProperty([]byte(`"baseColorFactor":[1,1,1,1]`), out)
		}
		out = sanitizeJSON(out)
	}
	return out, err
}

// UnmarshalJSON unmarshal the extensions with the supported extensions initialized.
func (ext *Extensions) UnmarshalJSON(data []byte) error {
	if len(*ext) == 0 {
		*ext = make(Extensions)
	}
	var raw map[string]json.RawMessage
	err := json.Unmarshal(data, &raw)
	if err == nil {
		for key, value := range raw {
			if extFactory, ok := queryExtension(key); ok {
				n, err := extFactory(value)
				if err != nil {
					(*ext)[key] = value
				} else {
					(*ext)[key] = n
				}
			} else {
				(*ext)[key] = value
			}
		}
	}

	return err
}

func removeProperty(str []byte, b []byte) []byte {
	b = bytes.Replace(b, str, []byte(""), 1)
	return bytes.Replace(b, []byte(`,,`), []byte(","), 1)
}

func sanitizeJSON(b []byte) []byte {
	b = bytes.Replace(b, []byte(`{,`), []byte("{"), 1)
	return bytes.Replace(b, []byte(`,}`), []byte("}"), 1)
}
