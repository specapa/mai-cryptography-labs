package rbtreealloc

type color uint8

const (
	red color = iota
	black
)

type node struct {
	offset, size        int
	c                   color
	left, right, parent *node
}

type Allocator struct {
	root *node
	heap []byte
}

func NewAllocator(heapSize int) *Allocator {
	n := &node{offset: 0, size: heapSize, c: black}
	return &Allocator{root: n, heap: make([]byte, heapSize)}
}

func (a *Allocator) Alloc(n int) (offset int, ok bool) {
	x := a.findFit(a.root, n)
	if x == nil {
		return 0, false
	}
	off := x.offset
	if x.size == n {
		a.delete(x)
	} else {
		x.offset += n
		x.size -= n
	}
	return off, true
}

func (a *Allocator) findFit(n *node, need int) *node {
	var best *node
	for n != nil {
		if n.size >= need {
			best = n
			n = n.left
		} else {
			n = n.right
		}
	}
	return best
}

func (a *Allocator) rotateLeft(x *node) {
	y := x.right
	x.right = y.left
	if y.left != nil {
		y.left.parent = x
	}
	y.parent = x.parent
	if x.parent == nil {
		a.root = y
	} else if x == x.parent.left {
		x.parent.left = y
	} else {
		x.parent.right = y
	}
	y.left = x
	x.parent = y
}

func (a *Allocator) rotateRight(x *node) {
	y := x.left
	x.left = y.right
	if y.right != nil {
		y.right.parent = x
	}
	y.parent = x.parent
	if x.parent == nil {
		a.root = y
	} else if x == x.parent.right {
		x.parent.right = y
	} else {
		x.parent.left = y
	}
	y.right = x
	x.parent = y
}

func (a *Allocator) insert(z *node) {
	var y *node
	x := a.root
	for x != nil {
		y = x
		if z.size < x.size {
			x = x.left
		} else {
			x = x.right
		}
	}
	z.parent = y
	if y == nil {
		a.root = z
	} else if z.size < y.size {
		y.left = z
	} else {
		y.right = z
	}
	z.c = red
	a.insertFix(z)
}

func (a *Allocator) insertFix(z *node) {
	for z.parent != nil && z.parent.c == red {
		if z.parent == z.parent.parent.left {
			y := z.parent.parent.right
			if y != nil && y.c == red {
				z.parent.c = black
				y.c = black
				z.parent.parent.c = red
				z = z.parent.parent
				continue
			}
			if z == z.parent.right {
				z = z.parent
				a.rotateLeft(z)
			}
			z.parent.c = black
			z.parent.parent.c = red
			a.rotateRight(z.parent.parent)
		} else {
			y := z.parent.parent.left
			if y != nil && y.c == red {
				z.parent.c = black
				y.c = black
				z.parent.parent.c = red
				z = z.parent.parent
				continue
			}
			if z == z.parent.left {
				z = z.parent
				a.rotateRight(z)
			}
			z.parent.c = black
			z.parent.parent.c = red
			a.rotateLeft(z.parent.parent)
		}
	}
	a.root.c = black
}

func (a *Allocator) transplant(u, v *node) {
	if u.parent == nil {
		a.root = v
	} else if u == u.parent.left {
		u.parent.left = v
	} else {
		u.parent.right = v
	}
	if v != nil {
		v.parent = u.parent
	}
}

func (a *Allocator) minimum(x *node) *node {
	for x.left != nil {
		x = x.left
	}
	return x
}

func (a *Allocator) delete(z *node) {
	y := z
	yOrig := y.c
	var x *node
	if z.left == nil {
		x = z.right
		a.transplant(z, z.right)
	} else if z.right == nil {
		x = z.left
		a.transplant(z, z.left)
	} else {
		y = a.minimum(z.right)
		yOrig = y.c
		x = y.right
		if y.parent == z {
			if x != nil {
				x.parent = y
			}
		} else {
			a.transplant(y, y.right)
			y.right = z.right
			y.right.parent = y
		}
		a.transplant(z, y)
		y.left = z.left
		y.left.parent = y
		y.c = z.c
	}
	if yOrig == black && x != nil {
		a.deleteFix(x)
	}
	if a.root != nil {
		a.root.c = black
	}
}

func (a *Allocator) deleteFix(x *node) {
	for x != a.root && x.c == black {
		if x == x.parent.left {
			w := x.parent.right
			if w != nil && w.c == red {
				w.c = black
				x.parent.c = red
				a.rotateLeft(x.parent)
				w = x.parent.right
			}
			if w == nil {
				x = x.parent
				continue
			}
			leftBlack := w.left == nil || w.left.c == black
			rightBlack := w.right == nil || w.right.c == black
			if leftBlack && rightBlack {
				w.c = red
				x = x.parent
				continue
			}
			if rightBlack {
				if w.left != nil {
					w.left.c = black
				}
				w.c = red
				a.rotateRight(w)
				w = x.parent.right
			}
			w.c = x.parent.c
			x.parent.c = black
			if w.right != nil {
				w.right.c = black
			}
			a.rotateLeft(x.parent)
			x = a.root
		} else {
			w := x.parent.left
			if w != nil && w.c == red {
				w.c = black
				x.parent.c = red
				a.rotateRight(x.parent)
				w = x.parent.left
			}
			if w == nil {
				x = x.parent
				continue
			}
			leftBlack := w.left == nil || w.left.c == black
			rightBlack := w.right == nil || w.right.c == black
			if leftBlack && rightBlack {
				w.c = red
				x = x.parent
				continue
			}
			if leftBlack {
				if w.right != nil {
					w.right.c = black
				}
				w.c = red
				a.rotateLeft(w)
				w = x.parent.left
			}
			w.c = x.parent.c
			x.parent.c = black
			if w.left != nil {
				w.left.c = black
			}
			a.rotateRight(x.parent)
			x = a.root
		}
	}
	x.c = black
}
