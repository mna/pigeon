package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/mna/pigeon/ast"
)

var g = &grammar{
	rules: []*rule{
		{
			name: "Grammar",
			pos:  position{line: 5, col: 1, offset: 22},
			expr: &actionExpr{
				pos: position{line: 5, col: 11, offset: 34},
				run: (*parser).callonGrammar1,
				expr: &seqExpr{
					pos: position{line: 5, col: 11, offset: 34},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 5, col: 11, offset: 34},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 5, col: 14, offset: 37},
							label: "initializer",
							expr: &zeroOrOneExpr{
								pos: position{line: 5, col: 28, offset: 51},
								expr: &seqExpr{
									pos: position{line: 5, col: 28, offset: 51},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 5, col: 28, offset: 51},
											name: "Initializer",
										},
										&ruleRefExpr{
											pos:  position{line: 5, col: 40, offset: 63},
											name: "__",
										},
									},
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 5, col: 46, offset: 69},
							label: "rules",
							expr: &oneOrMoreExpr{
								pos: position{line: 5, col: 54, offset: 77},
								expr: &seqExpr{
									pos: position{line: 5, col: 54, offset: 77},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 5, col: 54, offset: 77},
											name: "Rule",
										},
										&ruleRefExpr{
											pos:  position{line: 5, col: 59, offset: 82},
											name: "__",
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
			name: "Initializer",
			pos:  position{line: 24, col: 1, offset: 544},
			expr: &actionExpr{
				pos: position{line: 24, col: 15, offset: 560},
				run: (*parser).callonInitializer1,
				expr: &seqExpr{
					pos: position{line: 24, col: 15, offset: 560},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 24, col: 15, offset: 560},
							label: "code",
							expr: &ruleRefExpr{
								pos:  position{line: 24, col: 20, offset: 565},
								name: "CodeBlock",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 24, col: 30, offset: 575},
							name: "EOS",
						},
					},
				},
			},
		},
		{
			name: "Rule",
			pos:  position{line: 28, col: 1, offset: 609},
			expr: &actionExpr{
				pos: position{line: 28, col: 8, offset: 618},
				run: (*parser).callonRule1,
				expr: &seqExpr{
					pos: position{line: 28, col: 8, offset: 618},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 28, col: 8, offset: 618},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 28, col: 13, offset: 623},
								name: "IdentifierName",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 28, col: 28, offset: 638},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 28, col: 31, offset: 641},
							label: "display",
							expr: &zeroOrOneExpr{
								pos: position{line: 28, col: 41, offset: 651},
								expr: &seqExpr{
									pos: position{line: 28, col: 41, offset: 651},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 28, col: 41, offset: 651},
											name: "StringLiteral",
										},
										&ruleRefExpr{
											pos:  position{line: 28, col: 55, offset: 665},
											name: "__",
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 28, col: 61, offset: 671},
							name: "RuleDefOp",
						},
						&ruleRefExpr{
							pos:  position{line: 28, col: 71, offset: 681},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 28, col: 74, offset: 684},
							label: "expr",
							expr: &ruleRefExpr{
								pos:  position{line: 28, col: 79, offset: 689},
								name: "Expression",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 28, col: 90, offset: 700},
							name: "EOS",
						},
					},
				},
			},
		},
		{
			name: "Expression",
			pos:  position{line: 41, col: 1, offset: 997},
			expr: &ruleRefExpr{
				pos:  position{line: 41, col: 14, offset: 1012},
				name: "ChoiceExpr",
			},
		},
		{
			name: "ChoiceExpr",
			pos:  position{line: 43, col: 1, offset: 1026},
			expr: &actionExpr{
				pos: position{line: 43, col: 14, offset: 1041},
				run: (*parser).callonChoiceExpr1,
				expr: &seqExpr{
					pos: position{line: 43, col: 14, offset: 1041},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 43, col: 14, offset: 1041},
							label: "first",
							expr: &ruleRefExpr{
								pos:  position{line: 43, col: 20, offset: 1047},
								name: "ActionExpr",
							},
						},
						&labeledExpr{
							pos:   position{line: 43, col: 31, offset: 1058},
							label: "rest",
							expr: &zeroOrMoreExpr{
								pos: position{line: 43, col: 38, offset: 1065},
								expr: &seqExpr{
									pos: position{line: 43, col: 38, offset: 1065},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 43, col: 38, offset: 1065},
											name: "__",
										},
										&litMatcher{
											pos:        position{line: 43, col: 41, offset: 1068},
											val:        "/",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 43, col: 45, offset: 1072},
											name: "__",
										},
										&ruleRefExpr{
											pos:  position{line: 43, col: 48, offset: 1075},
											name: "ActionExpr",
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
			name: "ActionExpr",
			pos:  position{line: 58, col: 1, offset: 1495},
			expr: &actionExpr{
				pos: position{line: 58, col: 14, offset: 1510},
				run: (*parser).callonActionExpr1,
				expr: &seqExpr{
					pos: position{line: 58, col: 14, offset: 1510},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 58, col: 14, offset: 1510},
							label: "expr",
							expr: &ruleRefExpr{
								pos:  position{line: 58, col: 19, offset: 1515},
								name: "SeqExpr",
							},
						},
						&labeledExpr{
							pos:   position{line: 58, col: 27, offset: 1523},
							label: "code",
							expr: &zeroOrOneExpr{
								pos: position{line: 58, col: 34, offset: 1530},
								expr: &seqExpr{
									pos: position{line: 58, col: 34, offset: 1530},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 58, col: 34, offset: 1530},
											name: "__",
										},
										&ruleRefExpr{
											pos:  position{line: 58, col: 37, offset: 1533},
											name: "CodeBlock",
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
			name: "SeqExpr",
			pos:  position{line: 72, col: 1, offset: 1813},
			expr: &actionExpr{
				pos: position{line: 72, col: 11, offset: 1825},
				run: (*parser).callonSeqExpr1,
				expr: &seqExpr{
					pos: position{line: 72, col: 11, offset: 1825},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 72, col: 11, offset: 1825},
							label: "first",
							expr: &ruleRefExpr{
								pos:  position{line: 72, col: 17, offset: 1831},
								name: "LabeledExpr",
							},
						},
						&labeledExpr{
							pos:   position{line: 72, col: 29, offset: 1843},
							label: "rest",
							expr: &zeroOrMoreExpr{
								pos: position{line: 72, col: 36, offset: 1850},
								expr: &seqExpr{
									pos: position{line: 72, col: 36, offset: 1850},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 72, col: 36, offset: 1850},
											name: "__",
										},
										&ruleRefExpr{
											pos:  position{line: 72, col: 39, offset: 1853},
											name: "LabeledExpr",
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
			name: "LabeledExpr",
			pos:  position{line: 85, col: 1, offset: 2217},
			expr: &choiceExpr{
				pos: position{line: 85, col: 15, offset: 2233},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 85, col: 15, offset: 2233},
						run: (*parser).callonLabeledExpr2,
						expr: &seqExpr{
							pos: position{line: 85, col: 15, offset: 2233},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 85, col: 15, offset: 2233},
									label: "label",
									expr: &ruleRefExpr{
										pos:  position{line: 85, col: 21, offset: 2239},
										name: "Identifier",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 85, col: 32, offset: 2250},
									name: "__",
								},
								&litMatcher{
									pos:        position{line: 85, col: 35, offset: 2253},
									val:        ":",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 85, col: 39, offset: 2257},
									name: "__",
								},
								&labeledExpr{
									pos:   position{line: 85, col: 42, offset: 2260},
									label: "expr",
									expr: &ruleRefExpr{
										pos:  position{line: 85, col: 47, offset: 2265},
										name: "PrefixedExpr",
									},
								},
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 91, col: 5, offset: 2444},
						name: "PrefixedExpr",
					},
				},
			},
		},
		{
			name: "PrefixedExpr",
			pos:  position{line: 93, col: 1, offset: 2460},
			expr: &choiceExpr{
				pos: position{line: 93, col: 16, offset: 2477},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 93, col: 16, offset: 2477},
						run: (*parser).callonPrefixedExpr2,
						expr: &seqExpr{
							pos: position{line: 93, col: 16, offset: 2477},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 93, col: 16, offset: 2477},
									label: "op",
									expr: &ruleRefExpr{
										pos:  position{line: 93, col: 19, offset: 2480},
										name: "PrefixedOp",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 93, col: 30, offset: 2491},
									name: "__",
								},
								&labeledExpr{
									pos:   position{line: 93, col: 33, offset: 2494},
									label: "expr",
									expr: &ruleRefExpr{
										pos:  position{line: 93, col: 38, offset: 2499},
										name: "SuffixedExpr",
									},
								},
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 104, col: 5, offset: 2792},
						name: "SuffixedExpr",
					},
				},
			},
		},
		{
			name: "PrefixedOp",
			pos:  position{line: 106, col: 1, offset: 2808},
			expr: &actionExpr{
				pos: position{line: 106, col: 14, offset: 2823},
				run: (*parser).callonPrefixedOp1,
				expr: &choiceExpr{
					pos: position{line: 106, col: 16, offset: 2825},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 106, col: 16, offset: 2825},
							val:        "&",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 106, col: 22, offset: 2831},
							val:        "!",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "SuffixedExpr",
			pos:  position{line: 110, col: 1, offset: 2877},
			expr: &choiceExpr{
				pos: position{line: 110, col: 16, offset: 2894},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 110, col: 16, offset: 2894},
						run: (*parser).callonSuffixedExpr2,
						expr: &seqExpr{
							pos: position{line: 110, col: 16, offset: 2894},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 110, col: 16, offset: 2894},
									label: "expr",
									expr: &ruleRefExpr{
										pos:  position{line: 110, col: 21, offset: 2899},
										name: "PrimaryExpr",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 110, col: 33, offset: 2911},
									name: "__",
								},
								&labeledExpr{
									pos:   position{line: 110, col: 36, offset: 2914},
									label: "op",
									expr: &ruleRefExpr{
										pos:  position{line: 110, col: 39, offset: 2917},
										name: "SuffixedOp",
									},
								},
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 129, col: 5, offset: 3466},
						name: "PrimaryExpr",
					},
				},
			},
		},
		{
			name: "SuffixedOp",
			pos:  position{line: 131, col: 1, offset: 3482},
			expr: &actionExpr{
				pos: position{line: 131, col: 14, offset: 3497},
				run: (*parser).callonSuffixedOp1,
				expr: &choiceExpr{
					pos: position{line: 131, col: 16, offset: 3499},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 131, col: 16, offset: 3499},
							val:        "?",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 131, col: 22, offset: 3505},
							val:        "*",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 131, col: 28, offset: 3511},
							val:        "+",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "PrimaryExpr",
			pos:  position{line: 135, col: 1, offset: 3557},
			expr: &choiceExpr{
				pos: position{line: 135, col: 15, offset: 3573},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 135, col: 15, offset: 3573},
						name: "LitMatcher",
					},
					&ruleRefExpr{
						pos:  position{line: 135, col: 28, offset: 3586},
						name: "CharClassMatcher",
					},
					&ruleRefExpr{
						pos:  position{line: 135, col: 47, offset: 3605},
						name: "AnyMatcher",
					},
					&ruleRefExpr{
						pos:  position{line: 135, col: 60, offset: 3618},
						name: "RuleRefExpr",
					},
					&ruleRefExpr{
						pos:  position{line: 135, col: 74, offset: 3632},
						name: "SemanticPredExpr",
					},
					&actionExpr{
						pos: position{line: 135, col: 93, offset: 3651},
						run: (*parser).callonPrimaryExpr7,
						expr: &seqExpr{
							pos: position{line: 135, col: 93, offset: 3651},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 135, col: 93, offset: 3651},
									val:        "(",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 135, col: 97, offset: 3655},
									name: "__",
								},
								&labeledExpr{
									pos:   position{line: 135, col: 100, offset: 3658},
									label: "expr",
									expr: &ruleRefExpr{
										pos:  position{line: 135, col: 105, offset: 3663},
										name: "Expression",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 135, col: 116, offset: 3674},
									name: "__",
								},
								&litMatcher{
									pos:        position{line: 135, col: 119, offset: 3677},
									val:        ")",
									ignoreCase: false,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "RuleRefExpr",
			pos:  position{line: 138, col: 1, offset: 3709},
			expr: &actionExpr{
				pos: position{line: 138, col: 15, offset: 3725},
				run: (*parser).callonRuleRefExpr1,
				expr: &seqExpr{
					pos: position{line: 138, col: 15, offset: 3725},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 138, col: 15, offset: 3725},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 138, col: 20, offset: 3730},
								name: "IdentifierName",
							},
						},
						&notExpr{
							pos: position{line: 138, col: 35, offset: 3745},
							expr: &seqExpr{
								pos: position{line: 138, col: 38, offset: 3748},
								exprs: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 138, col: 38, offset: 3748},
										name: "__",
									},
									&zeroOrOneExpr{
										pos: position{line: 138, col: 43, offset: 3753},
										expr: &seqExpr{
											pos: position{line: 138, col: 43, offset: 3753},
											exprs: []interface{}{
												&ruleRefExpr{
													pos:  position{line: 138, col: 43, offset: 3753},
													name: "StringLiteral",
												},
												&ruleRefExpr{
													pos:  position{line: 138, col: 57, offset: 3767},
													name: "__",
												},
											},
										},
									},
									&ruleRefExpr{
										pos:  position{line: 138, col: 63, offset: 3773},
										name: "RuleDefOp",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "SemanticPredExpr",
			pos:  position{line: 143, col: 1, offset: 3894},
			expr: &actionExpr{
				pos: position{line: 143, col: 20, offset: 3915},
				run: (*parser).callonSemanticPredExpr1,
				expr: &seqExpr{
					pos: position{line: 143, col: 20, offset: 3915},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 143, col: 20, offset: 3915},
							label: "op",
							expr: &ruleRefExpr{
								pos:  position{line: 143, col: 23, offset: 3918},
								name: "SemanticPredOp",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 143, col: 38, offset: 3933},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 143, col: 41, offset: 3936},
							label: "code",
							expr: &ruleRefExpr{
								pos:  position{line: 143, col: 46, offset: 3941},
								name: "CodeBlock",
							},
						},
					},
				},
			},
		},
		{
			name: "SemanticPredOp",
			pos:  position{line: 154, col: 1, offset: 4229},
			expr: &actionExpr{
				pos: position{line: 154, col: 18, offset: 4248},
				run: (*parser).callonSemanticPredOp1,
				expr: &choiceExpr{
					pos: position{line: 154, col: 20, offset: 4250},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 154, col: 20, offset: 4250},
							val:        "&",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 154, col: 26, offset: 4256},
							val:        "!",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "RuleDefOp",
			pos:  position{line: 158, col: 1, offset: 4302},
			expr: &choiceExpr{
				pos: position{line: 158, col: 13, offset: 4316},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 158, col: 13, offset: 4316},
						val:        "=",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 158, col: 19, offset: 4322},
						val:        "<-",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 158, col: 26, offset: 4329},
						val:        "←",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 158, col: 37, offset: 4340},
						val:        "⟵",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "SourceChar",
			pos:  position{line: 160, col: 1, offset: 4352},
			expr: &anyMatcher{
				line: 160, col: 14, offset: 4367,
			},
		},
		{
			name: "Comment",
			pos:  position{line: 161, col: 1, offset: 4370},
			expr: &choiceExpr{
				pos: position{line: 161, col: 11, offset: 4382},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 161, col: 11, offset: 4382},
						name: "MultiLineComment",
					},
					&ruleRefExpr{
						pos:  position{line: 161, col: 30, offset: 4401},
						name: "SingleLineComment",
					},
				},
			},
		},
		{
			name: "MultiLineComment",
			pos:  position{line: 162, col: 1, offset: 4420},
			expr: &seqExpr{
				pos: position{line: 162, col: 20, offset: 4441},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 162, col: 20, offset: 4441},
						val:        "/*",
						ignoreCase: false,
					},
					&zeroOrMoreExpr{
						pos: position{line: 162, col: 27, offset: 4448},
						expr: &seqExpr{
							pos: position{line: 162, col: 27, offset: 4448},
							exprs: []interface{}{
								&notExpr{
									pos: position{line: 162, col: 27, offset: 4448},
									expr: &litMatcher{
										pos:        position{line: 162, col: 28, offset: 4449},
										val:        "*/",
										ignoreCase: false,
									},
								},
								&ruleRefExpr{
									pos:  position{line: 162, col: 33, offset: 4454},
									name: "SourceChar",
								},
							},
						},
					},
					&litMatcher{
						pos:        position{line: 162, col: 47, offset: 4468},
						val:        "*/",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "MultiLineCommentNoLineTerminator",
			pos:  position{line: 163, col: 1, offset: 4474},
			expr: &seqExpr{
				pos: position{line: 163, col: 36, offset: 4511},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 163, col: 36, offset: 4511},
						val:        "/*",
						ignoreCase: false,
					},
					&zeroOrMoreExpr{
						pos: position{line: 163, col: 43, offset: 4518},
						expr: &seqExpr{
							pos: position{line: 163, col: 43, offset: 4518},
							exprs: []interface{}{
								&notExpr{
									pos: position{line: 163, col: 43, offset: 4518},
									expr: &choiceExpr{
										pos: position{line: 163, col: 46, offset: 4521},
										alternatives: []interface{}{
											&litMatcher{
												pos:        position{line: 163, col: 46, offset: 4521},
												val:        "*/",
												ignoreCase: false,
											},
											&ruleRefExpr{
												pos:  position{line: 163, col: 53, offset: 4528},
												name: "EOL",
											},
										},
									},
								},
								&ruleRefExpr{
									pos:  position{line: 163, col: 59, offset: 4534},
									name: "SourceChar",
								},
							},
						},
					},
					&litMatcher{
						pos:        position{line: 163, col: 73, offset: 4548},
						val:        "*/",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "SingleLineComment",
			pos:  position{line: 164, col: 1, offset: 4554},
			expr: &seqExpr{
				pos: position{line: 164, col: 21, offset: 4576},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 164, col: 21, offset: 4576},
						val:        "//",
						ignoreCase: false,
					},
					&zeroOrMoreExpr{
						pos: position{line: 164, col: 28, offset: 4583},
						expr: &seqExpr{
							pos: position{line: 164, col: 28, offset: 4583},
							exprs: []interface{}{
								&notExpr{
									pos: position{line: 164, col: 28, offset: 4583},
									expr: &ruleRefExpr{
										pos:  position{line: 164, col: 29, offset: 4584},
										name: "EOL",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 164, col: 33, offset: 4588},
									name: "SourceChar",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Identifier",
			pos:  position{line: 166, col: 1, offset: 4605},
			expr: &ruleRefExpr{
				pos:  position{line: 166, col: 14, offset: 4620},
				name: "IdentifierName",
			},
		},
		{
			name: "IdentifierName",
			pos:  position{line: 167, col: 1, offset: 4636},
			expr: &actionExpr{
				pos: position{line: 167, col: 18, offset: 4655},
				run: (*parser).callonIdentifierName1,
				expr: &seqExpr{
					pos: position{line: 167, col: 18, offset: 4655},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 167, col: 18, offset: 4655},
							name: "IdentifierStart",
						},
						&zeroOrMoreExpr{
							pos: position{line: 167, col: 34, offset: 4671},
							expr: &ruleRefExpr{
								pos:  position{line: 167, col: 34, offset: 4671},
								name: "IdentifierPart",
							},
						},
					},
				},
			},
		},
		{
			name: "IdentifierStart",
			pos:  position{line: 170, col: 1, offset: 4756},
			expr: &charClassMatcher{
				pos:        position{line: 170, col: 19, offset: 4776},
				val:        "[a-z_]i",
				chars:      []rune{'_'},
				ranges:     []rune{'a', 'z'},
				ignoreCase: true,
				inverted:   false,
			},
		},
		{
			name: "IdentifierPart",
			pos:  position{line: 171, col: 1, offset: 4785},
			expr: &choiceExpr{
				pos: position{line: 171, col: 18, offset: 4804},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 171, col: 18, offset: 4804},
						name: "IdentifierStart",
					},
					&charClassMatcher{
						pos:        position{line: 171, col: 36, offset: 4822},
						val:        "[0-9]",
						ranges:     []rune{'0', '9'},
						ignoreCase: false,
						inverted:   false,
					},
				},
			},
		},
		{
			name: "LitMatcher",
			pos:  position{line: 173, col: 1, offset: 4831},
			expr: &actionExpr{
				pos: position{line: 173, col: 14, offset: 4846},
				run: (*parser).callonLitMatcher1,
				expr: &seqExpr{
					pos: position{line: 173, col: 14, offset: 4846},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 173, col: 14, offset: 4846},
							label: "lit",
							expr: &ruleRefExpr{
								pos:  position{line: 173, col: 18, offset: 4850},
								name: "StringLiteral",
							},
						},
						&labeledExpr{
							pos:   position{line: 173, col: 32, offset: 4864},
							label: "ignore",
							expr: &zeroOrOneExpr{
								pos: position{line: 173, col: 39, offset: 4871},
								expr: &litMatcher{
									pos:        position{line: 173, col: 39, offset: 4871},
									val:        "i",
									ignoreCase: false,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "StringLiteral",
			pos:  position{line: 183, col: 1, offset: 5107},
			expr: &actionExpr{
				pos: position{line: 183, col: 17, offset: 5125},
				run: (*parser).callonStringLiteral1,
				expr: &choiceExpr{
					pos: position{line: 183, col: 19, offset: 5127},
					alternatives: []interface{}{
						&seqExpr{
							pos: position{line: 183, col: 19, offset: 5127},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 183, col: 19, offset: 5127},
									val:        "\"",
									ignoreCase: false,
								},
								&zeroOrMoreExpr{
									pos: position{line: 183, col: 23, offset: 5131},
									expr: &ruleRefExpr{
										pos:  position{line: 183, col: 23, offset: 5131},
										name: "DoubleStringChar",
									},
								},
								&litMatcher{
									pos:        position{line: 183, col: 41, offset: 5149},
									val:        "\"",
									ignoreCase: false,
								},
							},
						},
						&seqExpr{
							pos: position{line: 183, col: 47, offset: 5155},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 183, col: 47, offset: 5155},
									val:        "'",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 183, col: 51, offset: 5159},
									name: "SingleStringChar",
								},
								&litMatcher{
									pos:        position{line: 183, col: 68, offset: 5176},
									val:        "'",
									ignoreCase: false,
								},
							},
						},
						&seqExpr{
							pos: position{line: 183, col: 74, offset: 5182},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 183, col: 74, offset: 5182},
									val:        "`",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 183, col: 78, offset: 5186},
									name: "RawStringChar",
								},
								&litMatcher{
									pos:        position{line: 183, col: 92, offset: 5200},
									val:        "`",
									ignoreCase: false,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "DoubleStringChar",
			pos:  position{line: 186, col: 1, offset: 5274},
			expr: &choiceExpr{
				pos: position{line: 186, col: 20, offset: 5295},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 186, col: 20, offset: 5295},
						exprs: []interface{}{
							&notExpr{
								pos: position{line: 186, col: 20, offset: 5295},
								expr: &choiceExpr{
									pos: position{line: 186, col: 23, offset: 5298},
									alternatives: []interface{}{
										&litMatcher{
											pos:        position{line: 186, col: 23, offset: 5298},
											val:        "\"",
											ignoreCase: false,
										},
										&litMatcher{
											pos:        position{line: 186, col: 29, offset: 5304},
											val:        "\\",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 186, col: 36, offset: 5311},
											name: "EOL",
										},
									},
								},
							},
							&ruleRefExpr{
								pos:  position{line: 186, col: 42, offset: 5317},
								name: "SourceChar",
							},
						},
					},
					&seqExpr{
						pos: position{line: 186, col: 55, offset: 5330},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 186, col: 55, offset: 5330},
								val:        "\\",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 186, col: 60, offset: 5335},
								name: "DoubleStringEscape",
							},
						},
					},
				},
			},
		},
		{
			name: "SingleStringChar",
			pos:  position{line: 187, col: 1, offset: 5355},
			expr: &choiceExpr{
				pos: position{line: 187, col: 20, offset: 5376},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 187, col: 20, offset: 5376},
						exprs: []interface{}{
							&notExpr{
								pos: position{line: 187, col: 20, offset: 5376},
								expr: &choiceExpr{
									pos: position{line: 187, col: 23, offset: 5379},
									alternatives: []interface{}{
										&litMatcher{
											pos:        position{line: 187, col: 23, offset: 5379},
											val:        "'",
											ignoreCase: false,
										},
										&litMatcher{
											pos:        position{line: 187, col: 29, offset: 5385},
											val:        "\\",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 187, col: 36, offset: 5392},
											name: "EOL",
										},
									},
								},
							},
							&ruleRefExpr{
								pos:  position{line: 187, col: 42, offset: 5398},
								name: "SourceChar",
							},
						},
					},
					&seqExpr{
						pos: position{line: 187, col: 55, offset: 5411},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 187, col: 55, offset: 5411},
								val:        "\\",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 187, col: 60, offset: 5416},
								name: "SingleStringEscape",
							},
						},
					},
				},
			},
		},
		{
			name: "RawStringChar",
			pos:  position{line: 188, col: 1, offset: 5436},
			expr: &seqExpr{
				pos: position{line: 188, col: 17, offset: 5454},
				exprs: []interface{}{
					&notExpr{
						pos: position{line: 188, col: 17, offset: 5454},
						expr: &litMatcher{
							pos:        position{line: 188, col: 18, offset: 5455},
							val:        "`",
							ignoreCase: false,
						},
					},
					&ruleRefExpr{
						pos:  position{line: 188, col: 22, offset: 5459},
						name: "SourceChar",
					},
				},
			},
		},
		{
			name: "DoubleStringEscape",
			pos:  position{line: 190, col: 1, offset: 5473},
			expr: &choiceExpr{
				pos: position{line: 190, col: 22, offset: 5496},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 190, col: 22, offset: 5496},
						val:        "'",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 190, col: 28, offset: 5502},
						name: "CommonEscapeSequence",
					},
				},
			},
		},
		{
			name: "SingleStringEscape",
			pos:  position{line: 191, col: 1, offset: 5524},
			expr: &choiceExpr{
				pos: position{line: 191, col: 22, offset: 5547},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 191, col: 22, offset: 5547},
						val:        "\"",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 191, col: 28, offset: 5553},
						name: "CommonEscapeSequence",
					},
				},
			},
		},
		{
			name: "CommonEscapeSequence",
			pos:  position{line: 193, col: 1, offset: 5577},
			expr: &choiceExpr{
				pos: position{line: 193, col: 24, offset: 5602},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 193, col: 24, offset: 5602},
						name: "SingleCharEscape",
					},
					&ruleRefExpr{
						pos:  position{line: 193, col: 43, offset: 5621},
						name: "OctalEscape",
					},
					&ruleRefExpr{
						pos:  position{line: 193, col: 57, offset: 5635},
						name: "HexEscape",
					},
					&ruleRefExpr{
						pos:  position{line: 193, col: 69, offset: 5647},
						name: "LongUnicodeEscape",
					},
					&ruleRefExpr{
						pos:  position{line: 193, col: 89, offset: 5667},
						name: "ShortUnicodeEscape",
					},
				},
			},
		},
		{
			name: "SingleCharEscape",
			pos:  position{line: 194, col: 1, offset: 5687},
			expr: &choiceExpr{
				pos: position{line: 194, col: 20, offset: 5708},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 194, col: 20, offset: 5708},
						val:        "a",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 194, col: 26, offset: 5714},
						val:        "b",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 194, col: 32, offset: 5720},
						val:        "n",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 194, col: 38, offset: 5726},
						val:        "f",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 194, col: 44, offset: 5732},
						val:        "r",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 194, col: 50, offset: 5738},
						val:        "t",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 194, col: 56, offset: 5744},
						val:        "v",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 194, col: 62, offset: 5750},
						val:        "\\",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "OctalEscape",
			pos:  position{line: 195, col: 1, offset: 5756},
			expr: &seqExpr{
				pos: position{line: 195, col: 15, offset: 5772},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 195, col: 15, offset: 5772},
						name: "OctalDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 195, col: 26, offset: 5783},
						name: "OctalDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 195, col: 37, offset: 5794},
						name: "OctalDigit",
					},
				},
			},
		},
		{
			name: "HexEscape",
			pos:  position{line: 196, col: 1, offset: 5806},
			expr: &seqExpr{
				pos: position{line: 196, col: 13, offset: 5820},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 196, col: 13, offset: 5820},
						val:        "x",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 196, col: 17, offset: 5824},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 196, col: 26, offset: 5833},
						name: "HexDigit",
					},
				},
			},
		},
		{
			name: "LongUnicodeEscape",
			pos:  position{line: 197, col: 1, offset: 5843},
			expr: &seqExpr{
				pos: position{line: 197, col: 21, offset: 5865},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 197, col: 21, offset: 5865},
						val:        "U",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 197, col: 25, offset: 5869},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 197, col: 34, offset: 5878},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 197, col: 43, offset: 5887},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 197, col: 52, offset: 5896},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 197, col: 61, offset: 5905},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 197, col: 70, offset: 5914},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 197, col: 79, offset: 5923},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 197, col: 88, offset: 5932},
						name: "HexDigit",
					},
				},
			},
		},
		{
			name: "ShortUnicodeEscape",
			pos:  position{line: 198, col: 1, offset: 5942},
			expr: &seqExpr{
				pos: position{line: 198, col: 22, offset: 5965},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 198, col: 22, offset: 5965},
						val:        "u",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 198, col: 26, offset: 5969},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 198, col: 35, offset: 5978},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 198, col: 44, offset: 5987},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 198, col: 53, offset: 5996},
						name: "HexDigit",
					},
				},
			},
		},
		{
			name: "OctalDigit",
			pos:  position{line: 200, col: 1, offset: 6008},
			expr: &charClassMatcher{
				pos:        position{line: 200, col: 14, offset: 6023},
				val:        "[0-7]",
				ranges:     []rune{'0', '7'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "DecimalDigit",
			pos:  position{line: 201, col: 1, offset: 6030},
			expr: &charClassMatcher{
				pos:        position{line: 201, col: 16, offset: 6047},
				val:        "[0-9]",
				ranges:     []rune{'0', '9'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "HexDigit",
			pos:  position{line: 202, col: 1, offset: 6054},
			expr: &charClassMatcher{
				pos:        position{line: 202, col: 12, offset: 6067},
				val:        "[0-9a-f]i",
				ranges:     []rune{'0', '9', 'a', 'f'},
				ignoreCase: true,
				inverted:   false,
			},
		},
		{
			name: "CharClassMatcher",
			pos:  position{line: 204, col: 1, offset: 6080},
			expr: &actionExpr{
				pos: position{line: 204, col: 20, offset: 6101},
				run: (*parser).callonCharClassMatcher1,
				expr: &seqExpr{
					pos: position{line: 204, col: 20, offset: 6101},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 204, col: 20, offset: 6101},
							val:        "[",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 204, col: 26, offset: 6107},
							expr: &choiceExpr{
								pos: position{line: 204, col: 26, offset: 6107},
								alternatives: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 204, col: 26, offset: 6107},
										name: "ClassCharRange",
									},
									&ruleRefExpr{
										pos:  position{line: 204, col: 43, offset: 6124},
										name: "ClassChar",
									},
									&seqExpr{
										pos: position{line: 204, col: 55, offset: 6136},
										exprs: []interface{}{
											&litMatcher{
												pos:        position{line: 204, col: 55, offset: 6136},
												val:        "\\",
												ignoreCase: false,
											},
											&ruleRefExpr{
												pos:  position{line: 204, col: 60, offset: 6141},
												name: "UnicodeClassEscape",
											},
										},
									},
								},
							},
						},
						&litMatcher{
							pos:        position{line: 204, col: 82, offset: 6163},
							val:        "]",
							ignoreCase: false,
						},
						&zeroOrOneExpr{
							pos: position{line: 204, col: 86, offset: 6167},
							expr: &litMatcher{
								pos:        position{line: 204, col: 86, offset: 6167},
								val:        "i",
								ignoreCase: false,
							},
						},
					},
				},
			},
		},
		{
			name: "ClassCharRange",
			pos:  position{line: 209, col: 1, offset: 6277},
			expr: &seqExpr{
				pos: position{line: 209, col: 18, offset: 6296},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 209, col: 18, offset: 6296},
						name: "ClassChar",
					},
					&litMatcher{
						pos:        position{line: 209, col: 28, offset: 6306},
						val:        "-",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 209, col: 32, offset: 6310},
						name: "ClassChar",
					},
				},
			},
		},
		{
			name: "ClassChar",
			pos:  position{line: 210, col: 1, offset: 6321},
			expr: &choiceExpr{
				pos: position{line: 210, col: 13, offset: 6335},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 210, col: 13, offset: 6335},
						exprs: []interface{}{
							&notExpr{
								pos: position{line: 210, col: 13, offset: 6335},
								expr: &choiceExpr{
									pos: position{line: 210, col: 16, offset: 6338},
									alternatives: []interface{}{
										&litMatcher{
											pos:        position{line: 210, col: 16, offset: 6338},
											val:        "]",
											ignoreCase: false,
										},
										&litMatcher{
											pos:        position{line: 210, col: 22, offset: 6344},
											val:        "\\",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 210, col: 29, offset: 6351},
											name: "EOL",
										},
									},
								},
							},
							&ruleRefExpr{
								pos:  position{line: 210, col: 35, offset: 6357},
								name: "SourceChar",
							},
						},
					},
					&seqExpr{
						pos: position{line: 210, col: 48, offset: 6370},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 210, col: 48, offset: 6370},
								val:        "\\",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 210, col: 53, offset: 6375},
								name: "CharClassEscape",
							},
						},
					},
				},
			},
		},
		{
			name: "CharClassEscape",
			pos:  position{line: 211, col: 1, offset: 6392},
			expr: &choiceExpr{
				pos: position{line: 211, col: 19, offset: 6412},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 211, col: 19, offset: 6412},
						val:        "]",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 211, col: 25, offset: 6418},
						name: "CommonEscapeSequence",
					},
				},
			},
		},
		{
			name: "UnicodeClassEscape",
			pos:  position{line: 213, col: 1, offset: 6442},
			expr: &seqExpr{
				pos: position{line: 213, col: 22, offset: 6465},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 213, col: 22, offset: 6465},
						val:        "p",
						ignoreCase: false,
					},
					&choiceExpr{
						pos: position{line: 213, col: 28, offset: 6471},
						alternatives: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 213, col: 28, offset: 6471},
								name: "SingleCharUnicodeClass",
							},
							&seqExpr{
								pos: position{line: 213, col: 53, offset: 6496},
								exprs: []interface{}{
									&litMatcher{
										pos:        position{line: 213, col: 53, offset: 6496},
										val:        "{",
										ignoreCase: false,
									},
									&ruleRefExpr{
										pos:  position{line: 213, col: 57, offset: 6500},
										name: "UnicodeClass",
									},
									&litMatcher{
										pos:        position{line: 213, col: 70, offset: 6513},
										val:        "}",
										ignoreCase: false,
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "SingleCharUnicodeClass",
			pos:  position{line: 214, col: 1, offset: 6520},
			expr: &charClassMatcher{
				pos:        position{line: 214, col: 26, offset: 6547},
				val:        "[LMNCPZS]",
				chars:      []rune{'L', 'M', 'N', 'C', 'P', 'Z', 'S'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "UnicodeClass",
			pos:  position{line: 215, col: 1, offset: 6558},
			expr: &oneOrMoreExpr{
				pos: position{line: 215, col: 16, offset: 6575},
				expr: &charClassMatcher{
					pos:        position{line: 215, col: 16, offset: 6575},
					val:        "[a-z_]i",
					chars:      []rune{'_'},
					ranges:     []rune{'a', 'z'},
					ignoreCase: true,
					inverted:   false,
				},
			},
		},
		{
			name: "AnyMatcher",
			pos:  position{line: 217, col: 1, offset: 6587},
			expr: &actionExpr{
				pos: position{line: 217, col: 14, offset: 6602},
				run: (*parser).callonAnyMatcher1,
				expr: &litMatcher{
					pos:        position{line: 217, col: 14, offset: 6602},
					val:        ".",
					ignoreCase: false,
				},
			},
		},
		{
			name: "CodeBlock",
			pos:  position{line: 222, col: 1, offset: 6682},
			expr: &actionExpr{
				pos: position{line: 222, col: 13, offset: 6696},
				run: (*parser).callonCodeBlock1,
				expr: &seqExpr{
					pos: position{line: 222, col: 13, offset: 6696},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 222, col: 13, offset: 6696},
							val:        "{",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 222, col: 17, offset: 6700},
							name: "Code",
						},
						&litMatcher{
							pos:        position{line: 222, col: 22, offset: 6705},
							val:        "}",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Code",
			pos:  position{line: 228, col: 1, offset: 6809},
			expr: &zeroOrMoreExpr{
				pos: position{line: 228, col: 10, offset: 6820},
				expr: &choiceExpr{
					pos: position{line: 228, col: 10, offset: 6820},
					alternatives: []interface{}{
						&oneOrMoreExpr{
							pos: position{line: 228, col: 12, offset: 6822},
							expr: &seqExpr{
								pos: position{line: 228, col: 12, offset: 6822},
								exprs: []interface{}{
									&notExpr{
										pos: position{line: 228, col: 12, offset: 6822},
										expr: &charClassMatcher{
											pos:        position{line: 228, col: 13, offset: 6823},
											val:        "[{}]",
											chars:      []rune{'{', '}'},
											ignoreCase: false,
											inverted:   false,
										},
									},
									&ruleRefExpr{
										pos:  position{line: 228, col: 18, offset: 6828},
										name: "SourceChar",
									},
								},
							},
						},
						&seqExpr{
							pos: position{line: 228, col: 34, offset: 6844},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 228, col: 34, offset: 6844},
									val:        "{",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 228, col: 38, offset: 6848},
									name: "Code",
								},
								&litMatcher{
									pos:        position{line: 228, col: 43, offset: 6853},
									val:        "}",
									ignoreCase: false,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "__",
			pos:  position{line: 230, col: 1, offset: 6863},
			expr: &zeroOrMoreExpr{
				pos: position{line: 230, col: 8, offset: 6872},
				expr: &choiceExpr{
					pos: position{line: 230, col: 8, offset: 6872},
					alternatives: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 230, col: 8, offset: 6872},
							name: "Whitespace",
						},
						&ruleRefExpr{
							pos:  position{line: 230, col: 21, offset: 6885},
							name: "EOL",
						},
						&ruleRefExpr{
							pos:  position{line: 230, col: 27, offset: 6891},
							name: "Comment",
						},
					},
				},
			},
		},
		{
			name: "_",
			pos:  position{line: 231, col: 1, offset: 6903},
			expr: &zeroOrMoreExpr{
				pos: position{line: 231, col: 7, offset: 6911},
				expr: &choiceExpr{
					pos: position{line: 231, col: 7, offset: 6911},
					alternatives: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 231, col: 7, offset: 6911},
							name: "Whitespace",
						},
						&ruleRefExpr{
							pos:  position{line: 231, col: 20, offset: 6924},
							name: "MultiLineCommentNoLineTerminator",
						},
					},
				},
			},
		},
		{
			name: "Whitespace",
			pos:  position{line: 233, col: 1, offset: 6963},
			expr: &charClassMatcher{
				pos:        position{line: 233, col: 14, offset: 6978},
				val:        "[ \\t\\r]",
				chars:      []rune{' ', '\t', '\r'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "EOL",
			pos:  position{line: 234, col: 1, offset: 6987},
			expr: &litMatcher{
				pos:        position{line: 234, col: 7, offset: 6995},
				val:        "\n",
				ignoreCase: false,
			},
		},
		{
			name: "EOS",
			pos:  position{line: 235, col: 1, offset: 7001},
			expr: &choiceExpr{
				pos: position{line: 235, col: 7, offset: 7009},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 235, col: 7, offset: 7009},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 235, col: 7, offset: 7009},
								name: "__",
							},
							&litMatcher{
								pos:        position{line: 235, col: 10, offset: 7012},
								val:        ";",
								ignoreCase: false,
							},
						},
					},
					&seqExpr{
						pos: position{line: 235, col: 16, offset: 7018},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 235, col: 16, offset: 7018},
								name: "_",
							},
							&zeroOrOneExpr{
								pos: position{line: 235, col: 18, offset: 7020},
								expr: &ruleRefExpr{
									pos:  position{line: 235, col: 18, offset: 7020},
									name: "SingleLineComment",
								},
							},
							&ruleRefExpr{
								pos:  position{line: 235, col: 37, offset: 7039},
								name: "EOL",
							},
						},
					},
					&seqExpr{
						pos: position{line: 235, col: 43, offset: 7045},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 235, col: 43, offset: 7045},
								name: "__",
							},
							&ruleRefExpr{
								pos:  position{line: 235, col: 46, offset: 7048},
								name: "EOF",
							},
						},
					},
				},
			},
		},
		{
			name: "EOF",
			pos:  position{line: 237, col: 1, offset: 7055},
			expr: &notExpr{
				pos: position{line: 237, col: 7, offset: 7063},
				expr: &anyMatcher{
					line: 237, col: 8, offset: 7064,
				},
			},
		},
	},
}

