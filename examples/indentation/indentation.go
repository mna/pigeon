package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"
)

func toIfaceSlice(v interface{}) []interface{} {
	if v == nil {
		return nil
	}
	return v.([]interface{})
}

var g = &grammar{
	rules: []*rule{
		{
			name: "Input",
			pos:  position{line: 13, col: 1, offset: 155},
			expr: &actionExpr{
				pos: position{line: 13, col: 15, offset: 171},
				run: (*parser).callonInput1,
				expr: &seqExpr{
					pos: position{line: 13, col: 15, offset: 171},
					exprs: []interface{}{
						&andCodeExpr{
							pos: position{line: 13, col: 15, offset: 171},
							run: (*parser).callonInput3,
						},
						&labeledExpr{
							pos:   position{line: 13, col: 64, offset: 220},
							label: "s",
							expr: &ruleRefExpr{
								pos:  position{line: 13, col: 66, offset: 222},
								name: "Statements",
							},
						},
						&labeledExpr{
							pos:   position{line: 13, col: 78, offset: 234},
							label: "r",
							expr: &ruleRefExpr{
								pos:  position{line: 13, col: 80, offset: 236},
								name: "ReturnOp",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 13, col: 89, offset: 245},
							name: "EOF",
						},
					},
				},
			},
		},
		{
			name: "Statements",
			pos:  position{line: 15, col: 1, offset: 357},
			expr: &actionExpr{
				pos: position{line: 15, col: 15, offset: 373},
				run: (*parser).callonStatements1,
				expr: &labeledExpr{
					pos:   position{line: 15, col: 15, offset: 373},
					label: "s",
					expr: &oneOrMoreExpr{
						pos: position{line: 15, col: 17, offset: 375},
						expr: &ruleRefExpr{
							pos:  position{line: 15, col: 17, offset: 375},
							name: "Line",
						},
					},
				},
			},
		},
		{
			name: "Line",
			pos:  position{line: 16, col: 1, offset: 435},
			expr: &actionExpr{
				pos: position{line: 16, col: 15, offset: 451},
				run: (*parser).callonLine1,
				expr: &seqExpr{
					pos: position{line: 16, col: 15, offset: 451},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 16, col: 15, offset: 451},
							name: "INDENTATION",
						},
						&labeledExpr{
							pos:   position{line: 16, col: 27, offset: 463},
							label: "s",
							expr: &ruleRefExpr{
								pos:  position{line: 16, col: 29, offset: 465},
								name: "Statement",
							},
						},
					},
				},
			},
		},
		{
			name: "ReturnOp",
			pos:  position{line: 17, col: 1, offset: 499},
			expr: &actionExpr{
				pos: position{line: 17, col: 15, offset: 515},
				run: (*parser).callonReturnOp1,
				expr: &seqExpr{
					pos: position{line: 17, col: 15, offset: 515},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 17, col: 15, offset: 515},
							val:        "return",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 17, col: 24, offset: 524},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 17, col: 26, offset: 526},
							label: "arg",
							expr: &ruleRefExpr{
								pos:  position{line: 17, col: 30, offset: 530},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 17, col: 41, offset: 541},
							name: "EOL",
						},
					},
				},
			},
		},
		{
			name: "Statement",
			pos:  position{line: 19, col: 1, offset: 594},
			expr: &choiceExpr{
				pos: position{line: 19, col: 15, offset: 610},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 19, col: 15, offset: 610},
						run: (*parser).callonStatement2,
						expr: &seqExpr{
							pos: position{line: 19, col: 15, offset: 610},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 19, col: 15, offset: 610},
									label: "s",
									expr: &ruleRefExpr{
										pos:  position{line: 19, col: 17, offset: 612},
										name: "Assignment",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 19, col: 28, offset: 623},
									name: "EOL",
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 20, col: 7, offset: 681},
						run: (*parser).callonStatement7,
						expr: &seqExpr{
							pos: position{line: 20, col: 7, offset: 681},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 20, col: 7, offset: 681},
									val:        "if",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 20, col: 12, offset: 686},
									name: "_",
								},
								&labeledExpr{
									pos:   position{line: 20, col: 14, offset: 688},
									label: "arg",
									expr: &ruleRefExpr{
										pos:  position{line: 20, col: 18, offset: 692},
										name: "LogicalExpression",
									},
								},
								&zeroOrOneExpr{
									pos: position{line: 20, col: 36, offset: 710},
									expr: &ruleRefExpr{
										pos:  position{line: 20, col: 36, offset: 710},
										name: "_",
									},
								},
								&litMatcher{
									pos:        position{line: 20, col: 39, offset: 713},
									val:        ":",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 20, col: 43, offset: 717},
									name: "EOL",
								},
								&ruleRefExpr{
									pos:  position{line: 20, col: 47, offset: 721},
									name: "INDENT",
								},
								&labeledExpr{
									pos:   position{line: 20, col: 54, offset: 728},
									label: "s",
									expr: &ruleRefExpr{
										pos:  position{line: 20, col: 56, offset: 730},
										name: "Statements",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 20, col: 67, offset: 741},
									name: "DEDENT",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Assignment",
			pos:  position{line: 24, col: 1, offset: 872},
			expr: &actionExpr{
				pos: position{line: 24, col: 14, offset: 887},
				run: (*parser).callonAssignment1,
				expr: &seqExpr{
					pos: position{line: 24, col: 14, offset: 887},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 24, col: 14, offset: 887},
							label: "lvalue",
							expr: &ruleRefExpr{
								pos:  position{line: 24, col: 21, offset: 894},
								name: "Identifier",
							},
						},
						&zeroOrOneExpr{
							pos: position{line: 24, col: 32, offset: 905},
							expr: &ruleRefExpr{
								pos:  position{line: 24, col: 32, offset: 905},
								name: "_",
							},
						},
						&litMatcher{
							pos:        position{line: 24, col: 35, offset: 908},
							val:        "=",
							ignoreCase: false,
						},
						&zeroOrOneExpr{
							pos: position{line: 24, col: 39, offset: 912},
							expr: &ruleRefExpr{
								pos:  position{line: 24, col: 39, offset: 912},
								name: "_",
							},
						},
						&labeledExpr{
							pos:   position{line: 24, col: 42, offset: 915},
							label: "rvalue",
							expr: &ruleRefExpr{
								pos:  position{line: 24, col: 49, offset: 922},
								name: "AdditiveExpression",
							},
						},
					},
				},
			},
		},
		{
			name: "LogicalExpression",
			pos:  position{line: 27, col: 1, offset: 1075},
			expr: &actionExpr{
				pos: position{line: 27, col: 23, offset: 1099},
				run: (*parser).callonLogicalExpression1,
				expr: &labeledExpr{
					pos:   position{line: 27, col: 23, offset: 1099},
					label: "arg",
					expr: &ruleRefExpr{
						pos:  position{line: 27, col: 27, offset: 1103},
						name: "PrimaryExpression",
					},
				},
			},
		},
		{
			name: "AdditiveExpression",
			pos:  position{line: 28, col: 1, offset: 1187},
			expr: &actionExpr{
				pos: position{line: 28, col: 23, offset: 1211},
				run: (*parser).callonAdditiveExpression1,
				expr: &seqExpr{
					pos: position{line: 28, col: 23, offset: 1211},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 28, col: 23, offset: 1211},
							label: "arg",
							expr: &ruleRefExpr{
								pos:  position{line: 28, col: 27, offset: 1215},
								name: "PrimaryExpression",
							},
						},
						&labeledExpr{
							pos:   position{line: 28, col: 45, offset: 1233},
							label: "rest",
							expr: &zeroOrMoreExpr{
								pos: position{line: 28, col: 50, offset: 1238},
								expr: &seqExpr{
									pos: position{line: 28, col: 52, offset: 1240},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 28, col: 52, offset: 1240},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 28, col: 54, offset: 1242},
											name: "AddOp",
										},
										&ruleRefExpr{
											pos:  position{line: 28, col: 60, offset: 1248},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 28, col: 62, offset: 1250},
											name: "PrimaryExpression",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "PrimaryExpression",
			pos:  position{line: 30, col: 1, offset: 1388},
			expr: &actionExpr{
				pos: position{line: 30, col: 23, offset: 1412},
				run: (*parser).callonPrimaryExpression1,
				expr: &labeledExpr{
					pos:   position{line: 30, col: 23, offset: 1412},
					label: "arg",
					expr: &choiceExpr{
						pos: position{line: 30, col: 28, offset: 1417},
						alternatives: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 30, col: 28, offset: 1417},
								name: "Integer",
							},
							&ruleRefExpr{
								pos:  position{line: 30, col: 38, offset: 1427},
								name: "Identifier",
							},
						},
					},
				},
			},
		},
		{
			name: "Integer",
			pos:  position{line: 33, col: 1, offset: 1529},
			expr: &actionExpr{
				pos: position{line: 33, col: 11, offset: 1541},
				run: (*parser).callonInteger1,
				expr: &oneOrMoreExpr{
					pos: position{line: 33, col: 11, offset: 1541},
					expr: &charClassMatcher{
						pos:        position{line: 33, col: 11, offset: 1541},
						val:        "[0-9]",
						ranges:     []rune{'0', '9'},
						ignoreCase: false,
						inverted:   false,
					},
				},
			},
		},
		{
			name: "Identifier",
			pos:  position{line: 34, col: 1, offset: 1618},
			expr: &actionExpr{
				pos: position{line: 34, col: 14, offset: 1633},
				run: (*parser).callonIdentifier1,
				expr: &seqExpr{
					pos: position{line: 34, col: 14, offset: 1633},
					exprs: []interface{}{
						&charClassMatcher{
							pos:        position{line: 34, col: 14, offset: 1633},
							val:        "[a-zA-Z]",
							ranges:     []rune{'a', 'z', 'A', 'Z'},
							ignoreCase: false,
							inverted:   false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 34, col: 23, offset: 1642},
							expr: &charClassMatcher{
								pos:        position{line: 34, col: 23, offset: 1642},
								val:        "[a-zA-Z0-9]",
								ranges:     []rune{'a', 'z', 'A', 'Z', '0', '9'},
								ignoreCase: false,
								inverted:   false,
							},
						},
					},
				},
			},
		},
		{
			name: "AddOp",
			pos:  position{line: 36, col: 1, offset: 1712},
			expr: &actionExpr{
				pos: position{line: 36, col: 9, offset: 1722},
				run: (*parser).callonAddOp1,
				expr: &choiceExpr{
					pos: position{line: 36, col: 11, offset: 1724},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 36, col: 11, offset: 1724},
							val:        "+",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 36, col: 17, offset: 1730},
							val:        "-",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "_",
			pos:  position{line: 38, col: 1, offset: 1791},
			expr: &oneOrMoreExpr{
				pos: position{line: 38, col: 5, offset: 1797},
				expr: &charClassMatcher{
					pos:        position{line: 38, col: 5, offset: 1797},
					val:        "[ \\t]",
					chars:      []rune{' ', '\t'},
					ignoreCase: false,
					inverted:   false,
				},
			},
		},
		{
			name: "EOL",
			pos:  position{line: 40, col: 1, offset: 1807},
			expr: &seqExpr{
				pos: position{line: 40, col: 7, offset: 1815},
				exprs: []interface{}{
					&zeroOrOneExpr{
						pos: position{line: 40, col: 7, offset: 1815},
						expr: &ruleRefExpr{
							pos:  position{line: 40, col: 7, offset: 1815},
							name: "_",
						},
					},
					&zeroOrOneExpr{
						pos: position{line: 40, col: 10, offset: 1818},
						expr: &ruleRefExpr{
							pos:  position{line: 40, col: 10, offset: 1818},
							name: "Comment",
						},
					},
					&choiceExpr{
						pos: position{line: 40, col: 20, offset: 1828},
						alternatives: []interface{}{
							&litMatcher{
								pos:        position{line: 40, col: 20, offset: 1828},
								val:        "\r\n",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 40, col: 29, offset: 1837},
								val:        "\n\r",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 40, col: 38, offset: 1846},
								val:        "\r",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 40, col: 45, offset: 1853},
								val:        "\n",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 40, col: 52, offset: 1860},
								name: "EOF",
							},
						},
					},
				},
			},
		},
		{
			name: "Comment",
			pos:  position{line: 42, col: 1, offset: 1868},
			expr: &seqExpr{
				pos: position{line: 42, col: 11, offset: 1880},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 42, col: 11, offset: 1880},
						val:        "//",
						ignoreCase: false,
					},
					&zeroOrMoreExpr{
						pos: position{line: 42, col: 16, offset: 1885},
						expr: &charClassMatcher{
							pos:        position{line: 42, col: 16, offset: 1885},
							val:        "[^\\r\\n]",
							chars:      []rune{'\r', '\n'},
							ignoreCase: false,
							inverted:   true,
						},
					},
				},
			},
		},
		{
			name: "EOF",
			pos:  position{line: 44, col: 1, offset: 1897},
			expr: &notExpr{
				pos: position{line: 44, col: 7, offset: 1905},
				expr: &anyMatcher{
					line: 44, col: 8, offset: 1906,
				},
			},
		},
		{
			name: "INDENTATION",
			pos:  position{line: 46, col: 1, offset: 1911},
			expr: &seqExpr{
				pos: position{line: 46, col: 15, offset: 1927},
				exprs: []interface{}{
					&labeledExpr{
						pos:   position{line: 46, col: 15, offset: 1927},
						label: "spaces",
						expr: &zeroOrMoreExpr{
							pos: position{line: 46, col: 22, offset: 1934},
							expr: &litMatcher{
								pos:        position{line: 46, col: 22, offset: 1934},
								val:        " ",
								ignoreCase: false,
							},
						},
					},
					&andCodeExpr{
						pos: position{line: 46, col: 27, offset: 1939},
						run: (*parser).callonINDENTATION5,
					},
				},
			},
		},
		{
			name: "INDENT",
			pos:  position{line: 48, col: 1, offset: 2017},
			expr: &andCodeExpr{
				pos: position{line: 48, col: 10, offset: 2028},
				run: (*parser).callonINDENT1,
			},
		},
		{
			name: "DEDENT",
			pos:  position{line: 50, col: 1, offset: 2111},
			expr: &andCodeExpr{
				pos: position{line: 50, col: 10, offset: 2122},
				run: (*parser).callonDEDENT1,
			},
		},
	},
}

