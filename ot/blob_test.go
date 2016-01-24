package ot

import (
	"bytes"
	"testing"
)

func TestBlobApply(t *testing.T) {
	cases := []struct {
		blob Blob
		ops  Ops
		exp  Blob
	}{
		{
			blob: Blob("Fbb"),
			ops:  Ops{RetainOp(1), InsertOp("oo"), RetainOp(1), InsertOp("ar"), RetainOp(1), InsertOp("az")},
			exp:  Blob("Foobarbaz"),
		},
		{
			blob: Blob("Fuubaz"),
			ops:  Ops{RetainOp(1), DeleteOp(2), InsertOp("oo"), InsertOp("bar"), RetainOp(3)},
			exp:  Blob("Foobarbaz"),
		},
	}
	for i, c := range cases {
		err := c.blob.Apply(c.ops)
		if err != nil {
			t.Fatalf("%d: expected err == nil, got %#v", i, err)
		}
		if !bytes.Equal(c.blob, c.exp) {
			t.Fatalf("%d: expected blob == %q, got %q", i, c.exp, c.blob)
		}
	}
}
