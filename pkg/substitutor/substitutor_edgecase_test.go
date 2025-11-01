package substitutor

import "testing"

// TestNavigate_EdgeCases tests edge cases in path navigation
func TestNavigate_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		data     interface{}
		path     string
		expected interface{}
	}{
		{
			name:     "empty path",
			data:     map[string]interface{}{"key": "value"},
			path:     "",
			expected: map[string]interface{}{"key": "value"},
		},
		{
			name:     "just dot",
			data:     map[string]interface{}{"key": "value"},
			path:     ".",
			expected: map[string]interface{}{"key": "value"},
		},
		{
			name:     "leading dot",
			data:     map[string]interface{}{"key": "value"},
			path:     ".key",
			expected: "value",
		},
		{
			name:     "no leading dot",
			data:     map[string]interface{}{"key": "value"},
			path:     "key",
			expected: "value",
		},
		{
			name: "nested path",
			data: map[string]interface{}{
				"level1": map[string]interface{}{
					"level2": map[string]interface{}{
						"level3": "deep",
					},
				},
			},
			path:     ".level1.level2.level3",
			expected: "deep",
		},
		{
			name: "map[interface{}]interface{} type",
			data: map[interface{}]interface{}{
				"key": "value",
			},
			path:     "key",
			expected: "value",
		},
		{
			name:     "non-existent key",
			data:     map[string]interface{}{"key": "value"},
			path:     "missing",
			expected: nil,
		},
		{
			name:     "non-map type",
			data:     "string value",
			path:     "key",
			expected: nil,
		},
		{
			name: "path stops at non-map",
			data: map[string]interface{}{
				"key": "value",
			},
			path:     "key.deeper",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := navigate(tt.data, tt.path)
			
			// For map comparisons, use string representation
			if result == nil && tt.expected == nil {
				return
			}
			
			if result == nil || tt.expected == nil {
				if result != tt.expected {
					t.Errorf("expected %v, got %v", tt.expected, result)
				}
				return
			}
			
			// Compare string representations for complex types
			if resStr, ok := result.(string); ok {
				if expStr, ok := tt.expected.(string); ok {
					if resStr != expStr {
						t.Errorf("expected %q, got %q", expStr, resStr)
					}
					return
				}
			}
			
			// For maps, just check they're not nil
			if _, ok := result.(map[string]interface{}); ok {
				if _, ok := tt.expected.(map[string]interface{}); !ok {
					t.Errorf("expected map but got different type")
				}
				return
			}
		})
	}
}

// TestValueToString_AllTypes tests all type conversions
func TestValueToString_AllTypes(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		expected string
	}{
		{"string", "hello", "hello"},
		{"int", int(42), "42"},
		{"int64", int64(9223372036854775807), "9223372036854775807"},
		{"float64 integer", float64(42.0), "42"},
		{"float64 decimal", float64(3.14), "3.14"},
		{"bool true", true, "true"},
		{"bool false", false, "false"},
		{"slice", []string{"a", "b"}, "[a b]"},
		{"nil", nil, "<nil>"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := valueToString(tt.value)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

// TestSubstitute_BoolValue tests boolean value substitution
func TestSubstitute_BoolValue(t *testing.T) {
	yamlContent := `
enabled: true
disabled: false
`
	input := "Enabled: ${.enabled}, Disabled: ${.disabled}"
	expected := "Enabled: true, Disabled: false"

	result, err := Substitute(input, yamlContent)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

// TestSubstitute_NoPlaceholders tests input without placeholders
func TestSubstitute_NoPlaceholders(t *testing.T) {
	yamlContent := `
name: John
`
	input := "Hello world, no placeholders here"

	result, err := Substitute(input, yamlContent)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != input {
		t.Errorf("expected %q, got %q", input, result)
	}
}

// TestSubstitute_ComplexPath tests paths with multiple levels
func TestSubstitute_ComplexPath(t *testing.T) {
	yamlContent := `
app:
  config:
    database:
      host: localhost
      port: 5432
`
	input := "Connect to ${.app.config.database.host}:${.app.config.database.port}"
	expected := "Connect to localhost:5432"

	result, err := Substitute(input, yamlContent)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}