func (c *current) onInput3() (bool, error) {
	c.state["Indentation"] = 0
	return true, nil
}

func (p *parser) callonInput3() (bool, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onInput3()
}

func (c *current) onInput1(s, r interface{}) (interface{}, error) {
	return newProgramNode(s.(StatementsNode), r.(ReturnNode))
}

func (p *parser) callonInput1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onInput1(stack["s"], stack["r"])
}

func (c *current) onStatements1(s interface{}) (interface{}, error) {
	return newStatementsNode(s)
}

func (p *parser) callonStatements1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onStatements1(stack["s"])
}

func (c *current) onLine1(s interface{}) (interface{}, error) {
	return s, nil
}

func (p *parser) callonLine1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onLine1(stack["s"])
}

func (c *current) onReturnOp1(arg interface{}) (interface{}, error) {
	return newReturnNode(arg.(IdentifierNode))
}

func (p *parser) callonReturnOp1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReturnOp1(stack["arg"])
}

func (c *current) onStatement2(s interface{}) (interface{}, error) {
	return s.(AssignmentNode), nil
}

func (p *parser) callonStatement2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onStatement2(stack["s"])
}

func (c *current) onStatement7(arg, s interface{}) (interface{}, error) {
	return newIfNode(arg.(LogicalExpressionNode), s.(StatementsNode))
}

