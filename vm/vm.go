package vm

// ϡmatchFailed is a sentinel value used to indicate a match failure.
var ϡmatchFailed = struct{}{}

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

// setOptions applies the options in sequence on the vm.
func (v *ϡvm) setOptions(opts []Option) {
	for _, opt := range opts {
		opt(v)
	}
}

// TODO : make run receive the list of instructions and the various lists,
// so it is easy to generate parsers on-the-fly for tests, without saving
// it to file.
func (v *ϡvm) run() (interface{}, error) {
	return nil, nil
}
