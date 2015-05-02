package vm

//+pigeon following code is part of the generated parser

// ϡpstack implements the Position stack. It stores savepoints.
type ϡpstack []ϡsvpt

// push adds a value on the stack.
func (p *ϡpstack) push(pt ϡsvpt) {
	*p = append(*p, pt)
}

// pop removes a value from the stack.
func (p *ϡpstack) pop() ϡsvpt {
	n := len(*p)
	if n == 0 {
		panic("pstack is empty")
	}
	v := (*p)[n-1]
	*p = (*p)[:n-1]
	return v
}

// ϡistack implements the Instruction index stack. It stores integers.
type ϡistack []int

// push adds a value on the stack.
func (i *ϡistack) push(v int) {
	*i = append(*i, v)
}

// pop removes a value from the stack.
func (i *ϡistack) pop() int {
	n := len(*i)
	if n == 0 {
		panic("istack is empty")
	}
	v := (*i)[n-1]
	*i = (*i)[:n-1]
	return v
}

// ϡvstack implements the Value stack. It stores empty interfaces.
type ϡvstack []interface{}

// push adds a value on the stack.
func (v *ϡvstack) push(i interface{}) {
	*v = append(*v, i)
}

// pop removes a value from the stack.
func (v *ϡvstack) pop() interface{} {
	i := v.peek()
	*v = (*v)[:len(*v)-1]
	return i
}

// peek returns the value at the top of the stack, leaving it there.
func (v *ϡvstack) peek() interface{} {
	n := len(*v)
	if n == 0 {
		panic("vstack is empty")
	}
	i := (*v)[n-1]
	return i
}

// ϡlstack implements the Loop stack. It stores slices of integers.
type ϡlstack [][]int

// push adds a value on the stack.
func (l *ϡlstack) push(a []int) {
	*l = append(*l, a)
}

// pop removes a value from the stack.
func (l *ϡlstack) pop() []int {
	n := len(*l)
	if n == 0 {
		panic("lstack is empty")
	}
	a := (*l)[n-1]
	*l = (*l)[:n-1]
	return a
}

// take removes the integer at index 0 from the slice at the top of the
// stack. It returns -1 if the slice is empty. The slice is left on the
// stack.
func (l *ϡlstack) take() int {
	n := len(*l)
	if n == 0 {
		panic("lstack is empty")
	}

	v := -1
	a := (*l)[n-1]
	if len(a) > 0 {
		v = a[0]
		(*l)[n-1] = a[1:]
	}
	return v
}

// ϡargsSet holds the list of arguments (key and value) to pass
// to the code blocks.
type ϡargsSet map[string]interface{}

// ϡastack is a stack of ϡargsSet.
type ϡastack []ϡargsSet

// push adds an empty ϡargsSet on top of the stack.
func (a *ϡastack) push() {
	*a = append(*a, ϡargsSet{})
}

// pop removes the top ϡargsSet from the stack.
func (a *ϡastack) pop() {
	n := len(*a)
	if n == 0 {
		panic("astack is empty")
	}
	*a = (*a)[:n-1]
}

// peek returns the current top ϡargsSet.
func (a *ϡastack) peek() ϡargsSet {
	n := len(*a)
	if n == 0 {
		panic("astack is empty")
	}
	as := (*a)[n-1]
	return as
}
