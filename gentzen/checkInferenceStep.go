package gentzen

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/adamay909/logicTools/ju"
)

type derivationLine struct {
	seq             sequent
	lineNo          int
	supportingLines []int
	infrule         inferenceRule
}

type derivationNode = ju.Node[derivationLine]

var (
	errDerivation       error = errors.New("")
	errBadLine          error = errors.New("badly formatted: need sequent followed by annotation")
	errIllegalLineRef   error = errors.New("illegal line reference")
	errMultiInfRule     error = errors.New("more than one inference rule specified?")
	errUnknownInfRule   error = errors.New("unknown inference rule")
	errIllegalEmptyLine error = errors.New("an empty line in the middle")
)

type infErr struct {
	line    int
	message string
}

type infError []infErr

func (err infError) Error() string {
	w := new(strings.Builder)
	for i := range err {
		w.WriteString("line ")
		w.WriteString(strconv.Itoa(err[i].line))
		w.WriteString(": ")
		w.WriteString(err[i].message)
		if i < len(err)-1 {
			w.WriteString("\n")
		}
	}
	return w.String()
}

func (err infError) add(e infErr) infError {
	err = append(err, e)
	return err
}

func newinferr(l int, m string) (e infErr) {
	e.line = l
	e.message = m
	return
}

// CheckDerivation is the entry point for proof checking frontends. Mode is the
// mode in which formulas are formatted, offset is the offset for line numbers.
// offset==1 means line numbering at 1. err shows any errors found during the checking.
func CheckDerivation(s []string, mode PrintMode, offset int) error {

	var derivation []derivationLine
	var err infError
	var derr error
	derivation, derr = toDerivation(s, mode, offset)

	if derr != nil {
		return derr
	}

	var series []sequent

	for ln, dline := range derivation {
		series = nil
		for _, i := range dline.supportingLines {
			series = append(series, derivation[i-offset].seq)
		}
		series = append(series, dline.seq)

		fmt.Print("\n\n\n")
		fmt.Println("check step:\n", printSequentSeries(normalizeDerivation(series)))

		e := inferenceStepSuccessful(series, dline.infrule)
		if e != nil {
			err = err.add(newinferr(ln+offset, e.Error()))
			return err
		}
	}

	return nil

}

func inferenceStepSuccessful(series []sequent, ir inferenceRule) error {

	var err error

	if ir.name == "rewrite" {
		return isRewriteStep(series)
	}

	if len(series) != len(ir.premises)+1 {
		err = errors.New(ir.fullName + " depends on " + strconv.Itoa(len(ir.premises)) + " lines")
		return err
	}

	if ir.isPLrule() {
		if !oPL {
			return errors.New(ir.fullName + " cannot be used with sentential logic")
		}
		return inferenceStepSuccessfulPL(series, ir)
	}

	series = normalizeDerivation(series, "s", "/G")

	if len(series) == 1 {
		if ok := matchInfrule(series, ir); ok {
			return nil
		}
		err = errors.New("step does not work")
		return err
	}

	//got through all possible permutations of the lines to be checked
	for _, prem := range permutations(series[:len(series)-1]) {
		fmt.Println()
		deriv := make([]sequent, 0, len(series))
		deriv = append(deriv, prem...)
		deriv = append(deriv, series[len(series)-1])
		fmt.Println("check permutation:\n", printSequentSeries(deriv))
		if matchInfrule(deriv, ir) {
			return nil
		}
	}

	err = errors.New("step does not work")
	return err
}

func toDerivation(s []string, mode PrintMode, offset int) ([]derivationLine, error) {

	var d []derivationLine
	var err infError

	lastEmptyLine, lastNonEmpty := -1, 0
	for i := len(s) - 1; i > -1; i-- {
		if len(strings.TrimSpace(s[i])) != 0 {
			lastNonEmpty = i
			break
		}
	}

	for i := range s[:lastNonEmpty+1] {
		if len(strings.TrimSpace(s[i])) == 0 {
			lastEmptyLine = i
		}
	}

	if lastEmptyLine > -1 && (lastNonEmpty > lastEmptyLine) {
		err = err.add(newinferr(lastEmptyLine+offset, "an empty line in the middle"))
		return nil, err
	}

	for i := range s[:lastNonEmpty+1] {
		dl, derr := toDerivationLine(s[i], mode)
		if derr != nil {
			err = err.add(newinferr(i+offset, derr.Error()))
			return d, err
		}
		for _, lr := range dl.supportingLines {
			if lr < offset || lr > i+offset {
				err = err.add(newinferr(i+offset, "illegal reference to line "+strconv.Itoa(lr)))
				return d, err
			}
		}
		d = append(d, dl)
	}
	return d, nil
}