func (c *current) onGrammar1(initializer, rules interface{}) (interface{}, error) {

	pos := c.astPos()

	// create the grammar, assign its initializer
	g := ast.NewGrammar(pos)
	initSlice := toIfaceSlice(initializer)
	if len(initSlice) > 0 {
		g.Init = initSlice[0].(*ast.CodeBlock)
	}

	rulesSlice := toIfaceSlice(rules)
	g.Rules = make([]*ast.Rule, len(rulesSlice))
	for i, duo := range rulesSlice {
		g.Rules[i] = duo.([]interface{})[0].(*ast.Rule)
	}

	return g, nil
}

func (p *parser) callonGrammar1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onGrammar1(stack["initializer"], stack["rules"])
}

func (c *current) onInitializer1(code interface{}) (interface{}, error) {

	return code, nil
}

func (p *parser) callonInitializer1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onInitializer1(stack["code"])
}

func (c *current) onRule1(name, display, expr interface{}) (interface{}, error) {

	pos := c.astPos()

	rule := ast.NewRule(pos, name.(*ast.Identifier))
	displaySlice := toIfaceSlice(display)
	if len(displaySlice) > 0 {
		rule.DisplayName = displaySlice[0].(*ast.StringLit)
	}
	rule.Expr = expr.(ast.Expression)

	return rule, nil
}

