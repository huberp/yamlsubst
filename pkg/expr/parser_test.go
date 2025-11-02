package expr

import (
	"fmt"
	"math"
	"testing"
)

// Helper resolver for tests
func testResolver(values map[string]float64) func(string) (float64, error) {
	return func(path string) (float64, error) {
		if val, ok := values[path]; ok {
			return val, nil
		}
		return 0, fmt.Errorf("reference not found: %s", path)
	}
}

func TestParse_Numbers(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"integer", "42", "42"},
		{"float", "3.14", "3.14"},
		{"large number", "1234567", "1234567"},
		{"small decimal", "0.001", "0.001"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, err := Parse(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got := node.String(); got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestParse_References(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"simple reference", ".value", ".value"},
		{"nested reference", ".app.config.port", ".app.config.port"},
		{"reference with numbers", ".item1.count", ".item1.count"},
		{"reference with underscore", ".my_value", ".my_value"},
		{"env var simple", "$PORT", "$PORT"},
		{"env var uppercase", "$DATABASE_HOST", "$DATABASE_HOST"},
		{"env var with underscore", "$MY_VAR", "$MY_VAR"},
		{"env var with numbers", "$VAR123", "$VAR123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, err := Parse(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got := node.String(); got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestParse_Addition(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"two numbers", "5 + 3", "(5 + 3)"},
		{"three numbers", "1 + 2 + 3", "((1 + 2) + 3)"},
		{"no spaces", "10+20", "(10 + 20)"},
		{"reference and number", ".x + 5", "(.x + 5)"},
		{"two references", ".a + .b", "(.a + .b)"},
		{"env var and number", "$PORT + 1000", "($PORT + 1000)"},
		{"yaml ref and env var", ".base + $OFFSET", "(.base + $OFFSET)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, err := Parse(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got := node.String(); got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestParse_Subtraction(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"two numbers", "10 - 3", "(10 - 3)"},
		{"three numbers", "100 - 20 - 5", "((100 - 20) - 5)"},
		{"reference and number", ".y - 10", "(.y - 10)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, err := Parse(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got := node.String(); got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestParse_Multiplication(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"two numbers", "4 * 5", "(4 * 5)"},
		{"three numbers", "2 * 3 * 4", "((2 * 3) * 4)"},
		{"reference and number", ".count * 2", "(.count * 2)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, err := Parse(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got := node.String(); got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestParse_Division(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"two numbers", "20 / 4", "(20 / 4)"},
		{"three numbers", "100 / 10 / 2", "((100 / 10) / 2)"},
		{"reference and number", ".total / 3", "(.total / 3)"},
		{"decimal division", "5.5 / 2.5", "(5.5 / 2.5)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, err := Parse(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got := node.String(); got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestParse_MixedOperations(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"add and multiply", "2 + 3 * 4", "(2 + (3 * 4))"},
		{"multiply and add", "3 * 4 + 2", "((3 * 4) + 2)"},
		{"all operators", "10 + 5 * 2 - 8 / 4", "((10 + (5 * 2)) - (8 / 4))"},
		{"subtract and divide", "20 - 10 / 2", "(20 - (10 / 2))"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, err := Parse(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got := node.String(); got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestParse_Parentheses(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"simple parens", "(5 + 3)", "(5 + 3)"},
		{"parens change precedence", "(2 + 3) * 4", "((2 + 3) * 4)"},
		{"nested parens", "((5 + 3) * 2)", "((5 + 3) * 2)"},
		{"complex parens", "(10 + 5) / (3 - 1)", "((10 + 5) / (3 - 1))"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, err := Parse(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got := node.String(); got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestParse_ComplexExpressions(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			"references and numbers",
			".x + .y * 2 - 5",
			"((.x + (.y * 2)) - 5)",
		},
		{
			"nested references",
			".app.width * .app.height",
			"(.app.width * .app.height)",
		},
		{
			"complex with parens",
			"(.base + .offset) * .multiplier / 2",
			"(((.base + .offset) * .multiplier) / 2)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, err := Parse(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got := node.String(); got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestParse_Errors(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"empty input", ""},
		{"only operator", "+"},
		{"missing operand", "5 +"},
		{"double operator", "5 + * 3"},
		{"unmatched left paren", "(5 + 3"},
		{"unmatched right paren", "5 + 3)"},
		{"invalid token at end", "5 + 3 @"},
		{"invalid env var dollar only", "$"},
		{"invalid env var dollar digit", "$123"},
		{"invalid env var dollar special", "$-VAR"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse(tt.input)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
		})
	}
}

func TestEval_Numbers(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  float64
	}{
		{"integer", "42", 42},
		{"float", "3.14", 3.14},
		{"zero", "0", 0},
		{"negative via subtraction", "0 - 5", -5},
	}

	resolver := testResolver(map[string]float64{})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseAndEval(tt.input, resolver)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEval_Addition(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  float64
	}{
		{"simple", "5 + 3", 8},
		{"three numbers", "1 + 2 + 3", 6},
		{"floats", "1.5 + 2.5", 4.0},
		{"mixed", "10 + 5.5", 15.5},
	}

	resolver := testResolver(map[string]float64{})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseAndEval(tt.input, resolver)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEval_Subtraction(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  float64
	}{
		{"simple", "10 - 3", 7},
		{"three numbers", "100 - 20 - 5", 75},
		{"floats", "5.5 - 2.2", 3.3},
	}

	resolver := testResolver(map[string]float64{})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseAndEval(tt.input, resolver)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if math.Abs(got-tt.want) > 0.0001 {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEval_Multiplication(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  float64
	}{
		{"simple", "4 * 5", 20},
		{"three numbers", "2 * 3 * 4", 24},
		{"floats", "2.5 * 4", 10.0},
		{"by zero", "5 * 0", 0},
	}

	resolver := testResolver(map[string]float64{})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseAndEval(tt.input, resolver)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEval_Division(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  float64
	}{
		{"simple", "20 / 4", 5},
		{"three numbers", "100 / 10 / 2", 5},
		{"floats", "7.5 / 2.5", 3.0},
		{"result is float", "7 / 2", 3.5},
	}

	resolver := testResolver(map[string]float64{})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseAndEval(tt.input, resolver)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEval_DivisionByZero(t *testing.T) {
	resolver := testResolver(map[string]float64{})

	_, err := ParseAndEval("10 / 0", resolver)
	if err == nil {
		t.Fatal("expected division by zero error, got nil")
	}
}

func TestEval_Precedence(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  float64
	}{
		{"multiply before add", "2 + 3 * 4", 14},
		{"divide before subtract", "20 - 10 / 2", 15},
		{"left to right same precedence", "10 - 5 - 2", 3},
		{"complex", "10 + 5 * 2 - 8 / 4", 18},
	}

	resolver := testResolver(map[string]float64{})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseAndEval(tt.input, resolver)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEval_Parentheses(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  float64
	}{
		{"simple", "(5 + 3)", 8},
		{"override precedence", "(2 + 3) * 4", 20},
		{"nested", "((5 + 3) * 2)", 16},
		{"complex", "(10 + 5) / (3 - 1)", 7.5},
	}

	resolver := testResolver(map[string]float64{})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseAndEval(tt.input, resolver)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEval_References(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		values map[string]float64
		want   float64
	}{
		{
			"single reference",
			".x",
			map[string]float64{".x": 42},
			42,
		},
		{
			"reference addition",
			".a + .b",
			map[string]float64{".a": 10, ".b": 20},
			30,
		},
		{
			"reference and number",
			".count * 2",
			map[string]float64{".count": 5},
			10,
		},
		{
			"nested reference",
			".app.width * .app.height",
			map[string]float64{".app.width": 10, ".app.height": 5},
			50,
		},
		{
			"complex expression",
			"(.base + .offset) * .multiplier / 2",
			map[string]float64{".base": 10, ".offset": 5, ".multiplier": 4},
			30,
		},
		{
			"env var simple",
			"$PORT",
			map[string]float64{"$PORT": 8080},
			8080,
		},
		{
			"env var addition",
			"$PORT + 1000",
			map[string]float64{"$PORT": 8080},
			9080,
		},
		{
			"mixed yaml and env",
			".base + $OFFSET",
			map[string]float64{".base": 100, "$OFFSET": 50},
			150,
		},
		{
			"complex with env vars",
			"($PORT + $OFFSET) * .multiplier",
			map[string]float64{"$PORT": 8080, "$OFFSET": 20, ".multiplier": 2},
			16200,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resolver := testResolver(tt.values)
			got, err := ParseAndEval(tt.input, resolver)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEval_ReferenceNotFound(t *testing.T) {
	resolver := testResolver(map[string]float64{})

	_, err := ParseAndEval(".missing", resolver)
	if err == nil {
		t.Fatal("expected error for missing reference, got nil")
	}
}

func TestFormatResult(t *testing.T) {
	tests := []struct {
		name  string
		value float64
		want  string
	}{
		{"integer", 42, "42"},
		{"zero", 0, "0"},
		{"negative integer", -10, "-10"},
		{"float with decimals", 3.14, "3.14"},
		{"float that is integer", 10.0, "10"},
		{"small decimal", 0.5, "0.5"},
		{"large number", 1234567, "1234567"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatResult(tt.value)
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestLexer_Whitespace(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"spaces around operators", "5 + 3", "(5 + 3)"},
		{"tabs", "5\t+\t3", "(5 + 3)"},
		{"multiple spaces", "5   +   3", "(5 + 3)"},
		{"newlines", "5\n+\n3", "(5 + 3)"},
		{"mixed whitespace", "  5  +  3  ", "(5 + 3)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, err := Parse(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got := node.String(); got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestComplexScenarios(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		values map[string]float64
		want   float64
	}{
		{
			"calculate area",
			".width * .height",
			map[string]float64{".width": 10.5, ".height": 20},
			210,
		},
		{
			"percentage calculation",
			".total * .percent / 100",
			map[string]float64{".total": 200, ".percent": 15},
			30,
		},
		{
			"average",
			"(.a + .b + .c) / 3",
			map[string]float64{".a": 10, ".b": 20, ".c": 30},
			20,
		},
		{
			"temperature conversion",
			"(.fahrenheit - 32) * 5 / 9",
			map[string]float64{".fahrenheit": 100},
			37.77777777777778,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resolver := testResolver(tt.values)
			got, err := ParseAndEval(tt.input, resolver)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if math.Abs(got-tt.want) > 0.0001 {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