func toDerivationLine(s string, mode PrintMode) (derivationLine, error) {

	var d derivationLine
	var err error

	parts := strings.Split(s, ".")
	if len(parts) != 2 {
		err = errBadLine
		return d, err
	}

	if len(parts[0]) == 0 || len(parts[1]) == 0 {
		err = errBadLine
		return d, err
	}

	d.seq, err = mkSequent(parts[0], mode)
	if err != nil {
		return d, err
	}

	annot := strings.Split(parts[1], ",")

	for k, a := range annot {
		ln, converr := strconv.Atoi(a)
		if converr == nil {
			d.supportingLines = append(d.supportingLines, ln)
			if k == len(annot)-1 {
				d.infrule = infrule["rewrite"]
			}
			continue
		}
		if k < len(annot)-1 {
			err = errMultiInfRule
			return d, err
		}
		if _, ok := infrule[a]; !ok {
			err = errUnknownInfRule
			return d, err
		}

		d.infrule = infrule[a]
	}
	return d, err
}

// This will normalize predicate logic sequents into sentential logic
// sequents. I.e., you will lose information about sub-sentential structures.
// If you need subsentential information, which you will for first-order specific rules,
// use normalizeDerivationPL.
func normalizeDerivation(series []sequent, l ...string) []sequent {

	p := "s_"
	if len(l) > 0 {
		p = l[0] + "_"
	}

	s := "/L_"
	if len(l) > 1 {
		s = l[1] + "_"
	}

	resp := make([]sequent, 0, len(series))
	sets := make(map[string]string)
	atoms := make(map[string]string)

	for i := range series {
		var nseq sequent
		for _, d := range series[i].front {
			n := Parse(d, allowGreekUpper)
			if n.Check(isSet) {
				if _, ok := sets[d]; !ok {
					sets[d] = s + strconv.Itoa(len(sets)+1)
				}
				nseq.front = append(nseq.front, sets[d])
				continue
			}
			for _, a := range sententialAtoms(n) {
				if _, ok := atoms[a.String()]; !ok {
					atoms[a.String()] = p + strconv.Itoa(len(atoms)+1)
				}
			}
			nseq.front = append(nseq.front, replaceAtomicWith(d, atoms))
		}
		n := Parse(series[i].back, !allowGreekUpper)
		for _, a := range sententialAtoms(n) {
			if _, ok := atoms[a.String()]; !ok {
				atoms[a.String()] = p + strconv.Itoa(len(atoms)+1)
			}
		}
		nseq.back = replaceAtomicWith(series[i].back, atoms)
		resp = append(resp, nseq)
	}

	return resp
}

// check if series is an instance of ir. Assumes series is already normalized.
// does not check for permutation of premises. series must have at least the
// right number of lines.
func matchInfrule(series []sequent, ir inferenceRule) bool {

	if ir.isPLrule() {
		return false
	}

	ropl := oPL
	defer func() {
		oPL = ropl
	}()
	SetPL(false)

	for _, c := range ir.conclusion {

		var canonical []sequent

		canonical = append(canonical, ir.premises...)
		canonical = append(canonical, c)

		canonical = normalizeDerivation(canonical, "t", "/L")
		fmt.Println("canonical:\n", printSequentSeries(canonical))
		if matchPattern(series, canonical, ir.isIntroductionRule(), []string{}) {
			return true
		}
	}
	return false
}

// check if the front of s2 can be gotten by rewriting front of s1 (i.e., reorder or removal of duplicates)
func isRewriteFront(s1, s2 sequent) bool {
	//back formulas must be identical
	if s1.back != s2.back {
		return false
	}
	// all elements of s1 must be present in s2 at least once and vice versa
	for i := range s1.front {
		if !slices.Contains(s2.front, s1.front[i]) {
			return false
		}
	}
	for i := range s2.front {
		if !slices.Contains(s1.front, s2.front[i]) {
			return false
		}
	}
	// make sure no new duplicates were added to s1
	dfront := make([]string, 0, len(s2.front))
	dfront = append(dfront, s2.front...)
	for i := range s1.front {
		dfront = slices.DeleteFunc(dfront, deleteOnceFunc(dfront, s1.front[i]))
	}
	return len(dfront) == 0
}

func isRewriteStep(series []sequent) error {
	if len(series) != 2 {
		return errors.New("sequent rewrite depends on one line")
	}

	if isRewriteFront(series[0], series[1]) {
		return nil
	}

	if isSubslice(series[1].front, series[0].front) {
		return nil
	}

	return errors.New("not a rewrite")
}

