package substitutor

import (
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/huberp/yamlsubst/pkg/expr"
	"gopkg.in/yaml.v3"
)

// placeholderRegex is compiled once for better performance
// Updated to match any content inside ${...}
var placeholderRegex = regexp.MustCompile(`\$\{([^}]+)\}`)

// Substitute replaces placeholders in the input string with values from the YAML content.
// Placeholders are in the format ${expression} where expression can be:
// - A simple YAML reference: ${.path.to.value}
// - An environment variable: ${$VAR}
// - An arithmetic expression: ${.width * .height}, ${$PORT + 1000}, ${.base + $OFFSET}
func Substitute(input, yamlContent string) (string, error) {
	// Parse YAML content
	var data interface{}
	if err := yaml.Unmarshal([]byte(yamlContent), &data); err != nil {
		return "", fmt.Errorf("failed to parse YAML: %w", err)
	}

	result := placeholderRegex.ReplaceAllStringFunc(input, func(match string) string {
		// Extract the expression (remove ${ and })
		expression := match[2 : len(match)-1] // Remove ${ and }

		// Try to evaluate as expression
		value, err := evaluateExpression(expression, data)
		if err != nil {
			// If evaluation fails, keep the placeholder as-is
			return match
		}

		return value
	})

	return result, nil
}

// evaluateExpression evaluates an expression which can be a simple reference or arithmetic expression
func evaluateExpression(expression string, yamlData interface{}) (string, error) {
	// Create a resolver function that can handle both YAML refs and env vars
	resolver := func(ref string) (float64, error) {
		if len(ref) == 0 {
			return 0, fmt.Errorf("empty reference")
		}

		switch ref[0] {
		case '.':
			// YAML reference
			value := navigate(yamlData, ref)
			if value == nil {
				return 0, fmt.Errorf("reference not found: %s", ref)
			}
			return valueToFloat(value)
		case '$':
			// Environment variable
			envVar := ref[1:] // Remove $
			envValue := os.Getenv(envVar)
			if envValue == "" {
				return 0, fmt.Errorf("env var not found: %s", envVar)
			}
			return strconv.ParseFloat(envValue, 64)
		}

		return 0, fmt.Errorf("invalid reference: %s", ref)
	}

	// Try to parse and evaluate as expression
	result, err := expr.ParseAndEval(expression, resolver)
	if err != nil {
		// If it's not a valid expression, it might be a simple string reference
		// Try to navigate directly (for backward compatibility with non-numeric values)
		if expression[0] == '.' {
			value := navigate(yamlData, expression)
			if value != nil {
				return valueToString(value), nil
			}
		}
		return "", err
	}

	// Format the numeric result
	return expr.FormatResult(result), nil
}

// valueToFloat converts a value to float64 for expression evaluation
func valueToFloat(value interface{}) (float64, error) {
	switch v := value.(type) {
	case int:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case float64:
		return v, nil
	case string:
		// Try to parse string as float
		return strconv.ParseFloat(v, 64)
	default:
		return 0, fmt.Errorf("cannot convert to float: %v", value)
	}
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
