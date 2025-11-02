package expr

import (
	"fmt"
	"strconv"
	"unicode"
)

// tokenType represents the type of a token
type tokenType int

const (
	tokenNumber tokenType = iota
	tokenReference
	tokenPlus
	tokenMinus
	tokenMultiply
	tokenDivide
	tokenLeftParen
	tokenRightParen
	tokenEOF
	tokenError
)

// token represents a lexical token
type token struct {
	typ   tokenType
	value string
}

// lexer tokenizes the input expression
type lexer struct {
	input string
	pos   int
}

// newLexer creates a new lexer
func newLexer(input string) *lexer {
	return &lexer{input: input, pos: 0}
}

// nextToken returns the next token from the input
func (l *lexer) nextToken() token {
	// Skip whitespace
	for l.pos < len(l.input) && unicode.IsSpace(rune(l.input[l.pos])) {
		l.pos++
	}

	// End of input
	if l.pos >= len(l.input) {
		return token{typ: tokenEOF}
	}

	ch := l.input[l.pos]

	// Single character tokens
	switch ch {
	case '+':
		l.pos++
		return token{typ: tokenPlus, value: "+"}
	case '-':
		l.pos++
		return token{typ: tokenMinus, value: "-"}
	case '*':
		l.pos++
		return token{typ: tokenMultiply, value: "*"}
	case '/':
		l.pos++
		return token{typ: tokenDivide, value: "/"}
	case '(':
		l.pos++
		return token{typ: tokenLeftParen, value: "("}
	case ')':
		l.pos++
		return token{typ: tokenRightParen, value: ")"}
	}

	// YAML reference (starts with .)
	if ch == '.' {
		start := l.pos
		l.pos++
		for l.pos < len(l.input) {
			ch := l.input[l.pos]
			if ch == '.' || unicode.IsLetter(rune(ch)) || unicode.IsDigit(rune(ch)) || ch == '_' {
				l.pos++
			} else {
				break
			}
		}
		return token{typ: tokenReference, value: l.input[start:l.pos]}
	}

	// Environment variable reference (starts with $)
	if ch == '$' {
		start := l.pos
		l.pos++
		// Env var name must start with letter or underscore
		if l.pos < len(l.input) {
			ch := l.input[l.pos]
			if unicode.IsLetter(rune(ch)) || ch == '_' {
				l.pos++
				for l.pos < len(l.input) {
					ch := l.input[l.pos]
					if unicode.IsLetter(rune(ch)) || unicode.IsDigit(rune(ch)) || ch == '_' {
						l.pos++
					} else {
						break
					}
				}
				return token{typ: tokenReference, value: l.input[start:l.pos]}
			}
		}
		// If $ is not followed by valid identifier, treat as error
		return token{typ: tokenError, value: "$"}
	}

	// Number (integer or float)
	if unicode.IsDigit(rune(ch)) {
		start := l.pos
		hasDecimal := false
		for l.pos < len(l.input) {
			ch := l.input[l.pos]
			if unicode.IsDigit(rune(ch)) {
				l.pos++
			} else if ch == '.' && !hasDecimal {
				hasDecimal = true
				l.pos++
			} else {
				break
			}
		}
		return token{typ: tokenNumber, value: l.input[start:l.pos]}
	}

	// Unknown character - return error token
	ch = l.input[l.pos]
	l.pos++
	return token{typ: tokenError, value: string(ch)}
}

// parser implements a recursive descent parser for arithmetic expressions
type parser struct {
	lexer   *lexer
	current token
}

// newParser creates a new parser
func newParser(input string) *parser {
	p := &parser{lexer: newLexer(input)}
	p.advance()
	return p
}

// advance moves to the next token
func (p *parser) advance() {
	p.current = p.lexer.nextToken()
}

// Parse parses the expression and returns an AST node
func Parse(input string) (Node, error) {
	p := newParser(input)
	if p.current.typ == tokenError {
		return nil, fmt.Errorf("unexpected character: %s", p.current.value)
	}
	node, err := p.parseExpression()
	if err != nil {
		return nil, err
	}
	if p.current.typ == tokenError {
		return nil, fmt.Errorf("unexpected character: %s", p.current.value)
	}
	if p.current.typ != tokenEOF {
		return nil, fmt.Errorf("unexpected token: %s", p.current.value)
	}
	return node, nil
}

