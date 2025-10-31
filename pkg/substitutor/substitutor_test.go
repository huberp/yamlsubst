package substitutor

import (
	"testing"
)

func TestSubstitute_SimpleValue(t *testing.T) {
	yamlContent := `
name: John
age: 30
`
	input := "Hello, my name is ${.name} and I am ${.age} years old"
	expected := "Hello, my name is John and I am 30 years old"

	result, err := Substitute(input, yamlContent)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestSubstitute_NestedValue(t *testing.T) {
	yamlContent := `
person:
  name: Alice
  address:
    city: Seattle
    state: WA
`
	input := "I live in ${.person.address.city}, ${.person.address.state}"
	expected := "I live in Seattle, WA"

	result, err := Substitute(input, yamlContent)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestSubstitute_MissingKey(t *testing.T) {
	yamlContent := `
name: John
`
	input := "Hello ${.missing}"

	result, err := Substitute(input, yamlContent)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Missing keys should be left as-is
	if result != input {
		t.Errorf("expected %q, got %q", input, result)
	}
}

func TestSubstitute_InvalidYAML(t *testing.T) {
	yamlContent := `
invalid: [unclosed
`
	input := "Test ${.invalid}"

	_, err := Substitute(input, yamlContent)
	if err == nil {
		t.Fatal("expected error for invalid YAML, got nil")
	}
}

func TestSubstitute_MultipleOccurrences(t *testing.T) {
	yamlContent := `
word: test
`
	input := "${.word} ${.word} ${.word}"
	expected := "test test test"

	result, err := Substitute(input, yamlContent)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestSubstitute_NumericValue(t *testing.T) {
	yamlContent := `
count: 42
price: 19.99
`
	input := "Count: ${.count}, Price: $${.price}"
	expected := "Count: 42, Price: $19.99"

	result, err := Substitute(input, yamlContent)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestSubstitute_EmptyInput(t *testing.T) {
	yamlContent := `
name: John
`
	input := ""
	expected := ""

	result, err := Substitute(input, yamlContent)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}
