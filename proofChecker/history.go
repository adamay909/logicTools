package main

import (
	"fmt"
	"strconv"
	"strings"
	"syscall/js"

	"github.com/adamay909/logicTools/proofChecker/internal/editor"
)

var history []string
var historyPosition int

const historyItemSeparator = `<!---->`

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
	return domDocument.Call("querySelector", "#editorWindow").Get("innerHTML").String()
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
	history = nil

	raw := js.Global().Get("localStorage").Call("getItem", "history2")
	if !raw.IsNull() {
		temphistory := strings.Split(editor.StripWhiteSpaceFromHTML(raw.String()), historyItemSeparator)
		for i := range temphistory {
			if len(temphistory[i]) == 0 {
				continue
			}
			history = append(history, temphistory[i]+`<!---->`)
		}
	} else {
		fmt.Println("appending to history")
		history = append(history, getCurrentConsoleState())
	}

	var err error
	historyPosition, err = strconv.Atoi(js.Global().Get("localStorage").Call("getItem", "historyPosition2").String())
	if err != nil {
		historyPosition = 0
	}
	snap := js.Global().Get("localStorage").Call("getItem", "snapshot")
	if snap.IsNull() {
		moveInHistoryTo(historyPosition)
		return
	}
	fmt.Println("restoring snapshot")
	setCurrentConsoleState(snap.String())
	saveHistory()
}

func moveInHistoryTo(pos int) {
	setCurrentConsoleState(history[pos])
	historyPosition = pos
	updatePageNumber()
}

func updatePageNumber() {
	js.Global().Get("localStorage").Call("setItem", "historyPosition2", strconv.Itoa(historyPosition))
	domDocument.Call("querySelector", "#pagenumber").Set("innerHTML", strconv.Itoa(historyPosition+1)+"/"+strconv.Itoa(len(history)))
}

func saveHistory() {
	if historyPosition == len(history) {
		history = append(history, "")
	}
	saveSnapshot()
	history[historyPosition] = getCurrentConsoleState() + historyItemSeparator

	js.Global().Get("localStorage").Call("setItem", "history2", strings.Join(history, ""))
	js.Global().Get("localStorage").Call("setItem", "historyPosition2", strconv.Itoa(historyPosition))
}

func moveBackInHistory() {
	newpos := historyPosition - 1
	if newpos < 0 {
		return
	}
	saveHistory()
	moveInHistoryTo(newpos)
}

func moveForwardInHistory() {
	newpos := historyPosition + 1
	if newpos > len(history)-1 {
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

	historyPosition++
	var hb, hf []string
	hb = append(hb, history[:historyPosition]...)
	hf = append(hf, history[historyPosition:]...)
	history = nil
	history = append(history, hb...)
	history = append(history, getCurrentConsoleState())
	history = append(history, hf...)
	saveSnapshot()
	updatePageNumber()
}

func duplicateHistoryItem() {
	p1 := dsp.editor.GetInnerHTML()
	p2 := dsp.title.GetInnerHTML()
	newpage()
	dsp.editor.SetInnerHTML(p1)
	dsp.title.SetInnerHTML(p2)
	saveSnapshot()
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
