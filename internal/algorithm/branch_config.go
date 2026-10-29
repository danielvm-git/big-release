package algorithm

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v3"
)

// branchConfigFields has BranchConfig's fields without its methods, so the
// decoders below can fall back to the default struct decoding without recursing.
type branchConfigFields BranchConfig

// UnmarshalYAML accepts both branch forms the README documents:
// a bare name (`- main`) and an object (`- name: beta`). See BUG-branch-string-form.
func (b *BranchConfig) UnmarshalYAML(node *yaml.Node) error {
	switch node.Kind {
	case yaml.ScalarNode:
		*b = BranchConfig{Name: node.Value}
		return nil
	case yaml.MappingNode:
		return node.Decode((*branchConfigFields)(b))
	default:
		return fmt.Errorf("line %d: branch must be a name or an object with a name, e.g. `- main` or `- name: main`", node.Line)
	}
}

// UnmarshalJSON accepts `"main"` as well as `{"name": "main"}`.
func (b *BranchConfig) UnmarshalJSON(data []byte) error {
	var name string
	if err := json.Unmarshal(data, &name); err == nil {
		*b = BranchConfig{Name: name}
		return nil
	}
	return json.Unmarshal(data, (*branchConfigFields)(b))
}
