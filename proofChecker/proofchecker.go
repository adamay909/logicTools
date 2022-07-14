package main

import (
	"strconv"
	"strings"

	"github.com/adamay909/logicTools/gentzen"
)

func checkDerivation() {
	if dsp.empty() {
		return
	}

	printMessage("")
	show("messages")
	gentzen.SetStrict(false)
	gentzen.ClearLog()

	arglines, ok := getArglines(dsp.Input)
	if !ok {
		printMessage("RESULTS:\n\n" + gentzen.ShowLog())
		return
	}

	displayDerivation()

	if gentzen.CheckDeriv(arglines, dsp.Offset) {
		printMessage("RESULTS:\n\n Good!" + gentzen.ShowLog())
		showPrettyDeriv(dsp)
		return
	}

	printMessage("RESULTS:\n\n" + gentzen.ShowLog())

	return
}

func printMessage(s string) {
	l := strings.Split(s, "\n")
	for i := range l {
		l[i] = "<p>" + l[i] + "</p>" + "\n"
	}
	s = strings.Join(l, "\n")
	setTextByID("messages", s)
}

func showPrettyDeriv(d *console) {

	d.html = nil

	var lines []string
	if arglines, ok := getArglines(dsp.Input); ok {
		lines = strings.Split(gentzen.PrintDerivText(arglines, dsp.Offset), "\n")
	} else {
		return
	}
	offset, _ := strconv.Atoi(lines[0][:strings.Index(lines[0], ".")])
	for i, l := range lines {
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
		d.html = append(d.html, r)
	}
	setTextByID("display", d.typeset())
}
