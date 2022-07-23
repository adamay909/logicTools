package main

import (
	"errors"
	"strconv"
	"strings"

	"github.com/adamay909/logicTools/gentzen"
)

var (
	turnstile = `\vdash`
	separator = `\ldots`
)

func getArglines(c []inputLine) (s []string, ok bool) {

	ok = true

	if dsp.empty() {
		gentzen.WriteLog("Nothing to do.", "")
		return
	}
	for i := 0; i <= dsp.lastLine(); i++ {
		line := dsp.Input[i]
		if len(line) == 0 || line == nil {
			gentzen.WriteLog("You seem to have an empty line in the middle.", "line "+strconv.Itoa(i+1)+": ")
			ok = false
			return
		}

		err := isArgline(line)
		if err != nil {
			gentzen.WriteLog(err.Error(), "line "+strconv.Itoa(i+1)+": ")
			ok = false
			continue
		}
		raw := true
		datum, succ, annot, err := parseLine(line, raw)
		if err != nil {
			gentzen.WriteLog(err.Error(), "line "+strconv.Itoa(i+1)+": ")
			ok = false
		}

		s = append(s, datum+";"+succ+";"+replaceInfrules(annot))

	}
	return
}

func (d *console) lastLine() int {

	for i := len(d.Input) - 1; i >= 0; i-- {
		if len(d.Input[i]) != 0 {
			return i
		}
	}
	return len(d.Input) - 1
}

func length(l []inputLine) int {

	for n := len(l) - 1; n > 0; n-- {
		if len(l[n]) != 0 {
			return n
		}
	}
	return 0
}
func parseLine(l []string, raw bool) (datum, succ, annot string, err error) {

	p1, _, p2, _, p3 := parseNsplit(l, raw)

	datum = spaceyStringOf(p1)
	succ = spaceyStringOf(p2)
	annot = spaceyStringOf(p3)

	if strings.TrimSpace(datum+succ+annot) == "" {
		err = errors.New("need to have at least succedent and annotation")
		return
	}

	if succ == "" {
		err = errors.New("need to have succedent")
		return
	}

	var formula *gentzen.Node

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

	return
}

func isArgline(l []string) error {

	var err error

	s := strings.Join(l, " ")

	if strings.Count(s, `\vdash`) != 1 {
		err = errors.New("malformed derivation line")
		return err
	}

	if strings.Count(s, `\ldots`) != 1 {
		err = errors.New("malformed derivation line")
		return err
	}

	if strings.Index(s, `\vdash`) > strings.Index(s, `\ldots`) {
		err = errors.New("Malformed derivation line")
		return err
	}
	return err
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
func index(raw []string, s string) int {

	for i := 0; i < len(raw); i++ {
		if raw[i] == s {
			return i
		}
	}

	return -1
}

func splitLineRaw(raw []string) (datum, tstl, succ, dots, annot []string, err error) {

	if len(raw) == 0 {
		return
	}

	tst := index(raw, turnstile)
	sep := index(raw, separator)

	if sep > 0 && sep < tst {
		datum = raw
		return
	}
	if tst == -1 {
		datum = raw
		return
	}
	datum = raw[:tst]
	tstl = raw[tst : tst+1]
	if sep != -1 {
		succ = raw[tst+1 : sep]
		dots = raw[sep : sep+1]
		annot = raw[sep+1:]
		return
	}
	succ = raw[tst+1:]
	return
}
func isGreekLetter(s string) bool {

	for _, c := range greekBindings {
		if s == c[1] {
			return true
		}
	}
	return false
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
