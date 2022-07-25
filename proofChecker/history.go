package main

import (
	"encoding/json"
	"strings"
	"syscall/js"
)

var history []string
var historyPosition int

func saveHistory() {

	c := marshalJson()

	history = append(history, c)

	historyPosition = len(history)

	json := strings.Join(history, "\n")

	js.Global().Get("localStorage").Call("setItem", "history", json)

}

func loadHistory() {

	json := js.Global().Get("localStorage").Call("getItem", "history").String()

	if json == "" {
		return
	}

	history = strings.Split(json, "\n")
	historyPosition = len(history)
	if historyPosition < 0 {
		historyPosition = 0
	}

}

func backHistory() {

	if historyPosition == 0 {
		return
	}

	historyPosition = historyPosition - 1

	moveInHistory()
}

func forwardHistory() {

	if historyPosition >= len(history) {
		return
	}

	historyPosition = historyPosition + 1

	moveInHistory()

}

func moveInHistory() {

	if historyPosition >= len(history) {
		return
	}

	json.Unmarshal([]byte(history[historyPosition]), dsp)

	dsp.xpos, dsp.ypos = 0, 0
	dsp.overhang = false
	dsp.modifier = ""
	if dsp.SystemPL != oPL {
		togglePL()
	}
	if dsp.Theorems != oTHM {
		toggleTheorems()
	}

	setTextByID("messages", "")
	hide("messages")
	display()
}

func saveState() {

	c := marshalJson()

	js.Global().Get("localStorage").Call("setItem", "current", c)

}

func recoverState() {

	data := js.Global().Get("localStorage").Call("getItem", "current").String()

	json.Unmarshal([]byte(data), dsp)

	dsp.xpos, dsp.ypos = 0, 0
	dsp.overhang = false
	dsp.modifier = ""
	if dsp.SystemPL != oPL {
		togglePL()
	}
	if dsp.Theorems != oTHM {
		toggleTheorems()
	}

	setTextByID("messages", "")
	hide("messages")
	display()
}

func rmFromHistory() {

	if historyPosition >= len(history) {
		return
	}

	h1 := history[:historyPosition]
	h2 := history[historyPosition+1:]

	history = nil

	history = append(history, h1...)
	history = append(history, h2...)

	js.Global().Get("localStorage").Call("setItem", "history", strings.Join(history, "\n"))

	if historyPosition == 0 {
		forwardHistory()
		return
	}
	backHistory()
}

func cleanHistory() {

	history = slicesCleanDuplicates(history)

	historyPosition = len(history)
	if historyPosition < 0 {
		historyPosition = 0
	}
}
