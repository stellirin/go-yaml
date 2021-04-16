package yaml_test

import (
	"reflect"
	"testing"

	. "czechia.dev/yaml"
)

type Object struct {
	A string `json:"a"`
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

func TestMarshal(t *testing.T) {

	t.Run("obj", func(t *testing.T) {
		t.Parallel()
		testMarshal(t, test{
			name: "obj",
			args: args{
				obj: Object{
					A: "x",
				},
			},
			want: []byte("a: x\n"),
		})
	})

}

func TestUnmarshal(t *testing.T) {

	t.Run("obj", func(t *testing.T) {
		t.Parallel()
		testUnmarshal(t, test{
			name: "obj",
			args: args{
				data: []byte("a: x\n"),
			},
			wantObj: Object{A: "x"},
		})
	})

	t.Run("obj-bad", func(t *testing.T) {
		t.Parallel()
		testUnmarshal(t, test{
			name: "obj-bad",
			args: args{
				data: []byte("a: 1\n"),
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
