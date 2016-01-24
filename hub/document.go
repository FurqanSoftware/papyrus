package hub

import (
	"errors"

	"github.com/gophergala2016/papyrus/ot"
)

var ErrBadChange = errors.New("hub: invalid change")

type Document struct {
	ID      string
	Blob    ot.Blob
	History []Change
}

func (d *Document) Head() Change {
	return Change{
		ID:   "",
		Root: len(d.History),
		Ops:  ot.Ops{ot.InsertOp(d.Blob)},
	}
}

func (d *Document) Apply(ch Change) (Change, error) {
	if ch.Root < 0 || ch.Root > len(d.History) {
		return Change{}, ErrBadChange
	}
	ops := ch.Ops
	for i := ch.Root; i < len(d.History); i++ {
		opp, _, err := ops.Transform(d.History[i].Ops)
		if err != nil {
			return Change{}, err
		}
		ops = opp
	}
	err := d.Blob.Apply(ops)
	if err != nil {
		return Change{}, err
	}
	d.History = append(d.History, Change{
		Root: ch.Root,
		Ops:  ops,
	})
	return Change{
		ID:   ch.ID,
		Root: len(d.History),
		Ops:  ops,
	}, nil
}

type Change struct {
	ID   string
	Root int
	Ops  ot.Ops
}
