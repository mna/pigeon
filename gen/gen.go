package gen

import (
	"errors"
	"io"

	"github.com/PuerkitoBio/pigeon/ast"
)

var (
	// errNoRule is returned when the grammar to generate has no rule.
	errNoRule = errors.New("grammar has no rule")
)

// Option is a function that can set an option on the code generator.
// It returns the previous setting as an Option.
type Option func(*generator) Option

// ReceiverName returns an option that specifies the receiver name to
// use for the "current" struct (which is the struct on which all code blocks
// except the initializer are generated).
func ReceiverName(nm string) Option {
	return func(g *generator) Option {
		prev := g.recvName
		g.recvName = nm
		return ReceiverName(prev)
	}
}

// Generate generates the PEG parser using the provided grammar. The code is
// written to the specified w.
func Generate(w io.Writer, gr *ast.Grammar, opts ...Option) error {
	g := &generator{w: w, recvName: "c"}
	return g.setOptions(opts).generate(gr)
}

// generator generates the PEG parser for a provided grammar.
type generator struct {
	w   io.Writer
	err error

	// options
	recvName string

	pg program
}

// generate generates the PEG parser's code to g.w for the provider
// grammar gr.
func (g *generator) generate(gr *ast.Grammar) error {
	if len(gr.Rules) == 0 {
		return errNoRule
	}

	g.pg.Init = gr.Init.Val
	// if err := g.bootstrap(gr); err != nil {
	// 	return err
	// }
	// for _, r := range gr.Rules {
	// 	g.rule(r)
	// }

	return g.err
}

// setOptions applies the options opts in sequence to the generator. It
// returns the generator so that calls can be chained.
func (g *generator) setOptions(opts []Option) *generator {
	for _, opt := range opts {
		opt(g)
	}
	return g
}
