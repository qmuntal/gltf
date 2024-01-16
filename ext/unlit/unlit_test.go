package unlit

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestUnlit_UnmarshalJSON(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *Unlit
		wantErr bool
	}{
		{"err", args{[]byte(`{"fake": {{"a":2}}`)}, &Unlit{}, true},
		{"empty", args{[]byte(`{}`)}, &Unlit{}, false},
		{"withProps", args{[]byte(`{"fake": {"a":2}}`)}, &Unlit{
			"fake": json.RawMessage(`{"a":2}`),
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := Unlit{}
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
		{"base", args{[]byte("{}")}, &Unlit{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Unmarshal(tt.args.data)
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
