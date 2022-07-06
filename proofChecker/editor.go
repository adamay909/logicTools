package main

import (
	_ "embed"
	"strconv"
	"strings"
	"syscall/js"

	"honnef.co/go/js/dom/v2"
)

func clearCanvas() {
	stopInput()
	var line []string
	canvas = nil
	canvas = append(canvas, line)
	xpos, ypos = 0, 0
	waitTurnstile = false
	waitESC = false
	waitSub = false
	dom.GetWindow().Document().GetElementByID("inputArea").SetAttribute("style", "counter-reset: line "+strconv.Itoa(0)+";")
	d, _ := assets.ReadFile("assets/html/startmessage.html")

	bottom := `<div id="overlay" contenteditable="true" tabindex=1 onclick="activateInput()"> dummy text </div>` + "\n"
	setTextByID("inputArea", string(d)+bottom)
	printMessage("")
	return
}

var mod int

func typeInput(key string) bool {
	if waitTurnstile {
		waitTurnstile = false
		mod++
		if ok, target := among(key, turnstileBindings); ok {
			addChar(target)
		}
		typesetCanvas()
		mod = 0
		return true
	}

	if waitESC {
		if key != "Shift" {
			waitESC = false
		}
		if ok, target := among(key, greekBindings); ok {
			addChar(target)
		}
		typesetCanvas()
		return true
	}

	if waitSub {
		if key != "Shift" {
			waitSub = false
		}
		if ok, target := among(key, keyBindings); ok {
			addSubscript(target)
		}
		typesetCanvas()
		return true
	}

	switch key {

	case "Backspace":
		backspace()

	case "ArrowLeft":
		arrowLeft()

	case "ArrowRight":
		arrowRight()

	case "ArrowUp":
		arrowUp()

	case "ArrowDown":
		arrowDown()

	case "Home":
		home()

	case "End":
		end()

	case "Enter":
		newLine()

	case "Delete":
		deletechar()

	case "|":
		waitTurnstile = true
		insertExtraChar(`|`)
		return true

	case `\`:
		waitESC = true
		insertExtraChar(`\`)
		return true

	case "_":
		waitSub = false //don't allow subscripts
		return false

	case ".":
		addChar(`\ldots`)

	default:

		if ok, target := among(key, keyBindings, punctBindings, logConstBindings); ok {
			addChar(target)
		}
	}

	typesetCanvas()
	return true
}

func addSubscript(s string) {
	if xpos == 0 {
		return
	}
	char := canvas[ypos][xpos-1]
	backspace()
	addChar(char + `_` + s[:1])
	return
}

func addChar(c string) {
	if c == "" {
		return
	}
	var l, l1, l2 []string
	//this seems to work but I'm not
	//sure about the workings of Go slices.
	l = canvas[ypos]
	l1 = append(l1, l[:xpos]...)
	l2 = append(l2, l[xpos:]...)

	ln := append(l1, c)
	ln = append(ln, l2...)
	canvas[ypos] = ln
	xpos = xpos + 1
	return
}

func arrowUp() {

	if ypos == 0 {
		return
	}
	ypos = ypos - 1
	if len(canvas[ypos]) < xpos+1 {
		xpos = len(canvas[ypos]) - 1
		if xpos < 1 {
			xpos = 0
		}
	}
}

func arrowDown() {
	if ypos == len(canvas)-1 {
		return
	}
	ypos = ypos + 1
	if len(canvas[ypos]) < xpos+1 {
		xpos = len(canvas[ypos]) - 1
		if xpos < 1 {
			xpos = 0
		}
	}
}
func arrowRight() {
	if emptyLine(canvas[ypos]) {
		return
	}

	if overEnd {
		return
	}

	if eol() {
		s := canvas[ypos][xpos]
		deletechar()
		addChar(s)
		return
	}
	xpos++
}

func home() {
	xpos = 0
	return
}

func end() {
	xpos = len(canvas[ypos]) - 1
	if xpos < 0 {
		xpos = 0
	}
	return
}

func arrowLeft() {
	if xpos == 0 {
		return
	}
	xpos = xpos - 1
}

func newLine() {
	var l1, l2 [][]string
	var l []string
	var f1, f2 []string

	l = canvas[ypos]

	l1 = canvas[:ypos]
	l2 = append(l2, canvas[ypos+1:]...)

	f1 = l[:xpos]
	f2 = l[xpos:]
	canvas = nil
	canvas = append(l1, f1)
	canvas = append(canvas, f2)
	canvas = append(canvas, l2...)
	xpos = 0
	ypos = ypos + 1
}

func backspace() {
	var l []string
	switch xpos {

	case 0:
		arrowUp()
		end()
		arrowRight()
		arrowRight()
		return

	case 1:
		if len(canvas[ypos]) < 2 {
			canvas[ypos] = nil
			xpos = 0
			//..deleteLine()
			//end()
			//arrowRight()
			//newLine()
		} else {
			canvas[ypos] = canvas[ypos][1:]
			xpos = 0
		}
	default:
		l = canvas[ypos]
		ln := append(l[:xpos-1], l[xpos:]...)
		canvas[ypos] = ln
		xpos = xpos - 1
	}
	typesetCanvas()
	return
}

func deletechar() {
	var l []string
	var l1 []string

	l = canvas[ypos]

	if len(l) == 0 {
		deleteLine()
		return
	}
	if overEnd {
		return
	}
	if eol() {
		canvas[ypos] = l[:len(l)-1]
	} else {
		l1 = append(l[:xpos], l[xpos+1:]...)
		canvas[ypos] = l1
	}
	if len(canvas[ypos]) == 0 {
		deleteLine()
	}
	return

}

func deleteLine() {
	var l1, l2 [][]string
	if len(canvas) > 1 {
		l1 = canvas[:ypos]
	}

	if !onLastLine() {
		l2 = canvas[ypos+1:]
	}
	canvas = nil

	canvas = append(l1, l2...)
	if len(canvas) == 0 {
		canvas = append(canvas, []string{})
	}
	xpos = 0
	arrowUp()
	return
}

func onLastLine() bool {
	return ypos+1 == len(canvas)
}

func formatLineNC(l []string) string {

	var r string
	datum, succ, annot, err := parseLineDisplay(l)

	var line string
	if err != nil {
		for _, e := range l {
			line = line + pad(plainText(e))
		}
		r = r + `<div class=formula">` + line + `</div>`
		return r
	}

	r = `<div class="pdatum">#datum#</div><div class="pturnstile">⊢</div><div class="psuccedent">#succ#</div><div class="pdots"></div><div class="pannotation">#annot#</div>`

	r = strings.Replace(r, "#datum#", datum, 1)
	r = strings.Replace(r, "#succ#", succ, 1)
	r = strings.Replace(r, "#annot#", annot, 1)

	return r
}

