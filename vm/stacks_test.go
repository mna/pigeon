package vm

import "testing"

func TestPStack(t *testing.T) {
	p := ϡpstack{}
	p.push(ϡsvpt{rn: 'a'})
	p.push(ϡsvpt{rn: 'b'})

	s := p.pop()
	if s.rn != 'b' {
		t.Errorf("want %c, got %c", 'b', s.rn)
	}
	s = p.pop()
	if s.rn != 'a' {
		t.Errorf("want %c, got %c", 'a', s.rn)
	}

	p.push(ϡsvpt{rn: 'c'})
	s = p.pop()
	if s.rn != 'c' {
		t.Errorf("want %c, got %c", 'c', s.rn)
	}

	ok := false
	func() {
		defer func() {
			if e := recover(); e != nil {
				ok = true
			}
		}()
		p.pop()
	}()
	if !ok {
		t.Errorf("want panic, got none")
	}
}

func TestIStack(t *testing.T) {
	i := ϡistack{}
	i.push(1)
	i.push(2)

	v := i.pop()
	if v != 2 {
		t.Errorf("want %d, got %d", 2, v)
	}
	v = i.pop()
	if v != 1 {
		t.Errorf("want %d, got %d", 1, v)
	}

	i.push(3)
	v = i.pop()
	if v != 3 {
		t.Errorf("want %d, got %d", 3, v)
	}

	ok := false
	func() {
		defer func() {
			if e := recover(); e != nil {
				ok = true
			}
		}()
		i.pop()
	}()
	if !ok {
		t.Errorf("want panic, got none")
	}
}

func TestVStack(t *testing.T) {
	v := ϡvstack{}
	v.push(1)
	v.push(2)

	vv := v.pop()
	if vv != 2 {
		t.Errorf("want %d, got %d", 2, vv)
	}
	vv = v.peek()
	if vv != 1 {
		t.Errorf("want %d, got %d", 1, vv)
	}
	vv = v.pop()
	if vv != 1 {
		t.Errorf("want %d, got %d", 1, vv)
	}

	v.push(3)
	vv = v.pop()
	if vv != 3 {
		t.Errorf("want %d, got %d", 3, vv)
	}

	ok := false
	func() {
		defer func() {
			if e := recover(); e != nil {
				ok = true
			}
		}()
		v.pop()
	}()
	if !ok {
		t.Errorf("want panic, got none")
	}
}

func TestLStack(t *testing.T) {
	l := ϡlstack{}
	l.push([]int{4})
	l.push([]int{2, 1})

	i := l.take()
	if i != 2 {
		t.Errorf("want %d, got %d", 2, i)
	}
	i = l.take()
	if i != 1 {
		t.Errorf("want %d, got %d", 1, i)
	}
	i = l.take()
	if i != -1 {
		t.Errorf("want %d, got %d", -1, i)
	}

	a := l.pop()
	if len(a) != 0 {
		t.Errorf("want empty array, got %v", a)
	}

	i = l.take()
	if i != 4 {
		t.Errorf("want %d, got %d", 4, i)
	}
	i = l.take()
	if i != -1 {
		t.Errorf("want %d, got %d", -1, i)
	}

	a = l.pop()
	if len(a) != 0 {
		t.Errorf("want empty array, got %v", a)
	}

	l.push([]int{3})
	a = l.pop()
	if len(a) != 1 {
		t.Errorf("want array of 1 element, got %v", a)
	} else if a[0] != 3 {
		t.Errorf("want %d, got %d", 3, a[0])
	}

	ok := false
	func() {
		defer func() {
			if e := recover(); e != nil {
				ok = true
			}
		}()
		l.pop()
	}()
	if !ok {
		t.Errorf("want panic, got none")
	}
}
