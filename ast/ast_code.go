package ast

import "fmt"

// CustomParserCodeExpr supports custom parser.
type CustomParserCodeExpr struct {
	p      Pos
	Code   *CodeBlock
	FuncIx int
	X      int
}

var _ Expression = (*CustomParserCodeExpr)(nil)

// NewCustomParserCodeExpr creates a new state (#) code expression at the specified
// position.
func NewCustomParserCodeExpr(p Pos) *CustomParserCodeExpr {
	return &CustomParserCodeExpr{p: p, X: 1}
}

// Pos returns the starting position of the node.
func (s *CustomParserCodeExpr) Pos() Pos { return s.p }

// String returns the textual representation of a node.
func (s *CustomParserCodeExpr) String() string {
	return fmt.Sprintf("%s: %T{Code: %v}", s.p, s, s.Code)
}

// NullableVisit recursively determines whether an object is nullable.
func (s *CustomParserCodeExpr) NullableVisit(rules map[string]*Rule) bool {
	return true
}

// IsNullable returns the nullable attribute of the node.
func (s *CustomParserCodeExpr) IsNullable() bool {
	return true
}

// InitialNames returns names of nodes with which an expression can begin.
func (s *CustomParserCodeExpr) InitialNames() map[string]struct{} {
	return make(map[string]struct{})
}
