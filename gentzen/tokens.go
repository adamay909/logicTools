package gentzen

import (
	"errors"
	"strings"
)

type tokenID int

type token struct {
	tokenType tokenID
	str       string
	variable  string   //for PL
	predicate string   //for PL
	term      []string //for PL

}

type tokenStr []token

const (
	tConj tokenID = 1 << iota
	tDisj
	tNeg
	tCond
	tAtomicSentence
	tUni
	tEx
	tNec
	tPos
	tPredicate
	tTerm
	tOpenb
	tCloseb
)

const (
	tConnective tokenID = tConj | tDisj | tNeg | tCond | tUni | tEx | tNec | tPos
)

const (
	tQuantifier tokenID = tUni | tEx
)

const (
	tModalOperator tokenID = tNec | tPos
)

const (
	tUnary tokenID = tNeg | tUni | tEx | tNec | tPos
)

const (
	tBinary tokenID = tConj | tCond | tDisj
)

func (t token) isQuantifier() bool {
	return t.tokenType == tUni || t.tokenType == tEx
}

func (t token) isModalOperator() bool {
	return t.tokenType == tNec || t.tokenType == tPos
}

func (t token) isPredicate() bool {
	return t.tokenType == tPredicate
}

func (t token) isTerm() bool {
	return t.tokenType == tTerm
}

func (t token) isConnective() bool {
	if oPL {
		return (t.tokenType & tConnective) != 0
	}
	return t.tokenType != tAtomicSentence
}

func (t token) isAtomicSentence() bool {
	return t.tokenType == tAtomicSentence
}

func (t token) String() string {
	if isGreekLowerCase(t.str) {
		return greekCharOf(t.str)
	}
	if t.isQuantifier() {
		return t.str + t.variable
	}
	return t.str
}

func (t token) StringMathJax() string {
	if t.isConnective() {
		return logicalConstant(t.str).StringMathJax()
	}
	return t.str
}

func (t token) isConj() bool {
	return t.tokenType == tConj
}

func (t token) isDisj() bool {
	return t.tokenType == tDisj
}

func (t token) isNeg() bool {
	return t.tokenType == tNeg
}

func (t token) isCond() bool {
	return t.tokenType == tCond
}

func (t token) isKCA() bool {
	return (t.tokenType & (tConj | tDisj | tCond)) != 0
}

func (t token) isUni() bool {
	return t.tokenType == tUni
}

func (t token) isEx() bool {
	return t.tokenType == tEx
}

func (t token) isNec() bool {
	return t.tokenType == tNec
}

func (t token) isPos() bool {
	return t.tokenType == tPos
}
func (t token) isUnary() bool {

	return t.isNeg() || t.isUni() || t.isEx() || t.isNec() || t.isPos()
}

func (t token) isBinary() bool {
	return t.isConj() || t.isDisj() || t.isCond()
}

func (t token) isOpenb() bool {
	return t.tokenType == tOpenb
}

func (t token) isCloseb() bool {
	return t.tokenType == tCloseb
}

func (t token) boundVariable() string {
	if !t.isQuantifier() {
		return ""
	}
	return t.variable
}

func (t token) terms() []string {
	return t.term
}

func (t tokenID) logicConstant() logicalConstant {
	switch t {
	case tNeg:
		return neg
	case tConj:
		return conj
	case tDisj:
		return disj
	case tCond:
		return cond
	case tUni:
		return uni
	case tEx:
		return ex
	case tNec:
		return nec
	case tPos:
		return pos
	default:
		return logicalConstant("")
	}
}
func tokenize(s string) (t tokenStr, err error) {

	//we igonre all spaces and tabs, etc
	s = cleanString(s)

	var e token
	for len(s) > 0 {
		e, s = nextToken(s)
		t = append(t, e)
	}

	//we are done if are not doing predicate logic
	if !oPL {

		return t, err
	}

	t, err = tokenizePLround2(t)

	return t, err
}

func (t tokenStr) String() string {

	var s string

	for _, e := range t {
		s = s + e.String()
	}
	return s
}

