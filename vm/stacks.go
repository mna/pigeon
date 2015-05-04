package vm

//+pigeon: stacks.go

// ϡpstack implements the Position stack. It stores savepoints.
type ϡpstack struct {
	ar []ϡsvpt
	sp int
}

// push adds a value on the stack.
func (p *ϡpstack) push(pt ϡsvpt) {
	if p.sp >= len(p.ar) {
		p.ar = append(p.ar, pt)
	} else {
		p.ar[p.sp] = pt
	}
	p.sp++
}

// pop removes a value from the stack.
func (p *ϡpstack) pop() ϡsvpt {
	p.sp--
	return p.ar[p.sp]
}

func (p *ϡpstack) len() int {
	return p.sp
}

func newPstack(cap int) *ϡpstack {
	return &ϡpstack{ar: make([]ϡsvpt, cap)}
}

// ϡistack implements the Instruction index stack. It stores integers.
type ϡistack struct {
	ar []int
	sp int
}

// push adds a value on the stack.
func (i *ϡistack) push(v int) {
	if i.sp >= len(i.ar) {
		i.ar = append(i.ar, v)
	} else {
		i.ar[i.sp] = v
	}
	i.sp++
}

// pop removes a value from the stack.
func (i *ϡistack) pop() int {
	i.sp--
	return i.ar[i.sp]
}

func (i *ϡistack) len() int {
	return i.sp
}

func newIstack(cap int) *ϡistack {
	return &ϡistack{ar: make([]int, cap)}
}

// ϡvstack implements the Value stack. It stores empty interfaces.
type ϡvstack struct {
	ar []interface{}
	sp int
}

// push adds a value on the stack.
func (v *ϡvstack) push(i interface{}) {
	if v.sp >= len(v.ar) {
		v.ar = append(v.ar, i)
	} else {
		v.ar[v.sp] = i
	}
	v.sp++
}

// pop removes a value from the stack.
func (v *ϡvstack) pop() interface{} {
	v.sp--
	return v.ar[v.sp]
}

// peek returns the value at the top of the stack, leaving it there.
func (v *ϡvstack) peek() interface{} {
	return v.ar[v.sp-1]
}

func (v *ϡvstack) len() int {
	return v.sp
}

func newVstack(cap int) *ϡvstack {
	return &ϡvstack{ar: make([]interface{}, cap)}
}

// ϡlstack implements the Loop stack. It stores slices of integers.
type ϡlstack struct {
	ar [][]int
	sp int
}

// push adds a value on the stack.
func (l *ϡlstack) push(a []int) {
	if l.sp >= len(l.ar) {
		l.ar = append(l.ar, a)
	} else {
		l.ar[l.sp] = a
	}
	l.sp++
}

// pop removes a value from the stack.
func (l *ϡlstack) pop() []int {
	l.sp--
	return l.ar[l.sp]
}

// take removes the integer at index 0 from the slice at the top of the
// stack. It returns -1 if the slice is empty. The slice is left on the
// stack.
func (l *ϡlstack) take() int {
	v := -1
	a := l.ar[l.sp-1]
	if len(a) > 0 {
		v = a[0]
		l.ar[l.sp-1] = a[1:]
	}
	return v
}

func (l *ϡlstack) len() int {
	return l.sp
}

func newLstack(cap int) *ϡlstack {
	return &ϡlstack{ar: make([][]int, cap)}
}

// ϡargsSet holds the list of arguments (key and value) to pass
// to the code blocks.
type ϡargsSet map[string]interface{}

// ϡastack is a stack of ϡargsSet.
type ϡastack struct {
	ar []ϡargsSet
	sp int
}

// push adds an empty ϡargsSet on top of the stack.
func (a *ϡastack) push() {
	if a.sp >= len(a.ar) {
		a.ar = append(a.ar, nil)
	} else {
		a.ar[a.sp] = nil
	}
	a.sp++
}

// pop removes the top ϡargsSet from the stack.
func (a *ϡastack) pop() {
	a.sp--
}

// peek returns the current top ϡargsSet.
func (a *ϡastack) peek() ϡargsSet {
	as := a.ar[a.sp-1]
	if as == nil {
		as = make(ϡargsSet)
		a.ar[a.sp-1] = as
	}
	return as
}

func (a *ϡastack) len() int {
	return a.sp
}

func newAstack(cap int) *ϡastack {
	return &ϡastack{ar: make([]ϡargsSet, cap)}
}