func (p *parser) callonRule1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onRule1(stack["name"], stack["display"], stack["expr"])
}

func (c *current) onChoiceExpr1(first, rest interface{}) (interface{}, error) {

	restSlice := toIfaceSlice(rest)
	if len(restSlice) == 0 {
		return first, nil
	}

	pos := c.astPos()
	choice := ast.NewChoiceExpr(pos)
	choice.Alternatives = []ast.Expression{first.(ast.Expression)}
	for _, sl := range restSlice {
		choice.Alternatives = append(choice.Alternatives, sl.([]interface{})[3].(ast.Expression))
	}
	return choice, nil
}

func (p *parser) callonChoiceExpr1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onChoiceExpr1(stack["first"], stack["rest"])
}

func (c *current) onActionExpr1(expr, code interface{}) (interface{}, error) {

	if code == nil {
		return expr, nil
	}

	pos := c.astPos()
	act := ast.NewActionExpr(pos)
	act.Expr = expr.(ast.Expression)
	codeSlice := toIfaceSlice(code)
	act.Code = codeSlice[1].(*ast.CodeBlock)

	return act, nil
}

func (p *parser) callonActionExpr1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onActionExpr1(stack["expr"], stack["code"])
}

func (c *current) onSeqExpr1(first, rest interface{}) (interface{}, error) {

	restSlice := toIfaceSlice(rest)
	if len(restSlice) == 0 {
		return first, nil
	}
	seq := ast.NewSeqExpr(c.astPos())
	seq.Exprs = []ast.Expression{first.(ast.Expression)}
	for _, sl := range restSlice {
		seq.Exprs = append(seq.Exprs, sl.([]interface{})[1].(ast.Expression))
	}
	return seq, nil
}

