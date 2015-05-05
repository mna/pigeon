package vm

import (
	"errors"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/pigeon/ast"
)

var (
	// errNoRule is returned when the grammar to generate has no rule.
	errNoRule = errors.New("grammar has no rule")
)

// NewGenerator creates a Generator that writes to w.
func NewGenerator(w io.Writer) *Generator {
	g := &Generator{w: w, RecvrName: "c"}
	return g
}

// Generator generates the PEG parser for a provided grammar.
type Generator struct {
	// options
	RecvrName string

	w   io.Writer
	err error

	pg program
}

// Generate generates the PEG parser's code to g.w for the provided
// grammar gr.
func (g *Generator) Generate(gr *ast.Grammar) error {
	pg, err := g.toProgram(gr)
	if err != nil {
		return err
	}

	return g.write(pg)
}

func (g *Generator) write(pg *program) error {
	// first write the template-generated code
	if err := tpl.Execute(g.w, pg); err != nil {
		return err
	}
	// then write the static code
	_, err := fmt.Fprint(g.w, staticCode)
	return err
}

func (g *Generator) toProgram(gr *ast.Grammar) (pg *program, err error) {
	if len(gr.Rules) == 0 {
		return nil, errNoRule
	}

	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
		}
	}()

	g.pg.RecvrNm = g.RecvrName
	g.pg.Now = time.Now()

	if gr.Init != nil {
		g.pg.Init = unwrapCode(gr.Init.Val)
	}
	g.pg.ruleNmStartIx = make(map[uint16]uint16)
	g.pg.ruleNmEntryIx = make(map[uint16]uint16)
	g.pg.ruleNmToDisNm = make(map[uint16]uint16)
	for i, r := range gr.Rules {
		if i == 0 {
			g.bootstrap(r)
		}
		g.rule(r)
		if g.err != nil {
			break
		}
	}

	g.fillPlaceholders()
	g.instrToRule()

	if g.err != nil {
		return nil, g.err
	}
	return &g.pg, nil
}

type thunkInfo struct {
	Parms  []string
	RuleNm string
	ExprIx uint16
	Code   string
}

type program struct {
	Now     time.Time
	RecvrNm string
	Init    string
	Instrs  []ϡinstr

	Ms []ast.Expression
	As []*thunkInfo
	Bs []*thunkInfo
	Ss []string

	mss map[string]uint16 // reverse map of string to index in Ss
	mms map[string]uint16 // reverse map of matcher's raw value to index in Ms

	ruleNmIx      uint16
	exprIx        uint16
	parmsSet      [][]string        // stack of parms set for code blocks
	ruleNmStartIx map[uint16]uint16 // rule name ix to first rule instr ix
	ruleNmEntryIx map[uint16]uint16 // rule name ix to entry point instr ix
	ruleNmToDisNm map[uint16]uint16 // rule name ix to rule display name ix
}

func (pg *program) pushParmsSet() {
	pg.parmsSet = append(pg.parmsSet, nil)
}

func (pg *program) popParmsSet() {
	pg.parmsSet = pg.parmsSet[:len(pg.parmsSet)-1]
}

func (pg *program) peekParmsSet() []string {
	return pg.parmsSet[len(pg.parmsSet)-1]
}

func (pg *program) pushParm(v string) {
	ix := len(pg.parmsSet) - 1
	pg.parmsSet[ix] = append(pg.parmsSet[ix], v)
}

func (pg *program) matcher(raw string, expr ast.Expression) uint16 {
	if pg.mms == nil {
		pg.mms = make(map[string]uint16)
	}
	ix, ok := pg.mms[raw]
	if !ok {
		if len(pg.Ms) >= math.MaxUint16 {
			panic("too many matchers")
		}
		pg.Ms = append(pg.Ms, expr)
		ix = uint16(len(pg.Ms) - 1)
		pg.mms[raw] = ix
	}
	return ix
}

func (pg *program) string(s string) uint16 {
	if pg.mss == nil {
		pg.mss = make(map[string]uint16)
	}
	ix, ok := pg.mss[s]
	if !ok {
		if len(pg.Ss) >= math.MaxUint16 {
			panic("too many strings")
		}
		pg.Ss = append(pg.Ss, s)
		ix = uint16(len(pg.Ss) - 1)
		pg.mss[s] = ix
	}
	return ix
}

func (g *Generator) rule(r *ast.Rule) uint16 {
	// store the rule's Identifier and Display name in the strings array
	g.pg.ruleNmIx = g.pg.string(r.Name.Val)
	disNmIx := g.pg.ruleNmIx
	if r.DisplayName != nil {
		disNmIx = g.pg.string(r.DisplayName.Val)
	}
	g.pg.exprIx = 0

	start := uint16(len(g.pg.Instrs))

	g.pg.pushParmsSet()
	entry := g.expr(r.Expr)
	g.pg.popParmsSet()

	g.pg.ruleNmEntryIx[g.pg.ruleNmIx] = entry
	g.pg.ruleNmStartIx[g.pg.ruleNmIx] = start
	g.pg.ruleNmToDisNm[g.pg.ruleNmIx] = disNmIx
	return start
}

