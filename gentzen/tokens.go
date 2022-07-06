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
	tPredicate
	tTerm
	tOpenb
	tCloseb
)

const (
	tConnective tokenID = tConj | tDisj | tNeg | tCond | tUni | tEx
)

const (
	tQuantifier tokenID = tUni | tEx
)

const (
	tUnary tokenID = tNeg | tUni | tEx
)

const (
	tBinary tokenID = tConj | tCond | tDisj
)

func (t token) isQuantifier() bool {
	return t.tokenType == tUni || t.tokenType == tEx
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

func (t token) isUnary() bool {

	return t.isNeg() || t.isUni() || t.isEx()
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

//nextToken returns the next token in s and
//returns the remainder string r. Works with
//predicate logic.
func nextToken(s string) (t token, r string) {

	pos := 1

	switch s[:pos] {
	case ldisj:
		t.tokenType = tDisj
	case lconj:
		t.tokenType = tConj
	case lcond:
		t.tokenType = tCond
	case lneg:
		t.tokenType = tNeg
	case luni:
		t.tokenType = tUni
	case lex:
		t.tokenType = tEx
	default:
		if isLowerCase(s[:pos]) {
			t.tokenType = tTerm
			if len(s) > 1+pos {
				if s[pos:pos+1] == `_` {
					pos = pos + 2
				}
			}
		} else {
			t.tokenType = tPredicate
			if len(s) > 1+pos {
				if s[pos:pos+1] == `\` {
					pos = pos + 2
				}
				if len(s) > 1+pos {
					if s[pos:pos+1] == `_` {
						pos = pos + 2
					}
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
	t.str = s[:pos]
	r = s[pos:]

	return t, r
}

func cleanString(s string) string {

	s = strings.ReplaceAll(s, " ", "")
	s = strings.ReplaceAll(s, "\t", "")
	s = strings.ReplaceAll(s, "(", "")
	s = strings.ReplaceAll(s, ")", "")

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

func tokenizePLround2(t tokenStr) (tokenStr, error) {
	var t2 tokenStr
	var e token
	var err error

	for i := 0; i < len(t); i++ {
		e = t[i]

		if e.isQuantifier() {
			if i > len(t)-2 {
				err = errors.New("quantifier without variable 1")
				return t, err
			}
			if !t[i+1].isTerm() {
				err = errors.New("quantifier without variable 2")
				return t, err
			}
			var n token
			n.tokenType = e.tokenType
			n.str = e.str
			n.variable = t[i+1].str
			t2 = append(t2, n)
			i++
			continue
		}
		if e.isPredicate() {
			if i == len(t)-1 {
				err = errors.New("predicate without term")
				return t, err
			}
			if !t[i+1].isTerm() {
				err = errors.New("predicate without term")
				return t, err
			}
			var n token
			n.tokenType = tAtomicSentence
			n.predicate = e.str
			n.str = e.str
			for j := i + 1; j < len(t); j++ {
				if !t[j].isTerm() {
					break
				}
				n.term = append(n.term, t[j].str)
				n.str = n.str + t[j].str
				i++
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

		err = errors.New("something wrong")
		return t, err
	}

	return t2, err

}