func (p *parser) callonStatement7() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onStatement7(stack["arg"], stack["s"])
}

func (c *current) onAssignment1(lvalue, rvalue interface{}) (interface{}, error) {
	return newAssignmentNode(lvalue.(IdentifierNode), rvalue.(AdditiveExpressionNode))
}

func (p *parser) callonAssignment1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onAssignment1(stack["lvalue"], stack["rvalue"])
}

func (c *current) onLogicalExpression1(arg interface{}) (interface{}, error) {
	return newLogicalExpressionNode(arg.(PrimaryExpressionNode))
}

func (p *parser) callonLogicalExpression1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onLogicalExpression1(stack["arg"])
}

func (c *current) onAdditiveExpression1(arg, rest interface{}) (interface{}, error) {
	return newAdditiveExpressionNode(arg.(PrimaryExpressionNode), rest)
}

func (p *parser) callonAdditiveExpression1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onAdditiveExpression1(stack["arg"], stack["rest"])
}

func (c *current) onPrimaryExpression1(arg interface{}) (interface{}, error) {
	return newPrimaryExpressionNode(arg)
}

func (p *parser) callonPrimaryExpression1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPrimaryExpression1(stack["arg"])
}

func (c *current) onInteger1() (interface{}, error) {
	return newIntegerNode(string(c.text))
}

