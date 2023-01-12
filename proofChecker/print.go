package main

import (
	"strconv"
	"strings"

	"github.com/adamay909/logicTools/gentzen"
)

func plainTextDeriv(withTitle bool) string {
	if dsp.empty() {
		return ""
	}

	output := ""

	if withTitle && dsp.Title != "" {
		output = output + "[[TITLE:" + dsp.Title + "]]" + "\n"
	}

	if arglines, ok := getArglines(dsp.Input); ok {
		output = output + gentzen.PrintDerivText(arglines, dsp.Offset)
		return output
	}

	for i, l := range dsp.Input {
		if len(l) == 0 {
			continue
		}
		output = output + strconv.Itoa(i+dsp.Offset) + `.` + plainOutput(l) + "\n"
	}
	return output + "\n"

}

func plainOutput(s []string) string {

	var r string

	for _, e := range s {
		r = r + plainHTML(e)
	}
	return r

}

func plainHTML(s string) string {

	for _, e := range allBindings {
		if s == e[tktex] {
			if isGreek(s) {
				return `<span class="greek">` + e[tktxt] + `</span>`
			}
			return e[tktxt]
		}
	}
	return s
}

func plainText(s string) string {
	for _, e := range greekBindings {
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

	output := ""
	if arglines, ok := getArglines(dsp.Input); ok {
		output = `\subsubsection*{` + dsp.Title + `}` + "\n"
		output = strings.ReplaceAll(output, "⊢", `$\lproves$`)
		output = output + gentzen.PrintDeriv(arglines, dsp.Offset)
		return output
	}
	ln := strconv.Itoa(dsp.Offset - 1)
	output = `\subsubsection*{` + dsp.Title + `}` + "\n"

	output = output + `\begin{enumerate}\setcounter{enumi}{` + ln + `}` + "\n"

	for _, l := range dsp.Input {
		if len(l) == 0 {
			continue
		}
		output = output + `\item ` + plainOutput(l) + "\n"
	}
	return output + `\end{enumerate}` + "\n"
}

func isGreek(s string) bool {

	for _, e := range greekBindings {
		if s == e[tktex] {
			return true
		}
	}
	return false
}
