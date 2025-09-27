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

func writeStateToHTML() {
	cl := domDocument.Call("querySelector", "#consoleState").Get("classList")
	if oPL {
		cl.Call("add", "oPL")
	} else {
		cl.Call("remove", "oPL")
	}
	if oTHM {
		cl.Call("add", "oTHM")
	} else {
		cl.Call("remove", "oTHM")
	}
	setupButtonLabels()
}

func setStateFromHTML() {

	cl := domDocument.Call("querySelector", "#consoleState").Get("classList")
	if cl.Call("contains", "oPL").Bool() {
		oPL = true
	} else {
		oPL = false
	}
	if cl.Call("contains", "oTHM").Bool() {
		oTHM = true
	} else {
		oTHM = false
	}
	setupButtonLabels()
}

func getCurrentConsoleState() string {
	return domDocument.Call("querySelector", "#editorWindow").Get("innerHTML").String()
}

func setCurrentConsoleState(html string) {
	dummy := domDocument.Call("createElement", "div")
	dummy.Set("innerHTML", html)
	state := dummy.Call("querySelector", "#consoleState").Get("classList")
	if state.Call("contains", "oPL").Bool() {
		oPL = true
	} else {
		oPL = false
	}
	if state.Call("contains", "oTHM").Bool() {
		oTHM = true
	} else {
		oTHM = false
	}
	writeStateToHTML()
	editorContent := dummy.Call("querySelector", "#editor").Get("innerHTML").String()
	titleContent := dummy.Call("querySelector", "#title").Get("innerHTML").String()
	dsp.editor.SetInnerHTML(editorContent)
	dsp.title.SetInnerHTML(titleContent)
}

func loadHistory() {
	history = nil
	//	history = append(history, getCurrentConsoleState())
	//	return

	raw := js.Global().Get("localStorage").Call("getItem", "history")
	if !raw.IsNull() {
		temphistory := strings.Split(editor.StripWhiteSpaceFromHTML(raw.String()), `<div id="consoleState"`)
		for i := range temphistory {
			if len(temphistory[i]) == 0 {
				continue
			}
			history = append(history, `<div id="consoleState"`+temphistory[i])
		}
	} else {
		fmt.Println("appending to history")
		history = append(history, getCurrentConsoleState())
	}

	var err error
	historyPosition, err = strconv.Atoi(js.Global().Get("localStorage").Call("getItem", "historyPosition").String())
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
	js.Global().Get("localStorage").Call("setItem", "historyPosition", strconv.Itoa(historyPosition))
	domDocument.Call("querySelector", "#pagenumber").Set("innerHTML", strconv.Itoa(historyPosition+1)+"/"+strconv.Itoa(len(history)))
}

func saveHistory() {
	if historyPosition == len(history) {
		history = append(history, "")
	}
	saveSnapshot()
	history[historyPosition] = getCurrentConsoleState()

	js.Global().Get("localStorage").Call("setItem", "history", strings.Join(history, ""))
	js.Global().Get("localStorage").Call("setItem", "historyPosition", strconv.Itoa(historyPosition))
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
	loadHistory()
}

func saveSnapshot() {
	js.Global().Get("localStorage").Call("setItem", "snapshot", getCurrentConsoleState())
}

func jsFuncSaveSnapshot(this js.Value, args ...any) {
	saveSnapshot()
}
