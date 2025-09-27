package editor

import (
	_ "embed" //embed template files
	"strings"
	"syscall/js"
	"unicode"
)

const (
	simpleeditor int = iota
	derivationeditor
)

type Editor struct {
	id         string
	elem       js.Value
	key        string
	editorType int
}

//go:embed editorTemplate.txt
var edTempl string

//go:embed rowTemplate.txt
var rowTemplate string

func init() {
	rowTemplate = StripWhiteSpaceFromHTML(rowTemplate)
	edTempl = StripWhiteSpaceFromHTML(edTempl)
}

/*
NewDerivationEditor makes the element with the given id into a
natural deduction style sequent calculus editor. The element itself
should not be content editable. Approriate sublements will be
set up automatically.
Edits do not fire input events. Instead a custom event of type "editorinput" is fired that bubbles up and can be used by event listeners.
*/
func NewDerivationEditor(id string) *Editor {
	e := new(Editor)
	e.editorType = derivationeditor
	e.id = id
	e.elem, _ = getElementByID(id)
	html := StripWhiteSpaceFromHTML(edTempl)
	e.elem.Set("innerHTML", html)
	deriv := e.elem.Get("firstChild")
	deriv.Set("innerHTML", rowTemplate)
	setupJSderivationEditor(e)
	sel := domWindow.Call("getSelection")
	sel.Call("empty")
	return e
}

/*
NewSimpleEditor makes the element with the given if into a simple editor. The element should be content-editable (that's your responsibility).
Edits do not fire input events. Instead a custom event of type "editorinput" is fired that bubbles up and can be used by event listeners.
*/
func NewSimpleEditor(id string) *Editor {
	e := new(Editor)
	e.editorType = simpleeditor
	e.id = id
	e.elem, _ = getElementByID(id)
	setupJSeditor(e)
	return e
}

func (e *Editor) AddEventListener(event string, f func(this js.Value, args ...any), params ...any) {
	addEventListener(e.elem, event, f, params...)
}

func (e *Editor) Clear() {
	switch e.editorType {

	case simpleeditor:
		e.elem.Set("innerHTML", "")

	case derivationeditor:
		e.elem.Set("innerHTML", edTempl)
		deriv := e.elem.Get("firstChild")
		deriv.Set("innerHTML", rowTemplate)
	}
}

func (e *Editor) GetInnerHTML() string {
	return e.elem.Get("innerHTML").String()
}

func (e *Editor) SetInnerHTML(s string) {
	e.elem.Set("innerHTML", s)
}

func setupJSderivationEditor(e *Editor) {

	addEventListener(e.elem, "click", clickHandler, e)
	addEventListener(e.elem, "focusout", focusOutHandler, e)
	addEventListener(domDocument, "selectionchange", caretHandler, e)
	addEventListener(e.elem, "beforeinput", beforeInputHandler, e)
	addEventListener(e.elem, "keydown", keyHandler, e)
	addEventListener(e.elem, "paste", pasteHandler, e)

}

func setupJSeditor(e *Editor) {
	addEventListener(e.elem, "beforeinput", simpleInputHandler, e)
	addEventListener(e.elem, "paste", pasteHandler, e)
}

var modifier string

func keyHandler(this js.Value, args ...any) {
	ed := args[0].(*Editor)
	switch this.Get("key").String() {

	case "ArrowUp":
		ed.key = "up"
		this.Call("preventDefault")
		ed.elem.Call("dispatchEvent", js.Global().Get("Event").New("selectionchange"))

	case "ArrowDown":
		ed.key = "down"
		this.Call("preventDefault")
		ed.elem.Call("dispatchEvent", js.Global().Get("Event").New("selectionchange"))

	case "ArrowRight":
		if caretAtEndOfCell() {
			ed.key = "right"
			this.Call("preventDefault")
			ed.elem.Call("dispatchEvent", js.Global().Get("Event").New("selectionchange"))
		}
	case "ArrowLeft":
		if caretAtStartOfCell() {
			ed.key = "left"
			this.Call("preventDefault")
			ed.elem.Call("dispatchEvent", js.Global().Get("Event").New("selectionchange"))
		}

	}
}

func StripWhiteSpaceFromHTML(html string) string {

	b := new(strings.Builder)
	b.Grow(len(html))

	inTag := false
	for _, r := range html {

		if r == '<' {
			inTag = true
		}

		if inTag {
			b.WriteRune(r)
		} else {
			if !unicode.IsSpace(r) {
				b.WriteRune(r)
			}
		}

		if r == '>' {
			inTag = false
		}

	}

	return b.String()
}

func focusOutHandler(this js.Value, args ...any) {

	ed := args[0].(*Editor)

	focusElement := domDocument.Get("activeElement")
	if !ed.elem.Call("contains", focusElement).Bool() {
		domDocument.Call("getSelection").Call("empty")
	}
	return
}

func clickHandler(this js.Value, args ...any) {
	focusElement := domDocument.Get("activeElement")
	ed := args[0].(*Editor)
	if ed.elem.Call("contains", focusElement).Bool() {
		return
	}
	e := ed.elem.Get("firstElementChild").Get("firstElementChild").Get("firstElementChild")
	e.Call("focus")
}
