package builder_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/mna/pigeon/ast"
	"github.com/mna/pigeon/bootstrap"
	"github.com/mna/pigeon/builder"
)

func TestLeftRecursive(t *testing.T) {
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
	if err != nil {
		t.Fatal(err)
	}
	haveLeftRecursion, err := builder.PrepareGrammar(grammar)
	if err != nil {
		t.Fatal(err)
	}
	if !haveLeftRecursion {
		t.Fatalf("Recursion not found")
	}

	mapRules := make(map[string]*ast.Rule, len(grammar.Rules))
	for _, rule := range grammar.Rules {
		mapRules[rule.Name.Val] = rule
	}
	if mapRules["start"].LeftRecursive {
		t.Fail()
	}
	if !mapRules["expr"].LeftRecursive {
		t.Fail()
	}
	if mapRules["term"].LeftRecursive {
		t.Fail()
	}
	if mapRules["foo"].LeftRecursive {
		t.Fail()
	}
	if mapRules["bar"].LeftRecursive {
		t.Fail()
	}
	if mapRules["baz"].LeftRecursive {
		t.Fail()
	}
}

func TestNullable(t *testing.T) {
	t.Parallel()

	text := `
	start = sign NUMBER
    sign = ('-' / '+')?
	`
	p := bootstrap.NewParser()
	grammar, err := p.Parse("", strings.NewReader(text))
	if err != nil {
		t.Fatal(err)
	}
	haveLeftRecursion, err := builder.PrepareGrammar(grammar)
	if err != nil {
		t.Fatal(err)
	}
	if haveLeftRecursion {
		t.Fatalf("Recursion found")
	}
	mapRules := make(map[string]*ast.Rule, len(grammar.Rules))
	for _, rule := range grammar.Rules {
		mapRules[rule.Name.Val] = rule
	}
	if mapRules["start"].Nullable {
		t.Fail()
	}
	if !mapRules["sign"].Nullable {
		t.Fail()
	}
}

func TestAdvancedLeftRecursive(t *testing.T) {
	t.Parallel()

	text := `
	start = NUMBER / sign start
    sign = '-'?
	`
	p := bootstrap.NewParser()
	grammar, err := p.Parse("", strings.NewReader(text))
	if err != nil {
		t.Fatal(err)
	}
	haveLeftRecursion, err := builder.PrepareGrammar(grammar)
	if err != nil {
		t.Fatal(err)
	}
	if !haveLeftRecursion {
		t.Fatalf("Recursion not found")
	}
	mapRules := make(map[string]*ast.Rule, len(grammar.Rules))
	for _, rule := range grammar.Rules {
		mapRules[rule.Name.Val] = rule
	}
	if mapRules["start"].Nullable {
		t.Fail()
	}
	if !mapRules["sign"].Nullable {
		t.Fail()
	}
	if !mapRules["start"].LeftRecursive {
		t.Fail()
	}
	if mapRules["sign"].LeftRecursive {
		t.Fail()
	}
}

func TestMutuallyLeftRecursive(t *testing.T) {
	t.Parallel()

	text := `
	start = foo 'E'
    foo = bar 'A' / 'B'
    bar = foo 'C' / 'D'
	`
	p := bootstrap.NewParser()
	grammar, err := p.Parse("", strings.NewReader(text))
	if err != nil {
		t.Fatal(err)
	}
	haveLeftRecursion, err := builder.PrepareGrammar(grammar)
	if err != nil {
		t.Fatal(err)
	}
	if !haveLeftRecursion {
		t.Fatalf("Recursion not found")
	}
	mapRules := make(map[string]*ast.Rule, len(grammar.Rules))
	for _, rule := range grammar.Rules {
		mapRules[rule.Name.Val] = rule
	}
	if mapRules["start"].LeftRecursive {
		t.Fail()
	}
	if !mapRules["foo"].LeftRecursive {
		t.Fail()
	}
	if !mapRules["bar"].LeftRecursive {
		t.Fail()
	}
}

func TestNastyMutuallyLeftRecursive(t *testing.T) {
	t.Parallel()

	text := `
	start = target '='
    target = maybe '+' / NAME
    maybe = maybe '-' / target
	`
	p := bootstrap.NewParser()
	grammar, err := p.Parse("", strings.NewReader(text))
	if err != nil {
		t.Fatal(err)
	}
	haveLeftRecursion, err := builder.PrepareGrammar(grammar)
	if err != nil {
		t.Fatal(err)
	}
	if !haveLeftRecursion {
		t.Fatalf("Recursion not found")
	}
	mapRules := make(map[string]*ast.Rule, len(grammar.Rules))
	for _, rule := range grammar.Rules {
		mapRules[rule.Name.Val] = rule
	}
	if mapRules["start"].LeftRecursive {
		t.Fail()
	}
	if !mapRules["target"].LeftRecursive {
		t.Fail()
	}
	if !mapRules["maybe"].LeftRecursive {
		t.Fail()
	}
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
	if err != nil {
		t.Fatal(err)
	}
	_, err = builder.PrepareGrammar(grammar)
	if !errors.Is(err, builder.ErrNoLeader) {
		t.Fatalf("Got %s, but expected %s", err, builder.ErrNoLeader)
	}
}
