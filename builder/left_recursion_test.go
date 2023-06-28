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
		t.Error("Rule 'start' does not contain left recursion")
	}
	if !mapRules["expr"].LeftRecursive {
		t.Error("Rule 'expr' contains left recursion")
	}
	if mapRules["term"].LeftRecursive {
		t.Error("Rule 'term' does not contain left recursion")
	}
	if mapRules["foo"].LeftRecursive {
		t.Error("Rule 'foo' does not contain left recursion")
	}
	if mapRules["bar"].LeftRecursive {
		t.Error("Rule 'bar' does not contain left recursion")
	}
	if mapRules["baz"].LeftRecursive {
		t.Error("Rule 'baz' does not contain left recursion")
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
		t.Error("Rule 'start' is not nullable")
	}
	if !mapRules["sign"].Nullable {
		t.Error("Rule 'sign' is nullable")
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
		t.Error("Rule 'start' is not Nullable")
	}
	if !mapRules["sign"].Nullable {
		t.Error("Rule 'sign' is Nullable")
	}
	if !mapRules["start"].LeftRecursive {
		t.Error("Rule 'start' does not contain left recursion")
	}
	if mapRules["sign"].LeftRecursive {
		t.Error("Rule 'sign' contains left recursion")
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
		t.Error("Rule 'start' does not contain left recursion")
	}
	if !mapRules["foo"].LeftRecursive {
		t.Error("Rule 'foo' contains left recursion")
	}
	if !mapRules["bar"].LeftRecursive {
		t.Error("Rule 'bar' contains left recursion")
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
		t.Error("Rule 'start' does not contain left recursion")
	}
	if !mapRules["target"].LeftRecursive {
		t.Error("Rule 'target' contains left recursion")
	}
	if !mapRules["maybe"].LeftRecursive {
		t.Error("Rule 'maybe' contains left recursion")
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
