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
	grammar := `
	start = expr NEWLINE
    expr = ('-' term / expr '+' term / term)
    term = NUMBER
    foo = NAME+
    bar = NAME*
    baz = NAME?
	`
	p := bootstrap.NewParser()
	g, err := p.Parse("", strings.NewReader(grammar))
	require.NoError(t, err)
	err = builder.PrepareGramma(g)
	require.ErrorIs(t, err, builder.ErrHaveLeftRecirsion)
	mapRules := make(map[string]*ast.Rule, len(g.Rules))
	for _, rule := range g.Rules {
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
	grammar := `
	start = sign NUMBER
    sign = ('-' / '+')?
	`
	p := bootstrap.NewParser()
	g, err := p.Parse("", strings.NewReader(grammar))
	require.NoError(t, err)
	err = builder.PrepareGramma(g)
	require.NoError(t, err)
	mapRules := make(map[string]*ast.Rule, len(g.Rules))
	for _, rule := range g.Rules {
		mapRules[rule.Name.Val] = rule
	}
	assert.False(t, mapRules["start"].Nullable)
	assert.True(t, mapRules["sign"].Nullable)
}

func TestAdvancedLeftRrecursive(t *testing.T) {
	t.Parallel()
	grammar := `
	start = NUMBER / sign start
    sign = '-'?
	`
	p := bootstrap.NewParser()
	g, err := p.Parse("", strings.NewReader(grammar))
	require.NoError(t, err)
	err = builder.PrepareGramma(g)
	require.ErrorIs(t, err, builder.ErrHaveLeftRecirsion)
	mapRules := make(map[string]*ast.Rule, len(g.Rules))
	for _, rule := range g.Rules {
		mapRules[rule.Name.Val] = rule
	}
	assert.False(t, mapRules["start"].Nullable)
	assert.True(t, mapRules["sign"].Nullable)
	assert.True(t, mapRules["start"].LeftRecursive)
	assert.False(t, mapRules["sign"].LeftRecursive)
}

func TestMutuallyLeftRrecursive(t *testing.T) {
	t.Parallel()
	grammar := `
	start = foo 'E'
    foo = bar 'A' / 'B'
    bar = foo 'C' / 'D'
	`
	p := bootstrap.NewParser()
	g, err := p.Parse("", strings.NewReader(grammar))
	require.NoError(t, err)
	err = builder.PrepareGramma(g)
	require.ErrorIs(t, err, builder.ErrHaveLeftRecirsion)
	mapRules := make(map[string]*ast.Rule, len(g.Rules))
	for _, rule := range g.Rules {
		mapRules[rule.Name.Val] = rule
	}
	assert.False(t, mapRules["start"].LeftRecursive)
	assert.True(t, mapRules["foo"].LeftRecursive)
	assert.True(t, mapRules["bar"].LeftRecursive)
}

func TestNastyMutuallyLeftRrecursive(t *testing.T) {
	t.Parallel()
	grammar := `
	start = target '='
    target = maybe '+' / NAME
    maybe = maybe '-' / target
	`
	p := bootstrap.NewParser()
	g, err := p.Parse("", strings.NewReader(grammar))
	require.NoError(t, err)
	err = builder.PrepareGramma(g)
	require.ErrorIs(t, err, builder.ErrHaveLeftRecirsion)
	mapRules := make(map[string]*ast.Rule, len(g.Rules))
	for _, rule := range g.Rules {
		mapRules[rule.Name.Val] = rule
	}
	assert.False(t, mapRules["start"].LeftRecursive)
	assert.True(t, mapRules["target"].LeftRecursive)
	assert.True(t, mapRules["maybe"].LeftRecursive)
}

func TestLeftRecursionTooComplex(t *testing.T) {
	t.Parallel()
	grammar := `
	start = foo
    foo = bar '+' / baz '+' / '+'
    bar = baz '-' / foo '-' / '-'
    baz = foo '*' / bar '*' / '*'
	`
	p := bootstrap.NewParser()
	g, err := p.Parse("", strings.NewReader(grammar))
	require.NoError(t, err)
	err = builder.PrepareGramma(g)
	require.ErrorIs(t, err, builder.ErrNoLeader)
}