func formatLine(l []string, p int) string {

	var r string
	overEnd = false
	for i, e := range l {
		if i == p {
			r = r + `<div class="formula" id="cursor"><div class="extra"></div><div class="normal" id="blink">` + pad(plainText(e)) + `</div> </div>`
			continue
		}
		r = r + `<div class="formula">` + pad(plainText(e)) + `</div>`
	}
	if p == len(l) {
		overEnd = true
		r = r + `<div class="formula" id="cursor"><div class="extra"></div><div class="normal" id="blink">&emsp;</div> </div>`
	}

	return r

}

func insertExtraChar(c string) {

	dom.GetWindow().Document().GetElementsByClassName("extra")[0].SetAttribute("id", "blink")
	dom.GetWindow().Document().GetElementsByClassName("extra")[0].SetInnerHTML(c)
	dom.GetWindow().Document().GetElementsByClassName("normal")[0].RemoveAttribute("id")

	return
}

func pad(e string) string {
	switch e {

	case "⊢":
		return `&emsp;` + "⊢" + `&emsp;`

	case "...":
		return `&emsp;` + "..." + `&emsp;`

	default:
		return e
	}
}

func among(s string, bindings ...[][3]string) (bool, string) {

	for _, b := range bindings {
		for _, e := range b {
			if s == e[0] {
				return true, e[1]
			}
		}
	}
	return false, ""
}

func typesetCanvas() {
	printMessage("")
	var s string

	for n := range canvas {

		sn := `<div class="pline" id="#ln#"><div class="plnum">#lnumber# .&emsp;</div>#argline#</div>` + "\n\n"

		sn = strings.Replace(sn, "#ln#", "line"+strconv.Itoa(n), 1)
		sn = strings.Replace(sn, "#lnumber#", strconv.Itoa(n+oOffset), 1)
		if n == ypos {
			sn = strings.Replace(sn, "#argline#", formatLine(canvas[n], xpos), 1)
		} else {
			sn = strings.Replace(sn, "#argline#", formatLineNC(canvas[n]), 1)
		}
		s = s + sn
	}

	if s == "" {
		s = `<div class="pline" id="line0"><div class="plnum"></div></div>`
	}

	bottom := `<div id="overlay" contenteditable="true" tabindex=1 onclick="activateInput()"> dummy text </div>` + "\n"

	setTextByID("inputArea", s+bottom)

	dom.GetWindow().Document().GetElementByID("line0").GetElementsByClassName("plnum")[0].SetAttribute("onclick", "setOffset()")

	dom.GetWindow().Document().GetElementByID("line0").GetElementsByClassName("plnum")[0].SetAttribute("style", "cursor: grab")

	js.Global().Get("overlay").Call("focus")
}

func typesetDeriv() {

	var s string
	for n, line := range canvas {
		if len(line) == 0 {
			continue
		}
		datum, succ, annot, _ := parseLineDisplay(line)

		sn := `<div class="linenumber">#ln#.&emsp;</div><div class="datum">#datum#</div><div class="turnstile">⊢</div><div class="subs">#succ#</div><div class="separator"></div><div class="annotation">#annot#</div>` + "\n"

		sn = strings.Replace(sn, "#ln#", strconv.Itoa(n+1), 1)
		sn = strings.Replace(sn, "#datum#", datum, 1)
		sn = strings.Replace(sn, "#succ#", succ, 1)
		sn = strings.Replace(sn, "#annot#", annot, 1)

		s = s + sn
	}
	s = `<div class="deriv">` + s + `</div>`
	bottom := `<div id="overlay" contenteditable="true" tabindex=1 onclick="activateInput()"> dummy text </div>` + "\n"
	setTextByID("inputArea", s+bottom)

	js.Global().Get("overlay").Call("focus")
	return
}

func emptyLine(l []string) bool {
	return len(l) == 0
}

func eol() bool {
	return xpos == len(canvas[ypos])-1
}
