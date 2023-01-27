package main

import (
	"embed"
	_ "embed"
	"encoding/json"
	"os"
	"sort"
	"strconv"
	"strings"
	"syscall/js"

	"github.com/adamay909/logicTools/gentzen"
	"honnef.co/go/js/dom/v2"
)

// set to true for debug log to stdout
var oDEBUG = false

//go:embed assets/html/* assets/files assets/samples
var assets embed.FS

// Enable some features for personal teaching material.
// Not useful for general consumption.

var oPRIVATE = true

var mainEditorFunc,
	titleEditorFunc,
	clickFunc,
	loadFunc js.Value

var indexHtml, helpHtml, styleCSS string

var (
	oPL = true

	oML = false

	oDR = false

	oTHM = true

	oHELP = false

	oMENU = true

	oABOUT = false

	oEXTHM = false

	oExercises = false

	oClipboard = 0

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

	initMessages()
	gentzen.SetStandardPolish(false)
	setupJS()
	resetDisplay()
	hideExtra()
	loadHistory()
	stopInput()

	<-make(chan bool)
}

func initMessages() {
	toggleDebug()
	debug("You can toggle verbose logging with CTRL-ALT-v in the editor.")
	toggleDebug()
	return
}

func toggleDebug() {
	oDEBUG = !oDEBUG
	gentzen.SetDebug(oDEBUG)
	if oDEBUG {
		gentzen.SetDebuglog(os.Stdout)
	}
	return
}

func debug(m ...any) {
	if !oDEBUG {
		return
	}
	dm := []any{"PC: "}
	m = append(dm, m...)
	gentzen.Debug(m...)
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

	d, _ = assets.ReadFile("assets/html/version")
	setTextByID("versionnumber", "v"+string(d)+"&emsp;&emsp;")

	if !oPRIVATE {
		setAttributeByID("loadExercise", "style", "display:none")
	}

	h := int(0.8 * float64(js.Global().Get("innerHeight").Int()))
	setAttributeByID("display", `style`, `height: `+strconv.Itoa(h)+`px`)

	dsp.fontSize = 120
	setAttributeByID("editor", `style`, `font-size:`+strconv.Itoa(dsp.fontSize)+`%;`)

	//Stuff for simplifying display
	hide("menuButton")
	hide("print")
}

func setupJS() {

	mainEditorFunc = js.FuncOf(jsWrap(typeformula)).Value
	titleEditorFunc = js.FuncOf(jsWrap(typetitle)).Value
	clickFunc = js.FuncOf(jsWrap(onClick)).Value
	//	loadFunc = js.FuncOf(jsWrap(readHistoryFromFile)).Value

	js.Global().Call("addEventListener", "keydown", mainEditorFunc, true)
	js.Global().Call("addEventListener", "click", clickFunc, true)
	// js.Global().Call("addEventListener", "change", loadFunc, true)
}

