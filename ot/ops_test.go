package ot

import (
	"reflect"
	"testing"
)

func TestOpsCompose(t *testing.T) {
	cases := []struct {
		in1 Ops
		in2 Ops
		out Ops
	}{
		// Fbb -> Foobbaz-> Foobarbar
		{
			in1: Ops{RetainOp(1), InsertOp("oo"), RetainOp(2), InsertOp("az")},
			in2: Ops{RetainOp(4), InsertOp("ar"), RetainOp(3)},
			out: Ops{RetainOp(1), InsertOp("oo"), RetainOp(1), InsertOp("ar"), RetainOp(1), InsertOp("az")},
		},
	}
	for i, c := range cases {
		ops, err := c.in1.Compose(c.in2)
		if err != nil {
			t.Fatalf("%d: expected err == nil, got %#v", i, err)
		}
		if !reflect.DeepEqual(ops, c.out) {
			t.Fatalf("%d: expected ops == %#v, got %#v", i, c.out, ops)
		}
	}
}
