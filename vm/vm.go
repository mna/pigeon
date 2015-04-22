package vm

// ϡsentinel is a type used to define sentinel values that shouldn't
// be equal to something else.
type ϡsentinel int

// ϡmatchFailed is a sentinel value used to indicate a match failure.
const ϡmatchFailed ϡsentinel = iota

type ϡmemoizedResult struct {
	v  interface{}
	pt ϡsvpt
}

type ϡvm struct {
	// input
	filename string
	parser   *ϡparser

	// options
	debug   bool
	memoize bool
	recover bool

	// error list
	errs *errList
}

// setOptions applies the options in sequence on the vm. It returns the
// vm to allow for chaining calls.
func (v *ϡvm) setOptions(opts []Option) *ϡvm {
	for _, opt := range opts {
		opt(v)
	}
	return v
}

// TODO : make run receive the list of instructions and the various lists,
// so it is easy to generate parsers on-the-fly for tests, without saving
// it to file.
func (v *ϡvm) run() (interface{}, error) {
	return nil, nil
}
