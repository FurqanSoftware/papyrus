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

func TestOpsTransform(t *testing.T) {
	cases := []struct {
		in1  Ops
		in2  Ops
		out1 Ops
		out2 Ops
	}{
		// Fbb -> (Foobbaz, Fbarb) -> Foobarbar
		{
			in1:  Ops{RetainOp(1), InsertOp("oo"), RetainOp(2), InsertOp("az")},
			in2:  Ops{RetainOp(2), InsertOp("ar"), RetainOp(1)},
			out1: Ops{RetainOp(1), InsertOp("oo"), RetainOp(1), RetainOp(2), RetainOp(1), InsertOp("az")},
			out2: Ops{RetainOp(1), RetainOp(2), RetainOp(1), InsertOp("ar"), RetainOp(1), RetainOp(2)},
		},
	}
	for i, c := range cases {
		ops1, ops2, err := c.in1.Transform(c.in2)
		if err != nil {
			t.Fatalf("%d: expected err == nil, got %#v", i, err)
		}
		if !reflect.DeepEqual(ops1, c.out1) {
			t.Fatalf("%d: expected ops1 == %#v, got %#v", i, c.out1, ops1)
		}
		if !reflect.DeepEqual(ops2, c.out2) {
			t.Fatalf("%d: expected ops2 == %#v, got %#v", i, c.out2, ops2)
		}
	}
}
