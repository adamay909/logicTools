package editor

import (
	_ "embed" //embed
	"slices"
	"syscall/js"
)

var simpleModifier string

func dispatchEditEvent() {
	ev := js.Global().Get("CustomEvent").New("editorinput", map[string]any{"bubbles": true})
	cellAtCaret().Call("dispatchEvent", ev)
}

func simpleInputHandler(this js.Value, args ...any) {

	defer func() {
		dispatchEditEvent()
	}()

	inputType := this.Get("inputType").String()
	if inputType != "insertText" {
		return
	}

	this.Call("preventDefault")
	char := this.Get("data").String()

	if char == `\` || char == `|` {
		if simpleModifier == "" {
			simpleModifier = char
			return
		}
	}
	char = simpleModifier + char
	insertChar(simpleModCharOf(char))
	simpleModifier = ""

}

func beforeInputHandler(this js.Value, args ...any) {

	defer func() {
		dispatchEditEvent()
	}()

	inputType := this.Get("inputType").String()
	switch inputType {

	case "deleteContentBackward":
		if caretAtStartOfLine() {
			moveCaretToEndOfPreviousLine()
			n := nextRowOfCaret()
			if isEmptyRow(n) {
				n.Call("remove")
			}
			return
		}
		if caretAtStartOfCell() {
			moveFocusToLeftOf(cellAtCaret())
		}

	case "deleteContentForward":
		if caretAtEndOfLine() {
			if onLastLine() {
				return
			}
			moveFocusToBelowOf(cellAtCaret())
			moveCaretToStartOfLine()
			p := previousRowOfCaret()
			if isEmptyRow(p) {
				p.Call("remove")
			}
			return
		}
		if caretAtEndOfCell() {
			moveFocusToRightOf(cellAtCaret())
		}

	case "insertLineBreak":
		this.Call("preventDefault")
		insertNewRow()
		return

	}

	if inputType != "insertText" {
		return
	}

	this.Call("preventDefault")
	char := this.Get("data").String()

	if !permitted(char) {
		return
	}

	if char == `\` || char == `|` {
		modifier = char
		return
	}

	char = modifier + char
	insertChar(modCharOf(char))
	modifier = ""

}

func insertChar(c string) {
	insertchar(c)
}

func insertchar(c string) {
	selection := domWindow.Call("getSelection")
	srange := selection.Call("getRangeAt", 0)
	textNode := domDocument.Call("createTextNode", c)
	srange.Call("deleteContents")
	srange.Call("insertNode", textNode)
	srange.Call("setStartAfter", textNode)
	selection.Call("addRange", srange)
}

func insertNewRow() {
	r := domDocument.Call("createElement", "div")
	cellAtCaret().Get("parentElement").Call("after", r)
	r.Set("outerHTML", rowTemplate)
	moveFocusToBelowOf(cellAtCaret())
	moveCaretToStartOfLine()
}

func permitted(c string) bool {
	return slices.Contains(permittedChars, c[:1])
}

func modCharOf(c string) string {
	v, ok := charMap[c]
	if ok {
		return v
	}
	return c
}

func simpleModCharOf(c string) string {
	v, ok := extraMap[c]
	if ok {
		return v
	}
	v, ok = charMap[c]
	if ok {
		return v
	}
	return c
}

func pasteHandler(this js.Value, args ...any) {
	this.Call("preventDefault")
}
