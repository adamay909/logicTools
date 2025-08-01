package main

import (
	"strconv"
	"strings"

	"github.com/adamay909/logicTools/gentzen"
)

const clean = true

func checkDerivation() {
	debug("start proof checker")
	if dsp.empty() {
		return
	}

	setAttributeByID("display", "class", "inactive-fail")
	printMessage("", !clean)
	show("messages")
	gentzen.SetStrict(false)
	gentzen.ClearLog()
	arglines, ok := getArglines(dsp.Input)
	if !ok {
		printMessage(gentzen.ShowLog(), clean)
		debug("error parsing derivation lines")
		return
	}

	displayDerivation()

	if gentzen.CheckDeriv(arglines, dsp.Offset) {
		printMessage("OK", !clean)
		showPrettyDeriv(dsp)
		setAttributeByID("display", "class", "inactive-success")
		return
	}

	printMessage(gentzen.ShowLog(), clean)

	return
}

func printMessage(s string, cleanup bool) {

	if cleanup {
		s = strings.ReplaceAll(s, `/`, `\`)

		for _, e := range allBindings {
			s = strings.ReplaceAll(s, e[0], e[2])
		}
	}

	l := strings.Split(s, "\n")
	for i := range l {
		l[i] = `<p class="message">` + l[i] + "</p>\n"
	}
	s = strings.Join(l, "\n")

	setTextByID("messages", s)
}

func showPrettyDeriv(d *console) {

	d.html = nil

	var lines []string
	if arglines, ok := getArglines(dsp.Input); ok {
		lines = strings.Split(gentzen.PrintDerivation(arglines, dsp.Offset, gentzen.O_ProofChecker), "\n")
	} else {
		return
	}
	offset, _ := strconv.Atoi(lines[0][:strings.Index(lines[0], ".")])
	for i, l := range lines {

		debug("received line:" + l)

		if strings.TrimSpace(l) == "" {
			break
		}
		ln := strconv.Itoa(i+offset) + ". "
		p := strings.Index(l, ".")
		t := strings.Index(l, "⊢")
		s := strings.Index(l, "...")
		datum := l[p+1 : t]
		succ := l[t+3 : s] //the turnstile is multibyte rune!
		annot := l[s+3:]

		r := `<div class="ln">#ln#</div><div class="ddat">#dat#</div><div class="dtstl">⊢</div><div class="succ">#succ#</div><div class="dsep">...</div><div class="dannot">#annot#</div>`

		r = strings.Replace(r, `#ln#`, ln, 1)
		r = strings.Replace(r, `#dat#`, datum, 1)
		r = strings.Replace(r, `#succ#`, succ, 1)
		r = strings.Replace(r, `#annot#`, annot, 1)

		d.html = append(d.html, prettyGreek(r))
	}
	setTextByID("display", d.typeset())
}

func prettyGreek(r string) string {

	var r2 string
	var found bool
	for _, c := range r {
		found = false
		for _, e := range greekBindings {
			if string(c) == e[tktxt] {
				r2 = r2 + `<span class="greek">` + e[tktxt] + `</span>`
				found = true
				break
			}
		}
		if !found {
			r2 = r2 + string(c)
		}
	}

	return strings.ReplaceAll(r2, "⊢", " ⊢ ")
}
