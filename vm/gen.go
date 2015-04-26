package vm

import (
	"errors"
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
}

func (g *Generator) rule(r *ast.Rule) {

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
