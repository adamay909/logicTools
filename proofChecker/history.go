package main

import (
	"encoding/json"
	"strings"
	"syscall/js"

	"honnef.co/go/js/dom/v2"
)

var history []string
var historyPosition int

var stash string

func saveHistory() {

	c := dsp.marshalJson()

	history = append(history, c)

	cleanHistory()

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

}

func backHistory() {

	if historyPosition == 0 {
		return
	}

	if historyPosition == len(history) {
		stashState()
	}

	historyPosition = historyPosition - 1

	moveInHistory()
}

func forwardHistory() {

	historyPosition = historyPosition + 1

	if historyPosition > len(history) {
		return
	}

	if historyPosition == len(history) {
		reloadStash()
		return
	}

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

	c := dsp.marshalJson()

	js.Global().Get("localStorage").Call("setItem", "current", c)

}

func stashState() {

	stash = dsp.marshalJson()

}

func reloadStash() {
	if stash == "" {
		return
	}

	json.Unmarshal([]byte(stash), dsp)
	stash = ""
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

func exportHistory() {

	obj := js.Global().Get("Blob").New([]any{strings.Join(history, "\n")})

	url := js.Global().Get("URL").Call("createObjectURL", obj).String()

	stopInput()
	hide("console")
	hide("controls2")
	show("extra")
	hide("txtinput")
	show("historyDialog")
	show("backButton")
	hide("console")

	setTextByID("historyDialog", `<a href="`+url+`">right-click to download history</a>`)

}

func importHistory() {
	stopInput()
	hide("console")
	hide("controls2")
	show("extra")
	hide("txtinput")
	show("historyDialog")
	show("backButton")
	hide("console")

	html := `<h3>Paste JSON</h3>
<textarea name="textarea" id="historyinputarea" rows="15" cols="40"></textarea>
	 <button id="importHistory">Import</button>`

	setTextByID("historyDialog", html)
}

func rewriteHistory() {
	stopInput()
	history = strings.Split(dom.GetWindow().Document().GetElementByID("historyinputarea").(*dom.HTMLTextAreaElement).Value(), "\n")
	dom.GetWindow().Document().GetElementByID("historyinputarea").(*dom.HTMLTextAreaElement).SetValue("")
	dsp.clear()
	display()
	printMessage("")
	hide("messages")
	hide("historyDialog")
	hide("extra")
	hide("backButton")
	show("console")
	stopInput()

}
