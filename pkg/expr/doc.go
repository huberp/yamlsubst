// Package expr provides an arithmetic expression parser and evaluator.
// It supports basic arithmetic operations (+, -, *, /) with proper operator
// precedence and parentheses. Expressions can contain hardcoded numbers
// (integers and floats) and YAML references starting with a dot.
//
// Example expressions:
//   - "5 + 3" -> 8
//   - "2 * (3 + 4)" -> 14
//   - ".width * .height" -> evaluates references
//   - "(.base + .offset) * 2" -> complex expression
//
// Usage:
//
//	// Parse an expression
//	node, err := expr.Parse("2 + 3 * 4")
//	if err != nil {
//		// Handle error
//	}
//
//	// Evaluate with a resolver function
//	resolver := func(path string) (float64, error) {
//		// Resolve YAML references to their numeric values
//		return value, nil
//	}
//	result, err := expr.Eval(node, resolver)
//
//	// Or use the convenience function
//	result, err := expr.ParseAndEval("2 + 3 * 4", resolver)
//
//	// Format the result
//	formatted := expr.FormatResult(result) // "14"
package expr
