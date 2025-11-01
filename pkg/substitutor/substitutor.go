package substitutor

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// placeholderRegex is compiled once for better performance
var placeholderRegex = regexp.MustCompile(`\$\{(\.[^}]+)\}`)

// Substitute replaces placeholders in the input string with values from the YAML content.
// Placeholders are in the format ${.path.to.value} where the path is a dot-separated
// sequence of keys to navigate the YAML structure.
func Substitute(input, yamlContent string) (string, error) {
	// Parse YAML content
	var data interface{}
	if err := yaml.Unmarshal([]byte(yamlContent), &data); err != nil {
		return "", fmt.Errorf("failed to parse YAML: %w", err)
	}

	result := placeholderRegex.ReplaceAllStringFunc(input, func(match string) string {
		// Extract the path (remove ${ and })
		path := match[2 : len(match)-1] // Remove ${ and }

		// Navigate the YAML structure
		value := navigate(data, path)
		if value == nil {
			// If value not found, keep the placeholder as-is
			return match
		}

		// Convert value to string
		return valueToString(value)
	})

	return result, nil
}

// navigate traverses the YAML data structure using the given path
func navigate(data interface{}, path string) interface{} {
	// Remove leading dot if present
	path = strings.TrimPrefix(path, ".")

	if path == "" {
		return data
	}

	// Split path into parts
	parts := strings.Split(path, ".")

	current := data
	for _, part := range parts {
		if part == "" {
			continue
		}

		switch v := current.(type) {
		case map[string]interface{}:
			var ok bool
			current, ok = v[part]
			if !ok {
				return nil
			}
		case map[interface{}]interface{}:
			var ok bool
			current, ok = v[part]
			if !ok {
				return nil
			}
		default:
			return nil
		}
	}

	return current
}

// valueToString converts a value to its string representation
func valueToString(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	case int64:
		return strconv.FormatInt(v, 10)
	case float64:
		// Format float without unnecessary trailing zeros
		s := strconv.FormatFloat(v, 'f', -1, 64)
		return s
	case bool:
		return strconv.FormatBool(v)
	default:
		return fmt.Sprintf("%v", v)
	}
}
