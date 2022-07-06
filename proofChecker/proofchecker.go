package main

import (
	"errors"
	"strconv"
	"strings"
	"syscall/js"

	"github.com/adamay909/logicTools/gentzen"
)

func checkDerivation() {
	if emptyCanvas() {
		return
	}

	printMessage("")
	gentzen.SetStrict(false)
	gentzen.ClearLog()

	arglines, ok := parseLines(canvas)
	if !ok {
		printMessage("RESULTS:\n\n" + gentzen.ShowLog())
		return
	}

	typesetDeriv()

	if gentzen.CheckDeriv(arglines, oOffset) {
		printMessage("RESULTS:\n\n Good!" + gentzen.ShowLog())
		if oPRIVATE {
			js.Global().Get("navigator").Get("clipboard").Call("writeText", gentzen.PrintDeriv(arglines, oOffset))
		}
		return
	}

	printMessage("RESULTS:\n\n" + gentzen.ShowLog())

	return
}

func emptyCanvas() bool {
	if len(canvas) == 0 {
		return true
	}
	for _, e := range canvas {
		if len(e) > 0 {
			return false
		}
	}
	return true
}

func parseLines(c [][]string) (s []string, ok bool) {

	ok = true

	if emptyCanvas() {
		gentzen.WriteLog("Nothing to do.", "")
		return
	}
	for i, line := range c {

		if len(line) == 0 {
			if i+1 < length(canvas) {
				gentzen.WriteLog("You seem to have an empty line in the middle.", "line "+strconv.Itoa(i+1)+": ")
				ok = false
				return
			}
			break
		}
		datum, succ, annot, err := parseLine(line)
		if err != nil {
			gentzen.WriteLog(err.Error(), "line "+strconv.Itoa(i+1)+": ")
			ok = false
			continue
		}

		s = append(s, datum+";"+succ+";"+annot)

	}
	return
}

func length(l [][]string) int {

	for n := len(l) - 1; n > 0; n-- {
		if !emptyLine(l[n]) {
			return n
		}
	}
	return 0
}

func parseLine(raw []string) (datum, succ, annot string, err error) {

	if len(raw) == 0 {
		return
	}
	var formula *gentzen.Node

	datum, succ, annot, err = splitLine(raw)

	if err != nil {
		return
	}

	//Deal with datum
	data := strings.Split(datum, ",")
	datum = ""
	for _, e := range data {
		f := strings.TrimSpace(e)
		if len(f) == 0 {
			continue
		}
		if isGreekLetter(f) {
			datum = datum + f + ","
			continue
		}
		formula, err = gentzen.InfixParser(tk(e))
		if err != nil {
			return
		}
		datum = datum + formula.String() + ","
	}
	datum = strings.TrimSuffix(datum, ",")

	//Deal with succedent

	formula, err = gentzen.InfixParser(tk(succ))
	if err != nil {
		return
	}
	succ = formula.String()

	//deal with annotation
	annot = replaceInfrules(annot)

	return
}

func parseLineDisplay(raw []string) (datum, succ, annot string, err error) {

	if len(raw) == 0 {
		err = errors.New("nothing there")
		return
	}
	var formula *gentzen.Node

	datum, succ, _, err = splitLine(raw)

	if err != nil {
		return
	}

	//Deal with datum
	data := strings.Split(datum, ",")
	datum = ""
	for _, e := range data {
		f := strings.TrimSpace(e)
		if len(f) == 0 {
			continue
		}
		if isGreekLetter(f) {
			datum = datum + plainText(f) + ","
			continue
		}
		formula, err = gentzen.InfixParser(tk(e))
		if err != nil {
			return
		}
		datum = datum + formula.StringPlain() + ","
	}
	datum = strings.TrimSuffix(datum, ",")

	//Deal with succedent

	formula, err = gentzen.InfixParser(tk(succ))
	if err != nil {
		return
	}
	succ = formula.StringPlain()

	an := split(raw, `\ldots`)[1]

	for _, e := range an {
		annot = annot + plainText(e)
	}
	return
}

func tk(s string) (t []string) {

	d := strings.Split(s, " ")

	for _, e := range d {
		if len(e) != 0 {
			t = append(t, e)
		}
	}

	return t

}

func splitLine(raw []string) (datum, succ, annot string, err error) {

	s := strings.Join(raw, " ")

	if strings.Count(s, `\vdash`) != 1 {
		err = errors.New("malformed derivation line")
		return
	}

	if strings.Count(s, `\ldots`) != 1 {
		err = errors.New("malformed derivation line")
		return
	}

	if strings.Index(s, `\vdash`) > strings.Index(s, `\ldots`) {
		err = errors.New("Malformed derivation line")
		return
	}

	datum = strings.Split(s, `\vdash`)[0]

	succ = strings.Split(strings.Split(s, `\vdash`)[1], `\ldots`)[0]

	annot = strings.Split(strings.Split(s, `\vdash`)[1], `\ldots`)[1]

	return
}

func replaceInfrules(s string) string {
	infrules := [][2]string{
		[2]string{`\negE`, `ne`},
		[2]string{`\negI`, `ni`},
		[2]string{`\veeE`, `de`},
		[2]string{`\veeI`, `di`},
		[2]string{`\wedgeE`, `ke`},
		[2]string{`\wedgeI`, `ki`},
		[2]string{`\supsetE`, `ce`},
		[2]string{`\supsetI`, `ci`},
		[2]string{`A`, `a`},
		[2]string{`\forallE`, `ue`},
		[2]string{`\forallI`, `ui`},
		[2]string{`\existsE`, `ee`},
		[2]string{`\existsI`, `ei`},
		[2]string{`=E`, `=e`},
		[2]string{`=I`, `=i`},
	}
	s = strings.ReplaceAll(s, " ", "")
	a := strings.Split(s, ",")

	for i := range a {
		for _, c := range infrules {
			if strings.TrimSpace(a[i]) == c[0] {
				a[i] = c[1]
			}
		}
	}
	return strings.Join(a, ",")
}

func printMessage(s string) {
	l := strings.Split(s, "\n")
	for i := range l {
		l[i] = "<p>" + l[i] + "</p>" + "\n"
	}
	s = strings.Join(l, "\n")
	setTextByID("messages", s)
}

func isGreekLetter(s string) bool {

	for _, c := range greekBindings {
		if s == c[1] {
			return true
		}
	}
	return false
}

func split(s []string, cut string) [][]string {

	var resp [][]string

	var sub []string
	for _, e := range s {
		if e == cut {
			resp = append(resp, sub)
			sub = nil
		} else {
			sub = append(sub, e)
		}
	}
	if len(sub) != 0 {
		resp = append(resp, sub)
	}
	return resp
}
