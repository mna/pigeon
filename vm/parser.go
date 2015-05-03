package vm

import (
	"fmt"
	"unicode/utf8"
)

//+pigeon: parser.go

// position records a position in the text. It is part of the supported
// API.
type position struct {
	// line is the 1-based index of the line of the current rune.
	line int
	// col is the 1-based index of the current rune on the line.
	col int
	// offset is the 0-based index of the starting byte of the current rune.
	offset int
}

// String formats a position as a string.
func (p position) String() string {
	return fmt.Sprintf("%d:%d (%d)", p.line, p.col, p.offset)
}

// current represents current matching data. It is the value on which
// action and predicate code blocks are generated as methods. It is
// part of the supported API.
type current struct {
	// pos holds the start position of the current match.
	pos position
	// text contains the raw text of the match. It is a slice in the
	// source data, so it should not be modified.
	text []byte
}

// ϡsvpt stores all state required to go back to a point in the
// parser.
type ϡsvpt struct {
	position
	rn rune
	w  int
}

// ϡparser parses the input text as rune code points.
type ϡparser struct {
	data []byte
	pt   ϡsvpt
	cur  current
}

// peek returns the current savepoint information.
func (p *ϡparser) peek() ϡsvpt {
	return p.pt
}

// read advances the parser to the next rune.
func (p *ϡparser) read() {
	p.pt.offset += p.pt.w
	rn, n := utf8.DecodeRune(p.data[p.pt.offset:])
	p.pt.rn = rn
	p.pt.w = n

	if rn == utf8.RuneError {
		if n > 0 {
			panic(errInvalidEncoding)
		}
	} else {
		p.pt.col++
		if rn == '\n' {
			p.pt.line++
			p.pt.col = 0
		}
	}
}

// sliceFrom gets the slice of bytes from the start savepoint to
// the current position, non inclusive.
func (p *ϡparser) sliceFrom(start ϡsvpt) []byte {
	return p.data[start.position.offset:p.pt.position.offset]
}
