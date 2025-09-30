package main

import (
	"syscall/js"
)

var (
	uint8ArrayConstructor,
	fileConstructor,
	arrayConstructor,
	blobConstructor,
	domWindow,
	domDocument,
	domHTML,
	domBody js.Value
)

func init() {

	domWindow = js.Global()

	domDocument = domWindow.Get("document")

}

// addEventListener adds an event listener. The function f is called
// with the event object as its first argument and the arguments given by params.
func addEventListener(elem js.Value, eventType string, f func(event js.Value, args ...any), params ...any) {

	elem.Call("addEventListener", eventType, js.FuncOf(func(this js.Value, margs []js.Value) any {

		f(margs[0], params...)
		return true
	}), true)

	return
}

func jsWrapper(f func(event js.Value, args ...any), params ...any) js.Func {

	return js.FuncOf(func(this js.Value, margs []js.Value) any {
		f(margs[0], params...)
		return true
	})

}

func removeEventListener(elem js.Value, eventType string, f func(event js.Value, args ...any), params ...any) {

	elem.Call("removeEventListener", eventType, js.FuncOf(func(this js.Value, margs []js.Value) any {
		f(margs[0], params...)
		return true
	}), true)

	return

}

func jsWrap(f func()) (fn func(this js.Value, args ...any)) {

	fn = func(this js.Value, args ...any) {
		f()
	}

	return fn
}
