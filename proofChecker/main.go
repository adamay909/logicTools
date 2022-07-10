package main

import (
	"embed"
	_ "embed"
	"strconv"
	"syscall/js"

	"github.com/adamay909/logicTools/gentzen"
	"honnef.co/go/js/dom/v2"
)

//go:embed assets/html/*
var assets embed.FS

//Enable some features for personal teaching material.
//Not useful for general consumption.
var oPRIVATE = true

var indexHtml, helpHtml, styleCSS string

var (
	oPL = true

	oTHM = true

	oHELP = false

	oMENU = false

	oABOUT = false

	oLatexOutput = false

	logConstBindings [][3]string

	acceptInput = true
)

var dsp *console

func main() {
	dsp = new(console)

	//load styles
	d, _ := assets.ReadFile("assets/html/main.css")

	dom.GetWindow().Document().GetElementsByTagName("style")[0].SetInnerHTML(string(d))

	//Populate the page
	d, _ = assets.ReadFile("assets/html/index.html")

	dom.GetWindow().Document().GetElementsByTagName("body")[0].SetInnerHTML(string(d))

	d, _ = assets.ReadFile("assets/html/help.html")

	setTextByID("help", string(d))
	dom.GetWindow().Document().GetElementByID("help").SetAttribute("style", "display: none")

	d, _ = assets.ReadFile("assets/html/README.html")

	setTextByID("readme", string(d))
	dom.GetWindow().Document().GetElementByID("readme").SetAttribute("style", "display: none")

	//	dom.GetWindow().Document().GetElementByID("display").SetAttribute("style", "counter-reset: line "+strconv.Itoa(0)+";")

	//setup JS stuff
	js.Global().Set("toggleSettings", js.FuncOf(jsWrap(toggleSettings)).Value)
	js.Global().Set("clearInput", js.FuncOf(jsWrap(clearInput)).Value)
	js.Global().Set("checkDerivation", js.FuncOf(jsWrap(checkDeriv)).Value)
	js.Global().Set("toggleSystem", js.FuncOf(jsWrap(togglePL)).Value)
	js.Global().Set("toggleTheorems", js.FuncOf(jsWrap(toggleTheorems)).Value)
	js.Global().Set("activateInput", js.FuncOf(jsWrap(startInput)).Value)
	js.Global().Set("toggleHelp", js.FuncOf(jsWrap(toggleHelp)).Value)
	js.Global().Set("toggleReadme", js.FuncOf(jsWrap(toggleReadme)).Value)
	js.Global().Set("setOffset", js.FuncOf(jsWrap(setOffset)).Value)
	js.Global().Set("toClipboard", js.FuncOf(jsWrap(toClipboard)).Value)
	js.Global().Set("toggleClipboardType", js.FuncOf(jsWrap(toggleClipboardType)).Value)

	js.Global().Call("addEventListener", "keydown", js.FuncOf(jsWrap(typeformula)).Value, true)
	//js.Global().Call("addEventListener", "input", js.FuncOf(jsWrap(typeformula)).Value, true)

	//finalize stuff
	dsp.clear()
	display()
	oPL = true
	toggleTheorems()
	togglePL()
	toggleSettings()
	setDisplay()
	focusInput()
	<-make(chan bool)
}

func display() {
	dsp.format()
	setTextByID("display", dsp.typeset())
}

func displayDerivation() {
	dsp.formatDerivation()
	setTextByID("display", dsp.typeset())
}

func typeformula() {
	if !acceptInput {
		return
	}
	o := js.Global().Get("event").Get("key")

	dsp.handleInput(o.String())
	dsp.format()
	setTextByID("display", dsp.typeset())
	focusInput()
	setTextByID("dummy", dsp.typeset())
	return
}

func focusInput() {
	js.Global().Get("dummy").Call("focus")
}

func clearInput() {
	dsp.clear()
	setTextByID("setoffset", "First Line: "+strconv.Itoa(dsp.offset))
	display()
	focusInput()
	stopInput()
}