func (p *parser) callonInteger1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onInteger1()
}

func (c *current) onIdentifier1() (interface{}, error) {
	return newIdentifierNode(string(c.text))
}

func (p *parser) callonIdentifier1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onIdentifier1()
}

func (c *current) onAddOp1() (interface{}, error) {
	return string(c.text), nil
}

func (p *parser) callonAddOp1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onAddOp1()
}

func (c *current) onINDENTATION5(spaces interface{}) (bool, error) {
	return len(toIfaceSlice(spaces)) == c.state["Indentation"].(int), nil
}

func (p *parser) callonINDENTATION5() (bool, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onINDENTATION5(stack["spaces"])
}

func (c *current) onINDENT1() (bool, error) {
	c.state["Indentation"] = c.state["Indentation"].(int) + 4
	return true, nil
}

func (p *parser) callonINDENT1() (bool, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onINDENT1()
}

func (c *current) onDEDENT1() (bool, error) {
	c.state["Indentation"] = c.state["Indentation"].(int) - 4
	return true, nil
}

func (p *parser) callonDEDENT1() (bool, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onDEDENT1()
}

var (
	// errNoRule is returned when the grammar to parse has no rule.
	errNoRule = errors.New("grammar has no rule")

	// errInvalidEncoding is returned when the source is not properly
	// utf8-encoded.
	errInvalidEncoding = errors.New("invalid encoding")
)

// Option is a function that can set an option on the parser. It returns
// the previous setting as an Option.
type Option func(*parser) Option

// Debug creates an Option to set the debug flag to b. When set to true,
// debugging information is printed to stdout while parsing.
//
// The default is false.
func Debug(b bool) Option {
	return func(p *parser) Option {
		old := p.debug
		p.debug = b
		return Debug(old)
	}
}

// Memoize creates an Option to set the memoize flag to b. When set to true,
// the parser will cache all results so each expression is evaluated only
// once. This guarantees linear parsing time even for pathological cases,
// at the expense of more memory and slower times for typical cases.
//
// The default is false.
func Memoize(b bool) Option {
	return func(p *parser) Option {
		old := p.memoize
		p.memoize = b
		return Memoize(old)
	}
}

// Recover creates an Option to set the recover flag to b. When set to
// true, this causes the parser to recover from panics and convert it
// to an error. Setting it to false can be useful while debugging to
// access the full stack trace.
//
// The default is true.
func Recover(b bool) Option {
	return func(p *parser) Option {
		old := p.recover
		p.recover = b
		return Recover(old)
	}
}

// GlobalStore creates an Option to set a key to a certain value in
// the globalStore.
func GlobalStore(key string, value interface{}) Option {
	return func(p *parser) Option {
		old := p.cur.globalStore[key]
		p.cur.globalStore[key] = value
		return GlobalStore(key, old)
	}
}

// ParseFile parses the file identified by filename.
func ParseFile(filename string, opts ...Option) (i interface{}, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil {
			err = closeErr
		}
	}()
	return ParseReader(filename, f, opts...)
}

// ParseReader parses the data from r using filename as information in the
// error messages.
func ParseReader(filename string, r io.Reader, opts ...Option) (interface{}, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return Parse(filename, b, opts...)
}

// Parse parses the data from b using filename as information in the
// error messages.
func Parse(filename string, b []byte, opts ...Option) (interface{}, error) {
	return newParser(filename, b, opts...).parse(g)
}

// position records a position in the text.
type position struct {
	line, col, offset int
}

func (p position) String() string {
	return fmt.Sprintf("%d:%d [%d]", p.line, p.col, p.offset)
}

