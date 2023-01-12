package gentzen

import (
	"errors"
	"strings"
)

//InfixParser parses a formula in infix form given as a slice of strings.
//Each element of the slice should represent either a logical constant, bracket,
//sentence letter, predicate letter, or term letter (the last three with or w/o
//subscripts). Identities can be fed in the normal order (i.e. x=y rather than
//=xy).

func InfixParser(s []string) (n *Node, err error) {
	var ts tokenStr

	ts, err = tokenizeLatex(s)
	if err != nil {
		return
	}

	if !bracketsOK(ts) {
		err = errors.New("not enough brackets")
		return
	}

	n, err = parseLatex(ts)

	if err != nil {
		return
	}

	return

}

func bracketsOK(s tokenStr) bool {

	s = removeOuterBrackets(s)

	c := 0
	for _, e := range s {
		if e.isBinary() {
			c++
		}
		if e.isOpenb() {
			c--
		}
	}

	return (c < 2)
}

func tokenizeLatex(s []string) (tokenStr, error) {

	var err error
	var ts tokenStr
	for _, e := range s {
		e = strings.TrimSpace(e)
		var t token

		switch {

		case isOpenb(e):
			t.tokenType = tOpenb
			t.str = "("

		case isCloseb(e):
			t.tokenType = tCloseb
			t.str = ")"

		case e == `\supset`:
			t.tokenType = tCond
			t.str = lcond

		case e == `\wedge`:
			t.tokenType = tConj
			t.str = lconj

		case e == `\vee`:
			t.tokenType = tDisj
			t.str = ldisj

		case e == `\neg`:
			t.tokenType = tNeg
			t.str = lneg

		case e == `\forall`:
			t.tokenType = tUni
			t.str = luni

		case e == `\exists`:
			t.tokenType = tEx
			t.str = lex

		case e == `\lnec`:
			t.tokenType = tNec
			t.str = lnec

		case e == `\lpos`:
			t.tokenType = tPos
			t.str = lpos

		case !oPL:
			t.tokenType = tAtomicSentence
			t.str = e

		case isGreekFormulaVar(e):
			t.tokenType = tPredicate
			t.str = e

		case isFormulaSet(e):
			t.tokenType = tAtomicSentence
			t.str = e

		case isLowerCase(e[:1]):
			t.tokenType = tTerm
			t.str = e
		default:
			t.tokenType = tPredicate
			t.str = e
		}
		ts = append(ts, t)
	}

	//	ts2 = ts
	if oPL {
		//		ts2 = nil
		ts = fixBrackets(ts)
		ts = fixIdentity(ts)
		//	ts2, err = tokenizePLround2(ts)
		ts, err = tokenizePLround2(ts)
		if err != nil {
			logger.Print(err)
		}

	}
	return ts, err
}

func fixBrackets(ts tokenStr) tokenStr {

	var rts tokenStr
	var q = false

	for i := 0; i < len(ts); i++ {
		e := ts[i]

		if isGreekFormulaVar(e.str) {
			rts = append(rts, e)
			if i < len(ts)-1 {
				if ts[i+1].tokenType == tOpenb {
					q = true
					i++
					continue
				}
			}
			continue
		}
		if e.tokenType == tCloseb {
			if q == true {
				q = false
				continue
			}
		}

		rts = append(rts, e)
	}

	return rts
}

func fixIdentity(ts tokenStr) tokenStr {

	for i := 1; i < len(ts); i++ {
		if ts[i].str == `=` || ts[i].str == `≠` {
			tt := ts[i]
			ts[i].tokenType = ts[i-1].tokenType
			ts[i].str = ts[i-1].str
			ts[i-1].tokenType = tt.tokenType
			ts[i-1].str = tt.str
		}
	}
	var tn, tn2 token
	tn.tokenType = tNeg
	tn.str = lneg
	tn2.tokenType = tPredicate
	tn2.str = "="
	var tr tokenStr

	for _, e := range ts {
		if e.str != `≠` {
			tr = append(tr, e)
			continue
		}
		tr = append(tr, tn, tn2)
	}
	return tr
}

