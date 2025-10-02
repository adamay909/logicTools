package main

import (
	"strconv"
	"syscall/js"
)

// we store state information about the proof checker (dsp values) as
// value data-* attributes of #consoleState which is a subelement
// of #editorWindow. We store the the innerHTML of #editorWindow as
// snapshots.
func writeStateToHTML() {
	cs := domDocument.Call("querySelector", "#consoleState")
	if dsp.oPL {
		cs.Get("dataset").Set("opl", "true")
	} else {
		cs.Get("dataset").Delete("opl")
	}
	if dsp.oTHM {
		cs.Get("dataset").Set("othm", "true")
	} else {
		cs.Get("dataset").Delete("othm")
	}

	cs.Get("dataset").Set("offset", strconv.Itoa(dsp.oOffset))
	setupButtonLabels()
	saveSnapshot()
}

// elem is the element from which to read the proofchecker state
func setStateFromHTML(elem js.Value) {

	cl := elem.Call("querySelector", "#consoleState").Get("classList")
	if cl.Get("length").Int() != 0 {
		if cl.Call("contains", "oPL").Bool() {
			dsp.oPL = true
			cl.Call("remove", "oPL")
		} else {
			dsp.oPL = false
		}
		if cl.Call("contains", "oTHM").Bool() {
			dsp.oTHM = true
			cl.Call("remove", "oTHM")
		} else {
			dsp.oTHM = false
		}
		return
	}
	cs := elem.Call("querySelector", "#consoleState")
	if cs.Get("dataset").Get("opl").String() == "true" {
		dsp.oPL = true
	} else {
		dsp.oPL = false
	}
	if cs.Get("dataset").Get("othm").String() == "true" {
		dsp.oTHM = true
	} else {
		dsp.oTHM = false
	}
	v, err := strconv.Atoi(cs.Get("dataset").Get("offset").String())
	if err != nil {
		dsp.oOffset = 1
	} else {
		dsp.oOffset = v
	}
	setupButtonLabels()
	setupGentzen()
}

func getCurrentConsoleState() string {
	return dsp.window.Get("outerHTML").String()
}

func setCurrentConsoleState(html string) {
	dummy := domDocument.Call("createElement", "div")
	dummy.Set("innerHTML", html)
	setStateFromHTML(dummy)
	writeStateToHTML()
	editorContent := dummy.Call("querySelector", "#editor").Get("innerHTML").String()
	titleContent := dummy.Call("querySelector", "#title").Get("innerHTML").String()
	dsp.editor.SetInnerHTML(editorContent)
	dsp.title.SetInnerHTML(titleContent)
	dsp.editor.SetOffset(dsp.oOffset)
	setupGentzen()
}

func loadHistory() {

	raw := js.Global().Get("localStorage").Call("getItem", "history3")
	if !raw.IsNull() {
		dsp.history.Set("innerHTML", raw.String())
	} else {
		dsp.history.Call("append", dsp.window.Call("cloneNode", "deep"))
	}

	dsp.historyItems = dsp.history.Get("children")

	var err error
	dsp.historyPosition, err = strconv.Atoi(js.Global().Get("localStorage").Call("getItem", "dsp.historyPosition2").String())
	if err != nil {
		dsp.historyPosition = 0
	}
	snap := js.Global().Get("localStorage").Call("getItem", "snapshot")
	if snap.IsNull() {
		moveInHistoryTo(dsp.historyPosition)
		return
	}
	setCurrentConsoleState(snap.String())
	saveHistory()
}

func saveHistory() {
	saveSnapshot()

	dsp.historyItems.Call("item", dsp.historyPosition).Set("outerHTML", getCurrentConsoleState())
	//	Call("replaceWith", dsp.window.Call("cloneNode", "deep"))
	js.Global().Get("localStorage").Call("setItem", "history3", dsp.history.Get("innerHTML"))
	js.Global().Get("localStorage").Call("setItem", "dsp.historyPosition2", strconv.Itoa(dsp.historyPosition))
}

func moveInHistoryTo(pos int) {
	setCurrentConsoleState(dsp.historyItems.Call("item", pos).Get("innerHTML").String())
	dsp.historyPosition = pos
	updatePageNumber()
}

func updatePageNumber() {
	js.Global().Get("localStorage").Call("setItem", "dsp.historyPosition2", strconv.Itoa(dsp.historyPosition))
	domDocument.Call("querySelector", "#pagenumber").Set("innerHTML", strconv.Itoa(dsp.historyPosition+1)+"/"+strconv.Itoa(dsp.historyItems.Get("length").Int()))
}

func moveBackInHistory() {
	newpos := dsp.historyPosition - 1
	if newpos < 0 {
		return
	}
	saveHistory()
	moveInHistoryTo(newpos)
}

func moveForwardInHistory() {
	newpos := dsp.historyPosition + 1
	if newpos > dsp.historyItems.Get("length").Int()-1 {
		return
	}
	saveHistory()
	moveInHistoryTo(newpos)
}

func newpage() {
	saveHistory()
	dsp.oOffset = 1
	dsp.editor.Clear()
	dsp.title.Clear()
	writeStateToHTML()

	dsp.historyItems.Call("item", dsp.historyPosition).Call("after", dsp.window.Call("cloneNode", "deep"))

	dsp.historyPosition++
	saveSnapshot()
	updatePageNumber()
}

func duplicateHistoryItem() {
	dsp.historyItems.Call("item", dsp.historyPosition).Call("after", dsp.window.Call("cloneNode", "deep"))
	dsp.historyPosition++
	saveSnapshot()
	updatePageNumber()
}

func clearHistory() {
	js.Global().Get("localStorage").Call("clear")
	dsp.editor.Clear()
	dsp.title.Clear()
	dsp.oOffset = 1
	dsp.oPL = false
	dsp.oTHM = false
	writeStateToHTML()
	loadHistory()
}

func saveSnapshot() {
	js.Global().Get("localStorage").Call("setItem", "snapshot", getCurrentConsoleState())
}

func cleanupEditorWindow() {
	s := domDocument.Call("querySelector", "#editorWindow").Get("classList")
	s.Call("remove", "fail")
	s.Call("remove", "success")
}
