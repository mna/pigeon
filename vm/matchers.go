package vm

import (
	"fmt"
	"strconv"
	"unicode"
	"unicode/utf8"
)

//+pigeon: matchers.go

// ϡpeekReader is the interface that defines the peek and read
// methods.
type ϡpeekReader interface {
	peek() rune
	read()
}

// ϡmatcher is the interface that defines the match method.
type ϡmatcher interface {
	match(ϡpeekReader) bool
	toDisplayMsg() string
}

// ϡanyMatcher is a matcher that matches any character but the
// EOF.
type ϡanyMatcher struct{}

// match tries to match a character in the peekReader.
func (a ϡanyMatcher) match(pr ϡpeekReader) bool {
	rn := pr.peek()
	pr.read()
	return rn != utf8.RuneError
}

func (a ϡanyMatcher) toDisplayMsg() string {
	return "<any>"
}

func (a ϡanyMatcher) String() string {
	return "."
}

// ϡstringMatcher is a matcher that matches a string.
type ϡstringMatcher struct {
	ignoreCase bool
	value      string // value must be lowercase if ignoreCase is true
}

// match tries to match the string in the peekReader.
func (s ϡstringMatcher) match(pr ϡpeekReader) bool {
	for _, want := range s.value {
		rn := pr.peek()
		pr.read()
		if s.ignoreCase {
			rn = unicode.ToLower(rn)
		}
		if rn != want {
			return false
		}
	}
	return true
}

func (s ϡstringMatcher) toDisplayMsg() string {
	return s.String()
}

func (s ϡstringMatcher) String() string {
	v := strconv.Quote(s.value)
	if s.ignoreCase {
		v += "i"
	}
	return v
}

// ϡcharClassMatcher is a matcher that matches classes of characters.
type ϡcharClassMatcher struct {
	raw     string
	chars   []rune // runes must be lowercase if ignoreCase is true
	ranges  []rune // same for ranges
	classes []*unicode.RangeTable

	ignoreCase bool
	inverted   bool
}

func (c ϡcharClassMatcher) toDisplayMsg() string {
	return c.raw
}

func (c ϡcharClassMatcher) String() string {
	return c.raw
}

// match tries to match classes of characters in the peekReader.
func (c ϡcharClassMatcher) match(pr ϡpeekReader) bool {
	rn := pr.peek()
	pr.read()

	if c.ignoreCase {
		rn = unicode.ToLower(rn)
	}

	// try to match in the list of available chars
	for _, ch := range c.chars {
		if rn == ch {
			return !c.inverted
		}
	}

	// try to match in the list of ranges
	for i := 0; i < len(c.ranges); i += 2 {
		if rn >= c.ranges[i] && rn <= c.ranges[i+1] {
			return !c.inverted
		}
	}

	// try to match in the list of Unicode classes
	for _, cl := range c.classes {
		if unicode.Is(cl, rn) {
			return !c.inverted
		}
	}

	return c.inverted
}

// ϡrangeTable returns the corresponding unicode range table from the
// provided class name.
func ϡrangeTable(class string) *unicode.RangeTable {
	if rt, ok := unicode.Categories[class]; ok {
		return rt
	}
	if rt, ok := unicode.Properties[class]; ok {
		return rt
	}
	if rt, ok := unicode.Scripts[class]; ok {
		return rt
	}

	// cannot happen
	panic(fmt.Sprintf("invalid Unicode class: %s", class))
}
