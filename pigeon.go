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
								pos: position{line: 5, col: 26, offset: 49},
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
								pos: position{line: 5, col: 52, offset: 75},
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
						&ruleRefExpr{
							pos:  position{line: 5, col: 65, offset: 88},
							name: "EOF",
						},
					},
				},
			},
		},
		{
			name: "Initializer",
			pos:  position{line: 24, col: 1, offset: 548},
			expr: &actionExpr{
				pos: position{line: 24, col: 15, offset: 564},
				run: (*parser).callonInitializer1,
				expr: &seqExpr{
					pos: position{line: 24, col: 15, offset: 564},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 24, col: 15, offset: 564},
							label: "code",
							expr: &ruleRefExpr{
								pos:  position{line: 24, col: 20, offset: 569},
								name: "CodeBlock",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 24, col: 30, offset: 579},
							name: "EOS",
						},
					},
				},
			},
		},
		{
			name: "Rule",
			pos:  position{line: 28, col: 1, offset: 613},
			expr: &actionExpr{
				pos: position{line: 28, col: 8, offset: 622},
				run: (*parser).callonRule1,
				expr: &seqExpr{
					pos: position{line: 28, col: 8, offset: 622},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 28, col: 8, offset: 622},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 28, col: 13, offset: 627},
								name: "IdentifierName",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 28, col: 28, offset: 642},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 28, col: 31, offset: 645},
							label: "display",
							expr: &zeroOrOneExpr{
								pos: position{line: 28, col: 39, offset: 653},
								expr: &seqExpr{
									pos: position{line: 28, col: 41, offset: 655},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 28, col: 41, offset: 655},
											name: "StringLiteral",
										},
										&ruleRefExpr{
											pos:  position{line: 28, col: 55, offset: 669},
											name: "__",
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 28, col: 61, offset: 675},
							name: "RuleDefOp",
						},
						&ruleRefExpr{
							pos:  position{line: 28, col: 71, offset: 685},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 28, col: 74, offset: 688},
							label: "expr",
							expr: &ruleRefExpr{
								pos:  position{line: 28, col: 79, offset: 693},
								name: "Expression",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 28, col: 90, offset: 704},
							name: "EOS",
						},
					},
				},
			},
		},
		{
			name: "Expression",
			pos:  position{line: 41, col: 1, offset: 1001},
			expr: &ruleRefExpr{
				pos:  position{line: 41, col: 14, offset: 1016},
				name: "ChoiceExpr",
			},
		},
		{
			name: "ChoiceExpr",
			pos:  position{line: 43, col: 1, offset: 1030},
			expr: &actionExpr{
				pos: position{line: 43, col: 14, offset: 1045},
				run: (*parser).callonChoiceExpr1,
				expr: &seqExpr{
					pos: position{line: 43, col: 14, offset: 1045},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 43, col: 14, offset: 1045},
							label: "first",
							expr: &ruleRefExpr{
								pos:  position{line: 43, col: 20, offset: 1051},
								name: "ActionExpr",
							},
						},
						&labeledExpr{
							pos:   position{line: 43, col: 31, offset: 1062},
							label: "rest",
							expr: &zeroOrMoreExpr{
								pos: position{line: 43, col: 36, offset: 1067},
								expr: &seqExpr{
									pos: position{line: 43, col: 38, offset: 1069},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 43, col: 38, offset: 1069},
											name: "__",
										},
										&litMatcher{
											pos:        position{line: 43, col: 41, offset: 1072},
											val:        "/",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 43, col: 45, offset: 1076},
											name: "__",
										},
										&ruleRefExpr{
											pos:  position{line: 43, col: 48, offset: 1079},
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
			pos:  position{line: 58, col: 1, offset: 1499},
			expr: &actionExpr{
				pos: position{line: 58, col: 14, offset: 1514},
				run: (*parser).callonActionExpr1,
				expr: &seqExpr{
					pos: position{line: 58, col: 14, offset: 1514},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 58, col: 14, offset: 1514},
							label: "expr",
							expr: &ruleRefExpr{
								pos:  position{line: 58, col: 19, offset: 1519},
								name: "SeqExpr",
							},
						},
						&labeledExpr{
							pos:   position{line: 58, col: 27, offset: 1527},
							label: "code",
							expr: &zeroOrOneExpr{
								pos: position{line: 58, col: 32, offset: 1532},
								expr: &seqExpr{
									pos: position{line: 58, col: 34, offset: 1534},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 58, col: 34, offset: 1534},
											name: "__",
										},
										&ruleRefExpr{
											pos:  position{line: 58, col: 37, offset: 1537},
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
			pos:  position{line: 72, col: 1, offset: 1817},
			expr: &actionExpr{
				pos: position{line: 72, col: 11, offset: 1829},
				run: (*parser).callonSeqExpr1,
				expr: &seqExpr{
					pos: position{line: 72, col: 11, offset: 1829},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 72, col: 11, offset: 1829},
							label: "first",
							expr: &ruleRefExpr{
								pos:  position{line: 72, col: 17, offset: 1835},
								name: "LabeledExpr",
							},
						},
						&labeledExpr{
							pos:   position{line: 72, col: 29, offset: 1847},
							label: "rest",
							expr: &zeroOrMoreExpr{
								pos: position{line: 72, col: 34, offset: 1852},
								expr: &seqExpr{
									pos: position{line: 72, col: 36, offset: 1854},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 72, col: 36, offset: 1854},
											name: "__",
										},
										&ruleRefExpr{
											pos:  position{line: 72, col: 39, offset: 1857},
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
			pos:  position{line: 85, col: 1, offset: 2221},
			expr: &choiceExpr{
				pos: position{line: 85, col: 15, offset: 2237},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 85, col: 15, offset: 2237},
						run: (*parser).callonLabeledExpr2,
						expr: &seqExpr{
							pos: position{line: 85, col: 15, offset: 2237},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 85, col: 15, offset: 2237},
									label: "label",
									expr: &ruleRefExpr{
										pos:  position{line: 85, col: 21, offset: 2243},
										name: "Identifier",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 85, col: 32, offset: 2254},
									name: "__",
								},
								&litMatcher{
									pos:        position{line: 85, col: 35, offset: 2257},
									val:        ":",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 85, col: 39, offset: 2261},
									name: "__",
								},
								&labeledExpr{
									pos:   position{line: 85, col: 42, offset: 2264},
									label: "expr",
									expr: &ruleRefExpr{
										pos:  position{line: 85, col: 47, offset: 2269},
										name: "PrefixedExpr",
									},
								},
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 91, col: 5, offset: 2448},
						name: "PrefixedExpr",
					},
				},
			},
		},
		{
			name: "PrefixedExpr",
			pos:  position{line: 93, col: 1, offset: 2464},
			expr: &choiceExpr{
				pos: position{line: 93, col: 16, offset: 2481},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 93, col: 16, offset: 2481},
						run: (*parser).callonPrefixedExpr2,
						expr: &seqExpr{
							pos: position{line: 93, col: 16, offset: 2481},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 93, col: 16, offset: 2481},
									label: "op",
									expr: &ruleRefExpr{
										pos:  position{line: 93, col: 19, offset: 2484},
										name: "PrefixedOp",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 93, col: 30, offset: 2495},
									name: "__",
								},
								&labeledExpr{
									pos:   position{line: 93, col: 33, offset: 2498},
									label: "expr",
									expr: &ruleRefExpr{
										pos:  position{line: 93, col: 38, offset: 2503},
										name: "SuffixedExpr",
									},
								},
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 104, col: 5, offset: 2796},
						name: "SuffixedExpr",
					},
				},
			},
		},
		{
			name: "PrefixedOp",
			pos:  position{line: 106, col: 1, offset: 2812},
			expr: &actionExpr{
				pos: position{line: 106, col: 14, offset: 2827},
				run: (*parser).callonPrefixedOp1,
				expr: &choiceExpr{
					pos: position{line: 106, col: 16, offset: 2829},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 106, col: 16, offset: 2829},
							val:        "&",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 106, col: 22, offset: 2835},
							val:        "!",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "SuffixedExpr",
			pos:  position{line: 110, col: 1, offset: 2881},
			expr: &choiceExpr{
				pos: position{line: 110, col: 16, offset: 2898},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 110, col: 16, offset: 2898},
						run: (*parser).callonSuffixedExpr2,
						expr: &seqExpr{
							pos: position{line: 110, col: 16, offset: 2898},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 110, col: 16, offset: 2898},
									label: "expr",
									expr: &ruleRefExpr{
										pos:  position{line: 110, col: 21, offset: 2903},
										name: "PrimaryExpr",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 110, col: 33, offset: 2915},
									name: "__",
								},
								&labeledExpr{
									pos:   position{line: 110, col: 36, offset: 2918},
									label: "op",
									expr: &ruleRefExpr{
										pos:  position{line: 110, col: 39, offset: 2921},
										name: "SuffixedOp",
									},
								},
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 129, col: 5, offset: 3470},
						name: "PrimaryExpr",
					},
				},
			},
		},
		{
			name: "SuffixedOp",
			pos:  position{line: 131, col: 1, offset: 3486},
			expr: &actionExpr{
				pos: position{line: 131, col: 14, offset: 3501},
				run: (*parser).callonSuffixedOp1,
				expr: &choiceExpr{
					pos: position{line: 131, col: 16, offset: 3503},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 131, col: 16, offset: 3503},
							val:        "?",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 131, col: 22, offset: 3509},
							val:        "*",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 131, col: 28, offset: 3515},
							val:        "+",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "PrimaryExpr",
			pos:  position{line: 135, col: 1, offset: 3561},
			expr: &choiceExpr{
				pos: position{line: 135, col: 15, offset: 3577},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 135, col: 15, offset: 3577},
						name: "LitMatcher",
					},
					&ruleRefExpr{
						pos:  position{line: 135, col: 28, offset: 3590},
						name: "CharClassMatcher",
					},
					&ruleRefExpr{
						pos:  position{line: 135, col: 47, offset: 3609},
						name: "AnyMatcher",
					},
					&ruleRefExpr{
						pos:  position{line: 135, col: 60, offset: 3622},
						name: "RuleRefExpr",
					},
					&ruleRefExpr{
						pos:  position{line: 135, col: 74, offset: 3636},
						name: "SemanticPredExpr",
					},
					&actionExpr{
						pos: position{line: 135, col: 93, offset: 3655},
						run: (*parser).callonPrimaryExpr7,
						expr: &seqExpr{
							pos: position{line: 135, col: 93, offset: 3655},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 135, col: 93, offset: 3655},
									val:        "(",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 135, col: 97, offset: 3659},
									name: "__",
								},
								&labeledExpr{
									pos:   position{line: 135, col: 100, offset: 3662},
									label: "expr",
									expr: &ruleRefExpr{
										pos:  position{line: 135, col: 105, offset: 3667},
										name: "Expression",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 135, col: 116, offset: 3678},
									name: "__",
								},
								&litMatcher{
									pos:        position{line: 135, col: 119, offset: 3681},
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
			pos:  position{line: 138, col: 1, offset: 3713},
			expr: &actionExpr{
				pos: position{line: 138, col: 15, offset: 3729},
				run: (*parser).callonRuleRefExpr1,
				expr: &seqExpr{
					pos: position{line: 138, col: 15, offset: 3729},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 138, col: 15, offset: 3729},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 138, col: 20, offset: 3734},
								name: "IdentifierName",
							},
						},
						&notExpr{
							pos: position{line: 138, col: 35, offset: 3749},
							expr: &seqExpr{
								pos: position{line: 138, col: 38, offset: 3752},
								exprs: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 138, col: 38, offset: 3752},
										name: "__",
									},
									&zeroOrOneExpr{
										pos: position{line: 138, col: 41, offset: 3755},
										expr: &seqExpr{
											pos: position{line: 138, col: 43, offset: 3757},
											exprs: []interface{}{
												&ruleRefExpr{
													pos:  position{line: 138, col: 43, offset: 3757},
													name: "StringLiteral",
												},
												&ruleRefExpr{
													pos:  position{line: 138, col: 57, offset: 3771},
													name: "__",
												},
											},
										},
									},
									&ruleRefExpr{
										pos:  position{line: 138, col: 63, offset: 3777},
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
			pos:  position{line: 143, col: 1, offset: 3898},
			expr: &actionExpr{
				pos: position{line: 143, col: 20, offset: 3919},
				run: (*parser).callonSemanticPredExpr1,
				expr: &seqExpr{
					pos: position{line: 143, col: 20, offset: 3919},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 143, col: 20, offset: 3919},
							label: "op",
							expr: &ruleRefExpr{
								pos:  position{line: 143, col: 23, offset: 3922},
								name: "SemanticPredOp",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 143, col: 38, offset: 3937},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 143, col: 41, offset: 3940},
							label: "code",
							expr: &ruleRefExpr{
								pos:  position{line: 143, col: 46, offset: 3945},
								name: "CodeBlock",
							},
						},
					},
				},
			},
		},
		{
			name: "SemanticPredOp",
			pos:  position{line: 154, col: 1, offset: 4233},
			expr: &actionExpr{
				pos: position{line: 154, col: 18, offset: 4252},
				run: (*parser).callonSemanticPredOp1,
				expr: &choiceExpr{
					pos: position{line: 154, col: 20, offset: 4254},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 154, col: 20, offset: 4254},
							val:        "&",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 154, col: 26, offset: 4260},
							val:        "!",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "RuleDefOp",
			pos:  position{line: 158, col: 1, offset: 4306},
			expr: &choiceExpr{
				pos: position{line: 158, col: 13, offset: 4320},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 158, col: 13, offset: 4320},
						val:        "=",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 158, col: 19, offset: 4326},
						val:        "<-",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 158, col: 26, offset: 4333},
						val:        "←",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 158, col: 37, offset: 4344},
						val:        "⟵",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "SourceChar",
			pos:  position{line: 160, col: 1, offset: 4356},
			expr: &anyMatcher{
				line: 160, col: 14, offset: 4371,
			},
		},
		{
			name: "Comment",
			pos:  position{line: 161, col: 1, offset: 4374},
			expr: &choiceExpr{
				pos: position{line: 161, col: 11, offset: 4386},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 161, col: 11, offset: 4386},
						name: "MultiLineComment",
					},
					&ruleRefExpr{
						pos:  position{line: 161, col: 30, offset: 4405},
						name: "SingleLineComment",
					},
				},
			},
		},
		{
			name: "MultiLineComment",
			pos:  position{line: 162, col: 1, offset: 4424},
			expr: &seqExpr{
				pos: position{line: 162, col: 20, offset: 4445},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 162, col: 20, offset: 4445},
						val:        "/*",
						ignoreCase: false,
					},
					&zeroOrMoreExpr{
						pos: position{line: 162, col: 25, offset: 4450},
						expr: &seqExpr{
							pos: position{line: 162, col: 27, offset: 4452},
							exprs: []interface{}{
								&notExpr{
									pos: position{line: 162, col: 27, offset: 4452},
									expr: &litMatcher{
										pos:        position{line: 162, col: 28, offset: 4453},
										val:        "*/",
										ignoreCase: false,
									},
								},
								&ruleRefExpr{
									pos:  position{line: 162, col: 33, offset: 4458},
									name: "SourceChar",
								},
							},
						},
					},
					&litMatcher{
						pos:        position{line: 162, col: 47, offset: 4472},
						val:        "*/",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "MultiLineCommentNoLineTerminator",
			pos:  position{line: 163, col: 1, offset: 4478},
			expr: &seqExpr{
				pos: position{line: 163, col: 36, offset: 4515},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 163, col: 36, offset: 4515},
						val:        "/*",
						ignoreCase: false,
					},
					&zeroOrMoreExpr{
						pos: position{line: 163, col: 41, offset: 4520},
						expr: &seqExpr{
							pos: position{line: 163, col: 43, offset: 4522},
							exprs: []interface{}{
								&notExpr{
									pos: position{line: 163, col: 43, offset: 4522},
									expr: &choiceExpr{
										pos: position{line: 163, col: 46, offset: 4525},
										alternatives: []interface{}{
											&litMatcher{
												pos:        position{line: 163, col: 46, offset: 4525},
												val:        "*/",
												ignoreCase: false,
											},
											&ruleRefExpr{
												pos:  position{line: 163, col: 53, offset: 4532},
												name: "EOL",
											},
										},
									},
								},
								&ruleRefExpr{
									pos:  position{line: 163, col: 59, offset: 4538},
									name: "SourceChar",
								},
							},
						},
					},
					&litMatcher{
						pos:        position{line: 163, col: 73, offset: 4552},
						val:        "*/",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "SingleLineComment",
			pos:  position{line: 164, col: 1, offset: 4558},
			expr: &seqExpr{
				pos: position{line: 164, col: 21, offset: 4580},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 164, col: 21, offset: 4580},
						val:        "//",
						ignoreCase: false,
					},
					&zeroOrMoreExpr{
						pos: position{line: 164, col: 26, offset: 4585},
						expr: &seqExpr{
							pos: position{line: 164, col: 28, offset: 4587},
							exprs: []interface{}{
								&notExpr{
									pos: position{line: 164, col: 28, offset: 4587},
									expr: &ruleRefExpr{
										pos:  position{line: 164, col: 29, offset: 4588},
										name: "EOL",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 164, col: 33, offset: 4592},
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
			pos:  position{line: 166, col: 1, offset: 4609},
			expr: &actionExpr{
				pos: position{line: 166, col: 14, offset: 4624},
				run: (*parser).callonIdentifier1,
				expr: &labeledExpr{
					pos:   position{line: 166, col: 14, offset: 4624},
					label: "ident",
					expr: &ruleRefExpr{
						pos:  position{line: 166, col: 20, offset: 4630},
						name: "IdentifierName",
					},
				},
			},
		},
		{
			name: "IdentifierName",
			pos:  position{line: 174, col: 1, offset: 4857},
			expr: &actionExpr{
				pos: position{line: 174, col: 18, offset: 4876},
				run: (*parser).callonIdentifierName1,
				expr: &seqExpr{
					pos: position{line: 174, col: 18, offset: 4876},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 174, col: 18, offset: 4876},
							name: "IdentifierStart",
						},
						&zeroOrMoreExpr{
							pos: position{line: 174, col: 34, offset: 4892},
							expr: &ruleRefExpr{
								pos:  position{line: 174, col: 34, offset: 4892},
								name: "IdentifierPart",
							},
						},
					},
				},
			},
		},
		{
			name: "IdentifierStart",
			pos:  position{line: 177, col: 1, offset: 4977},
			expr: &charClassMatcher{
				pos:        position{line: 177, col: 19, offset: 4997},
				val:        "[\\pL_]",
				chars:      []rune{'_'},
				classes:    []*unicode.RangeTable{rangeTable("L")},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "IdentifierPart",
			pos:  position{line: 178, col: 1, offset: 5005},
			expr: &choiceExpr{
				pos: position{line: 178, col: 18, offset: 5024},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 178, col: 18, offset: 5024},
						name: "IdentifierStart",
					},
					&charClassMatcher{
						pos:        position{line: 178, col: 36, offset: 5042},
						val:        "[\\p{Nd}]",
						classes:    []*unicode.RangeTable{rangeTable("Nd")},
						ignoreCase: false,
						inverted:   false,
					},
				},
			},
		},
		{
			name: "LitMatcher",
			pos:  position{line: 180, col: 1, offset: 5054},
			expr: &actionExpr{
				pos: position{line: 180, col: 14, offset: 5069},
				run: (*parser).callonLitMatcher1,
				expr: &seqExpr{
					pos: position{line: 180, col: 14, offset: 5069},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 180, col: 14, offset: 5069},
							label: "lit",
							expr: &ruleRefExpr{
								pos:  position{line: 180, col: 18, offset: 5073},
								name: "StringLiteral",
							},
						},
						&labeledExpr{
							pos:   position{line: 180, col: 32, offset: 5087},
							label: "ignore",
							expr: &zeroOrOneExpr{
								pos: position{line: 180, col: 39, offset: 5094},
								expr: &litMatcher{
									pos:        position{line: 180, col: 39, offset: 5094},
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
			pos:  position{line: 193, col: 1, offset: 5506},
			expr: &choiceExpr{
				pos: position{line: 193, col: 17, offset: 5524},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 193, col: 17, offset: 5524},
						run: (*parser).callonStringLiteral2,
						expr: &choiceExpr{
							pos: position{line: 193, col: 19, offset: 5526},
							alternatives: []interface{}{
								&seqExpr{
									pos: position{line: 193, col: 19, offset: 5526},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 193, col: 19, offset: 5526},
											val:        "\"",
											ignoreCase: false,
										},
										&zeroOrMoreExpr{
											pos: position{line: 193, col: 23, offset: 5530},
											expr: &ruleRefExpr{
												pos:  position{line: 193, col: 23, offset: 5530},
												name: "DoubleStringChar",
											},
										},
										&litMatcher{
											pos:        position{line: 193, col: 41, offset: 5548},
											val:        "\"",
											ignoreCase: false,
										},
									},
								},
								&seqExpr{
									pos: position{line: 193, col: 47, offset: 5554},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 193, col: 47, offset: 5554},
											val:        "'",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 193, col: 51, offset: 5558},
											name: "SingleStringChar",
										},
										&litMatcher{
											pos:        position{line: 193, col: 68, offset: 5575},
											val:        "'",
											ignoreCase: false,
										},
									},
								},
								&seqExpr{
									pos: position{line: 193, col: 74, offset: 5581},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 193, col: 74, offset: 5581},
											val:        "`",
											ignoreCase: false,
										},
										&zeroOrMoreExpr{
											pos: position{line: 193, col: 78, offset: 5585},
											expr: &ruleRefExpr{
												pos:  position{line: 193, col: 78, offset: 5585},
												name: "RawStringChar",
											},
										},
										&litMatcher{
											pos:        position{line: 193, col: 93, offset: 5600},
											val:        "`",
											ignoreCase: false,
										},
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 195, col: 5, offset: 5675},
						run: (*parser).callonStringLiteral18,
						expr: &choiceExpr{
							pos: position{line: 195, col: 7, offset: 5677},
							alternatives: []interface{}{
								&seqExpr{
									pos: position{line: 195, col: 9, offset: 5679},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 195, col: 9, offset: 5679},
											val:        "\"",
											ignoreCase: false,
										},
										&zeroOrMoreExpr{
											pos: position{line: 195, col: 13, offset: 5683},
											expr: &ruleRefExpr{
												pos:  position{line: 195, col: 13, offset: 5683},
												name: "DoubleStringChar",
											},
										},
										&choiceExpr{
											pos: position{line: 195, col: 33, offset: 5703},
											alternatives: []interface{}{
												&ruleRefExpr{
													pos:  position{line: 195, col: 33, offset: 5703},
													name: "EOL",
												},
												&ruleRefExpr{
													pos:  position{line: 195, col: 39, offset: 5709},
													name: "EOF",
												},
											},
										},
									},
								},
								&seqExpr{
									pos: position{line: 195, col: 51, offset: 5721},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 195, col: 51, offset: 5721},
											val:        "'",
											ignoreCase: false,
										},
										&zeroOrOneExpr{
											pos: position{line: 195, col: 55, offset: 5725},
											expr: &ruleRefExpr{
												pos:  position{line: 195, col: 55, offset: 5725},
												name: "SingleStringChar",
											},
										},
										&choiceExpr{
											pos: position{line: 195, col: 75, offset: 5745},
											alternatives: []interface{}{
												&ruleRefExpr{
													pos:  position{line: 195, col: 75, offset: 5745},
													name: "EOL",
												},
												&ruleRefExpr{
													pos:  position{line: 195, col: 81, offset: 5751},
													name: "EOF",
												},
											},
										},
									},
								},
								&seqExpr{
									pos: position{line: 195, col: 91, offset: 5761},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 195, col: 91, offset: 5761},
											val:        "`",
											ignoreCase: false,
										},
										&zeroOrMoreExpr{
											pos: position{line: 195, col: 95, offset: 5765},
											expr: &ruleRefExpr{
												pos:  position{line: 195, col: 95, offset: 5765},
												name: "RawStringChar",
											},
										},
										&ruleRefExpr{
											pos:  position{line: 195, col: 110, offset: 5780},
											name: "EOF",
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
			name: "DoubleStringChar",
			pos:  position{line: 199, col: 1, offset: 5886},
			expr: &choiceExpr{
				pos: position{line: 199, col: 20, offset: 5907},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 199, col: 20, offset: 5907},
						exprs: []interface{}{
							&notExpr{
								pos: position{line: 199, col: 20, offset: 5907},
								expr: &choiceExpr{
									pos: position{line: 199, col: 23, offset: 5910},
									alternatives: []interface{}{
										&litMatcher{
											pos:        position{line: 199, col: 23, offset: 5910},
											val:        "\"",
											ignoreCase: false,
										},
										&litMatcher{
											pos:        position{line: 199, col: 29, offset: 5916},
											val:        "\\",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 199, col: 36, offset: 5923},
											name: "EOL",
										},
									},
								},
							},
							&ruleRefExpr{
								pos:  position{line: 199, col: 42, offset: 5929},
								name: "SourceChar",
							},
						},
					},
					&seqExpr{
						pos: position{line: 199, col: 55, offset: 5942},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 199, col: 55, offset: 5942},
								val:        "\\",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 199, col: 60, offset: 5947},
								name: "DoubleStringEscape",
							},
						},
					},
				},
			},
		},
		{
			name: "SingleStringChar",
			pos:  position{line: 200, col: 1, offset: 5967},
			expr: &choiceExpr{
				pos: position{line: 200, col: 20, offset: 5988},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 200, col: 20, offset: 5988},
						exprs: []interface{}{
							&notExpr{
								pos: position{line: 200, col: 20, offset: 5988},
								expr: &choiceExpr{
									pos: position{line: 200, col: 23, offset: 5991},
									alternatives: []interface{}{
										&litMatcher{
											pos:        position{line: 200, col: 23, offset: 5991},
											val:        "'",
											ignoreCase: false,
										},
										&litMatcher{
											pos:        position{line: 200, col: 29, offset: 5997},
											val:        "\\",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 200, col: 36, offset: 6004},
											name: "EOL",
										},
									},
								},
							},
							&ruleRefExpr{
								pos:  position{line: 200, col: 42, offset: 6010},
								name: "SourceChar",
							},
						},
					},
					&seqExpr{
						pos: position{line: 200, col: 55, offset: 6023},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 200, col: 55, offset: 6023},
								val:        "\\",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 200, col: 60, offset: 6028},
								name: "SingleStringEscape",
							},
						},
					},
				},
			},
		},
		{
			name: "RawStringChar",
			pos:  position{line: 201, col: 1, offset: 6048},
			expr: &seqExpr{
				pos: position{line: 201, col: 17, offset: 6066},
				exprs: []interface{}{
					&notExpr{
						pos: position{line: 201, col: 17, offset: 6066},
						expr: &litMatcher{
							pos:        position{line: 201, col: 18, offset: 6067},
							val:        "`",
							ignoreCase: false,
						},
					},
					&ruleRefExpr{
						pos:  position{line: 201, col: 22, offset: 6071},
						name: "SourceChar",
					},
				},
			},
		},
		{
			name: "DoubleStringEscape",
			pos:  position{line: 203, col: 1, offset: 6085},
			expr: &choiceExpr{
				pos: position{line: 203, col: 22, offset: 6108},
				alternatives: []interface{}{
					&choiceExpr{
						pos: position{line: 203, col: 24, offset: 6110},
						alternatives: []interface{}{
							&litMatcher{
								pos:        position{line: 203, col: 24, offset: 6110},
								val:        "\"",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 203, col: 30, offset: 6116},
								name: "CommonEscapeSequence",
							},
						},
					},
					&actionExpr{
						pos: position{line: 204, col: 7, offset: 6146},
						run: (*parser).callonDoubleStringEscape5,
						expr: &choiceExpr{
							pos: position{line: 204, col: 9, offset: 6148},
							alternatives: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 204, col: 9, offset: 6148},
									name: "SourceChar",
								},
								&ruleRefExpr{
									pos:  position{line: 204, col: 22, offset: 6161},
									name: "EOL",
								},
								&ruleRefExpr{
									pos:  position{line: 204, col: 28, offset: 6167},
									name: "EOF",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "SingleStringEscape",
			pos:  position{line: 207, col: 1, offset: 6235},
			expr: &choiceExpr{
				pos: position{line: 207, col: 22, offset: 6258},
				alternatives: []interface{}{
					&choiceExpr{
						pos: position{line: 207, col: 24, offset: 6260},
						alternatives: []interface{}{
							&litMatcher{
								pos:        position{line: 207, col: 24, offset: 6260},
								val:        "'",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 207, col: 30, offset: 6266},
								name: "CommonEscapeSequence",
							},
						},
					},
					&actionExpr{
						pos: position{line: 208, col: 7, offset: 6296},
						run: (*parser).callonSingleStringEscape5,
						expr: &choiceExpr{
							pos: position{line: 208, col: 9, offset: 6298},
							alternatives: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 208, col: 9, offset: 6298},
									name: "SourceChar",
								},
								&ruleRefExpr{
									pos:  position{line: 208, col: 22, offset: 6311},
									name: "EOL",
								},
								&ruleRefExpr{
									pos:  position{line: 208, col: 28, offset: 6317},
									name: "EOF",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "CommonEscapeSequence",
			pos:  position{line: 212, col: 1, offset: 6387},
			expr: &choiceExpr{
				pos: position{line: 212, col: 24, offset: 6412},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 212, col: 24, offset: 6412},
						name: "SingleCharEscape",
					},
					&ruleRefExpr{
						pos:  position{line: 212, col: 43, offset: 6431},
						name: "OctalEscape",
					},
					&ruleRefExpr{
						pos:  position{line: 212, col: 57, offset: 6445},
						name: "HexEscape",
					},
					&ruleRefExpr{
						pos:  position{line: 212, col: 69, offset: 6457},
						name: "LongUnicodeEscape",
					},
					&ruleRefExpr{
						pos:  position{line: 212, col: 89, offset: 6477},
						name: "ShortUnicodeEscape",
					},
				},
			},
		},
		{
			name: "SingleCharEscape",
			pos:  position{line: 213, col: 1, offset: 6497},
			expr: &choiceExpr{
				pos: position{line: 213, col: 20, offset: 6518},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 213, col: 20, offset: 6518},
						val:        "a",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 213, col: 26, offset: 6524},
						val:        "b",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 213, col: 32, offset: 6530},
						val:        "n",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 213, col: 38, offset: 6536},
						val:        "f",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 213, col: 44, offset: 6542},
						val:        "r",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 213, col: 50, offset: 6548},
						val:        "t",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 213, col: 56, offset: 6554},
						val:        "v",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 213, col: 62, offset: 6560},
						val:        "\\",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "OctalEscape",
			pos:  position{line: 214, col: 1, offset: 6566},
			expr: &choiceExpr{
				pos: position{line: 214, col: 15, offset: 6582},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 214, col: 15, offset: 6582},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 214, col: 15, offset: 6582},
								name: "OctalDigit",
							},
							&ruleRefExpr{
								pos:  position{line: 214, col: 26, offset: 6593},
								name: "OctalDigit",
							},
							&ruleRefExpr{
								pos:  position{line: 214, col: 37, offset: 6604},
								name: "OctalDigit",
							},
						},
					},
					&actionExpr{
						pos: position{line: 215, col: 7, offset: 6622},
						run: (*parser).callonOctalEscape6,
						expr: &seqExpr{
							pos: position{line: 215, col: 7, offset: 6622},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 215, col: 7, offset: 6622},
									name: "OctalDigit",
								},
								&choiceExpr{
									pos: position{line: 215, col: 20, offset: 6635},
									alternatives: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 215, col: 20, offset: 6635},
											name: "SourceChar",
										},
										&ruleRefExpr{
											pos:  position{line: 215, col: 33, offset: 6648},
											name: "EOL",
										},
										&ruleRefExpr{
											pos:  position{line: 215, col: 39, offset: 6654},
											name: "EOF",
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
			name: "HexEscape",
			pos:  position{line: 218, col: 1, offset: 6718},
			expr: &choiceExpr{
				pos: position{line: 218, col: 13, offset: 6732},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 218, col: 13, offset: 6732},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 218, col: 13, offset: 6732},
								val:        "x",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 218, col: 17, offset: 6736},
								name: "HexDigit",
							},
							&ruleRefExpr{
								pos:  position{line: 218, col: 26, offset: 6745},
								name: "HexDigit",
							},
						},
					},
					&actionExpr{
						pos: position{line: 219, col: 7, offset: 6761},
						run: (*parser).callonHexEscape6,
						expr: &seqExpr{
							pos: position{line: 219, col: 7, offset: 6761},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 219, col: 7, offset: 6761},
									val:        "x",
									ignoreCase: false,
								},
								&choiceExpr{
									pos: position{line: 219, col: 13, offset: 6767},
									alternatives: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 219, col: 13, offset: 6767},
											name: "SourceChar",
										},
										&ruleRefExpr{
											pos:  position{line: 219, col: 26, offset: 6780},
											name: "EOL",
										},
										&ruleRefExpr{
											pos:  position{line: 219, col: 32, offset: 6786},
											name: "EOF",
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
			name: "LongUnicodeEscape",
			pos:  position{line: 222, col: 1, offset: 6856},
			expr: &choiceExpr{
				pos: position{line: 223, col: 5, offset: 6884},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 223, col: 5, offset: 6884},
						run: (*parser).callonLongUnicodeEscape2,
						expr: &seqExpr{
							pos: position{line: 223, col: 5, offset: 6884},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 223, col: 5, offset: 6884},
									val:        "U",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 223, col: 9, offset: 6888},
									name: "HexDigit",
								},
								&ruleRefExpr{
									pos:  position{line: 223, col: 18, offset: 6897},
									name: "HexDigit",
								},
								&ruleRefExpr{
									pos:  position{line: 223, col: 27, offset: 6906},
									name: "HexDigit",
								},
								&ruleRefExpr{
									pos:  position{line: 223, col: 36, offset: 6915},
									name: "HexDigit",
								},
								&ruleRefExpr{
									pos:  position{line: 223, col: 45, offset: 6924},
									name: "HexDigit",
								},
								&ruleRefExpr{
									pos:  position{line: 223, col: 54, offset: 6933},
									name: "HexDigit",
								},
								&ruleRefExpr{
									pos:  position{line: 223, col: 63, offset: 6942},
									name: "HexDigit",
								},
								&ruleRefExpr{
									pos:  position{line: 223, col: 72, offset: 6951},
									name: "HexDigit",
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 226, col: 7, offset: 7056},
						run: (*parser).callonLongUnicodeEscape13,
						expr: &seqExpr{
							pos: position{line: 226, col: 7, offset: 7056},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 226, col: 7, offset: 7056},
									val:        "U",
									ignoreCase: false,
								},
								&choiceExpr{
									pos: position{line: 226, col: 13, offset: 7062},
									alternatives: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 226, col: 13, offset: 7062},
											name: "SourceChar",
										},
										&ruleRefExpr{
											pos:  position{line: 226, col: 26, offset: 7075},
											name: "EOL",
										},
										&ruleRefExpr{
											pos:  position{line: 226, col: 32, offset: 7081},
											name: "EOF",
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
			name: "ShortUnicodeEscape",
			pos:  position{line: 229, col: 1, offset: 7147},
			expr: &choiceExpr{
				pos: position{line: 230, col: 5, offset: 7176},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 230, col: 5, offset: 7176},
						run: (*parser).callonShortUnicodeEscape2,
						expr: &seqExpr{
							pos: position{line: 230, col: 5, offset: 7176},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 230, col: 5, offset: 7176},
									val:        "u",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 230, col: 9, offset: 7180},
									name: "HexDigit",
								},
								&ruleRefExpr{
									pos:  position{line: 230, col: 18, offset: 7189},
									name: "HexDigit",
								},
								&ruleRefExpr{
									pos:  position{line: 230, col: 27, offset: 7198},
									name: "HexDigit",
								},
								&ruleRefExpr{
									pos:  position{line: 230, col: 36, offset: 7207},
									name: "HexDigit",
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 233, col: 7, offset: 7312},
						run: (*parser).callonShortUnicodeEscape9,
						expr: &seqExpr{
							pos: position{line: 233, col: 7, offset: 7312},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 233, col: 7, offset: 7312},
									val:        "u",
									ignoreCase: false,
								},
								&choiceExpr{
									pos: position{line: 233, col: 13, offset: 7318},
									alternatives: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 233, col: 13, offset: 7318},
											name: "SourceChar",
										},
										&ruleRefExpr{
											pos:  position{line: 233, col: 26, offset: 7331},
											name: "EOL",
										},
										&ruleRefExpr{
											pos:  position{line: 233, col: 32, offset: 7337},
											name: "EOF",
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
			name: "OctalDigit",
			pos:  position{line: 237, col: 1, offset: 7405},
			expr: &charClassMatcher{
				pos:        position{line: 237, col: 14, offset: 7420},
				val:        "[0-7]",
				ranges:     []rune{'0', '7'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "DecimalDigit",
			pos:  position{line: 238, col: 1, offset: 7427},
			expr: &charClassMatcher{
				pos:        position{line: 238, col: 16, offset: 7444},
				val:        "[0-9]",
				ranges:     []rune{'0', '9'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "HexDigit",
			pos:  position{line: 239, col: 1, offset: 7451},
			expr: &charClassMatcher{
				pos:        position{line: 239, col: 12, offset: 7464},
				val:        "[0-9a-f]i",
				ranges:     []rune{'0', '9', 'a', 'f'},
				ignoreCase: true,
				inverted:   false,
			},
		},
		{
			name: "CharClassMatcher",
			pos:  position{line: 241, col: 1, offset: 7477},
			expr: &choiceExpr{
				pos: position{line: 241, col: 20, offset: 7498},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 241, col: 20, offset: 7498},
						run: (*parser).callonCharClassMatcher2,
						expr: &seqExpr{
							pos: position{line: 241, col: 20, offset: 7498},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 241, col: 20, offset: 7498},
									val:        "[",
									ignoreCase: false,
								},
								&zeroOrMoreExpr{
									pos: position{line: 241, col: 24, offset: 7502},
									expr: &choiceExpr{
										pos: position{line: 241, col: 26, offset: 7504},
										alternatives: []interface{}{
											&ruleRefExpr{
												pos:  position{line: 241, col: 26, offset: 7504},
												name: "ClassCharRange",
											},
											&ruleRefExpr{
												pos:  position{line: 241, col: 43, offset: 7521},
												name: "ClassChar",
											},
											&seqExpr{
												pos: position{line: 241, col: 55, offset: 7533},
												exprs: []interface{}{
													&litMatcher{
														pos:        position{line: 241, col: 55, offset: 7533},
														val:        "\\",
														ignoreCase: false,
													},
													&ruleRefExpr{
														pos:  position{line: 241, col: 60, offset: 7538},
														name: "UnicodeClassEscape",
													},
												},
											},
										},
									},
								},
								&litMatcher{
									pos:        position{line: 241, col: 82, offset: 7560},
									val:        "]",
									ignoreCase: false,
								},
								&zeroOrOneExpr{
									pos: position{line: 241, col: 86, offset: 7564},
									expr: &litMatcher{
										pos:        position{line: 241, col: 86, offset: 7564},
										val:        "i",
										ignoreCase: false,
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 245, col: 5, offset: 7675},
						run: (*parser).callonCharClassMatcher15,
						expr: &seqExpr{
							pos: position{line: 245, col: 5, offset: 7675},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 245, col: 5, offset: 7675},
									val:        "[",
									ignoreCase: false,
								},
								&zeroOrMoreExpr{
									pos: position{line: 245, col: 9, offset: 7679},
									expr: &seqExpr{
										pos: position{line: 245, col: 11, offset: 7681},
										exprs: []interface{}{
											&notExpr{
												pos: position{line: 245, col: 11, offset: 7681},
												expr: &ruleRefExpr{
													pos:  position{line: 245, col: 14, offset: 7684},
													name: "EOL",
												},
											},
											&ruleRefExpr{
												pos:  position{line: 245, col: 20, offset: 7690},
												name: "SourceChar",
											},
										},
									},
								},
								&choiceExpr{
									pos: position{line: 245, col: 36, offset: 7706},
									alternatives: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 245, col: 36, offset: 7706},
											name: "EOL",
										},
										&ruleRefExpr{
											pos:  position{line: 245, col: 42, offset: 7712},
											name: "EOF",
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
			name: "ClassCharRange",
			pos:  position{line: 249, col: 1, offset: 7826},
			expr: &seqExpr{
				pos: position{line: 249, col: 18, offset: 7845},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 249, col: 18, offset: 7845},
						name: "ClassChar",
					},
					&litMatcher{
						pos:        position{line: 249, col: 28, offset: 7855},
						val:        "-",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 249, col: 32, offset: 7859},
						name: "ClassChar",
					},
				},
			},
		},
		{
			name: "ClassChar",
			pos:  position{line: 250, col: 1, offset: 7870},
			expr: &choiceExpr{
				pos: position{line: 250, col: 13, offset: 7884},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 250, col: 13, offset: 7884},
						exprs: []interface{}{
							&notExpr{
								pos: position{line: 250, col: 13, offset: 7884},
								expr: &choiceExpr{
									pos: position{line: 250, col: 16, offset: 7887},
									alternatives: []interface{}{
										&litMatcher{
											pos:        position{line: 250, col: 16, offset: 7887},
											val:        "]",
											ignoreCase: false,
										},
										&litMatcher{
											pos:        position{line: 250, col: 22, offset: 7893},
											val:        "\\",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 250, col: 29, offset: 7900},
											name: "EOL",
										},
									},
								},
							},
							&ruleRefExpr{
								pos:  position{line: 250, col: 35, offset: 7906},
								name: "SourceChar",
							},
						},
					},
					&seqExpr{
						pos: position{line: 250, col: 48, offset: 7919},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 250, col: 48, offset: 7919},
								val:        "\\",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 250, col: 53, offset: 7924},
								name: "CharClassEscape",
							},
						},
					},
				},
			},
		},
		{
			name: "CharClassEscape",
			pos:  position{line: 251, col: 1, offset: 7941},
			expr: &choiceExpr{
				pos: position{line: 251, col: 19, offset: 7961},
				alternatives: []interface{}{
					&choiceExpr{
						pos: position{line: 251, col: 21, offset: 7963},
						alternatives: []interface{}{
							&litMatcher{
								pos:        position{line: 251, col: 21, offset: 7963},
								val:        "]",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 251, col: 27, offset: 7969},
								name: "CommonEscapeSequence",
							},
						},
					},
					&actionExpr{
						pos: position{line: 252, col: 7, offset: 7999},
						run: (*parser).callonCharClassEscape5,
						expr: &seqExpr{
							pos: position{line: 252, col: 7, offset: 7999},
							exprs: []interface{}{
								&notExpr{
									pos: position{line: 252, col: 7, offset: 7999},
									expr: &litMatcher{
										pos:        position{line: 252, col: 8, offset: 8000},
										val:        "p",
										ignoreCase: false,
									},
								},
								&choiceExpr{
									pos: position{line: 252, col: 14, offset: 8006},
									alternatives: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 252, col: 14, offset: 8006},
											name: "SourceChar",
										},
										&ruleRefExpr{
											pos:  position{line: 252, col: 27, offset: 8019},
											name: "EOL",
										},
										&ruleRefExpr{
											pos:  position{line: 252, col: 33, offset: 8025},
											name: "EOF",
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
			name: "UnicodeClassEscape",
			pos:  position{line: 256, col: 1, offset: 8095},
			expr: &seqExpr{
				pos: position{line: 256, col: 22, offset: 8118},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 256, col: 22, offset: 8118},
						val:        "p",
						ignoreCase: false,
					},
					&choiceExpr{
						pos: position{line: 257, col: 7, offset: 8132},
						alternatives: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 257, col: 7, offset: 8132},
								name: "SingleCharUnicodeClass",
							},
							&actionExpr{
								pos: position{line: 258, col: 7, offset: 8162},
								run: (*parser).callonUnicodeClassEscape5,
								expr: &seqExpr{
									pos: position{line: 258, col: 7, offset: 8162},
									exprs: []interface{}{
										&notExpr{
											pos: position{line: 258, col: 7, offset: 8162},
											expr: &litMatcher{
												pos:        position{line: 258, col: 8, offset: 8163},
												val:        "{",
												ignoreCase: false,
											},
										},
										&choiceExpr{
											pos: position{line: 258, col: 14, offset: 8169},
											alternatives: []interface{}{
												&ruleRefExpr{
													pos:  position{line: 258, col: 14, offset: 8169},
													name: "SourceChar",
												},
												&ruleRefExpr{
													pos:  position{line: 258, col: 27, offset: 8182},
													name: "EOL",
												},
												&ruleRefExpr{
													pos:  position{line: 258, col: 33, offset: 8188},
													name: "EOF",
												},
											},
										},
									},
								},
							},
							&actionExpr{
								pos: position{line: 259, col: 7, offset: 8260},
								run: (*parser).callonUnicodeClassEscape13,
								expr: &seqExpr{
									pos: position{line: 259, col: 7, offset: 8260},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 259, col: 7, offset: 8260},
											val:        "{",
											ignoreCase: false,
										},
										&labeledExpr{
											pos:   position{line: 259, col: 11, offset: 8264},
											label: "ident",
											expr: &ruleRefExpr{
												pos:  position{line: 259, col: 17, offset: 8270},
												name: "IdentifierName",
											},
										},
										&litMatcher{
											pos:        position{line: 259, col: 32, offset: 8285},
											val:        "}",
											ignoreCase: false,
										},
									},
								},
							},
							&actionExpr{
								pos: position{line: 265, col: 7, offset: 8468},
								run: (*parser).callonUnicodeClassEscape19,
								expr: &seqExpr{
									pos: position{line: 265, col: 7, offset: 8468},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 265, col: 7, offset: 8468},
											val:        "{",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 265, col: 11, offset: 8472},
											name: "IdentifierName",
										},
										&choiceExpr{
											pos: position{line: 265, col: 28, offset: 8489},
											alternatives: []interface{}{
												&litMatcher{
													pos:        position{line: 265, col: 28, offset: 8489},
													val:        "]",
													ignoreCase: false,
												},
												&ruleRefExpr{
													pos:  position{line: 265, col: 34, offset: 8495},
													name: "EOL",
												},
												&ruleRefExpr{
													pos:  position{line: 265, col: 40, offset: 8501},
													name: "EOF",
												},
											},
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
			name: "SingleCharUnicodeClass",
			pos:  position{line: 269, col: 1, offset: 8588},
			expr: &charClassMatcher{
				pos:        position{line: 269, col: 26, offset: 8615},
				val:        "[LMNCPZS]",
				chars:      []rune{'L', 'M', 'N', 'C', 'P', 'Z', 'S'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "AnyMatcher",
			pos:  position{line: 271, col: 1, offset: 8628},
			expr: &actionExpr{
				pos: position{line: 271, col: 14, offset: 8643},
				run: (*parser).callonAnyMatcher1,
				expr: &litMatcher{
					pos:        position{line: 271, col: 14, offset: 8643},
					val:        ".",
					ignoreCase: false,
				},
			},
		},
		{
			name: "CodeBlock",
			pos:  position{line: 276, col: 1, offset: 8723},
			expr: &choiceExpr{
				pos: position{line: 276, col: 13, offset: 8737},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 276, col: 13, offset: 8737},
						run: (*parser).callonCodeBlock2,
						expr: &seqExpr{
							pos: position{line: 276, col: 13, offset: 8737},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 276, col: 13, offset: 8737},
									val:        "{",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 276, col: 17, offset: 8741},
									name: "Code",
								},
								&litMatcher{
									pos:        position{line: 276, col: 22, offset: 8746},
									val:        "}",
									ignoreCase: false,
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 280, col: 5, offset: 8849},
						run: (*parser).callonCodeBlock7,
						expr: &seqExpr{
							pos: position{line: 280, col: 5, offset: 8849},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 280, col: 5, offset: 8849},
									val:        "{",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 280, col: 9, offset: 8853},
									name: "Code",
								},
								&ruleRefExpr{
									pos:  position{line: 280, col: 14, offset: 8858},
									name: "EOF",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Code",
			pos:  position{line: 284, col: 1, offset: 8927},
			expr: &zeroOrMoreExpr{
				pos: position{line: 284, col: 8, offset: 8936},
				expr: &choiceExpr{
					pos: position{line: 284, col: 10, offset: 8938},
					alternatives: []interface{}{
						&oneOrMoreExpr{
							pos: position{line: 284, col: 10, offset: 8938},
							expr: &seqExpr{
								pos: position{line: 284, col: 12, offset: 8940},
								exprs: []interface{}{
									&notExpr{
										pos: position{line: 284, col: 12, offset: 8940},
										expr: &charClassMatcher{
											pos:        position{line: 284, col: 13, offset: 8941},
											val:        "[{}]",
											chars:      []rune{'{', '}'},
											ignoreCase: false,
											inverted:   false,
										},
									},
									&ruleRefExpr{
										pos:  position{line: 284, col: 18, offset: 8946},
										name: "SourceChar",
									},
								},
							},
						},
						&seqExpr{
							pos: position{line: 284, col: 34, offset: 8962},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 284, col: 34, offset: 8962},
									val:        "{",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 284, col: 38, offset: 8966},
									name: "Code",
								},
								&litMatcher{
									pos:        position{line: 284, col: 43, offset: 8971},
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
			pos:  position{line: 286, col: 1, offset: 8981},
			expr: &zeroOrMoreExpr{
				pos: position{line: 286, col: 6, offset: 8988},
				expr: &choiceExpr{
					pos: position{line: 286, col: 8, offset: 8990},
					alternatives: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 286, col: 8, offset: 8990},
							name: "Whitespace",
						},
						&ruleRefExpr{
							pos:  position{line: 286, col: 21, offset: 9003},
							name: "EOL",
						},
						&ruleRefExpr{
							pos:  position{line: 286, col: 27, offset: 9009},
							name: "Comment",
						},
					},
				},
			},
		},
		{
			name: "_",
			pos:  position{line: 287, col: 1, offset: 9021},
			expr: &zeroOrMoreExpr{
				pos: position{line: 287, col: 5, offset: 9027},
				expr: &choiceExpr{
					pos: position{line: 287, col: 7, offset: 9029},
					alternatives: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 287, col: 7, offset: 9029},
							name: "Whitespace",
						},
						&ruleRefExpr{
							pos:  position{line: 287, col: 20, offset: 9042},
							name: "MultiLineCommentNoLineTerminator",
						},
					},
				},
			},
		},
		{
			name: "Whitespace",
			pos:  position{line: 289, col: 1, offset: 9081},
			expr: &charClassMatcher{
				pos:        position{line: 289, col: 14, offset: 9096},
				val:        "[ \\t\\r]",
				chars:      []rune{' ', '\t', '\r'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "EOL",
			pos:  position{line: 290, col: 1, offset: 9105},
			expr: &litMatcher{
				pos:        position{line: 290, col: 7, offset: 9113},
				val:        "\n",
				ignoreCase: false,
			},
		},
		{
			name: "EOS",
			pos:  position{line: 291, col: 1, offset: 9119},
			expr: &choiceExpr{
				pos: position{line: 291, col: 7, offset: 9127},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 291, col: 7, offset: 9127},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 291, col: 7, offset: 9127},
								name: "__",
							},
							&litMatcher{
								pos:        position{line: 291, col: 10, offset: 9130},
								val:        ";",
								ignoreCase: false,
							},
						},
					},
					&seqExpr{
						pos: position{line: 291, col: 16, offset: 9136},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 291, col: 16, offset: 9136},
								name: "_",
							},
							&zeroOrOneExpr{
								pos: position{line: 291, col: 18, offset: 9138},
								expr: &ruleRefExpr{
									pos:  position{line: 291, col: 18, offset: 9138},
									name: "SingleLineComment",
								},
							},
							&ruleRefExpr{
								pos:  position{line: 291, col: 37, offset: 9157},
								name: "EOL",
							},
						},
					},
					&seqExpr{
						pos: position{line: 291, col: 43, offset: 9163},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 291, col: 43, offset: 9163},
								name: "__",
							},
							&ruleRefExpr{
								pos:  position{line: 291, col: 46, offset: 9166},
								name: "EOF",
							},
						},
					},
				},
			},
		},
		{
			name: "EOF",
			pos:  position{line: 293, col: 1, offset: 9173},
			expr: &notExpr{
				pos: position{line: 293, col: 7, offset: 9181},
				expr: &anyMatcher{
					line: 293, col: 8, offset: 9182,
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

func (c *current) onIdentifier1(ident interface{}) (interface{}, error) {

	astIdent := ast.NewIdentifier(c.astPos(), string(c.text))
	if reservedWords[astIdent.Val] {
		return astIdent, errors.New("identifier is a reserved word")
	}
	return astIdent, nil
}

func (p *parser) callonIdentifier1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onIdentifier1(stack["ident"])
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
		// an invalid string literal raises an error in the escape rules,
		// so simply replace the literal with an empty string here to
		// avoid a cascade of errors.
		s = ""
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

func (c *current) onStringLiteral2() (interface{}, error) {

	return ast.NewStringLit(c.astPos(), string(c.text)), nil
}

func (p *parser) callonStringLiteral2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onStringLiteral2()
}

func (c *current) onStringLiteral18() (interface{}, error) {

	return ast.NewStringLit(c.astPos(), "``"), errors.New("string literal not terminated")
}

func (p *parser) callonStringLiteral18() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onStringLiteral18()
}

func (c *current) onDoubleStringEscape5() (interface{}, error) {

	return nil, errors.New("invalid escape character")
}

func (p *parser) callonDoubleStringEscape5() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onDoubleStringEscape5()
}

