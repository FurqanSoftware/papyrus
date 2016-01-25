var assert = require('assert')
var ot = require('./ot')

describe('opsCompose()', function() {
	it('Fbb -> Foobaz -> Foobarbaz', function() {
		assert.deepEqual(ot.opsCompose([1, 'oo', 2, 'az'], [4, 'ar', 3]), [1, 'oo', 1, 'ar', 1, 'az'])
	})

	it('Fuubaz -> Fuubarbaz -> Foobarbaz', function() {
		assert.deepEqual(ot.opsCompose([3, 'bar', 3], [1, -2, 'oo', 6]), [1, -2, 'oo', 'bar', 3])
	})

	it('(209)sd ada ad (55) -> (181)(-40)(54) -> (181)(-28)(-1)(54)', function() {
		assert.deepEqual(ot.opsCompose([209, 'sd a da ad ', 55], [181, -40, 54]), [181, -28, -1, 54])
	})
})

describe('opsTransform()', function() {
	it('Fbb -> (Foobbaz, Fbarb) -> Foobarbaz', function() {
		assert.deepEqual(ot.opsTransform([1, 'oo', 2, 'az'], [2, 'ar', 1]), [[1, 'oo', 1, 2, 1, 'az'], [1, 2, 1, 'ar', 1, 2]])
	})

	it('Fuubbaz -> (Foobbaz, Fuubarbaz) -> Foobarbaz', function() {
		assert.deepEqual(ot.opsTransform([1, -2, 'oo', 4], [4, 'ar', 3]), [[1, -2, 'oo', 1, 2, 3], [1, 2, 1, 'ar', 3]])
	})
})