// parseExpression parses addition and subtraction (lowest precedence)
func (p *parser) parseExpression() (Node, error) {
	left, err := p.parseTerm()
	if err != nil {
		return nil, err
	}

	for p.current.typ == tokenPlus || p.current.typ == tokenMinus {
		op := p.current.value
		p.advance()
		right, err := p.parseTerm()
		if err != nil {
			return nil, err
		}
		left = &BinaryOpNode{Left: left, Op: op, Right: right}
	}

	return left, nil
}

// parseTerm parses multiplication and division (higher precedence)
func (p *parser) parseTerm() (Node, error) {
	left, err := p.parseFactor()
	if err != nil {
		return nil, err
	}

	for p.current.typ == tokenMultiply || p.current.typ == tokenDivide {
		op := p.current.value
		p.advance()
		right, err := p.parseFactor()
		if err != nil {
			return nil, err
		}
		left = &BinaryOpNode{Left: left, Op: op, Right: right}
	}

	return left, nil
}

// parseFactor parses numbers, references, and parenthesized expressions
func (p *parser) parseFactor() (Node, error) {
	switch p.current.typ {
	case tokenNumber:
		value := p.current.value
		p.advance()
		// Parse as float to support both int and float
		f, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid number: %s", value)
		}
		return &NumberNode{Value: f}, nil

	case tokenReference:
		path := p.current.value
		p.advance()
		return &ReferenceNode{Path: path}, nil

	case tokenLeftParen:
		p.advance()
		node, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		if p.current.typ != tokenRightParen {
			return nil, fmt.Errorf("expected ')', got %s", p.current.value)
		}
		p.advance()
		return node, nil

	default:
		return nil, fmt.Errorf("unexpected token: %s", p.current.value)
	}
}

// Node represents a node in the AST
type Node interface {
	String() string
}

// formatFloat formats a float64, converting to int if it's a whole number
func formatFloat(value float64) string {
	// Check if the value is within safe integer range
	const maxSafeInt = 1<<53 - 1
	if value >= -maxSafeInt && value <= maxSafeInt {
		// Check if it's effectively an integer (no fractional part)
		if value == float64(int64(value)) {
			return strconv.FormatInt(int64(value), 10)
		}
	}
	// Format as float, removing trailing zeros
	s := strconv.FormatFloat(value, 'f', -1, 64)
	return s
}

// NumberNode represents a numeric literal
type NumberNode struct {
	Value float64
}

func (n *NumberNode) String() string {
	return formatFloat(n.Value)
}

// ReferenceNode represents a YAML reference or environment variable reference
type ReferenceNode struct {
	Path string
}

func (n *ReferenceNode) String() string {
	return n.Path
}

// BinaryOpNode represents a binary operation
type BinaryOpNode struct {
	Left  Node
	Op    string
	Right Node
}

func (n *BinaryOpNode) String() string {
	return fmt.Sprintf("(%s %s %s)", n.Left.String(), n.Op, n.Right.String())
}

// Eval evaluates the expression with the given resolver function
// The resolver function takes a reference path and returns its value
func Eval(node Node, resolver func(string) (float64, error)) (float64, error) {
	switch n := node.(type) {
	case *NumberNode:
		return n.Value, nil

	case *ReferenceNode:
		return resolver(n.Path)

	case *BinaryOpNode:
		left, err := Eval(n.Left, resolver)
		if err != nil {
			return 0, err
		}
		right, err := Eval(n.Right, resolver)
		if err != nil {
			return 0, err
		}

		switch n.Op {
		case "+":
			return left + right, nil
		case "-":
			return left - right, nil
		case "*":
			return left * right, nil
		case "/":
			if right == 0 {
				return 0, fmt.Errorf("division by zero")
			}
			return left / right, nil
		default:
			return 0, fmt.Errorf("unknown operator: %s", n.Op)
		}

	default:
		return 0, fmt.Errorf("unknown node type")
	}
}

// ParseAndEval is a convenience function that parses and evaluates an expression
func ParseAndEval(input string, resolver func(string) (float64, error)) (float64, error) {
	node, err := Parse(input)
	if err != nil {
		return 0, err
	}
	return Eval(node, resolver)
}

// FormatResult formats a float64 result, removing unnecessary decimal points
func FormatResult(value float64) string {
	return formatFloat(value)
}
