package ot

import "errors"

var (
	ErrBadSpans  = errors.New("ot: base span doesn't match target span")
	ErrTooShort  = errors.New("ot: operations are too short")
	ErrTooLong   = errors.New("ot: operations are too long")
	ErrBadOpPair = errors.New("ot: invalid pair of operations")
)

type Ops []Op

func (u Ops) SpanBase() int {
	l := 0
	for _, o := range u {
		if o.Type() == OpInsert {
			continue
		}
		l += o.Span()
	}
	return l
}

func (u Ops) SpanTarget() int {
	l := 0
	for _, o := range u {
		if o.Type() == OpDelete {
			continue
		}
		l += o.Span()
	}
	return l
}

func (u Ops) Compact() Ops {
	// return u

	z := Ops{}
	for i := 0; i < len(u); i++ {
		t := u[i]
	s:
		switch a := t.(type) {
		case RetainOp:
			for j := i + 1; j < len(u); j++ {
				switch b := u[j].(type) {
				case RetainOp:
					a += b
					t = a
					i++

				default:
					break s
				}
			}

		case InsertOp:
			for j := i + 1; j < len(u); j++ {
				switch b := u[j].(type) {
				case InsertOp:
					a += b
					t = a
					i++

				default:
					break s
				}
			}

		case DeleteOp:
			for j := i + 1; j < len(u); j++ {
				switch b := u[j].(type) {
				case DeleteOp:
					a += b
					t = a
					i++

				default:
					break s
				}
			}
		}
		if t.IsZero() {
			continue
		}
		z = append(z, t)
	}
	return z
}

func (u Ops) Compose(v Ops) (z Ops, err error) {
	if u.SpanTarget() != v.SpanBase() {
		return nil, ErrBadSpans
	}

	l := Cursor{Ops: u}
	r := Cursor{Ops: v}
	for !l.Fin() || !r.Fin() {
		a := l.Op()
		b := r.Op()

		switch {
		case a.Type() == OpDelete:
			z = append(z, a)
			l.Next(Noop)
			continue

		case b.Type() == OpInsert:
			z = append(z, b)
			r.Next(Noop)
			continue
		}

		if l.Fin() {
			return nil, ErrTooShort
		}
		if r.Fin() {
			return nil, ErrTooLong
		}

		var c, p, q Op
		switch b := b.(type) {
		case RetainOp:
			x, ok := a.(RetainComposer)
			if !ok {
				return nil, ErrBadOpPair
			}
			c, p, q = x.ComposeRetain(b)

		case DeleteOp:
			x, ok := a.(DeleteComposer)
			if !ok {
				return nil, ErrBadOpPair
			}
			c, p, q = x.ComposeDelete(b)
		}

		z = append(z, c)
		l.Next(p)
		r.Next(q)
	}

	return z.Compact(), nil
}

func (u Ops) Transform(v Ops) (up, vp Ops, err error) {
	if u.SpanBase() != v.SpanBase() {
		return nil, nil, ErrBadSpans
	}

	l := Cursor{Ops: u}
	r := Cursor{Ops: v}
	for !l.Fin() || !r.Fin() {
		a := l.Op()
		b := r.Op()

		switch {
		case a.Type() == OpInsert:
			up = append(up, a)
			vp = append(vp, RetainOp(a.Span()))
			l.Next(Noop)
			continue

		case b.Type() == OpInsert:
			up = append(up, RetainOp(b.Span()))
			vp = append(vp, b)
			r.Next(Noop)
			continue
		}

		if l.Fin() {
			return nil, nil, ErrTooShort
		}
		if r.Fin() {
			return nil, nil, ErrTooLong
		}

		var c, d, p, q Op
		switch b := b.(type) {
		case RetainOp:
			x, ok := a.(RetainTransformer)
			if !ok {
				return nil, nil, ErrBadOpPair
			}
			c, d, p, q = x.TransformRetain(b)

		case DeleteOp:
			x, ok := a.(DeleteTransformer)
			if !ok {
				return nil, nil, ErrBadOpPair
			}
			c, d, p, q = x.TransformDelete(b)
		}

		if c != Noop {
			up = append(up, c)
		}
		if d != Noop {
			vp = append(vp, d)
		}
		l.Next(p)
		r.Next(q)
	}

	return up.Compact(), vp.Compact(), nil
}

type Cursor struct {
	Ops
	nxt Op
	i   int
}

func (c *Cursor) Next(p Op) {
	c.nxt = p
	if c.nxt == Noop {
		c.i++
	}
}

func (c *Cursor) Op() Op {
	if c.i >= len(c.Ops) {
		return Noop
	}
	if c.nxt != Noop && c.nxt != nil {
		return c.nxt
	}
	return c.Ops[c.i]
}

func (c *Cursor) Fin() bool {
	return c.i >= len(c.Ops)
}
