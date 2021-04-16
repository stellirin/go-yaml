package yaml_test

import (
	"reflect"
	"testing"

	. "czechia.dev/yaml"
)

func TestYAMLToJSON(t *testing.T) {

	t.Run("null", func(t *testing.T) {
		t.Parallel()
		testYAMLToJSON(t, test{
			name: "null",
			args: args{
				data: []byte("a: nullTag\nb: null\n"),
			},
			want: []byte(`{"a":"nullTag","b":null}`),
		})
	})

	t.Run("bool", func(t *testing.T) {
		t.Parallel()
		testYAMLToJSON(t, test{
			name: "bool",
			args: args{
				data: []byte("a: boolTag\nb: true\nc: false\n"),
			},
			want: []byte(`{"a":"boolTag","b":true,"c":false}`),
		})
	})

	t.Run("string", func(t *testing.T) {
		t.Parallel()
		testYAMLToJSON(t, test{
			name: "string",
			args: args{
				data: []byte("a: strTag\nb: \"x\"\n"),
			},
			want: []byte(`{"a":"strTag","b":"x"}`),
		})
	})

	t.Run("int", func(t *testing.T) {
		t.Parallel()
		testYAMLToJSON(t, test{
			name: "int",
			args: args{
				data: []byte("a: strTag\nb: 1\n"),
			},
			want: []byte(`{"a":"strTag","b":1}`),
		})
	})

	t.Run("float", func(t *testing.T) {
		t.Parallel()
		testYAMLToJSON(t, test{
			name: "float",
			args: args{
				// strconv.FormatFloat(math.MaxFloat32, 'g', -1, 32) = 3.4028235e+38
				// strconv.FormatFloat(math.MaxFloat32, 'g', -1, 64) = 3.4028234663852886e+38
				// strconv.FormatFloat(math.MaxFloat64, 'g', -1, 64) = 1.7976931348623157e+308
				data: []byte("a: floatTag\nb: 3.4028235e+38\nc: 3.4028234663852886e+38\n"),
			},
			want: []byte(`{"a":"floatTag","b":3.4028234663852886e+38,"c":3.4028234663852886e+38}`),
		})
	})

	t.Run("timestamp", func(t *testing.T) {
		t.Parallel()
		testYAMLToJSON(t, test{
			name: "timestamp",
			args: args{
				data: []byte("a: timestampTag\nb: 2001-02-03T04:05:06.07Z\n"),
			},
			want: []byte(`{"a":"timestampTag","b":"2001-02-03T04:05:06.07Z"}`),
		})
	})

	t.Run("sequence", func(t *testing.T) {
		t.Parallel()
		testYAMLToJSON(t, test{
			name: "sequence",
			args: args{
				data: []byte("a:\n- b\n- c\n"),
			},
			want: []byte(`{"a":["b","c"]}`),
		})
	})

	t.Run("map-bool", func(t *testing.T) {
		t.Parallel()
		testYAMLToJSON(t, test{
			name: "map-bool",
			args: args{
				data: []byte("true:\n- b: x1\n  c: y1\n- b: x2\n  c: y2\n"),
			},
			want: []byte(`{"true":[{"b":"x1","c":"y1"},{"b":"x2","c":"y2"}]}`),
		})
	})

	t.Run("map-string", func(t *testing.T) {
		t.Parallel()
		testYAMLToJSON(t, test{
			name: "map-string",
			args: args{
				data: []byte("a:\n- b: x1\n  c: y1\n- b: x2\n  c: y2\n"),
			},
			want: []byte(`{"a":[{"b":"x1","c":"y1"},{"b":"x2","c":"y2"}]}`),
		})
	})

	t.Run("map-int", func(t *testing.T) {
		t.Parallel()
		testYAMLToJSON(t, test{
			name: "map-int",
			args: args{
				data: []byte("1:\n- b: x1\n  c: y1\n- b: x2\n  c: y2\n"),
			},
			want: []byte(`{"1":[{"b":"x1","c":"y1"},{"b":"x2","c":"y2"}]}`),
		})
	})

	t.Run("map-float", func(t *testing.T) {
		t.Parallel()
		testYAMLToJSON(t, test{
			name: "map-float",
			args: args{
				data: []byte("3.4028235e+38:\n- b: x1\n  c: y1\n- b: x2\n  c: y2\n"),
			},
			want: []byte(`{"3.4028235e+38":[{"b":"x1","c":"y1"},{"b":"x2","c":"y2"}]}`),
		})
	})

	t.Run("map-timestamp", func(t *testing.T) {
		t.Parallel()
		testYAMLToJSON(t, test{
			name: "map-int",
			args: args{
				data: []byte("2001-02-03T04:05:06.07Z:\n- b: x1\n  c: y1\n- b: x2\n  c: y2\n"),
			},
			want: []byte(`{"2001-02-03T04:05:06.07Z":[{"b":"x1","c":"y1"},{"b":"x2","c":"y2"}]}`),
		})
	})

	t.Run("template", func(t *testing.T) {
		t.Parallel()
		testYAMLToJSON(t, test{
			name: "template",
			args: args{
				data: []byte("a1: &t\n  b: x\n  c: y1\na2:\n  <<: *t\n  c: y2"),
			},
			want: []byte(`{"a1":{"b":"x","c":"y1"},"a2":{"b":"x","c":"y2"}}`),
		})
	})

	t.Run("multiline-literal", func(t *testing.T) {
		t.Parallel()
		testYAMLToJSON(t, test{
			name: "multiline-literal",
			args: args{
				data: []byte("a:\n  - x\n  - |\n    y1\n    y2\n  - z\n"),
			},
			want: []byte(`{"a":["x","y1\ny2\n","z"]}`),
		})
	})

	t.Run("multiline-literal-chomp", func(t *testing.T) {
		t.Parallel()
		testYAMLToJSON(t, test{
			name: "multiline-literal",
			args: args{
				data: []byte("a:\n  - x\n  - |-\n    y1\n    y2\n  - z\n"),
			},
			want: []byte(`{"a":["x","y1\ny2","z"]}`),
		})
	})

	t.Run("multiline-fold", func(t *testing.T) {
		t.Parallel()
		testYAMLToJSON(t, test{
			name: "multiline-fold",
			args: args{
				data: []byte("a:\n  - x\n  - >\n    y1\n    y2\n  - z\n"),
			},
			want: []byte(`{"a":["x","y1 y2\n","z"]}`),
		})
	})

	t.Run("multiline-fold-chomp", func(t *testing.T) {
		t.Parallel()
		testYAMLToJSON(t, test{
			name: "multiline-fold-chomp",
			args: args{
				data: []byte("a:\n  - x\n  - >-\n    y1\n    y2\n  - z\n"),
			},
			want: []byte(`{"a":["x","y1 y2","z"]}`),
		})
	})

	t.Run("bad-format", func(t *testing.T) {
		t.Parallel()
		testYAMLToJSON(t, test{
			name: "bad-format",
			args: args{
				data: []byte("a: nullTag\n b: null\n  b: null\n"),
			},
			wantErr: true,
		})
	})
}

func testYAMLToJSON(t *testing.T, tt test) {
	got, err := YAMLToJSON(tt.args.data)
	if (err != nil) != tt.wantErr {
		t.Errorf("YAMLToJSON(%s) error = %v, wantErr %v", tt.name, err, tt.wantErr)
		return
	}
	if !reflect.DeepEqual(got, tt.want) {
		t.Errorf("YAMLToJSON(%s) = %v, want %v", tt.name, string(got), string(tt.want))
	}
}
