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
)

type console struct {
	editor *editor.Editor
	title  *editor.Editor
}

// set to true for debug log to stdout

//go:embed assets/html/*
var assets embed.FS

// Enable some features for personal teaching material.
// Not useful for general consumption.

var oPRIVATE = true

/*
var inputFunc,

	mainEditorFunc,
	titleEditorFunc,
	clickFunc,
	snapshotFunc,
	loadFunc js.Value
*/
var indexHtml, helpHtml, styleCSS string

var (
	oPL              = false
	oDR              = false
	oTHM             = false
	oDEBUG           = false
	oOffset          = 1
	logConstBindings [][3]string
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
	checkForOldFormat()
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

	domDocument.Call("getElementsByTagName", "style").Call("item", 0).Set("innerHTML", string(d))
	d, _ = assets.ReadFile("assets/html/version")
	setTextByID("versionnumber", "v"+string(d)+"&emsp;&emsp;")

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

	oTHM = !oTHM

	setBasicInferenceRules()
	if oTHM {
		setupTheorems(oDR)
	}
	return
}

func setupGentzen() {
	gentzen.SetPL(oPL)
	setBasicInferenceRules()
	if oTHM {
		setupTheorems(oDR)
	}
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
	setTextByID("offset", strconv.Itoa(oOffset))
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
	e.Set("value", strconv.Itoa(oOffset))
	toggleVisibilityInline("offset")
	toggleVisibilityInline("inputoffset")
	e.Call("focus")

}

func cancelInputOffset() {
	e := domDocument.Call("querySelector", "#inputoffset")
	e.Set("value", strconv.Itoa(oOffset))
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
		v = oOffset
	}
	if v < 1 || v > 999 {
		v = oOffset
	}
	inputelem.Set("value", strconv.Itoa(v))
	setOffset(v)
	inputelem.Call("blur")
}

func setOffset(v int) {
	oOffset = v
	dsp.editor.SetOffset(v)
	domDocument.Call("querySelector", "#offset").Set("innerHTML", strconv.Itoa(v))
	domDocument.Call("querySelector", "#consoleState").Call("setAttribute", "data-offset", strconv.Itoa(v))
}
