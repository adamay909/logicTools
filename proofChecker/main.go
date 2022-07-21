package main

import (
	"embed"
	_ "embed"
	"encoding/json"
	"sort"
	"strconv"
	"syscall/js"

	"github.com/adamay909/logicTools/gentzen"
	"honnef.co/go/js/dom/v2"
)

//go:embed assets/html/* assets/files assets/samples
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

	oExercises = false

	oClipboard = 2

	oAdvanced = false

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

	setupJS()
	resetDisplay()

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
	case "dummy":
		startInput()

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
	case "togglereadme":
		toggleReadme()

	case "loadExercise":
		toggleExercises()
	case "loadSamples":
		toggleSamples()
	case "reset":
		resetDisplay()

	case "toggleadvanced":
		toggleAdvanced()
	case "cliptype":
		toggleClipboardType()
	case "textInput":
		inputFromText()
	case "showCursorKeys":
		showArrowKeys()

	case "backButton":
		backToNormal()

	case "submitInput":
		getInput()

	case "nextExercise":
		nextProblem()
	case "prevExercise":
		prevProblem()
	case "quitExercise":
		endRandomExercise()

	case "arrowUp":
		handleInput("ArrowUp")
	case "arrowDown":
		handleInput("ArrowDown")
	case "arrowLeft":
		handleInput("ArrowLeft")
	case "arrowRight":
		handleInput("ArrowRight")
	case "delete":
		handleInput("Delete")

	default:
		if target.Get("className").String() == "fileLink" {
			loadFile(target.Get("innerHTML").String(), "exercises")
		}
		if target.Get("className").String() == "sampleLink" {
			loadFile(target.Get("innerHTML").String(), "samples")
		}
	}
}

func display() {
	dsp.format()
	setTextByID("display", dsp.typeset())
	setTextByID("dummy", dsp.typeset())
}

func displayDerivation() {
	dsp.formatDerivation()
	setTextByID("display", dsp.typeset())
	setTextByID("dummy", dsp.typeset())

}

func typeformula() {
	if !acceptInput {
		return
	}
	input := js.Global().Get("event").Get("key")

	handleInput(input.String())
}

func handleInput(s string) {

	dsp.handleInput(s)
	dsp.format()
	setTextByID("display", dsp.typeset())
	focusInput()
	setTextByID("dummy", dsp.typeset())
	scrollDisplay()
	return
}

func focusInput() {
	js.Global().Get("dummy").Call("focus")
}

func clearInput() {
	dsp.clear()
	setTextByID("setOffset", "First Line: "+strconv.Itoa(dsp.Offset))
	display()
	printMessage("")
	hide("messages")
	focusInput()
	stopInput()
}

func resetDisplay() {

	oPL = true
	oTHM = true
	oHELP = false
	oMENU = false
	oABOUT = false
	oExercises = false
	oClipboard = 2
	oAdvanced = false

	var cons console

	dsp = &cons

	setupPage()

	dsp.SystemPL = oPL
	dsp.Theorems = oTHM
	dsp.overhang = false
	dsp.Offset = 1
	dsp.viewTop = 0
	dsp.viewBottom = 20
	//finalize stuff
	display()
	toggleClipboardType()
	toggleTheorems()
	togglePL()
	toggleSettings()
	setDisplay()
	focusInput()
	exercises = nil
	expos = 0
}

func toggleTheorems() {
	stopInput()
	oTHM = !oTHM
	if oTHM {
		setTextByID("togglethm", "With Theorems")
	} else {
		setTextByID("togglethm", "No Theorems")
	}
	dsp.Theorems = oTHM
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
	dsp.SystemPL = oPL
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
	hide("console")
	show("extra")
	show("readme")
	show("backButton")
	return
}