func isOpenb(s string) bool {

	for _, c := range brackets[1:] {
		if s == c[0] {
			return true
		}
	}
	return false
}

func isCloseb(s string) bool {

	for _, c := range brackets[1:] {
		if s == c[1] {
			return true
		}
	}
	return false
}

func parseLatex(ts tokenStr) (*Node, error) {

	var m Node
	n := &m

	mc, sub1, sub2, err := findMC(ts)

	if err != nil {
		return n, err
	}

	if mc.isAtomicSentence() {
		n.SetAtomic()
		n.SetFormula(mc.String())
		if oPL {
			n.predicateLetter = mc.predicate
			n.term = mc.term
		}
		return n, err
	}

	var ns1, ns2 *Node

	n.SetConnective(mc.tokenType.logicConstant())
	if mc.isQuantifier() {
		n.SetBoundVar(mc.variable)
	}
	ns1, err = parseLatex(sub1)
	if err != nil {
		return n, err
	}
	n.SetChild1(ns1)
	n.subnode1.parent = n

	if len(sub2) > 0 {
		ns2, err = parseLatex(sub2)
		if err != nil {
			return n, err
		}
		n.SetChild2(ns2)
		n.subnode2.parent = n
	}

	n, err = ParseStrict(n.String())

	return n, err
}

func _populateNodes(n *Node) {

	ns := getSubnodes(n)

	for _, e := range ns {
		e.SetFormula(e.String())
	}
	return
}

func findMC(ts tokenStr) (mc token, sub1, sub2 tokenStr, err error) {

	ts = removeOuterBrackets(ts)
	if ts[0].isQuantifier() {
	}
	if len(ts) == 0 {
		return
	}

	sub1, err = nextSentence(ts)
	if err != nil {
		return
	}

	if len(sub1) == len(ts) {
		mc = ts[0]
		sub1 = sub1[1:]
		return
	}

	if len(sub1)+2 > len(ts) {
		err = errors.New("malformed")
		return
	}

	if !ts[len(sub1)].isBinary() {
		err = errors.New("malformed")
		return
	}

	mc = ts[len(sub1)] //.tokenType.logicConstant()
	sub2, err = nextSentence(ts[len(sub1)+1:])

	if err != nil {
		logger.Print(err.Error())
		return
	}

	return
}

func removeOuterBrackets(ts tokenStr) tokenStr {
	if len(ts) < 3 {
		return ts
	}
	if ts[0].tokenType != tOpenb {
		return ts
	}

	c := findMatchingBracket(ts, 0)
	if c < len(ts)-1 {
		return ts
	}
	ts = ts[1:c]
	return removeOuterBrackets(ts)
}

func findMatchingBracket(ts tokenStr, start int) (end int) {
	var e token

	counter := 0

	for end, e = range ts {
		if e.isOpenb() {
			counter++
		}
		if e.isCloseb() {
			counter--
		}
		if counter == 0 {
			break
		}
	}
	if counter != 0 {
		return -1
	}

	return end
}

func nextSentence(ts tokenStr) (tn tokenStr, err error) {
	if len(ts) == 0 {
		return
	}
	if ts[0].isAtomicSentence() {
		if len(ts) == 1 {
			return ts[:1], err
		}
		if ts[1].isConnective() {
			return ts[:1], err
		}
	}

	if ts[0].isNeg() {
		ns, err := nextSentence(ts[1:])
		return append(ts[:1], ns...), err
	}

	if ts[0].isQuantifier() {
		ns, err := nextSentence(ts[1:])
		return append(ts[:1], ns...), err
	}

	if ts[0].isModalOperator() {
		ns, err := nextSentence(ts[1:])
		return append(ts[:1], ns...), err
	}
	if ts[0].isOpenb() {
		e := findMatchingBracket(ts, 0)
		if e == -1 {
			err = errors.New("Check your brackets")
			return ts, err
		}
		return ts[:findMatchingBracket(ts, 0)+1], err
	}

	err = errors.New("Something wrong. Check commas, brackets, etc." + ts[0].str)
	return
}
