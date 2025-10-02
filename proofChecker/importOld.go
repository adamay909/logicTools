package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"syscall/js"

	"github.com/adamay909/logicTools/proofChecker/internal/editor"
)

type inputLine []string

type oldconsole struct {
	Title        string
	Input        []inputLine
	SystemPL     bool
	SystemML     bool
	Theorems     bool
	DerivedRules bool
}

func checkForOldFormat() {
	if !js.Global().Get("localStorage").Call("getItem", "history").IsNull() {
		fmt.Println("detected material from previous version. Importing.")
		importOldHistory()
		return
	}
	if !js.Global().Get("localStorage").Call("getItem", "history2").IsNull() {
		fmt.Println("detected material from previous version. Importing.")
		updateHistory2()
	}
}

func updateHistory2() {

	raw := strings.ReplaceAll(js.Global().Get("localStorage").Call("getItem", "history2").String(), `<!---->`, "")
	js.Global().Get("localStorage").Call("setItem", "history3", raw)
	js.Global().Get("localStorage").Call("removeItem", "history2")
}

func importOldHistory() {

	dsp.historyItems = history.Get("children")
	oldhistory := loadOldHistory()
	historyPosition = 0
	for i := range oldhistory {
		page := new(oldconsole)
		err := json.Unmarshal([]byte(oldhistory[i]), page)
		if err != nil {
			log.Fatal("json falied")
		}
		newpage()
		inputlines := convertOldHistoryHTML(page.Input)
		dsp.editor.SetInnerHTML(inputlines)
		dsp.title.SetInnerHTML(page.Title)
		dsp.oPL = page.SystemPL
		dsp.oTHM = page.Theorems
		writeStateToHTML()
	}
	saveHistory()
}

func loadOldHistory() (oldhistory []string) {

	json := js.Global().Get("localStorage").Call("getItem", "history").String()

	oldhistory = strings.Split(json, "\n")

	oldhp, err := strconv.Atoi(js.Global().Get("localStorage").Call("getItem", "historyPosition").String())

	if err != nil {
		return
	}

	cur := js.Global().Get("localStorage").Call("getItem", "current")
	if cur.IsNull() {
		return
	}
	oldhistory[oldhp] = cur.String()
	js.Global().Get("localStorage").Call("setItem", "historyBackup1", strings.Join(oldhistory, "\n"))
	js.Global().Get("localStorage").Call("removeItem", "history")
	return
}

func convertOldHistoryHTML(input []inputLine) (html string) {

	var dat, succ, annot []string
	for _, l := range input {

		dat, succ, annot = parseNsplit(l)

		r := `
<div class="row"> 
<div class="ddat" contenteditable="plaintext-only">
#DAT
</div>
<div class="dtstl" contenteditable="false">
&vdash;
</div>
<div class="dsucc" contenteditable="plaintext-only">
#SUCC
</div>
<div class="dsep" contenteditable="false">
&#x2026;
</div>
<div class="dannot" contenteditable="plaintext-only">
#ANNOT
</div>
</div>
`

		r = strings.Replace(r, `#DAT`, stringOf(dat), 1)
		r = strings.Replace(r, `#SUCC`, stringOf(succ), 1)
		r = strings.Replace(r, `#ANNOT`, stringOf(annot), 1)

		r = strings.ReplaceAll(r, ",", ",\u00a0")

		html = html + r
	}

	if len(input) == 0 {
		dsp.editor.Clear()
		return dsp.editor.GetInnerHTML()
	}
	html = `<div class="derivation">` + html + `</div>`
	return editor.StripWhiteSpaceFromHTML(html)
}

func parseNsplit(l []string) (dat, succ, annot []string) {

	tst := getTstIdx(l)
	dot := getDotIdx(l)

	for i, e := range l {
		var text string
		text = plainHTML(e)

		if i < tst {
			dat = append(dat, text)
			continue
		}

		if i == tst {
			continue
		}

		if i < dot {
			succ = append(succ, text)
			continue
		}

		if i == dot {
			continue
		}

		annot = append(annot, text)
	}
	return
}

func getTstIdx(l []string) int {

	tst := index(l, `\vdash`)
	dot := index(l, `\ldots`)

	if tst == -1 {
		return len(l)
	}

	if dot > -1 && dot < tst {
		return len(l)
	}

	return tst
}

func getDotIdx(l []string) int {

	if getTstIdx(l) == len(l) {
		return -1
	}

	idx := index(l, `\ldots`)
	if idx == -1 {
		idx = len(l)
	}
	return idx
}

func stringOf(src []string) string {

	return strings.Join(src, "")
}

func plainHTML(s string) string {

	for _, e := range allBindings {
		if s == e[tktex] {
			return e[tktxt]
		}
	}
	return s
}

func index(s []string, t string) int {
	for i := range s {
		if s[i] == t {
			return i
		}
	}
	return -1
}

func (d *oldconsole) unmarshalJson(data string) {

	err := json.Unmarshal([]byte(data), d)

	if err != nil {
		log.Fatal("json falied")
	}
}