func (c *current) onSingleStringEscape5() (interface{}, error) {

	return nil, errors.New("invalid escape character")
}

func (p *parser) callonSingleStringEscape5() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onSingleStringEscape5()
}

func (c *current) onOctalEscape6() (interface{}, error) {

	return nil, errors.New("invalid octal escape")
}

func (p *parser) callonOctalEscape6() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOctalEscape6()
}

func (c *current) onHexEscape6() (interface{}, error) {

	return nil, errors.New("invalid hexadecimal escape")
}

func (p *parser) callonHexEscape6() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onHexEscape6()
}

func (c *current) onLongUnicodeEscape2() (interface{}, error) {

	return validateUnicodeEscape(string(c.text), "invalid Unicode escape")

}

func (p *parser) callonLongUnicodeEscape2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onLongUnicodeEscape2()
}

func (c *current) onLongUnicodeEscape13() (interface{}, error) {

	return nil, errors.New("invalid Unicode escape")
}

func (p *parser) callonLongUnicodeEscape13() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onLongUnicodeEscape13()
}

func (c *current) onShortUnicodeEscape2() (interface{}, error) {

	return validateUnicodeEscape(string(c.text), "invalid Unicode escape")

}

func (p *parser) callonShortUnicodeEscape2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onShortUnicodeEscape2()
}