func toggleTheorems() {
	stopInput()
	oTHM = !oTHM
	if oTHM {
		setTextByID("togglethm", "With Theorems")
	} else {
		setTextByID("togglethm", "No Theorems")
	}
	gentzen.SetAllowTheorems(oTHM)
	return
}

func togglePL() {
	stopInput()
	oPL = !oPL
	if oPL {
		logConstBindings = append(connBindings, plBindings...)
		setTextByID("toggle", "Predicate Logic")
	} else {
		logConstBindings = connBindings
		setTextByID("toggle", "Sentential Logic")
	}
	gentzen.SetPL(oPL)
	return
}

func toggleHelp() {
	if oABOUT {
		return
	}

	oHELP = !oHELP
	stopInput()
	if !oHELP {
		setTextByID("toggleHelp", "Show Help")
		hide("help")
	} else {
		setTextByID("toggleHelp", "Hide Help")
		show("help")
	}
	setDisplay()
	return
}

func toggleSettings() {
	if oABOUT {
		toggleReadme()
	}
	stopInput()
	oMENU = !oMENU
	if !oMENU {
		hide("settingsMenu")
	} else {
		show("settingsMenu")
	}
	setDisplay()
	return
}

func toggleReadme() {
	stopInput()
	oABOUT = !oABOUT
	if oABOUT {
		hide("controls")
		hide("editor")
		hide("help")
		show("readme")
		return
	}
	show("editor")
	show("controls")
	hide("readme")
	if oHELP {
		show("help")
	}
	setDisplay()
}

func toggleClipboardType() {
	stopInput()
	oLatexOutput = !oLatexOutput
	if oLatexOutput {
		setTextByID("cliptype", "Clipboard: Latex")
	} else {
		setTextByID("cliptype", "Clipboard: text")
	}
	return
}

func setDisplay() {

	if oABOUT {
		setAttributeByID("mainArea", "style", "grid-template-columns: 1fr 10fr")
		return
	}

	if oMENU && oHELP {
		setAttributeByID("mainArea", "style", "grid-template-columns: 1fr 6fr 4fr")
		return
	}

	if oMENU && !oHELP {
		setAttributeByID("mainArea", "style", "grid-template-columns: 1fr 10fr")
		return
	}

	if !oMENU && oHELP {
		setAttributeByID("mainArea", "style", "grid-template-columns: 6fr 4fr")
		return
	}

	setAttributeByID("mainArea", "style", "grid-template-columns: 100%")
	return
}

func checkDeriv() {
	if oABOUT {
		return
	}
	stopInput()
	checkDerivation()
	return
}

func setOffset() {

	n, err := strconv.Atoi(js.Global().Call("prompt", "Number of first line", strconv.Itoa(dsp.offset)).String())
	if err != nil {
		return
	}
	dsp.setOffset(n)
	setTextByID("setoffset", "First Line: "+strconv.Itoa(dsp.offset))
	display()
}

func toClipboard() {
	if oABOUT {
		return
	}
	stopInput()
	if oLatexOutput {
		copyToClipboard(latexOutput())
	} else {
		copyToClipboard(plainTextDeriv())
	}
	return
}

func startInput() {

	acceptInput = true
	dom.GetWindow().Document().GetElementByID("cursor").SetAttribute("class", "active")
	display()

}

func stopInput() {

	acceptInput = false
	dom.GetWindow().Document().GetElementByID("cursor").SetAttribute("class", "inactive")

}

func setTextByID(elem string, content string) {
	dom.GetWindow().Document().GetElementByID(elem).SetInnerHTML(content)
}

func setAttributeByID(elem string, attrName, attrCont string) {
	dom.GetWindow().Document().GetElementByID(elem).SetAttribute(attrName, attrCont)
	return
}

func jsWrap(f func()) (fn func(this js.Value, args []js.Value) any) {

	fn = func(this js.Value, args []js.Value) any {
		f()
		return true
	}

	return fn
}

func show(elem string) {
	setAttributeByID(elem, "style", "display: inline-block")
}

func hide(elem string) {
	setAttributeByID(elem, "style", "display: none")
}

func copyToClipboard(s string) {
	js.Global().Get("navigator").Get("clipboard").Call("writeText", s)
	return
}
