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
		if o.Type() == OpInsert || o.Type() == OpDelete {
			continue
		}
		l += o.Span()
	}
	return l
}

func (u Ops) SpanTarget() int {
	l := 0
	for _, o := range u {
		l += o.Span()
	}
	return l
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

		switch b := b.(type) {
		case RetainOp:
			x, ok := a.(RetainComposer)
			if !ok {
				return nil, ErrBadOpPair
			}
			c, p, q := x.ComposeRetain(b)
			z = append(z, c)
			l.Next(p)
			r.Next(q)

		case DeleteOp:
			x, ok := a.(DeleteComposer)
			if !ok {
				return nil, ErrBadOpPair
			}
			c, p, q := x.ComposeDelete(b)
			z = append(z, c)
			l.Next(p)
			r.Next(q)
		}
	}

	return
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
