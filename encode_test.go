package gltf

import (
	"bytes"
	"io"
	"io/ioutil"
	"testing"

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
		chunks[uri] = bytes.NewBuffer(make([]byte, size))
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := saveMemory(tt.args.doc, true)
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