func (p *parser) callonSeqExpr1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onSeqExpr1(stack["first"], stack["rest"])
}

func (c *current) onLabeledExpr2(label, expr interface{}) (interface{}, error) {

	pos := c.astPos()
	lab := ast.NewLabeledExpr(pos)
	lab.Label = label.(*ast.Identifier)
	lab.Expr = expr.(ast.Expression)
	return lab, nil
}

func (p *parser) callonLabeledExpr2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onLabeledExpr2(stack["label"], stack["expr"])
}

func (c *current) onPrefixedExpr2(op, expr interface{}) (interface{}, error) {

	pos := c.astPos()
	opStr := op.(string)
	if opStr == "&" {
		and := ast.NewAndExpr(pos)
		and.Expr = expr.(ast.Expression)
		return and, nil
	}
	not := ast.NewNotExpr(pos)
	not.Expr = expr.(ast.Expression)
	return not, nil
}

func (p *parser) callonPrefixedExpr2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPrefixedExpr2(stack["op"], stack["expr"])
}

func (c *current) onPrefixedOp1() (interface{}, error) {

	return string(c.text), nil
}

func (p *parser) callonPrefixedOp1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPrefixedOp1()
}

func (c *current) onSuffixedExpr2(expr, op interface{}) (interface{}, error) {

	pos := c.astPos()
	opStr := op.(string)
	switch opStr {
	case "?":
		zero := ast.NewZeroOrOneExpr(pos)
		zero.Expr = expr.(ast.Expression)
		return zero, nil
	case "*":
		zero := ast.NewZeroOrMoreExpr(pos)
		zero.Expr = expr.(ast.Expression)
		return zero, nil
	case "+":
		one := ast.NewOneOrMoreExpr(pos)
		one.Expr = expr.(ast.Expression)
		return one, nil
	default:
		return nil, errors.New("unknown operator: " + opStr)
	}
}