// savepoint stores all state required to go back to this point in the
// parser.
type savepoint struct {
	position
	rn rune
	w  int
}

type current struct {
	pos  position // start position of the match
	text []byte   // raw text of the match

	// the state allows the parser to store arbitrary values and rollback them if needed
	state statedict
	// the globalStore allows the parser to store arbitrary values
	globalStore map[string]interface{}
}

type statedict map[string]interface{}

// the AST types...

type grammar struct {
	pos   position
	rules []*rule
}

type rule struct {
	pos         position
	name        string
	displayName string
	expr        interface{}
}

type choiceExpr struct {
	pos          position
	alternatives []interface{}
}

type actionExpr struct {
	pos  position
	expr interface{}
	run  func(*parser) (interface{}, error)
}

type seqExpr struct {
	pos   position
	exprs []interface{}
}

type labeledExpr struct {
	pos   position
	label string
	expr  interface{}
}

type expr struct {
	pos  position
	expr interface{}
}

type andExpr expr
type notExpr expr
type zeroOrOneExpr expr
type zeroOrMoreExpr expr
type oneOrMoreExpr expr

type ruleRefExpr struct {
	pos  position
	name string
}

type andCodeExpr struct {
	pos position
	run func(*parser) (bool, error)
}

type notCodeExpr struct {
	pos position
	run func(*parser) (bool, error)
}

type litMatcher struct {
	pos        position
	val        string
	ignoreCase bool
}

type charClassMatcher struct {
	pos             position
	val             string
	basicLatinChars [128]bool
	chars           []rune
	ranges          []rune
	classes         []*unicode.RangeTable
	ignoreCase      bool
	inverted        bool
}

type anyMatcher position

// errList cumulates the errors found by the parser.
type errList []error

func (e *errList) add(err error) {
	*e = append(*e, err)
}

func (e errList) err() error {
	if len(e) == 0 {
		return nil
	}
	e.dedupe()
	return e
}

func (e *errList) dedupe() {
	var cleaned []error
	set := make(map[string]bool)
	for _, err := range *e {
		if msg := err.Error(); !set[msg] {
			set[msg] = true
			cleaned = append(cleaned, err)
		}
	}
	*e = cleaned
}

func (e errList) Error() string {
	switch len(e) {
	case 0:
		return ""
	case 1:
		return e[0].Error()
	default:
		var buf bytes.Buffer

		for i, err := range e {
			if i > 0 {
				buf.WriteRune('\n')
			}
			buf.WriteString(err.Error())
		}
		return buf.String()
	}
}

// parserError wraps an error with a prefix indicating the rule in which
// the error occurred. The original error is stored in the Inner field.
type parserError struct {
	Inner    error
	pos      position
	prefix   string
	expected []string
}

// Error returns the error message.
func (p *parserError) Error() string {
	return p.prefix + ": " + p.Inner.Error()
}

// newParser creates a parser with the specified input source and options.
func newParser(filename string, b []byte, opts ...Option) *parser {
	p := &parser{
		filename: filename,
		errs:     new(errList),
		data:     b,
		pt:       savepoint{position: position{line: 1}},
		recover:  true,
		cur: current{
			state:       make(statedict),
			globalStore: make(map[string]interface{}),
		},
		maxFailPos:      position{col: 1, line: 1},
		maxFailExpected: make([]string, 0, 20),
	}
	p.setOptions(opts)
	return p
}

// setOptions applies the options to the parser.
func (p *parser) setOptions(opts []Option) {
	for _, opt := range opts {
		opt(p)
	}
}

type resultTuple struct {
	v   interface{}
	b   bool
	end savepoint
}

type parser struct {
	filename string
	pt       savepoint
	cur      current

	data []byte
	errs *errList

	depth   int
	recover bool
	debug   bool

	memoize bool
	// memoization table for the packrat algorithm:
	// map[offset in source] map[expression or rule] {value, match}
	memo map[int]map[interface{}]resultTuple

	// rules table, maps the rule identifier to the rule node
	rules map[string]*rule
	// variables stack, map of label to value
	vstack []map[string]interface{}
	// rule stack, allows identification of the current rule in errors
	rstack []*rule

	// stats
	exprCnt int

	// parse fail
	maxFailPos            position
	maxFailExpected       []string
	maxFailInvertExpected bool
}

// push a variable set on the vstack.
func (p *parser) pushV() {
	if cap(p.vstack) == len(p.vstack) {
		// create new empty slot in the stack
		p.vstack = append(p.vstack, nil)
	} else {
		// slice to 1 more
		p.vstack = p.vstack[:len(p.vstack)+1]
	}

	// get the last args set
	m := p.vstack[len(p.vstack)-1]
	if m != nil && len(m) == 0 {
		// empty map, all good
		return
	}

	m = make(map[string]interface{})
	p.vstack[len(p.vstack)-1] = m
}

// pop a variable set from the vstack.
func (p *parser) popV() {
	// if the map is not empty, clear it
	m := p.vstack[len(p.vstack)-1]
	if len(m) > 0 {
		// GC that map
		p.vstack[len(p.vstack)-1] = nil
	}
	p.vstack = p.vstack[:len(p.vstack)-1]
}

