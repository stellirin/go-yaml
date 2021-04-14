package yaml

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"gopkg.in/yaml.v3"
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
	var i interface{}
	var err error

	switch n.Kind {
	case yaml.DocumentNode:
		i, err = convertDocumentNode(n)
	case yaml.SequenceNode:
		i, err = convertSequenceNode(n.Content)
	case yaml.MappingNode:
		i, err = convertMappingNode(n.Content)
	case yaml.ScalarNode:
		i, err = convertScalarNode(n)
	case yaml.AliasNode:
		return nil, fmt.Errorf("YAMLToJSON: alias node type not implemented")
	default:
		return nil, fmt.Errorf("YAMLToJSON: unknown node type: %d", n.Kind)
	}

	if err != nil {
		return nil, fmt.Errorf("YAMLToJSON: %s", err)
	}

	return i, nil
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

func convertMappingNode(n []*yaml.Node) (map[string]interface{}, error) {
	m := make(map[string]interface{})

	// if n.Content[0].Tag == "!!merge" {
	// 	m, err = convertNode(n.Content[1].Alias)
	// }

	for c := 0; c < len(n); c = c + 2 {
		k, err := convertNode(n[c])
		if err != nil {
			return nil, fmt.Errorf("YAMLToJSON: %s", err)
		}
		v, err := convertNode(n[c+1])
		if err != nil {
			return nil, fmt.Errorf("YAMLToJSON: %s", err)
		}
		m[k.(string)] = v
	}

	return m, nil
}

func convertScalarNode(n *yaml.Node) (interface{}, error) {
	switch n.Tag {
	case "!!null":
		return nil, nil
	case "!!bool":
		return strconv.ParseBool(n.Value)
	case "!!str":
		return n.Value, nil
	case "!!int":
		return strconv.Atoi(n.Value)
	case "!!float":
		return strconv.ParseFloat(n.Value, 32)
	case "!!timestamp":
		return time.Parse(time.RFC3339, n.Value)
	case "!!seq":
		// preempted by by yaml.SequenceNode
		return nil, fmt.Errorf("YAMLToJSON: '!!seq' node type should not be processed as scalar node type")
	case "!!map":
		// preempted by yaml.MappingNode
		return nil, fmt.Errorf("YAMLToJSON: '!!map' node type should not be processed as scalar node type")
	case "!!binary":
		return nil, fmt.Errorf("YAMLToJSON: '!!binary' node tag not implemented in scalar node type")
	case "!!merge":
		// i, err = convertNode(n.Content)
		return nil, fmt.Errorf("YAMLToJSON: '!!merge' node tag not implemented in  scalar node type")
	default:
		return nil, fmt.Errorf("YAMLToJSON: unknown node tag on scalar node type")
	}
}
