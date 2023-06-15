package builder_test

import (
	"strings"
	"testing"

	"github.com/mna/pigeon/ast"
	"github.com/mna/pigeon/bootstrap"
	"github.com/mna/pigeon/builder"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLeftRrecursive(t *testing.T) {
	t.Parallel()
	text := `
	start = expr NEWLINE
    expr = ('-' term / expr '+' term / term)
    term = NUMBER
    foo = NAME+
    bar = NAME*
    baz = NAME?
	`
	p := bootstrap.NewParser()
	grammar, err := p.Parse("", strings.NewReader(text))
	require.NoError(t, err)
	haveLeftRecursion, err := builder.PrepareGrammar(grammar)
	require.NoError(t, err)
	require.True(t, haveLeftRecursion)
	mapRules := make(map[string]*ast.Rule, len(grammar.Rules))
	for _, rule := range grammar.Rules {
		mapRules[rule.Name.Val] = rule
	}
	assert.False(t, mapRules["start"].LeftRecursive)
	assert.True(t, mapRules["expr"].LeftRecursive)
	assert.False(t, mapRules["term"].LeftRecursive)
	assert.False(t, mapRules["foo"].LeftRecursive)
	assert.False(t, mapRules["bar"].LeftRecursive)
	assert.False(t, mapRules["baz"].LeftRecursive)
}

func TestNullable(t *testing.T) {
	t.Parallel()
	text := `
	start = sign NUMBER
    sign = ('-' / '+')?
	`
	p := bootstrap.NewParser()
	grammar, err := p.Parse("", strings.NewReader(text))
	require.NoError(t, err)
	haveLeftRecursion, err := builder.PrepareGrammar(grammar)
	require.NoError(t, err)
	require.False(t, haveLeftRecursion)
	mapRules := make(map[string]*ast.Rule, len(grammar.Rules))
	for _, rule := range grammar.Rules {
		mapRules[rule.Name.Val] = rule
	}
	assert.False(t, mapRules["start"].Nullable)
	assert.True(t, mapRules["sign"].Nullable)
}

func TestAdvancedLeftRrecursive(t *testing.T) {
	t.Parallel()
	text := `
	start = NUMBER / sign start
    sign = '-'?
	`
	p := bootstrap.NewParser()
	grammar, err := p.Parse("", strings.NewReader(text))
	require.NoError(t, err)
	haveLeftRecursion, err := builder.PrepareGrammar(grammar)
	require.NoError(t, err)
	require.True(t, haveLeftRecursion)
	mapRules := make(map[string]*ast.Rule, len(grammar.Rules))
	for _, rule := range grammar.Rules {
		mapRules[rule.Name.Val] = rule
	}
	assert.False(t, mapRules["start"].Nullable)
	assert.True(t, mapRules["sign"].Nullable)
	assert.True(t, mapRules["start"].LeftRecursive)
	assert.False(t, mapRules["sign"].LeftRecursive)
}

func TestMutuallyLeftRrecursive(t *testing.T) {
	t.Parallel()
	text := `
	start = foo 'E'
    foo = bar 'A' / 'B'
    bar = foo 'C' / 'D'
	`
	p := bootstrap.NewParser()
	grammar, err := p.Parse("", strings.NewReader(text))
	require.NoError(t, err)
	haveLeftRecursion, err := builder.PrepareGrammar(grammar)
	require.NoError(t, err)
	require.True(t, haveLeftRecursion)
	mapRules := make(map[string]*ast.Rule, len(grammar.Rules))
	for _, rule := range grammar.Rules {
		mapRules[rule.Name.Val] = rule
	}
	assert.False(t, mapRules["start"].LeftRecursive)
	assert.True(t, mapRules["foo"].LeftRecursive)
	assert.True(t, mapRules["bar"].LeftRecursive)
}

func TestNastyMutuallyLeftRrecursive(t *testing.T) {
	t.Parallel()
	text := `
	start = target '='
    target = maybe '+' / NAME
    maybe = maybe '-' / target
	`
	p := bootstrap.NewParser()
	grammar, err := p.Parse("", strings.NewReader(text))
	require.NoError(t, err)
	haveLeftRecursion, err := builder.PrepareGrammar(grammar)
	require.NoError(t, err)
	require.True(t, haveLeftRecursion)
	mapRules := make(map[string]*ast.Rule, len(grammar.Rules))
	for _, rule := range grammar.Rules {
		mapRules[rule.Name.Val] = rule
	}
	assert.False(t, mapRules["start"].LeftRecursive)
	assert.True(t, mapRules["target"].LeftRecursive)
	assert.True(t, mapRules["maybe"].LeftRecursive)
}

func TestLeftRecursionTooComplex(t *testing.T) {
	t.Parallel()
	text := `
	start = foo
    foo = bar '+' / baz '+' / '+'
    bar = baz '-' / foo '-' / '-'
    baz = foo '*' / bar '*' / '*'
	`
	p := bootstrap.NewParser()
	grammar, err := p.Parse("", strings.NewReader(text))
	require.NoError(t, err)
	_, err = builder.PrepareGrammar(grammar)
	require.ErrorIs(t, err, builder.ErrNoLeader)
}
