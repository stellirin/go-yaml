package yaml_test

import (
	"reflect"
	"testing"

	. "czechia.dev/yaml"
)

func TestYAMLToJSON(t *testing.T) {
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
				data: []byte("a: nullTag\nb: null\n"),
			},
			want: []byte(`{"a":"nullTag","b":null}`),
		},
		{
			name: "bool",
			args: args{
				data: []byte("a: boolTag\nb: true\nc: false\n"),
			},
			want: []byte(`{"a":"boolTag","b":true,"c":false}`),
		},
		{
			name: "string",
			args: args{
				data: []byte("a: strTag\nb: 1\n"),
			},
			want: []byte(`{"a":"strTag","b":1}`),
		},
		{
			name: "float",
			args: args{
				// strconv.FormatFloat(math.MaxFloat32, 'g', -1, 32) = 3.4028235e+38
				// strconv.FormatFloat(math.MaxFloat32, 'g', -1, 64) = 3.4028234663852886e+38
				// strconv.FormatFloat(math.MaxFloat64, 'g', -1, 64) = 1.7976931348623157e+308
				data: []byte("a: floatTag\nb: 3.4028235e+38\nc: 3.4028234663852886e+38\n"),
			},
			want: []byte(`{"a":"floatTag","b":3.4028234663852886e+38,"c":3.4028234663852886e+38}`),
		},
		{
			name: "timestamp",
			args: args{
				data: []byte("a: timestampTag\nb: 2001-12-15T02:59:43.1Z\n"),
			},
			want: []byte(`{"a":"timestampTag","b":"2001-12-15T02:59:43.1Z"}`),
		},
		{
			name: "sequence",
			args: args{
				data: []byte("a:\n- b\n- c\n"),
			},
			want: []byte(`{"a":["b","c"]}`),
		},
		{
			name: "map",
			args: args{
				data: []byte("a:\n- b: x1\n  c: y1\n- b: x2\n  c: y2\n"),
			},
			want: []byte(`{"a":[{"b":"x1","c":"y1"},{"b":"x2","c":"y2"}]}`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := YAMLToJSON(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Marshaler.YAMLToJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Marshaler.YAMLToJSON() = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}
