package yaml_test

import (
	"reflect"
	"testing"

	. "czechia.dev/yaml"
)

func TestJSONToYAML(t *testing.T) {

	t.Run("null", func(t *testing.T) {
		t.Parallel()
		testJSONToYAML(t, test{
			name: "null",
			args: args{
				data: []byte(`{"a":null}`),
			},
			want: []byte("a: null\n"),
		})
	})

	t.Run("bool", func(t *testing.T) {
		t.Parallel()
		testJSONToYAML(t, test{
			name: "bool",
			args: args{
				data: []byte(`{"a":true,"b":false}`),
			},
			want: []byte("a: true\nb: false\n"),
		})
	})

	t.Run("string", func(t *testing.T) {
		t.Parallel()
		testJSONToYAML(t, test{
			name: "string",
			args: args{
				data: []byte(`{"a":"x"}`),
			},
			want: []byte("a: x\n"),
		})
	})

	t.Run("int", func(t *testing.T) {
		t.Parallel()
		testJSONToYAML(t, test{
			name: "int",
			args: args{
				data: []byte(`{"a":1}`),
			},
			want: []byte("a: 1\n"),
		})
	})

	t.Run("float", func(t *testing.T) {
		t.Parallel()
		testJSONToYAML(t, test{
			name: "float",
			args: args{
				data: []byte(`{"a":3.4028234663852886e+38}`),
			},
			want: []byte("a: 3.4028234663852886e+38\n"),
		})
	})

	t.Run("timestamp", func(t *testing.T) {
		t.Parallel()
		testJSONToYAML(t, test{
			name: "timestamp",
			args: args{
				data: []byte(`{"a":"2001-02-03T04:05:06.07Z"}`),
			},
			want: []byte("a: \"2001-02-03T04:05:06.07Z\"\n"),
		})
	})

	t.Run("sequence", func(t *testing.T) {
		t.Parallel()
		testJSONToYAML(t, test{
			name: "sequence",
			args: args{
				data: []byte(`{"a":["b","c"]}`),
			},
			want: []byte("a:\n  - b\n  - c\n"),
		})
	})

	t.Run("map", func(t *testing.T) {
		t.Parallel()
		testJSONToYAML(t, test{
			name: "map",
			args: args{
				data: []byte(`{"a":[{"b":"x1","c":"y1"},{"b":"x2","c":"y2"}]}`),
			},
			want: []byte("a:\n  - b: x1\n    c: y1\n  - b: x2\n    c: y2\n"),
		})
	})

	t.Run("bad-format", func(t *testing.T) {
		t.Parallel()
		testJSONToYAML(t, test{
			name: "bad-format",
			args: args{
				data: []byte("[}"),
			},
			wantErr: true,
		})
	})
}

func testJSONToYAML(t *testing.T, tt test) {
	got, err := JSONToYAML(tt.args.data)
	if (err != nil) != tt.wantErr {
		t.Errorf("JSONToYAML() error = %v, wantErr %v", err, tt.wantErr)
		return
	}
	if !reflect.DeepEqual(got, tt.want) {
		t.Errorf("JSONToYAML() = '%v', want '%v'", string(got), string(tt.want))
	}
}
