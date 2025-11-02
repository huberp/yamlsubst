// Package expr provides an arithmetic expression parser and evaluator.
// It supports basic arithmetic operations (+, -, *, /) with proper operator
// precedence and parentheses. Expressions can contain hardcoded numbers
// (integers and floats), YAML references starting with a dot, and environment
// variable references starting with a dollar sign.
//
// Example expressions:
//   - "5 + 3" -> 8
//   - "2 * (3 + 4)" -> 14
//   - ".width * .height" -> evaluates YAML references
//   - "$PORT + 1000" -> evaluates environment variable
//   - "(.base + $OFFSET) * 2" -> complex expression with both types
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
//	resolver := func(ref string) (float64, error) {
//		// Resolve YAML references (starting with .) or env vars (starting with $)
//		// to their numeric values
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
