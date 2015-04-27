package vm

import (
	"errors"
	"fmt"
	"io"
	"strconv"

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
	pg, err := g.toProgram(gr)
	if err != nil {
		return err
	}

	// TODO : return g.write(pg)
	_ = pg
	return nil
}

func (g *Generator) toProgram(gr *ast.Grammar) (*program, error) {
	if len(gr.Rules) == 0 {
		return nil, errNoRule
	}

	g.pg.Init = gr.Init.Val
	if g.bootstrap(); g.err == nil {
		for _, r := range gr.Rules {
			g.rule(r)
			if g.err != nil {
				break
			}
		}
	}
	return &g.pg, g.err
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

func (g *Generator) rule(r *ast.Rule) int {
	// store the rule's Identifier or Display name in the strings array
	s := r.Name.Val
	if r.DisplayName != nil {
		s = r.DisplayName.Val
	}
	g.pg.ruleNmIx = g.pg.string(s)
	g.pg.exprIx = 0

	return g.expr(r.Expr)
}

// TODO : return index of start instruction
func (g *Generator) expr(expr ast.Expression) int {
	g.pg.exprIx++
	switch expr := expr.(type) {
	case *ast.ActionExpr:
	case *ast.AndCodeExpr:
	case *ast.AndExpr:
	case *ast.AnyMatcher:
		g.anyMatcher(expr)
	case *ast.CharClassMatcher:
		g.charClassMatcher(expr)
	case *ast.ChoiceExpr:
	case *ast.LabeledExpr:
	case *ast.LitMatcher:
		g.litMatcher(expr)
	case *ast.NotCodeExpr:
	case *ast.NotExpr:
	case *ast.OneOrMoreExpr:
	case *ast.RuleRefExpr:
	case *ast.SeqExpr:
		g.sequence()
	case *ast.ZeroOrMoreExpr:
	case *ast.ZeroOrOneExpr:
	default:
		g.err = fmt.Errorf("unknown expression type %T", expr)
	}
	return 0
}

func (g *Generator) anyMatcher(e *ast.AnyMatcher) int {
	// Val is a dot `.` for the any matcher
	mIx := g.pg.matcher(e.Val, e)
	return g.matcher(mIx)
}

func (g *Generator) charClassMatcher(e *ast.CharClassMatcher) int {
	// Val is the raw char class literal, including [] and optional `i`
	mIx := g.pg.matcher(e.Val, e)
	return g.matcher(mIx)
}

func (g *Generator) litMatcher(e *ast.LitMatcher) int {
	raw := strconv.Quote(e.Val)
	if e.IgnoreCase {
		raw += "i"
	}

	// raw is the quoted string with the optional `i`, so it can't conflict
	// with any and char class, and different literals are normalized to the
	// same form, so `"a"` and `'a'` is equivalent and the same, single matcher
	// will be used.
	mIx := g.pg.matcher(raw, e)
	return g.matcher(mIx)
}

func (g *Generator) sequence() int {
	return 0
}

// matcher generates the instructions to call the matcher at index mIx.
func (g *Generator) matcher(mIx int) int {
	start := g.encode(ϡopPush, ϡpstackID)
	g.encode(ϡopMatch, mIx)
	g.encode(ϡopRestoreIfF)
	g.encode(ϡopReturn)
	return start
}

// bootstrap adds the bootstrapping opcode sequence to the program's
// instructions.
func (g *Generator) bootstrap() {
	// TODO : that's no good, entry point might not be instr 3 (e.g. if first
	// rule is a sequence)
	g.encode(ϡopPush, ϡistackID, 3)
	g.encode(ϡopCall)
	g.encode(ϡopExit)
}

func (g *Generator) encode(op ϡop, args ...int) int {
	if g.err == nil {
		instr, err := ϡencodeInstr(ϡopExit)
		g.err = err
		start := len(g.pg.Instrs)
		g.pg.Instrs = append(g.pg.Instrs, instr...)
		return start
	}
	return 0
}