/*
check if series is an instance of ir where ir must be a frontManipulation rule. Assumes series is already normalized does not check for permutation of premises. series must have at least the	right number of lines.
*/
func matchInfruleFM(series []sequent, ir inferenceRule) bool {

	if ir.name == "rewrite" {
		for i := range series[0].front {
			if !slices.Contains(series[1].front, series[0].front[i]) {
				return false
			}
		}
		return true
	}

	// no other front manipulation rules implemented yet.
	return false
}

func toSequentSeries(d []derivationLine) []sequent {
	r := make([]sequent, 0, len(d))
	for i := range d {
		r = append(r, d[i].seq)
	}
	return r
}

func printDerivation(derivation []derivationLine) string {
	w := new(strings.Builder)
	for i := range derivation {
		w.WriteString(derivation[i].seq.frontString())
		w.WriteString(" ⊢ ")
		w.WriteString(derivation[i].seq.backString())
		w.WriteString("...")
		w.WriteString(derivation[i].infrule.name)
		w.WriteString("\n")
	}
	return w.String()
}

func printSequentSeries(series []sequent) string {

	s := make([]string, 0, len(series))

	for i := range series {
		s = append(s, series[i].String())
	}

	return strings.Join(s, "\n")
}

func sententialAtoms(n *Node) []*Node {

	var resp []*Node

	ingressFunc := func(e *Node) {

		if !e.Check(isSententialAtom) {
			return
		}
		for f := e; f != nil; f = f.Parent() {
			if f.HasFlag("done") {
				return
			}
		}
		resp = append(resp, e)
		e.SetFlag("done")
	}

	donothing := func(e *Node) {}

	n.Walk(ingressFunc, donothing, donothing)

	return resp
}

func normalizeDerivationPL(series []sequent, l ...string) []sequent {

	s := "/L_"
	if len(l) > 2 {
		s = l[2] + "_"
	}

	resp := make([]sequent, 0, len(series))

	sets := make(map[string]string)
	preds := make(map[string]string)
	terms := make(map[string]string)

	//fmt.Println("check series\n", printDeriv(series))
	for i := range series {
		var nseq sequent
		for _, d := range series[i].front {
			n := Parse(d, allowGreekUpper)
			if n.Check(isSet) {
				if _, ok := sets[d]; !ok {
					sets[d] = s + strconv.Itoa(len(sets)+1)
				}
				nseq.front = append(nseq.front, sets[d])
				continue
			}
			safifyNames(n, preds, terms)
			toNormalNames(n, l...)
			nseq.front = append(nseq.front, n.String())
		}
		n := Parse(series[i].back, !allowGreekUpper)
		safifyNames(n, preds, terms)
		toNormalNames(n, l...)
		nseq.back = n.String()
		resp = append(resp, nseq)
	}

	return resp
}

func inferenceStepSuccessfulPL(series []sequent, ir inferenceRule) error {

	fmt.Println("PL rule checking")

	if !ir.isPLrule() {
		panic("do not call inferenceStepSuccessfulPL directly!")
	}

	series = normalizeDerivationPL(series, "G", "c", "/G")

	//got through all possible permutations of the lines to be checked
	for _, prem := range permutations(series[:len(series)-1]) {
		deriv := make([]sequent, 0, len(series))
		deriv = append(deriv, prem...)
		deriv = append(deriv, series[len(series)-1])
		fmt.Println("check permutation:", printSequentSeries(deriv))
		if matchInfrulePL(deriv, ir) {
			return nil
		}
	}

	return errors.New("step does not work")
}

// check if series is an instance of ir. Assumes series is already normalized.
// does not check for permutation of premises. series must have at least the
// right number of lines.
func matchInfrulePL(series []sequent, ir inferenceRule) bool {

	if !ir.isPLrule() {
		panic("do not call matchinfrulePL directly!")
	}

	for _, c := range ir.conclusion {

		var canonical []sequent

		canonical = append(canonical, ir.premises...)
		canonical = append(canonical, c)

		var exclude []string

		if strings.Contains(ir.spec, "constants unique") {
			exclude = getConstants(series)
		}
		//exclude = []string{}

		//		fmt.Println("constants to exclude:", exclude)

		canonical = normalizeDerivationPL(canonical, "F", "k", "/L")

		if matchPattern(series, canonical, ir.isIntroductionRule(), exclude) {
			return true
		}

	}

	return false
}

