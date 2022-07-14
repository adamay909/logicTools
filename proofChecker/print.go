package main

import (
	"strconv"

	"github.com/adamay909/logicTools/gentzen"
)

func plainTextDeriv() string {
	if dsp.empty() {
		return ""
	}

	if arglines, ok := getArglines(dsp.Input); ok {
		return gentzen.PrintDerivText(arglines, dsp.Offset)
	}
	output := ""

	for i, l := range dsp.Input {
		if len(l) == 0 {
			continue
		}
		output = output + strconv.Itoa(i+dsp.Offset) + plainOutput(l) + "\n"
	}
	return output + "\n"

}

func plainOutput(s []string) string {

	var r string

	for _, e := range s {
		r = r + plainText(e)
	}
	return r

}

func plainText(s string) string {

	for _, e := range allBindings {
		if s == e[tktex] {
			return e[tktxt]
		}
	}
	return s
}

func latexOutput() string {
	if dsp.empty() {
		return ""
	}

	if arglines, ok := getArglines(dsp.Input); ok {
		return gentzen.PrintDeriv(arglines, dsp.Offset)
	}
	output := ""
	ln := strconv.Itoa(dsp.Offset - 1)
	output = `\begin{enumerate}\setcounter{enumi}{` + ln + `}` + "\n"

	for _, l := range dsp.Input {
		if len(l) == 0 {
			continue
		}
		output = output + `\item ` + plainOutput(l) + "\n"
	}
	return output + `\end{enumerate}` + "\n"
}