// nextToken returns the next token in s and
// returns the remainder string r. Works with
// predicate logic.
func nextToken(s string) (t token, r string) {

	pos := 1
	sr := []rune(s)

	switch {
	case string(sr[:pos]) == ldisj:
		t.tokenType = tDisj
	case string(sr[:pos]) == lconj:
		t.tokenType = tConj
	case string(sr[:pos]) == lcond:
		t.tokenType = tCond
	case string(sr[:pos]) == lneg:
		t.tokenType = tNeg
	case string(sr[:pos]) == luni:
		t.tokenType = tUni
	case string(sr[:pos]) == lex:
		t.tokenType = tEx
	case string(sr[:pos]) == lnec:
		t.tokenType = tNec
	case string(sr[:pos]) == lpos:
		t.tokenType = tPos
	case isGreekFormulaVar(string(sr[:pos])):
		t.tokenType = tPredicate
	case isFormulaSet(string(sr[:pos])):
		t.tokenType = tAtomicSentence
		if len(s) > 1+pos {
			if string(sr[pos:pos+1]) == `_` {
				pos = pos + 2
			}
		}
	case isLowerCase(string(sr[:pos])):
		t.tokenType = tTerm
		if len(s) > 1+pos {
			if string(sr[pos:pos+1]) == `_` {
				pos = pos + 2
			}
		}
	default:
		t.tokenType = tPredicate
		if len(sr) > 1+pos {
			if string(sr[pos:pos+1]) == `\` {
				pos = pos + 2
			}
			if len(sr) > 1+pos {
				if string(sr[pos:pos+1]) == `_` {
					pos = pos + 2
				}
			}
		}

	}
	if !oPL {
		if t.tokenType == tTerm || t.tokenType == tPredicate {
			t.tokenType = tAtomicSentence
		}
		if t.tokenType == tUni || t.tokenType == tEx {
			t.tokenType = tAtomicSentence
		}
	}
	t.str = string(sr[:pos])
	r = string(sr[pos:])

	return t, r
}

func cleanString(s string) string {

	s = strings.ReplaceAll(s, " ", "")
	s = strings.ReplaceAll(s, "\t", "")
	s = strings.ReplaceAll(s, "(", "")
	s = strings.ReplaceAll(s, ")", "")

	for _, g := range greekLCBindings {
		s = strings.ReplaceAll(s, g[0], g[2])
	}

	return s
}

func isLowerCase(s string) bool {
	if len(s) > 1 {
		return false
	}

	if s == "=" {
		return false
	}

	if s == strings.ToUpper(s) {
		return false
	}

	return true
}

func isGreekLowerCase(s string) bool {

	for _, e := range greekLCBindings {
		if s == e[1] {
			return true
		}
	}

	return false
}

func isGreekFormulaVar(s string) bool {

	for _, e := range greekLCBindings {
		if s == e[2] {
			return true
		}
	}

	return false
}

func tokenizePLround2(t tokenStr) (tokenStr, error) {
	var t2 tokenStr
	var e token
	var err error
	for i := 0; i < len(t); i++ {
		e = t[i]

		var n token
		if e.isQuantifier() {
			if i > len(t)-2 {
				err = errors.New("quantifier without variable")
				return t, err
			}
			if !t[i+1].isTerm() {
				err = errors.New("quantifier without variable")
				return t, err
			}
			n.tokenType = e.tokenType
			n.str = e.str
			n.variable = t[i+1].str
			t2 = append(t2, n)
			i++
			continue
		}
		if e.isAtomicSentence() {
			t2 = append(t2, e)
			continue
		}

		if e.isPredicate() {
			if isGreekFormulaVar(e.str) {
				if i == len(t)-1 {
					n.tokenType = tAtomicSentence
					n.predicate = e.str
					n.str = e.str
					t2 = append(t2, n)
					continue
				}
				if !t[i+1].isTerm() {
					n.tokenType = tAtomicSentence
					n.predicate = e.str
					n.str = e.str
					t2 = append(t2, n)
					continue
				}
			}
			if i == len(t)-1 {
				err = errors.New("predicate without term")
				return t, err
			}
			if !t[i+1].isTerm() {
				err = errors.New("predicate without term, " + t[i+1].str)
				return t, err
			}
			n.tokenType = tAtomicSentence
			n.predicate = e.str
			n.str = e.str
			for j := i + 1; j < len(t); j++ {
				if !t[j].isTerm() {
					break
				}
				n.term = append(n.term, t[j].str)
				i++
			}
			if isGreekFormulaVar(e.str) {
				n.str = n.str + "("
			}
			for _, e := range n.term {
				n.str = n.str + e
			}

			if isGreekFormulaVar(e.str) {
				n.str = n.str + ")"
			}
			t2 = append(t2, n)
			continue
		}
		if e.isConnective() {
			t2 = append(t2, e)
			continue
		}
		if e.isOpenb() || e.isCloseb() {
			t2 = append(t2, e)
			continue
		}
		err = errors.New("something wrong: " + e.str)
		return t, err
	}
	return t2, err

}
