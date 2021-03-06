package yaml

import (
	"bytes"
	"fmt"

	"gopkg.in/yaml.v3"
)

// Indent is the default YAML indent for this package.
// yaml.v3 defaults to 4 but industry standard is 2.
var Indent = 2

// JSONToYAML Converts JSON to YAML.
func JSONToYAML(data []byte) ([]byte, error) {
	n := new(yaml.Node)
	if err := yaml.Unmarshal(data, n); err != nil {
		return nil, err
	}

	if err := formatYAML(n); err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	enc := yaml.NewEncoder(buf)
	enc.SetIndent(Indent)

	if err := enc.Encode(n); err != nil {
		return nil, fmt.Errorf("marshal formated: %w", err)
	}

	return buf.Bytes(), nil
}

func formatYAML(n *yaml.Node) error {
	if n == nil {
		return nil
	}

	switch n.Kind {
	case yaml.DocumentNode:
		// NOOP - Document doesn't need styling
	case yaml.SequenceNode:
		n.Style = yaml.LiteralStyle
	case yaml.MappingNode:
		n.Style = yaml.LiteralStyle
	case yaml.ScalarNode:
		n.Style = yaml.FlowStyle
		// if n.Style == yaml.DoubleQuotedStyle {
		// 	n.Style = yaml.FlowStyle
		// 	if strings.Contains(n.Value, "\n") {
		// 		n.Style = yaml.LiteralStyle
		// 	}
		// }
	case yaml.AliasNode:
		return fmt.Errorf("formatYAML: alias node type not implemented")
	}

	for _, c := range n.Content {
		if err := formatYAML(c); err != nil {
			return err
		}
	}

	return nil
}
