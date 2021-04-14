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
		i, err = convertNode(n.Content[0])
	case yaml.SequenceNode:
		var s []interface{}
		for _, c := range n.Content {
			v, e := convertNode(c)
			if e != nil {
				err = e
				break
			}
			s = append(s, v)
		}
		i = &s
	case yaml.MappingNode:
		m := make(map[string]interface{})
		for c := 0; c < len(n.Content); c = c + 2 {
			k, e := convertNode(n.Content[c])
			if e != nil {
				err = e
				break
			}
			v, e := convertNode(n.Content[c+1])
			if e != nil {
				return nil, fmt.Errorf("YAMLToJSON: %s", e)
			}
			m[k.(string)] = v
		}
		i = &m
	case yaml.ScalarNode:
		switch n.Tag {
		case "!!null":
			// NOOP
		case "!!bool":
			i, err = strconv.ParseBool(n.Value)
		case "!!str":
			i = n.Value
		case "!!int":
			i, err = strconv.Atoi(n.Value)
		case "!!float":
			i, err = strconv.ParseFloat(n.Value, 32)
		case "!!timestamp":
			i, err = time.Parse(time.RFC3339, n.Value)
		case "!!seq":
			// preempted by by yaml.SequenceNode
			return nil, fmt.Errorf("YAMLToJSON: '!!seq' node type should not be processed as scalar node type")
		case "!!map":
			// preempted by yaml.MappingNode
			return nil, fmt.Errorf("YAMLToJSON: '!!map' node type should not be processed as scalar node type")
		case "!!binary":
			return nil, fmt.Errorf("YAMLToJSON: '!!binary' node tag not implemented in scalar node type")
		case "!!merge":
			return nil, fmt.Errorf("YAMLToJSON: '!!merge' node tag not implemented in  scalar node type")
		default:
			return nil, fmt.Errorf("YAMLToJSON: unknown node tag on scalar node type")
		}
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
