package unlit_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/qmuntal/gltf/ext/unlit"
)

func TestUnlit_UnmarshalJSON(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *unlit.Unlit
		wantErr bool
	}{
		{"err", args{[]byte(`{"fake": {{"a":2}}`)}, &unlit.Unlit{}, true},
		{"empty", args{[]byte(`{}`)}, &unlit.Unlit{}, false},
		{"withProps", args{[]byte(`{"fake": {"a":2}}`)}, &unlit.Unlit{
			"fake": json.RawMessage(`{"a":2}`),
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := unlit.Unlit{}
			if err := u.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Unlit.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(&u, tt.want) {
				t.Errorf("Unlit.UnmarshalJSON() = %v, want %v", &u, tt.want)
			}
		})
	}
}

func TestUnmarshal(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    any
		wantErr bool
	}{
		{"base", args{[]byte("{}")}, &unlit.Unlit{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := unlit.Unmarshal(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Unmarshal() = %v, want %v", got, tt.want)
			}
		})
	}
}
