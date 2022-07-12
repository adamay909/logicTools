package main

import (
	"strings"

	"github.com/adamay909/logicTools/gentzen"
)

func checkDerivation() {
	if dsp.empty() {
		return
	}

	printMessage("")
	gentzen.SetStrict(false)
	gentzen.ClearLog()

	arglines, ok := parseLines(dsp.Input)
	if !ok {
		printMessage("RESULTS:\n\n" + gentzen.ShowLog())
		return
	}

	displayDerivation()

	if gentzen.CheckDeriv(arglines, dsp.offset) {
		printMessage("RESULTS:\n\n Good!" + gentzen.ShowLog())
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
