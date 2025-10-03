package editor

import (
	"fmt"
	"strings"
	"syscall/js"
)

type caretPos struct {
	x    int //the character count from start of line
	line int //line number
}

func currentCaretlocation() caretPos {
	var c caretPos
	return c
}

func caretHandler(this js.Value, args ...any) {

	ed := args[0].(*Editor)

	focusElement := domDocument.Get("activeElement")
	rejectSelection()
	this.Call("preventDefault")

	switch ed.key {
	case "up":
		moveFocusToAboveOf(focusElement)
		ed.key = ""

	case "down":
		moveFocusToBelowOf(focusElement)
		ed.key = ""

	case "left":
		if caretAtStartOfLine() && !onFirstLine() {
			moveFocusToAboveOf(focusElement)
			moveCaretToEndOfLine()
			return
		}
		moveFocusToLeftOf(focusElement)
		ed.key = ""

	case "right":
		if caretAtEndOfLine() && !onLastLine() {
			moveFocusToBelowOf(focusElement)
			moveCaretToStartOfLine()
			return
		}
		moveFocusToRightOf(focusElement)
		ed.key = ""

	case "home":
		moveCaretToStartOfLine()

	case "ctrlHome":
		moveCaretToFirstLine()
		moveCaretToStartOfLine()

	case "end":
		moveCaretToEndOfLine()

	case "ctrlEnd":
		moveCaretToLastLine()
		moveCaretToEndOfLine()
	}

}

func moveFocusToRightOf(e js.Value) {
	t := e.Get("nextElementSibling")
	if t.IsNull() {
		return
	}
	switch valOf(t, "class") {
	case "dtstl", "dsep":
		t = t.Get("nextElementSibling")
	}
	t.Call("focus")
}

func moveFocusToLeftOf(e js.Value) {
	t := e.Get("previousElementSibling")
	if t.IsNull() {
		return
	}
	if valOf(t, "class") == "dtstl" || valOf(t, "class") == "dsep" {
		t = t.Get("previousElementSibling")
	}
	t.Call("focus")
	moveCaretToEndOfCell()
}

func moveCaretToEndOfCell() {
	sel := domWindow.Call("getSelection")
	sel.Call("modify", "move", "forward", "line")
}

func lineOfCaret() (r js.Value) {
	e := cellAtCaret()
	return e.Get("parentElement")
}

func moveCaretToFirstLine() {
	for ; !onFirstLine(); moveFocusToAboveOf(cellAtCaret()) {
	}
}

func moveCaretToLastLine() {
	for ; !onLastLine(); moveFocusToBelowOf(cellAtCaret()) {
	}
}

func moveFocusToAboveOf(e js.Value) {
	r := e.Get("parentElement")
	if r.IsNull() {
		return
	}
	r = r.Get("previousElementSibling")
	if r.IsNull() {
		return
	}
	t := childOfClass(r, e.Get("className").String())
	if t.IsNull() {
		fmt.Println("couldn't find corresponding cell")
		return
	}
	c := charCountToCaret()
	t.Call("focus")
	moveCaretForwardBy(c)
}

func moveFocusToBelowOf(e js.Value) {
	r := e.Get("parentElement")
	if r.IsNull() {
		return
	}
	r = r.Get("nextElementSibling")
	if r.IsNull() {
		return
	}
	t := childOfClass(r, e.Get("className").String())
	if t.IsNull() {
		fmt.Println("couldn't find corresponding cell")
		return
	}
	c := charCountToCaret()
	t.Call("focus")
	moveCaretForwardBy(c)
}

func moveCaretForwardBy(count int) {

	sel := domWindow.Call("getSelection")

	for _ = range count {
		sel.Call("modify", "move", "forward", "character")
	}
}

func cellIsEmpty(e js.Value) bool {
	return textOfCell(e) == ""
}

func cellToLeftEmpty(e js.Value) bool {
	t := e.Get("previousElementSibling")
	if t.IsNull() {
		return false
	}
	return textOfCell(t) == ""
}

func lineEmptyToRight(e js.Value) bool {

	for c := cellAtCaret(); !c.IsNull(); c = c.Get("nextElementSibling") {
		if textOfCell(c) != "" {
			return false
		}
	}
	return true
}

