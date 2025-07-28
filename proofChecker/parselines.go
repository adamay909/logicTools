package main

import (
	"errors"
	"fmt"
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
		debug("parsing :", strings.Join(line, "; "))

		datum, succ, annot, err := parseLine(line, raw)
		if err != nil {
			gentzen.WriteLog(err.Error(), "line "+strconv.Itoa(i+1)+": ")
			ok = false
		}

		s = append(s, datum+";"+succ+";"+replaceInfrules(annot))

		fmt.Println("got line: ", s[len(s)-1])

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

func subscript(s string) string {

	if !strings.Contains(s, "_") {
		return ""
	}

	i := strings.Index(s, "_")

	return strings.TrimSpace(s[i+1:])

}

func removeSubscript(s string) string {

	return strings.TrimSuffix(s, "_"+subscript(s))
}

func fixSubscripts(l []string) (o []string, err error) {

	if len(l) == 0 {
		return
	}

	o = append(o, l[0])

	for i := 1; i < len(l); i++ {

		if strings.HasPrefix(l[i], `<sub>`) {
			o[len(o)-1] = o[len(o)-1] + "_" + strings.TrimSuffix(strings.TrimPrefix(l[i], `<sub>`), `</sub>`)
			continue
		}
		o = append(o, l[i])
	}
	return
}

func parseLine(l []string, raw bool) (datum, succ, annot string, err error) {

	lf, _ := fixSubscripts(l)

	p1, _, p2, _, p3 := parseNsplit(lf, raw)

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

		if len(e) == 0 {
			continue
		}
		debug("datum check: ", e)
		// formula, err = gentzen.InfixParser(tk(e))
		formula, err = gentzen.ParseInfix(e, true)

		if err != nil {
			fmt.Println("initial parse", err)
			return
		}
		datum = datum + formula.String() + ","
	}
	datum = strings.TrimSuffix(datum, ",")

	debug("datum is ", datum)

	//Deal with succedent

	formula, err = gentzen.ParseInfix(succ, false)
	if err != nil {
		return
	}
	succ = formula.String()

	debug("succ is ", succ)

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
			if subscript(e) != "" {
				e = removeSubscript(e) + "_" + subscript(e)
			}
			t = append(t, e)
		}
	}

	return t

}

func replaceInfrules(s string) string {
	infrules := [][2]string{
		[2]string{`-E`, `ne`},
		[2]string{`-I`, `ni`},
		[2]string{`vE`, `de`},
		[2]string{`vI`, `di`},
		[2]string{`^E`, `ke`},
		[2]string{`^I`, `ki`},
		[2]string{`>E`, `ce`},
		[2]string{`>I`, `ci`},
		[2]string{`A`, `a`},
		[2]string{`M`, `m`},
		[2]string{`UE`, `ue`},
		[2]string{`UI`, `ui`},
		[2]string{`XE`, `ee`},
		[2]string{`XI`, `ei`},
		[2]string{`=E`, `=e`},
		[2]string{`=I`, `=i`},
		//	[2]string{`\lnecE`, `le`},
		//	[2]string{`\lnecI`, `li`},
		//	[2]string{`S5\lnecI`, `mli`},
		//	[2]string{`S4\lnecI`, `pli`},
		//	[2]string{`T\lnecI`, `tli`},
		//	[2]string{`\lposE`, `me`},
		//	[2]string{`S5\lposE`, `mme`},
		//	[2]string{`\lposI`, `mi`},
		//	[2]string{`SC`, `sc`},
		//	[2]string{`logic`, `sl`},
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
