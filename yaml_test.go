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
			name: "string",
			args: args{
				data: []byte(`{"t":"a"}`),
			},
			want: []byte("t: a\n"),
		},
		{
			name: "number",
			args: args{
				data: []byte(`{"t":"1"}`),
			},
			want: []byte("t: \"1\"\n"),
		},
		{
			name: "null",
			args: args{
				data: []byte(`{"t":null}`),
			},
			want: []byte("t: null\n"),
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
				t.Errorf("JSONToYAML() = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}
