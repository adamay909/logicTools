package main

import (
	"fmt"
	"strings"
	"syscall/js"
)

var avail [][3]string

func editTitle() {
	var cursor = `<div id="cursor2">&thinsp;</div>`

	dsp.Title = strings.TrimSuffix(dsp.Title, cursor)

	dsp.Title = dsp.Title + cursor

	stopInput()

	setAttributeByID("extitle", "style", "background-color: lightblue")

	js.Global().Call("removeEventListener", "keydown", mainEditorFunc, true)

	js.Global().Call("addEventListener", "keydown", titleEditorFunc, true)

	avail = combineBindings(keyBindings, punctBindings, connBindings, plBindings, turnstileBindings, greekBindings, extraBindings)
	fmt.Println("ready")
	focusInput()
}

func endEditTitle() {
	var cursor = `<div id="cursor2">&thinsp;</div>`
	dsp.Title = strings.TrimSuffix(dsp.Title, cursor)
	if historyPosition != len(history) {
		history[historyPosition] = dsp.marshalJson()
	}

	setAttributeByID("extitle", "style", "background-color: white")
	js.Global().Call("removeEventListener", "keydown", titleEditorFunc, true)
	js.Global().Call("addEventListener", "keydown", mainEditorFunc, true)
	stopInput()
}

var tpos int
var tmodifier string
var titleString []rune

func typetitle() {
	var cursor = `<div id="cursor2">&thinsp;</div>`
	fmt.Println("editing title")
	titleString = []rune(strings.TrimSuffix(dsp.Title, cursor))

	key := getInputString()
	if key == "Shift" {
		return
	}

	if key == enter {
		tmodifier = ""
		endEditTitle()
		return
	}

	if key == backspace {
		tmodifier = ""
		if string(titleString) == "" {
			setAttributeByID("extitle", "style", "background-color: lightblue")
			return
		}
		titleString = titleString[:len(titleString)-1]
		dsp.Title = string(titleString) + cursor
		stopInput()
		setAttributeByID("extitle", "style", "background-color: lightblue")
		return
	}

	if key == `\` || key == `|` {
		tmodifier = key
		return
	}

	if tmodifier != "" {
		key = tmodifier + key
	}

	inp, err := tkOf(key, tkraw, tktxt, avail)
	if err == nil {
		titleString = append(titleString, []rune(inp)...)
		dsp.Title = string(titleString) + cursor
		stopInput()
		setAttributeByID("extitle", "style", "background-color: lightblue")
	}
	tmodifier = ""
}

func getInputString() string {

	return js.Global().Get("event").Get("key").String()

}
