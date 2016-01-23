package ot

type OpType int

const (
	OpNoop OpType = iota
	OpRetain
	OpInsert
	OpDelete
)

type Op interface {
	Type() OpType
	Span() int
}

type RetainOp int

func (p RetainOp) Type() OpType {
	return OpRetain
}

func (p RetainOp) Span() int {
	return int(p)
}

func (p RetainOp) ComposeRetain(q RetainOp) (Op, Op, Op) {
	switch {
	case p > q:
		return q, p - q, Noop

	default:
		return p, Noop, q - p
	}
}

func (p RetainOp) ComposeDelete(q DeleteOp) (Op, Op, Op) {
	switch {
	case int(p) < int(q):
		return DeleteOp(p), Noop, q - DeleteOp(p)

	default:
		return q, p - RetainOp(q), Noop
	}
}

var Noop = RetainOp(0)

type InsertOp string

func (p InsertOp) Type() OpType {
	return OpInsert
}

func (p InsertOp) Span() int {
	return len(p)
}

func (p InsertOp) ComposeRetain(q RetainOp) (Op, Op, Op) {
	switch {
	case len(p) > int(q):
		return p[:int(q)], p[int(q):], Noop

	default:
		return p, Noop, q - RetainOp(len(p))
	}
}

func (p InsertOp) ComposeDelete(q DeleteOp) (Op, Op, Op) {
	switch {
	case len(p) > int(q):
		return Noop, p[int(q):], Noop

	case len(p) < int(q):
		return Noop, Noop, q - DeleteOp(len(p))

	default:
		return Noop, Noop, Noop
	}
}

type DeleteOp int

func (p DeleteOp) Type() OpType {
	return OpDelete
}

func (p DeleteOp) Span() int {
	return int(p)
}

type RetainComposer interface {
	ComposeRetain(q RetainOp) (Op, Op, Op)
}

type DeleteComposer interface {
	ComposeDelete(q DeleteOp) (Op, Op, Op)
}
