package vm

import (
	"errors"
	"fmt"
	"io"

	"github.com/PuerkitoBio/pigeon/ast"
)

var (
	// errNoRule is returned when the grammar to generate has no rule.
	errNoRule = errors.New("grammar has no rule")
)

// NewGenerator creates a Generator that writes to w.
func NewGenerator(w io.Writer) *Generator {
	g := &Generator{w: w, RecvName: "c"}
	return g
}

// Generator generates the PEG parser for a provided grammar.
type Generator struct {
	// options
	RecvName string

	w   io.Writer
	err error

	pg program
}

// Generate generates the PEG parser's code to g.w for the provider
// grammar gr.
func (g *Generator) Generate(gr *ast.Grammar) error {
	if len(gr.Rules) == 0 {
		return errNoRule
	}

	g.pg.Init = gr.Init.Val
	g.bootstrap()
	for _, r := range gr.Rules {
		g.rule(r)
	}

	return g.err
}

type program struct {
	Instrs []ϡinstr
	Init   string
	Ms     []ast.Expression
	Ss     []string

	ruleNmIx int
	exprIx   int

	mss map[string]int
	mms map[string]int
}

func (pg *program) matcher(raw string, expr ast.Expression) int {
	if pg.mms == nil {
		pg.mms = make(map[string]int)
	}
	ix, ok := pg.mms[raw]
	if !ok {
		pg.Ms = append(pg.Ms, expr)
		ix = len(pg.Ms) - 1
		pg.mms[raw] = ix
	}
	return ix
}

func (pg *program) string(s string) int {
	if pg.mss == nil {
		pg.mss = make(map[string]int)
	}
	ix, ok := pg.mss[s]
	if !ok {
		pg.Ss = append(pg.Ss, s)
		ix = len(pg.Ss) - 1
		pg.mss[s] = ix
	}
	return ix
}

func (g *Generator) rule(r *ast.Rule) {
	// store the rule's Identifier or Display name in the strings array
	s := r.Name.Val
	if r.DisplayName != nil {
		s = r.DisplayName.Val
	}
	g.pg.ruleNmIx = g.pg.string(s)
	g.pg.exprIx = 0

	g.expr(r.Expr)
}

func (g *Generator) expr(expr ast.Expression) {
	g.pg.exprIx++
	switch expr := expr.(type) {
	case *ast.ActionExpr:
	case *ast.AndCodeExpr:
	case *ast.AndExpr:
	case *ast.AnyMatcher:
		g.anyMatcher(expr)
	case *ast.CharClassMatcher:
	case *ast.ChoiceExpr:
	case *ast.LabeledExpr:
	case *ast.LitMatcher:
	case *ast.NotCodeExpr:
	case *ast.NotExpr:
	case *ast.OneOrMoreExpr:
	case *ast.RuleRefExpr:
	case *ast.SeqExpr:
	case *ast.ZeroOrMoreExpr:
	case *ast.ZeroOrOneExpr:
	default:
		g.err = fmt.Errorf("unknown expression type %T", expr)
	}
}

func (g *Generator) anyMatcher(e *ast.AnyMatcher) {
	mIx := g.pg.matcher(e.Val, e)
	g.matcher(mIx)
}

// matcher generates the instructions to call the matcher at index mIx.
func (g *Generator) matcher(mIx int) {
	g.encode(ϡopPush, ϡpstackID)
	g.encode(ϡopMatch, mIx)
	g.encode(ϡopRestoreIfF)
	g.encode(ϡopReturn)
}

// bootstrap adds the bootstrapping opcode sequence to the program's
// instructions.
func (g *Generator) bootstrap() {
	g.encode(ϡopPush, ϡistackID, 3)
	g.encode(ϡopCall)
	g.encode(ϡopExit)
}

func (g *Generator) encode(op ϡop, args ...int) {
	if g.err == nil {
		instr, err := ϡencodeInstr(ϡopExit)
		g.err = err
		g.pg.Instrs = append(g.pg.Instrs, instr...)
	}
}
