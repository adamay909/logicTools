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
		output = output + `\item ` + attemptLatex(l) + "\n"
	}
	return output + `\end{enumerate}` + "\n"
}

func attemptLatex(l []string) string {

	var ret []string

	for _, e := range l {
		ret = append(ret, plainText(e))
	}

	if len(ret) == 0 {
		return ""
	}

	var i int
	var formula *gentzen.Node
	var err error

	for i = len(ret); i > 0; i-- {
		txt := spaceyStringOf(ret[:i])
		formula, err = gentzen.InfixParser(tk(txt))
		if err == nil {
			break
		}
	}

	if i > 0 {

		return `\p{` + formula.StringLatex() + `}` + plainOutput(l[i:])

	}

	return plainOutput(l)
}

func isGreek(s string) bool {

	for _, e := range greekBindings {
		if s == e[tktex] {
			return true
		}
	}
	return false
}
