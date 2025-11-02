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

func TestSubstitute_SimpleExpression(t *testing.T) {
	yamlContent := `
width: 10
height: 5
`
	input := "Area: ${.width * .height}"
	expected := "Area: 50"

	result, err := Substitute(input, yamlContent)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestSubstitute_ExpressionWithAddition(t *testing.T) {
	yamlContent := `
base: 100
offset: 25
`
	input := "Total: ${.base + .offset}"
	expected := "Total: 125"

	result, err := Substitute(input, yamlContent)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestSubstitute_ComplexExpression(t *testing.T) {
	yamlContent := `
width: 10
height: 5
padding: 2
`
	input := "Result: ${(.width + .padding) * (.height + .padding)}"
	expected := "Result: 84"

	result, err := Substitute(input, yamlContent)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestSubstitute_ExpressionWithLiteral(t *testing.T) {
	yamlContent := `
base: 100
`
	input := "Result: ${.base * 2 + 50}"
	expected := "Result: 250"

	result, err := Substitute(input, yamlContent)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestSubstitute_ExpressionWithFloats(t *testing.T) {
	yamlContent := `
price: 19.99
quantity: 3
`
	input := "Total: ${.price * .quantity}"
	expected := "Total: 59.97"

	result, err := Substitute(input, yamlContent)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestSubstitute_InvalidExpression(t *testing.T) {
	yamlContent := `
value: 10
`
	input := "Result: ${.value +}"
	expected := "Result: ${.value +}" // Should keep placeholder as-is

	result, err := Substitute(input, yamlContent)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestSubstitute_MixedSimpleAndExpression(t *testing.T) {
	yamlContent := `
name: Test
width: 10
height: 5
`
	input := "Name: ${.name}, Area: ${.width * .height}"
	expected := "Name: Test, Area: 50"

	result, err := Substitute(input, yamlContent)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestSubstitute_EnvVariable(t *testing.T) {
// Set up test env var
t.Setenv("TEST_PORT", "8080")

yamlContent := `
base: 1000
`
input := "Result: ${$TEST_PORT + .base}"
expected := "Result: 9080"

result, err := Substitute(input, yamlContent)
if err != nil {
t.Fatalf("unexpected error: %v", err)
}

if result != expected {
t.Errorf("expected %q, got %q", expected, result)
}
}

func TestSubstitute_EnvVariableOnly(t *testing.T) {
// Set up test env var
t.Setenv("TEST_VALUE", "42")

yamlContent := `
name: test
`
input := "Value: ${$TEST_VALUE}"
expected := "Value: 42"

result, err := Substitute(input, yamlContent)
if err != nil {
t.Fatalf("unexpected error: %v", err)
}

if result != expected {
t.Errorf("expected %q, got %q", expected, result)
}
}

func TestSubstitute_EnvVariableMissing(t *testing.T) {
yamlContent := `
value: 10
`
input := "Result: ${$MISSING_VAR + .value}"
expected := "Result: ${$MISSING_VAR + .value}" // Should keep placeholder as-is

result, err := Substitute(input, yamlContent)
if err != nil {
t.Fatalf("unexpected error: %v", err)
}

if result != expected {
t.Errorf("expected %q, got %q", expected, result)
}
}