func (p *parser) callonSuffixedExpr2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onSuffixedExpr2(stack["expr"], stack["op"])
}

func (c *current) onSuffixedOp1() (interface{}, error) {

	return string(c.text), nil
}

func (p *parser) callonSuffixedOp1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onSuffixedOp1()
}

func (c *current) onPrimaryExpr7(expr interface{}) (interface{}, error) {

	return expr, nil
}

func (p *parser) callonPrimaryExpr7() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPrimaryExpr7(stack["expr"])
}

func (c *current) onRuleRefExpr1(name interface{}) (interface{}, error) {

	ref := ast.NewRuleRefExpr(c.astPos())
	ref.Name = name.(*ast.Identifier)
	return ref, nil
}

func (p *parser) callonRuleRefExpr1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onRuleRefExpr1(stack["name"])
}

func (c *current) onSemanticPredExpr1(op, code interface{}) (interface{}, error) {

	opStr := op.(string)
	if opStr == "&" {
		and := ast.NewAndCodeExpr(c.astPos())
		and.Code = code.(*ast.CodeBlock)
		return and, nil
	}
	not := ast.NewNotCodeExpr(c.astPos())
	not.Code = code.(*ast.CodeBlock)
	return not, nil
}

func (p *parser) callonSemanticPredExpr1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onSemanticPredExpr1(stack["op"], stack["code"])
}

