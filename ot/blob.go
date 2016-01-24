package ot

import (
	"bytes"
	"errors"
	"fmt"
)

var (
	ErrBadSpan = errors.New("ot: base span doesn't match blob length")
)

type Blob []byte

func (b *Blob) Apply(u Ops) error {
	if u.SpanBase() != len(*b) {
		fmt.Printf("%#v %q\n", u, *b)
		return ErrBadSpan
	}

	z := bytes.Buffer{}

	d := *b
	i := 0
	for _, o := range u {
		switch o := o.(type) {
		case RetainOp:
			z.Write(d[i : i+int(o)])
			i += int(o)

		case InsertOp:
			z.Write([]byte(o))

		case DeleteOp:
			i += int(o)
		}
	}
	if i != len(d) {
		return ErrTooShort
	}

	*b = z.Bytes()
	return nil
}
