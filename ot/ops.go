package ot

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
