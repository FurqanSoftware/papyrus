var errBadType = new Error('ot.js: invalid operation type')
var errBadSpans = new Error('ot.js: base span doesn\'t match target span')
var errTooShort = new Error('ot.js: operations are too short')
var errTooLong = new Error('ot.js: operations are too long')
var errBadOpPair = new Error('ot.js: invalid pair of operations')
var errBadSpan = new Error("ot.js: base span doesn't match blob length")

var noop = 0

function opType(p) {
	switch(true) {
	case typeof p === 'string':
		return 'insert'

	case typeof p === 'number' && p >= 0:
		return 'retain'

	case typeof p === 'number' && p < 0:
		return 'delete'

	default:
		throw new errBadType
	}
}

function opSpan(p) {
	switch(opType(p)) {
	case 'insert':
		return p.length

	case 'retain':
		return p

	case 'delete':
		return -p

	default:
		throw new errBadType
	}
}

function opComposeRetainRetain(p, q) {
	switch(true) {
	case p > q:
		return [q, p - q, noop]

	default:
		return [p, noop, q - p]
	}
}

function opComposeRetainDelete(p, q) {
	switch(true) {
	case p < -q:
		return [-p, noop, p - q]

	default:
		return [q, p + q, noop]
	}
}

function opTransformRetainRetain(p, q) {
	switch(true) {
	case p > q:
		return [q, q, p - q, noop]

	default:
		return [p, p, noop, q - p]
	}
}

function opTransformRetainDelete(p, q) {
	switch(true) {
	case p > -q:
		return [noop, -q, p + q, noop]

	case p < -q:
		return [noop, -p, noop, q + p]

	default:
		return [noop, -p, noop, noop]
	}
}

function opComposeInsertRetain(p, q) {
	switch(true) {
	case p.length > q:
		return [p.substr(0, q), p.substr(q), noop]

	default:
		return [p, noop, q - p.length]
	}
}

function opComposeInsertDelete(p, q) {
	switch(true) {
	case p.length > -q:
		return [noop, p.substr(-q), noop]

	case p.length < -q:
		return [noop, noop, p.length + q]

	default:
		return [noop, noop, noop]
	}
}

function opTransformDeleteRetain(p, q) {
	switch(true) {
	case -p > q:
		return [-q, noop, q - p, noop]

	default:
		return [p, noop, noop, q + p]
	}
}

function opTransformDeleteDelete(p, q) {
	switch(true) {
	case -p > -q:
		return [noop, noop, q - p, noop]

	case -p < -q:
		return [noop, noop, noop, p - q]

	default:
		return [noop, noop, noop, noop]
	}
}

function opsSpanBase(u) {
	var l = 0
	u.forEach(function(o) {
		if(opType(o) == 'insert') {
			return
		}
		l += opSpan(o)
	})
	return l
}

function opsSpanTarget(u) {
	var l = 0
	u.forEach(function(o) {
		if(opType(o) == 'delete') {
			return
		}
		l += opSpan(o)
	})
	return l
}

function opsCompact(u) {
	var z = []
	for(var i = 0; i < u.length; i++) {
		var t = u[i]
		for(j = i+1; j < u.length && opType(u[i]) === opType(u[j]); j++) {
			t += u[j]
			i++
		}
		if(!t) {
			continue
		}
		z.push(t)
	}
	return z
}

function opsCompose(u, v) {
	if(opsSpanTarget(u) != opsSpanBase(v)) {
		throw errBadSpans
	}

	var z = []

	var l = newCursor(u)
	var r = newCursor(v)
	while(!cursorFin(l) || !cursorFin(r)) {
		var a = cursorOp(l)
		var b = cursorOp(r)

		switch(true) {
		case opType(a) === 'delete':
			z.push(a)
			cursorNext(l, noop)
			continue

		case opType(b) === 'insert':
			z.push(b)
			cursorNext(r, noop)
			continue
		}

		if(cursorFin(l)) {
			throw errTooShort
		}
		if(cursorFin(r)) {
			throw errTooLong
		}

		var c, p, q
		switch(opType(b)) {
		case 'retain':
			var e
			switch(opType(a)) {
			case 'retain':
				 e = opComposeRetainRetain(a, b)
				 break

			case 'insert':
				 e = opComposeInsertRetain(a, b)
				 break
			}
			c = e[0], p = e[1], q = e[2]
			break

		case 'delete':
			var e
			switch(opType(a)) {
			case 'retain':
				e = opComposeRetainDelete(a, b)
				break

			case 'insert':
				e = opComposeInsertDelete(a, b)
				break
			}
			c = e[0], p = e[1], q = e[2]
			break
		}

		z.push(c)
		cursorNext(l, p)
		cursorNext(r, q)
	}

	return opsCompact(z)
}

