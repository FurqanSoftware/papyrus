package ot

type OpType int

const (
	OpRetain OpType = iota
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

func (p RetainOp) TransformRetain(q RetainOp) (Op, Op, Op, Op) {
	switch {
	case p > q:
		return q, q, p - q, Noop

	default:
		return p, p, Noop, q - p
	}
}

func (p RetainOp) TransformDelete(q DeleteOp) (Op, Op, Op, Op) {
	switch {
	case int(p) > int(q):
		return Noop, q, p - RetainOp(q), Noop

	default:
		return Noop, p, Noop, RetainOp(q) - p
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

func (p DeleteOp) TransformRetain(q RetainOp) (Op, Op, Op, Op) {
	switch {
	case int(p) > int(q):
		return DeleteOp(q), Noop, p - DeleteOp(q), Noop

	default:
		return p, Noop, Noop, q - RetainOp(p)
	}
}

func (p DeleteOp) TransformDelete(q DeleteOp) (Op, Op, Op, Op) {
	switch {
	case p > q:
		return Noop, Noop, p - q, Noop

	default:
		return Noop, Noop, Noop, q - p
	}
}

type RetainComposer interface {
	ComposeRetain(q RetainOp) (Op, Op, Op)
}

type DeleteComposer interface {
	ComposeDelete(q DeleteOp) (Op, Op, Op)
}

type RetainTransformer interface {
	TransformRetain(q RetainOp) (Op, Op, Op, Op)
}

type DeleteTransformer interface {
	TransformDelete(q DeleteOp) (Op, Op, Op, Op)
}
