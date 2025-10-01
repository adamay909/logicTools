package main

import (
	_ "embed"
	"strconv"
	"strings"
	"syscall/js"

	"github.com/adamay909/logicTools/gentzen"
	"github.com/adamay909/logicTools/proofChecker/internal/editor"
)

// console is a struct for holding global variables that control
// state of the proof checker window. You could have multiple proof
// checker windows.
type console struct {
	editor *editor.Editor
	title  *editor.Editor
	oPL,
	oDR,
	oTHM bool
	oOffset int
}

var dsp *console

func main() {
	dsp = new(console)
	gentzen.SetSpecialConn(true)
	setBasicInferenceRules()
	setupPage()
	setupJS()
	editor.AddCSS()
	checkForOldFormat()
	loadHistory()
	updatePageNumber()
	<-make(chan bool)
}

//go:embed assets/html/main.css
var styleCSS string

//go:embed assets/html/version
var vnum string

func setupPage() {

	domDocument.Call("getElementsByTagName", "style").Call("item", 0).Set("innerHTML", styleCSS)

	setTextByID("versionnumber", "v"+vnum+"&emsp;&emsp;")

	domDocument.Call("querySelector", "#description").Get("style").Call("setProperty", "color", "black")

}

func setupJS() {
	dsp.editor = editor.NewDerivationEditor("editor")
	dsp.title = editor.NewSimpleEditor("title")
	addEventListener(domDocument, "click", jsWrap(onClick))
	addEventListener(domDocument.Call("querySelector", "#editorWindow"), "editorinput", jsWrap(saveSnapshot))
	addEventListener(domDocument.Call("querySelector", "#editorWindow"), "focus", jsWrap(cleanupEditorWindow))
	addEventListener(domDocument.Call("querySelector", "#inputoffset"), "keydown", commitOffset)
	addEventListener(domDocument.Call("querySelector", "#inputoffset"), "blur", jsWrap(cancelInputOffset))

}

func onClick() {

	target := js.Global().Get("event").Get("target")
	switch target.Get("id").String() {

	case "toggleSettings":
		toggleSettingsMenu()
	case "check":
		checkDerivation()
	case "toLatex":
		copyToClipboard(latexOutput())
	case "printTree":
		copyToClipboard(printTree())
	case "toggleHelp":
		toggleHelp()

	case "toggleSystem":
		togglePL()
	case "setOffset":
		inputOffset()
	case "offset":
		inputOffset()
	case "toggleTheorem":
		toggleTheorems()
	case "newPage":
		newpage()
	case "backHistory":
		moveBackInHistory()
	case "forwardHistory":
		moveForwardInHistory()
	case "duplicateScreen":
		duplicateHistoryItem()
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

	dsp.oTHM = !dsp.oTHM

	setBasicInferenceRules()
	if dsp.oTHM {
		setupTheorems(dsp.oDR)
	}
	writeStateToHTML()
	return
}

func setupGentzen() {
	gentzen.SetPL(dsp.oPL)
	setBasicInferenceRules()
	if dsp.oTHM {
		setupTheorems(dsp.oDR)
	}
}

func togglePL() {
	dsp.oPL = !dsp.oPL
	gentzen.SetPL(dsp.oPL)
	writeStateToHTML()
	return
}

func setupButtonLabels() {
	if dsp.oPL {
		setTextByID("toggleSystem", "Predicate Logic")
	} else {
		setTextByID("toggleSystem", "Sentential Logic")
	}
	if dsp.oTHM {
		setupTheorems(dsp.oDR)
		setTextByID("toggleTheorem", "Theorems ✓")
	} else {
		setTextByID("toggleTheorem", "Theorems")
	}
	setTextByID("offset", strconv.Itoa(dsp.oOffset))
}

func toggleHelp() {
	toggleVisibilityInline("help")
}

func toggleReadme() {
	toggleVisibility("description")
	toggleVisibility("proofChecker")
}

func inputOffset() {
	e := domDocument.Call("querySelector", "#inputoffset")
	e.Set("value", strconv.Itoa(dsp.oOffset))
	toggleVisibilityInline("offset")
	toggleVisibilityInline("inputoffset")
	e.Call("focus")

}

func cancelInputOffset() {
	e := domDocument.Call("querySelector", "#inputoffset")
	e.Set("value", strconv.Itoa(dsp.oOffset))
	toggleVisibilityInline("offset")
	toggleVisibilityInline("inputoffset")
}

func setTextByID(elem string, content string) {
	domDocument.Call("querySelector", "#"+elem).Set("innerHTML", content)
}

func copyToClipboard(s string) {
	js.Global().Get("navigator").Get("clipboard").Call("writeText", s)
	return
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
func makeVisible(id string) {
	s := domDocument.Call("querySelector", "#"+id).Get("style")
	s.Call("setProperty", "display", "block")
}

func hide(id string) {
	s := domDocument.Call("querySelector", "#"+id).Get("style")
	s.Call("setProperty", "display", "none")
}

func latexOutput() string {
	arglines := dsp.editor.GetArglines()
	return gentzen.PrintDerivation(arglines, 1, gentzen.O_ProofChecker, gentzen.O_Latex)
}

func printTree() string {
	arglines := dsp.editor.GetArglines()
	r, err := gentzen.PrintDerivationTree(arglines, gentzen.O_ProofChecker, 1)
	if err != nil {
		return err.Error()
	}
	return r
}

func commitOffset(e js.Value, args ...any) {
	if e.Get("key").String() != "Enter" {
		return
	}
	inputelem := domDocument.Call("querySelector", "#inputoffset")
	v, err := strconv.Atoi(inputelem.Get("value").String())
	if err != nil {
		v = dsp.oOffset
	}
	if v < 1 || v > 999 {
		v = dsp.oOffset
	}
	inputelem.Set("value", strconv.Itoa(v))
	setOffset(v)
	writeStateToHTML()
	inputelem.Call("blur")
}

func setOffset(v int) {
	dsp.oOffset = v
	dsp.editor.SetOffset(v)
	domDocument.Call("querySelector", "#offset").Set("innerHTML", strconv.Itoa(v))
}
