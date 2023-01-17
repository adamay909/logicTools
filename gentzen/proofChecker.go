package gentzen

import (
	"errors"
	"strconv"
	"strings"
)

//arugment lines are given as semicolon separated lists:
//datum;succedent;lines used;inference rule
//
//inference rules are given by:ni,ne,ki,ke,di,de,ci,ce,a
//

type argLine struct {
	seq   sequent
	lines []int
	inf   string
}

const (
	ni  = "ni"
	ne  = "ne"
	ki  = "ki"
	ke  = "ke"
	di  = "di"
	de  = "de"
	ci  = "ci"
	ce  = "ce"
	a   = "a"
	ui  = "ui"
	ue  = "ue"
	ei  = "ei"
	ee  = "ee"
	ii  = "=i"
	ie  = "=e"
	li  = `li` //use l and m for necessity and possibility
	le  = `le`
	mi  = `mi`
	me  = `me`
	mli = `mli`
	pli = `pli`
	tli = `tli`
	mme = "mme"
	sc  = "sc"
	sl  = "sl"
)

var strictCheck bool

func checkDerivation(lines []string, offset int) bool {

	deriv, ok := getDerivation(lines, offset)
	if !ok {
		return false
	}

	if !lineRefsOK(deriv, offset) {
		return false
	}

	aE := func(v bool) {
		if !v {
			ok = false
		}
	}

	for n := range deriv {
		logger.SetPrefix("line " + strconv.Itoa(n+1) + ": ")
		aE(checkStep(getDerivTree(deriv, n)))

	}
	return ok
}

func parseArgline(s string) (al argLine, err error) {

	s = strings.ReplaceAll(s, " ", "")
	s = strings.ReplaceAll(s, "\t", "")

	fields := strings.Split(s, ";")
	//check we have enough fields
	if len(fields) < 3 {
		err = errors.New("you need: datum, succedent and at least one of: line references, inference rule")
		return
	}

	//check datum is ok
	datums := strings.Split(fields[0], ",")
	for _, d := range datums {
		if len(d) < 1 {
			continue
		}
		if isFormulaSet(d) {
			continue
		}

		if containsFormulaSet(d) {
			err = errors.New("datum: place holders for sets of formulas cannot appear inside a formula")
			return
		}
		_, err = ParseStrict(d)
		if err != nil {
			return
		}
	}

	//check succedent is ok
	if len(fields[1]) < 1 {
		err = errors.New("Not a sequent")
		return
	}
	if containsFormulaSet(fields[1]) {
		err = errors.New("Cannot have place holders for sets of formulas in succedent")
		return
	}
	_, err = ParseStrict(fields[1])
	if err != nil {
		return
	}

	al.seq.d = datum(strings.TrimSpace(fields[0]))
	al.seq.s = plshFormula(strings.TrimSpace(fields[1]))
	if len(fields) == 4 {
		al.inf = fields[3]
	}

	if len(strings.TrimSpace(fields[2])) == 0 {
		err = errors.New("you need: datum, succedent and at least one of: line references, inference rule")
		return
	}

	ln := strings.Split(fields[2], ",")

	var i int
	var e string

	for i = 0; i < len(ln); i++ {
		e = ln[i]
		n, err := strconv.Atoi(e)
		if err != nil {
			break
		}
		al.lines = append(al.lines, n)
	}

	if len(ln[i:]) > 1 {
		err = errors.New("You must have no more than one inference rule")
		return
	}
	if len(ln[i:]) != 0 {
		al.inf = strings.TrimSpace(ln[i])
	} else {
		al.inf = "rewrite"
	}
	return
}

func hasGreek(s string) bool {
	if !oPL {
		return strings.Contains(s, `\`)
	}
	r := []rune(s)

	for i := range r {
		for _, e := range greekUpperCaseLetters {
			if strings.HasPrefix(string(r[i:]), e) {
				return true
			}
		}
	}
	return false
}

func printArgLine(s string, m printMode) string {

	al, _ := parseArgline(s)
	datumstring := ""

	if len(al.seq.d) != 0 {

		datums := al.seq.datumSlice()
		for _, d := range datums {
			if d[:1] == `\` {
				datumstring = datumstring + runeOf(d, m) + `, `
			} else {
				datumstring = datumstring + printNodeInfix(Parse(d), m) + `, `
			}
		}
		datumstring = strings.TrimRight(datumstring, ", ")
	}

	succstring := printNodeInfix(Parse(al.seq.succedent().String()), m)
	annotation := ""
	if len(al.lines) > 0 {
		for _, i := range al.lines {
			annotation = annotation + strconv.Itoa(i) + `,`
		}
	}

	annotation = annotation + symb(al.inf, m)

	annotation = strings.TrimRight(annotation, ",")

	var resp string

	if m == mLatex {
		resp = `\ai{` + datumstring + `}{` + succstring + `}{` + annotation + `}` + "\n\n"
	}
	if m == mPlainText {
		resp = strings.ReplaceAll(datumstring+`⊢`+succstring+`...`+annotation, " ", "")
	}
	return resp
}

func symb(s string, m printMode) string {
	s = strings.TrimSpace(s)
	for _, i := range infRules {
		if i[0] == s {
			return i[m]
		}
	}
	return s
}

func runeOf[str ~string](s str, m printMode) string {
	var n int
	if m == mLatex {
		n = 1
	} else {
		n = 2
	}
	for _, e := range greekLCBindings {
		if e[1] == string(s) {
			return e[n]
		}
	}
	for _, e := range greekUCBindings {
		if e[1] == string(s) {
			return e[n]
		}
	}
	return string(s)
}

func isFormulaSet(s string) bool {

	for _, e := range greekUCBindings {
		if e[2] == s {
			return true
		}
	}
	return false
}

func containsFormulaSet(s string) bool {

	r := []rune(s)

	for _, e := range r {
		if isFormulaSet(string(e)) {
			return true
		}
	}
	return false
}
