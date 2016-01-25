function codemirrorDocLength(doc) {
	return doc.indexFromPos({ line: doc.lastLine(), ch: 0 }) + doc.getLine(doc.lastLine()).length
}

function codemirrorTextLen(text) {
	return text.reduce(function(m, v) { return m+v.length+1 }, -1)
}

function codemirrorChangeLen(ch) {
	return codemirrorTextLen(ch.text) - codemirrorTextLen(ch.removed)
}

function codemirrorChangeToOps(doc, docLen, ch) {
	var rl = doc.indexFromPos(ch.from)
	var d = -codemirrorTextLen(ch.removed)
	var i = ch.text.join('\n')
	var rr = docLen - rl - codemirrorTextLen(ch.text)
	return [rl, d, i, rr].filter(function(v) {
		return !!v
	})
}

var editor = CodeMirror.fromTextArea($('#editor').val('')[0], {
	lineNumbers: true
})
var opsBuf = null
var opsWait = null
var opsWaitID = ""
var opsRoot = 0
var skipChange = false
editor.on('changes', function(event, changes) {
	if(skipChange) {
		skipChange = false
		return
	}
	var opsTmp = null
	var docLen = codemirrorDocLength(event.doc)
	changes.slice().reverse().forEach(function(ch) {
		if(ch.origin === 'setValue') {
			return
		}
		if(opsTmp === null) {
			opsTmp = codemirrorChangeToOps(event.doc, docLen, ch)
		} else {
			opsTmp = opsCompose(codemirrorChangeToOps(event.doc, docLen, ch), opsTmp)
		}
		docLen -= codemirrorChangeLen(ch)
	})
	if(opsBuf === null) {
		opsBuf = opsTmp
	} else {
		opsBuf = opsCompose(opsBuf, opsTmp)
	}
	sync()
})

var socket = glue()
socket.send('attach '+token)
socket.onMessage(function(data) {
	var fields = [data.substr(0, data.indexOf(' ')), data.substr(data.indexOf(' ')+1)]
	switch(fields[0]) {
	case 'change':
		var data = JSON.parse(fields[1])

		opsRoot = data.root

		if(data.id === "") {
			editor.getDoc().setValue(blobApply(editor.getDoc().getValue(), data.ops))

		} else if(data.id === opsWaitID) {
			opsWait = null
			sync()

		} else {

			var ops = data.ops
			if(opsWait !== null) {
				var e = opsTransform(opsWait, ops)
				opsWait = e[0]
				ops = e[1]
			}
			if(opsBuf !== null) {
				var e = opsTransform(opsBuf, ops)
				opsBuf = e[0]
				ops = e[1]
			}

			var doc = editor.getDoc()
			skipChange = true
			editor.operation(function() {
				var i = 0
				ops.forEach(function(o) {
					switch(opType(o)) {
					case 'retain':
						i += o
						break

					case 'insert':
						editor.replaceRange(o, editor.posFromIndex(i))
						i += o.length
						break

					case 'delete':
						editor.replaceRange('', editor.posFromIndex(i), editor.posFromIndex(i + (-o)))
						break
					}
				})
			})
		}
		break
	}
})
socket.on("connected", function(){
	$('#connectionStatus').removeClass("fa-ban").addClass("fa-flash")
})
socket.on("error", function(){
	$('#connectionStatus').removeClass("fa-flash").addClass("fa-ban")
})
socket.on("disconnected", function(){
	$('#connectionStatus').removeClass("fa-flash").addClass("fa-ban")
})
socket.on("connect_timeout", function(){
	$('#connectionStatus').removeClass("fa-flash").addClass("fa-ban")
})
socket.on("timeout", function(){
	$('#connectionStatus').removeClass("fa-flash").addClass("fa-ban")
})

function sync() {
	if(opsWait !== null || opsBuf === null) {
		return
	}

	opsWait = opsBuf
	opsBuf = null
	opsWaitID = new Date().getTime()+'-'+Math.round(1e9*Math.random())
	socket.send('change '+JSON.stringify({
		id: opsWaitID,
		root: opsRoot,
		ops: opsWait
	}))
}


var converter = new showdown.Converter()
editor.on('change', _.throttle(function() {
	$('#preview').html(converter.makeHtml(editor.getValue()))
}, 500))