func backToNormal() {
	stopInput()
	dom.GetWindow().Document().GetElementByID("textinputarea").(*dom.HTMLTextAreaElement).SetValue("")
	hide("backButton")
	hide("readme")
	hide("textinput")
	hide("exerciseList")
	hide("extra")
	show("console")
	show("dummy")
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
		setTextByID("cliptype", "Clipboard: LaTeX")
	case oJsonOutput:
		setTextByID("cliptype", "Clipboard: JSON")
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

	n, err := strconv.Atoi(js.Global().Call("prompt", "Number of first line", strconv.Itoa(dsp.Offset)).String())
	if err != nil {
		return
	}
	dsp.setOffset(n)
	setTextByID("setOffset", "First Line: "+strconv.Itoa(dsp.Offset))
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
		withTitle := true
		copyToClipboard(plainTextDeriv(withTitle))

	case oJsonOutput:
		copyToClipboard(marshalJson())
	}

	return
}

func startInput() {

	acceptInput = true
	display()
	setAttributeByID("cursor", "class", "active")
	setAttributeByID("display", "style", "border-color: blue")

}

func stopInput() {

	acceptInput = false
	display()
	setAttributeByID("cursor", "class", "inactive")
	setAttributeByID("display", "style", "border-color:lightgrey")

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
	hide("console")
	hide("dummy")
	show("extra")
	show("exerciseList")
	show("backButton")
	files, err := assets.ReadDir("assets/files")
	if err != nil {
		//fmt.Println(err)
		return
	}

	if len(files) == 0 {
		return
	}
	h := "<h3>Pick one to load</h3>"
	for _, e := range files {
		//		fmt.Println(e.Name())
		h = h + `<button class="fileLink" tabindex=0>` + e.Name() + `</button>`
	}
	h = h + `<p><button class="fileLink" id="nextExercise" tabindex=0>Random Exercises</button></p>`

	setTextByID("exerciseList", h)
}

func toggleSamples() {
	stopInput()
	hide("console")
	show("extra")
	show("exerciseList")
	show("backButton")
	files, err := assets.ReadDir("assets/samples")
	if err != nil {
		//		fmt.Println(err)
		return
	}

	if len(files) == 0 {
		return
	}

	order := func(i, j int) bool {
		return sort.StringsAreSorted([]string{files[i].Name(), files[j].Name()})
	}

	sort.Slice(files, order)

	h := "<h3>Pick one to load</h3>"
	for _, e := range files {
		//fmt.Println(e.Name())
		h = h + `<button class="sampleLink" tabindex=0>` + e.Name() + `</button>`
	}
	setTextByID("exerciseList", h)
}

func loadFile(name string, t string) {
	stopInput()
	if t == "exercises" {
		name = "assets/files/" + name
	}
	if t == "samples" {
		name = "assets/samples/" + name
	}

	d, err := assets.ReadFile(name)

	if err != nil {
		panic(err)
	}

	dsp.clear()

	json.Unmarshal(d, dsp)
	dsp.xpos, dsp.ypos = 0, 0
	dsp.overhang = false
	dsp.modifier = ""
	if dsp.SystemPL != oPL {
		togglePL()
	}
	if dsp.Theorems != oTHM {
		toggleTheorems()
	}
	show("console")
	show("dummy")
	hide("backButton")
	hide("exerciseList")
	hide("extra")
	printMessage("")
	hide("messages")
	display()
	stopInput()
}

func inputFromText() {

	stopInput()
	hide("console")
	show("extra")
	show("textinput")
	show("backButton")

}

func getInput() {
	stopInput()
	s := dom.GetWindow().Document().GetElementByID("textinputarea").(*dom.HTMLTextAreaElement).Value()
	lines, title, err := text2data(s)
	if err != nil {
		js.Global().Call("alert", err.Error())
		return
	}
	dom.GetWindow().Document().GetElementByID("textinputarea").(*dom.HTMLTextAreaElement).SetValue("")
	dsp.clear()
	dsp.Title = title
	dsp.Input = lines
	display()
	printMessage("")
	hide("messages")
	hide("textinput")
	hide("extra")
	hide("backButton")
	show("console")
	stopInput()

}

func toggleAdvanced() {

	stopInput()
	oAdvanced = !oAdvanced

	if oAdvanced {
		show("advancedstuff")
	} else {
		hide("advancedstuff")
	}
	return
}

func showArrowKeys() {
	show("cursorControls")
}
