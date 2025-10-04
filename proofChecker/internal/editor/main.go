package editor

import (
	_ "embed" //embed template files
	"strconv"
	"strings"
	"syscall/js"
	"unicode"
)

const (
	simpleeditor int = iota
	derivationeditor
	axiomaticeditor
)

type Editor struct {
	id            string
	elem          js.Value
	key           string
	editorType    int
	rowTemplate   string
	derivTemplate string
}

//go:embed derivationTemplates.html
var mainTemplate string

/* AddCSS adds the CSS stylesheets needed for all the derivation formats as internal style sheets to the end of the head element. Should be called after any other manipulation of style sheets.
 */
func AddCSS() {
	dummy := domDocument.Call("createElement", "div")

	dummy.Set("innerHTML", mainTemplate)
	styleSheets := dummy.Call("getElementsByTagName", "style")

	for i := range styleSheets.Get("length").Int() {
		domDocument.Get("head").Call("appendChild", styleSheets.Call("item", i).Call("cloneNode", "deep"))
	}
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
	e.id = id
	e.MakeDerivationEditor()
	e.elem.Set("innerHTML", e.derivTemplate)

	return e
}

/*
NewAxiomaticEditor makes the element with the given id into a
natural deduction style sequent calculus editor. The element itself
should not be content editable. Approriate sublements will be
set up automatically.
Edits do not fire input events. Instead a custom event of type "editorinput" is fired that bubbles up and can be used by event listeners.
*/
func NewAxiomaticEditor(id string) *Editor {
	e := new(Editor)
	e.id = id
	e.MakeAxiomaticEditor()
	e.elem.Set("innerHTML", e.derivTemplate)
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
	e.removeEventListeners()
	setupJSeditor(e)
	return e
}

// MakeDerivationEditor turns e into a derivation editor but leaves its contents intact.
// So do not use this if you need to set up the markup of the inside of the editor.
func (e *Editor) MakeDerivationEditor() {
	e.editorType = derivationeditor
	e.removeEventListeners()
	e.setupEditorTemplates("ndseqcalc")
	setupJSderivationEditor(e)
	sel := domWindow.Call("getSelection")
	sel.Call("empty")
}

func (e *Editor) removeEventListeners() {
	elem, _ := getElementByID(e.id)
	cloneElem := elem.Call("cloneNode", "true")
	elem.Call("replaceWith", cloneElem)
	e.elem = cloneElem
	elem.Set("outerHTML", "")
}

func (e *Editor) MakeAxiomaticEditor() {
	e.editorType = axiomaticeditor
	e.removeEventListeners()
	e.setupEditorTemplates("axiomatic")
	setupJSderivationEditor(e)
	sel := domWindow.Call("getSelection")
	sel.Call("empty")
}

/*name is the CSS selector for the relevant portion of mainTemplate. The function looks for .derivation.{name} and .row.{name}.
 */
func (e *Editor) setupEditorTemplates(name string) {

	templ := domDocument.Call("createElement", "html")
	templ.Set("innerHTML", mainTemplate)
	e.derivTemplate = StripWhiteSpaceFromHTML(templ.Call("querySelector", ".derivation."+name).Get("outerHTML").String())
	e.rowTemplate = StripWhiteSpaceFromHTML(templ.Call("querySelector", ".row."+name).Get("outerHTML").String())

}

func (e *Editor) AddEventListener(event string, f func(this js.Value, args ...any), params ...any) {
	addEventListener(e.elem, event, f, params...)
}

func (e *Editor) Clear() {

	e.elem.Set("innerHTML", e.derivTemplate)

}

func (e *Editor) GetInnerHTML() string {
	return e.elem.Get("innerHTML").String()
}

func (e *Editor) SetInnerHTML(s string) {
	e.elem.Set("innerHTML", s)
}

func (e *Editor) SetOffset(v int) {
	d := e.elem.Call("querySelector", ".derivation")
	s := d.Get("style")
	s.Call("setProperty", "counter-set", "linecounter "+strconv.Itoa(v-1))
}

func setupJSderivationEditor(e *Editor) {

	addEventListener(e.elem, "click", clickHandler, e)
	addEventListener(e.elem, "focusout", focusOutHandler, e)
	addEventListener(e.elem, "selectionchange", caretHandler, e)
	addEventListener(e.elem, "customselectionchange", caretHandler, e)
	addEventListener(e.elem, "beforeinput", beforeInputHandler, e)
	addEventListener(e.elem, "keydown", keyHandler, e)
	addEventListener(e.elem, "paste", pasteHandler, e)

}

func setupJSeditor(e *Editor) {
	addEventListener(e.elem, "beforeinput", simpleInputHandler, e)
	addEventListener(e.elem, "paste", pasteHandler, e)
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
