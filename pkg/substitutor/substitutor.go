package substitutor

import (
	"fmt"
	"regexp"
	"strconv"

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
	if len(path) > 0 && path[0] == '.' {
		path = path[1:]
	}

	if path == "" {
		return data
	}

	// Navigate through path segments without allocating a slice
	current := data
	start := 0
	for i := 0; i <= len(path); i++ {
		if i == len(path) || path[i] == '.' {
			if i > start {
				part := path[start:i]

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
			start = i + 1
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
		return strconv.FormatFloat(v, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(v)
	default:
		return fmt.Sprintf("%v", v)
	}
}
