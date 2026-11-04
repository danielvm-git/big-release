// story: e01s05
package algorithm

import (
	"encoding/json"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestBranchConfig_UnmarshalYAML(t *testing.T) {
	var cfg struct {
		Branches []BranchConfig `yaml:"branches"`
	}
	input := `
branches:
  - main
  - name: beta
    type: prerelease
    channel: beta
    prerelease: beta
`
	if err := yaml.Unmarshal([]byte(input), &cfg); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(cfg.Branches) != 2 {
		t.Fatalf("got %d branches, want 2", len(cfg.Branches))
	}
	if cfg.Branches[0].Name != "main" {
		t.Errorf("scalar branch name = %q, want main", cfg.Branches[0].Name)
	}
	beta := cfg.Branches[1]
	if beta.Name != "beta" || beta.Type != "prerelease" ||
		beta.Channel != "beta" || beta.Prerelease != "beta" {
		t.Errorf("mapping branch decoded incorrectly: %+v", beta)
	}
}

func TestBranchConfig_UnmarshalYAML_RejectsSequence(t *testing.T) {
	var cfg struct {
		Branches []BranchConfig `yaml:"branches"`
	}
	err := yaml.Unmarshal([]byte("branches:\n  - [a, b]\n"), &cfg)
	if err == nil {
		t.Fatal("expected an error for a sequence branch entry")
	}
}

func TestBranchConfig_UnmarshalJSON(t *testing.T) {
	cases := []struct {
		name string
		data string
		want BranchConfig
	}{
		{"string form", `"main"`, BranchConfig{Name: "main"}},
		{"object form", `{"name":"beta","type":"prerelease"}`,
			BranchConfig{Name: "beta", Type: "prerelease"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var got BranchConfig
			if err := json.Unmarshal([]byte(tc.data), &got); err != nil {
				t.Fatalf("unmarshal: %v", err)
			}
			if got != tc.want {
				t.Fatalf("got %+v, want %+v", got, tc.want)
			}
		})
	}
}

func TestBranchConfig_UnmarshalJSON_RejectsNumber(t *testing.T) {
	var got BranchConfig
	if err := json.Unmarshal([]byte("123"), &got); err == nil {
		t.Fatal("expected an error for a numeric branch entry")
	}
}