function opsTransform(u, v) {
	if(opsSpanBase(u) != opsSpanBase(v)) {
		throw errBadSpans
	}

	var up = []
	var vp = []

	var l = newCursor(u)
	var r = newCursor(v)
	while(!cursorFin(l) || !cursorFin(r)) {
		var a = cursorOp(l)
		var b = cursorOp(r)

		switch(true) {
		case opType(a) === 'insert':
			up.push(a)
			vp.push(opSpan(a))
			cursorNext(l, noop)
			continue

		case opType(b) === 'insert':
			up.push(opSpan(b))
			vp.push(b)
			cursorNext(r, noop)
			continue
		}

		if(cursorFin(l)) {
			throw errTooShort
		}
		if(cursorFin(r)) {
			throw errTooLong
		}

		var c, d, p, q
		switch(opType(b)) {
		case 'retain':
			var e
			switch(opType(a)) {
			case 'retain':
				 e = opTransformRetainRetain(a, b)
				 break

			case 'delete':
				 e = opTransformDeleteRetain(a, b)
				 break
			}
			c = e[0], d = e[1], p = e[2], q = e[3]
			break

		case 'delete':
			var e
			switch(opType(a)) {
			case 'retain':
				 e = opTransformRetainDelete(a, b)
				 break

			case 'delete':
				 e = opTransformDeleteDelete(a, b)
				 break
			}
			c = e[0], d = e[1], p = e[2], q = e[3]
			break
		}

		if(c != noop) {
			up.push(c)
		}
		if(d != noop) {
			vp.push(d)
		}
		cursorNext(l, p)
		cursorNext(r, q)
	}

	return [opsCompact(up), opsCompact(vp)]
}

function newCursor(ops) {
	return {
		ops: ops,
		nxt: null,
		i: 0
	}
}

function cursorNext(c, p) {
	c.nxt = p
	if(c.nxt == noop) {
		c.i++
	}
}

function cursorOp(c) {
	if(c.i >= c.ops.length) {
		return noop
	}
	if(c.nxt != noop && c.nxt != null) {
		return c.nxt
	}
	return c.ops[c.i]
}

function cursorFin(c) {
	return c.i >= c.ops.length
}

function blobApply(b, ops) {
	if(opsSpanBase(ops) != b.length) {
		throw errBadSpan
	}

	var z = ''

	var i = 0
	ops.forEach(function(o) {
		switch(opType(o)) {
		case 'retain':
			z += b.substr(i, i+o)
			i += o
			break

		case 'insert':
			z += o
			break

		case 'delete':
			i += -o
			break
		}
	})
	if(i !== b.length) {
		throw errTooShort
	}

	return z
}

try {
	module.exports = {
		opType: opType,
		opSpan: opSpan,
		opComposeRetainRetain: opComposeRetainRetain,
		opComposeRetainDelete: opComposeRetainDelete,
		opTransformRetainRetain: opTransformRetainRetain,
		opTransformRetainDelete: opTransformRetainDelete,
		opComposeInsertRetain: opComposeInsertRetain,
		opComposeInsertDelete: opComposeInsertDelete,
		opTransformDeleteRetain: opTransformDeleteRetain,
		opTransformDeleteDelete: opTransformDeleteDelete,
		opsSpanBase: opsSpanBase,
		opsSpanTarget: opsSpanTarget,
		opsCompose: opsCompose,
		opsTransform: opsTransform
	}
} catch(e) {}
