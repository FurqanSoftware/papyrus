package ot

import (
	"testing"
)

func TestRetainOpComposeRetain(t *testing.T) {
	cases := []struct {
		ret1 RetainOp
		ret2 RetainOp
		exp  [3]Op
	}{
		{
			ret1: 6,
			ret2: 6,
			exp:  [3]Op{RetainOp(6), Noop, Noop},
		},
		{
			ret1: 6,
			ret2: 3,
			exp:  [3]Op{RetainOp(3), RetainOp(3), Noop},
		},
		{
			ret1: 6,
			ret2: 9,
			exp:  [3]Op{RetainOp(6), Noop, RetainOp(3)},
		},
	}
	for i, c := range cases {
		x, p, q := c.ret1.ComposeRetain(c.ret2)
		if x != c.exp[0] {
			t.Fatalf("%d: expected x == %#v, got %#v", i, c.exp[0], x)
		}
		if p != c.exp[1] {
			t.Fatalf("%d: expected p == %#v, got %#v", i, c.exp[1], p)
		}
		if q != c.exp[2] {
			t.Fatalf("%d: expected q == %#v, got %#v", i, c.exp[2], q)
		}
	}
}

func TestInsertOpComposeRetain(t *testing.T) {
	cases := []struct {
		ins InsertOp
		ret RetainOp
		exp [3]Op
	}{
		{
			ins: "foobar",
			ret: 6,
			exp: [3]Op{InsertOp("foobar"), Noop, Noop},
		},
		{
			ins: "foobar",
			ret: 3,
			exp: [3]Op{InsertOp("foo"), InsertOp("bar"), Noop},
		},
		{
			ins: "foobar",
			ret: 9,
			exp: [3]Op{InsertOp("foobar"), Noop, RetainOp(3)},
		},
	}
	for i, c := range cases {
		x, p, q := c.ins.ComposeRetain(c.ret)
		if x != c.exp[0] {
			t.Fatalf("%d: expected x == %#v, got %#v", i, c.exp[0], x)
		}
		if p != c.exp[1] {
			t.Fatalf("%d: expected p == %#v, got %#v", i, c.exp[1], p)
		}
		if q != c.exp[2] {
			t.Fatalf("%d: expected q == %#v, got %#v", i, c.exp[2], q)
		}
	}
}

func TestRetainOpComposeDelete(t *testing.T) {
	cases := []struct {
		ret RetainOp
		del DeleteOp
		exp [3]Op
	}{
		{
			ret: 6,
			del: 6,
			exp: [3]Op{DeleteOp(6), Noop, Noop},
		},
		{
			ret: 6,
			del: 3,
			exp: [3]Op{DeleteOp(3), RetainOp(3), Noop},
		},
		{
			ret: 6,
			del: 9,
			exp: [3]Op{DeleteOp(6), Noop, DeleteOp(3)},
		},
	}
	for i, c := range cases {
		x, p, q := c.ret.ComposeDelete(c.del)
		if x != c.exp[0] {
			t.Fatalf("%d: expected x == %#v, got %#v", i, c.exp[0], x)
		}
		if p != c.exp[1] {
			t.Fatalf("%d: expected p == %#v, got %#v", i, c.exp[1], p)
		}
		if q != c.exp[2] {
			t.Fatalf("%d: expected q == %#v, got %#v", i, c.exp[2], q)
		}
	}
}

func TestInsertOpComposeDelete(t *testing.T) {
	cases := []struct {
		ins InsertOp
		del DeleteOp
		exp [3]Op
	}{
		{
			ins: "foobar",
			del: 6,
			exp: [3]Op{Noop, Noop, Noop},
		},
		{
			ins: "foobar",
			del: 3,
			exp: [3]Op{Noop, InsertOp("bar"), Noop},
		},
		{
			ins: "foobar",
			del: 9,
			exp: [3]Op{Noop, Noop, DeleteOp(3)},
		},
	}
	for i, c := range cases {
		x, p, q := c.ins.ComposeDelete(c.del)
		if x != c.exp[0] {
			t.Fatalf("%d: expected x == %#v, got %#v", i, c.exp[0], x)
		}
		if p != c.exp[1] {
			t.Fatalf("%d: expected p == %#v, got %#v", i, c.exp[1], p)
		}
		if q != c.exp[2] {
			t.Fatalf("%d: expected q == %#v, got %#v", i, c.exp[2], q)
		}
	}
}

func TestRetainOpTransformRetain(t *testing.T) {
	cases := []struct {
		ret1 RetainOp
		ret2 RetainOp
		exp  [4]Op
	}{
		{
			ret1: 6,
			ret2: 6,
			exp:  [4]Op{RetainOp(6), RetainOp(6), Noop, Noop},
		},
		{
			ret1: 6,
			ret2: 3,
			exp:  [4]Op{RetainOp(3), RetainOp(3), RetainOp(3), Noop},
		},
		{
			ret1: 6,
			ret2: 9,
			exp:  [4]Op{RetainOp(6), RetainOp(6), Noop, RetainOp(3)},
		},
	}
	for i, c := range cases {
		x, y, p, q := c.ret1.TransformRetain(c.ret2)
		if x != c.exp[0] {
			t.Fatalf("%d: expected x == %#v, got %#v", i, c.exp[0], x)
		}
		if y != c.exp[1] {
			t.Fatalf("%d: expected y == %#v, got %#v", i, c.exp[1], y)
		}
		if p != c.exp[2] {
			t.Fatalf("%d: expected p == %#v, got %#v", i, c.exp[1], p)
		}
		if q != c.exp[3] {
			t.Fatalf("%d: expected q == %#v, got %#v", i, c.exp[2], q)
		}
	}
}

func TestDeleteOpTransformRetain(t *testing.T) {
	cases := []struct {
		del DeleteOp
		ret RetainOp
		exp [4]Op
	}{
		{
			del: 6,
			ret: 6,
			exp: [4]Op{DeleteOp(6), Noop, Noop, Noop},
		},
		{
			del: 6,
			ret: 3,
			exp: [4]Op{DeleteOp(3), Noop, DeleteOp(3), Noop},
		},
		{
			del: 6,
			ret: 9,
			exp: [4]Op{DeleteOp(6), Noop, Noop, RetainOp(3)},
		},
	}
	for i, c := range cases {
		x, y, p, q := c.del.TransformRetain(c.ret)
		if x != c.exp[0] {
			t.Fatalf("%d: expected x == %#v, got %#v", i, c.exp[0], x)
		}
		if y != c.exp[1] {
			t.Fatalf("%d: expected y == %#v, got %#v", i, c.exp[1], y)
		}
		if p != c.exp[2] {
			t.Fatalf("%d: expected p == %#v, got %#v", i, c.exp[1], p)
		}
		if q != c.exp[3] {
			t.Fatalf("%d: expected q == %#v, got %#v", i, c.exp[2], q)
		}
	}
}
