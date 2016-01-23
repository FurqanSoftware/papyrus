var errBadType = new Error('ot.js: invalid operation type')
var errBadSpans = new Error('ot.js: base span doesn\'t match target span')
var errTooShort = new Error('ot.js: operations are too short')
var errTooLong = new Error('ot.js: operations are too long')
var errBadOpPair = new Error('ot.js: invalid pair of operations')

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
