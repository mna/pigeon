package vm

import (
	"errors"
	"fmt"
	"io"
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

func (g *Generator) toProgram(gr *ast.Grammar) (*program, error) {
	if len(gr.Rules) == 0 {
		return nil, errNoRule
	}

	g.pg.RecvrNm = g.RecvrName
	g.pg.Now = time.Now()

	if gr.Init != nil {
		g.pg.Init = unwrapCode(gr.Init.Val)
	}
	g.pg.ruleNmStartIx = make(map[int]int)
	g.pg.ruleNmEntryIx = make(map[int]int)
	g.pg.ruleNmToDisNm = make(map[int]int)
	for i, r := range gr.Rules {
		if i == 0 {
			g.bootstrap(r)
		}
		g.rule(r)
		if g.err != nil {
			break
		}
	}

	if g.err != nil {
		return nil, g.err
	}

	g.fillPlaceholders()
	g.instrToRule()
	return &g.pg, nil
}

type thunkInfo struct {
	Parms  []string
	RuleNm string
	ExprIx int
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

	InstrToRule []int

	mss map[string]int // reverse map of string to index in Ss
	mms map[string]int // reverse map of matcher's raw value to index in Ms

	ruleNmIx      int
	exprIx        int
	parmsSet      [][]string  // stack of parms set for code blocks
	ruleNmStartIx map[int]int // rule name ix to first rule instr ix
	ruleNmEntryIx map[int]int // rule name ix to entry point instr ix
	ruleNmToDisNm map[int]int // rule name ix to rule display name ix
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
	disNmIx := g.pg.ruleNmIx
	if r.DisplayName != nil {
		disNmIx = g.pg.string(r.DisplayName.Val)
	}
	g.pg.exprIx = 0

	start := len(g.pg.Instrs)

	g.pg.pushParmsSet()
	entry := g.expr(r.Expr)
	g.pg.popParmsSet()

	g.pg.ruleNmEntryIx[g.pg.ruleNmIx] = entry
	g.pg.ruleNmStartIx[g.pg.ruleNmIx] = start
	g.pg.ruleNmToDisNm[g.pg.ruleNmIx] = disNmIx
	return start
}

func (g *Generator) expr(expr ast.Expression) int {
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
		Parms:  g.pg.peekParmsSet(),
		RuleNm: g.pg.Ss[g.pg.ruleNmIx],
		ExprIx: g.pg.exprIx,
		Code:   unwrapCode(code.Val),
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
	actIx := g.pg.exprIx
	ix := g.expr(e.Expr)

	th := &thunkInfo{
		Parms:  g.pg.peekParmsSet(),
		RuleNm: g.pg.Ss[g.pg.ruleNmIx],
		ExprIx: actIx,
		Code:   unwrapCode(e.Code.Val),
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

func (g *Generator) andNot(subExpr ast.Expression, and bool) int {
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

func (g *Generator) ruleRef(e *ast.RuleRefExpr) int {
	nm := e.Name.Val
	ix := g.pg.string(nm)

	start := g.encode(ϡopPlaceholder, ix)
	g.encode(ϡopPush, ϡastackID)
	g.encode(ϡopCall)
	g.encode(ϡopPop, ϡastackID)
	g.encode(ϡopReturn)
	return start
}

func (g *Generator) zeroOrOne(e *ast.ZeroOrOneExpr) int {
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

func (g *Generator) repetition(subExpr ast.Expression, zeroOk bool) int {
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

func (g *Generator) choice(e *ast.ChoiceExpr) int {
	// first generate code for each of the choice's expressions
	indices := make([]int, len(e.Alternatives)+1)
	for i, se := range e.Alternatives {
		g.pg.pushParmsSet()
		indices[i+1] = g.expr(se)
		g.pg.popParmsSet()
	}

	// then generate the sequence's instructions
	indices[0] = ϡlstackID
	start := g.encode(ϡopPush, indices...)
	g.encodeJumpDelta(ϡopTakeLOrJump, +6)
	g.encode(ϡopPush, ϡastackID)
	g.encode(ϡopCall)
	g.encode(ϡopPop, ϡastackID)
	g.encodeJumpDelta(ϡopJumpIfT, +2)
	g.encodeJumpDelta(ϡopJump, -5)
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
	g.encode(ϡopPush, ϡastackID)
	g.encode(ϡopCall)
	g.encode(ϡopExit)
}

func (g *Generator) fillPlaceholders() {
	var op ϡop
	var n, nmIx int
	for i, instr := range g.pg.Instrs {
		if n > 0 {
			n -= 4
			continue
		}
		op, n, nmIx, _, _ = instr.decode()
		n -= 3
		if op == ϡopPlaceholder {
			ix := g.pg.ruleNmEntryIx[nmIx]
			newInstrs, _ := ϡencodeInstr(ϡopPush, ϡistackID, ix)
			g.pg.Instrs[i] = newInstrs[0]
		}
	}
}

func (g *Generator) instrToRule() {
	g.pg.InstrToRule = make([]int, len(g.pg.Instrs))
	// rule start index is necessarily > 0 because of the bootstrap sequence
	startStartIx := 0
	for ruleNmIx, startIx := range g.pg.ruleNmStartIx {
		if startIx < startStartIx || startStartIx == 0 {
			startStartIx = startIx
		}
		g.pg.InstrToRule[startIx] = g.pg.ruleNmToDisNm[ruleNmIx]
	}

	// fill the blanks
	fillIx := -1
	for i := 0; i < len(g.pg.InstrToRule); i++ {
		if ruleNmIx := g.pg.InstrToRule[i]; ruleNmIx != 0 || i == startStartIx {
			fillIx = ruleNmIx
		}
		g.pg.InstrToRule[i] = fillIx
	}
}

func (g *Generator) encodeJumpDelta(op ϡop, delta int) int {
	return g.encode(op, delta+len(g.pg.Instrs))
}

func (g *Generator) encode(op ϡop, args ...int) int {
	if g.err == nil {
		instr, err := ϡencodeInstr(op, args...)
		g.err = err
		start := len(g.pg.Instrs)
		g.pg.Instrs = append(g.pg.Instrs, instr...)
		return start
	}
	return 0
}

func unwrapCode(val string) string {
	return strings.TrimSpace(val)[1 : len(val)-1]
}
