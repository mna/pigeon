package vm

import (
	"bytes"
	"errors"
)

var (
	// errInvalidEncoding is returned when the source is not properly
	// utf8-encoded.
	errInvalidEncoding = errors.New("invalid encoding")

	// errNoMatch is returned if no match could be found and no other
	// error has been raised.
	errNoMatch = errors.New("no match found")
)

// errList cumulates the errors found by the parser. It is part
// of the supported API.
type errList []error

// ϡadd adds err to the list of errors.
func (e *errList) ϡadd(err error) {
	*e = append(*e, err)
}

// ϡerr returns the error list as an error, or nil if the list is empty.
func (e errList) ϡerr() error {
	if len(e) == 0 {
		return nil
	}
	e.ϡdedupe()
	return e
}

// ϡdedupe removes duplicate error messages from the list.
func (e *errList) ϡdedupe() {
	var cleaned []error
	set := make(map[string]bool)
	for _, err := range *e {
		if msg := err.Error(); !set[msg] {
			set[msg] = true
			cleaned = append(cleaned, err)
		}
	}
	*e = cleaned
}

// Error returns the error message for the errList. It implements the
// error interface.
func (e errList) Error() string {
	var buf bytes.Buffer

	for i, err := range e {
		if i > 0 {
			buf.WriteRune('\n')
		}
		buf.WriteString(err.Error())
	}
	return buf.String()
}

// parserError wraps an error with a prefix indicating the rule in which
// the error occurred. The original error is stored in the Inner field.
// It is part of the supported API.
type parserError struct {
	Inner   error
	ϡpos    position
	ϡprefix string
}

// Error returns the prefixed error message. It implements the error
// interface.
func (p *parserError) Error() string {
	return p.ϡprefix + ": " + p.Inner.Error()
}
