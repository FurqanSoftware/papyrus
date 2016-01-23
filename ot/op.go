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
}

type RetainOp int

func (o RetainOp) Type() OpType {
	return OpRetain
}

type InsertOp string

func (o InsertOp) Type() OpType {
	return OpInsert
}

type DeleteOp int

func (o DeleteOp) Type() OpType {
	return OpDelete
}