func (g *Generator) expr(expr ast.Expression) uint16 {
	g.pg.exprIx++

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
		return g.repetition(expr.Expr, false)
	case *ast.RuleRefExpr:
		return g.ruleRef(expr)
	case *ast.SeqExpr:
		return g.sequence(expr)
	case *ast.ZeroOrMoreExpr:
		return g.repetition(expr.Expr, true)
	case *ast.ZeroOrOneExpr:
		return g.zeroOrOne(expr)
	default:
		g.err = fmt.Errorf("unknown expression type %T", expr)
	}
	return 0
}

func (g *Generator) anyMatcher(e *ast.AnyMatcher) uint16 {
	// Val is a dot `.` for the any matcher
	mIx := g.pg.matcher(e.Val, e)
	return g.matcher(mIx)
}

func (g *Generator) charClassMatcher(e *ast.CharClassMatcher) uint16 {
	// Val is the raw char class literal, including [] and optional `i`
	mIx := g.pg.matcher(e.Val, e)
	return g.matcher(mIx)
}

func (g *Generator) litMatcher(e *ast.LitMatcher) uint16 {
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

func (g *Generator) andNotCode(code *ast.CodeBlock, and bool) uint16 {
	th := &thunkInfo{
		Parms:  g.pg.peekParmsSet(),
		RuleNm: g.pg.Ss[g.pg.ruleNmIx],
		ExprIx: g.pg.exprIx,
		Code:   unwrapCode(code.Val),
	}
	if len(g.pg.Bs) >= math.MaxUint16 {
		panic("too many code predicates")
	}
	g.pg.Bs = append(g.pg.Bs, th)

	start := g.encode(ϡopCallB, uint16(len(g.pg.Bs)-1))
	if and {
		g.encode(ϡopNilIfT)
	} else {
		g.encode(ϡopNilIfF)
	}
	g.encode(ϡopReturn)
	return start
}

func (g *Generator) action(e *ast.ActionExpr) uint16 {
	actIx := g.pg.exprIx
	ix := g.expr(e.Expr)

	th := &thunkInfo{
		Parms:  g.pg.peekParmsSet(),
		RuleNm: g.pg.Ss[g.pg.ruleNmIx],
		ExprIx: actIx,
		Code:   unwrapCode(e.Code.Val),
	}
	if len(g.pg.As) >= math.MaxUint16 {
		panic("too many actions")
	}
	g.pg.As = append(g.pg.As, th)

	start := g.encode(ϡopPush, ϡpstackID)
	g.encode(ϡopPush, ϡistackID, ix)
	g.encode(ϡopCall)
	g.encodeJumpDelta(ϡopJumpIfF, +3)
	g.encode(ϡopCallA, uint16(len(g.pg.As)-1))
	g.encode(ϡopReturn)
	g.encode(ϡopPop, ϡpstackID)
	g.encode(ϡopReturn)
	return start
}

func (g *Generator) labeled(e *ast.LabeledExpr) uint16 {
	lbl := e.Label.Val
	lblIx := g.pg.string(lbl)
	g.pg.pushParm(lbl)

	g.pg.pushParmsSet()
	ix := g.expr(e.Expr)
	g.pg.popParmsSet()

	start := g.encode(ϡopPush, ϡistackID, ix)
	g.encode(ϡopPush, ϡastackID)
	g.encode(ϡopCall)
	g.encode(ϡopPop, ϡastackID)
	g.encode(ϡopStoreIfT, lblIx)
	g.encode(ϡopReturn)
	return start
}

func (g *Generator) andNot(subExpr ast.Expression, and bool) uint16 {
	g.pg.pushParmsSet()
	ix := g.expr(subExpr)
	g.pg.popParmsSet()

	start := g.encode(ϡopPush, ϡpstackID)
	g.encode(ϡopPush, ϡistackID, ix)
	g.encode(ϡopPush, ϡastackID)
	g.encode(ϡopCall)
	g.encode(ϡopPop, ϡastackID)
	if and {
		g.encode(ϡopNilIfT)
	} else {
		g.encode(ϡopNilIfF)
	}
	g.encode(ϡopRestore)
	g.encode(ϡopReturn)
	return start
}

func (g *Generator) ruleRef(e *ast.RuleRefExpr) uint16 {
	nm := e.Name.Val
	ix := g.pg.string(nm)

	start := g.encode(ϡopPlaceholder, ix, 0)
	g.encode(ϡopPush, ϡastackID)
	g.encode(ϡopCall)
	g.encode(ϡopPop, ϡastackID)
	g.encode(ϡopReturn)
	return start
}

func (g *Generator) zeroOrOne(e *ast.ZeroOrOneExpr) uint16 {
	g.pg.pushParmsSet()
	ix := g.expr(e.Expr)
	g.pg.popParmsSet()

	start := g.encode(ϡopPush, ϡistackID, ix)
	g.encode(ϡopPush, ϡastackID)
	g.encode(ϡopCall)
	g.encode(ϡopPop, ϡastackID)
	g.encodeJumpDelta(ϡopPopVJumpIfF, +2)
	g.encode(ϡopReturn)
	g.encode(ϡopPush, ϡvstackID, ϡvValNil)
	g.encode(ϡopReturn)
	return start
}

func (g *Generator) repetition(subExpr ast.Expression, zeroOk bool) uint16 {
	g.pg.pushParmsSet()
	ix := g.expr(subExpr)
	g.pg.popParmsSet()

	vVal := ϡvValFailed
	if zeroOk {
		vVal = ϡvValEmpty
	}
	start := g.encode(ϡopPush, ϡvstackID, vVal)
	g.encode(ϡopPush, ϡistackID, ix)
	g.encode(ϡopPush, ϡastackID)
	g.encode(ϡopCall)
	g.encode(ϡopPop, ϡastackID)
	g.encodeJumpDelta(ϡopPopVJumpIfF, +3)
	g.encode(ϡopCumulOrF)
	g.encodeJumpDelta(ϡopJump, -6)
	g.encode(ϡopReturn)
	return start
}

func (g *Generator) choice(e *ast.ChoiceExpr) uint16 {
	// first generate code for each of the choice's expressions
	indices := make([]uint16, len(e.Alternatives)+1)
	for i, se := range e.Alternatives {
		g.pg.pushParmsSet()
		indices[i+1] = g.expr(se)
		g.pg.popParmsSet()
	}

	// then generate the sequence's instructions
	indices[0] = ϡlstackID
	start := g.encode(ϡopPush, indices...)
	g.encodeJumpDelta(ϡopTakeLOrJump, +7)
	g.encode(ϡopPush, ϡastackID)
	g.encode(ϡopCall)
	g.encode(ϡopPop, ϡastackID)
	g.encodeJumpDelta(ϡopJumpIfT, +4)
	g.encode(ϡopPop, ϡvstackID)
	g.encodeJumpDelta(ϡopJump, -6)
	g.encode(ϡopPush, ϡvstackID, ϡvValFailed)
	g.encode(ϡopPop, ϡlstackID)
	g.encode(ϡopReturn)
	return start
}

func (g *Generator) sequence(e *ast.SeqExpr) uint16 {
	// first generate code for each of the sequence's expressions
	indices := make([]uint16, len(e.Exprs)+1)
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
func (g *Generator) matcher(mIx uint16) uint16 {
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

	g.encode(ϡopPlaceholder, ix, 0)
	g.encode(ϡopPush, ϡastackID)
	g.encode(ϡopCall)
	g.encode(ϡopExit)
}

func (g *Generator) fillPlaceholders() {
	if g.err != nil {
		return
	}

	for i, instr := range g.pg.Instrs {
		if instr.op == ϡopPlaceholder {
			ix := g.pg.ruleNmEntryIx[instr.args[0]]
			if ix == 0 {
				g.err = fmt.Errorf("undefined rule %q", g.pg.Ss[instr.args[0]])
			}
			instr.op = ϡopPush
			instr.args[0] = ϡistackID
			instr.args[1] = ix
			g.pg.Instrs[i] = instr
		}
	}
}

func (g *Generator) instrToRule() {
	if g.err != nil {
		return
	}

	// rule start index is necessarily > 0 because of the bootstrap sequence
	var startStartIx uint16
	for ruleNmIx, startIx := range g.pg.ruleNmStartIx {
		if startIx < startStartIx || startStartIx == 0 {
			startStartIx = startIx
		}
		g.pg.Instrs[startIx].ruleNmIx = g.pg.ruleNmToDisNm[ruleNmIx]
	}

	// fill the blanks
	fillIx := -1
	for i := uint16(0); i < uint16(len(g.pg.Instrs)); i++ {
		if ruleNmIx := g.pg.Instrs[i].ruleNmIx; ruleNmIx != 0 || i == startStartIx {
			fillIx = int(ruleNmIx)
		}
		g.pg.Instrs[i].ruleNmIx = uint16(fillIx)
	}
}

func (g *Generator) encodeJumpDelta(op ϡop, delta int) uint16 {
	return g.encode(op, uint16(delta+len(g.pg.Instrs)))
}

func (g *Generator) encode(op ϡop, args ...uint16) uint16 {
	if g.err == nil {
		start := len(g.pg.Instrs)
		if start >= math.MaxUint16 {
			panic("too many instructions")
		}
		g.pg.Instrs = append(g.pg.Instrs, ϡinstr{op: op, args: args})
		return uint16(start)
	}
	return 0
}

func unwrapCode(val string) string {
	return strings.TrimSpace(val)[1 : len(val)-1]
}