func (c *current) onShortUnicodeEscape9() (interface{}, error) {

	return nil, errors.New("invalid Unicode escape")
}

func (p *parser) callonShortUnicodeEscape9() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onShortUnicodeEscape9()
}

func (c *current) onCharClassMatcher2() (interface{}, error) {

	pos := c.astPos()
	cc := ast.NewCharClassMatcher(pos, string(c.text))
	return cc, nil
}

func (p *parser) callonCharClassMatcher2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onCharClassMatcher2()
}

func (c *current) onCharClassMatcher15() (interface{}, error) {

	return ast.NewCharClassMatcher(c.astPos(), "[]"), errors.New("character class not terminated")
}

func (p *parser) callonCharClassMatcher15() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onCharClassMatcher15()
}

func (c *current) onCharClassEscape5() (interface{}, error) {

	return nil, errors.New("invalid escape character")
}

func (p *parser) callonCharClassEscape5() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onCharClassEscape5()
}

func (c *current) onUnicodeClassEscape5() (interface{}, error) {
	return nil, errors.New("invalid Unicode class escape")
}

func (p *parser) callonUnicodeClassEscape5() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onUnicodeClassEscape5()
}

func (c *current) onUnicodeClassEscape13(ident interface{}) (interface{}, error) {

	if !unicodeClasses[ident.(*ast.Identifier).Val] {
		return nil, errors.New("invalid Unicode class escape")
	}
	return nil, nil

}

func (p *parser) callonUnicodeClassEscape13() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onUnicodeClassEscape13(stack["ident"])
}

func (c *current) onUnicodeClassEscape19() (interface{}, error) {

	return nil, errors.New("Unicode class not terminated")

}

func (p *parser) callonUnicodeClassEscape19() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onUnicodeClassEscape19()
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

func (c *current) onCodeBlock2() (interface{}, error) {

	pos := c.astPos()
	cb := ast.NewCodeBlock(pos, string(c.text))
	return cb, nil
}

func (p *parser) callonCodeBlock2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onCodeBlock2()
}

func (c *current) onCodeBlock7() (interface{}, error) {

	return nil, errors.New("code block not terminated")
}

func (p *parser) callonCodeBlock7() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onCodeBlock7()
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

func rangeTable(class string) *unicode.RangeTable {
	if rt, ok := unicode.Categories[class]; ok {
		return rt
	}
	if rt, ok := unicode.Properties[class]; ok {
		return rt
	}
	if rt, ok := unicode.Scripts[class]; ok {
		return rt
	}

	// cannot happen
	panic(fmt.Sprintf("invalid Unicode class: %s", class))
}