func (p *parser) print(prefix, s string) string {
	if !p.debug {
		return s
	}

	fmt.Printf("%s %d:%d:%d: %s [%#U]\n",
		prefix, p.pt.line, p.pt.col, p.pt.offset, s, p.pt.rn)
	return s
}

func (p *parser) in(s string) string {
	p.depth++
	return p.print(strings.Repeat(" ", p.depth)+">", s)
}

func (p *parser) out(s string) string {
	p.depth--
	return p.print(strings.Repeat(" ", p.depth)+"<", s)
}

func (p *parser) addErr(err error) {
	p.addErrAt(err, p.pt.position, []string{})
}

func (p *parser) addErrAt(err error, pos position, expected []string) {
	var buf bytes.Buffer
	if p.filename != "" {
		buf.WriteString(p.filename)
	}
	if buf.Len() > 0 {
		buf.WriteString(":")
	}
	buf.WriteString(fmt.Sprintf("%d:%d (%d)", pos.line, pos.col, pos.offset))
	if len(p.rstack) > 0 {
		if buf.Len() > 0 {
			buf.WriteString(": ")
		}
		rule := p.rstack[len(p.rstack)-1]
		if rule.displayName != "" {
			buf.WriteString("rule " + rule.displayName)
		} else {
			buf.WriteString("rule " + rule.name)
		}
	}
	pe := &parserError{Inner: err, pos: pos, prefix: buf.String(), expected: expected}
	p.errs.add(pe)
}

func (p *parser) failAt(fail bool, pos position, want string) {
	// process fail if parsing fails and not inverted or parsing succeeds and invert is set
	if fail == p.maxFailInvertExpected {
		if pos.offset < p.maxFailPos.offset {
			return
		}

		if pos.offset > p.maxFailPos.offset {
			p.maxFailPos = pos
			p.maxFailExpected = p.maxFailExpected[:0]
		}

		if p.maxFailInvertExpected {
			want = "!" + want
		}
		p.maxFailExpected = append(p.maxFailExpected, want)
	}
}

// read advances the parser to the next rune.
func (p *parser) read() {
	p.pt.offset += p.pt.w
	rn, n := utf8.DecodeRune(p.data[p.pt.offset:])
	p.pt.rn = rn
	p.pt.w = n
	p.pt.col++
	if rn == '\n' {
		p.pt.line++
		p.pt.col = 0
	}

	if rn == utf8.RuneError {
		if n == 1 {
			p.addErr(errInvalidEncoding)
		}
	}
}

// restore parser position to the savepoint pt.
func (p *parser) restore(pt savepoint) {
	if p.debug {
		defer p.out(p.in("restore"))
	}
	if pt.offset == p.pt.offset {
		return
	}
	p.pt = pt
}

// copy and return parser current state.
func (p *parser) copyState() (state statedict) {
	if p.debug {
		defer p.out(p.in("copyState"))
	}
	state = make(statedict)
	for k, v := range p.cur.state {
		state[k] = v
	}
	return state
}

// restore parser current state to the state statedict.
func (p *parser) restoreState(state statedict) {
	if p.debug {
		defer p.out(p.in("restoreState"))
	}
	p.cur.state = state
}

// get the slice of bytes from the savepoint start to the current position.
func (p *parser) sliceFrom(start savepoint) []byte {
	return p.data[start.position.offset:p.pt.position.offset]
}

func (p *parser) getMemoized(node interface{}) (resultTuple, bool) {
	if len(p.memo) == 0 {
		return resultTuple{}, false
	}
	m := p.memo[p.pt.offset]
	if len(m) == 0 {
		return resultTuple{}, false
	}
	res, ok := m[node]
	return res, ok
}

func (p *parser) setMemoized(pt savepoint, node interface{}, tuple resultTuple) {
	if p.memo == nil {
		p.memo = make(map[int]map[interface{}]resultTuple)
	}
	m := p.memo[pt.offset]
	if m == nil {
		m = make(map[interface{}]resultTuple)
		p.memo[pt.offset] = m
	}
	m[node] = tuple
}

func (p *parser) buildRulesTable(g *grammar) {
	p.rules = make(map[string]*rule, len(g.rules))
	for _, r := range g.rules {
		p.rules[r.name] = r
	}
}

func (p *parser) parse(g *grammar) (val interface{}, err error) {
	if len(g.rules) == 0 {
		p.addErr(errNoRule)
		return nil, p.errs.err()
	}

	// TODO : not super critical but this could be generated
	p.buildRulesTable(g)

	if p.recover {
		// panic can be used in action code to stop parsing immediately
		// and return the panic as an error.
		defer func() {
			if e := recover(); e != nil {
				if p.debug {
					defer p.out(p.in("panic handler"))
				}
				val = nil
				switch e := e.(type) {
				case error:
					p.addErr(e)
				default:
					p.addErr(fmt.Errorf("%v", e))
				}
				err = p.errs.err()
			}
		}()
	}

	// start rule is rule [0]
	p.read() // advance to first rune
	val, ok := p.parseRule(g.rules[0])
	if !ok {
		if len(*p.errs) == 0 {
			// If parsing fails, but no errors have been recorded, the expected values
			// for the farthest parser position are returned as error.
			maxFailExpectedMap := make(map[string]struct{}, len(p.maxFailExpected))
			for _, v := range p.maxFailExpected {
				maxFailExpectedMap[v] = struct{}{}
			}
			expected := make([]string, 0, len(maxFailExpectedMap))
			eof := false
			if _, ok := maxFailExpectedMap["!."]; ok {
				delete(maxFailExpectedMap, "!.")
				eof = true
			}
			for k := range maxFailExpectedMap {
				expected = append(expected, k)
			}
			sort.Strings(expected)
			if eof {
				expected = append(expected, "EOF")
			}
			p.addErrAt(errors.New("no match found, expected: "+listJoin(expected, ", ", "or")), p.maxFailPos, expected)
		}
		return nil, p.errs.err()
	}
	return val, p.errs.err()
}

