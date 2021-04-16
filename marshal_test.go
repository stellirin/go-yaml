package yaml_test

import (
	"reflect"
	"testing"
	"time"

	. "czechia.dev/yaml"
)

type Object struct {
	A *Object   `json:"a"`
	B bool      `json:"b"`
	C string    `json:"c"`
	D int       `json:"d"`
	E float32   `json:"e"`
	F time.Time `json:"f"`
}

type args struct {
	data []byte
	obj  Object
}

type test struct {
	name    string
	args    args
	want    []byte
	wantObj Object
	wantErr bool
}

var testTime, _ = time.Parse(time.RFC3339, "2001-02-03T04:05:06.07Z")

func TestMarshal(t *testing.T) {

	t.Run("obj", func(t *testing.T) {
		t.Parallel()
		testMarshal(t, test{
			name: "obj",
			args: args{
				obj: Object{
					A: nil,
					B: true,
					C: "x",
					D: 42,
					E: 3.14159,
					F: testTime,
				},
			},
			want: []byte("a: null\nb: true\nc: x\nd: 42\ne: 3.14159\nf: \"2001-02-03T04:05:06.07Z\"\n"),
		})
	})

}

func TestUnmarshal(t *testing.T) {

	t.Run("obj", func(t *testing.T) {
		t.Parallel()
		testUnmarshal(t, test{
			name: "obj",
			args: args{
				data: []byte("a: null\nb: true\nc: x\nd: 42\ne: 3.14159\nf: \"2001-02-03T04:05:06.07Z\"\n"),
			},
			wantObj: Object{A: nil, B: true, C: "x", D: 42, E: 3.14159, F: testTime},
		})
	})

	t.Run("obj-bad", func(t *testing.T) {
		t.Parallel()
		testUnmarshal(t, test{
			name: "obj-bad",
			args: args{
				data: []byte("f: \"time\"\n"),
			},
			wantErr: true,
		})
	})

}

func testMarshal(t *testing.T, tt test) {
	got, err := Marshal(tt.args.obj)
	if (err != nil) != tt.wantErr {
		t.Errorf("Marshal(%s) error = %v, wantErr %v", tt.name, err, tt.wantErr)
		return
	}
	if !reflect.DeepEqual(got, tt.want) {
		t.Errorf("Marshal(%s) = %v, want %v", tt.name, string(got), string(tt.want))
	}
}

func testUnmarshal(t *testing.T, tt test) {
	err := Unmarshal(tt.args.data, &tt.args.obj)
	if (err != nil) != tt.wantErr {
		t.Errorf("Unmarshal(%s) error = %v, wantErr %v", tt.name, err, tt.wantErr)
		return
	}
	if !reflect.DeepEqual(tt.args.obj, tt.wantObj) {
		t.Errorf("Unmarshal(%s) = %v, want %v", tt.name, tt.args.obj, tt.wantObj)
	}
}