func onClick() {

	target := js.Global().Get("event").Get("target")
	//	fmt.Println(target.Get("id"))
	//	fmt.Println(target.Get("outerHTML"))
	switch target.Get("id").String() {
	case "dummy":
		endEditTitle()
		startInput()
	case "title":
		editTitle()

	case "toggleSettings":
		toggleSettings()
	case "check":
		checkDeriv()
	case "clearInput":
		clearInput()
	case "toClipboard":
		toClipboard()
	case "toLatex":
		toClipboardLatex()
	case "print":
		toPrinter()
	case "toggleHelp":
		toggleHelp()

	case "setTitle":
		editTitle()

	case "toggleSystem":
		togglePL()
	case "setOffset":
		setOffset()
	case "toggleDR":
		toggleDR()
	case "toggleML":
		toggleML()
	case "togglereadme":
		toggleReadme()

	case "imp-export":
		exportHistory()

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
	case "backHistory":
		backHistory()
	case "forwardHistory":
		forwardHistory()
	case "removeFromHistory":
		rmFromHistory()
	case "exportHistory":
		exportHistory()
	case "importHistory":
		rewriteHistory()
	case "duplicateScreen":
		duplicateHistoryItem()
	case "deleteAllHistory":
		clearHistory()

	case "backButton":
		backToNormal()

	case "submitInput":
		getInput()

	case "randomTheorem":
		oEXTHM = true
		nextProblem()

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

	case "sizeUp":
		sizeUp()
	case "sizeDown":
		sizeDown()

	case "loadExamples":
		loadSamples()

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
	setTextByID("dummy", `<h3 id="title">`+prettyGreek(dsp.Title)+`</h3>`)
	setTextByID("pagenumber", "page: "+strconv.Itoa(historyPosition+1)+"/"+strconv.Itoa(len(history)))

	// setTextByID("dummy", "")
}

func displayDerivation() {
	dsp.formatDerivation()
	setTextByID("display", dsp.typeset())
	setTextByID("dummy", "")

}

func typeformula() {
	if !acceptInput {
		return
	}
	input := js.Global().Get("event").Get("key")

	handleInput(input.String())
	saveState()
}

func handleInput(s string) {

	dsp.handleInput(s)
	dsp.format()
	setTextByID("display", dsp.typeset())
	focusInput()
	setTextByID("dummy", "")
	scrollDisplay()
	return
}

func focusInput() {
	js.Global().Get("dummy").Call("focus")
}

func clearInput() {
	saveHistory()
	insertEmptyHistoryItem()
	clearScreen()
	saveState()
	focusInput()
	stopInput()
}

func clearScreen() {
	dsp.clear()
	setTextByID("setOffset", "First Line: "+strconv.Itoa(dsp.Offset))
	display()
	printMessage("")
	hide("messages")
}

func resetDisplay() {

	oEXTHM = false
	oPL = true
	oDR = false
	oML = false
	oHELP = false
	oMENU = true
	oABOUT = false
	oExercises = false
	oClipboard = 0
	oAdvanced = false

	var cons console

	dsp = &cons

	setupPage()

	dsp.SystemPL = oPL
	dsp.SystemML = oML
	dsp.Theorems = oTHM
	dsp.overhang = false
	dsp.Offset = 1
	dsp.viewTop = 0
	dsp.viewBottom = 20
	//finalize stuff
	display()
	toggleClipboardType()
	toggleTheorems()
	toggleDR()
	toggleML()
	togglePL()
	toggleSettings()
	toggleMenuButton()
	setDisplay()
	focusInput()
	exercises = nil
	expos = 0
}

func toggleTheorems() {
	stopInput()
	oTHM = true
	/**
	oTHM = !oTHM
	if oTHM {
		setTextByID("togglethm", "With Theorems")
	} else {
		setTextByID("togglethm", "No Theorems")
	}**/
	dsp.Theorems = oTHM
	gentzen.SetAllowTheorems(oTHM)
	return
}

func togglePL() {
	stopInput()
	oPL = !oPL
	logConstBindings = nil
	if oPL {
		logConstBindings = append(connBindings, plBindings...)
		setTextByID("toggleSystem", "Predicate Logic")
	} else {
		logConstBindings = connBindings
		setTextByID("toggleSystem", "Sentential Logic")
	}
	if oML {
		logConstBindings = append(logConstBindings, mlBindings...)
	}
	dsp.SystemPL = oPL
	gentzen.SetPL(oPL)
	return
}

func toggleDR() {
	stopInput()
	oDR = !oDR
	if oDR {
		setTextByID("toggleDR", `Derived Rules&emsp; &#x2713;`)
	} else {
		setTextByID("toggleDR", `Derived Rules`)
	}
	dsp.DerivedRules = oDR
	gentzen.SetDR(oDR)
}
func toggleML() {
	stopInput()
	oML = !oML
	if oML {
		setTextByID("toggleML", `Modal Logic&emsp; &#x2713;`)
	} else {
		setTextByID("toggleML", `Modal Logic`)
	}
	dsp.SystemML = oML
	gentzen.SetML(oML)
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
	hideExtra()
	show("console")
	show("dummy")
}

func toggleClipboardType() {
	stopInput()
	oClipboard = oTextOutput
	return

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

func setTitle() {

	title := js.Global().Call("prompt", "Title:").String()
	dsp.Title = convert(title)

	display()
}

func convert(s string) string {

	words := strings.Split(s, " ")

	var wn []string
	for _, w := range words {
		r := ""
		t := w
		for i := 0; i < len(t); {
			for _, e := range allBindings {
				if strings.HasPrefix(string(t[i:]), e[tkraw]) {
					r = r + e[tktxt]
					i = i + len(e[tkraw]) - 1
					break
				}
			}
			i++
		}
		wn = append(wn, r)
	}

	return strings.Join(wn, " ")

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
		copyToClipboard(dsp.marshalJson())
	}

	return
}
func toClipboardLatex() {
	if oABOUT {
		return
	}
	stopInput()

	copyToClipboard(latexOutput())

	return
}
func startInput() {

	acceptInput = true
	display()
	setAttributeByID("cursor", "class", "active")
	setAttributeByID("display", "class", "active")

}

func stopInput() {

	acceptInput = false
	display()
	setAttributeByID("cursor", "class", "inactive")
	setAttributeByID("display", "class", "inactive")

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
	//removeClass(elem, "hide")
	//addClass(elem, "show")
}

func hide(elem string) {
	setAttributeByID(elem, "style", "display: none")
	// addClass(elem, "hide")
}

func addClass(elem string, nc string) {
	class := dom.GetWindow().Document().GetElementByID(elem).GetAttribute("class")
	class = class + " " + nc
	dom.GetWindow().Document().GetElementByID(elem).SetAttribute("class", class)
}

func removeClass(elem string, nc string) {
	class := dom.GetWindow().Document().GetElementByID(elem).GetAttribute("class")
	ic := strings.Split(class, " ")
	class = ""
	for _, c := range ic {
		if c == nc {
			continue
		}
		class = class + c
	}
	dom.GetWindow().Document().GetElementByID(elem).SetAttribute("class", class)
}

func copyToClipboard(s string) {
	js.Global().Get("navigator").Get("clipboard").Call("writeText", s)
	return
}

func hideExtra() {

	hide("backButton")
	hide("txtinput")
	hide("exerciseList")
	hide("historyDialog")
	hide("readme")
	hide("extra")
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
		h = h + `<button class="fileLink" tabindex=0>` + e.Name() + `</button>`
	}
	h = h + `<p><button class="fileLink" id="nextExercise" tabindex=0>Random Exercises (Derivations)</button></p>`
	h = h + `<p><button class="fileLink" id="randomTheorem" tabindex=0>Random Exercises (Theorems)</button></p>`

	setTextByID("exerciseList", h)
}

