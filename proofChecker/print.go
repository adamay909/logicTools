package main

import (
	"strconv"

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
		output = output + strconv.Itoa(i+dsp.Offset) + `. ` + plainTextOutput(l) + "\n"
	}
	return output + "\n"

}

func plainHTMLOutput(s []string) string {

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

func plainTextOutput(s []string) string {

	var r string

	for _, e := range s {
		r = r + plainText(e)
	}
	return r

}

func plainLatextOutput(s []string) string {

	var r string

	for _, e := range s {
		r = r + plainLatex(e)
	}
	return r

}

func plainLatex(s string) string {

	for _, e := range greekBindings {
		if s == e[tktxt] {
			return e[tktex]
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
	if s == `\vdash` {
		s = `⊢`
	}

	if s == `\ldots` {
		s = `...`
	}
	return s
}

func latexOutput() string {
	if dsp.empty() {
		return ""
	}

	output := ""
	output = `%title:%` + dsp.Title + "\n"

	if arglines, ok := getArglines(dsp.Input); ok {
		for c, k := range arglines {
			debug(c, ": ", k)
		}

		output = output + gentzen.PrintDeriv(arglines, dsp.Offset)
		return safeLtx(output)
	}

	ln := strconv.Itoa(dsp.Offset - 1)
	output = output + `\begin{enumerate}\setcounter{enumi}{` + ln + `}` + "\n"

	for _, l := range dsp.Input {
		if len(l) == 0 {
			continue
		}
		output = output + `\item ` + safeLtx(latexfy(l)) + "\n"
	}
	return output + `\end{enumerate}` + "\n"
}

func findFormula(l []string) (string, int) {

	var ret []string

	for _, e := range l {
		ret = append(ret, plainText(e))
	}

	if len(ret) == 0 {
		return "", 0
	}

	var i int
	var formula *gentzen.Node
	var err error

	for i = len(ret); i > 1; i-- {
		txt := spaceyStringOf(ret[:i])
		formula, err = gentzen.InfixParser(tk(txt))
		if err == nil {
			debug("findFormula: returning ", formula.StringLatex())
			return formula.StringLatex(), len(ret[:i])
		}
	}

	return "", 0
}

func latexfy(l []string) string {

	if len(l) == 0 {
		return ""
	}

	var out string

	for start := 0; start < len(l); {
		f, span := findFormula(l[start:])
		if f != "" {
			out = out + `\p{` + f + `}`
			start = start + span
			continue
		}
		out = out + plainTextOutput(l[start:start+1])
		start++
	}

	return out
}

func isGreek(s string) bool {

	for _, e := range greekBindings {
		if s == e[tktex] {
			return true
		}
	}
	return false
}

func tknz(s string) []string {

	var ret []string

	ab := combineBindings(allBindings, extraBindings, connBindings, plBindings, mlBindings)

	for _, c := range s {
		brk := false
		for _, t := range ab {
			if string(c) == t[tktxt] {
				ret = append(ret, t[tktex]+" ")
				brk = true
				break
			}
		}
		if !brk {
			ret = append(ret, string(c))
			brk = false
		}
	}
	return ret
}

func safeLtx(s string) string {
	var ret string
	for _, e := range s {

		ret = ret + ltxof(string(e))

	}

	return ret
}

func ltxof(e string) string {

	ab := combineBindings(greekBindings, extraBindings, connBindings, plBindings, mlBindings, turnstileBindings)

	for _, c := range ab {
		if e == c[tktxt] {
			return c[tktex] + " "
		}
	}

	return e
}

func printTree() {

	gentzen.SetStrict(false)
	gentzen.ClearLog()
	arglines, ok := getArglines(dsp.Input)
	if !ok {
		printMessage(gentzen.ShowLog())
		debug("error parsing derivation lines")
		return
	}

	copyToClipboard(gentzen.PrintDerivTree(arglines, dsp.Offset))

}
