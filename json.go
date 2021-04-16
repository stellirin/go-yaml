package yaml

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	nullTag      = "!!null"
	boolTag      = "!!bool"
	strTag       = "!!str"
	intTag       = "!!int"
	floatTag     = "!!float"
	timestampTag = "!!timestamp"
	binaryTag    = "!!binary"

	seqTag   = "!!seq"
	mapTag   = "!!map"
	mergeTag = "!!merge"
)

// YAMLToJSON converts YAML to JSON.
func YAMLToJSON(data []byte) ([]byte, error) {
	n := new(yaml.Node)
	err := yaml.Unmarshal(data, n)
	if err != nil {
		return nil, fmt.Errorf("YAMLToJSON: %s", err)
	}

	obj, err := convertNode(n)
	if err != nil {
		return nil, fmt.Errorf("YAMLToJSON: %s", err)
	}

	return json.Marshal(obj)
}

func convertNode(n *yaml.Node) (interface{}, error) {
	switch n.Kind {
	case yaml.DocumentNode:
		return convertDocumentNode(n)
	case yaml.SequenceNode:
		return convertSequenceNode(n.Content)
	case yaml.MappingNode:
		m := make(map[string]interface{})
		if err := convertMappingNode(n.Content, m); err != nil {
			return nil, err
		}
		return m, nil
	case yaml.ScalarNode:
		return convertScalarNode(n)
	case yaml.AliasNode:
		return nil, fmt.Errorf("alias node type not (yet) implemented")
	default:
		return nil, fmt.Errorf("unknown node type: %d", n.Kind)
	}
}

func convertDocumentNode(n *yaml.Node) (interface{}, error) {
	return convertNode(n.Content[0])
}

func convertSequenceNode(n []*yaml.Node) ([]interface{}, error) {
	var s []interface{}
	for _, c := range n {
		v, err := convertNode(c)
		if err != nil {
			return nil, err
		}
		s = append(s, v)
	}

	return s, nil
}

func convertMappingNode(n []*yaml.Node, m map[string]interface{}) error {
	if len(n) <= 0 {
		return nil
	}

	var a int
	if n[0].Tag == mergeTag {
		a = 2
		if err := convertMappingNode(n[1].Alias.Content, m); err != nil {
			return err
		}
	}

	for c := a; c < len(n); c = c + 2 {
		// JSON allows only string keys
		if n[c].Tag == boolTag || n[c].Tag == intTag || n[c].Tag == floatTag || n[c].Tag == timestampTag {
			n[c].Tag = strTag
		}

		k, err := convertNode(n[c])
		if err != nil {
			return err
		}
		v, err := convertNode(n[c+1])
		if err != nil {
			return err
		}
		m[k.(string)] = v
	}

	return nil
}

func convertScalarNode(n *yaml.Node) (interface{}, error) {
	switch n.Tag {
	case nullTag:
		return nil, nil
	case boolTag:
		return strconv.ParseBool(n.Value)
	case strTag:
		return n.Value, nil
	case intTag:
		return strconv.Atoi(n.Value)
	case floatTag:
		return strconv.ParseFloat(n.Value, 32)
	case timestampTag:
		return time.Parse(time.RFC3339, n.Value)
	case binaryTag:
		return nil, fmt.Errorf("'!!binary' node tag not (yet) implemented in scalar node type")
	default:
		return nil, fmt.Errorf("'%s' tag should not be processed as a scalar node type", n.Tag)
	}
}
