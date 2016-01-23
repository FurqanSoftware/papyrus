package hub

import (
	"errors"

	"github.com/gophergala2016/papyrus/ot"
)

var ErrBadChange = errors.New("hub: invalid change")

type Document struct {
	Blob    ot.Blob
	History []Change
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
	ch = Change{
		Root: len(d.History),
		Ops:  ops,
	}
	d.History = append(d.History, ch)
	return ch, nil
}

type Change struct {
	Root int
	Ops  ot.Ops
}
