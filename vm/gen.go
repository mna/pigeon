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
	for i, r := range gr.Rules {
		if i == 0 {
			g.bootstrap(r)
		}
		g.rule(r)
		if g.err != nil {
			break
		}
	}
	return &g.pg, g.err
}

// this Generator-generated program will be used to write the
// runtime ϡtheProgram variable, so it needs to have all the
// information required to build this:
//
// type ϡprogram struct {
// 	instrs []ϡinstr
//
// 	ms []ϡmatcher
// 	as []func(*ϡvm) (interface{}, error)
// 	bs []func(*ϡvm) (bool, error)
// 	ss []string
//
// 	instrToRule []int
// }

type thunkInfo struct {
	Parms  []string
	RuleNm string
	ExprIx int
}

type program struct {
	Init   string
	Instrs []ϡinstr

	Ms  []ast.Expression
	Ss  []string
	As  []*thunkInfo
	Bs  []*thunkInfo
	mss map[string]int
	mms map[string]int

	ruleNmIx        int
	ruleDisplayNmIx int
	exprIx          int
	parmsSet        [][]string
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
	// store the rule's Identifier and Display name in the strings array
	g.pg.ruleNmIx = g.pg.string(r.Name.Val)
	g.pg.ruleDisplayNmIx = g.pg.ruleNmIx
	if r.DisplayName != nil {
		g.pg.ruleDisplayNmIx = g.pg.string(r.DisplayName.Val)
	}
	g.pg.exprIx = 0

	return g.expr(r.Expr)
}

func (g *Generator) expr(expr ast.Expression) int {
	g.pg.exprIx++
	if len(g.pg.parmsSet) < g.pg.exprIx {
		g.pg.parmsSet = append(g.pg.parmsSet, nil)
	}
	defer func() {
		g.pg.exprIx--
	}()

	switch expr := expr.(type) {
	case *ast.ActionExpr:
		return g.action(expr)
	case *ast.AndCodeExpr:
		return g.andNotCode(expr.Code, true)
	case *ast.AndExpr:
		return g.andNot(expr.Expr, true)
	case *ast.AnyMatcher:
		return g.anyMatcher(expr)
	case *ast.CharClassMatcher:
		return g.charClassMatcher(expr)
	case *ast.ChoiceExpr:
		return g.choice(expr)
	case *ast.LabeledExpr:
		return g.labeled(expr)
	case *ast.LitMatcher:
		return g.litMatcher(expr)
	case *ast.NotCodeExpr:
		return g.andNotCode(expr.Code, false)
	case *ast.NotExpr:
		return g.andNot(expr.Expr, false)
	case *ast.OneOrMoreExpr:
		return g.oneOrMore(expr)
	case *ast.RuleRefExpr:
		return g.ruleRef(expr)
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

func (g *Generator) andNotCode(code *ast.CodeBlock, and bool) int {
	th := &thunkInfo{
		Parms:  g.pg.parmsSet[g.pg.exprIx],
		RuleNm: g.pg.Ss[g.pg.ruleNmIx],
		ExprIx: g.pg.exprIx,
	}
	g.pg.Bs = append(g.pg.Bs, th)

	start := g.encode(ϡopCallB, len(g.pg.Bs)-1)
	if and {
		g.encode(ϡopNilIfT)
	} else {
		g.encode(ϡopNilIfF)
	}
	g.encode(ϡopReturn)
	return start
}

func (g *Generator) action(e *ast.ActionExpr) int {
	ix := g.expr(e.Expr)

	th := &thunkInfo{
		Parms:  g.pg.parmsSet[g.pg.exprIx],
		RuleNm: g.pg.Ss[g.pg.ruleNmIx],
		ExprIx: g.pg.exprIx,
	}
	g.pg.As = append(g.pg.As, th)

	start := g.encode(ϡopPush, ϡpstackID)
	g.encode(ϡopPush, ϡistackID, ix)
	g.encode(ϡopCall)
	g.encodeJumpDelta(ϡopJumpIfF, +3)
	g.encode(ϡopCallA, len(g.pg.As)-1)
	g.encode(ϡopReturn)
	g.encode(ϡopPop, ϡpstackID)
	g.encode(ϡopReturn)
	return start
}

func (g *Generator) labeled(e *ast.LabeledExpr) int {
	ix := g.expr(e.Expr)
	lbl := e.Label.Val
	lblIx := g.pg.string(lbl)

	setIx := g.pg.exprIx - 1
	g.pg.parmsSet[setIx] = append(g.pg.parmsSet[setIx], lbl)

	start := g.encode(ϡopPush, ϡistackID, ix)
	g.encode(ϡopCall)
	g.encode(ϡopStoreIfT, lblIx)
	g.encode(ϡopReturn)
	return start
}

func (g *Generator) andNot(subExpr ast.Expression, and bool) int {
	ix := g.expr(subExpr)
	start := g.encode(ϡopPush, ϡpstackID)
	g.encode(ϡopPush, ϡistackID, ix)
	g.encode(ϡopCall)
	if and {
		g.encode(ϡopNilIfT)
	} else {
		g.encode(ϡopNilIfF)
	}
	g.encode(ϡopRestore)
	g.encode(ϡopReturn)
	return start
}

func (g *Generator) ruleRef(e *ast.RuleRefExpr) int {
	nm := e.Name.Val
	ix := g.pg.string(nm)

	start := g.encode(ϡopPlaceholder, ix)
	g.encode(ϡopCall)
	g.encode(ϡopReturn)
	return start
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
func (g *Generator) bootstrap(r *ast.Rule) {
	nm := r.Name.Val
	ix := g.pg.string(nm)

	g.encode(ϡopPlaceholder, ix)
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
