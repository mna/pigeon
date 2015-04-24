package vm

import (
	"unicode"
	"unicode/utf8"
)

type ϡpeekReader interface {
	peek() ϡsvpt
	read()
}

type ϡmatcher interface {
	match(ϡpeekReader) bool
}

type ϡanyMatcher struct{}

func (a ϡanyMatcher) match(pr ϡpeekReader) bool {
	pt := pr.peek()
	pr.read()
	return pt.rn != utf8.RuneError
}

type ϡstringMatcher struct {
	ignoreCase bool
	value      string // value must be lowercase if ignoreCase is true
}

func (s ϡstringMatcher) match(pr ϡpeekReader) bool {
	for _, want := range s.value {
		pt := pr.peek()
		if s.ignoreCase {
			pt.rn = unicode.ToLower(pt.rn)
		}
		if pt.rn != want {
			return false
		}
		pr.read()
	}
	return true
}

type ϡcharClassMatcher struct {
	chars   []rune
	ranges  []rune
	classes []*unicode.RangeTable

	ignoreCase bool
	inverted   bool
}

func (c ϡcharClassMatcher) match(pr ϡpeekReader) bool {
	return false
}
