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

func (v *ϡvm) setOptions(opts []Option) {
	for _, opt := range opts {
		opt(v)
	}
}

func (v *ϡvm) run() (interface{}, error) {
	return nil, nil
}
