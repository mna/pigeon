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

func (g *Generator) expr(expr ast.Expression) int {
	g.pg.exprIx++
	switch expr := expr.(type) {
	case *ast.ActionExpr:
	case *ast.AndCodeExpr:
	case *ast.AndExpr:
	case *ast.AnyMatcher:
		return g.anyMatcher(expr)
	case *ast.CharClassMatcher:
		return g.charClassMatcher(expr)
	case *ast.ChoiceExpr:
		return g.choice(expr)
	case *ast.LabeledExpr:
	case *ast.LitMatcher:
		return g.litMatcher(expr)
	case *ast.NotCodeExpr:
	case *ast.NotExpr:
	case *ast.OneOrMoreExpr:
		return g.oneOrMore(expr)
	case *ast.RuleRefExpr:
	case *ast.SeqExpr:
		return g.sequence(expr)
	case *ast.ZeroOrMoreExpr:
		return g.zeroOrMore(expr)
	case *ast.ZeroOrOneExpr:
		return g.zeroOrOne(expr)
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

func (g *Generator) zeroOrOne(e *ast.ZeroOrOneExpr) int {
	ix := g.expr(e.Expr)

	start := g.encode(ϡopPush, ϡistackID, ix)
	g.encode(ϡopCall)
	g.encodeJumpDelta(ϡopPopVJumpIfF, +2)
	g.encode(ϡopReturn)
	g.encode(ϡopPush, ϡvstackID, ϡvValNil)
	g.encode(ϡopReturn)
	return start
}

func (g *Generator) oneOrMore(e *ast.OneOrMoreExpr) int {
	ix := g.expr(e.Expr)

	start := g.encode(ϡopPush, ϡvstackID, ϡvValFailed)
	g.encode(ϡopPush, ϡistackID, ix)
	g.encode(ϡopCall)
	g.encodeJumpDelta(ϡopPopVJumpIfF, +3)
	g.encode(ϡopCumulOrF)
	g.encodeJumpDelta(ϡopJump, -4)
	g.encode(ϡopReturn)
	return start
}

func (g *Generator) zeroOrMore(e *ast.ZeroOrMoreExpr) int {
	ix := g.expr(e.Expr)

	start := g.encode(ϡopPush, ϡvstackID, ϡvValEmpty)
	g.encode(ϡopPush, ϡistackID, ix)
	g.encode(ϡopCall)
	g.encodeJumpDelta(ϡopPopVJumpIfF, +3)
	g.encode(ϡopCumulOrF)
	g.encodeJumpDelta(ϡopJump, -4)
	g.encode(ϡopReturn)
	return start
}

func (g *Generator) choice(e *ast.ChoiceExpr) int {
	// first generate code for each of the choice's expressions
	indices := make([]int, len(e.Alternatives)+1)
	for i, se := range e.Alternatives {
		indices[i+1] = g.expr(se)
	}

	// then generate the sequence's instructions
	indices[0] = ϡlstackID
	start := g.encode(ϡopPush, indices...)
	g.encodeJumpDelta(ϡopTakeLOrJump, +4)
	g.encode(ϡopCall)
	g.encodeJumpDelta(ϡopJumpIfT, +2)
	g.encodeJumpDelta(ϡopJump, -3)
	g.encode(ϡopPop, ϡlstackID)
	g.encode(ϡopReturn)
	return start
}

func (g *Generator) sequence(e *ast.SeqExpr) int {
	// first generate code for each of the sequence's expressions
	indices := make([]int, len(e.Exprs)+1)
	for i, se := range e.Exprs {
		indices[i+1] = g.expr(se)
	}

	// then generate the sequence's instructions
	start := g.encode(ϡopPush, ϡpstackID)
	g.encode(ϡopPush, ϡvstackID, ϡvValFailed)
	indices[0] = ϡlstackID
	g.encode(ϡopPush, indices...)
	g.encodeJumpDelta(ϡopTakeLOrJump, +5)
	g.encode(ϡopCall)
	g.encode(ϡopCumulOrF)
	g.encodeJumpDelta(ϡopJumpIfF, +2)
	g.encodeJumpDelta(ϡopJump, -4)
	g.encode(ϡopPop, ϡlstackID)
	g.encode(ϡopRestoreIfF)
	g.encode(ϡopReturn)
	return start
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

func (g *Generator) encodeJumpDelta(op ϡop, delta int) int {
	return g.encode(op, delta+len(g.pg.Instrs))
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
