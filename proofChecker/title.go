package main

import (
	"fmt"
	"strings"
	"syscall/js"
)

var avail [][3]string

var tCursorPos int

var tcursor = `<div id="cursor2">&thinsp;</div>`

func editTitle() {

	//	dsp.Title = strings.TrimSuffix(dsp.Title, cursor)

	dsp.Title = rmString(dsp.Title, tcursor)
	tCursorPos = len([]rune(dsp.Title))

	dsp.Title = dsp.Title + tcursor

	stopInput()

	js.Global().Call("removeEventListener", "keydown", mainEditorFunc, true)

	js.Global().Call("addEventListener", "keydown", titleEditorFunc, true)

	avail = combineBindings(keyBindings, punctBindings, connBindings, plBindings, turnstileBindings, greekBindings, extraBindings)
	focusInput()
}

func endEditTitle() {
	dsp.Title = rmString(dsp.Title, tcursor)
	if historyPosition != len(history) {
		history[historyPosition] = dsp.marshalJson()
	}

	js.Global().Call("removeEventListener", "keydown", titleEditorFunc, true)
	js.Global().Call("addEventListener", "keydown", mainEditorFunc, true)
	stopInput()
}

var tpos int
var tmodifier string
var titleString []rune

func typetitle() {
	typetitleNew()
}

func typetitleNew() {
	dsp.Title = rmString(dsp.Title, tcursor)
	titleString = []rune(dsp.Title)
	fmt.Println(titleString)
	key := getInputString()
	if key == "Shift" {
		return
	}

	if key == enter {
		tmodifier = ""
		endEditTitle()
		return
	}

	if key == left {
		if tCursorPos > 0 {
			tCursorPos--
		}

		dsp.Title = insertString(dsp.Title, tcursor, tCursorPos)
		stopInput()
		return
	}

	if key == right {
		if tCursorPos < len(titleString) {
			tCursorPos++
		}

		dsp.Title = insertString(dsp.Title, tcursor, tCursorPos)
		stopInput()
		return
	}

	if key == backspace {
		tmodifier = ""
		if len(titleString) == 0 {
			dsp.Title = tcursor

			return
		}

		if tCursorPos == 0 {
			dsp.Title = tcursor + dsp.Title
			return
		}
		var nts []rune
		nts = titleString[:tCursorPos-1]
		if tCursorPos < len(titleString) {
			nts = append(nts, titleString[tCursorPos:]...)
		}
		dsp.Title = string(nts)
		tCursorPos--
		dsp.Title = insertString(dsp.Title, tcursor, tCursorPos)
		stopInput()
		return
	}

	if key == `\` || key == `|` {
		tmodifier = key
		dsp.Title = insertString(dsp.Title, tcursor, tCursorPos)
		return
	}

	if tmodifier != "" {
		key = tmodifier + key
	}

	inp, err := tkOf(key, tkraw, tktxt, avail)
	if err == nil {
		dsp.Title = insertString(dsp.Title, inp, tCursorPos)
		tCursorPos++
	}
	tmodifier = ""
	dsp.Title = insertString(dsp.Title, tcursor, tCursorPos)
	stopInput()
}

func typetitleOld() {
	titleString = []rune(strings.TrimSuffix(dsp.Title, tcursor))
	dsp.Title = string(titleString) + tcursor

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
			return
		}
		titleString = titleString[:len(titleString)-1]
		dsp.Title = string(titleString) + tcursor
		stopInput()
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
		dsp.Title = string(titleString) + tcursor
		stopInput()
	}
	tmodifier = ""
}
func getInputString() string {

	return js.Global().Get("event").Get("key").String()

}
