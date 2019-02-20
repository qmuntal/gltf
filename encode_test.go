package gltf

import (
	"bytes"
	"io"
	"io/ioutil"
	"testing"
	"fmt"

	"github.com/go-test/deep"
)

type writeCloser struct {
	io.Writer
}

func (w *writeCloser) Close() error { return nil }

func saveMemory(doc *Document, asBinary bool) (*Decoder, error) {
	buff := new(bytes.Buffer)
	chunks := make(map[string]*bytes.Buffer)
	wcb := func(uri string, size int) (io.WriteCloser, error) {
		chunks[uri] = bytes.NewBuffer(make([]byte, 0, size))
		return &writeCloser{chunks[uri]}, nil
	}
	if err := NewEncoder(buff, wcb, asBinary).Encode(doc); err != nil {
		return nil, err
	}
	rcb := func(uri string) (io.ReadCloser, error) {
		return ioutil.NopCloser(chunks[uri]), nil
	}
	return NewDecoder(buff, rcb), nil
}

func TestEncoder_Encode(t *testing.T) {
	type args struct {
		doc *Document
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"empty", args{&Document{}}, false},
		{"withExtensions", args{&Document{Extensions: 8.0, Extras: map[string]interface{}{"a": "b"}, ExtensionsUsed: []string{"c"}, ExtensionsRequired: []string{"d", "e"}}}, false},
		{"withAsset", args{&Document{Asset: Asset{Extensions: 8.0, Extras: map[string]interface{}{"a": "b"}, Copyright: "@2019", Generator: "qmuntal/gltf", Version: "2.0", MinVersion: "1.0"}}}, false},
		{"withAccessors", args{&Document{Accessors: []Accessor{
			{Extensions: 8.0, Extras: map[string]interface{}{"a": "b"}, Name: "acc_1", BufferView: 0, ByteOffset: 50, ComponentType: Byte, Normalized: true, Count: 5, Type: Vec3, Max: []float32{1, 2}, Min: []float32{2.4}},
			{BufferView: 0, Normalized: false, Count: 50, Type: Vec4, Sparse: &Sparse{Extensions: 8.0, Extras: map[string]interface{}{"a": "b"}, Count: 2,
				Values:  SparseValues{Extensions: 8.0, Extras: map[string]interface{}{"a": "b"}, BufferView: 1, ByteOffset: 2},
				Indices: SparseIndices{Extensions: 8.0, Extras: map[string]interface{}{"a": "b"}, BufferView: 1, ByteOffset: 2, ComponentType: Float}},
			},
		}}}, false},
		{"withAnimations", args{&Document{Animations: []Animation{
			{Extensions: 8.0, Extras: map[string]interface{}{"a": "b"}, Name: "an_1", Channels: []Channel{
				{Extensions: 8.0, Extras: map[string]interface{}{"a": "b"}, Sampler: 1, Target: ChannelTarget{Extensions: 8.0, Extras: map[string]interface{}{"a": "b"}, Node: 10, Path: Rotation}},
				{Extensions: 8.0, Extras: map[string]interface{}{"a": "b"}, Sampler: 2, Target: ChannelTarget{Extensions: 8.0, Extras: map[string]interface{}{"a": "b"}, Node: 10, Path: Scale}},
			}},
			{Extensions: 8.0, Extras: map[string]interface{}{"a": "b"}, Name: "an_2", Channels: []Channel{
				{Extensions: 8.0, Extras: map[string]interface{}{"a": "b"}, Sampler: 1, Target: ChannelTarget{Extensions: 8.0, Extras: map[string]interface{}{"a": "b"}, Node: 3, Path: Weights}},
				{Extensions: 8.0, Extras: map[string]interface{}{"a": "b"}, Sampler: 2, Target: ChannelTarget{Extensions: 8.0, Extras: map[string]interface{}{"a": "b"}, Node: 5, Path: Translation}},
			}},
		}}}, false},
		{"withBuffer", args{&Document{Buffers: []Buffer{
			{Extensions: 8.0, Extras: map[string]interface{}{"a": "b"}, Name: "binary", ByteLength: 3, URI: "a.bin", Data: []uint8{1,2,3}},
			{Extensions: 8.0, Extras: map[string]interface{}{"a": "b"}, Name: "embedded", ByteLength: 2, URI: "data:application/octet-stream;base64,YW55ICsgb2xkICYgZGF0YQ==", Data: []byte("any + old & data")},
			{Extensions: 8.0, Extras: map[string]interface{}{"a": "b"}, Name: "external", ByteLength: 4, URI: "b.bin", Data: []uint8{4,5,6,7}},
		}}}, false},
	}
	for _, tt := range tests {
		for _, method := range []string{"json", "binary"} {
			t.Run(fmt.Sprintf("%s_%s", tt.name, method), func(t *testing.T) {
				var asBinary bool
				if method == "binary" {
					asBinary = true
					for i := 1; i < len(tt.args.doc.Buffers); i++ {
						tt.args.doc.Buffers[i].EmbeddedResource()
					}
				}
				d, err := saveMemory(tt.args.doc, asBinary)
				if (err != nil) != tt.wantErr {
					t.Errorf("Encoder.Encode() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				doc := new(Document)
				d.Decode(doc)
				if diff := deep.Equal(doc, tt.args.doc); diff != nil {
					t.Errorf("Encoder.Encode() = %v", diff)
					return
				}
			})
		}
	}
}
