var assert = require('assert')
var ot = require('./ot')

describe('opsCompose()', function() {
	it('Fbb -> Foobaz -> Foobarbaz', function() {
		assert.deepEqual(ot.opsCompose([1, 'oo', 2, 'az'], [4, 'ar', 3]), [1, 'oo', 1, 'ar', 1, 'az'])
	})

	it('Fuubaz -> Fuubarbaz -> Foobarbaz', function() {
		assert.deepEqual(ot.opsCompose([3, 'bar', 3], [1, -2, 'oo', 6]), [1, -2, 'oo', 'bar', 3])
	})
})
