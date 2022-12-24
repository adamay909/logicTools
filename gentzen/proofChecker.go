package gentzen

import (
	"errors"
	"log"
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
	mme = "mme"
	sc  = "sc"
	sl  = "sl"
)

type infRule struct {
	abbr     string
	fullName string
	latex    string
	mathjax  string
	text     string
	spec     int
	rule     func() bool
}

var checkLog strings.Builder
var logger *log.Logger
var strictCheck bool

func checkDerivation(lines []string, offset int) bool {
	var al []argLine
	var hasE bool

	//record if there is an error
	aE := func(v bool) {
		if v {
			hasE = true
		}
	}

	hasE = false
	for k, i := range lines {
		logger.SetPrefix("line " + strconv.Itoa(k+1) + ": ")
		pl, err := parseArgline(i)
		if err != nil {
			logger.Print(err.Error())
			hasE = true
		}
		al = append(al, pl)
	}
	if hasE {
		return false
	}

	for i, l := range al {
		logger.SetPrefix("line " + strconv.Itoa(i+1) + ": ")

		//check the line references first
		if !checkLineRef(l.inf, i+offset, offset, l.lines) {
			hasE = true
			continue
		}

		//if we are dealing with a derived rule
		if oDR {
			if strings.HasSuffix(l.inf, "R") {
				aE(!derivR(l.inf, al[l.lines[0]-offset].seq, l.seq))
				continue
			}
		}

		switch l.inf {
		case a: //assumption
			aE(!assumption(l.seq))

		case ki: //conjunction intro
			aE(!conjI(al[l.lines[0]-offset].seq, al[l.lines[1]-offset].seq, l.seq))

		case ke: //conjunction elim
			aE(!conjE(al[l.lines[0]-offset].seq, l.seq))

		case di: //disjunction intro
			aE(!disjI(al[l.lines[0]-offset].seq, l.seq))

		case de: //disjunction elim
			aE(!disjE(al[l.lines[0]-offset].seq, al[l.lines[1]-offset].seq, al[l.lines[2]-offset].seq, l.seq))

		case ci: //conditional intro
			aE(!condI(al[l.lines[0]-offset].seq, l.seq))

		case ce: //conditional elim
			aE(!condE(al[l.lines[0]-offset].seq, al[l.lines[1]-offset].seq, l.seq))

		case ni: //negation intro
			aE(!negI(al[l.lines[0]-offset].seq, al[l.lines[1]-offset].seq, l.seq))

		case ne: //negation elim
			aE(!negE(al[l.lines[0]-offset].seq, l.seq))

		case ue: //universal elim
			aE(!uniE(al[l.lines[0]-offset].seq, l.seq))

		case ui: //universal intro
			aE(!uniI(al[l.lines[0]-offset].seq, l.seq))

		case ei: //existential intro
			aE(!exI(al[l.lines[0]-offset].seq, l.seq))

		case ee: //existential elimo
			aE(!exE(al[l.lines[0]-offset].seq, al[l.lines[1]-offset].seq, l.seq))

		case ii: //identity introduction
			aE(!idI(l.seq))

		case ie: //identity introduction
			aE(!idE(l.seq))

		case li: //necessity introduction
			aE(!necI(al[l.lines[0]-offset].seq, l.seq))

		case pli: //necessity introduction
			aE(!necI_S4(al[l.lines[0]-offset].seq, l.seq))

		case mli: //necessity introduction
			aE(!necI_S5(al[l.lines[0]-offset].seq, l.seq))

		case le: //necessity elim
			aE(!necE(al[l.lines[0]-offset].seq, l.seq))

		case mi: //possibility intro
			aE(!posI(al[l.lines[0]-offset].seq, l.seq))

		case me: //possibility elim
			aE(!posE(al[l.lines[0]-offset].seq, al[l.lines[1]-offset].seq, l.seq))

		case mme: //possibility elim
			aE(!posE_S5(al[l.lines[0]-offset].seq, al[l.lines[1]-offset].seq, l.seq))

		case sc: //possibility elim
			aE(!scopeReplacement(al[l.lines[0]-offset].seq, l.seq))

		case sl: //possibility elim
			aE(!sententialLogic(al[l.lines[0]-offset].seq, l.seq))

		case "premise": //premise

		case "": //sequent rewrite
			aE(!seqRewrite(l.seq, al[l.lines[0]-offset].seq, l.lines[0]))

		default: //check if we are dealing with a theorem
			aE(!oTHM)
			aE(!theorem(l.seq, l.inf))
		}
	}
	return !hasE
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

func lineSpec(infRule string) int {
	switch infRule {
	case a:
		return 0
	case ni:
		return 2
	case ne:
		return 1
	case ki:
		return 2
	case ke:
		return 1
	case di:
		return 1
	case de:
		return 3
	case ci:
		return 1
	case ce:
		return 2
	case ue:
		return 1
	case ui:
		return 1
	case ee:
		return 2
	case ei:
		return 1
	case ii:
		return 0
	case ie:
		return 0
	case le:
		return 1
	case li:
		return 1
	case mli:
		return 1
	case pli:
		return 1
	case me:
		return 2
	case mme:
		return 2
	case mi:
		return 1
	case sc:
		return 1
	case sl:
		return 1
	case "premise":
		return 0
	case "":
		return 1
	case "theorem":
		return 0
	default:
		return -1
	}
}

func checkLineRef(infRule string, cur int, offset int, lines []int) bool {

	for n := range lines {
		if lines[n] >= cur || lines[n] < offset {
			logger.Print("illegel reference to line ", lines[n])
			return false
		}
	}
	if oDR {
		if strings.HasSuffix(infRule, "R") {
			if len(lines) != 1 {
				logger.Print("derived rule must refer to one other line")
				return false
			}
			return true
		}
	}

	thm := theorems
	if oML {
		thm = append(thm, modalTheorems...)
	}

	if oPL {
		thm = append(thm, quantifierRules...)
	}

	if oTHM {
		for i := range thm {
			if infRule == thm[i][0] || infRule == thm[i][1] {
				infRule = "theorem"
				break
			}
		}
	}

	if lineSpec(infRule) == -1 {
		logger.Print("unknown inference rule or theorem: ", infRule)
		return false
	}

	if len(lines) != lineSpec(infRule) {
		if lineSpec(infRule) != 1 {
			logger.Print(fullName(infRule), " should refer to ", lineSpec(infRule), " lines")
		} else {
			logger.Print(fullName(infRule), " should refer to ", lineSpec(infRule), " line")
		}

		return false
	}

	return true
}

func fullName(i string) string {
	switch i {
	case a:
		return "Assumption"
	case ne:
		return "Negation Elimination"
	case ni:
		return "Negation Introduction"
	case de:
		return "Disjunction Elimination"
	case di:
		return "Disjunction  Introduction"
	case ke:
		return "Conjunction Elimination"
	case ki:
		return "Conjunction Introduction"
	case ce:
		return "Conditional Elimination"
	case ci:
		return "Conditional Introduction"
	case ue:
		return "Universal Quantifier Elimination"
	case ui:
		return "Universal Quantifier Introduction"
	case ee:
		return "Existential Quantifier Elimination"
	case ei:
		return "Existential Quantifier Introduction"
	case ie:
		return "Identity Elimination"
	case ii:
		return "Identity Introduction"
	case li:
		return "Necessity Introduction"
	case mli:
		return "Metaphysical Necessity Introduction"
	case pli:
		return "Physical Necessity Introduction"
	case le:
		return "Necessity Elimination"
	case mi:
		return "Possibility Introduction"
	case me:
		return "Possibility Elimination"
	case mme:
		return "Metaphysical Possibility Introduction"
	case sc:
		return "Scope Replacement"
	case sl:
		return "Sentential Logic"
	default:
		return i
	}
}

func isTheorem(s sequent) bool {
	return s.d == ""
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