func listJoin(list []string, sep string, lastSep string) string {
	switch len(list) {
	case 0:
		return ""
	case 1:
		return list[0]
	default:
		return fmt.Sprintf("%s %s %s", strings.Join(list[:len(list)-1], sep), lastSep, list[len(list)-1])
	}
}

func (p *parser) parseRule(rule *rule) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseRule " + rule.name))
	}

	if p.memoize {
		res, ok := p.getMemoized(rule)
		if ok {
			p.restore(res.end)
			return res.v, res.b
		}
	}

	start := p.pt
	p.rstack = append(p.rstack, rule)
	p.pushV()
	val, ok := p.parseExpr(rule.expr)
	p.popV()
	p.rstack = p.rstack[:len(p.rstack)-1]
	if ok && p.debug {
		p.print(strings.Repeat(" ", p.depth)+"MATCH", string(p.sliceFrom(start)))
	}

	if p.memoize {
		p.setMemoized(start, rule, resultTuple{val, ok, p.pt})
	}
	return val, ok
}

func (p *parser) parseExpr(expr interface{}) (interface{}, bool) {
	var pt savepoint

	if p.memoize {
		res, ok := p.getMemoized(expr)
		if ok {
			p.restore(res.end)
			return res.v, res.b
		}
		pt = p.pt
	}

	p.exprCnt++
	var val interface{}
	var ok bool
	switch expr := expr.(type) {
	case *actionExpr:
		val, ok = p.parseActionExpr(expr)
	case *andCodeExpr:
		val, ok = p.parseAndCodeExpr(expr)
	case *andExpr:
		val, ok = p.parseAndExpr(expr)
	case *anyMatcher:
		val, ok = p.parseAnyMatcher(expr)
	case *charClassMatcher:
		val, ok = p.parseCharClassMatcher(expr)
	case *choiceExpr:
		val, ok = p.parseChoiceExpr(expr)
	case *labeledExpr:
		val, ok = p.parseLabeledExpr(expr)
	case *litMatcher:
		val, ok = p.parseLitMatcher(expr)
	case *notCodeExpr:
		val, ok = p.parseNotCodeExpr(expr)
	case *notExpr:
		val, ok = p.parseNotExpr(expr)
	case *oneOrMoreExpr:
		val, ok = p.parseOneOrMoreExpr(expr)
	case *ruleRefExpr:
		val, ok = p.parseRuleRefExpr(expr)
	case *seqExpr:
		val, ok = p.parseSeqExpr(expr)
	case *zeroOrMoreExpr:
		val, ok = p.parseZeroOrMoreExpr(expr)
	case *zeroOrOneExpr:
		val, ok = p.parseZeroOrOneExpr(expr)
	default:
		panic(fmt.Sprintf("unknown expression type %T", expr))
	}
	if p.memoize {
		p.setMemoized(pt, expr, resultTuple{val, ok, p.pt})
	}
	return val, ok
}

func (p *parser) parseActionExpr(act *actionExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseActionExpr"))
	}

	start := p.pt
	val, ok := p.parseExpr(act.expr)
	if ok {
		p.cur.pos = start.position
		p.cur.text = p.sliceFrom(start)
		actVal, err := act.run(p)
		if err != nil {
			p.addErrAt(err, start.position, []string{})
		}
		val = actVal
	}
	if ok && p.debug {
		p.print(strings.Repeat(" ", p.depth)+"MATCH", string(p.sliceFrom(start)))
	}
	return val, ok
}

func (p *parser) parseAndCodeExpr(and *andCodeExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAndCodeExpr"))
	}

	ok, err := and.run(p)
	if err != nil {
		p.addErr(err)
	}
	return nil, ok
}

func (p *parser) parseAndExpr(and *andExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAndExpr"))
	}

	pt := p.pt
	state := p.copyState()
	p.pushV()
	_, ok := p.parseExpr(and.expr)
	p.popV()
	p.restore(pt)
	p.restoreState(state)
	return nil, ok
}

func (p *parser) parseAnyMatcher(any *anyMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAnyMatcher"))
	}

	if p.pt.rn != utf8.RuneError {
		start := p.pt
		p.read()
		p.failAt(true, start.position, ".")
		return p.sliceFrom(start), true
	}
	p.failAt(false, p.pt.position, ".")
	return nil, false
}

