package main

import (
	"syscall/js"
)

var avail [][3]string

func editTitle() {

	stopInput()

	if dsp.Title == "" {
		dsp.Title = "&ensp;"
	}
	stopInput()

	setAttributeByID("extitle", "style", "background-color: lightblue")

	js.Global().Call("removeEventListener", "keydown", mainEditorFunc, true)

	js.Global().Call("addEventListener", "keydown", titleEditorFunc, true)

	avail = combineBindings(keyBindings, punctBindings, connBindings, plBindings, turnstileBindings, greekBindings, extraBindings)
}

func endEditTitle() {

	if dsp.Title == "&ensp;" {
		dsp.Title = ""
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
	stopInput()
	titleString = []rune(dsp.Title)

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
		if string(titleString) == "&ensp;" {
			setAttributeByID("extitle", "style", "background-color: lightblue")
			return
		}
		titleString = titleString[:len(titleString)-1]
		dsp.Title = string(titleString)
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
		dsp.Title = string(titleString)
		stopInput()
		setAttributeByID("extitle", "style", "background-color: lightblue")
	}
	tmodifier = ""

}

func getInputString() string {

	return js.Global().Get("event").Get("key").String()

}
