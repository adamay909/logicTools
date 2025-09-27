package main

import (
	"embed"
	_ "embed"
	"os"
	"strconv"
	"strings"
	"syscall/js"

	"github.com/adamay909/logicTools/gentzen"
	"github.com/adamay909/logicTools/proofChecker/internal/editor"
	"honnef.co/go/js/dom/v2"
)

// set to true for debug log to stdout

//go:embed assets/html/*
var assets embed.FS

// Enable some features for personal teaching material.
// Not useful for general consumption.

var oPRIVATE = true

var inputFunc,
	mainEditorFunc,
	titleEditorFunc,
	clickFunc,
	snapshotFunc,
	loadFunc js.Value

var indexHtml, helpHtml, styleCSS string

var (
	oPL              = false
	oML              = false
	oDR              = false
	oTHM             = false
	oHELP            = false
	oMENU            = false
	oABOUT           = false
	oEXTHM           = false
	oDEBUG           = false
	oExercises       = false
	oClipboard       = 0
	oAdvanced        = false
	logConstBindings [][3]string
	acceptInput      = true
)

const (
	oLatexOutput = 1
	oTextOutput  = 0
	oJsonOutput  = 2
)

var dsp *console

func main() {
	dsp = new(console)
	initMessages()
	gentzen.SetSpecialConn(true)
	setBasicInferenceRules()
	setupPage()
	setupJS()

	//**********
	//clearHistory()
	//**********
	loadHistory()
	updatePageNumber()
	//writeStateToHTML()
	//	loadHistory()
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

	d, _ = assets.ReadFile("assets/html/version")
	setTextByID("versionnumber", "v"+string(d)+"&emsp;&emsp;")

}

func setupJS() {
	clickFunc = js.FuncOf(jsWrap(onClick)).Value
	js.Global().Call("addEventListener", "click", clickFunc, true)
	dsp.editor = editor.NewDerivationEditor("editor")
	dsp.title = editor.NewSimpleEditor("title")
	addEventListener(domDocument.Call("querySelector", "#editorWindow"), "editorinput", jsFuncSaveSnapshot)
}

func onClick() {

	target := js.Global().Get("event").Get("target")
	switch target.Get("id").String() {

	case "toggleSettings":
		toggleSettingsMenu()
	case "check":
		checkDeriv()
	case "toClipboard":
		toClipboard()
	case "toLatex":
		toClipboardLatex()
	case "printTree":
		//printTree()
	case "toggleHelp":
		toggleHelp()

	case "toggleSystem":
		togglePL()
	case "setOffset":
		setOffset()
	case "toggleTheorem":
		toggleTheorems()
	case "newPage":
		newpage()
	case "backHistory":
		moveBackInHistory()
	case "forwardHistory":
		moveForwardInHistory()
	case "removeFromHistory":
		//		rmFromHistory()
	case "duplicateScreen":
		//		duplicateHistoryItem()
	case "clearHistory":
		clearHistory()
	case "backButton":
		toggleReadme()

	case "sizeUp":
		sizeUp()
	case "sizeDown":
		sizeDown()
	}
	if target.Get("classList").Call("contains", "togglereadme").Bool() {
		toggleReadme()
	}
	writeStateToHTML()
	saveSnapshot()
}

func toggleTheorems() {

	oTHM = !oTHM

	setBasicInferenceRules()
	if oTHM {
		setupTheorems(oDR)
	}
	return
}

func togglePL() {
	oPL = !oPL
	logConstBindings = nil
	if oPL {
		logConstBindings = append(connBindings, plBindings...)
	} else {
		logConstBindings = connBindings
	}
	gentzen.SetPL(oPL)
	return
}

func setupButtonLabels() {
	if oPL {
		setTextByID("toggleSystem", "Predicate Logic")
	} else {
		setTextByID("toggleSystem", "Sentential Logic")
	}
	if oTHM {
		setupTheorems(oDR)
		setTextByID("toggleTheorem", "Theorems ✓")
	} else {
		setTextByID("toggleTheorem", "Theorems")
	}
}

func toggleHelp() {
	toggleVisibilityInline("help")
}

func toggleReadme() {
	toggleVisibility("description")
	toggleVisibility("proofChecker")
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
	//checkDerivation()
	return
}

func setOffset() {
	/*
	   		n, err := strconv.Atoi(js.Global().Call("prompt", "Number of first line", strconv.Itoa(dsp.offset)).String())
	   		if err != nil {
	   			return
	   		}
	   		dsp.setOffset(n)
	   		setTextByID("setOffset", "First Line: "+strconv.Itoa(dsp.offset))
	   	}

	   func setTitle() {

	   	title := js.Global().Call("prompt", "Title:").String()
	   	dsp.title = convert(title)
	*/
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
	switch oClipboard {

	case oLatexOutput:
		//	copyToClipboard(latexOutput())

	}

	return
}
func toClipboardLatex() {
	if oABOUT {
		return
	}

	//	copyToClipboard(latexOutput())

	return
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

func sizeUp() {
	s := domDocument.Call("querySelector", "#editorWindow").Get("style")
	v := strings.TrimSuffix(s.Call("getPropertyValue", "font-size").String(), "%")
	iv, err := strconv.Atoi(v)
	if err != nil {
		iv = 120
	}
	v = strconv.Itoa(iv+20) + "%"
	s.Call("setProperty", "font-size", v)
}

func sizeDown() {
	s := domDocument.Call("querySelector", "#editorWindow").Get("style")
	v := strings.TrimSuffix(s.Call("getPropertyValue", "font-size").String(), "%")
	iv, err := strconv.Atoi(v)
	if err != nil {
		iv = 120
	}
	v = strconv.Itoa(iv-20) + "%"
	s.Call("setProperty", "font-size", v)
}

var screenStash string

func stashScreen() {

	screenStash = dom.GetWindow().Document().GetElementsByTagName("body")[0].InnerHTML()

}

func restoreScreen() {

	dom.GetWindow().Document().GetElementsByTagName("body")[0].SetInnerHTML(screenStash)

}

func toggleSettingsMenu() {
	toggleVisibilityInline("settingsMenu")
}

func toggleVisibilityInline(id string) {
	s := domDocument.Call("querySelector", "#"+id).Get("style")
	if s.Call("getPropertyValue", "display").String() == "inline-block" {
		s.Call("setProperty", "display", "none")
	} else {
		s.Call("setProperty", "display", "inline-block")
	}
}

func toggleVisibility(id string) {
	s := domDocument.Call("querySelector", "#"+id).Get("style")
	if s.Call("getPropertyValue", "display").String() == "block" {
		s.Call("setProperty", "display", "none")
	} else {
		s.Call("setProperty", "display", "block")
	}
}