func (p *parser) parseCharClassMatcher(chr *charClassMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseCharClassMatcher"))
	}

	cur := p.pt.rn
	start := p.pt

	// can't match EOF
	if cur == utf8.RuneError {
		p.failAt(false, start.position, chr.val)
		return nil, false
	}

	if chr.ignoreCase {
		cur = unicode.ToLower(cur)
	}

	// try to match in the list of available chars
	for _, rn := range chr.chars {
		if rn == cur {
			if chr.inverted {
				p.failAt(false, start.position, chr.val)
				return nil, false
			}
			p.read()
			p.failAt(true, start.position, chr.val)
			return p.sliceFrom(start), true
		}
	}

	// try to match in the list of ranges
	for i := 0; i < len(chr.ranges); i += 2 {
		if cur >= chr.ranges[i] && cur <= chr.ranges[i+1] {
			if chr.inverted {
				p.failAt(false, start.position, chr.val)
				return nil, false
			}
			p.read()
			p.failAt(true, start.position, chr.val)
			return p.sliceFrom(start), true
		}
	}

	// try to match in the list of Unicode classes
	for _, cl := range chr.classes {
		if unicode.Is(cl, cur) {
			if chr.inverted {
				p.failAt(false, start.position, chr.val)
				return nil, false
			}
			p.read()
			p.failAt(true, start.position, chr.val)
			return p.sliceFrom(start), true
		}
	}

	if chr.inverted {
		p.read()
		p.failAt(true, start.position, chr.val)
		return p.sliceFrom(start), true
	}
	p.failAt(false, start.position, chr.val)
	return nil, false
}

func (p *parser) parseChoiceExpr(ch *choiceExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseChoiceExpr"))
	}

	for _, alt := range ch.alternatives {
		p.pushV()
		val, ok := p.parseExpr(alt)
		p.popV()
		if ok {
			return val, ok
		}
	}
	return nil, false
}

func (p *parser) parseLabeledExpr(lab *labeledExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseLabeledExpr"))
	}

	p.pushV()
	val, ok := p.parseExpr(lab.expr)
	p.popV()
	if ok && lab.label != "" {
		m := p.vstack[len(p.vstack)-1]
		m[lab.label] = val
	}
	return val, ok
}

func (p *parser) parseLitMatcher(lit *litMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseLitMatcher"))
	}

	ignoreCase := ""
	if lit.ignoreCase {
		ignoreCase = "i"
	}
	val := fmt.Sprintf("%q%s", lit.val, ignoreCase)
	start := p.pt
	for _, want := range lit.val {
		cur := p.pt.rn
		if lit.ignoreCase {
			cur = unicode.ToLower(cur)
		}
		if cur != want {
			p.failAt(false, start.position, val)
			p.restore(start)
			return nil, false
		}
		p.read()
	}
	p.failAt(true, start.position, val)
	return p.sliceFrom(start), true
}

func (p *parser) parseNotCodeExpr(not *notCodeExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseNotCodeExpr"))
	}

	ok, err := not.run(p)
	if err != nil {
		p.addErr(err)
	}
	return nil, !ok
}

func (p *parser) parseNotExpr(not *notExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseNotExpr"))
	}

	pt := p.pt
	state := p.copyState()
	p.pushV()
	p.maxFailInvertExpected = !p.maxFailInvertExpected
	_, ok := p.parseExpr(not.expr)
	p.maxFailInvertExpected = !p.maxFailInvertExpected
	p.popV()
	p.restore(pt)
	p.restoreState(state)
	return nil, !ok
}

func (p *parser) parseOneOrMoreExpr(expr *oneOrMoreExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseOneOrMoreExpr"))
	}

	var vals []interface{}

	for {
		p.pushV()
		val, ok := p.parseExpr(expr.expr)
		p.popV()
		if !ok {
			if len(vals) == 0 {
				// did not match once, no match
				return nil, false
			}
			return vals, true
		}
		vals = append(vals, val)
	}
}

func (p *parser) parseRuleRefExpr(ref *ruleRefExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseRuleRefExpr " + ref.name))
	}

	if ref.name == "" {
		panic(fmt.Sprintf("%s: invalid rule: missing name", ref.pos))
	}

	rule := p.rules[ref.name]
	if rule == nil {
		p.addErr(fmt.Errorf("undefined rule: %s", ref.name))
		return nil, false
	}
	return p.parseRule(rule)
}

func (p *parser) parseSeqExpr(seq *seqExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseSeqExpr"))
	}

	vals := make([]interface{}, 0, len(seq.exprs))

	pt := p.pt
	state := p.copyState()
	for _, expr := range seq.exprs {
		val, ok := p.parseExpr(expr)
		if !ok {
			p.restore(pt)
			p.restoreState(state)
			return nil, false
		}
		vals = append(vals, val)
	}
	return vals, true
}

func (p *parser) parseZeroOrMoreExpr(expr *zeroOrMoreExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseZeroOrMoreExpr"))
	}

	var vals []interface{}

	for {
		p.pushV()
		val, ok := p.parseExpr(expr.expr)
		p.popV()
		if !ok {
			return vals, true
		}
		vals = append(vals, val)
	}
}

func (p *parser) parseZeroOrOneExpr(expr *zeroOrOneExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseZeroOrOneExpr"))
	}

	p.pushV()
	val, _ := p.parseExpr(expr.expr)
	p.popV()
	// whether it matched or not, consider it a match
	return val, true
}
