package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"syscall/js"

	"honnef.co/go/js/dom/v2"
)

var history []string
var historyPosition int

var stash string

func readHistoryFromFile() {

	blob := js.Global().Get("document").Call("getElementById", "inputfile").Get("files").Index(0).Call("text").Call("then").String()

	//	s := js.Global().Get("Promise").Call("resolve", blob).String()

	fmt.Println("file read @")
	fmt.Println(blob)
}

func insertEmptyHistoryItem() {

	saveHistory()
	if historyPosition == len(history) {
		h2 := ""
		history = append(history, h2)
		//		cleanHistory()
	} else {

		h1 := history[:historyPosition]
		h2 := history[historyPosition:]

		history = nil

		history = append(history, h1...)
		history = append(history, "")
		history = append(history, h2...)
	}

	saveHistory()
	forwardHistory()
	moveInHistory()
	display()
	return

}

func duplicateHistoryItem() {
	insertEmptyHistoryItem()
	history[historyPosition] = history[historyPosition-1]
	saveHistory()
	moveInHistory()
	display()
	return
}

/*
*
func _duplicateHistoryItem() {

	appendHistory()
	history[historyPosition] = history[historyPosition-1]
	display()

}
*
*/
func saveHistory() {

	c := dsp.marshalJson()

	if historyPosition == len(history) {
		history = append(history, c)
	} else {
		history[historyPosition] = c
	}

	//cleanHistory()

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
	cleanHistory()
	historyPosition = len(history)

}

func backHistory() {

	saveHistory()
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

	saveHistory()

	historyPosition = historyPosition + 1

	if historyPosition > len(history) {
		historyPosition--
		return
	}

	if historyPosition == len(history) {
		reloadStash()
		historyPosition--
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
	if acceptInput {
		setAttributeByID("display", "class", "active")
	} else {
		setAttributeByID("display", "class", "inactive")
	}

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
		dsp.Input = nil
		dsp.setTitle("")
		saveState()
		stashState()
		display()
		return
	}

	h1 := history[:historyPosition]
	h2 := history[historyPosition+1:]

	history = nil

	history = append(history, h1...)
	history = append(history, h2...)

	js.Global().Get("localStorage").Call("setItem", "history", strings.Join(history, "\n"))
	/**
	if historyPosition == 0 {
		moveInHistory()
		return
	}
	**/
	saveState()
	stashState()
	moveInHistory()
}

func cleanHistory() {

	if historyPosition < len(history) {
		return
	}

	history = slicesCleanDuplicates(history)
	//	return
	var newhist []string
	dummy := new(console)
	dummy.Title = ""
	line := inputLine([]string{""})
	dummy.Input = []inputLine{line}

	for _, j := range history {
		json.Unmarshal([]byte(j), dummy)
		if dummy.Title != "" {
			newhist = append(newhist, j)
			continue
		}
		if len(dummy.Input) != 0 {
			newhist = append(newhist, j)
			continue
		}
	}

	history = nil
	history = append(history, newhist...)

	historyPosition = len(history)
	saveHistory()
}

func exportHistory() {

	obj := js.Global().Get("Blob").New([]any{strings.Join(history, "\n")})
	url := js.Global().Get("URL").Call("createObjectURL", obj).String()

	//obj2 := js.Global().Get("Blob").New([]any{historyToLatex()})
	//url2 := js.Global().Get("URL").Call("createObjectURL", obj2).String()

	stopInput()
	hide("console")
	hide("controls2")
	show("extra")
	hide("txtinput")
	show("historyDialog")
	show("backButton")
	hide("console")

	html := `<h3> Export History</h3>
<a href="` + url + `">right-click to download history as JSON</a>`

	/**+ "<br /><br />" + `<a href="` + url2 + `">right-click to download history as LaTeX</a></ br></ br>`**/

	html = html + `<h3>Import History</h3>
	Paste JSON into box.

<textarea name="textarea" id="historyinputarea" rows="15" cols="40"></textarea>
	 <button id="importHistory">Import</button>

`
	setTextByID("historyDialog", html)
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
	 <button id="importHistory">Import</button>

	<input type="file" id="inputfile />
	 `

	setTextByID("historyDialog", html)
}

func rewriteHistory() {
	stopInput()
	input := dom.GetWindow().Document().GetElementByID("historyinputarea").(*dom.HTMLTextAreaElement).Value()
	dom.GetWindow().Document().GetElementByID("historyinputarea").(*dom.HTMLTextAreaElement).SetValue("")

	if len(strings.TrimSpace(input)) > 0 {
		history = nil
		history = strings.Split(input, "\n")
	}
	historyPosition = len(history)
	cleanHistory()
	dsp.clear()
	display()
	printMessage("")
	hide("messages")
	hideExtra()
	show("console")
	stopInput()

}

func historyToLatex() string {
	var out string

	stashState()

	for p := range history {

		json.Unmarshal([]byte(history[p]), dsp)

		out = out + latexOutput()

	}

	reloadStash()

	return out
}

func clearHistory() {
	dsp.clear()
	history = nil
	stash = ""
	historyPosition = len(history)
	saveHistory()
	display()
	printMessage("")
	hide("messages")
	hideExtra()
	show("console")
	stopInput()
	return
}
