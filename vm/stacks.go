package vm

type ϡpstack []ϡsvpt

func (p *ϡpstack) push(pt ϡsvpt) {
	*p = append(*p, pt)
}

func (p *ϡpstack) pop() ϡsvpt {
	n := len(*p)
	if n == 0 {
		panic("pstack is empty")
	}
	v := (*p)[n-1]
	*p = (*p)[:n-1]
	return v
}

type ϡistack []int

func (i *ϡistack) push(v int) {
	*i = append(*i, v)
}

func (i *ϡistack) pop() int {
	n := len(*i)
	if n == 0 {
		panic("istack is empty")
	}
	v := (*i)[n-1]
	*i = (*i)[:n-1]
	return v
}

type ϡvstack []interface{}

func (v *ϡvstack) push(i interface{}) {
	*v = append(*v, i)
}

func (v *ϡvstack) pop() interface{} {
	i := v.peek()
	*v = (*v)[:len(*v)-1]
	return i
}

func (v *ϡvstack) peek() interface{} {
	n := len(*v)
	if n == 0 {
		panic("vstack is empty")
	}
	i := (*v)[n-1]
	return i
}

type ϡlstack [][]int

func (l *ϡlstack) push(a []int) {
	*l = append(*l, a)
}

func (l *ϡlstack) pop() []int {
	n := len(*l)
	if n == 0 {
		panic("lstack is empty")
	}
	a := (*l)[n-1]
	*l = (*l)[:n-1]
	return a
}

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
