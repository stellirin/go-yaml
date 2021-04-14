package yaml_test

import (
	"reflect"
	"testing"

	. "czechia.dev/yaml"
)

func TestJSONToYAML(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "null",
			args: args{
				data: []byte(`{"a":null}`),
			},
			want: []byte("a: null\n"),
		},
		{
			name: "bool",
			args: args{
				data: []byte(`{"a":true,"b":false}`),
			},
			want: []byte("a: true\nb: false\n"),
		},
		{
			name: "string",
			args: args{
				data: []byte(`{"a":"x"}`),
			},
			want: []byte("a: x\n"),
		},
		{
			name: "int",
			args: args{
				data: []byte(`{"a":1}`),
			},
			want: []byte("a: 1\n"),
		},
		{
			name: "float",
			args: args{
				data: []byte(`{"a":3.4028234663852886e+38}`),
			},
			want: []byte("a: 3.4028234663852886e+38\n"),
		},
		{
			name: "timestamp",
			args: args{
				data: []byte(`{"a":"2001-12-15T02:59:43.1Z"}`),
			},
			want: []byte("a: \"2001-12-15T02:59:43.1Z\"\n"),
		},
		{
			name: "sequence",
			args: args{
				data: []byte(`{"a":["b","c"]}`),
			},
			want: []byte("a:\n  - b\n  - c\n"),
		},
		{
			name: "map",
			args: args{
				data: []byte(`{"a":[{"b":"x1","c":"y1"},{"b":"x2","c":"y2"}]}`),
			},
			want: []byte("a:\n  - b: x1\n    c: y1\n  - b: x2\n    c: y2\n"),
		},
		{
			name: "bad-format",
			args: args{
				data: []byte("[}"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := JSONToYAML(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("JSONToYAML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JSONToYAML() = '%v', want '%v'", string(got), string(tt.want))
			}
		})
	}
}
