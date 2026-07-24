// story: e24s03

package nodeutil

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// npmNamePattern matches valid npm package names per the npm registry spec.
// Does not include scope prefix; for scoped names use IsValidPackageName.
var npmNamePattern = regexp.MustCompile(`^[a-z0-9][-a-z0-9._]*$`)

// npmScopePattern matches valid npm scope names (after the initial @ and before /).
var npmScopePattern = regexp.MustCompile(`^@[a-z0-9][-a-z0-9._]*$`)

// IsValidPackageName validates an npm-compatible package name.
// Supports both unscoped ("my-package") and scoped ("@scope/my-package") formats.
// Enforces name rules to prevent flag injection via crafted names.
func IsValidPackageName(name string) bool {
	if len(name) == 0 || len(name) > 214 {
		return false
	}
	if name[0] == '@' {
		if !strings.Contains(name, "/") {
			return false
		}
		scopeEnd := strings.Index(name, "/")
		scope := name[:scopeEnd]
		if !npmScopePattern.MatchString(scope) {
			return false
		}
		rest := name[scopeEnd+1:]
		if len(rest) == 0 {
			return false
		}
		return npmNamePattern.MatchString(rest)
	}
	return npmNamePattern.MatchString(name)
}

// ReadPackageJSON reads and parses package.json in the working directory.
// prefix is used in error messages (e.g. "npm", "pnpm").
func ReadPackageJSON(prefix string) (map[string]interface{}, error) {
	data, err := os.ReadFile("package.json")
	if err != nil {
		return nil, fmt.Errorf("%s: failed to read package.json: %w", prefix, err)
	}

	var pkg map[string]interface{}
	if err := json.Unmarshal(data, &pkg); err != nil {
		return nil, fmt.Errorf("%s: failed to parse package.json: %w", prefix, err)
	}

	return pkg, nil
}

// WritePackageJSON marshals and writes package.json in the working directory.
func WritePackageJSON(prefix string, pkg map[string]interface{}) error {
	data, err := json.MarshalIndent(pkg, "", "  ")
	if err != nil {
		return fmt.Errorf("%s: failed to marshal package.json: %w", prefix, err)
	}

	if err := os.WriteFile("package.json", data, 0644); err != nil {
		return fmt.Errorf("%s: failed to write package.json: %w", prefix, err)
	}

	return nil
}

// ReadPackageName reads and validates the package name from package.json.
func ReadPackageName(prefix string) (string, error) {
	pkg, err := ReadPackageJSON(prefix)
	if err != nil {
		return "", err
	}

	name, ok := pkg["name"].(string)
	if !ok || name == "" {
		return "", fmt.Errorf("%s: package name not found or not a string in package.json", prefix)
	}
	if !IsValidPackageName(name) {
		return "", fmt.Errorf("%s: invalid package name %q in package.json", prefix, name)
	}

	return name, nil
}
