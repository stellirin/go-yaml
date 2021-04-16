package yaml

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type Options struct {
	Indent int
}

// default indent for yaml.v3 is 4 but industry standard is 2.
var defaultOptions = Options{
	Indent: 2,
}

// Marshal converts an object into YAML, via an intermediate marshal to JSON.
func Marshal(obj interface{}) ([]byte, error) {
	j, err := json.Marshal(obj)
	if err != nil {
		return nil, fmt.Errorf("error marshaling into JSON: %v", err)
	}

	y, err := JSONToYAML(j)
	if err != nil {
		return nil, fmt.Errorf("error converting JSON to YAML: %v", err)
	}

	return y, nil
}

// JSONOpt is a decoding option for decoding from JSON format.
type JSONOpt func(*json.Decoder) *json.Decoder

// Unmarshal converts YAML to an object, via an intermediate marshal to JSON.
func Unmarshal(data []byte, obj interface{}, opts ...JSONOpt) error {
	j, err := YAMLToJSON(data)
	if err != nil {
		return fmt.Errorf("error converting YAML to JSON: %v", err)
	}

	r := bytes.NewReader(j)
	d := json.NewDecoder(r)
	for _, opt := range opts {
		d = opt(d)
	}

	if err := d.Decode(&obj); err != nil {
		return fmt.Errorf("error while decoding JSON: %v", err)
	}

	return nil
}