/*PrintDerivation formats lines a derivation. informat species how the input is formatted. O_PlainASCII, O_ProofChecker imply infix formatting of input. outputFormat species the formatting of the output.*/
func PrintDerivation(lines []string, offset int, informat PrintMode, outputFormat PrintMode) string {

	derivation, derr := toDerivation(lines, informat, offset)

	if derr != nil {
		panic("input is not a derivation")
	}
	w := new(strings.Builder)

	if outputFormat == O_Latex {
		return printDerivLatex(derivation, offset)
	}

	for i := range derivation {

		w.WriteString(strconv.Itoa(i + offset))
		w.WriteString(`. `)
		w.WriteString(printDerivline(derivation[i], outputFormat))
		w.WriteString("\n")

	}

	return w.String()
}

func printDerivLatex(derivation []derivationLine, offset int) string {

	w := new(strings.Builder)

	w.WriteString(`\begin{argumentN}[`)
	w.WriteString(strconv.Itoa(offset))
	w.WriteString("]\n")
	w.WriteString(`%generated by gentzen`)
	w.WriteString("\n\n")

	for i := range derivation {

		w.WriteString(printDerivline(derivation[i], O_Latex))
		w.WriteString("\n\n")
	}

	w.WriteString(`\end{argumentN}`)
	w.WriteString("\n\n")

	return w.String()
}

func printDerivline(l derivationLine, m PrintMode) string {

	w := new(strings.Builder)

	var frontString, backString, annotString string

	for i, d := range l.seq.front {
		w.WriteString(StringF(Parse(d, allowGreekUpper), m))
		if i < len(l.seq.front)-1 {
			w.WriteString(`,`)
		}
	}
	frontString = w.String()

	w.Reset()

	backString = StringF(Parse(l.seq.back, !allowGreekUpper), m)

	for i, n := range l.supportingLines {
		w.WriteString(strconv.Itoa(n))
		if i < len(l.supportingLines)-1 {
			w.WriteString(`,`)
		}
	}
	if strings.TrimSpace(l.infrule.displayName) != "" {
		if len(l.supportingLines) != 0 {
			w.WriteString(",")
		}
		if m == O_Latex {
			w.WriteString(l.infrule.latexName)
		} else {
			w.WriteString(l.infrule.displayName)
		}
	}
	annotString = w.String()

	w.Reset()

	if m == O_Latex {
		w.WriteString(`\ai{`)
		w.WriteString(frontString)
		w.WriteString(`}{`)
		w.WriteString(backString)
		w.WriteString(`}{`)
		w.WriteString(annotString)
		w.WriteString("}")
	}

	if m == O_PlainText || m == O_ProofChecker {
		w.WriteString(frontString)
		w.WriteString(`⊢`)
		w.WriteString(backString)
		w.WriteString("...")
		w.WriteString(annotString)
	}

	if m == O_PlainASCII || m == O_English || m == O_Polish {
		w.WriteString(frontString)
		w.WriteString(":")
		w.WriteString(backString)
		w.WriteString("...")
		w.WriteString(annotString)
	}

	return w.String()

}

func getDerivationTree(deriv []derivationLine, offset int) *derivationNode {

	d := new(derivationNode)

	d.Val = deriv[len(deriv)-1]

	for _, l := range d.Val.supportingLines {

		d.AddChild(getDerivationTree(deriv[:l-offset+1], offset))

	}

	return d
}

// PrintDerivationTree outputs Latex code for representing the derivation given by s in tree form.
func PrintDerivationTree(s []string, mode PrintMode, offset int) (string, error) {

	w := new(strings.Builder)

	ingressFunc := func(e *derivationNode) {

		w.WriteString(`[ \p{`)
		w.WriteString(e.Val.seq.stringF(O_Latex))
		w.WriteString(`}`)
		w.WriteString("\n")
	}

	pivotFunc := func(m *derivationNode) {

		w.WriteString(` ] `)
		w.WriteString("\n")
	}

	egressFunc := func(m *derivationNode) {

		if !m.IsLastSibling() {
			return
		}

		w.WriteString(` ] `)
		w.WriteString("\n")
	}

	derivation, derr := toDerivation(s, mode, offset)

	if derr != nil {
		return "", derr
	}

	dn := getDerivationTree(derivation, offset)
	w.WriteString(`\begin{forest}{for tree={grow=north}}`)
	w.WriteString("\n%generated by gentzen\n")

	dn.Walk(ingressFunc, pivotFunc, egressFunc)

	w.WriteString(`\end{forest}`)
	w.WriteString("\n")

	return w.String(), nil
}