func cellAtCaret() js.Value {
	e := domDocument.Get("activeElement")
	for f := e; !f.IsNull(); f = e.Get("parentElement") {
		if f.Get("parentElement").Get("classList").Call("contains", "row").Bool() {
			return f
		}
	}
	return js.ValueOf(nil)
}

func caretAtStartOfCell() bool {
	return textToCaret() == ""
}

func caretAtEndOfCell() bool {
	cellrange := domDocument.Call("createRange")
	cellrange.Call("selectNodeContents", cellAtCaret())
	return textToCaret() == strings.Join(strings.Fields(cellrange.Call("toString").String()), "")
}

func caretAtStartOfLine() bool {
	if !cellAtCaret().Get("previousElementSibling").IsNull() {
		return false
	}
	return caretAtStartOfCell()
}

func caretAtEndOfLine() bool {
	if !cellAtCaret().Get("nextElementSibling").IsNull() {
		return false
	}
	return caretAtEndOfCell()
}

func onLastLine() bool {
	return nextRowOfCaret().IsNull()
}

func onFirstLine() bool {
	return previousRowOfCaret().IsNull()
}

func moveCaretToEndOfPreviousLine() {
	r := cellAtCaret().Get("parentElement").Get("previousElementSibling")
	if r.IsNull() {
		return
	}
	e := r.Get("lastElementChild")
	//for ; isEditable(e) && textOfCell(e) == ""; e = e.Get("previousElementSibling") {
	//}
	e.Call("focus")
	moveCaretToEndOfCell()
}

func moveCaretToEndOfLine() {
	e := cellAtCaret().Get("parentElement").Get("lastElementChild")
	for ; e.Get("textContent").String() == ""; e = e.Get("previousElementSibling") {
	}
	e.Call("focus")
	moveCaretToEndOfCell()
}

func moveCaretToStartOfLine() {
	e := cellAtCaret().Get("parentElement").Get("firstElementChild")
	e.Call("focus")
}

// Returns text content of cell up to end of row
func textToCaret() string {
	w := domDocument.Call("createTreeWalker", cellAtCaret(), 4)

	txt := ""
	r := domWindow.Call("getSelection").Call("getRangeAt", 0)
	s := r.Get("startContainer")
	for e := w.Get("root"); !e.IsNull(); e = w.Call("nextNode") {
		if e.Call("isSameNode", s).Bool() {
			sc := string([]rune(s.Get("textContent").String())[:r.Get("startOffset").Int()])
			txt = txt + sc
			break
		}
		if e.Get("nodeType").Int() != 3 {
			continue
		}

		txt = txt + e.Get("textContent").String()
	}
	return strings.Join(strings.Fields(txt), "")
}

func textOfCell(e js.Value) string {
	w := domDocument.Call("createTreeWalker", e, 4)
	txt := ""

	for n := w.Get("root"); !n.IsNull(); n = w.Call("nextNode") {
		if n.Get("nodeType").Int() != 3 {
			continue
		}
		txt = txt + n.Get("textContent").String()
	}
	return strings.Join(strings.Fields(txt), "")
}

func charCountToCaret() int {
	return len([]rune(textToCaret()))
}

func textOfLine() string {
	t := cellAtCaret().Get("parentElement").Get("textContent").String()
	return strings.Join(strings.Fields(t), "")
}

func lineIsEmpty() bool {
	return textOfLine() == turnstile+ldots
}

func nextLineIsEmpty() bool {
	r := cellAtCaret().Get("parentElement").Get("nextElementSibling")
	if r.IsNull() {
		return false
	}
	return strings.Join(strings.Fields(r.Get("textContent").String()), "") == turnstile+ldots
}

func isEmptyRow(r js.Value) bool {
	if r.IsNull() {
		return false
	}
	return strings.Join(strings.Fields(r.Get("textContent").String()), "") == turnstile+ldots
}

func nextRowOfCaret() js.Value {
	return cellAtCaret().Get("parentElement").Get("nextElementSibling")
}

func previousRowOfCaret() js.Value {
	return cellAtCaret().Get("parentElement").Get("previousElementSibling")
}
