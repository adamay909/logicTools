package main

import (
	"strconv"

	"github.com/adamay909/logicTools/gentzen"
)

func plainTextDeriv() string {

	var r string

	for n, l := range dsp.Input {
		if len(l) == 0 {
			continue
		}

		r = r + strconv.Itoa(n+1) + ". " + plainOutput(l) + "\n"
	}

	return r
}

func plainOutput(s []string) string {

	var r string

	for _, e := range s {
		r = r + plainText(e)
	}
	return r

}

func plainText(s string) string {

	allBindings := [][][3]string{
		keyBindings,
		punctBindings,
		connBindings,
		plBindings,
		turnstileBindings,
		greekBindings,
	}

	for _, b := range allBindings {
		for _, e := range b {
			if s == e[1] {
				if e[2] != "" {
					return e[2]
				}
				return s
			}
		}
	}
	return s
}

func latexOutput() string {
	if dsp.empty() {
		return ""
	}

	if arglines, ok := parseLines(dsp.Input); ok {
		return gentzen.PrintDeriv(arglines, dsp.offset)
	}
	output := ""
	ln := strconv.Itoa(dsp.offset - 1)
	output = `\begin{enumerate}\setcounter{enumi}{` + ln + `}` + "\n"

	for _, l := range dsp.Input {
		if len(l) == 0 {
			continue
		}
		output = output + `\item ` + plainOutput(l) + "\n"
	}
	return output + `\end{enumerate}` + "\n"
}