func (c *current) onSemanticPredOp1() (interface{}, error) {

	return string(c.text), nil
}

func (p *parser) callonSemanticPredOp1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onSemanticPredOp1()
}

func (c *current) onIdentifierName1() (interface{}, error) {

	return ast.NewIdentifier(c.astPos(), string(c.text)), nil
}

func (p *parser) callonIdentifierName1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onIdentifierName1()
}

func (c *current) onLitMatcher1(lit, ignore interface{}) (interface{}, error) {

	rawStr := lit.(*ast.StringLit).Val
	s, err := strconv.Unquote(rawStr)
	if err != nil {
		return nil, err
	}
	m := ast.NewLitMatcher(c.astPos(), s)
	m.IgnoreCase = ignore != nil
	return m, nil
}

func (p *parser) callonLitMatcher1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onLitMatcher1(stack["lit"], stack["ignore"])
}

func (c *current) onStringLiteral1() (interface{}, error) {

	return ast.NewStringLit(c.astPos(), string(c.text)), nil
}

func (p *parser) callonStringLiteral1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onStringLiteral1()
}

func (c *current) onCharClassMatcher1() (interface{}, error) {

	pos := c.astPos()
	cc := ast.NewCharClassMatcher(pos, string(c.text))
	return cc, nil
}

func (p *parser) callonCharClassMatcher1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onCharClassMatcher1()
}

func (c *current) onAnyMatcher1() (interface{}, error) {

	any := ast.NewAnyMatcher(c.astPos(), ".")
	return any, nil
}

func (p *parser) callonAnyMatcher1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onAnyMatcher1()
}

func (c *current) onCodeBlock1() (interface{}, error) {

	pos := c.astPos()
	cb := ast.NewCodeBlock(pos, string(c.text))
	return cb, nil
}

func (p *parser) callonCodeBlock1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onCodeBlock1()
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
	p.cur.state = make(statedict)
	for k, v := range state {
		p.cur.state[k] = v
	}
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
