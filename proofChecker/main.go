package main

import (
	"embed"
	_ "embed"
	"encoding/json"
	"fmt"
	"strconv"
	"syscall/js"

	"github.com/adamay909/logicTools/gentzen"
	"honnef.co/go/js/dom/v2"
)

//go:embed assets/html/* assets/files
var assets embed.FS

//Enable some features for personal teaching material.
//Not useful for general consumption.
var oPRIVATE = false

var indexHtml, helpHtml, styleCSS string

var (
	oPL = true

	oTHM = true

	oHELP = false

	oMENU = false

	oABOUT = false

	oClipboard = 2

	logConstBindings [][3]string

	acceptInput = true
)

const (
	oLatexOutput = 1
	oTextOutput  = 0
	oJsonOutput  = 2
)

var dsp *console

func main() {

	setupPage()
	setupJS()

	dsp = new(console)

	//finalize stuff
	dsp.clear()
	display()
	toggleClipboardType()
	toggleTheorems()
	togglePL()
	toggleSettings()
	setDisplay()
	focusInput()
	<-make(chan bool)
}

func setupPage() {

	//load styles
	d, _ := assets.ReadFile("assets/html/main.css")

	dom.GetWindow().Document().GetElementsByTagName("style")[0].SetInnerHTML(string(d))

	//Populate the page
	d, _ = assets.ReadFile("assets/html/body.html")

	dom.GetWindow().Document().GetElementsByTagName("body")[0].SetInnerHTML(string(d))

	d, _ = assets.ReadFile("assets/html/help.html")

	setTextByID("help", string(d))
	dom.GetWindow().Document().GetElementByID("help").SetAttribute("style", "display: none")

	d, _ = assets.ReadFile("assets/html/README.html")

	setTextByID("readme", string(d))
	dom.GetWindow().Document().GetElementByID("readme").SetAttribute("style", "display: none")

	if !oPRIVATE {
		setAttributeByID("loadExercise", "style", "display:none")
	}
}

func setupJS() {

	js.Global().Call("addEventListener", "keydown", js.FuncOf(jsWrap(typeformula)).Value, true)
	js.Global().Call("addEventListener", "click", js.FuncOf(jsWrap(onClick)).Value, true)
}

func onClick() {

	target := js.Global().Get("event").Get("target")
	//	fmt.Println(target.Get("id"))
	//	fmt.Println(target.Get("outerHTML"))
	switch target.Get("id").String() {
	case "toggleSettings":
		toggleSettings()
	case "check":
		checkDeriv()
	case "clearInput":
		clearInput()
	case "toClipboard":
		toClipboard()
	case "toggleHelp":
		toggleHelp()
	case "toggleSystem":
		togglePL()
	case "setOffset":
		setOffset()
	case "togglethm":
		toggleTheorems()
	case "cliptype":
		toggleClipboardType()
	case "togglereadme":
		toggleReadme()
	case "dummy":
		startInput()
	case "loadExercise":
		toggleExercises()
	default:
		if target.Get("className").String() == "fileLink" {
			loadFile(target.Get("innerHTML").String())
		}
	}
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
	setTextByID("setOffset", "First Line: "+strconv.Itoa(dsp.offset))
	display()
	printMessage("")
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
		setTextByID("toggleSystem", "Predicate Logic")
	} else {
		logConstBindings = connBindings
		setTextByID("toggleSystem", "Sentential Logic")
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
	if oPRIVATE {
		oClipboard = (oClipboard + 1) % 3
	} else {
		oClipboard = (oClipboard + 1) % 2
	}
	switch oClipboard {
	case oTextOutput:
		setTextByID("cliptype", "Clipboard: text")
	case oLatexOutput:
		setTextByID("cliptype", "Clipboard: Latex")
	case oJsonOutput:
		setTextByID("cliptype", "Clipboard: json")
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
	setTextByID("setOffset", "First Line: "+strconv.Itoa(dsp.offset))
	display()
}

func toClipboard() {
	if oABOUT {
		return
	}
	stopInput()
	switch oClipboard {

	case oLatexOutput:
		copyToClipboard(latexOutput())

	case oTextOutput:
		copyToClipboard(plainTextDeriv())

	case oJsonOutput:
		copyToClipboard(marshalJson())
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

func toggleExercises() {
	stopInput()
	hide("controls")
	hide("editor")
	show("exerciseList")
	files, err := assets.ReadDir("assets/files")
	if err != nil {
		fmt.Println(err)
		return
	}

	if len(files) == 0 {
		return
	}
	h := "<h3>Pick one to load</h3>"
	for _, e := range files {
		fmt.Println(e.Name())
		h = h + `<div class="fileLink">` + e.Name() + `</div>`
	}
	setTextByID("exerciseList", h)
}

func loadFile(name string) {
	stopInput()
	name = "assets/files/" + name

	d, err := assets.ReadFile(name)

	if err != nil {
		panic(err)
	}

	dsp.clear()

	json.Unmarshal(d, dsp)
	dsp.xpos, dsp.ypos = 0, 0
	dsp.overhang = false
	dsp.modifier = ""
	show("controls")
	show("editor")
	hide("exerciseList")
	display()
	stopInput()
}