func toggleSamples() {
	stopInput()
	hide("console")
	hide("controls2")
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

func loadSamples() {
	stopInput()
	d, err := assets.ReadFile("assets/samples/samples.txt")
	if err != nil {
		panic(err)
	}
	json := string(d)
	history = strings.Split(json, "\n")
	historyPosition = 0
	moveInHistory()
	display()

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
	show("txtinput")
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

func sizeUp() {
	dsp.fontSize = dsp.fontSize + 20
	setAttributeByID("editor", `style`, `font-size:`+strconv.Itoa(dsp.fontSize)+`%;`)
}

func sizeDown() {
	dsp.fontSize = dsp.fontSize - 20
	setAttributeByID("editor", `style`, `font-size:`+strconv.Itoa(dsp.fontSize)+`%;`)
}

func toPrinter() {

	checkDeriv()

	stashScreen()

	message := dom.GetWindow().Document().GetElementByID("messages").OuterHTML()

	dom.GetWindow().Document().GetElementsByTagName("body")[0].SetInnerHTML(dsp.typeset() + message)

	js.Global().Call("print", "")

	restoreScreen()
}

var screenStash string

func stashScreen() {

	screenStash = dom.GetWindow().Document().GetElementsByTagName("body")[0].InnerHTML()

}

func restoreScreen() {

	dom.GetWindow().Document().GetElementsByTagName("body")[0].SetInnerHTML(screenStash)

}
