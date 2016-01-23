package hub

import (
	"bytes"
	"testing"

	"github.com/gophergala2016/papyrus/ot"
)

func TestDocumentApply(t *testing.T) {
	cases := []struct {
		doc  Document
		chgs []Change
		exp  ot.Blob
	}{
		{
			doc: Document{
				Blob: ot.Blob("Lem sum lor"),
			},
			chgs: []Change{
				{
					Root: 0,
					Ops:  ot.Ops{ot.RetainOp(1), ot.InsertOp("or"), ot.RetainOp(10)},
				},
				{
					Root: 0,
					Ops:  ot.Ops{ot.RetainOp(8), ot.InsertOp("do"), ot.RetainOp(3)},
				},
				{
					Root: 1,
					Ops:  ot.Ops{ot.RetainOp(6), ot.InsertOp("ip"), ot.RetainOp(7)},
				},
			},
			exp: ot.Blob("Lorem ipsum dolor"),
		},
	}
	for i, c := range cases {
		for _, ch := range c.chgs {
			_, err := c.doc.Apply(ch)
			if err != nil {
				t.Fatalf("%d: expected err == nil, got %#v", i, err)
			}
		}
		if !bytes.Equal(c.doc.Blob, c.exp) {
			t.Fatalf("%d: expected blob == %q, got %q", i, c.exp, c.doc.Blob)
		}
	}
}
