package editor

import "syscall/js"

func keyHandler(this js.Value, args ...any) {
	ed := args[0].(*Editor)
	focusElement := domDocument.Get("activeElement")
	switch this.Get("key").String() {

	case "ArrowUp":
		ed.key = "up"
		this.Call("preventDefault")
		focusElement.Call("dispatchEvent", js.Global().Get("CustomEvent").New("customselectionchange", map[string]any{"bubbles": true}))

	case "ArrowDown":
		ed.key = "down"
		this.Call("preventDefault")
		focusElement.Call("dispatchEvent", js.Global().Get("CustomEvent").New("customselectionchange", map[string]any{"bubbles": true}))

	case "ArrowRight":
		if caretAtEndOfCell() {
			ed.key = "right"
			this.Call("preventDefault")
			focusElement.Call("dispatchEvent", js.Global().Get("CustomEvent").New("customselectionchange", map[string]any{"bubbles": true}))
		}

	case "ArrowLeft":
		if caretAtStartOfCell() {
			ed.key = "left"
			this.Call("preventDefault")
			focusElement.Call("dispatchEvent", js.Global().Get("CustomEvent").New("customselectionchange", map[string]any{"bubbles": true}))
		}

	case "Home":
		ed.key = "home"
		if this.Get("ctrlKey").Bool() {
			ed.key = "ctrlHome"
		}
		this.Call("preventDefault")
		focusElement.Call("dispatchEvent", js.Global().Get("CustomEvent").New("customselectionchange", map[string]any{"bubbles": true}))

	case "End":
		ed.key = "end"
		if this.Get("ctrlKey").Bool() {
			ed.key = "ctrlEnd"
		}
		this.Call("preventDefault")
		focusElement.Call("dispatchEvent", js.Global().Get("CustomEvent").New("customselectionchange", map[string]any{"bubbles": true}))

	}

}
