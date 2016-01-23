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

func (o RetainOp) Type() OpType {
	return OpRetain
}

func (o RetainOp) Span() int {
	return int(o)
}

type InsertOp string

func (o InsertOp) Type() OpType {
	return OpInsert
}

func (o InsertOp) Span() int {
	return len(o)
}

type DeleteOp int

func (o DeleteOp) Type() OpType {
	return OpDelete
}

func (o DeleteOp) Span() int {
	return int(o)
}
